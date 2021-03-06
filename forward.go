package main

import (
	"github.com/rhino1998/cluster/forward"
	"github.com/syncthing/syncthing/lib/upnp"
	"log"
	"time"
)

var (
	nat upnp.IGD
)

func initForward() {
	nat = upnp.Discover(5 * time.Second)[0]
	err := forward.Forward(nat, forward.Mapping{Ports: []int{port}, Protocols: []upnp.Protocol{upnp.TCP}, Description: "Cluster"})
	if err != nil {
		log.Printf("Port forwarding failed: %v", err)
		log.Println(`You may have to mainually port forward or choose different ports`)
	} else {
		log.Println("Port Forwading Success!")
	}
}
