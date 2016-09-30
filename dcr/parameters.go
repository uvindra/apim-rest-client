package dcr

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomClientName() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}


func SetDCRParameters(request *DCRRequest) {
	request.CallbackURL = "www.google.lk"
	request.ClientName = generateRandomClientName()
	request.TokenScope = "Production"
	request.Owner = "admin"
	request.SupportedGrantTypes = "password refresh_token"
	request.IsSaaSApp = true
}
