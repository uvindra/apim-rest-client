package cmd

import (
	"encoding/json"
	"log"
	"bytes"
	"apim_rest_client/comm"
)

type ProductMetaData struct {
	Name    	string `json:"name"`
	Description     string `json:"description"`
	Version    	string `json:"version"`
	Provider    	string `json:"provider"`
	ThrottlingTier  []string `json:"throttlingTier"`
	Visibility 	string `json:"visibility"`
}

func createProductDefinition() *ProductMetaData {
	productMetaData := ProductMetaData{}

	productMetaData.Name = "Test"
	productMetaData.Version = "1.0.0"
	productMetaData.ThrottlingTier = []string{"Unlimited"}
	productMetaData.Visibility = "PUBLIC"

	return &productMetaData
}


func publisherCreateProduct(flags *Flags, productURL string, token string) {

	productInfo := createProductDefinition()

	data, err := json.Marshal(productInfo)

	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	req := comm.CreatePost(productURL, bytes.NewBuffer(data))

	comm.SetRestAPIHeaders(token, req)

	resp := comm.SendHTTPRequest(req)

	defer resp.Body.Close()

	comm.PrintResponse(resp)
}
