package wsockets

import (
	"encoding/json"
	"github.com/zacharyworks/huddle-gateway/auth"
	dataLayer "github.com/zacharyworks/huddle-gateway/data-layer"
	types "github.com/zacharyworks/huddle-shared/data"
	"log"
)

// Board alias
type Board types.Board

func (b Board) newJoinCode(client *Client) {
	code, err := auth.GetRandomString(2)
	if err != nil {
		log.Fatal(err)
	}
	newJoinCode := types.BoardJoinCode{
		b.BoardID, code}

	dataLayer.NewBoardJoinCode(newJoinCode)

	action, err := newAction("Board", "JoinCode", newJoinCode).build()

	if err != nil {
		log.Fatal(err)
	}

	client.send <- action
}

func joinBoard(actionPayload json.RawMessage, client *Client) {
	var joinBoard types.BoardJoin
	if err := json.Unmarshal(actionPayload, &joinBoard); err != nil {
		log.Fatal(err)
	}

	board := dataLayer.JoinBoard(joinBoard)

	// Build action to be dispatched
	action, err := json.Marshal(types.Action{
		ActionSubset:  "Board",
		ActionType:    "Create",
		ActionPayload: board,
	})

	if err != nil {
		log.Print(err)
	}

	// Dispatch action
	client.send <- action

}

func createBoard(actionPayload json.RawMessage, client *Client) {
	var newBoard types.NewBoard
	if err := json.Unmarshal(actionPayload, &newBoard); err != nil {
		log.Fatal(err)
	}

	board := dataLayer.NewBoard(newBoard)

	// Build action to be dispatched
	action, err := json.Marshal(types.Action{
		ActionSubset:  "Board",
		ActionType:    "Create",
		ActionPayload: board,
	})

	if err != nil {
		log.Print(err)
	}

	// Dispatch action
	client.send <- action
}

//}
//
//func (b Board) update(message []byte, ah actionHandler) {
//	// Update & Forward message
//	dataLayer.UpdateTodo(types.Todo(t))
//	for client := range ah.hub.clients {
//		if client.boardID == t.BoardFK {
//			client.send <- message
//		}
//	}
//}
//
//func (b Board) delete(message []byte, ah actionHandler) {
//	// Delete & Forward message
//	dataLayer.DeleteTodo(types.Todo(t))
//	for client := range ah.hub.clients {
//		if client.boardID == t.BoardFK {
//			client.send <- message
//		}
//	}
//}
