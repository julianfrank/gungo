package gungo

import (
	"net/http"

	"golang.org/x/net/websocket"
)

//GunServer Gun server Structure
type GunServer struct {
	pattern   string
	messages  []*[]byte
	gunPeers  map[int]*GunPeer
	addCh     chan *GunPeer
	delCh     chan *GunPeer
	sendAllCh chan *[]byte
	doneCh    chan bool
	errCh     chan error
}

//NewGunServer Create new Gun server
func NewGunServer(pattern string) *GunServer {
	messages := make([]*[]byte, channelBufSize)
	gunPeers := make(map[int]*GunPeer)
	addCh := make(chan *GunPeer)
	delCh := make(chan *GunPeer)
	sendAllCh := make(chan *[]byte)
	doneCh := make(chan bool)
	errCh := make(chan error)
	return &GunServer{pattern, messages, gunPeers, addCh, delCh, sendAllCh, doneCh, errCh}
}

//Add GunServer.Add Method to add new GunPeer
func (gunServer *GunServer) Add(gunPeer *GunPeer) {
	gunServer.addCh <- gunPeer
}

//Del Delete GunPeer
func (gunServer *GunServer) Del(gunPeer *GunPeer) {
	gunServer.delCh <- gunPeer
}

//SendAll Flush All Messages
func (gunServer *GunServer) SendAll(message *[]byte) {
	gunServer.sendAllCh <- message
}

//Done Signal Completion
func (gunServer *GunServer) Done() {
	gunServer.doneCh <- true
}

//Err Error
func (gunServer *GunServer) Err(err error) {
	gunServer.errCh <- err
}

func (gunServer *GunServer) sendPastMessages(gunPeer *GunPeer) {
	for _, message := range gunServer.messages {
		gunPeer.Write(message)
	}
}

func (gunServer *GunServer) sendAll(message *[]byte) {
	for _, gunPeer := range gunServer.gunPeers {
		gunPeer.Write(message)
	}
}

// Listen and serve.
// It serves gunPeer connection and broadcast request.
func (gunServer *GunServer) Listen() {

	gunLog("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				gunServer.errCh <- err
			}
		}()

		gunPeer := NewGunPeer(ws, gunServer)
		gunServer.Add(gunPeer)
		gunPeer.Listen()
	}
	http.Handle(gunServer.pattern, websocket.Handler(onConnected))
	gunLog("Created handler")

	for {
		select {

		// Add new gunPeer
		case gunPeer := <-gunServer.addCh:
			gunLog("Added new gunPeer")
			gunServer.gunPeers[gunPeer.id] = gunPeer
			gunLog("Now %d gunPeers connected", len(gunServer.gunPeers))
			gunServer.sendPastMessages(gunPeer)

		// del a gunPeer
		case gunPeer := <-gunServer.delCh:
			gunLog("Delete gunPeer")
			delete(gunServer.gunPeers, gunPeer.id)

		// broadcast message for all gunPeers
		case message := <-gunServer.sendAllCh:
			gunLog("Send all:", message)
			gunServer.messages = append(gunServer.messages, message)
			gunServer.sendAll(message)

		case err := <-gunServer.errCh:
			gunLog("Error:", err.Error())

		case <-gunServer.doneCh:
			return
		}
	}
}
