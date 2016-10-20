package cmd

import(
	"fmt"
	"os"
	"apim-rest-client/constants"
)

type APIOptions struct {
	Resource string
	URLParams *FlagMap
	QueryParams *FlagMap
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

func InvokeAPI(apiOptions *APIOptions, publisherAPI string, storeAPI string, token string) {

	switch apiOptions.Resource {
	case "publisher:view:apis":
		publisherGetAPIs(apiOptions, publisherAPI + "/apis", token)
	case "publisher:create:api":
		publisherCreateAPI(apiOptions, publisherAPI + "/apis", token)
	case "publisher:view:products":
		publisherGetProducts(apiOptions, publisherAPI + "/products", token)
	case "publisher:create:product":
		publisherCreateProduct(apiOptions, publisherAPI + "/products", token)
	case constants.UNDEFINED_STRING:  // Flag not specified
		return
	default:
		fmt.Println("Unsupported resource specified")
	}

}
