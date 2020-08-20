package auth

import (
	"bytes"
	"encoding/json"
	"github.com/zacharyworks/huddle-shared/data"
	"net/http"
	"strings"
)

var serverURL = "http://localhost:8000"
var httpClient = &http.Client{}

func postOauthUser(oUser types.Response) {
	newUser := types.User{
		OauthID:    oUser.ID,
		Email:      oUser.Email,
		Picture:    oUser.Picture,
		Name:       oUser.Name,
		GivenName:  oUser.GivenName,
		FamilyName: oUser.FamilyName,
	}
	newUserJSON, err := json.Marshal(newUser)
	if err != nil {
		println(err)
	}

	// Build the url
	var postURL strings.Builder
	postURL.WriteString(serverURL)
	postURL.WriteString("/user")

	// Build the request
	req, err := http.NewRequest(
		http.MethodPost, postURL.String(), bytes.NewBuffer(newUserJSON))

	if err != nil {
		println(err)
	}

	// Execute the request, fetch response
	_, err = httpClient.Do(req)
	if err != nil {
		println(err)
	}
}

func putOauthUser(oUser types.Response) {
	newUser := types.User{
		OauthID:    oUser.ID,
		Email:      oUser.Email,
		Picture:    oUser.Picture,
		Name:       oUser.Name,
		GivenName:  oUser.GivenName,
		FamilyName: oUser.FamilyName,
	}
	newUserJSON, err := json.Marshal(newUser)
	if err != nil {
		println(err)
	}

	// Build the url
	var putURL strings.Builder
	putURL.WriteString(serverURL)
	putURL.WriteString("/user")

	// Build the request
	req, err := http.NewRequest(
		http.MethodPut, putURL.String(), bytes.NewBuffer(newUserJSON))

	if err != nil {
		println(err)
	}

	// Execute the request, fetch response
	_, err = httpClient.Do(req)
	if err != nil {
		println(err)
	}
}
