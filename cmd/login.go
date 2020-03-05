/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/kochie/guardian/utils"
	"github.com/spf13/cobra"
	"log"
	"math/big"
	"os"
	"strings"
)

const (
	ClientId   = "hdc5e6d3gg84k00v6vdsirnd0"
	UserPoolId = "ap-southeast-2_HvB7nliGK"
	InitN      = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64" +
		"ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7" +
		"ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6B" +
		"F12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB31" +
		"43DB5BFCE0FD108E4B82D120A93AD2CAFFFFFFFFFFFFFFFF"
)

func checkUserExists(email string, cid *cognitoidentityprovider.CognitoIdentityProvider) bool {
	listUsersResponse, err := cid.ListUsers(&cognitoidentityprovider.ListUsersInput{
		AttributesToGet: nil,
		Filter:          aws.String(fmt.Sprintf("username=\"%s\"", email)),
		Limit:           nil,
		PaginationToken: nil,
		UserPoolId:      aws.String(UserPoolId),
	})

	if err != nil {
		log.Fatal(err)
	}

	return len(listUsersResponse.Users) > 0
}

func createUser(email string, cid *cognitoidentityprovider.CognitoIdentityProvider) {
	password := utils.GeneratePassword()
	signUpResponse, err := cid.SignUp(&cognitoidentityprovider.SignUpInput{
		Username: aws.String(email),
		Password: aws.String(password),
		ClientId: aws.String(ClientId),
		UserAttributes: []*cognitoidentityprovider.AttributeType{{
			Name:  aws.String("email"),
			Value: aws.String(email),
		}},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(signUpResponse.String())
}

func login(email string, cid *cognitoidentityprovider.CognitoIdentityProvider) *cognitoidentityprovider.InitiateAuthOutput {
	loginResponse, err := cid.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AnalyticsMetadata: nil,
		AuthFlow:          aws.String(cognitoidentityprovider.AuthFlowTypeCustomAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(email),
		},
		ClientId: aws.String(ClientId),
		ClientMetadata: map[string]*string{
			"email": aws.String(email),
		},
		UserContextData: nil,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(loginResponse.String())

	return loginResponse
}

func verifyLogin(email string, session *string, cid *cognitoidentityprovider.CognitoIdentityProvider) *cognitoidentityprovider.RespondToAuthChallengeOutput {
	reader := bufio.NewReader(os.Stdin)

	var authChallengeResponse = &cognitoidentityprovider.RespondToAuthChallengeOutput{}

	for authChallengeResponse.AuthenticationResult == nil {
		fmt.Print("Enter Code: ")
		secretLoginCode, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		secretLoginCode = strings.Replace(secretLoginCode, "\n", "", -1)

		authChallengeResponse, err = cid.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
			AnalyticsMetadata: nil,
			ChallengeName:     aws.String("CUSTOM_CHALLENGE"),
			ChallengeResponses: map[string]*string{
				"USERNAME": aws.String(email),
				"ANSWER":   aws.String(secretLoginCode),
			},
			ClientId:        aws.String(ClientId),
			ClientMetadata:  nil,
			Session:         session,
			UserContextData: nil,
		})

		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println(authChallengeResponse.String())

	return authChallengeResponse
}

func padHex(bigInt *big.Int) string {
	hashStr := bigInt.Text(16)
	if len(hashStr)%2 == 1 {
		hashStr = fmt.Sprintf("0%s", hashStr)
	} else if strings.IndexByte("89ABCDEFabcdef", hashStr[0]) != -1 {
		hashStr = fmt.Sprintf("00%s", hashStr)
	}
	return hashStr
}

func generateRandomPassword() string {
	password := make([]byte, 40)
	_, err := rand.Read(password)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(password)
}

func generateSalt() []byte {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}

	return salt
}

func authenticateDevice(email string, accessToken, deviceKey, deviceGroup *string, cid *cognitoidentityprovider.CognitoIdentityProvider) {
	g := big.NewInt(2)
	N := &big.Int{}
	N.SetString(InitN, 16)

	password := generateRandomPassword()
	salt := generateSalt()

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fp := fmt.Sprintf("%s%s:%s", *deviceGroup, email, password)
	fullPassword := sha256.Sum256([]byte(fp))

	saltAndFull := &big.Int{}
	saltAndFull.SetBytes(append(salt, fullPassword[:]...))

	g.Exp(g, saltAndFull, N)

	passwordVerifier := g

	saltBigInt := &big.Int{}
	saltBigInt.SetBytes(salt)

	confirmDeviceResponse, err := cid.ConfirmDevice(&cognitoidentityprovider.ConfirmDeviceInput{
		AccessToken: accessToken,
		DeviceKey:   deviceKey,
		DeviceName:  aws.String(hostname),
		DeviceSecretVerifierConfig: &cognitoidentityprovider.DeviceSecretVerifierConfigType{
			PasswordVerifier: aws.String(base64.StdEncoding.EncodeToString([]byte(padHex(passwordVerifier)))),
			Salt:             aws.String(base64.StdEncoding.EncodeToString([]byte(padHex(saltBigInt)))),
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(confirmDeviceResponse.String())
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			fmt.Println("email was not provided")
			return
		}
		cid := cognitoidentityprovider.New(session.Must(session.NewSession()))

		userExists := checkUserExists(email, cid)
		if !userExists {
			createUser(email, cid)
		}

		loginResponse := login(email, cid)

		verifyResponse := verifyLogin(email, loginResponse.Session, cid)

		authenticateDevice(
			email,
			verifyResponse.AuthenticationResult.AccessToken,
			verifyResponse.AuthenticationResult.NewDeviceMetadata.DeviceKey,
			verifyResponse.AuthenticationResult.NewDeviceMetadata.DeviceGroupKey,
			cid,
		)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	loginCmd.PersistentFlags().String("email", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
