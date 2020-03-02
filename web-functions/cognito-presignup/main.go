package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(events.CognitoEventUserPoolsPreSignupRequest) (events.CognitoEventUserPoolsPreSignupResponse, error) {
	return events.CognitoEventUserPoolsPreSignupResponse{
		AutoConfirmUser: true,
		AutoVerifyEmail: true,
		AutoVerifyPhone: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
