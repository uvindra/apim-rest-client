# **A**PIM **R**EST **C**lient

An intuitive tool for interacting with WSO2 API Manager via its REST APIs.

## Features

1. Automatically handles dynamic client regsitration, token generation & token refreshing without user intervention.

2. Defines generic methods for invoking WSO2 API Managers REST interface.


## Setup

_**Prerequisite**_ : Install [Go](https://golang.org/) 

1. Clone the repository

2. Run the command `go build arc.go` to build the `arc` executable.


## Configure

1. Execute `arc init` to generate the *config.json* file

2. The *config.json* file will contain the following default values

     ```json
     {
       "dcrURL": "https://localhost:9443/client-registration/{version}/register",
       "publisherAPI": "https://localhost:9443/api/am/publisher/{version}",
       "storeAPI": "https://localhost:9443/api/am/store/{version}",
       "adminAPI": "https://localhost:9443/api/am/admin/{version}",
       "userName": "admin",
       "password": "admin",
       "tokenURL": "https://localhost:8243/token",
       "scope": "apim:api_view apim:api_create apim:api_publish apim:subscribe 
     }
     ```
    Please consult the product REST API documentation of the respective WSO2 API Manager version you are using and replace the {version} tag with the releavent version string.

    For example for API Manager 3.0.0 the releavent versions will be as follows,

     ```json
     {
       "dcrURL": "https://localhost:9443/client-registration/v0.15/register",
       "publisherAPI": "https://localhost:9443/api/am/publisher/v1.0",
       "storeAPI": "https://localhost:9443/api/am/store/v1.0",
       "adminAPI": "https://localhost:9443/api/am/admin/v0.15",
       "userName": "admin",
       "password": "admin",
       "tokenURL": "https://localhost:8243/token",
       "scope": "apim:api_view apim:api_create apim:api_publish apim:subscribe 
     }
     ```


3. Execute `arc call` with the following arguments to do a REST call to WSO2 API Manager
    
     ```
     arc call --api ("publisher"|"store"|"admin") 
              --method ("GET"|"POST"|"PUT"|"DELETE")           
              --resource "<resource-path>" 
              [--header "<header-name>:<header-value>"]
              [--query-param "<param-name>:<param-value>"] 
              [--form-data "<key>:<value>|@<file-path>"]
              [--body "<value>|@<file-path>"]
     ```

     ### --api ("publisher"|"store"|"admin")

     Specifies the REST API being invoked,

          

     

     Store Function | Usage 
    ---------- | -------
    To get list of existing applications | `./arc call --api "store" --method "GET" --resource "/applications"`
    To get details of an aplication | `./arc call --api "store" --method "GET" --resource "/applications/{uuid}"`
    To delete an application | `./arc call --api "store" --method "DELETE" --resource "/applications/{uuid}"`
    To get key details of an application | `./arc call --api "store" --method "GET" --resource "/applications/{uuid}/keys/PRODUCTION"`
    To update grant types & callback URL of an application | `./arc call --api "store" --method "PUT" --resource "/applications/{uuid}/keys/PRODUCTION" --body "./data.json"`
    To update an application | `./arc call --api "store" --method "PUT" --resource "/applications/{uuid}" --body "./data.json"`
    To generate keys of an application | `./arc call --api "store" --method "POST" --resource "/applications/generate-keys" --query-param "applicationId:{uuid}" --body "./data.json"`
    Importing an api via openapi definition | `./arc call --api "publisher" --method "POST" --resource "/apis/import-openapi" --form-data "file:@/path/api-definition.yaml" --form-data "additionalProperties:@/path/additional-properties.json"`
