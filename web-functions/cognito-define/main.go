package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func handler(event events.CognitoEventUserPoolsDefineAuthChallenge) (events.CognitoEventUserPoolsDefineAuthChallenge, error) {
	l := len(event.Request.Session)
	if l >= 3 && event.Request.Session[l-1].ChallengeResult == false {
		log.Println("authentication failed")
		// The user provided a wrong answer 3 times; fail auth
		event.Response.IssueTokens = false
		event.Response.FailAuthentication = true
	} else if l > 0 && event.Request.Session[l-1].ChallengeResult == true {
		log.Println("authentication succeeded")
		// The user provided the right answer; succeed auth
		event.Response.IssueTokens = true
		event.Response.FailAuthentication = false
	} else {
		log.Println("code was incorrect")
		// The user did not provide a correct answer yet; present challenge
		event.Response.IssueTokens = false
		event.Response.FailAuthentication = false
		event.Response.ChallengeName = "CUSTOM_CHALLENGE"
	}

	return event, nil
}

func main() {
	lambda.Start(handler)
}
