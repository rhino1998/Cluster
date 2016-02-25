package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/rhino1998/cluster/bench"
	"github.com/rhino1998/cluster/db"
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/node"
	"github.com/rhino1998/cluster/util"
	"github.com/rhino1998/god/dhash"
	"log"
	"net/http"
)

var (
	This *node.Node
)

func init_node() {
	specs, err := bench.LoadSpecs("./specs.json")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	description := &info.Info{Compute: Config.Compute, Specs: *specs}
	extip, err := util.GetExternalIP()
	if err != nil {
		log.Println(err)
	}
	locip, err := util.GetLocalIP()
	if err != nil {
		log.Println(err)
	}
	kvstore := dhash.NewNodeDir(fmt.Sprintf("%v:%v", "0.0.0.0", Config.Mappings["DHT"].Port), fmt.Sprintf("%v:%v", This.Addr, Config.Mappings["DHT"].Port), "")
	layer := db.NewTransactionLayer(kvstore)
	This = node.NewNode(extip.String(), locip.String(), *description, layer)
	This.DB.DB.Start()
	if Config.DHTSeed != "" {
		This.DB.DB.MustJoin(Config.DHTSeed)
	}
}

func main() {
	//initForward()
	init_node()
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	s.RegisterService(This, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	log.Println("whee")
	http.ListenAndServe(":1234", r)
}
