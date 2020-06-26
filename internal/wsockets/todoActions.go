package wsockets

import (
	"bytes"
	"encoding/json"
	"github.com/zacbriggssagecom/huddle/server/sharedinternal/data"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Todo alias
type Todo types.Todo

func (t Todo) update(actionMap map[string]json.RawMessage, message []byte, ah actionHandler) {

	// Build the url
	var putURL strings.Builder
	putURL.WriteString(serverURL)
	putURL.WriteString("/todos/")
	putURL.WriteString(strconv.Itoa(t.TodoID))

	// Build the request
	req, err := http.NewRequest(
		http.MethodPut, putURL.String(), bytes.NewBuffer(actionMap["ActionPayload"]))

	if err != nil {
		log.Fatal(err)
	}

	// Execute the request
	_, err = httpClient.Do(req)
	ah.hub.broadcast <- message
}

func (t Todo) create(actionMap map[string]json.RawMessage, message []byte, ah actionHandler) {

	// Build the url
	var putURL strings.Builder
	putURL.WriteString(serverURL)
	putURL.WriteString("/todos")

	// Build the request
	req, err := http.NewRequest(
		http.MethodPost, putURL.String(), bytes.NewBuffer(actionMap["ActionPayload"]))

	if err != nil {
		log.Fatal(err)
	}

	// Execute the request, fetch response
	resp, err := httpClient.Do(req)

	// Build todo from response
	var todoPayload types.Todo
	todoJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	json.Unmarshal(todoJSON, &todoPayload)

	// Build action to be dispatched
	action, err := json.Marshal(types.TodoActionSingle{
		ActionSubset:  "Todo",
		ActionType:    "Create",
		ActionPayload: todoPayload,
	})
	if err != nil {
		log.Print(err)
	}

	// Dispatch action
	ah.hub.broadcast <- action
}

func (t Todo) delete(actionMap map[string]json.RawMessage, message []byte, ah actionHandler) {

	// Build the URL
	var putURL strings.Builder
	putURL.WriteString(serverURL)
	putURL.WriteString("/todos/")
	putURL.WriteString(strconv.Itoa(t.TodoID))

	// Build the request
	req, err := http.NewRequest(
		http.MethodDelete, putURL.String(), bytes.NewBuffer(actionMap["ActionPayload"]))

	if err != nil {
		log.Fatal(err)
	}

	// Execute the request & dispatch message
	_, err = httpClient.Do(req)
	ah.hub.broadcast <- message
}
