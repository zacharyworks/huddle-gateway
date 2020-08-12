package main

import (
	"encoding/json"
	"flag"
	"github.com/zacharyworks/huddle-gateway/auth"
	"github.com/zacharyworks/huddle-gateway/wsockets"
	"io/ioutil"
	"log"
	"net/http"
)

type credentials struct {
	ClientID     string
	ClientSecret string
	RestService  string
	LocalService string
}

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	credentials := readAuthCredentials()
	flag.Parse()
	// Create web socket service
	hub := wsockets.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsockets.ServeWs(hub, w, r)
	})

	//Create auth service
	oAuth := auth.NewAuth(credentials.ClientID, credentials.ClientSecret, credentials.LocalService)
	http.HandleFunc("/login", oAuth.HandleLogin)
	http.HandleFunc("/callback", oAuth.HandleCallback)

	// Start database
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}

func readAuthCredentials() credentials {
	file, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatal(err)
	}
	credsmap := make(map[string]string)
	err = json.Unmarshal([]byte(file), &credsmap)
	if err != nil {
		log.Fatal(err)
	}

	return credentials{
		ClientID:     credsmap["ClientID"],
		ClientSecret: credsmap["ClientSecret"],
		RestService:  credsmap["RestService"],
		LocalService: credsmap["Local"]}
}
