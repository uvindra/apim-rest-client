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
	QueryParams *FlagMap
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

	switch apiOptions.Method {
	case "GET":
		req = comm.CreateGet(fullPath)
	case "DELETE":
		req = comm.CreateDelete(fullPath)
	case "POST":
		body := getBodyContent(apiOptions)

		if body == nil {
			req = comm.CreatePostEmptyBody(fullPath)
		} else {
			req = comm.CreatePost(fullPath, body)
		}
	case "PUT":
		body := getBodyContent(apiOptions)

		if body == nil {
			req = comm.CreatePutEmptyBody(fullPath)
		} else {
			req = comm.CreatePut(fullPath, body)
		}
	}

	comm.SetRestAPIHeaders(token, req)

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

func getBodyContent(apiOptions *APIOptions) (body *bytes.Buffer) {
	if apiOptions.Body != constants.UNDEFINED_STRING {
		content, err := ioutil.ReadFile(apiOptions.Body)

		if err != nil {
			panic(err)
		}

		return bytes.NewBuffer(content)
	}

	return nil
}
