package wsockets

import (
	"encoding/json"
	dataLayer "github.com/zacharyworks/huddle-gateway/data-layer"
	"github.com/zacharyworks/huddle-shared/data"
	"log"
)

// Todo alias
type Todo types.Todo

func (t Todo) create(ah actionHandler) {

	// Build to-do from response
	todo := dataLayer.NewTodo(types.Todo(t))

	// Build action to be dispatched
	action, err := json.Marshal(types.Action{
		Subset:  "Todo",
		Type:    "Create",
		Payload: todo,
	})
	if err != nil {
		log.Print(err)
	}

	// Dispatch action
	for client := range ah.hub.clients {
		if client.selectedBoardID == t.BoardFK {
			client.send <- action
		}
	}
}

func (t Todo) update(message []byte, ah actionHandler) {
	// Update & Forward message
	dataLayer.UpdateTodo(types.Todo(t))
	for client := range ah.hub.clients {
		if client.selectedBoardID == t.BoardFK {
			client.send <- message
		}
	}
}

func (t Todo) delete(message []byte, ah actionHandler) {
	// Delete & Forward message
	dataLayer.DeleteTodo(types.Todo(t))
	for client := range ah.hub.clients {
		if client.selectedBoardID == t.BoardFK {
			client.send <- message
		}
	}
}
