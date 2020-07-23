package wsockets

import (
	"../auth"
	"../data-layer"
	"encoding/json"
	"github.com/zacharyworks/huddle-shared/data"
	"log"
)

func openBoard(actionPayload json.RawMessage, client *Client) {
	var board types.Board
	json.Unmarshal(actionPayload, &board)
	// register client as being subscribed to this board
	client.boardID = board.BoardID
	todos := dataLayer.GetBoardTodos(board)

	initBoardAction, err := newAction("Session", "Board", board).build()
	if err != nil {
		log.Fatal(err)
	}
	client.send <- initBoardAction

	initBoardTodoAction, err := newAction("Todo", "Init", todos).build()
	if err != nil {
		log.Fatal(err)
	}
	client.send <- initBoardTodoAction
}

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
	session := dataLayer.RetrieveSession(sessionID)
	if session.UserFK != "" {
		// A user was found!!
		sessionJSON, err := json.Marshal(session)
		client.authorised = true
		if err != nil {
			panic(err)
		}
		client.send <- sessionJSON
	} else {
		// No user was found
		return
	}

	// Get the users information
	user := dataLayer.GetUser(session.UserFK)

	// Retrieve the boards the user is a member of
	boards := dataLayer.GetUserBoards(user.OauthID)

	userAction, err := newAction("Session", "User", user).build()
	client.send <- userAction

	boardAction, err := newAction("Board", "Init", boards).build()
	client.send <- boardAction
}
