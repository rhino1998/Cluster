package main

import (
	"fmt"
	"github.com/rhino1998/cluster/bench"
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/node"
	"net"
	"net/rpc"
	//"github.com/rhino1998/cluster/peer"
	"flag"
	"github.com/rhino1998/cluster/util"
	"log"
	//"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

var (
	This *node.Node
	port int
)

func init_node() {
	log.Println(port)
	specs, err := bench.LoadSpecs("./specs.json")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	description := &info.Info{Compute: Config.Compute, Specs: *specs}
	extip, err := util.GetExternalIP()
	log.Println(extip, port)
	if err != nil {
		log.Println(err)
	}
	locip, err := util.GetLocalIP()
	if err != nil {
		log.Println(err)
	}
	This = node.NewNode(fmt.Sprintf("%v:%v", extip.String(), port), fmt.Sprintf("%v:%v", locip.String(), port), *description, 20*time.Second, Config.MaxTasks)
}

func startrpc() {
	server := rpc.NewServer()
	server.Register(This)
	l, e := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go server.ServeConn(conn)
	}
}

func main() {
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
	flag.IntVar(&port, "port", Config.Port, "port")
	forwardport := flag.Bool("f", false, "forward or not")
	flag.Parse()
	runtime.GOMAXPROCS(1)
	log.Println(*forwardport)
	if *forwardport {
		initForward()
	}
	init_node()
	go startrpc()
	log.Println("whee")
	if Config.PeerSeed != "" {
		This.Peers.AddPeer(Config.PeerSeed)
	}
	for {
		time.Sleep(2500 * time.Millisecond)
		//if This.Peers.Length() == 0 {
		/*if Config.PeerSeed != "" {
			This.Peers.AddPeer(Config.PeerSeed)
		}
		//}*/
		//log.Println(This.Peers.Length())
		go This.Peers.Update()
	}

	select {}
}
