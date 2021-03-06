package dcr

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomClientName() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}

func SetDCRParameters(request *DCRRequest, username string) {
	request.CallbackURL = "www.google.lk"
	request.ClientName = generateRandomClientName()
	request.TokenScope = "Production"
	request.Owner = username
	request.SupportedGrantTypes = "password refresh_token"
	request.IsSaaSApp = true
}
