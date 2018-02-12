package persist

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"fmt"
)

const OAUTH_CREDENTIAL_FILE = "auth_info.json";
const CONFIG_FILE = "config/config.json";

type OAuthCredentials struct {
	ClientID string
	ClientSecret string
	AccessToken string
	RefreshToken string
}

type Config struct {
	dcrURL string
	publisherAPI string
	storeAPI string
	userName string
	password string
	tokenURL string
	scope string
}


func SaveAppCredentials(credentials *OAuthCredentials) {
	content, _ := json.MarshalIndent(*credentials, "", "    ")
	err := ioutil.WriteFile(OAUTH_CREDENTIAL_FILE, content, 0644)

	if err != nil {
		panic(err)
	}
}

func GenerateConfig(version *string) {
	var apiVersion string

	switch *version {
	case "2.0.0":
		apiVersion = "v0.10"
	case "2.1.0":
		apiVersion = "v0.11"
	case "2.2.0":
		apiVersion = "v0.11"
	default:
		fmt.Println("Unhanlded version specified")
		return
	}

	config := Config{}
	config.dcrURL = "http://localhost:9763/client-registration/" + apiVersion +"/register"
	config.publisherAPI = "https://localhost:9443/api/am/publisher/" + apiVersion
	config.storeAPI = "https://localhost:9443/api/am/store/" + apiVersion
	config.userName = "admin"
	config.password = "admin"
	config.tokenURL = "https://localhost:8243/token"
	config.scope = "apim:subscribe"

	content, _ := json.MarshalIndent(config, "", "    ")
	err := ioutil.WriteFile(CONFIG_FILE, content, 0644)

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
