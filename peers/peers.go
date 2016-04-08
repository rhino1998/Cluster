package peers

import (
	"log"
	//"net"
	//"fmt"
	"github.com/rhino1998/cluster/common"
	"github.com/rhino1998/cluster/tasks"
	"github.com/rhino1998/cluster/util"
	// /"runtime"
	//"sync"
	//"time"
)

//A time indexed queryable map of peers
type Peers struct {
	table *Table
	peers chan *Peer
	this  *Peer
}

//Creates an empty list of peers
func NewPeers(this *Peer) *Peers {
	peers := &Peers{peers: make(chan *Peer), this: this, table: NewTable(this.Key())}
	return peers
}

func (self *Peers) Length() int {
	return self.table.Len()
}

func (self *Peers) Update() {
	peers, err := self.GetPeers(self.GetAPeer(), 12)
	if err == nil {
		for _, peeraddr := range peers {
			self.AddPeer(peeraddr)
		}
	}

}

func (self *Peers) GetXPeers(x int) []*Peer {
	temp := make([]*Peer, x, x)
	for i := 0; i < x; i++ {
		temp[i] = self.GetAPeer()
	}
	return temp
}

func (self *Peers) ClosestPeer(key uint64) *Peer {
	peer := self.table.GetClosest(key)
	return peer
}

func (self *Peers) GetAPeer() *Peer {
	peernode := <-self.peers
	//log.Println(peernode.Addr, "got")
	for peernode.isDead() {
		self.table.Delete(peernode.Key())
		peernode = <-self.peers
	}
	//self.peers <- peernode
	return peernode
}

//Adds a peer to the set of peers
func (self *Peers) AddPeer(remaddr string) {
	if remaddr == self.this.Addr || !self.table.Fits(util.IpValue(remaddr)) {

		return
	}
	newpeer, err := NewPeer(self.this.IntAddr, self.this.Addr, remaddr)
	if err == nil && newpeer != nil {
		if self.table.Insert(newpeer) {
			go func() {
				for !newpeer.isDead() {
					self.peers <- newpeer
				}
				self.table.Delete(newpeer.Key())
			}()
		}

		peeraddrs, err := self.GetPeers(newpeer, 48)
		if err == nil {
			for _, peeraddr := range peeraddrs {
				self.AddPeer(peeraddr)
			}
		}
	}
	return
}

//PeerWrappers For Failure

func (self *Peers) Ping(peer *Peer) (err error) {
	err = peer.ping()
	if err != nil {
		self.failure(peer, err)
	}
	return err
}

func (self *Peers) GetPeers(peer *Peer, x int) (peers []string, err error) {
	peers, err = peer.getpeers(x)
	if err != nil {
		self.failure(peer, err)
	}
	return peers, err
}

func (self *Peers) Get(peer *Peer, key string) (data []byte, err error) {
	data, err = peer.get(key)
	if err != nil {
		self.failure(peer, err)
	}
	return data, err
}

func (self *Peers) Put(peer *Peer, item *common.Item) (success bool, err error) {
	success, err = peer.put(item)
	if err != nil {
		self.failure(peer, err)
	}
	return success, err
}

func (self *Peers) AllocateTask(peer *Peer, task *tasks.Task) (result []byte, err error) {
	log.Printf("Allocated %v from %v to %v", string(task.Id), self.this.Addr, peer.Addr)
	defer log.Printf("Recieved %v from %v", string(task.Id), peer.Addr)
	result, err = peer.AllocateTask(task)
	if err != nil {
		self.failure(peer, err)
	}
	return result, err

}

func (self *Peers) failure(peer *Peer, err error) {
	log.Println("FAIL FUCK", err)
	peer.kill()
	self.table.Delete(peer.Key())
}
