package dataLayer

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

var httpClient = &http.Client{}
var restEndpoint = "http://localhost:8000/"

func processResponse(response *http.Response, data interface{}) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		println(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		println(err)
	}
}

func makeRequest(method string, url string, body io.Reader) *http.Response {
	request, err := http.NewRequest(
		method,
		url,
		body)
	if err != nil {
		println(err)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		println(err)
	}

	return response
}
