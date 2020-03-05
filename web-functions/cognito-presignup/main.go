package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	event.Response.AutoConfirmUser = true
	event.Response.AutoVerifyEmail = true
	event.Response.AutoVerifyPhone = false

	return event, nil
}

func main() {
	lambda.Start(handler)
}
