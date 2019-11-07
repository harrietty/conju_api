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

	langData := parseLanguageJSON(s3ObjectBytes)

	if verbsProvided {
		langData = extractRelevantVerbs(langData, verbsArr)
	}

	jsonString, err := json.Marshal(langData)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Credentials"] = "true"
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(jsonString), Headers: headers}, nil
}

func extractRelevantVerbs(l verbs.LanguageData, v []string) verbs.LanguageData {
	lc := verbs.LanguageData{}
	copier.Copy(&lc, &l)

	var selectedVerbs []verbs.Verb

	// How to do something like this?
	// lc.Verbs.Basic = []verbs.Verb
	for _, elem := range l.Verbs.Basic {
		if contains(v, elem.Infinitive) {
			selectedVerbs = append(selectedVerbs, elem)
		}
	}

	lc.Verbs.Basic = selectedVerbs
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

func parseLanguageJSON(jsonData []byte) verbs.LanguageData {
	var langData verbs.LanguageData
	json.Unmarshal(jsonData, &langData)
	return langData
}

func extractRelevantVerbs(l verbs.LanguageData, v []string) verbs.LanguageData {
	lc := verbs.LanguageData{}
	copier.Copy(&lc, &l)

	var selectedVerbs []verbs.Verb

	// How to do something like this?
	// lc.Verbs.Basic = []verbs.Verb
	for _, elem := range l.Verbs.Basic {
		if contains(v, elem.Infinitive) {
			fmt.Println("contains verb ", elem.Infinitive)
			selectedVerbs = append(selectedVerbs, elem)
			fmt.Println(len(lc.Verbs.Basic))
		}
	}

	lc.Verbs.Basic = selectedVerbs
	return lc
}

func contains(sl []string, s string) bool {
	for _, elem := range sl {
		if elem == s {
			return true
		}
	}
	return false
}
