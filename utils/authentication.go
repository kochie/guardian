package utils

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/kochie/guardian/config"
	"gopkg.in/square/go-jose.v2/jwt"
	"gopkg.in/square/go-jose.v2"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

type CognitoHelper struct {
	cid *cognitoidentityprovider.CognitoIdentityProvider
}

func NewCognitoHelper(cid *cognitoidentityprovider.CognitoIdentityProvider) *CognitoHelper {
	return &CognitoHelper{cid: cid}
}

func (cog CognitoHelper) AuthenticateDevice(email, accessToken, deviceKey, deviceGroup *string) {
	g := big.NewInt(2)
	N := &big.Int{}
	N.SetString(config.InitN, 16)

	password := generateRandomPassword()
	salt := generateSalt()

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fp := fmt.Sprintf("%s%s:%s", *deviceGroup, *email, password)
	fullPassword := sha256.Sum256([]byte(fp))

	saltAndFull := &big.Int{}
	saltAndFull.SetBytes(append(salt, fullPassword[:]...))

	g.Exp(g, saltAndFull, N)

	passwordVerifier := g

	saltBigInt := &big.Int{}
	saltBigInt.SetBytes(salt)

	_, err = cog.cid.ConfirmDevice(&cognitoidentityprovider.ConfirmDeviceInput{
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
}


func (cog CognitoHelper) CheckUserExists(email string) bool {
	listUsersResponse, err := cog.cid.ListUsers(&cognitoidentityprovider.ListUsersInput{
		AttributesToGet: nil,
		Filter:          aws.String(fmt.Sprintf("username=\"%s\"", email)),
		Limit:           nil,
		PaginationToken: nil,
		UserPoolId:      aws.String(config.UserPoolId),
	})

	if err != nil {
		log.Fatal(err)
	}

	return len(listUsersResponse.Users) > 0
}

func (cog CognitoHelper) CreateUser(email string) {
	password := GeneratePassword()
	_, err := cog.cid.SignUp(&cognitoidentityprovider.SignUpInput{
		Username: aws.String(email),
		Password: aws.String(password),
		ClientId: aws.String(config.ClientId),
		UserAttributes: []*cognitoidentityprovider.AttributeType{{
			Name:  aws.String("email"),
			Value: aws.String(email),
		}},
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (cog CognitoHelper) Login(email string) *cognitoidentityprovider.InitiateAuthOutput {
	loginResponse, err := cog.cid.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AnalyticsMetadata: nil,
		AuthFlow:          aws.String(cognitoidentityprovider.AuthFlowTypeCustomAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(email),
		},
		ClientId: aws.String(config.ClientId),
		ClientMetadata: map[string]*string{
			"email": aws.String(email),
		},
		UserContextData: nil,
	})

	if err != nil {
		log.Fatal(err)
	}

	return loginResponse
}


func (cog CognitoHelper) VerifyLogin(email string, session *string) *cognitoidentityprovider.RespondToAuthChallengeOutput {
	reader := bufio.NewReader(os.Stdin)

	var authChallengeResponse = &cognitoidentityprovider.RespondToAuthChallengeOutput{}

	for authChallengeResponse.AuthenticationResult == nil {
		fmt.Print("Enter Code: ")
		secretLoginCode, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		secretLoginCode = strings.Replace(secretLoginCode, "\n", "", -1)

		authChallengeResponse, err = cog.cid.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
			AnalyticsMetadata: nil,
			ChallengeName:     aws.String("CUSTOM_CHALLENGE"),
			ChallengeResponses: map[string]*string{
				"USERNAME": aws.String(email),
				"ANSWER":   aws.String(secretLoginCode),
			},
			ClientId:        aws.String(config.ClientId),
			ClientMetadata:  nil,
			Session:         session,
			UserContextData: nil,
		})

		if err != nil {
			log.Fatal(err)
		}

	}

	return authChallengeResponse
}

type Key struct {
	Alg string `json:"alg"`
	E string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N string `json:"n"`
	Use string `json:"use"`
}

func (cog CognitoHelper) IsTokenExpired() bool {
	credentials, err := RetrieveAuth()
	if err != nil { panic(err) }

	tok, err := jwt.ParseSigned(credentials.AccessToken)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", config.DefaultRegion, config.UserPoolId))
	if err != nil { panic(err) }

	keys := struct {
		Keys []jose.JSONWebKey `json:"keys"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&keys)
	if err != nil { panic(err) }

	var i int
	for i = 0; i < len(keys.Keys); i++ {
		if keys.Keys[i].KeyID == tok.Headers[0].KeyID {break}
	}

	out := jwt.Claims{}
	if err := tok.Claims(keys.Keys[i].Key, &out); err != nil {
		panic(err)
	}

	return *out.Expiry < *jwt.NewNumericDate(time.Now())
}

func (cog CognitoHelper) RefreshToken() {
	credentials, err := RetrieveAuth()
	if err != nil { panic(err) }

	output, err := cog.cid.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AnalyticsMetadata: nil,
		AuthFlow:          aws.String(cognitoidentityprovider.AuthFlowTypeRefreshTokenAuth),
		AuthParameters:    map[string]*string{
			"REFRESH_TOKEN": aws.String(credentials.RefreshToken),
			"DEVICE_KEY": aws.String(credentials.DeviceKey),
		},
		ClientId:          aws.String(config.ClientId),
		ClientMetadata:    nil,
		UserContextData:   nil,
	})
	if err != nil { panic(err) }

	credentials.AccessToken = *output.AuthenticationResult.AccessToken
	credentials.RefreshToken = *output.AuthenticationResult.RefreshToken

	StoreAuth(credentials)
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
