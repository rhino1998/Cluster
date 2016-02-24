package node

import (
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/peers"
	"net/http"
	"sync"
	"time"
)

type Node struct {
	sync.RWMutex
	Peers                *peers.Peers
	Addr                 string `json:"addr"`
	LocalIP              string
	lastroutetableupdate time.Time
	info.Info
}

func NewNode() *Node {
	return &Node{Peers: peers.NewPeers()}
}

func (self *Node) GetPeers(r *http.Request, start time.Time, peerList []string) error {
	temp, err := self.Peers.After(start)
	peerList = make([]string, len(temp))
	for i, peer := range temp {
		peerList[i] = peer.Addr
	}
	return err
}

func (self *Node) Ping(start time.Time, end *time.Time) error {
	*end = time.Now()
	return nil
}

func (self *Node) Describe(n bool, desciption *info.Info) error {
	desciption = &info.Info{Specs: self.Specs, Compute: self.Compute}
	return nil
}
