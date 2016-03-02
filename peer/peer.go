package peer

import (
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/lib/jsonrpc"
	"github.com/rhino1998/cluster/tasks"
	//"log"
	"sync"
	"time"
)

//A peer node
type Peer struct {
	sync.RWMutex
	Connection *jsonrpc.Client
	Addr       string `json:"addr"`
	info.Info
}

func ThisPeer(addr string, description info.Info) *Peer {
	return &Peer{Addr: addr, Connection: jsonrpc.NewClient(addr), Info: description}
}

//initializes a new peer
func NewPeer(locaddr, remaddr string) (*Peer, error) {
	connection := jsonrpc.NewClient(remaddr)
	var desciption info.Info
	err := connection.Call("Node.Greet", &locaddr, &desciption)
	return &Peer{Addr: remaddr, Info: desciption, Connection: connection}, err
}

func (self *Peer) Ping() (time.Time, error) {
	return time.Now(), nil
}

func (self *Peer) GetPeers(x int) (peers []string, err error) {
	err = self.Connection.Call("Node.GetPeers", &x, &peers)
	if err != nil {
		return nil, err
	}
	return peers, nil

}

func (self *Peer) AllocateTask(task *tasks.Task) (result []byte, err error) {
	err = self.Connection.Call("Node.AllocateTask", task, &result)
	return result, err
}
