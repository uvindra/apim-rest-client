**A**PIM **R**EST **C**lient

Currently you can add and get API Products

1. Install [Go](https://golang.org/) if you havent already and clone the repository

2. Run the command `go build arc.go` to build

3. Configure the client using the *config/config.json* file

4. Run the command line executable `arc` as follows,
    
     Function | Usage 
    ---------- | -------
    To generate template *data/product.json* for creating a product | `./arc -create-data=product` 
    To create a new product using the data in *data/product.json* |  `./arc -resource=publisher:create:product` 
    To get list of existing products | `./arc -resource=publisher:view:products` 

     Store Function | Usage 
    ---------- | -------
    To get list of existing applications | `./arc call --api "store" --method "GET" --resource "/applications"`
    To get details of an aplication | `./arc call --api "store" --method "GET" --resource "/applications/{uuid}"`
    To delete an application | `./arc call --api "store" --method "DELETE" --resource "/applications/{uuid}"`
    To get key details of an application | `./arc call --api "store" --method "GET" --resource "/applications/{uuid}/keys/PRODUCTION"`
    To update grant types & callback URL of an application | `./arc call --api "store" --method "PUT" --resource "/applications/{uuid}/keys/PRODUCTION" --body "./data.json"`
    To update an application | `./arc call --api "store" --method "PUT" --resource "/applications/{uuid}" --body "./data.json"`
    To generate keys of an application | `./arc call --api "store" --method "POST" --resource "/applications/generate-keys" --query-param "applicationId:{uuid}" --body "./data.json"`
