package auth

import (
	"bytes"
	"encoding/json"
	"github.com/zacbriggssagecom/huddle/server/sharedinternal/data"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var serverURL = "http://localhost:8081"
var httpClient = &http.Client{}

func SaveNewSession(sessionID string, randomState string) {
	// Build the url
	var putURL strings.Builder
	putURL.WriteString(serverURL)
	putURL.WriteString("/session")

	newSession := types.Session{
		SessionID: sessionID,
		State:     randomState}

	newSessionJSON, err := json.Marshal(newSession)
	if err != nil {
		log.Fatal(err)
	}

	// Build the request
	req, err := http.NewRequest(
		http.MethodPost, putURL.String(), bytes.NewBuffer(newSessionJSON))

	if err != nil {
		log.Fatal(err)
	}

	// Execute the request, fetch response
	_, err = httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

}

func RetreiveSessionByID(id string) types.Session {
	// Build the url
	var putURL strings.Builder
	putURL.WriteString(serverURL)
	putURL.WriteString("/session/id/")
	putURL.WriteString(id)

	resp, err := http.Get(putURL.String())
	if err != nil {
		log.Fatal(err)
	}
	var session types.Session
	sessionJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(sessionJSON, &session)
	if err != nil {
		log.Fatal(err)
	}

	return session
}

func RetreiveSessionByState(state string) (types.Session, error) {
	// Build the url
	var putURL strings.Builder
	putURL.WriteString(serverURL)
	putURL.WriteString("/session/state/")
	putURL.WriteString(state)

	resp, err := http.Get(putURL.String())
	if err != nil {
		return types.Session{}, err
	}
	var session types.Session
	sessionJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(sessionJSON, &session)
	if err != nil {
		log.Fatal(err)
	}

	return session, nil
}

func updateSession(session types.Session, userFK string) {
	// Build the url
	var putURL strings.Builder
	putURL.WriteString(serverURL)
	putURL.WriteString("/session/adduser")

	session.UserFK = userFK
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		panic(err)
	}
	// Build the request
	req, err := http.NewRequest(
		http.MethodPut, putURL.String(), bytes.NewBuffer(sessionJSON))

	if err != nil {
		log.Fatal(err)
	}

	// Execute the request
	_, err = httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

func postOauthUser(oUser types.Response) {
	newUser := types.User{
		OauthID: oUser.Id,
		Email:   oUser.Email}
	newUserJSON, err := json.Marshal(newUser)
	if err != nil {
		log.Fatal(err)
	}

	// Build the url
	var putURL strings.Builder
	putURL.WriteString(serverURL)
	putURL.WriteString("/user")

	// Build the request
	req, err := http.NewRequest(
		http.MethodPost, putURL.String(), bytes.NewBuffer(newUserJSON))

	if err != nil {
		log.Fatal(err)
	}

	// Execute the request, fetch response
	_, err = httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

func getUser(oauthid string) string {
	// Build the url
	var getURL strings.Builder
	getURL.WriteString(serverURL)
	getURL.WriteString("/user/")
	getURL.WriteString(oauthid)

	resp, err := http.Get(getURL.String())

	if err != nil {
		log.Fatal(err)
	}

	var user types.User
	userJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		log.Fatal(err)
	}
	return user.OauthID
}
