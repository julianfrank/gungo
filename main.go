package main

import (
	"log"
	"net/http"

	"github.com/julianfrank/gungo/gungo"
)

func main() {
	// websocket server
	server := gungo.NewGunServer("/gun")
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":7777", nil))
}

/*
var myDB Gun
opts := make(map[string]interface{})
	opts["debug"] = "true"
	opts["peerURL"] = url.URL{Scheme: "wss", Host: "gunjs.herokuapp.com", Path: "/gun"}
	opts["origin"] = url.URL{Scheme: "http", Host: "localhost", Path: "/"}
	myDB.Init(opts)
*/
