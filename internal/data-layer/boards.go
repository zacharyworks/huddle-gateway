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

func NewBoard(board types.NewBoard) types.Board {
	// build URL
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("board")

	body, err := json.Marshal(board)
	if err != nil {
		log.Fatal(err)
	}

	// Make the request
	response := makeRequest(
		http.MethodPost,
		url.String(),
		bytes.NewBuffer(body))

	var newBoard types.Board
	processResponse(response, &newBoard)

	return newBoard
}

func NewBoardJoinCode(boardJoinCode types.BoardJoinCode) {
	// build URL
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("board/code")

	body, err := json.Marshal(boardJoinCode)
	if err != nil {
		log.Fatal(err)
	}

	// Make the request
	makeRequest(
		http.MethodPost,
		url.String(),
		bytes.NewBuffer(body))
}

func JoinBoard(boardJoin types.BoardJoin) (board types.Board) {
	// build URL
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("board/join")

	body, err := json.Marshal(boardJoin)
	if err != nil {
		log.Fatal(err)
	}

	// Make the request
	response := makeRequest(
		http.MethodPost,
		url.String(),
		bytes.NewBuffer(body))

	processResponse(response, &board)
	return
}
