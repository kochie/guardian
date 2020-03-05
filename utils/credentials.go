package utils

import (
	"encoding/json"
	"fmt"
	"github.com/kochie/guardian/definitions"
	"os"
	"path"
)

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		panic(err)
	}
}

func getCredentialLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	credentialLocation := fmt.Sprintf("%s/.config/guardian/credentials.json", homeDir)
	return credentialLocation
}

func StoreAuth(credentials *definitions.Credentials) {
	credFileLocation := getCredentialLocation()
	var file *os.File
	var err error

	if file, err = os.Open(credFileLocation); os.IsNotExist(err) {
		err = os.MkdirAll(path.Dir(credFileLocation), 0700)
		if err != nil {panic(err)}
		file, err = os.Create(credFileLocation)
	}

	if err != nil {
		panic (err)
	}

	defer closeFile(file)
	err = json.NewEncoder(file).Encode(credentials)
	if err != nil {panic(err)}
}

func RetrieveAuth() (credentials *definitions.Credentials, err error) {
	credentials = &definitions.Credentials{}
	credFileLocation := getCredentialLocation()
	var file *os.File

	if file, err = os.Open(credFileLocation); os.IsNotExist(err) {
		return nil, err
	}

	if err != nil {panic(err)}
	defer closeFile(file)

	err = json.NewDecoder(file).Decode(credentials)
	if err != nil {panic(err)}
	return
}

func RemoveAuth() error {
	credFileLocation := getCredentialLocation()
	if _, err := os.Stat(credFileLocation); os.IsNotExist(err) {
		return err
	}

	err := os.Remove(credFileLocation)
	return err
}
