package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/julianfrank/gungo/gungo"
)

func main() {

	// websocket server
	server := gungo.NewGunServer("/gun")
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	opts := make(map[string]interface{})
	opts["debug"] = "true"
	peerURL := url.URL{Scheme: "wss", Host: "gunjs.herokuapp.com", Path: "/gun"}
	originURL := url.URL{Scheme: "http", Host: "localhost", Path: "/"}
	ws, _ := gungo.NewWS(peerURL, originURL)
	c := gungo.NewGunPeer(ws, server)
	server.Add(c)
	c.Listen()

	for index := 0; index < 4; index++ {
		time.Sleep(time.Second)
		msg := []byte(`network`)
		c.Write(&msg)
	}

	log.Fatal(http.ListenAndServe(":7777", nil))
}
