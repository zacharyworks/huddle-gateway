package wsockets

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// ActionHandler struct maintains a reference to the hub
type actionHandler struct {
	hub *Hub
}

func newActionHandler(h *Hub) *actionHandler {
	return &actionHandler{
		hub: h}
}

var serverURL = "http://localhost:8081"
var httpClient = &http.Client{}

func (ah actionHandler) handle(message []byte, client *Client) {
	// Transform message into a map
	var actionMap map[string]json.RawMessage
	err := json.Unmarshal(message, &actionMap)
	if err != nil {
		log.Fatal(err)
	}

	actionType, err := strconv.Unquote(string(actionMap["ActionType"]))
	actionSubset, err := strconv.Unquote(string(actionMap["ActionSubset"]))

	if err != nil {
		log.Fatal(err)
	}

	// Convert todo from JSON into an structure
	var todo Todo
	switch actionSubset {
	case "Todo":
		err := json.Unmarshal(actionMap["ActionPayload"], &todo)
		if err != nil {
			log.Fatal(err)
		}
		switch actionType {
		case "Update":
			todo.update(actionMap, message, ah)
		case "Create":
			todo.create(actionMap, message, ah)
		case "Delete":
			todo.delete(actionMap, message, ah)
		}

	case "Session":
		switch actionType {
		case "RequestNew":
			requestSession(client)
		case "Exists":
			sessionExists(actionMap, client)
		}
	}
}
