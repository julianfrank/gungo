package main

import (
	"net/url"
	"os"
	"os/signal"
)

//Gun Structure containing the key object
type Gun struct {
	DB   map[string]interface{} //Hold the Database in this map
	peer GunPeer                //Hold the Peer

	interrupt chan os.Signal
	done      chan struct{}
}

//Init Initialize the DB Manager
func (gun Gun) Init(opts map[string]interface{}) {
	//startTime := time.Now()
	if opts["debug"] == "true" {
		gunDebug = true
		gunLog("gungo.go::Gun.Init opts:%v", opts)
	} else {
		gunDebug = false
	}

	gun.interrupt = make(chan os.Signal, 1)
	signal.Notify(gun.interrupt, os.Interrupt)

	gun.done = make(chan struct{})
	defer close(gun.done)

	gun.DB = make(map[string]interface{})

	if opts["peerURL"] != nil {
		gun.peer.Open(opts["peerURL"].(url.URL), opts["origin"].(url.URL))
	}
}

/*




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
