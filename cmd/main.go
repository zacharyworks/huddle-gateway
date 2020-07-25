package main

import (
	"encoding/json"
	"flag"
	"github.com/zacharyworks/huddle-gateway/auth"
	"github.com/zacharyworks/huddle-gateway/wsockets"
	"github.com/zacharyworks/huddle-shared/db"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8000", "http service address")

func main() {
	flag.Parse()
	// Create web socket service
	hub := wsockets.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsockets.ServeWs(hub, w, r)
	})

	//Create auth service
	oAuth := auth.NewAuth(readAuthCredentials())
	http.HandleFunc("/login", oAuth.HandleLogin)
	http.HandleFunc("/callback", oAuth.HandleCallback)

	// Start database
	db.ConnectDB()
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}

func readAuthCredentials() (clientID string, clientSecret string) {
	file, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatal(err)
	}
	credentials := make(map[string]string)
	err = json.Unmarshal([]byte(file), &credentials)
	if err != nil {
		log.Fatal(err)
	}
	return credentials["ClientID"], credentials["ClientSecret"]
}
