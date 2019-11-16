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
	"github.com/jinzhu/copier"
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

	fmt.Println("Requested verbs: ", verbsArr)

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

	var langDataStr []byte
	if language == "english" {
		langData := parseEnglishJSON(s3ObjectBytes)
		langDataStr, err = json.Marshal(langData)
		if err != nil {
			log.Println("Error marshalling JSON: ", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}
	} else {
		langData := parseLanguageJSON(s3ObjectBytes)
		if verbsProvided {
			langData = extractRelevantVerbs(langData, verbsArr)
		}
		langDataStr, err = json.Marshal(langData)
		if err != nil {
			log.Println("Error marshalling JSON: ", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}
	}

	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Credentials"] = "true"
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(langDataStr), Headers: headers}, nil
}

func extractRelevantVerbs(l verbs.LanguageData, v []string) verbs.LanguageData {
	lc := verbs.LanguageData{}
	copier.Copy(&lc, &l)

	lc.Verbs.Basic = []verbs.Verb{}

	for _, elem := range l.Verbs.Basic {
		if contains(v, elem.Infinitive) {
			lc.Verbs.Basic = append(lc.Verbs.Basic, elem)
		}
	}

	return lc
}

func contains(v []string, elem string) bool {
	for _, val := range v {
		if val == elem {
			return true
		}
	}
	return false
}

func parseEnglishJSON(jsonData []byte) verbs.EnglishData {
	vbs := verbs.EnglishData{}
	json.Unmarshal(jsonData, &vbs)
	return vbs
}

func parseLanguageJSON(jsonData []byte) verbs.LanguageData {
	var langData verbs.LanguageData
	json.Unmarshal(jsonData, &langData)
	return langData
}
