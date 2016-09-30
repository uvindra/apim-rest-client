package persist

import (
	"io/ioutil"
	"encoding/json"
	"os"
)

const OAUTH_CREDENTIAL_FILE = "auth_info.json";

type OAuthCredentials struct {
	ClientID string
	ClientSecret string
	AccessToken string
	RefreshToken string
}

func SaveAppCredentials(credentials *OAuthCredentials) {
	content, _ := json.Marshal(*credentials)
	err := ioutil.WriteFile(OAUTH_CREDENTIAL_FILE, content, 0644)

	if err != nil {
		panic(err)
	}
}

func IsAppCredentialsExist() bool {
	_, err := os.Stat(OAUTH_CREDENTIAL_FILE)
	return !os.IsNotExist(err)
}

func ReadAppCredentials() OAuthCredentials {
	b, err := ioutil.ReadFile(OAUTH_CREDENTIAL_FILE)

	if err != nil {
		panic(err)
	}

	var credentials OAuthCredentials

	err = json.Unmarshal(b, &credentials)

	if err != nil {
		panic(err)
	}

	return credentials
}

