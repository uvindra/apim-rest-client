package persist

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const OAUTH_CREDENTIAL_FILE = ".auth_info.json"
const CONFIG_FILE = "config.json"

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
	AdminAPI     string `json:"adminAPI"`
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

func DeleteAppCredentials() {
	err := os.Remove(OAUTH_CREDENTIAL_FILE)

	if err != nil {
		panic(err)
	}
}

func GenerateConfig() {
	config := Config{}
	config.DcrURL = "http://localhost:9763/client-registration/{version}/register"
	config.PublisherAPI = "https://localhost:9443/api/am/publisher/{version}"
	config.StoreAPI = "https://localhost:9443/api/am/store/{version}"
	config.AdminAPI = "https://localhost:9443/api/am/admin/{version}"
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

func IsConfigExists() bool {
	_, err := os.Stat(CONFIG_FILE)
	return !os.IsNotExist(err)
}

func ReadConfig() Config {
	b, err := ioutil.ReadFile(CONFIG_FILE)

	if err != nil {
		panic(err)
	}

	var config Config

	err = json.Unmarshal(b, &config)

	if err != nil {
		panic(err)
	}

	return config
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
