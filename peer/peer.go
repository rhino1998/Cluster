package peer

import (
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/lib/jsonrpc"
	"github.com/rhino1998/cluster/tasks"
	//"log"
	"errors"
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
func NewPeer(locaddr, remaddr string, timeout time.Duration) (*Peer, error) {
	connection := jsonrpc.NewClient(remaddr)
	errchan := make(chan error)
	resultchan := make(chan info.Info)
	var description info.Info
	go func() {
		err := connection.Call("Node.Greet", &locaddr, &description)
		errchan <- err
		resultchan <- description
	}()
	var err error
	select {
	case err = <-errchan:
		description = <-resultchan
	case <-time.After(timeout):
		return nil, errors.New("Timeout")
	}
	return &Peer{Addr: remaddr, Info: description, Connection: connection}, err
}

func (self *Peer) Ping() (time.Time, error) {
	return time.Now(), nil
}

func (self *Peer) GetPeers(x int, timeout time.Duration) (peers []string, err error) {
	errchan := make(chan error)
	resultchan := make(chan []string)
	go func() {
		err := self.Connection.Call("Node.GetPeers", &x, &peers)
		errchan <- err
		resultchan <- peers
	}()
	select {
	case err = <-errchan:
		peers = <-resultchan
	case <-time.After(timeout):
		return nil, errors.New("Timeout")
	}

	if err != nil {
		return nil, err
	}
	return peers, nil

}

func (self *Peer) AllocateTask(task *tasks.Task) (result []byte, err error) {
	err = self.Connection.Call("Node.AllocateTask", task, &result)
	return result, err
}
