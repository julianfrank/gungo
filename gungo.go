package main

import (
	"encoding/json"

	"log"
)

//GunGraph DB Instance
type GunGraph struct {
	Graph map[string]interface{}
}

func main() {

	//Instantiating the DB here
	var db GunGraph
	db.Graph = make(map[string]interface{})

	//Creating the First Node here
	n1ID := make(map[string]string)
	n1ID["#"] = "ASDF"
	n1Node := make(map[string]interface{})
	n1Node["_"] = n1ID
	n1Node["Name"] = "JF"
	//Creating the Second Node here
	n2ID := make(map[string]string)
	n2ID["#"] = "FDSA"
	n2Node := make(map[string]interface{})
	n2Node["_"] = n2ID
	n2Node["Name"] = "Fluffy"
	n2Node["Species"] = "Cat"
	//Connecting both nodes here
	n1Node["boss"] = n2ID
	n2Node["slave"] = n1ID
	//Adding the nodes to the DB
	db.Graph["ASDF"] = n1Node
	db.Graph["FDSA"] = n2Node

	//Trying out Traversals
	log.Print("\n############### Traversal before conversion to JSON ############")
	log.Print("db.Graph[ASDF]\t", db.Graph["ASDF"])
	t := db.Graph["ASDF"].(map[string]interface{})
	boss := t["boss"].(map[string]string)
	log.Print("n1 boss\t", t["boss"], db.Graph[boss["#"]])

	//Converting DB to JSON for communication with Peers
	log.Print("\n###############  JSON format ############")
	dbjson, err := json.Marshal(db)
	log.Print("string(dbjson), err\t", string(dbjson), err)
	//Checking Reverse MApping back to DB from json format
	var intermediateinterface interface{}
	err = json.Unmarshal(dbjson, &intermediateinterface)
	m := intermediateinterface.(map[string]interface{})

	var tempDB GunGraph
	tempDB.Graph = m["Graph"].(map[string]interface{})
	log.Printf("tempDB Parsed from JSON\tType:%T\tValue:%v\n", tempDB, tempDB)

	log.Print("\n############### Traversal after parsing from JSON ############")
	//Trying out Traversals
	log.Print("tempDB.Graph[ASDF]\t", tempDB.Graph["ASDF"])
	t = tempDB.Graph["ASDF"].(map[string]interface{})
	newboss := t["boss"].(map[string]interface{})
	newbossnode := newboss["#"].(string)
	log.Print("n1 boss\t", t["boss"], tempDB.Graph[newbossnode])
}
