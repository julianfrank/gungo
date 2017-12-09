package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (self *Message) String() string {
	return self.Author + " says " + self.Body
}

// Chat server.
type Server struct {
	pattern   string
	messages  []*Message
	gunPeers  map[int]*GunPeer
	addCh     chan *GunPeer
	delCh     chan *GunPeer
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

// Create new chat server.
func NewServer(pattern string) *Server {
	messages := []*Message{}
	gunPeers := make(map[int]*GunPeer)
	addCh := make(chan *GunPeer)
	delCh := make(chan *GunPeer)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		messages,
		gunPeers,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s *Server) Add(c *GunPeer) {
	s.addCh <- c
}

func (s *Server) Del(c *GunPeer) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *GunPeer) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}

func (s *Server) sendAll(msg *Message) {
	for _, c := range s.gunPeers {
		c.Write(msg)
	}
}

// Listen and serve.
// It serves gunPeer connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		gunPeer := NewGunPeer(ws, s)
		s.Add(gunPeer)
		gunPeer.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a gunPeer
		case c := <-s.addCh:
			log.Println("Added new gunPeer")
			s.gunPeers[c.id] = c
			log.Println("Now", len(s.gunPeers), "gunPeers connected.")
			s.sendPastMessages(c)

		// del a gunPeer
		case c := <-s.delCh:
			log.Println("Delete gunPeer")
			delete(s.gunPeers, c.id)

		// broadcast message for all gunPeers
		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
