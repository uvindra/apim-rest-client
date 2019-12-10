package cmd

import (
	"apim-rest-client/comm"
	"apim-rest-client/constants"
	"apim-rest-client/dcr"
	"apim-rest-client/persist"
	"apim-rest-client/token"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type APIOptions struct {
	API         string
	Method      string
	Resource    string
	Headers     *FlagMap
	QueryParams *FlagMap
	FormData    *FlagMap
	Body        string
	IsVerbose   bool
}

type BasePaths struct {
	PublisherAPI string
	StoreAPI     string
	AdminAPI     string
}

func RefreshExistingTokens(confJSON *persist.Config, isVerbose bool) {
	if isVerbose {
		fmt.Println("Credentials already exist")
	}

	credentials := persist.ReadAppCredentials()

	tokenResp, error := token.RefreshToken(confJSON.TokenURL, credentials.ClientID, credentials.ClientSecret,
		credentials.RefreshToken, confJSON.Scope, isVerbose)

	if error != nil {
		fmt.Printf("Error returned in when refreshing token. error : %s, error_description : %s\n",
			error.ErrorType, error.ErrorDescription)

		if error.ErrorType == "invalid_client" {
			fmt.Println("\nRegistered client does not exist, please execute `arc clear` and then rerun the desired command.")
			os.Exit(1)
		}

		tokenResp, error = token.RequestToken_PasswordGrant(confJSON.TokenURL, credentials.ClientID,
			credentials.ClientSecret, confJSON.UserName, confJSON.Password, confJSON.Scope, isVerbose)
	}

	if error == nil {
		// Store new access token and refresh token
		credentials.AccessToken = tokenResp.AccessToken
		credentials.RefreshToken = tokenResp.RefreshToken

		persist.SaveAppCredentials(&credentials)
	} else {
		fmt.Printf("Error returned in when requesting new token. error : %s, error_description : %s\n",
			error.ErrorType, error.ErrorDescription)
	}
}

func RegisterClient(confJSON *persist.Config, isVerbose bool) persist.OAuthCredentials {
	if isVerbose {
		fmt.Println("Credentials do not exist")
	}

	var dcrRequest dcr.DCRRequest
	dcr.SetDCRParameters(&dcrRequest, confJSON.UserName)

	dcrResp := dcr.Register(confJSON.DcrURL, confJSON.UserName, confJSON.Password, dcrRequest, isVerbose)

	var credentials persist.OAuthCredentials
	credentials.ClientID = dcrResp.ClientId
	credentials.ClientSecret = dcrResp.ClientSecret

	return credentials
}

func GetTokens(credentials *persist.OAuthCredentials, confJSON *persist.Config, isVerbose bool) {
	tokenResp, error := token.RequestToken_PasswordGrant(confJSON.TokenURL, credentials.ClientID,
		credentials.ClientSecret, confJSON.UserName, confJSON.Password, confJSON.Scope, isVerbose)

	if error == nil {
		credentials.AccessToken = tokenResp.AccessToken
		credentials.RefreshToken = tokenResp.RefreshToken

		persist.SaveAppCredentials(credentials)
	} else {
		fmt.Printf("Error returned in when requesting new token. error : %s, error_description : %s\n",
			error.ErrorType, error.ErrorDescription)
	}
}

func InvokeAPI(apiOptions *APIOptions, basePaths *BasePaths, token string) {
	var basePath string

	switch apiOptions.API {
	case "publisher":
		basePath = basePaths.PublisherAPI
	case "store":
		basePath = basePaths.StoreAPI
	case "admin":
		basePath = basePaths.AdminAPI
	default:
		fmt.Println("Unsupported API base path")
		return
	}

	fullPath := basePath + apiOptions.Resource

	var req *http.Request
	var body *bytes.Buffer
	var contentType string

	switch apiOptions.Method {
	case "GET":
		req = comm.CreateGet(fullPath)
	case "DELETE":
		req = comm.CreateDelete(fullPath)
	case "POST":
		body, contentType = getBodyContent(apiOptions)

		if body == nil {
			req = comm.CreatePostEmptyBody(fullPath)
		} else {
			req = comm.CreatePost(fullPath, body)
		}
	case "PUT":
		body, contentType = getBodyContent(apiOptions)

		if body == nil {
			req = comm.CreatePutEmptyBody(fullPath)
		} else {
			req = comm.CreatePut(fullPath, body)
		}
	}

	comm.SetDefaultRestAPIHeaders(token, contentType, req)

	headers := http.Header{}

	for k, v := range *apiOptions.Headers {
		headers.Add(k, v)
	}

	comm.AddHeaders(&headers, req)

	values := req.URL.Query()
	for k, v := range *apiOptions.QueryParams {
		values.Add(k, v)
	}

	comm.AddQueryParams(&values, req)

	if apiOptions.IsVerbose {
		comm.PrintRequest(constants.REST_API_REQUEST_LOG_STRING, req)
	}

	resp := comm.SendHTTPRequest(req)

	defer resp.Body.Close()

	comm.PrintResponse(constants.REST_API_RESPONSE_LOG_STRING, resp)
}

func getBodyContent(apiOptions *APIOptions) (body *bytes.Buffer, contentType string) {
	if apiOptions.Body != constants.UNDEFINED_STRING {
		return bytes.NewBuffer(readData(apiOptions.Body)), constants.UNDEFINED_STRING
	}

	data := map[string]string{}
	for k, v := range *apiOptions.FormData {
		data[k] = v
	}

	if len(*apiOptions.FormData) != 0 {
		return comm.CreateMultipartFormData(&data)
	}

	return nil, constants.UNDEFINED_STRING
}

func readData(data string) []byte {
	if data[0] == '@' {
		content, err := ioutil.ReadFile(data[1:])

		if err != nil {
			panic(err)
		}

		return content
	}

	return []byte(data)
}
