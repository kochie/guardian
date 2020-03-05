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
	"fmt"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/kochie/guardian/config"
	"github.com/kochie/guardian/definitions"
	"github.com/kochie/guardian/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login or create an account",
	Long: `login to your guardian account and set up access to resources. 
If the email is not associated with the account then a new account will be created.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := utils.RetrieveAuth()
		if err == nil {
			fmt.Println("You are already logged in!")
			return
		}
		//fmt.Println("login called")
		email, err := cmd.Flags().GetString("email")
		if err != nil || email == "" {
			fmt.Print("email was not provided, please enter it now\nEmail Address: ")
			reader := bufio.NewReader(os.Stdin)
			email, err = reader.ReadString('\n')
			if err != nil {
				log.Panic("email wasn't provided")
			}
			email = strings.Replace(email, "\n", "", -1)
		}
		cid := cognitoidentityprovider.New(config.GetDefaultSession())
		cog := utils.NewCognitoHelper(cid)

		userExists := cog.CheckUserExists(email)
		if !userExists {
			cog.CreateUser(email)
		}

		loginResponse := cog.Login(email)

		verifyResponse := cog.VerifyLogin(email, loginResponse.Session)

		cog.AuthenticateDevice(
			&email,
			verifyResponse.AuthenticationResult.AccessToken,
			verifyResponse.AuthenticationResult.NewDeviceMetadata.DeviceKey,
			verifyResponse.AuthenticationResult.NewDeviceMetadata.DeviceGroupKey,
		)

		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		utils.StoreAuth(&definitions.Credentials{
			DeviceName:   hostname,
			DeviceKey:    *verifyResponse.AuthenticationResult.NewDeviceMetadata.DeviceKey,
			DeviceGroup:  *verifyResponse.AuthenticationResult.NewDeviceMetadata.DeviceGroupKey,
			AccessToken:  *verifyResponse.AuthenticationResult.AccessToken,
			RefreshToken: *verifyResponse.AuthenticationResult.RefreshToken,
			KeyID:        *verifyResponse.AuthenticationResult.IdToken,
		})
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
