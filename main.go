package main

import (
	"log"
)

var (
	This *node.Node
)

func main() {
	//initForward()
	initDHT()
	log.Println("whee")
	select {}
}
