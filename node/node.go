package node

import (
	"github.com/muesli/cache2go"
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/peer"
	"github.com/rhino1998/cluster/peers"
	"github.com/rhino1998/cluster/tasks"
	"github.com/rhino1998/god/dhash"
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
	DB          *dhash.Node
	TaskValue   int64
	RoutedTasks int64
	MaxTasks    int
	MaxRouted   int
	TTL         time.Duration
	Peers       *peers.Peers
	peer.Peer
	peerCache            *cache2go.CacheTable
	LocalIP              string
	lastroutetableupdate time.Time
}

func NewNode(extip, locip string, description info.Info, kvstore *dhash.Node, ttl time.Duration, maxtasks int) *Node {
	return &Node{Peers: peers.NewPeers(ttl), Peer: *peer.ThisPeer(extip, description), LocalIP: locip, lastroutetableupdate: time.Now(), TaskValue: 0, DB: kvstore, TTL: ttl, peerCache: cache2go.Cache("PeerCache"), MaxTasks: maxtasks, processlock: &sync.RWMutex{}, RoutedTasks: 0}
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
	log.Printf("Added %v to Processing Queue", string(task.Id))
	atomic.AddInt64(&self.TaskValue, int64(task.Value))
	self.processlock.Lock()
	defer self.processlock.Unlock()
	defer atomic.AddInt64(&self.TaskValue, int64(-task.Value))
	log.Printf("Processing %v", task.Id, self.TaskValue)
	return exec.Command(fmt.Sprintf("%v", task.FileName)).Output()
}

func (self *Node) NewTask(task tasks.Task) error {
	log.Printf("Added %v to queue", string(task.Id))
	go func() {
		temp := time.Now()
		for true {
			peernode, err := self.Peers.GetAPeer()
			if err == nil {
				log.Println(string(task.Id))
				result, err := peernode.AllocateTask(&task)
				if err == nil {
					log.Println("Allocated", err, string(result), string(task.Id), time.Now().Sub(temp).Seconds())
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
		if !self.Compute || int(self.TaskValue+1) > 10000 {
			peernode, err := self.Peers.GetAPeer()
			if err == nil {
				log.Printf("Allocated %v from %v to %v", string(task.Id), self.Addr, peernode.Addr)
				if int(self.RoutedTasks+1) > self.MaxTasks {
					self.processlock.Lock()
				}
				atomic.AddInt64(&self.RoutedTasks, 1)
				*result, err = peernode.AllocateTask(task)
				atomic.AddInt64(&self.RoutedTasks, -1)
				if int(self.RoutedTasks+1) > self.MaxTasks {
					self.processlock.Unlock()
				}
				if err == nil {
					log.Printf("Recieved %v from %v", string(task.Id), peernode.Addr)
					return nil
				}
			}
		} else {
			data, err := self.process(task)
			*result = data
			if err == nil {
				log.Println("Successful Process")
				return err
			}
			log.Println(err)
		}
	}
	*result = []byte("Error")
	return nil
}
