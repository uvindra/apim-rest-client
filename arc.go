package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"apim_rest_client/dcr"
	"apim_rest_client/token"
	"apim_rest_client/persist"
	"apim_rest_client/cmd"
	"apim_rest_client/constants"
	)

const CONFIG_FILE_PATH = "config" + string(filepath.Separator) + "config.json";

type Config struct {
	DcrURL string
	PublisherAPI string
	StoreAPI string
	UserName string
	Password string
	TokenURL string
	Scope string
}


func readConfig() Config {
	b, err := ioutil.ReadFile(CONFIG_FILE_PATH)

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


func refreshExistingTokens(confJSON *Config) {
	credentials := persist.ReadAppCredentials()

	tokenResp, error := token.RefreshToken(confJSON.TokenURL, credentials.ClientID, credentials.ClientSecret,
		credentials.RefreshToken, confJSON.Scope)

	if error != nil {
		fmt.Printf("Error returned in when refreshing token. error : %s, error_description : %s\n",
			error.ErrorType, error.ErrorDescription)

		tokenResp, error = token.RequestToken_PasswordGrant(confJSON.TokenURL, credentials.ClientID,
			credentials.ClientSecret, confJSON.UserName, confJSON.Password, confJSON.Scope)
	}

	if (error == nil) {
		// Store new access token and refresh token
		credentials.AccessToken = tokenResp.AccessToken
		credentials.RefreshToken = tokenResp.RefreshToken

		persist.SaveAppCredentials(&credentials)
	} else {
		fmt.Printf("Error returned in when requesting new token. error : %s, error_description : %s\n",
			error.ErrorType, error.ErrorDescription)
	}
}

func registerClient(confJSON *Config) persist.OAuthCredentials {
	var dcrRequest dcr.DCRRequest
	dcr.SetDCRParameters(&dcrRequest)

	dcrResp := dcr.Register(confJSON.DcrURL, confJSON.UserName, confJSON.Password, dcrRequest)

	var credentials persist.OAuthCredentials
	credentials.ClientID = dcrResp.ClientId
	credentials.ClientSecret = dcrResp.ClientSecret

	return credentials
}

func getTokens(credentials *persist.OAuthCredentials, confJSON *Config) {
	tokenResp, error := token.RequestToken_PasswordGrant(confJSON.TokenURL, credentials.ClientID,
		credentials.ClientSecret, confJSON.UserName, confJSON.Password, confJSON.Scope)

	if (error == nil) {
		credentials.AccessToken = tokenResp.AccessToken
		credentials.RefreshToken = tokenResp.RefreshToken

		persist.SaveAppCredentials(credentials)
	} else {
		fmt.Printf("Error returned in when requesting new token. error : %s, error_description : %s\n",
			error.ErrorType, error.ErrorDescription)
	}
}


func main() {
	flags := cmd.Flags{}

	// Customize flag usage output to prevent default values being printed
	flag.Usage = func() {
		//fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		//flag.PrintDefaults()
	}


	flag.StringVar(&flags.Resource, "resource", constants.UNDEFINED_STRING, "Desired resource in the format " +
			"<location of resource>:<resource name> (example: apis resource in publisher = publisher:apis)")

	flag.StringVar(&flags.Query, "query", constants.UNDEFINED_STRING, "Search query")

	flag.IntVar(&flags.Limit, "limit", constants.UNDEFINED_INT, "Maximum size of resource array to return")

	flag.IntVar(&flags.Offset, "offset", constants.UNDEFINED_INT, "Starting point within the complete list of items qualified")

	flag.Parse()

	confJSON := readConfig()

	if persist.IsAppCredentialsExist() {
		fmt.Println("Credentials already exist")
		refreshExistingTokens(&confJSON)
	} else {
		fmt.Println("Credentials do not exist")
		credentials := registerClient(&confJSON)

		getTokens(&credentials, &confJSON)
	}

	credentials := persist.ReadAppCredentials()

	cmd.ProcessArgs(&flags, confJSON.PublisherAPI, confJSON.StoreAPI, credentials.AccessToken)
}
