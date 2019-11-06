package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/harrietty/conju_api/verbs"
	"github.com/harrietty/conju_api/verbsbucket"
)

// Handler struct
type Handler struct {
	stage        string
	GetVerbsFile verbsbucket.VerbsFileGetter
}

// New creates a new handler
func New(stage string, vfg verbsbucket.VerbsFileGetter) Handler {
	return Handler{
		stage:        stage,
		GetVerbsFile: vfg,
	}
}

// HandleRequest handles an API Request and responds with an array of verb conjugations
func (h Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Received request: ", request.HTTPMethod, request.Path, request.QueryStringParameters)

	language, ok := request.QueryStringParameters["language"]
	if !ok {
		return events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}

	verbs, verbsProvided := request.QueryStringParameters["verbs"]
	var verbsArr []string
	if verbsProvided {
		verbsArr = strings.Split(verbs, ",")
	}

	fmt.Println(verbsArr)

	languageFileName := language + ".json"
	if h.stage == "dev" {
		languageFileName = language + ".dev.json"
	}

	s3ObjectBytes, err := h.GetVerbsFile(languageFileName)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				log.Println(s3.ErrCodeNoSuchKey, aerr.Error())
				return events.APIGatewayProxyResponse{StatusCode: 404}, nil
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}

		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	fmt.Println("Fetches language data")
	langData := parseLanguageJSON(s3ObjectBytes)
	fmt.Println("got langData as JSON")
	jsonString, err := json.Marshal(langData)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}
	fmt.Println("Marshalled JSON")

	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Credentials"] = "true"
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(jsonString), Headers: headers}, nil
}

func parseLanguageJSON(jsonData []byte) verbs.LanguageData {
	var langData verbs.LanguageData
	json.Unmarshal(jsonData, &langData)
	return langData
}
