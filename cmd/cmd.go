package cmd

import(
	//"apim_rest_client/config"
	"fmt"
)

type Flags struct {
	Resource string
	Query string
	Limit int
	Offset int
}


func ProcessArgs(flags *Flags, publisherAPI string, storeAPI string, token string) {
	fmt.Println("resource:", flags)

	switch flags.Resource {
		case "publisher:view:apis":
			getAPIsPublisher(flags, publisherAPI + "/apis", token)
		case "publisher:create:api":
			publisherCreateAPI(flags, publisherAPI + "/apis", token)
		case "publisher:create:product":
			publisherCreateProduct(flags, publisherAPI + "/products", token)
		default:
			fmt.Println("Unsupported resource specified")
	}

}
