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

func gunLog(msgFormat string, msg ...interface{}) {
	if gunDebug {
		fmt.Printf(msgFormat+"\n", msg...)
	}
}
func gunErr(msgFormat string, msg ...interface{}) {
	log.Printf(msgFormat, msg...)
}
func gunTimed(msgFormat string, startTime time.Time, msg ...interface{}) {
	if gunDebug {
		msg = append(msg, time.Since(startTime))
		log.Printf(msgFormat+"\t%s", msg...)
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
func (gunPeer GunPeer) Open(peerURL url.URL, origin url.URL) {
	startTime := time.Now()
	if gunPeer.Connected == true {
		gunErr("gungo.go::GunPeer.Open Error: Peer Already Exists -> url:%s", gunPeer.url.String())
	} else {
		var err error
		gunLog("gungo.go::GunPeer.Open\tpeerURL: %s\torigin: %s", peerURL.String(), origin.String())
		gunPeer.Wire, err = websocket.Dial(peerURL.String(), "", origin.String())
		if err != nil {
			gunErr("gungo.go::GunPeer.Open::websocket.Dial(peerURL.String(),... Error:%s", err)
		}
		defer gunPeer.Wire.Close()
		gunTimed("gungo.go::GunPeer.Open Success", startTime)

		for {
			msg := make([]byte, 512)
			n, err := gunPeer.Wire.Read(msg)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Receive: %s\n", msg[:n])
		}
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
