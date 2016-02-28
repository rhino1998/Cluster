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
	//"github.com/rhino1998/cluster/peer"
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
	This = node.NewNode(fmt.Sprintf("%v:%v", extip.String(), Config.Mappings["RPC"].Port), fmt.Sprintf("%v:%v", locip.String(), Config.Mappings["RPC"].Port), *description, layer, 20*time.Second)
	if Config.DHTSeed != "" {
		This.DB.DB.MustJoin(Config.DHTSeed)
	}
}

func main() {
	initForward()
	init_node()
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	s.RegisterService(This, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	r.HandleFunc("/api/peers", api_peers)
	r.HandleFunc("/api/db/put", api_db_put)
	r.HandleFunc("/api/db/get", api_db_get)
	r.HandleFunc("/api/db/del", api_db_del)
	r.HandleFunc("/api/task", api_task)
	log.Println("whee")
	go http.ListenAndServe(fmt.Sprintf(":%v", Config.Mappings["RPC"].Port), r)
	log.Println(s.HasMethod("Node.RouteTask"))
	if Config.PeerSeed != "" {
		This.GreetPeer(Config.PeerSeed)
	}
	for {
		time.Sleep(5 * time.Second)
		go func() {

			peernode, err := This.Peers.GetAPeer()
			if err != nil {
				log.Println(err)
				return
			}
			This.GreetPeer(peernode.Addr)
			log.Println("sup")
		}()
	}

	select {}
}
