package peer

import (
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/lib/jsonrpc"
	"github.com/rhino1998/cluster/tasks"
	"sync"
	"time"
)

//A peer node
type Peer struct {
	sync.RWMutex
	connection *jsonrpc.Client
	Addr       string    `json:"addr"`
	Timestamp  time.Time `json:"timestamp"`
	info.Info
}

//initializes a new peer
func NewPeer(addr string) (*Peer, error) {
	connection := jsonrpc.NewClient(addr)
	var desciption info.Info
	err := connection.Call("Node.Describe", true, &desciption)
	return &Peer{Addr: addr, Info: desciption, Timestamp: time.Now().UTC(), connection: connection}, err
}

func (self *Peer) Ping() (time.Time, error) {
	return time.Now(), nil
}

func (self *Peer) AllocateTask(task *tasks.Task) (result *[]byte, err error) {
	err = self.connection.Call("Node.Allocate", true, &result)
	return result, err
}
