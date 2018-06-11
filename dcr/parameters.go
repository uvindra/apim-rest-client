package dcr

import (
	"crypto/rand"
	"encoding/base32"
)

func generateRandomClientName() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		panic(err)
	}

	return base32.StdEncoding.EncodeToString(b)
}

func SetDCRParameters(request *DCRRequest, username string) {
	request.CallbackURL = "www.google.lk"
	request.ClientName = generateRandomClientName()
	request.TokenScope = "Production"
	request.Owner = username
	request.SupportedGrantTypes = "password refresh_token"
	request.IsSaaSApp = true
}
