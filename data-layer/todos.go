package dataLayer

import (
	"bytes"
	"encoding/json"
	types "github.com/zacharyworks/huddle-shared/data"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func NewTodo(todo types.Todo) types.Todo {

	// Build the url
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("todos")

	body, err := json.Marshal(todo)
	if err != nil {
		log.Fatal(err)
	}

	// Make the request
	response := makeRequest(
		http.MethodPost,
		url.String(),
		bytes.NewBuffer(body))

	processResponse(response, &todo)
	return todo
}

func AllTodos() (todos []types.Todo) {

	// Build the url
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("todos")

	// Make the request
	response, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}

	processResponse(response, &todos)
	return todos
}

func UpdateTodo(todo types.Todo) {

	// Build the url
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("todos/")
	url.WriteString(strconv.Itoa(todo.TodoID))

	body, err := json.Marshal(todo)
	if err != nil {
		log.Fatal(err)
	}

	makeRequest(
		http.MethodPut,
		url.String(),
		bytes.NewBuffer(body))
}

func DeleteTodo(todo types.Todo) {

	// Build the URL
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("todos/")
	url.WriteString(strconv.Itoa(todo.TodoID))

	body, err := json.Marshal(todo)
	if err != nil {
		log.Fatal(err)
	}

	// Make the request
	makeRequest(
		http.MethodDelete,
		url.String(),
		bytes.NewBuffer(body))
}
