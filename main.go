package main

import (
	"log"
)

func init() {
	initForward()
	initDHT()
	log.Println("initialized")
}

func main() {
	//aliveLoop()
	log.Println("whee")
	select {}
}
