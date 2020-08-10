package auth

import (
	"crypto/rand"
	"encoding/base64"
)

//TODO: Must check if generated string already exists in DB, if so: regenerate
func GetRandomString(length int) (string, error) {
	//https://golang.org/pkg/crypto/rand/
	// Generates random bytes of length 'c'
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(bytes)

	// Turn random bytes into URL safe string
	return randomString, nil
}
