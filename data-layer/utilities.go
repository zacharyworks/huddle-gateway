package dataLayer

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var httpClient = &http.Client{}
var restEndpoint = "http://localhost:8081/"

func processResponse(response *http.Response, data interface{}) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
}

func makeRequest(method string, url string, body io.Reader) *http.Response {
	request, err := http.NewRequest(
		method,
		url,
		body)
	if err != nil {
		log.Fatal(err)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	return response
}
