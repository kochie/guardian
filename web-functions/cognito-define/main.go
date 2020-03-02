package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.CognitoEventUserPoolsDefineAuthChallengeRequest) (events.CognitoEventUserPoolsDefineAuthChallengeResponse, error) {
	response := events.CognitoEventUserPoolsDefineAuthChallengeResponse{}

	if len(request.Session) >= 3 &&
			request.Session[len(request.Session)-1].ChallengeResult == false {
			// The user provided a wrong answer 3 times; fail auth
			response.IssueTokens = false
			response.FailAuthentication = true
		} else if len(request.Session) > 0 &&
					request.Session[len(request.Session)].ChallengeResult == true {
			// The user provided the right answer; succeed auth
			response.IssueTokens = true
			response.FailAuthentication = false
		} else {
			// The user did not provide a correct answer yet; present challenge
			response.IssueTokens = false
			response.FailAuthentication = false
			response.ChallengeName = "CUSTOM_CHALLENGE"
		}

		return response, nil
}

func main() {
	lambda.Start(handler)
}
