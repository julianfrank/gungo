package gungo

import (
	"fmt"
	"io"
	"log"
	"net/url"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxID int

//GunPeer Structure to Hold Each Gun Peer
type GunPeer struct {
	id     int
	ws     *websocket.Conn
	server *GunServer
	ch     chan *[]byte
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
	ch := make(chan *[]byte, channelBufSize)
	doneCh := make(chan bool)

	return &GunPeer{maxID, ws, gunServer, ch, doneCh}
}

//NewWS Create New WebSocket Client Directly
func NewWS(peerURL url.URL, origin url.URL) (*websocket.Conn, error) {
	gunLog("gungo.go::GunPeer.Open\tpeerURL: %s\torigin: %s", peerURL.String(), origin.String())
	ws, err := websocket.Dial(peerURL.String(), "", origin.String())
	if err != nil {
		gunErr("gungo.go::GunPeer.Open::websocket.Dial(peerURL.String(),... Error:%s", err)
		return nil, err
	}
	defer ws.Close()
	return ws, nil
}

//Conn Method for GunPeer
func (gunPeer *GunPeer) Conn() *websocket.Conn {
	return gunPeer.ws
}

//Write Method for GunPeer
func (gunPeer *GunPeer) Write(message *[]byte) {
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
	go gunPeer.listenRead()

}

// Listen write request via chanel
// [TODO] This needs Work
func (gunPeer *GunPeer) listenWrite() {
	log.Println("Listening write to GunPeer")
	for {
		select {

		// send message to the GunPeer
		case message := <-gunPeer.ch:
			if message != nil {
				log.Println("Send:", string(*message))
				websocket.Message.Send(gunPeer.ws, *message)
			}
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
			var message []byte
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
