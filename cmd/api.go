package cmd

import (
	"net/url"
	"encoding/json"
	"log"
	"bytes"
	"apim_rest_client/comm"
	"apim_rest_client/constants"
)


type EndpointConfig struct {
	ProductionEndpoints *EndpointData `json:"production_endpoints"`
	SandboxEndpoints *EndpointData `json:"sandbox_endpoints"`
	EndpointType string `json:"endpoint_type"`
}

type EndpointData struct {
	URL string `json:"url"`
	Config string `json:"config"`
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


func getAPIsPublisher(flags *Flags, apiURL string, token string) {
	req := comm.CreateGet(apiURL)

	comm.SetRestAPIHeaders(token, req)

	values := url.Values{}

	if flags.Limit != constants.UNDEFINED_INT {
		values.Add(constants.LIMIT_KEY, string(flags.Limit))
	}

	if flags.Offset != constants.UNDEFINED_INT {
		values.Add(constants.OFFSET_KEY, string(flags.Offset))
	}

	if flags.Query != constants.UNDEFINED_STRING {
		values.Add(constants.QUERY_KEY, flags.Query)
	}

	comm.AddQueryParams(&values, req)

	resp := comm.SendHTTPRequest(req)

	defer resp.Body.Close()

	comm.PrintResponse(resp)
}


func createAPIDefinition() *ApiMetaData {
	apiMetaData := ApiMetaData{}

	apiMetaData.Name = "two"
	apiMetaData.Context = "two"
	apiMetaData.Version = "1.0.0"
	apiMetaData.ApiDefinition = "{\"swagger\":\"2.0\",\"paths\":{\"/what\":{\"get\":{\"responses\":{\"200\":{\"description\":\"\"}},\"x-auth-type\":\"Application & Application User\",\"x-throttling-tier\":\"Unlimited\"}},\"/how\":{\"get\":{\"responses\":{\"200\":{\"description\":\"\"}},\"x-auth-type\":\"Application & Application User\",\"x-throttling-tier\":\"Unlimited\"}},\"/where\":{\"get\":{\"responses\":{\"200\":{\"description\":\"\"}},\"x-auth-type\":\"Application & Application User\",\"x-throttling-tier\":\"Unlimited\"}}},\"info\":{\"title\":\"one\",\"version\":\"1.0.0\"}}"
	apiMetaData.EndpointConfig = "{\"production_endpoints\":{\"url\":\"https://localhost:9443/am/sample/pizzashack/v1/api/\",\"config\":null},\"sandbox_endpoints\":{\"url\":\"https://localhost:9443/am/sample/pizzashack/v1/api/\",\"config\":null},\"endpoint_type\":\"http\"}"
	apiMetaData.Tiers = []string{"Unlimited"}
	apiMetaData.Visibility = "PUBLIC"
	apiMetaData.SubscriptionAvailability = "current_tenant"
	apiMetaData.Transport = []string{"http", "https"}

	return &apiMetaData
}

func publisherCreateAPI(flags *Flags, apiURL string, token string) {

	apiInfo := createAPIDefinition()

	data, err := json.Marshal(apiInfo)

	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	req := comm.CreatePost(apiURL, bytes.NewBuffer(data))

	comm.SetRestAPIHeaders(token, req)

	resp := comm.SendHTTPRequest(req)

	defer resp.Body.Close()

	comm.PrintResponse(resp)
}
