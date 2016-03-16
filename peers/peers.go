package peers

import (
	"log"
	//"net"
	"github.com/rhino1998/cluster/tasks"
	"github.com/rhino1998/cluster/util"
	"github.com/spaolacci/murmur3"
	"sync"
)

//A time indexed queryable map of peers
type Peers struct {
	sync.RWMutex
	history map[string]uint64
	peers   chan *Peer
	route   chan *Peer
	addr    string
}

//Creates an empty list of peers
func NewPeers(addr string) *Peers {
	return &Peers{peers: make(chan *Peer, 20), route: make(chan *Peer, 2000), addr: addr, history: make(map[string]uint64, 20)}
}

func (self *Peers) Length() int {
	return len(self.peers)
}

func (self *Peers) Update() {
	peernode := <-self.route
	self.route <- peernode
}

func (self *Peers) GetXPeers(x int) []*Peer {
	temp := make([]*Peer, x, x)
	for i := 0; i < x; i++ {
		select {
		case peernode := <-self.route:
			temp = append(temp, peernode)
			self.route <- peernode
		default:
		}
	}
	return temp
}

func (self *Peers) ClosestPeer(key []byte) *Peer {
	hash := murmur3.Sum64(key)
	closest := ^uint64(0)
	peer := <-self.route
	self.route <- peer
	for i := 0; i < len(self.route); i++ {
		peernode := <-self.route
		self.route <- peernode
		val := hash ^ self.history[peernode.Addr]
		if val < closest {
			closest = val
			peer = peernode
		}
	}
	return peer
}

func (self *Peers) GetAPeer() *Peer {
	peernode := <-self.peers
	log.Println(peernode.IsDead())
	for peernode.IsDead() {
		self.Lock()
		delete(self.history, peernode.Addr)
		self.Unlock()
		peernode = <-self.peers
	}
	self.peers <- peernode
	return peernode
}

func (self *Peers) Reconnect(remaddr string) {
	NewPeer(self.addr, remaddr)
}

func (self *Peers) AllocateTask(task *tasks.Task) (result []byte, err error) {
	peernode := self.GetAPeer()
	for task.Visited(peernode.Addr) {
		peernode = self.GetAPeer()
	}
	log.Printf("Allocated %v from %v to %v", string(task.Id), self.addr, peernode.Addr)
	defer log.Printf("Recieved %v from %v", string(task.Id), peernode.Addr)
	result, err = peernode.AllocateTask(task)
	if err != nil {
		peernode.Kill()
		log.Println(err)
	}
	return

}

//Adds a peer to the set of peers
func (self *Peers) AddPeer(remaddr string) {
	self.RLock()
	if _, exists := self.history[remaddr]; exists {
		self.RUnlock()
		return
	}
	self.RUnlock()
	self.Lock()
	self.history[remaddr] = util.IpValue(remaddr)
	self.Unlock()
	newpeer, err := NewPeer(self.addr, remaddr)
	if err == nil {
		self.peers <- newpeer
		self.route <- newpeer
		peeraddrs, err := newpeer.GetPeers(12)
		if err == nil {
			for _, peeraddr := range peeraddrs {
				self.AddPeer(peeraddr)
			}
		}
	} else {
		self.Lock()
		delete(self.history, remaddr)
		self.Unlock()
	}
	return
}
