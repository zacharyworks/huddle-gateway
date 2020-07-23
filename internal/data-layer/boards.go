package dataLayer

import (
	types "github.com/zacharyworks/huddle-shared/data"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetBoardTodos(board types.Board) []types.Todo {

	// build URL
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("board/")
	url.WriteString(strconv.Itoa(board.BoardID))
	url.WriteString("/todos")
	println(url.String())
	// make request
	response, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}

	var todos []types.Todo
	processResponse(response, &todos)

	return todos
}

func GetUserBoards(userID string) []types.Board {
	// build URL
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("user/")
	url.WriteString(userID)
	url.WriteString("/boards")

	// make request
	response, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}

	var boards []types.Board
	processResponse(response, &boards)

	return boards
}
