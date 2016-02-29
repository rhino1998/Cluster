package node

import (
	"github.com/muesli/cache2go"
	"github.com/rhino1998/cluster/db"
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/peer"
	"github.com/rhino1998/cluster/peers"
	"github.com/rhino1998/cluster/tasks"
	"log"
	"net/http"
	"os/exec"
	//"reflect"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	sync.RWMutex
	DB    *db.TransactionLayer
	Tasks int64
	TTL   time.Duration
	Peers *peers.Peers
	peer.Peer
	peerCache            *cache2go.CacheTable
	LocalIP              string
	lastroutetableupdate time.Time
}

func NewNode(extip, locip string, description info.Info, layer *db.TransactionLayer, ttl time.Duration) *Node {
	return &Node{Peers: peers.NewPeers(ttl), Peer: *peer.ThisPeer(extip, description), LocalIP: locip, lastroutetableupdate: time.Now(), Tasks: 0, DB: layer, TTL: ttl, peerCache: cache2go.Cache("PeerCache")}
}

func (self *Node) GreetPeer(addr string) error {
	self.peerCache.Add(addr, 2*self.TTL, nil)
	newpeer, err := peer.NewPeer(self.Addr, addr)
	if err != nil {
		return err
	}
	self.Peers.AddPeer(newpeer)
	self.Peers.Clean(self.Addr)
	newpeernodeaddrs, err := newpeer.GetPeers(self.TTL)
	if err == nil {
		for _, newpeernodeaddr := range newpeernodeaddrs {
			if !self.Peers.Exists(newpeernodeaddr) && newpeernodeaddr != self.Addr {
				self.GreetPeer(newpeernodeaddr)
			}

		}
	}
	return nil
}

func (self *Node) GetPeers(r *http.Request, start *time.Duration, peerList *[]string) error {
	temp := self.Peers.Items()
	peers := make([]string, len(temp), cap(temp))
	for i, peer := range temp {
		peers[i] = peer.Addr
	}
	*peerList = peers
	return nil
}

func (self *Node) Ping(start time.Time, end *time.Time) error {
	*end = time.Now()
	return nil
}

func (self *Node) Greet(r *http.Request, remaddr *string, desciption *info.Info) error {
	*desciption = info.Info{Compute: self.Compute, Specs: self.Specs}
	if !self.peerCache.Exists(*remaddr) && !self.Peers.Exists(*remaddr) {
		self.GreetPeer(*remaddr)
	}
	return nil
}

func (self *Node) process(task *tasks.Task) ([]byte, error) {
	log.Println("Processing")
	return exec.Command(fmt.Sprintf("%v", task.FileName)).Output()
}

func (self *Node) NewTask(task tasks.Task) error {
	peernode, err := self.Peers.GetAPeer()
	if err != nil {
		return err
	}
	log.Println("Allocate Init")
	go func(peernode *peer.Peer) {
		log.Println("Allocate Init2")
		result, err := peernode.AllocateTask(&task)
		log.Println("Allocated", err, string(result))
	}(peernode)
	return nil
}

func (self *Node) AllocateTask(r *http.Request, task *tasks.Task, result *[]byte) error {
	log.Println("allocinit4")
	// /task.Jumps[self.Addr] = len(task.Jumps) + 1
	/*for _, req := range task.Reqs {
		if ok, err := req.Comp(req.Value(), reflect.ValueOf(self).FieldByName(req.Name())); !ok || !self.Compute {
			if err != nil {
				return err
			}
			peernode, err := self.Peers.BestMatch(task.Reqs)
			if err != nil {
				return err
			}
			result, err = peernode.AllocateTask(task)
			if err != nil {
				return err
			}
		}
	}*/
	if !self.Compute || self.Tasks+1 > 1 {
		peernode, err := self.Peers.GetAPeer()
		if err != nil {
			return err
		}
		*result, err = peernode.AllocateTask(task)
		if err != nil {
			return err
		}
	}
	atomic.AddInt64(&self.Tasks, 1)
	data, err := self.process(task)
	log.Println(string(data), err)
	*result = data
	atomic.AddInt64(&self.Tasks, -1)
	if err != nil {
		return err
	}
	return nil

}
