package dataLayer

import (
	types "github.com/zacharyworks/huddle-shared/data"
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
		println(err)
	}

	var user types.User
	// no user was found
	if response.StatusCode == http.StatusNotFound {
		return user
	}

	processResponse(response, &user)

	return user
}
