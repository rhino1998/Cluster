package peers

import (
	//"github.com/Workiva/go-datastructures/trie/yfast"
	"log"
	"sync"
	//"time"
)

type Table struct {
	//trie   *yfast.YFastTrie
	peers  []*Peer
	center uint64
	sync.RWMutex
}

func NewTable(center uint64) *Table {
	arr := make([]*Peer, 48, 48)
	return &Table{peers: arr, center: center}
}

func (self *Table) GetClosest(key uint64) *Peer {
	return self.peers[self.bitslot(key)]
}

func (self *Table) bitslot(key uint64) int {
	pos := 0
	min := ^uint64(0) >> 16
	for i, peer := range self.peers {
		self.RLock()
		if peer != nil {
			dist := key ^ (self.center ^ (uint64(1) << uint(i)))
			if dist < min {
				pos = i
				min = dist
			}
		}
		self.RUnlock()
	}
	return pos
}

func (self *Table) Fits(key uint64) bool {
	for i := len(self.peers) - 1; i >= 0; i-- {
		self.RLock()
		peernode := self.peers[i]
		if peernode == nil {
			self.RUnlock()
			return true
		}
		if peernode.Key() != key {
			mask := (uint64(1) << uint(i))
			if peernode.isDead() || peernode.Key()^(self.center^mask) > key^(self.center^mask) {
				self.RUnlock()
				return true
			}
		} else {
			self.RUnlock()
			return false
		}
		self.RUnlock()
	}
	return false
}

func (self *Table) Len() int {
	self.RLock()
	defer self.RUnlock()
	sum := 0
	for _, peernode := range self.peers {
		if peernode != nil {
			sum++
		}
	}
	return sum
}

func (self *Table) items() []*Peer {
	self.RLock()
	defer self.RUnlock()
	return self.peers[1:]
}

func (self *Table) Delete(key uint64) {
	for i := len(self.peers) - 1; i >= 0; i-- {
		self.Lock()
		peernode := self.peers[i]
		if peernode != nil && peernode.Key() == key {
			peernode.kill()
			self.peers[i] = nil
		}
		self.Unlock()
	}
}

func (self *Table) Addrs() {
	for i, peernode := range self.peers {
		if peernode != nil {
			log.Printf("%v %02v %048b", peernode.Addr, i, peernode.Key())
		} else {
			log.Println(nil)
		}
		//time.Sleep(50 * time.Millisecond)
	}
	log.Printf("                       %048b", self.center)
}

func (self *Table) Exists(addr string) bool {
	self.RLock()
	defer self.RUnlock()
	for _, peernode := range self.peers {
		if peernode != nil && peernode.Addr == addr {
			return true
		}
	}
	return false
}

func (self *Table) insert(peer *Peer, offset int) bool {
	if peer == nil {
		log.Printf("PeerNil %v", offset)
		return false
	}
	if !peer.isDead() && offset < len(self.peers) {
		for i := len(self.peers) - 1 - offset; i >= 0; i-- {
			peernode := self.peers[i]

			mask := (uint64(1) << uint(i))
			if peernode == nil || (peernode.isDead() || peernode.Key()^(self.center^mask) > peer.Key()^(self.center^mask)) {
				self.insert(self.peers[i], len(self.peers)-i)
				self.peers[i] = peer
				//log.Println(peer.Addr, pos)
				return true
			} else if peer.Key() == peernode.Key() {
				if peer != peernode {
					log.Println("killed because same", peer.Addr, offset)
					peer.kill()
				}
				return false
			}
		}
	}
	log.Println("kill cause no fit", peer.Addr)
	peer.kill()
	return false
}

func (self *Table) Insert(peer *Peer) bool {
	//log.Printf("Inserting: %v", peer.Addr)
	//self.Addrs()
	if !self.Exists(peer.Addr) {
		self.Lock()
		defer self.Unlock()
		res := self.insert(peer, 0)
		return res
	}
	return false
}
