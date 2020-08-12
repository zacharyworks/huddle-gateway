package wsockets

import (
	"encoding/json"
	types "github.com/zacharyworks/huddle-shared/data"
	"log"
)

// hubSession keeps track of the current state
// of Huddle. e.g. what boards users have open
// and other miscellaneous information.
type hubSession struct {
	boardClientMap map[int]map[*Client]bool
	hub            *Hub
}

func newSession(hub *Hub) *hubSession {
	hubSession := hubSession{
		make(map[int]map[*Client]bool),
		hub,
	}
	return &hubSession
}

func (hs hubSession) clientLeft(client *Client) {
	action := types.Action{Subset: "Todo", Type: "PeerLeft", Payload: client.user}

	for id, bMap := range hs.boardClientMap {
		if bMap[client] == true {
			hs.notifyBoard(id, action)
			delete(bMap, client)
		}
	}
}

func (hs hubSession) notifyBoard(boardID int, a types.Action) {
	action, err := json.Marshal(types.Action{
		a.Subset,
		a.Type,
		a.Payload,
	})
	if err != nil {
		log.Fatal(err)
	}
	for c := range hs.boardClientMap[boardID] {
		c.send <- action
	}
}

func (hs hubSession) addUserToBoard(boardID int, client *Client) {
	// Does the board have a set of users, create one if not
	if hs.boardClientMap[boardID] == nil {
		hs.boardClientMap[boardID] = make(map[*Client]bool)
	}

	// Remove the user from their current board
	hs.removeUserFromBoard(boardID, client)

	// Add user to the new board
	action := types.Action{Subset: "Todo", Type: "Peer", Payload: client.user}
	hs.notifyBoard(boardID, action)

	// Add the user to the board
	hs.boardClientMap[boardID][client] = true
	client.selectedBoardID = boardID
}

func (hs hubSession) removeUserFromBoard(boardId int, client *Client) {
	if client.selectedBoardID == 0 {
		return
	}

	// remove from map
	boardClients := hs.boardClientMap[client.selectedBoardID]
	if boardClients[client] == true {
		delete(hs.boardClientMap[client.selectedBoardID], client)
	}

	// notify peers
	action := types.Action{Subset: "Todo", Type: "PeerLeft", Payload: client.user}
	hs.notifyBoard(client.selectedBoardID, action)

	client.selectedBoardID = 0
	return
}

func (hs hubSession) notifyClientOfPeers(client *Client) {
	var peers []types.User
	for client := range hs.boardClientMap[client.selectedBoardID] {
		peers = append(peers, client.user)
	}

	action, err := newAction("Todo", "Peers", peers).build()
	if err != nil {
		log.Fatal(err)
	}

	client.send <- action
}
