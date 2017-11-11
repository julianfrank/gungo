package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type test struct {
	Item1 interface{}
}

func main() {

	fmt.Print("fmt")

	log.Print("log")

	x := test{Item1: "gsdfgs"}

	log.Print(x)

	m := make(map[interface{}]interface{})

	m["_"] = "hello"

	log.Print(m)

	myjson, err := json.Marshal(m)

	log.Print(myjson, err)
}
