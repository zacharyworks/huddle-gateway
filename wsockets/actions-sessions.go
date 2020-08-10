package wsockets

import (
	"encoding/json"
	"github.com/zacharyworks/huddle-gateway/auth"
	dataLayer "github.com/zacharyworks/huddle-gateway/data-layer"
	"github.com/zacharyworks/huddle-shared/data"
	"log"
)

func openBoard(Payload json.RawMessage, client *Client, ah actionHandler) {
	var board types.Board
	json.Unmarshal(Payload, &board)
	// register client as being subscribed to this board

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

	// Assuming board is now open, update sessions
	ah.hub.hubSession.addUserToBoard(board.BoardID, client)
	ah.hub.hubSession.notifyClientOfPeers(client)
}

func requestSession(client *Client) {
	var err error
	// Get a session id to provide the
	// user so they can hithertoforth assert
	// who they are in our state.
	client.SessionID, err = auth.GetRandomString(8)
	if err != nil {
		log.Fatal(err)
	}

	action, err := json.Marshal(types.StringAction{
		Subset:  "Session",
		Type:    "SessionID",
		Payload: client.SessionID,
	})

	client.send <- action
}

func sessionExists(actionMap map[string]json.RawMessage, client *Client) {
	// Get the provided session ID
	var sessionID string
	err := json.Unmarshal(actionMap["payload"], &sessionID)
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
			log.Fatal(err)
		}
		client.send <- sessionJSON
	} else {
		// No user was found
		requestSession(client)
		return
	}

	// Get the users information
	user := dataLayer.GetUser(session.UserFK)
	client.user = user

	// Retrieve the boards the user is a member of
	boards := dataLayer.GetUserBoards(user.OauthID)

	// Add the users authorised boards to their list of membership
	for _, v := range boards {
		client.boardMembership = append(client.boardMembership, v.BoardID)
	}

	userAction, err := newAction("Session", "User", user).build()
	client.send <- userAction

	boardAction, err := newAction("Board", "Init", boards).build()
	client.send <- boardAction
}
