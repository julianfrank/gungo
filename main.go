package main

import (
	"net/url"
)

func main() {
	var myDB Gun

	opts := make(map[string]interface{})
	opts["debug"] = "true"
	opts["peerURL"] = url.URL{Scheme: "wss", Host: "gunjs.herokuapp.com", Path: "/gun"}
	myDB.Init(opts)
}
