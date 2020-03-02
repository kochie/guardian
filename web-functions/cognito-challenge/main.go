package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/kochie/guardian/utils"
	"log"
	"os"
	"regexp"
)

var sesClient ses.SES

func handler(request events.CognitoEventUserPoolsCreateAuthChallengeRequest) (events.CognitoEventUserPoolsCreateAuthChallengeResponse, error) {
	var secretLoginCode string
	if len(request.Session) > 0 {

		// This is a new auth session
		// Generate a new secret login code and mail it to the user
		secretLoginCode = utils.GenerateDigitCode(6)
		err := sendEmail(request.UserAttributes["email"], secretLoginCode)
		if err != nil {
			log.Print(err)
			return events.CognitoEventUserPoolsCreateAuthChallengeResponse{}, err
		}

	} else {

		re, err := regexp.Compile(`CODE-(\d*)`)
		if err != nil {
			log.Print(err)
			return events.CognitoEventUserPoolsCreateAuthChallengeResponse{}, err
		}

		// There's an existing session. Don't generate new digits but
		// re-use the code from the current session. This allows the user to
		// make a mistake when keying in the code and to then retry, rather
		// the needing to e-mail the user an all new code again.
		previousChallenge := request.Session[len(request.Session)-1]
		secretLoginCode = re.FindString(previousChallenge.ChallengeMetadata)[4:]
		//secretLoginCode = previousChallenge.ChallengeMetadata!.match(/CODE-(\d*)/)![1];
	}

	// This is sent back to the client app
	publicChallengeParameters := map[string]string {
		"email": request.UserAttributes["email"],
	}

	// Add the secret login code to the private challenge parameters
	// so it can be verified by the "Verify Auth Challenge Response" trigger
	privateChallengeParameters := map[string]string{ "secretLoginCode": secretLoginCode }

	// Add the secret login code to the session so it is available
	// in a next invocation of the "Create Auth Challenge" trigger
	challengeMetadata := fmt.Sprintf("CODE-%s}", secretLoginCode)

	return events.CognitoEventUserPoolsCreateAuthChallengeResponse{
		PrivateChallengeParameters: privateChallengeParameters,
		PublicChallengeParameters: publicChallengeParameters,
		ChallengeMetadata: challengeMetadata,
	}, nil
}

func sendEmail(emailAddress, secretLoginCode string) error {
	fromAddress := os.Getenv("SES_FROM_ADDRESS")
	encoding := "UTF-8"
	subject := "Your secret login code"
	loginCodeHtml := fmt.Sprintf(`<html><body><p>This is your secret login code:</p>
                           <h3>%s</h3></body></html>`, secretLoginCode)
	loginCodeText := fmt.Sprintf(`Your secret login code: %s`, secretLoginCode)
 	params := ses.SendEmailInput{
		Destination: &ses.Destination{ ToAddresses: []*string{&emailAddress} },
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: &encoding,
					Data: &loginCodeHtml,
				},
				Text: &ses.Content{
					Charset: &encoding,
					Data: &loginCodeText,
				},
			},
			Subject: &ses.Content{
				Charset: &encoding,
				Data: &subject,
			},
		},
		Source: &fromAddress,
	}

	_, err := sesClient.SendEmail(&params)
	return err
}

func main() {
	sesClient = ses.SES{}
	lambda.Start(handler)
}