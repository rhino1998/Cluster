package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/rhino1998/cluster/node"
	"log"
	"net/http"
)

var (
	This *node.Node
)

func main() {
	//initForward()
	initDHT()
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	s.RegisterService(This, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	log.Println("whee")
	http.ListenAndServe(":1234", r)
}
