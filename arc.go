package main

import (
	"apim-rest-client/cmd"
	"apim-rest-client/constants"
	"apim-rest-client/persist"
	"flag"
	"fmt"
	"os"
)

func main() {
	apiOptions := cmd.APIOptions{}
	headers := cmd.FlagMap{}
	queryParams := cmd.FlagMap{}
	formData := cmd.FlagMap{}

	// Customize flag usage output to prevent default values being printed
	flag.Usage = func() {
		usageText := "APIM REST Client.\n" +
			"\n" +
			"Usage: \n" +
			"  arc init \n" +
			"  arc clear \n" +
			"  arc call --api (\"publisher\"|\"store\"|\"admin\") --method (\"GET\"|\"POST\"|\"PUT\"|\"DELETE\") --resource \"<resource-path>\" \n" +
			"           [--query-param \"<param-name>:<param-value>\"] [--body \"<file-path>\"] [-v|--verbose]\n" +
			"\n" +
			"Options: \n" +
			"  -h --help     Show this screen. \n"

		fmt.Println(usageText)
	}

	initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	clearCommand := flag.NewFlagSet("clear", flag.ExitOnError)
	callCommand := flag.NewFlagSet("call", flag.ExitOnError)

	callCommand.StringVar(&apiOptions.API, "api", constants.UNDEFINED_STRING, "REST API to invoked(example: publisher|store|admin)")
	callCommand.StringVar(&apiOptions.Method, "method", constants.UNDEFINED_STRING, "HTTP Method(example: GET)")
	callCommand.StringVar(&apiOptions.Resource, "resource", constants.UNDEFINED_STRING, "Desired resource path(example: /apis)")
	callCommand.StringVar(&apiOptions.Body, "body", constants.UNDEFINED_STRING, "File path to content of HTTP body(example: ./body.json)")
	callCommand.BoolVar(&apiOptions.IsVerbose, "v", false, "Outputs details of the communications done from/to the client")
	callCommand.BoolVar(&apiOptions.IsVerbose, "verbose", false, "Outputs details of the communications done from/to the client")

	callCommand.Var(&headers, "header", "")
	callCommand.Var(&queryParams, "query-param", "")
	callCommand.Var(&formData, "form-data", "")

	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		initCommand.Parse(os.Args[2:])
	case "clear":
		clearCommand.Parse(os.Args[2:])
	case "call":
		callCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("Unspported argument '%s' provided\n", os.Args[1])
		fmt.Println()
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	if initCommand.Parsed() {
		persist.GenerateConfig()
		os.Exit(0)
	}

	if clearCommand.Parsed() {
		persist.DeleteAppCredentials()
		os.Exit(0)
	}

	if callCommand.Parsed() {
		apiOptions.Headers = &headers
		apiOptions.QueryParams = &queryParams
		apiOptions.FormData = &formData

		cmd.Validate(&apiOptions)

		confJSON := persist.ReadConfig()

		if persist.IsAppCredentialsExist() {
			cmd.RefreshExistingTokens(&confJSON, apiOptions.IsVerbose)
		} else {
			credentials := cmd.RegisterClient(&confJSON, apiOptions.IsVerbose)

			cmd.GetTokens(&credentials, &confJSON, apiOptions.IsVerbose)
		}

		credentials := persist.ReadAppCredentials()

		basePaths := cmd.BasePaths{}
		basePaths.PublisherAPI = confJSON.PublisherAPI
		basePaths.StoreAPI = confJSON.StoreAPI
		basePaths.AdminAPI = confJSON.AdminAPI

		cmd.InvokeAPI(&apiOptions, &basePaths, credentials.AccessToken)
	}

}
