package main

import (
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

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

		//-------------------Loop in future----------------
		message := []byte("hi")
		gunPeer.Send(message)
		gunPeer.read()
		//-----------------------------------------------
		gunTimed("gungo.go::GunPeer.Open Success", startTime)
	}
}

//Send Perform Write on the websocket
func (gunPeer GunPeer) Send(message []byte) error {
	startTime := time.Now()
	_, err := gunPeer.Wire.Write(message)
	if err != nil {
		gunErr("gungo.go::GunPeer.Send::gunPeer.Wire.Write(message) Error:%s", err)
	}
	gunTimed("gungo.go::GunPeer.Send Success", startTime)
	return nil
}

//read Perform read on the websocket
func (gunPeer GunPeer) read() ([]byte, error) {
	startTime := time.Now()
	var msg = make([]byte, 512)
	_, err := gunPeer.Wire.Read(msg)
	if err != nil {
		gunErr("gungo.go::GunPeer.read::gunPeer.Wire.Read(msg) Error:%s", err)
		return msg, err
	}
	gunTimed("gungo.go::GunPeer.read msg: %s", startTime, msg)
	return msg, nil
}
