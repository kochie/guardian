package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var cup *cognitoidentityprovider.CognitoIdentityProvider

func handler(event events.CognitoEventUserPoolsPostAuthentication) (events.CognitoEventUserPoolsPostAuthentication, error) {
	if event.Request.UserAttributes["email_verified"] != "true" {
		params := cognitoidentityprovider.AdminUpdateUserAttributesInput{
			UserAttributes: []*cognitoidentityprovider.AttributeType{
				{
					Name:  aws.String("email_verified"),
					Value: aws.String("true"),
				},
			},
			UserPoolId: aws.String(event.UserPoolID),
			Username:   aws.String(event.UserName),
		}

		_, err := cup.AdminUpdateUserAttributes(&params)
		if err != nil {
			return event, err
		}
	}
	return event, nil
}

func main() {
	cup = cognitoidentityprovider.New(session.Must(session.NewSession()))
	lambda.Start(handler)
}
