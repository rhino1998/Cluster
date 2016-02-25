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
	"github.com/rhino1998/cluster/peer"
	"github.com/rhino1998/cluster/util"
	"github.com/rhino1998/god/dhash"
	"log"
	"net/http"
	"time"
)

var (
	This *node.Node
)

func init_node() {
	log.Println(Config.Mappings["DHT"].Port)
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
	kvstore := dhash.NewNodeDir(fmt.Sprintf("%v:%v", "0.0.0.0", Config.Mappings["DHT"].Port), fmt.Sprintf("%v:%v", extip, Config.Mappings["DHT"].Port), "")
	kvstore.Start()
	layer := db.NewTransactionLayer(kvstore)
	This = node.NewNode(fmt.Sprintf("%v:%v", extip.String(), Config.Mappings["DHT"].Port), fmt.Sprintf("%v:%v", locip.String(), Config.Mappings["DHT"].Port), *description, layer)
	if Config.DHTSeed != "" {
		This.DB.DB.MustJoin(Config.DHTSeed)
	}
}

func print() {
	for {
		time.Sleep(3 * time.Second)
		vals, _ := This.Peers.After(time.Now().UTC().Add(-500 * time.Minute))
		if len(vals) > 0 {
			log.Println("peers")
			for _, peernode := range vals {
				log.Println(peernode.Addr)
			}
		}
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
	go http.ListenAndServe(fmt.Sprintf(":%v", Config.Mappings["RPC"].Port), r)
	log.Println(s.HasMethod("Node.RouteTask"))
	if Config.PeerSeed != "" {
		log.Println("sup")
		newpeer, err := peer.NewPeer(This.Addr, This.Info, Config.PeerSeed)
		log.Println(newpeer)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		This.Peers.AddPeer(newpeer)
	}
	go print()
	select {}
}
