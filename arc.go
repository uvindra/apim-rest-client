package main

import (
	"fmt"
	"flag"
	"os"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"apim-rest-client/dcr"
	"apim-rest-client/token"
	"apim-rest-client/persist"
	"apim-rest-client/cmd"
	"apim-rest-client/constants"
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
	dcr.SetDCRParameters(&dcrRequest, confJSON.UserName)

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
	apiOptions := cmd.APIOptions{}
	queryParams := cmd.FlagMap{}

	// Customize flag usage output to prevent default values being printed
	//flag.Usage = func() {
	//	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	//	flag.PrintDefaults()
	//}

	initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	callCommand := flag.NewFlagSet("call", flag.ExitOnError)
	createDataCommand := flag.NewFlagSet("create-data", flag.ExitOnError)

	productVersion := initCommand.String("version", constants.UNDEFINED_STRING, "APIM product version being used(example: 2.1.0)")

	callCommand.StringVar(&apiOptions.API, "api", constants.UNDEFINED_STRING, "REST API to invoked(example: publisher|store|admin)")
	callCommand.StringVar(&apiOptions.Method, "method", constants.UNDEFINED_STRING, "HTTP Method(example: GET)")
	callCommand.StringVar(&apiOptions.Resource, "resource", constants.UNDEFINED_STRING, "Desired resource path(example: /apis)")
	callCommand.StringVar(&apiOptions.Body, "body", constants.UNDEFINED_STRING, "File path to content of HTTP body(example: ./body.json)")

	callCommand.Var(&queryParams, "query-param", "")

	dataTemplate := createDataCommand.String("create-data", constants.UNDEFINED_STRING, "Create specified data template to be sent in request")

	//flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Mandatory arguments missing.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		initCommand.Parse(os.Args[2:])
	case "call":
		callCommand.Parse(os.Args[2:])
	case "create-data":
		createDataCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if initCommand.Parsed() {
		persist.GenerateConfig(productVersion)
	}

	if createDataCommand.Parsed() {
		cmd.CreateData(*dataTemplate)
	}

	if callCommand.Parsed() {
		apiOptions.QueryParams = &queryParams

		fmt.Printf("apiOptions: %+v\n", apiOptions)
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


		basePaths := cmd.BasePaths{}
		basePaths.PublisherAPI = confJSON.PublisherAPI
		basePaths.StoreAPI = confJSON.StoreAPI

		cmd.InvokeAPI(&apiOptions, &basePaths, credentials.AccessToken)
	}
}
