package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(event events.CognitoEventUserPoolsVerifyAuthChallenge) (events.CognitoEventUserPoolsVerifyAuthChallenge, error) {
	expectedAnswer := event.Request.PrivateChallengeParameters["secretLoginCode"]
	answer := event.Request.ChallengeAnswer.(string)

	if answer == expectedAnswer {
		event.Response.AnswerCorrect = true
	} else {
		event.Response.AnswerCorrect = false
	}
	return event, nil
}

func main() {
	lambda.Start(handler)
}
