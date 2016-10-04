package cmd

import (
	"net/url"
	"encoding/json"
	"fmt"
	"os"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"apim-rest-client/comm"
	"apim-rest-client/constants"
)

const DATA_FOLDER = "data"
const DATA_FILE_PATH = DATA_FOLDER + string(filepath.Separator) + "product.json";

type ProductMetaData struct {
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Version        string  `json:"version"`
	Provider       string   `json:"provider"`
	ThrottlingTier []string `json:"throttlingTiers"`
	Visibility     string   `json:"visibility"`
}

func readProductDataFile() ProductMetaData {
	_, noFileErr := os.Stat(DATA_FILE_PATH)

	if os.IsNotExist(noFileErr) {
		fmt.Println("\nTemplate data file : %s does not exist, please run `./arc -create-data=product` to create it", DATA_FILE_PATH)
		fmt.Println()
		os.Exit(0)
	}

	b, err := ioutil.ReadFile(DATA_FILE_PATH)

	if err != nil {
		panic(err)
	}

	var product ProductMetaData
	err = json.Unmarshal(b, &product)

	if err != nil {
		panic(err)
	}

	return product
}

func createProductDataFile() {
	var productMetaData ProductMetaData

	productMetaData.Name = "Test"
	productMetaData.Version = "1.0.0"
	productMetaData.Description = "Test Description"
	productMetaData.Provider = "admin"
	productMetaData.ThrottlingTier = []string{"Unlimited", "Gold"}
	productMetaData.Visibility = "PUBLIC"

	_ = os.Mkdir(DATA_FOLDER, 0777)

	content, _ := json.MarshalIndent(productMetaData, "", "    ")
	err := ioutil.WriteFile(DATA_FILE_PATH, content, 0644)

	if err != nil {
		panic(err)
	}
}


func publisherGetProducts(apiOptions *APIOptions, productURL string, token string) {
	req := comm.CreateGet(productURL)

	comm.SetRestAPIHeaders(token, req)

	values := url.Values{}

	if apiOptions.Limit != constants.UNDEFINED_INT {
		values.Add(constants.LIMIT_KEY, string(apiOptions.Limit))
	}

	if apiOptions.Offset != constants.UNDEFINED_INT {
		values.Add(constants.OFFSET_KEY, string(apiOptions.Offset))
	}

	if apiOptions.Query != constants.UNDEFINED_STRING {
		values.Add(constants.QUERY_KEY, apiOptions.Query)
	}

	comm.AddQueryParams(&values, req)

	comm.PrintRequest(constants.REST_API_REQUEST_LOG_STRING, req)

	resp := comm.SendHTTPRequest(req)

	defer resp.Body.Close()

	comm.PrintResponse(constants.REST_API_RESPONSE_LOG_STRING, resp)
}


func publisherCreateProduct(apiOptions *APIOptions, productURL string, token string) {
	productInfo := readProductDataFile()

	data, err := json.Marshal(productInfo)

	if err != nil {
		fmt.Println("JSON marshaling failed: %s", err)
	}

	req := comm.CreatePost(productURL, bytes.NewBuffer(data))

	comm.SetRestAPIHeaders(token, req)

	comm.PrintRequest(constants.REST_API_REQUEST_LOG_STRING, req)

	resp := comm.SendHTTPRequest(req)

	defer resp.Body.Close()

	comm.PrintResponse(constants.REST_API_RESPONSE_LOG_STRING, resp)
}
