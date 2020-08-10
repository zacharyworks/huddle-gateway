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
	code, err := auth.GetRandomString(5)
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

func (b Board) deleteMember(client *Client) {
	boardMember := types.BoardMember{
		BoardFK: b.BoardID,
		UserFK:  client.user.OauthID,
	}
	err := dataLayer.LeaveBoard(boardMember)
	if err != nil {
		log.Fatal(err)
	}

	action, err := json.Marshal(types.Action{
		Subset:  "Board",
		Type:    "Remove",
		Payload: b,
	})

	if err != nil {
		log.Print(err)
	}

	// Dispatch action
	client.send <- action
}

func joinBoard(Payload json.RawMessage, client *Client) {
	var joinBoard types.BoardJoin
	if err := json.Unmarshal(Payload, &joinBoard); err != nil {
		log.Fatal(err)
	}

	board := dataLayer.JoinBoard(joinBoard)

	// Build action to be dispatched
	action, err := json.Marshal(types.Action{
		Subset:  "Board",
		Type:    "Create",
		Payload: board,
	})

	if err != nil {
		log.Print(err)
	}

	// Dispatch action
	client.send <- action
}

func createBoard(Payload json.RawMessage, client *Client) {
	var newBoard types.NewBoard
	if err := json.Unmarshal(Payload, &newBoard); err != nil {
		log.Fatal(err)
	}

	board := dataLayer.NewBoard(newBoard)
	client.boardMembership = append(client.boardMembership, board.BoardID)

	// Build action to be dispatched
	action, err := json.Marshal(types.Action{
		Subset:  "Board",
		Type:    "Create",
		Payload: board,
	})

	if err != nil {
		log.Print(err)
	}

	// Dispatch action
	client.send <- action
}
