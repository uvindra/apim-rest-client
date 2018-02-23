package persist

import (
	"apim-rest-client/constants"
	"encoding/json"
	"io/ioutil"
	"os"
)

const OAUTH_CREDENTIAL_FILE = "auth_info.json"
const CONFIG_FILE = "config/config.json"

type OAuthCredentials struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
}

type Config struct {
	DcrURL       string `json:"dcrURL"`
	PublisherAPI string `json:"publisherAPI"`
	StoreAPI     string `json:"storeAPI"`
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	TokenURL     string `json:"tokenURL"`
	Scope        string `json:"scope"`
}

func SaveAppCredentials(credentials *OAuthCredentials) {
	content, _ := json.MarshalIndent(*credentials, "", "    ")
	err := ioutil.WriteFile(OAUTH_CREDENTIAL_FILE, content, 0644)

	if err != nil {
		panic(err)
	}
}

func GenerateConfig(version *string) {
	if *version == constants.UNDEFINED_STRING {
		*version = "v0.11"
	}

	config := Config{}
	config.DcrURL = "http://localhost:9763/client-registration/" + *version + "/register"
	config.PublisherAPI = "https://localhost:9443/api/am/publisher/" + *version
	config.StoreAPI = "https://localhost:9443/api/am/store/" + *version
	config.UserName = "admin"
	config.Password = "admin"
	config.TokenURL = "https://localhost:8243/token"
	config.Scope = "apim:api_view apim:api_create apim:api_publish apim:subscribe"

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
