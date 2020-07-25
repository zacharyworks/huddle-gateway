package wsockets

import (
	"encoding/json"
	"log"
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

func (ah actionHandler) handle(message []byte, client *Client) {
	// Transform message into a map
	var actionMap map[string]json.RawMessage
	err := json.Unmarshal(message, &actionMap)
	if err != nil {
		log.Fatal(err)
	}

	actionSubset, err := strconv.Unquote(string(actionMap["ActionSubset"]))
	actionType, err := strconv.Unquote(string(actionMap["ActionType"]))
	actionPayload := actionMap["ActionPayload"]

	if err != nil {
		log.Fatal(err)
	}

	switch actionSubset {
	case "Todo":
		var todo Todo
		if err := json.Unmarshal(actionMap["ActionPayload"], &todo); err != nil {
			log.Fatal(err)
		}

		switch actionType {
		case "Create":
			todo.create(ah)
		case "Update":
			todo.update(message, ah)
		case "Delete":
			todo.delete(message, ah)
		}

	case "Session":
		switch actionType {
		case "RequestNew":
			requestSession(client)
		case "Exists":
			sessionExists(actionMap, client)
		case "OpenBoard":
			openBoard(actionPayload, client)
		}

	case "Board":
		var board Board
		if actionType != "NewBoard" && actionType != "JoinBoard" {
			if err := json.Unmarshal(actionMap["ActionPayload"], &board); err != nil {
				log.Fatal(err)
			}
		}
		switch actionType {
		case "JoinBoard":
			joinBoard(actionPayload, client)
		case "GetJoinCode":
			board.newJoinCode(client)
		case "NewBoard":
			createBoard(actionPayload, client)
		}

	}
}
