package main

import (
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	languageData := parseLanguageJson(s3ObjectBytes)
	log.Println(languageData)

	return events.APIGatewayProxyResponse{Body: "You requested " + language, StatusCode: 200}, nil
}

type LanguageData struct {
	Pronouns []string `json:"pronouns"`
	Verbs Verbs `json:"verbs"`
}

type Verbs struct {
	Basic []Verb `json:"basic"`
}

type Verb struct {
	Infinitive string `json:"infinitive"`
	Translations []string `json:"translations"`
	Type []string `json:"type"`
	Conjugations Conjugations `json:"conjugations"`
}

type Conjugations struct {
	Present []string `json:"present"`
	Preterite []string `json:"preterite"`
	Imperfect []string `json:"imperfect"`
	Conditional []string `json:"conditional"`
	Future []string `json:"future"`
}

func parseLanguageJson(jsonData []byte) LanguageData {
	var langData LanguageData
	json.Unmarshal(jsonData, &langData)
	return langData
}

func main() {
	lambda.Start(Handler)
}
