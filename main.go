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
	_ "net/http/pprof"
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
	This = node.NewNode(fmt.Sprintf("%v:%v", extip.String(), Config.Port), fmt.Sprintf("%v:%v", locip.String(), Config.Port), *description, 20*time.Second, Config.MaxTasks)
}

func startrpc() {
	server := rpc.NewServer()
	server.Register(This)
	l, e := net.Listen("tcp", fmt.Sprintf(":%v", Config.Port))
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
		go func() {
			if This.Peers.Length() == 0 {
				if Config.PeerSeed != "" {
					This.Peers.AddPeer(Config.PeerSeed)
				}
			} else {
				This.Peers.Update()
			}
		}()
	}

	select {}
}
