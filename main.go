package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/rhino1998/cluster/bench"
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/node"
	//"github.com/rhino1998/cluster/peer"
	"flag"
	"github.com/rhino1998/cluster/util"
	"github.com/rhino1998/god/dhash"
	"log"
	"net/http"
	"runtime"
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
	kvstoreaddr := fmt.Sprintf("%v:%v", extip, Config.Mappings["DHT"].Port)
	kvstore := dhash.NewNodeDir(fmt.Sprintf("%v:%v", "0.0.0.0", Config.Mappings["DHT"].Port), fmt.Sprintf("%v:%v", extip, Config.Mappings["DHT"].Port), "")
	kvstore.Start()
	kvstore.StartJson()
	if Config.DHTSeed != "" {
		kvstore.MustJoin(Config.DHTSeed)
	}
	This = node.NewNode(fmt.Sprintf("%v:%v", extip.String(), Config.Mappings["RPC"].Port), fmt.Sprintf("%v:%v", locip.String(), Config.Mappings["RPC"].Port), *description, kvstoreaddr, 20*time.Second, Config.MaxTasks)
}

func main() {
	forwardport := flag.Bool("f", false, "forward or not")
	flag.Parse()
	runtime.GOMAXPROCS(1)
	log.Println(*forwardport)
	log.Println(flag.Args())
	if *forwardport {
		initForward()
	}
	init_node()
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	s.RegisterService(This, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	r.HandleFunc("/api/peers", api_peers)
	r.HandleFunc("/api/task", api_task)
	log.Println("whee")
	go http.ListenAndServe(fmt.Sprintf(":%v", Config.Mappings["RPC"].Port), r)
	log.Println(s.HasMethod("Node.AllocateTask"))
	if Config.PeerSeed != "" {
		This.GreetPeer(Config.PeerSeed)
	}
	for {
		time.Sleep(5 * time.Second)
		go func() {
			peernode, err := This.Peers.GetAPeer()
			if err != nil {
				log.Println(err)
				if Config.PeerSeed != "" {
					This.GreetPeer(Config.PeerSeed)
				} else {
					log.Println("waiting for peers")
				}
			} else {
				This.GreetPeer(peernode.Addr)
			}
		}()
	}

	select {}
}
