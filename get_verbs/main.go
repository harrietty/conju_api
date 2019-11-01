package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Received request: ", request.HTTPMethod, request.Path, request.QueryStringParameters)

	language, ok := request.QueryStringParameters["language"]
	if !ok {
		return events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}

	svc := s3.New(session.New())
	input := &s3.GetObjectInput{
		Bucket: aws.String("conjugator-verb-data"),
		Key:    aws.String(language + ".json"),
	}
	result, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				log.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}

		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	s3ObjectBytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Println("Error reading S3 result body: ", err)
	}

	langData := parseLanguageJSON(s3ObjectBytes)
	// langDataJSON, err := json.Marshal(langData)
	inf, err := json.Marshal(extractInfinitives(langData))
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
	}

	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Credentials"] = "true"
	return events.APIGatewayProxyResponse{Body: string(inf), StatusCode: 200, Headers: headers}, nil
}

type infinitives []string

type languageData struct {
	Pronouns []string `json:"pronouns"`
	Verbs    verbs    `json:"verbs"`
}

type verbs struct {
	Basic []verb `json:"basic"`
}

type verb struct {
	Infinitive   string       `json:"infinitive"`
	Translations []string     `json:"translations"`
	Type         []string     `json:"type"`
	Conjugations conjugations `json:"conjugations"`
}

type conjugations struct {
	Present     []string `json:"present"`
	Preterite   []string `json:"preterite"`
	Imperfect   []string `json:"imperfect"`
	Conditional []string `json:"conditional"`
	Future      []string `json:"future"`
}

func parseLanguageJSON(jsonData []byte) languageData {
	var langData languageData
	json.Unmarshal(jsonData, &langData)
	return langData
}

func extractInfinitives(languageData languageData) []string {
	var res []string
	for _, val := range languageData.Verbs.Basic {
		res = append(res, val.Infinitive)
	}
	return res
}

func main() {
	lambda.Start(handler)
}
