package wsockets

import (
	"encoding/json"
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
		println(err)
	}

	Subset, err := strconv.Unquote(string(actionMap["subset"]))
	Type, err := strconv.Unquote(string(actionMap["type"]))
	Payload := actionMap["payload"]

	if err != nil {
		println(err)
	}

	switch Subset {
	case "Todo":
		var todo Todo
		if err := json.Unmarshal(actionMap["payload"], &todo); err != nil {
			println(err)
		}
		switch Type {
		case "Create":
			todo.create(ah)
		case "Update":
			todo.update(message, ah)
		case "Delete":
			todo.delete(message, ah)
		}

	case "Session":
		switch Type {
		case "RequestNew":
			requestSession(client)
		case "Exists":
			sessionExists(actionMap, client)
		case "OpenBoard":
			openBoard(Payload, client, ah)
		case "Select":
			selectTodo(actionMap, client, ah)
		}

	case "Board":
		var board Board
		if err := json.Unmarshal(actionMap["payload"], &board); err != nil {
			println(err)
		}

		switch Type {
		case "GetJoinCode":
			board.newJoinCode(client)
		case "Leave":
			board.deleteMember(client)
		case "BoardJoin":
			joinBoard(Payload, client)
		case "BoardNew":
			createBoard(Payload, client)
		}
	}
}
