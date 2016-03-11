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

type MyServer struct {
	r *mux.Router
}

func addDefaultHeadersFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		fn(w, r)
	}
}

func addDefaultHeadersHand(han http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		han.ServeHTTP(w, r)
	}
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
	r.HandleFunc("/rpc", addDefaultHeadersHand(s))
	r.HandleFunc("/api/peers", api_peers)
	r.HandleFunc("/api/task", api_task)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	log.Println("whee")
	go http.ListenAndServe(fmt.Sprintf(":%v", Config.Mappings["RPC"].Port), r)
	log.Println(s.HasMethod("Node.AllocateTask"))
	if Config.PeerSeed != "" {
		This.GreetPeer(Config.PeerSeed)
	}
	for {
		time.Sleep(2500 * time.Millisecond)
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
