package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/harrietty/conju_api/lambdas/infinitives/handler"
)

func main() {
	stage, exists := os.LookupEnv("STAGE")
	if !exists {
		stage = "dev"
	}
	h := handler.New(stage)
	lambda.Start(h.HandleRequest)
}
