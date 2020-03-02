package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.CognitoEventUserPoolsVerifyAuthChallengeRequest) (events.CognitoEventUserPoolsVerifyAuthChallengeResponse, error) {
	expectedAnswer := request.PrivateChallengeParameters["secretLoginCode"]

	response := events.CognitoEventUserPoolsVerifyAuthChallengeResponse{}

	if request.ChallengeAnswer == expectedAnswer {
		response.AnswerCorrect = true
	} else {
		response.AnswerCorrect = false
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
