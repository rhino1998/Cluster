package peers

import (
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/tasks"
	//"log"

	"net/rpc"
	"sync/atomic"
)

//A peer node
type Peer struct {
	client *rpc.Client
	Addr   string `json:"addr"`
	info.Info
	dead uint32
}

func ThisPeer(addr string, description info.Info) *Peer {
	return &Peer{Addr: addr, Info: description}
}

//initializes a new peer
func NewPeer(locaddr, remaddr string) (*Peer, error) {
	client, err := rpc.Dial("tcp", remaddr)
	if err != nil {
		return nil, err
	}
	var description info.Info
	err = client.Call("Node.Greet", &locaddr, &description)
	if err != nil {
		client.Close()
		return nil, err
	}
	return &Peer{Addr: remaddr, Info: description, client: client, dead: 0}, err
}

/*func (self *Peer) Ping() (time.Time, error) {
	return time.Now(), nil
}*/

func (self *Peer) IsDead() bool {
	return atomic.LoadUint32(&self.dead) == 1
}

func (self *Peer) Kill() {
	atomic.StoreUint32(&self.dead, 1)
	self.client.Close()
}
func (self *Peer) Describe(remaddr string) (description info.Info, err error) {
	err = self.client.Call("Node.Describe", &remaddr, &description)
	return description, err
}

func (self *Peer) GetPeers(x int) (peers []string, err error) {
	err = self.client.Call("Node.GetPeers", &x, &peers)
	return peers, err
}

func (self *Peer) AllocateTask(task *tasks.Task) (result []byte, err error) {
	err = self.client.Call("Node.AllocateTask", task, &result)
	return result, err
}
