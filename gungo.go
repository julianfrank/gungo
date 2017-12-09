package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/websocket"
)

var (
	gunDebug bool
)

func gunLog(msgFormat string, msg interface{}) {
	if gunDebug {
		fmt.Printf(msgFormat+"\n", msg)
	}
}
func gunErr(msgFormat string, msg interface{}) {
	log.Printf(msgFormat, msg)
}
func gunTimed(msgFormat string, msg interface{}, startTime time.Time) {
	if gunDebug {
		log.Printf(msgFormat+"\t%s", msg, time.Since(startTime))
	}
}

//Gun Structure containing the key object
type Gun struct {
	DB   map[string]interface{} //Hold the Database in this map
	peer GunPeer                //Hold the Peer

	interrupt chan os.Signal
	done      chan struct{}
}

//GunPeer Structure to Hold the Connection with Peer
type GunPeer struct {
	url       url.URL         //URL Object with Scheme, Host and Path
	Wire      *websocket.Conn //Hold the Connections in this Connection
	Connected bool            //Connection status
}

//Open Open a New Peer
func (gunPeer GunPeer) Open(peerURL url.URL) {
	//startTime := time.Now()
	if gunPeer.Connected == true {
		gunErr("gungo.go::GunPeer.Open Error: Peer Already Exists -> url:%s", gunPeer.url.String())
		panic("")
	} else {
		var err error
		var peerStr = peerURL.String()
		var origin = "http://localhost/"
		gunLog("gungo.go::GunPeer.Open peerURL:%s", peerStr)
		peerStr = "wss://gunjs.herokuapp.com/gun"
		gunPeer.Wire, err = websocket.Dial(peerStr, "", origin)
		if err != nil {
			gunErr("gungo.go::GunPeer.Open::websocket.Dial(peerURL.String(),... Error:%s", err)
		}
		defer gunPeer.Wire.Close()
	}
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
		gun.peer.Open(opts["peerURL"].(url.URL))
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
