package node

import (
	"github.com/rhino1998/cluster/db"
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/peers"
	"github.com/rhino1998/cluster/tasks"
	"net/http"
	"os/exec"
	"reflect"
	"sync"
	"time"
)

type Node struct {
	sync.RWMutex
	DB                   *db.TransactionLayer
	Tasks                int64
	Peers                *peers.Peers
	Addr                 string `json:"addr"`
	LocalIP              string
	lastroutetableupdate time.Time
	info.Info
}

func NewNode(extip, locip string, description info.Info, layer *db.TransactionLayer) *Node {
	return &Node{Peers: peers.NewPeers(), Addr: extip, LocalIP: locip, Info: description, lastroutetableupdate: time.Now(), Tasks: 0, DB: layer}
}

func (self *Node) GetPeers(r *http.Request, start time.Time, peerList []string) error {
	temp, err := self.Peers.After(start)
	peerList = make([]string, len(temp))
	for i, peer := range temp {
		peerList[i] = peer.Addr
	}
	return err
}

func (self *Node) Ping(start time.Time, end *time.Time) error {
	*end = time.Now()
	return nil
}

func (self *Node) Describe(r *http.Request, n *bool, desciption *info.Info) error {
	desciption = &info.Info{Specs: self.Specs, Compute: self.Compute}
	return nil
}

func (self *Node) process(task *tasks.Task) ([]byte, error) {
	return exec.Command(task.Loc).Output()
}

func (self *Node) RouteTask(r *http.Request, task *tasks.Task, result *[]byte) error {
	task.Jumps[self.Addr] = len(task.Jumps) + 1
	for _, req := range task.Reqs {
		if ok, err := req.Comp(req.Value(), reflect.ValueOf(self).FieldByName(req.Name())); !ok || !self.Compute {
			if err != nil {
				return err
			}
			peer, err := self.Peers.BestMatch(task.Reqs)
			if err != nil {
				return err
			}
			result, err = peer.AllocateTask(task)
			if err != nil {
				return err
			}
		}
	}
	data, err := self.process(task)
	if err != nil {
		return err
	}
	result = &data
	return nil

}
