package wsockets

import (
	"encoding/json"
	"github.com/zacbriggssagecom/huddle/server/sharedinternal/data"
	"log"
	// "io/ioutil"
	// "net/http"
)

// Hub struct maintains the list of clients which are active
type Hub struct {
	// Clients that are registered
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister request from clients
	unregister chan *Client

	// Action Handler
	actionHandler *actionHandler
}

// NewHub creates a hub
func NewHub() *Hub {
	hub := Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool)}
	hub.actionHandler = newActionHandler(&hub)
	return &hub
}

// Run runs the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register: // On a clients first connection
			var err error
			h.clients[client] = true
			action, err := json.Marshal(types.StringAction{
				ActionSubset:  "Session",
				ActionType:    "Connected",
				ActionPayload: "",
			})
			if err != nil {
				log.Panic(err)
			}
			client.send <- action

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
