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
	processlock *sync.RWMutex
	sync.RWMutex
	DB       *db.TransactionLayer
	Tasks    int64
	MaxTasks int
	TTL      time.Duration
	Peers    *peers.Peers
	peer.Peer
	peerCache            *cache2go.CacheTable
	LocalIP              string
	lastroutetableupdate time.Time
}

func NewNode(extip, locip string, description info.Info, layer *db.TransactionLayer, ttl time.Duration, maxtasks int) *Node {
	return &Node{Peers: peers.NewPeers(ttl), Peer: *peer.ThisPeer(extip, description), LocalIP: locip, lastroutetableupdate: time.Now(), Tasks: 0, DB: layer, TTL: ttl, peerCache: cache2go.Cache("PeerCache"), MaxTasks: maxtasks, processlock: &sync.RWMutex{}}
}

func (self *Node) GreetPeer(addr string) error {
	self.peerCache.Add(addr, 2*self.TTL, nil)
	newpeer, err := peer.NewPeer(self.Addr, addr)
	if err != nil {
		return err
	}
	self.Peers.AddPeer(newpeer)
	self.Peers.Clean(self.Addr)
	newpeernodeaddrs, err := newpeer.GetPeers(12)
	if err == nil {
		for _, newpeernodeaddr := range newpeernodeaddrs {
			if !self.Peers.Exists(newpeernodeaddr) && newpeernodeaddr != self.Addr {
				self.GreetPeer(newpeernodeaddr)
			}

		}
	}
	return nil
}

func (self *Node) GetPeers(r *http.Request, x *int, peerList *[]string) error {
	temp := self.Peers.FirstX(*x)
	peers := make([]string, 0, cap(temp))
	if len(temp) > 0 {
		for _, peernode := range temp {
			if peernode != nil {
				peers = append(peers, peernode.Addr)
			}
		}
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
	log.Printf("Added %v to Processing Queue", task.Name)
	self.processlock.Lock()
	defer self.processlock.Unlock()
	log.Printf("Processing %v", task.Name)
	return exec.Command(fmt.Sprintf("%v", task.FileName)).Output()
}

func (self *Node) NewTask(task tasks.Task) error {
	log.Println("Added to queue")
	go func() {
		for true {
			peernode, err := self.Peers.GetAPeer()
			if err == nil {
				result, err := peernode.AllocateTask(&task)
				if err == nil {
					log.Println("Allocated", err, string(result))
					return
				}
			}
		}
	}()
	return nil
}

func (self *Node) AllocateTask(r *http.Request, task *tasks.Task, result *[]byte) error {
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
	for true {
		if !self.Compute || int(self.Tasks+1) > self.MaxTasks {
			peernode, err := self.Peers.GetAPeer()
			if err == nil {
				log.Printf("Allocated from %v to %v", self.Addr, peernode.Addr)
				*result, err = peernode.AllocateTask(task)
				if err == nil {
					log.Printf("Recieved from %v", peernode.Addr)
					return nil
				}
			}
		} else {
			atomic.AddInt64(&self.Tasks, 1)
			data, err := self.process(task)
			*result = data
			atomic.AddInt64(&self.Tasks, -1)
			if err == nil {
				log.Println("Successful Process")
				return err
			}
			log.Println(err)
		}
	}
	return nil
}
