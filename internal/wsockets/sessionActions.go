package wsockets

import (
	"encoding/json"
	"github.com/zacbriggssagecom/huddle/server/gateway/internal/auth"
	"github.com/zacbriggssagecom/huddle/server/sharedinternal/data"
	"io/ioutil"
	"log"
	"net/http"
)

func requestSession(client *Client) {
	var err error
	// Get a session id to provide the
	// user so they can hithertoforth assert
	// who they are in our state.
	client.SessionID, err = auth.GetRandomState()
	if err != nil {
		log.Fatal(err)
	}

	action, err := json.Marshal(types.StringAction{
		ActionSubset:  "Session",
		ActionType:    "SessionID",
		ActionPayload: client.SessionID,
	})

	client.send <- action
}

func sessionExists(actionMap map[string]json.RawMessage, client *Client) {
	// Get the provided session ID
	var sessionID string
	err := json.Unmarshal(actionMap["ActionPayload"], &sessionID)
	if err != nil {
		log.Fatal(err)
	}

	// Update client to reflect provided session ID
	client.SessionID = sessionID

	// Lookup the session in the DB
	session := auth.RetreiveSessionByID(sessionID)
	if session.UserFK != "" {
		// A user was found
		sessionJSON, err := json.Marshal(session)
		if err != nil {
			panic(err)
		}
		client.send <- sessionJSON
	} else {
		// No user was found
		println("not authenticated")
		return
	}

	// Send todos to the user
	payload, err := http.Get("http://localhost:8081/todos")
	var actionPayload []types.Todo
	body, err := ioutil.ReadAll(payload.Body)
	if err != nil {
		log.Print(err)
	}
	json.Unmarshal(body, &actionPayload)

	action, err := json.Marshal(types.TodoAction{
		ActionSubset:  "Todo",
		ActionType:    "Init",
		ActionPayload: actionPayload,
	})
	if err != nil {
		log.Print(err)
	}

	client.send <- action
}
