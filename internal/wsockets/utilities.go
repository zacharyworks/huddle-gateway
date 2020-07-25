package wsockets

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	types "github.com/zacharyworks/huddle-shared/data"
)

type action struct {
	Subset  string
	Type    string
	Payload interface{}
}

func newAction(actionSubset string, actionType string, actionPayload interface{}) *action {
	return &action{
		actionSubset,
		actionType,
		actionPayload,
	}
}

func (a action) build() ([]byte, error) {
	action, err := json.Marshal(types.Action{
		a.Subset,
		a.Type,
		a.Payload,
	})

	return action, err
}

func GetRandomString() (string, error) {
	//https://golang.org/pkg/crypto/rand/
	// Generates random bytes of length 'c'
	c := 4
	bytes := make([]byte, c)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Turn random bytes into URL safe string
	return base64.URLEncoding.EncodeToString(bytes), nil
}
