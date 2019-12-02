package comm

import (
	"apim-rest-client/constants"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

func CreatePostEmptyBody(url string) *http.Request {
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		panic(err)
	}

	return req
}

func CreatePutEmptyBody(url string) *http.Request {
	req, err := http.NewRequest("PUT", url, nil)

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

func AddHeaders(headers *http.Header, request *http.Request) {
	isOverrideContentType := false

	for k, v := range *headers {
		// If Content-Type header exists among the supplied headers,
		// this shows intent to set the value of the Content-Type
		// explicitly. Therefore remove the default Content-Type
		// that has been set
		if isOverrideContentType == false && k == "Content-Type" {
			request.Header.Del(k)
			isOverrideContentType = true
		}

		for _, value := range v {
			request.Header.Add(k, value)
		}
	}
}

func AddQueryParams(params *url.Values, request *http.Request) {
	request.URL.RawQuery = params.Encode()
}

func CreateMultipartFormData(formData *map[string]string) (body *bytes.Buffer, contentType string) {
	buffer := new(bytes.Buffer)
	w := multipart.NewWriter(buffer)
	var err error
	var fw io.Writer

	for k, v := range *formData {
		if v[0] == '@' {
			file, err := os.Open(v[1:])

			if err != nil {
				panic(err)
			}

			if fw, err = w.CreateFormFile(k, file.Name()); err != nil {
				panic(err)
			}

			if _, err = io.Copy(fw, file); err != nil {
				panic(err)
			}

		} else {
			if fw, err = w.CreateFormField(k); err != nil {
				panic(err)
			}

			fw.Write([]byte(v))
		}

	}

	w.Close()

	return buffer, w.FormDataContentType()
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

	fmt.Printf("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> %s >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n", logString)
	fmt.Printf("\n%s\n", dump)
	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n")
}

func PrintResponse(logString string, response *http.Response) {
	dump, err := httputil.DumpResponse(response, true)

	if err != nil {
		panic(err)
	}

	//content, _ := json.MarshalIndent(dump, "", "    ")

	fmt.Printf("\n<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< %s <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n", logString)
	fmt.Printf("\n%s\n", dump)
	fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n")
}

func SetDefaultRestAPIHeaders(token string, contentType string, request *http.Request) {
	request.Header.Set(constants.AUTH_HEADER, "Bearer "+token)

	if contentType == constants.UNDEFINED_STRING {
		request.Header.Set("Content-Type", "application/json")
	} else {
		request.Header.Set("Content-Type", contentType)
	}
}

func SetTokenAPIHeaders(clientID string, clientSecret string, request *http.Request) {
	var authHeader = clientID + ":" + clientSecret
	encoded := base64.StdEncoding.EncodeToString([]byte(authHeader))

	request.Header.Set(constants.AUTH_HEADER, "Basic "+encoded)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

func SetDCRHeaders(userName string, password string, request *http.Request) {
	var authHeader = userName + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(authHeader))

	request.Header.Set(constants.AUTH_HEADER, "Basic "+encoded)
	request.Header.Set("Content-Type", "application/json")
}
