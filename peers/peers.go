package peers

import (
	"errors"
	"github.com/fatih/structs"
	"github.com/rhino1998/cluster/lib/querymap"
	"github.com/rhino1998/cluster/lib/timelist"
	"github.com/rhino1998/cluster/peer"
	"github.com/rhino1998/cluster/reqs"
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
func NewPeers() *Peers {
	return &Peers{index: timelist.New(), data: querymap.New()}
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
	for _, addr := range self.data.Keys() {
		return self.peers[addr], nil
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

//Returns peers after time
func (self *Peers) After(start time.Time) ([]*peer.Peer, error) {
	self.Lock()
	defer self.Unlock()
	var found bool
	addrs := self.index.After(start)
	temp := make([]*peer.Peer, addrs.Length(), addrs.Length())
	for i, addr := range addrs.Items() {
		temp[i], found = self.peers[string(addr.Value())]
		if !found {
			return nil, NoPeerFound
		}
	}
	return temp, nil
}

//Adds a peer to the set of peers
func (self *Peers) AddPeer(peer peer.Peer) {
	self.Lock()
	defer self.Unlock()
	temp := structs.Map(peer)

	delete(temp, "timestamp")
	self.peers[peer.Addr] = &peer
	self.data.Assign(peer.Addr, temp)
	self.index.Insert([]byte(peer.Addr), peer.Timestamp)
	return
}
