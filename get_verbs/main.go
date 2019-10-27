package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received request: ", request.HTTPMethod, request.Path, request.QueryStringParameters)

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
            fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
        default:
            fmt.Println(aerr.Error())
        }
    } else {
        fmt.Println(err.Error())
		}
	}
		
	fmt.Println(result)

	return events.APIGatewayProxyResponse{Body: "You requested " + language, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
