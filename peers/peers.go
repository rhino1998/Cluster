package peers

import (
	"errors"
	"github.com/fatih/structs"
	"github.com/rhino1998/cluster/lib/querymap"
	"github.com/rhino1998/cluster/lib/timelist"
	"github.com/rhino1998/cluster/peer"
	"github.com/rhino1998/cluster/reqs"
	//"log"
	"math/rand"
	"sync"
	"time"
	//"net"
)

//A time indexed queryable map of peers
type Peers struct {
	sync.RWMutex
	index *timelist.TimeList
	data  *querymap.QueryMap
	peers map[string]*peer.Peer
	ttl   time.Duration
}

var (
	NoPeerFound  error = errors.New("Peer not found")
	NoValidPeers error = errors.New("No valid peers")
)

//Creates an empty list of peers
func NewPeers(ttl time.Duration) *Peers {
	return &Peers{index: timelist.NewTimeList(), data: querymap.NewQueryMap(), peers: make(map[string]*peer.Peer), ttl: ttl}
}

func (self *Peers) Length() int {
	return len(self.peers)
}

func (self *Peers) GetPeer(addr string) (*peer.Peer, error) {
	self.RLock()
	defer self.RUnlock()
	peerNode, ok := self.peers[addr]
	if !ok {
		return nil, NoPeerFound
	}
	return peerNode, nil
}

func (self *Peers) LastX(x int) []*peer.Peer {
	self.RLock()
	defer self.RUnlock()
	temp := make([]*peer.Peer, 0, x)
	itemsfound := make(map[string]int)
	for _, addr := range self.index.LastX(x).Items() {
		if _, found := itemsfound[string(addr.Value())]; !found {
			temp = append(temp, self.peers[string(addr.Value())])
			itemsfound[string(addr.Value())] = 1
		}
	}
	return temp
}

func (self *Peers) FirstX(x int) []*peer.Peer {
	self.RLock()
	defer self.RUnlock()
	temp := make([]*peer.Peer, 0, x)
	itemsfound := make(map[string]int)
	for _, addr := range self.index.FirstX(x).Items() {
		if _, found := itemsfound[string(addr.Value())]; !found {
			temp = append(temp, self.peers[string(addr.Value())])
			itemsfound[string(addr.Value())] = 1
		}
	}
	return temp
}

//checks whether peer exists by address
func (self *Peers) Exists(addr string) bool {
	self.RLock()
	defer self.RUnlock()
	_, found := self.peers[addr]
	return found
}

func (self *Peers) GetAPeer() (*peer.Peer, error) {
	self.RLock()
	defer self.RUnlock()
	if len(self.peers) >= 1 {
		return self.peers[self.data.Keys()[rand.Intn(len(self.peers))]], nil
	}
	return nil, NoValidPeers
}

func (self *Peers) applyReqs(reqs []reqs.Req) (*querymap.QueryMap, error) {
	var err error
	data := self.data
	for _, req := range reqs {
		data, err = self.data.Mask(req.Name(), req.Comp, req.Value())
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

//Applies requirements and returns a value which satisfies those requirements
//returns an error if no peer meets requirements
func (self *Peers) BestMatch(reqs []reqs.Req) (*peer.Peer, error) {
	self.RLock()
	defer self.RUnlock()
	data, err := self.applyReqs(reqs)
	if err != nil {
		return nil, err
	}
	for _, addr := range data.Keys() {
		return self.peers[addr], nil
	}
	return nil, NoValidPeers
}

func (self *Peers) cleanworker(locaddr, addr string) {
	peernode, err := peer.NewPeer(locaddr, addr, 1*time.Second)
	if err == nil {
		self.peers[peernode.Addr] = peernode
		self.data.Assign(peernode.Addr, structs.Map(peernode))
		self.index.Insert([]byte(peernode.Addr), time.Now().UTC())
	} else {
		delete(self.peers, addr)
		self.data.Delete(addr)
	}
	return
}

func (self *Peers) Clean(locaddr string) error {
	self.Lock()
	temp, err := self.index.PopBefore(time.Now().UTC().Add(-self.ttl))
	self.Unlock()
	if err != nil {
		return err
	}
	for _, index := range temp.Items() {
		if !self.index.Exists(index.Value()) && temp.Exists(index.Value()) {
			go self.cleanworker(locaddr, string(index.Value()))

		}
	}
	return nil
}

func (self *Peers) Items() []*peer.Peer {
	self.Lock()
	defer self.Unlock()
	i := 0
	temp := make([]*peer.Peer, len(self.peers), len(self.peers))
	for _, val := range self.peers {
		temp[i] = val
		i++
	}
	return temp
}

//Returns peers after time ****Broken****
func (self *Peers) After(start time.Time) ([]*peer.Peer, error) {
	self.Lock()
	defer self.Unlock()
	//addrs := self.index.After(start)
	if true {
		temp := make([]*peer.Peer, self.index.Length(), self.index.Length())
		for i, addr := range self.index.Items() {
			temp[i] = self.peers[string(addr.Value())]
		}
		return temp, nil
	}
	return nil, nil
}

//Adds a peer to the set of peers
func (self *Peers) AddPeer(peernode *peer.Peer) {
	self.Lock()
	defer self.Unlock()
	self.peers[peernode.Addr] = peernode
	self.data.Assign(peernode.Addr, structs.Map(peernode))
	self.index.Insert([]byte(peernode.Addr), time.Now().UTC())
	return
}
