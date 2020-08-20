package dataLayer

import (
	"bytes"
	"encoding/json"
	types "github.com/zacharyworks/huddle-shared/data"
	"net/http"
	"strings"
)

func UpdateSession(session types.Session, userFK string) {
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("session/adduser")

	session.UserFK = userFK
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		println(err)
	}

	// Build the request
	makeRequest(
		http.MethodPut,
		url.String(),
		bytes.NewBuffer(sessionJSON))

}

func RetrieveSession(ID string) types.Session {
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("session/id/")
	url.WriteString(ID)

	response, err := http.Get(url.String())
	if err != nil {
		println(err)
	}

	var session types.Session

	if response.StatusCode == http.StatusNotFound {
		return session
	}

	processResponse(response, &session)
	return session
}

func RetrieveSessionByState(state string) (types.Session, error) {
	// Build the url
	var putURL strings.Builder
	putURL.WriteString(restEndpoint)
	putURL.WriteString("session/state/")
	putURL.WriteString(state)

	response, err := http.Get(putURL.String())
	if err != nil {
		return types.Session{}, err
	}

	var session types.Session
	processResponse(response, &session)

	return session, nil
}

func SaveNewSession(sessionID string, randomState string) {

	// Build the url
	var putURL strings.Builder
	putURL.WriteString(restEndpoint)
	putURL.WriteString("session")

	newSession := types.Session{
		SessionID: sessionID,
		State:     randomState}

	newSessionJSON, err := json.Marshal(newSession)
	if err != nil {
		println(err)
	}

	// Build the request
	makeRequest(
		http.MethodPost,
		putURL.String(),
		bytes.NewBuffer(newSessionJSON))
}
