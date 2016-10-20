package cmd

import (
	"net/url"
	"encoding/json"
	"fmt"
	"bytes"
	"os"
	"io/ioutil"
	"path/filepath"
	"apim-rest-client/comm"
	"apim-rest-client/constants"
	"apim-rest-client/swagger"
)


const SWAGGER_DATA_FILE_PATH = constants.DATA_FOLDER + string(filepath.Separator) + "swagger.json";
const API_DATA_FILE_PATH = constants.DATA_FOLDER + string(filepath.Separator) + "api.json";

type EndpointConfig_t struct {
	ProductionEndpoints *EndpointData_t `json:"production_endpoints"`
	SandboxEndpoints *EndpointData_t `json:"sandbox_endpoints"`
	EndpointType string `json:"endpoint_type"`
}

type EndpointData_t struct {
	URL string `json:"url"`
	Config *string `json:"config"`
}

type Security struct {
	UserName string `json:"username"`
	Type string `json:"type"`
	Password string `json:"password"`
}

type BusinessInfo struct {
	BusinessOwnerEmail string `json:"businessOwnerEmail"`
	TechnicalOwnerEmail string `json:"technicalOwnerEmail"`
	TechnicalOwner string `json:"technicalOwner"`
	BusinessOwner string `json:"businessOwner"`
}

type CorsConfig struct {
	AccessControlAllowOrigins []string `json:"accessControlAllowOrigins"`
	AccessControlAllowHeaders []string `json:"accessControlAllowHeaders"`
	AccessControlAllowMethods []string `json:"accessControlAllowMethods"`
	AccessControlAllowCredentials bool `json:"accessControlAllowCredentials"`
	CorsConfigurationEnabled bool `json:"corsConfigurationEnabled"`
}

type ApiMaxTps struct {
	Sandbox    int32 `json:"sandbox"`
	Production    int32 `json:"production"`
}

type ApiMetaData struct {
	Name    	string `json:"name"`
	Description     string `json:"description"`
	Context     	string `json:"context"`
	Version    	string `json:"version"`
	Provider    	string `json:"provider"`
	ApiDefinition  	string `json:"apiDefinition"`
	WsdlUri    	string `json:"wsdlUri"`
	Status    	string `json:"status"`
	ResponseCaching   string `json:"responseCaching"`
	CacheTimeout    int32 `json:"cacheTimeout"`
	DestinationStatsEnabled    bool `json:"destinationStatsEnabled"`
	IsDefaultVersion    bool `json:"isDefaultVersion"`
	Transport    	[]string `json:"transport"`
	Tags    	[]string `json:"tags"`
	Tiers    	[]string `json:"tiers"`
	MaxTps    	*ApiMaxTps `json:"maxTps"`
	ThumbnailUri    string `json:"thumbnailUri"`
	Visibility 	string `json:"visibility"`
	VisibleRoles 	[]string `json:"visibleRoles"`
	VisibleTenants  []string `json:"visibleTenants"`
	EndpointConfig    string `json:"endpointConfig"`
	EndpointSecurity    *Security `json:"endpointSecurity"`
	GatewayEnvironments    string `json:"gatewayEnvironments"`
	Sequences    			[]string `json:"sequences"`
	SubscriptionAvailability  	string `json:"subscriptionAvailability"`
	SubscriptionAvailableTenants    []string `json:"subscriptionAvailableTenants"`
	BusinessInformation    	*BusinessInfo `json:"businessInformation"`
	CorsConfiguration    	*CorsConfig `json:"corsConfiguration"`
}

