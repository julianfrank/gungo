package main

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxID int

//GunPeer Structure to Hold Each Gun Peer
type GunPeer struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *Message
	doneCh chan bool
}

//NewGunPeer Call This to Create new GunPeer
func NewGunPeer(ws *websocket.Conn, server *Server) *GunPeer {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxID++
	ch := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)

	return &GunPeer{maxID, ws, server, ch, doneCh}
}

func (c *GunPeer) Conn() *websocket.Conn {
	return c.ws
}

func (c *GunPeer) Write(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("GunPeer %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

func (c *GunPeer) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *GunPeer) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *GunPeer) listenWrite() {
	log.Println("Listening write to GunPeer")
	for {
		select {

		// send message to the GunPeer
		case msg := <-c.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *GunPeer) listenRead() {
	log.Println("Listening read from GunPeer")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.server.SendAll(&msg)
			}
		}
	}
}
