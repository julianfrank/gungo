package gungo

import (
	"log"

	"golang.org/x/net/websocket"
)

//Try sdlkal;kd
func Try() {

	var origin = "http://localhost/"
	//var url = "ws://gunjs.herokuapp.com/gun"
	var url = "ws://localhost:7777"

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte("hi")
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	gunLog("Send: %s", message)

	var msg = make([]byte, 16)
	_, err = ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	gunLog("Receive: %s", msg)
	_, err = ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	gunLog("Receive: %s\n", msg)
}