func readSwaggerDataFile() string {
	_, noFileErr := os.Stat(SWAGGER_DATA_FILE_PATH)

	if os.IsNotExist(noFileErr) {
		fmt.Printf("\nTemplate data file : %s does not exist, please run `./arc -create-data=api` to create it", SWAGGER_DATA_FILE_PATH)
		fmt.Println()
		os.Exit(0)
	}

	b, err := ioutil.ReadFile(SWAGGER_DATA_FILE_PATH)

	if err != nil {
		panic(err)
	}

	var x map[string]interface{}
	err = json.Unmarshal(b, &x)

	if err != nil {
		panic(err)
	}

	str, err := json.Marshal(x)

	if err != nil {
		panic(err)
	}

	return string(str)
}

func readApiDataFile() ApiMetaData {
	_, noFileErr := os.Stat(API_DATA_FILE_PATH)

	if os.IsNotExist(noFileErr) {
		fmt.Printf("\nTemplate data file : %s does not exist, please run `./arc -create-data=api` to create it", API_DATA_FILE_PATH)
		fmt.Println()
		os.Exit(0)
	}

	b, err := ioutil.ReadFile(API_DATA_FILE_PATH)

	if err != nil {
		panic(err)
	}

	var api ApiMetaData
	err = json.Unmarshal(b, &api)

	if err != nil {
		panic(err)
	}

	return api
}


func createSwaggerDataFile()  {
	swagger := swagger.GetSwagger("Test", "1.0.0")

	_ = os.Mkdir(constants.DATA_FOLDER, 0777)

	content, _ := json.MarshalIndent(swagger, "", "    ")
	err := ioutil.WriteFile(SWAGGER_DATA_FILE_PATH, content, 0644)

	if err != nil {
		panic(err)
	}
}

func createApiDataFile()  {
	endpointConfig := EndpointConfig_t{}
	endpointConfig.EndpointType = "http"
	endpointConfig.ProductionEndpoints = &EndpointData_t{"http://localhost", nil}
	endpointConfig.SandboxEndpoints = &EndpointData_t{"http://localhost", nil}

	endpointJSON, err := json.Marshal(endpointConfig)

	if err != nil {
		panic(err)
	}

	apiMetaData := ApiMetaData{}
	apiMetaData.Name = "Test"
	apiMetaData.Version = "1.0.0"
	apiMetaData.Context = "test"
	apiMetaData.Description = "Test description"
	apiMetaData.EndpointConfig = string(endpointJSON)
	apiMetaData.Tiers = []string{"Unlimited"}
	apiMetaData.Visibility = "PUBLIC"
	apiMetaData.SubscriptionAvailability = "current_tenant"
	apiMetaData.Transport = []string{"http", "https"}


	_ = os.Mkdir(constants.DATA_FOLDER, 0777)

	content, _ := json.MarshalIndent(apiMetaData, "", "    ")
	err = ioutil.WriteFile(API_DATA_FILE_PATH, content, 0644)

	if err != nil {
		panic(err)
	}
}


func publisherGetAPIs(apiOptions *APIOptions, apiURL string, token string) {
	req := comm.CreateGet(apiURL)

	comm.SetRestAPIHeaders(token, req)

	values := url.Values{}

	comm.AddQueryParams(&values, req)

	comm.PrintRequest(constants.REST_API_REQUEST_LOG_STRING, req)

	resp := comm.SendHTTPRequest(req)

	defer resp.Body.Close()

	comm.PrintResponse(constants.REST_API_RESPONSE_LOG_STRING, resp)
}


func publisherCreateAPI(apiOptions *APIOptions, apiURL string, token string) {
	apiInfo := readApiDataFile()
	apiInfo.ApiDefinition = readSwaggerDataFile()

	data, err := json.Marshal(apiInfo)

	if err != nil {
		fmt.Printf("JSON marshaling failed: %s", err)
	}

	req := comm.CreatePost(apiURL, bytes.NewBuffer(data))

	comm.SetRestAPIHeaders(token, req)

	comm.PrintRequest(constants.REST_API_REQUEST_LOG_STRING, req)

	resp := comm.SendHTTPRequest(req)

	defer resp.Body.Close()

	comm.PrintResponse(constants.REST_API_RESPONSE_LOG_STRING, resp)
}
