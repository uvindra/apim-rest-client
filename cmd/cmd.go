package cmd

import (
	"apim-rest-client/comm"
	"apim-rest-client/constants"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type APIOptions struct {
	API         string
	Method      string
	Resource    string
	Headers     *FlagMap
	QueryParams *FlagMap
	FormData    *FlagMap
	Body        string
}

type BasePaths struct {
	PublisherAPI string
	StoreAPI     string
	AdminAPI     string
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

	comm.PrintRequest(constants.REST_API_REQUEST_LOG_STRING, req)

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
