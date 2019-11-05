package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/harrietty/conju_api/lambdas/conjugations/handler"
	"github.com/harrietty/conju_api/verbsbucket"
)

func main() {
	stage, exists := os.LookupEnv("STAGE")
	if !exists {
		stage = "dev"
	}

	vb := verbsbucket.New("conjugator-verb-data")

	h := handler.New(stage, vb.GetFile)
	lambda.Start(h.HandleRequest)
}
