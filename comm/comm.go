package comm

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"crypto/tls"
	"encoding/base64"
	"io"
	"apim-rest-client/constants"
)

func CreateGet(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	return req
}

func CreatePost(url string, body io.Reader) *http.Request {
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		panic(err)
	}

	return req
}

func CreatePut(url string, body io.Reader) *http.Request {
	req, err := http.NewRequest("PUT", url, body)

	if err != nil {
		panic(err)
	}

	return req
}

func CreateDelete(url string) *http.Request {
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		panic(err)
	}

	return req
}

func AddQueryParams(params *url.Values, request *http.Request) {
	request.URL.RawQuery = params.Encode()
}


func SendHTTPRequest(request *http.Request) *http.Response {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	return resp
}

func PrintRequest(logString string, request *http.Request) {
	dump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> %s >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n", logString);
	fmt.Printf("\n%s\n", dump)
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n");
}

func PrintResponse(logString string, response *http.Response) {
	dump, err := httputil.DumpResponse(response, true)

	if err != nil {
		panic(err)
	}

	//content, _ := json.MarshalIndent(dump, "", "    ")

	fmt.Printf("\n<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< %s <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n", logString);
	fmt.Printf("\n%s\n", dump)
	fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n");
}


func SetRestAPIHeaders(token string, request *http.Request) {
	request.Header.Set(constants.AUTH_HEADER, "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
}


func SetTokenAPIHeaders(clientID string, clientSecret string, request *http.Request) {
	var authHeader = clientID + ":" + clientSecret
	encoded := base64.StdEncoding.EncodeToString([]byte(authHeader))

	request.Header.Set(constants.AUTH_HEADER, "Basic " + encoded)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}


func SetDCRHeaders(userName string, password string, request *http.Request) {
	var authHeader = userName + ":" + password;
	encoded := base64.StdEncoding.EncodeToString([]byte(authHeader))

	request.Header.Set(constants.AUTH_HEADER, "Basic " + encoded)
	request.Header.Set("Content-Type", "application/json")
}
