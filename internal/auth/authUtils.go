package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func GetRandomState() (string, error) {
	//https://golang.org/pkg/crypto/rand/

	// Generates random bytes of length 'c'
	c := 8
	bytes := make([]byte, c)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Turn random bytes into URL safe string
	return base64.URLEncoding.EncodeToString(bytes), nil

}
