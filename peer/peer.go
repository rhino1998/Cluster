package peer

import (
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/lib/jsonrpc"
	"github.com/rhino1998/cluster/tasks"
	"log"
	"sync"
	"time"
)

//A peer node
type Peer struct {
	sync.RWMutex
	Connection *jsonrpc.Client
	Addr       string    `json:"addr"`
	Timestamp  time.Time `json:"timestamp"`
	info.Info
}

func ThisPeer(addr string, description info.Info) *Peer {
	connection := jsonrpc.NewClient(addr)
	return &Peer{Addr: addr, Info: description, Timestamp: time.Now().UTC(), Connection: connection}
}

//initializes a new peer
func NewPeer(locaddr string, thisDescription info.Info, remaddr string) (*Peer, error) {
	connection := jsonrpc.NewClient(remaddr)
	var desciption info.Info
	err := connection.Call("Node.Greet", ThisPeer(locaddr, thisDescription), &desciption)
	log.Println(desciption)
	return &Peer{Addr: remaddr, Info: desciption, Timestamp: time.Now().UTC(), Connection: connection}, err
}

func (self *Peer) Ping() (time.Time, error) {
	return time.Now(), nil
}

func (self *Peer) AllocateTask(task *tasks.Task) (result *[]byte, err error) {
	err = self.Connection.Call("Node.Allocate", true, &result)
	return result, err
}
