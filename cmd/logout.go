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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/kochie/guardian/config"
	"github.com/kochie/guardian/utils"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("logout called")
		credentials, err := utils.RetrieveAuth()
		if err != nil {panic(err)}

		global, err := cmd.Flags().GetBool("global")
		if err != nil {panic(err)}

		cid := cognitoidentityprovider.New(config.GetDefaultSession())
		cog := utils.NewCognitoHelper(cid)
		if cog.IsTokenExpired() { cog.RefreshToken() }

		forget, err := cmd.Flags().GetBool("forget")
		if err != nil {panic(err)}

		if forget {
			_, err := cid.ForgetDevice(&cognitoidentityprovider.ForgetDeviceInput{
				AccessToken: aws.String(credentials.AccessToken),
				DeviceKey:   aws.String(credentials.DeviceKey),
			})
			if err != nil {panic(err)}
		}

		if global {
			cid.GlobalSignOutRequest(&cognitoidentityprovider.GlobalSignOutInput{
				AccessToken: aws.String(credentials.AccessToken),
			})
		}

		if err := utils.RemoveAuth(); err != nil {panic(err)}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	logoutCmd.PersistentFlags().Bool("forget", true, "forget device")
	logoutCmd.PersistentFlags().Bool("global", false, "log out of all devices")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
