package dataLayer

import (
	types "github.com/zacharyworks/huddle-shared/data"
	"log"
	"net/http"
	"strings"
)

func GetUser(ID string) types.User {
	// build URL
	var url strings.Builder
	url.WriteString(restEndpoint)
	url.WriteString("user/")
	url.WriteString(ID)

	// make request
	response, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}

	var user types.User
	processResponse(response, &user)

	return user
}
