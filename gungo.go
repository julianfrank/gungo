package main

import (
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
}
