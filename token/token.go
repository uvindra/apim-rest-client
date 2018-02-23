package token

import (
	"apim-rest-client/comm"
	"apim-rest-client/constants"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type TokenResponse struct {
	Scope        string `json:"scope"`
	TokeType     string `json:"token_type"`
	ValidTime    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type Error struct {
	ErrorType        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func invokeTokenAPI(request *http.Request) (TokenResponse, *Error) {
	comm.PrintRequest(constants.TOKEN_API_REQUEST_LOG_STRING, request)

	resp := comm.SendHTTPRequest(request)

	comm.PrintResponse(constants.TOKEN_API_RESPONSE_LOG_STRING, resp)

	defer resp.Body.Close()

	var jsonResp TokenResponse
	var error Error

	switch code := resp.StatusCode; code {
	case 200:
		json.NewDecoder(resp.Body).Decode(&jsonResp)
	case 400:
	case 401:
		json.NewDecoder(resp.Body).Decode(&error)
		return jsonResp, &error
	default:
	}

	return jsonResp, nil
}

func RequestToken_PasswordGrant(tokenURL string, clientID string, clientSecret string,
	userName string, password string, scope string) (TokenResponse, *Error) {
	fmt.Println("Request new token with password grant")

	req := comm.CreatePost(tokenURL, nil)

	comm.SetTokenAPIHeaders(clientID, clientSecret, req)

	values := url.Values{}
	values.Add(constants.GRANT_TYPE_HEADER, constants.PASSWORD_GRANT_TYPE)
	values.Add(constants.USER_NAME_HEADER, userName)
	values.Add(constants.PASSWORD_HEADER, password)
	values.Add(constants.SCOPE_HEADER, scope)

	comm.AddQueryParams(&values, req)

	return invokeTokenAPI(req)
}

func RefreshToken(tokenURL string, clientID string, clientSecret string,
	refreshToken string, scope string) (TokenResponse, *Error) {
	fmt.Println("Refresh token")

	req := comm.CreatePost(tokenURL, nil)

	comm.SetTokenAPIHeaders(clientID, clientSecret, req)

	values := url.Values{}
	values.Add(constants.GRANT_TYPE_HEADER, constants.REFRESH_TOKEN_GRANT_TYPE)
	values.Add(constants.REFRESH_TOKEN_HEADER, refreshToken)
	values.Add(constants.SCOPE_HEADER, scope)
	comm.AddQueryParams(&values, req)

	return invokeTokenAPI(req)
}
