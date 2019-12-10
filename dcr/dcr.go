package dcr

import (
	"apim-rest-client/comm"
	"apim-rest-client/constants"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type DCRRequest struct {
	CallbackURL         string `json:"callbackUrl"`
	ClientName          string `json:"clientName"`
	TokenScope          string `json:"tokenScope"`
	Owner               string `json:"owner"`
	SupportedGrantTypes string `json:"grantType"`
	IsSaaSApp           bool   `json:"saasApp"`
}

type JsonString struct {
	UserName     string `json:"username"`
	ClientName   string `json:"client_name"`
	RedirectURIs string `json:"redirect_uris"`
	GrantTypes   string `json:"grant_types"`
}

type DCRResponse struct {
	CallbackURL  string      `json:"callBackURL"`
	ClientName   string      `json:"clientName"`
	JsonString   *JsonString `json:"jsonString"`
	ClientId     string      `json:"clientId"`
	ClientSecret string      `json:"clientSecret"`
	IsSaaSApp    bool        `json:"isSaasApplication"`
	Owner        string      `json:"appOwner"`
}

func Register(dcrURL string, userName string, password string, regInfo DCRRequest, isVerbose bool) DCRResponse {
	if isVerbose {
		fmt.Println("Register " + regInfo.ClientName + " with DCR")
	}

	data, err := json.Marshal(regInfo)

	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	req := comm.CreatePost(dcrURL, bytes.NewBuffer(data))

	comm.SetDCRHeaders(userName, password, req)

	if isVerbose {
		comm.PrintRequest(constants.DCR_REQUEST_LOG_STRING, req)
	}

	resp := comm.SendHTTPRequest(req)

	if isVerbose {
		comm.PrintResponse(constants.DCR_RESPONSE_LOG_STRING, resp)
	}

	defer resp.Body.Close()

	var jsonResp DCRResponse
	json.NewDecoder(resp.Body).Decode(&jsonResp)

	if isVerbose {
		fmt.Println()
		fmt.Println(jsonResp)

		fmt.Println()
		fmt.Printf("ClientId : %s\n", jsonResp.ClientId)
		fmt.Printf("ClientSecret : %s\n", jsonResp.ClientSecret)
	}

	return jsonResp
}
