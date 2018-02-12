package cmd

import(
	"bytes"
	"io/ioutil"
	"fmt"
	"net/http"
	"os"
	"apim-rest-client/comm"
	"apim-rest-client/constants"
)

type APIOptions struct {
	API string
	Method string
	Resource string
	QueryParams *FlagMap
	Body string
}

type BasePaths struct {
	PublisherAPI string
	StoreAPI string
	AdminAPI string
}

func CreateData(dataTemplate string) {

	switch dataTemplate {
	case "product":
		createProductDataFile()
	case "api":
		createSwaggerDataFile()
		createApiDataFile()
	case constants.UNDEFINED_STRING: // Flag not specified
		return
	default:
		fmt.Println("Unsupported data template specified")
	}

	// Since this is a data setup function exit to allow users to customize data values
	os.Exit(0);
}

func InvokeAPI(apiOptions *APIOptions, basePaths *BasePaths, token string) {

	switch apiOptions.Resource {
	case "publisher:view:apis":
		publisherGetAPIs(apiOptions, basePaths.PublisherAPI + "/apis", token)
	case "publisher:create:api":
		publisherCreateAPI(apiOptions, basePaths.PublisherAPI + "/apis", token)
	case "publisher:view:products":
		publisherGetProducts(apiOptions, basePaths.PublisherAPI + "/products", token)
	case "publisher:create:product":
		publisherCreateProduct(apiOptions, basePaths.PublisherAPI + "/products", token)
	case constants.UNDEFINED_STRING:  // Flag not specified
		return
	default:
		genericInvoke(apiOptions, basePaths, token)
	}

}

func genericInvoke(apiOptions *APIOptions, basePaths *BasePaths, token string) {
	var basePath string

	switch apiOptions.API {
	case "publisher":
		basePath = basePaths.PublisherAPI
	case "store":
		basePath = basePaths.StoreAPI
	case "admin":
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
		req = comm.CreatePost(fullPath, getBodyContent(apiOptions))
	case "PUT":
		req = comm.CreatePut(fullPath, getBodyContent(apiOptions))
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
