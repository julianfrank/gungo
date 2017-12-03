package main

func main() {
	var myDB Gun

	opts := make(map[string]string)
	opts["peerURL"] = "gunjs.herokuapp.com"
	opts["gunPath"] = "/gun"
	myDB.Init(opts)
}
