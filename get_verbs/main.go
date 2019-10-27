package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received request: ", request.HTTPMethod, request.Path, request.QueryStringParameters)

	return events.APIGatewayProxyResponse{Body: "successful, yay", StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
