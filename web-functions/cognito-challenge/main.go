package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/kochie/guardian/utils"
	"log"
	"os"
	"regexp"
)

var sesClient *ses.SES

func handler(event events.CognitoEventUserPoolsCreateAuthChallenge) (events.CognitoEventUserPoolsCreateAuthChallenge, error) {
	var secretLoginCode string
	if len(event.Request.Session) == 0 {
		// This is a new auth session
		// Generate a new secret login code and mail it to the user
		secretLoginCode = utils.GenerateDigitCode(6)
		if _, ok := event.Request.UserAttributes["email"]; !ok {
			return event, errors.New("no email defined")
		}
		err := sendEmail(event.Request.UserAttributes["email"], secretLoginCode)
		if err != nil {
			log.Print(err)
			return event, err
		}
	} else {
		re, err := regexp.Compile(`CODE-(\d*)`)
		if err != nil {
			log.Print(err)
			return event, err
		}

		// There's an existing session. Don't generate new digits but
		// re-use the code from the current session. This allows the user to
		// make a mistake when keying in the code and to then retry, rather
		// the needing to e-mail the user an all new code again.
		previousChallenge := event.Request.Session[len(event.Request.Session)-1]
		secretLoginCode = re.FindString(previousChallenge.ChallengeMetadata)[4:]
		//secretLoginCode = previousChallenge.ChallengeMetadata!.match(/CODE-(\d*)/)![1];
	}
	log.Println("Secret Login Code:", secretLoginCode, len(secretLoginCode))

	// This is sent back to the client app
	event.Response.PublicChallengeParameters = map[string]string{
		"email": event.Request.UserAttributes["email"],
	}

	// Add the secret login code to the private challenge parameters
	// so it can be verified by the "Verify Auth Challenge Response" trigger
	event.Response.PrivateChallengeParameters = map[string]string{"secretLoginCode": secretLoginCode}

	// Add the secret login code to the session so it is available
	// in a next invocation of the "Create Auth Challenge" trigger
	event.Response.ChallengeMetadata = fmt.Sprintf("CODE-%s", secretLoginCode)

	return event, nil
}

func sendEmail(emailAddress, secretLoginCode string) error {
	loginCodeHtml := fmt.Sprintf(`<html><body><p>This is your secret login code:</p>
                           <h3>%s</h3></body></html>`, secretLoginCode)
	loginCodeText := fmt.Sprintf(`Your secret login code: %s`, secretLoginCode)
	params := ses.SendEmailInput{
		Destination: &ses.Destination{ToAddresses: []*string{&emailAddress}},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    &loginCodeHtml,
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    &loginCodeText,
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Your secret login code"),
			},
		},
		Source: aws.String(os.Getenv("SES_FROM_ADDRESS")),
	}

	_, err := sesClient.SendEmail(&params)
	return err
}

func main() {
	sesClient = ses.New(session.Must(session.NewSession()))
	lambda.Start(handler)
}
