package node

import (
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/peers"
	"github.com/rhino1998/cluster/tasks"
	"github.com/rhino1998/cluster/util"
	"log"
	"os/exec"
	//"reflect"
	//"errors"
	"fmt"

	"github.com/rhino1998/cluster/lib/multimutex"
	"sync/atomic"
	"time"
)

type Node struct {
	queuelock    *multimutex.MultiMutex
	allocatelock *multimutex.MultiMutex
	processlock  *multimutex.MultiMutex
	routelock    *multimutex.MultiMutex
	TaskValue    int64
	MaxTasks     int
	MaxRouted    int
	TTL          time.Duration
	Peers        *peers.Peers
	peers.Peer
	LocalIP              string
	lastroutetableupdate time.Time
}

func NewNode(extip, locip string, description info.Info, ttl time.Duration, maxtasks int) *Node {
	return &Node{
		Peers:                peers.NewPeers(extip),
		Peer:                 *peers.ThisPeer(extip, description),
		LocalIP:              locip,
		lastroutetableupdate: time.Now(),
		TaskValue:            0,
		TTL:                  ttl,
		MaxTasks:             maxtasks,
		queuelock:            multimutex.NewMultiMutex(500),
		allocatelock:         multimutex.NewMultiMutex(50),
		processlock:          multimutex.NewMultiMutex(1),
		routelock:            multimutex.NewMultiMutex(maxtasks),
	}
}

func (self *Node) GetPeers(x *int, peerList *[]string) error {
	temp := self.Peers.GetXPeers(*x)
	peers := make([]string, 0, cap(temp))
	if len(temp) > 0 {
		for _, peernode := range temp {
			if peernode != nil {
				peers = append(peers, peernode.Addr)
			}
		}
	}
	return nil
}

func (self *Node) Greet(remaddr *string, desciption *info.Info) error {
	*desciption = info.Info{Compute: self.Compute, Specs: self.Specs}
	self.Peers.AddPeer(*remaddr)
	return nil
}

//Processes a given task
func (self *Node) process(task *tasks.Task) ([]byte, error) {
	//Times the time spent in queue and updates statistic
	//temp := time.Now()

	log.Printf("Added %v to Processing Queue", string(task.Id))

	//Ensures only one task executes at a time
	self.processlock.Lock()
	defer self.processlock.Unlock()
	//self.timeQueue(temp, string(task.Id))
	//defer self.timeExecution(time.Now(), string(task.Id))

	log.Printf("Processing %v %v %v", string(task.Id), atomic.LoadInt64(&self.TaskValue), task.Value)

	//executes actual task
	return exec.Command(fmt.Sprintf("%v", task.FileName), task.Args...).Output()
}

func (self *Node) Put(data *[]byte, success *bool) error {
	log.Println(self.Peers.ClosestPeer(*data).Addr)
	*success = true
	return nil
}

func (self *Node) AllocateTask(task *tasks.Task, result *[]byte) error {
	if len(task.Id) == 0 {
		task.Id = []byte(util.NewUUID())
	}
	//temp := time.Now()
	//self.timeQueue(temp, string(task.Id))
	task.Add(self.Addr)
	for true {
		if self.Compute && (int(atomic.LoadInt64(&self.TaskValue)+int64(task.Value)) < 10000 || len(task.Jumps) >= self.Peers.Length()-1) {
			//Updates current processing value to reflect task queue
			atomic.AddInt64(&self.TaskValue, int64(task.Value))
			//Processes Task
			data, err := self.process(task)
			//Updates current processing value to reflect task queue
			atomic.AddInt64(&self.TaskValue, int64(-task.Value))
			*result = data
			if err == nil {
				log.Println("Successful Process")
				return err
			}
			//self.incrementTotalTaskFailures(string(task.Id))
			log.Println(err)
		} else {
			var err error
			*result, err = self.Peers.AllocateTask(task)
			if err == nil {
				//self.incrementTotalRoutedTasks(string(task.Id))
				return nil
			}
			log.Println(err)
		}
	}
	*result = []byte("Error")
	return nil
}
