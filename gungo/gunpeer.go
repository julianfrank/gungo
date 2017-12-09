package gungo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxID int

//Message Structur to hold each Message between Peers
type Message []byte

//String Convert Message to String when Needed
func (message *Message) String() string {
	out, err := json.Marshal(message)
	if err != nil {
		return gunErr("gunpeer.go::Message.String::json.Marshal(message) Error: %s", err.Error())
	}
	return string(out)
}

//GunPeer Structure to Hold Each Gun Peer
type GunPeer struct {
	id     int
	ws     *websocket.Conn
	server *GunServer
	ch     chan *Message
	doneCh chan bool
}

//NewGunPeer Call This to Create new GunPeer
func NewGunPeer(ws *websocket.Conn, gunServer *GunServer) *GunPeer {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if gunServer == nil {
		panic("gunServer cannot be nil")
	}

	maxID++
	ch := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)

	return &GunPeer{maxID, ws, gunServer, ch, doneCh}
}

//Conn Method for GunPeer
func (gunPeer *GunPeer) Conn() *websocket.Conn {
	return gunPeer.ws
}

//Write Method for GunPeer
func (gunPeer *GunPeer) Write(message *Message) {
	select {
	case gunPeer.ch <- message:
	default:
		gunPeer.server.Del(gunPeer)
		err := fmt.Errorf("gunPeer %d is disconnected", gunPeer.id)
		gunPeer.server.Err(err)
	}
}

//Done gunPeer.Done()
func (gunPeer *GunPeer) Done() {
	gunPeer.doneCh <- true
}

// Listen Write and Read request via chanel
func (gunPeer *GunPeer) Listen() {
	go gunPeer.listenWrite()
	gunPeer.listenRead()
}

// Listen write request via chanel
// [TODO] This needs Work
func (gunPeer *GunPeer) listenWrite() {
	log.Println("Listening write to GunPeer")
	for {
		select {

		// send message to the GunPeer
		case message := <-gunPeer.ch:
			log.Println("Send:", message)
			websocket.JSON.Send(gunPeer.ws, message)

		// receive done request
		case <-gunPeer.doneCh:
			gunPeer.server.Del(gunPeer)
			gunPeer.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (gunPeer *GunPeer) listenRead() {
	log.Println("Listening read from GunPeer")
	for {
		select {

		// receive done request
		case <-gunPeer.doneCh:
			gunPeer.server.Del(gunPeer)
			gunPeer.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var message Message
			err := websocket.Message.Receive(gunPeer.ws, &message)
			if err == io.EOF {
				gunPeer.doneCh <- true
			} else if err != nil {
				gunPeer.server.Err(err)
			} else {
				gunPeer.server.SendAll(&message)
			}
		}
	}
}
