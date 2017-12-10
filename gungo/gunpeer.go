package gungo

import (
	"golang.org/x/net/websocket"
)

//Try sdlkal;kd
func Try() {
	gunLog("hi")
	ws, err := websocket.Dial("ws://gunjs.herokuapp.com/gun", "", "localhost")
	gunLog("%s %s", ws, err)
}
