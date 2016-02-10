package peers

import (
	"cluster/bench"
	"cluster/lib/querymap"
	"cluster/lib/timelist"
	"cluster/reqs"
	"errors"
	"github.com/fatih/structs"
	"sync"
	"time"
	//"net"
)

//A peer node
type Peer struct {
	Addr      string    `json:"addr"`
	Timestamp time.Time `json:"timestamp"`
	Compute   bool      `json:"compute"`
	bench.Specs
}

//A time indexed queryable map of peers
type Peers struct {
	sync.RWMutex
	index *timelist.TimeList
	data  *querymap.QueryMap
	peers map[string]*Peer
	ttl   time.Duration
}

var (
	NoPeerFound  error = errors.New("No peer found")
	NoValidPeers error = errors.New("No valid peers")
)

//Creates an empty list of peers
func NewPeers() *Peers {
	return &Peers{index: timelist.New(), data: querymap.New()}
}

func (self *Peers) GetPeer(key string) (*Peer, error) {
	self.RLock()
	defer self.RUnlock()
	peer, ok := self.peers[key]
	if !ok {
		return nil, NoPeerFound
	}
	return peer, nil
}

func (self *Peers) GetAPeer() (*Peer, error) {
	self.RLock()
	defer self.RUnlock()
	for _, key := range self.data.Keys() {
		return self.peers[key], nil
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
func (self *Peers) BestMatch(reqs []reqs.Req) (*Peer, error) {
	self.RLock()
	defer self.RUnlock()
	data, err := self.applyReqs(reqs)
	if err != nil {
		return nil, err
	}
	for _, key := range data.Keys() {
		return self.peers[key], nil
	}
	return nil, NoValidPeers
}

//Adds a peer to the set of peers
func (self *Peers) AddPeer(peer Peer) {
	self.Lock()
	defer self.Unlock()
	temp := structs.Map(peer)

	delete(temp, "timestamp")
	self.peers[peer.Addr] = &peer
	self.data.Assign(peer.Addr, temp)
	self.index.Insert(peer.Addr, peer.Timestamp)
	return
}
