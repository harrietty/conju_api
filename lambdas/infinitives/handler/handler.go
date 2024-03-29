package handler

import (
	"encoding/json"
	"log"

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

// HandleRequest handles an API Request and responds with an array of infinitive verbs
func (h Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Received request: ", request.HTTPMethod, request.Path, request.QueryStringParameters)

	language, ok := request.QueryStringParameters["language"]
	if !ok {
		return events.APIGatewayProxyResponse{StatusCode: 404}, nil
	}

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

	langData := parseLanguageJSON(s3ObjectBytes)
	inf, err := json.Marshal(extractInfinitives(langData))
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
	}

	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Credentials"] = "true"
	return events.APIGatewayProxyResponse{Body: string(inf), StatusCode: 200, Headers: headers}, nil
}

func parseLanguageJSON(jsonData []byte) verbs.LanguageData {
	var langData verbs.LanguageData
	json.Unmarshal(jsonData, &langData)
	return langData
}

func extractInfinitives(languageData verbs.LanguageData) []string {
	var res []string
	for _, val := range languageData.Verbs.Basic {
		res = append(res, val.Infinitive)
	}
	return res
}
