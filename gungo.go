package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

//Gun Structure containing the key object
type Gun struct {
	DB    map[string]string //Hold the Database in this map
	peers []struct {
		URLPath string         //URL Path (url + Path including ws/wss)
		Wire    websocket.Conn //Hold the Connections in this Connection Slice
		/*		Open      interface{}    //Open Method
				OnOpen    interface{}
				Close     interface{}
				onClose   interface{}
				Send      interface{}
				onMessage interface{}
				onError   interface{}*/
	}
}

//Init Initialize the DB Manager
func (gun Gun) Init(opts map[string]string) {
	log.Println(opts)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	done := make(chan struct{})

	defer close(done)

}

/*
	u := url.URL{Scheme: "wss", Host: opts["peerURL"], Path: opts["gunPath"]}
	log.Printf("gungo Going to work with \tpeerURL:%s\tgunPath:%s\tusing String:%s", opts["peerURL"], opts["gunPath"], u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial err:", err)
	}
	//defer c.Close()

	go func() {
		//defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read err:", err)
				return
			}
			log.Printf("recv: %s %T", string(message), message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write err:", err)
				return
			}
			log.Println("write: ", t.String())
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			//c.Close()
			return
		}
	}

*/
