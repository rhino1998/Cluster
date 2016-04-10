package node

import (
	"github.com/rhino1998/cluster/common"
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/peers"
	"github.com/rhino1998/cluster/tasks"
	//"github.com/rhino1998/cluster/util"
	"github.com/spaolacci/murmur3"
	"log"
	"os/exec"
	//"reflect"
	//"errors"
	"fmt"
	"github.com/rhino1998/cluster/lib/multimutex"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func abssub(a, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}

type Node struct {
	sync.RWMutex
	data         map[string][]byte
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
	lastroutetableupdate time.Time
}

func NewNode(extip, locip string, description info.Info, ttl time.Duration, maxtasks int) *Node {
	thispeer := peers.ThisPeer(locip, extip, description)
	return &Node{
		data:                 make(map[string][]byte),
		Peers:                peers.NewPeers(thispeer),
		Peer:                 *thispeer,
		lastroutetableupdate: time.Now(),
		TaskValue:            0,
		TTL:                  ttl,
		MaxTasks:             maxtasks,
		queuelock:            multimutex.NewMultiMutex(500),
		allocatelock:         multimutex.NewMultiMutex(500),
		processlock:          multimutex.NewMultiMutex(1),
		routelock:            multimutex.NewMultiMutex(500),
	}
}

func (self *Node) GetPeers(x *int, peers *[]string) error {
	temp := self.Peers.GetXPeers(*x)
	if len(temp) > 0 {
		for _, peernode := range temp {
			if peernode != nil {
				*peers = append(*peers, peernode.Addr)
			}
		}
	}
	return nil
}

func (self *Node) Greet(remaddr *string, desciption *info.Info) error {
	*desciption = self.Info
	go self.Peers.AddPeer(*remaddr)
	return nil
}

func (self *Node) Ping(a *struct{}, b *struct{}) error {
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
	if _, err := os.Stat(task.FileName); os.IsNotExist(err) {
		// path/to/whatever does not exist
		self.download_task(task)
	}

	return exec.Command(fmt.Sprintf("%v", task.FileName), task.Args...).Output()
}

func (self *Node) Put(item *common.Item, success *bool) error {
	*success = false
	closestpeer := self.Peers.ClosestPeer(item.Value())
	log.Println(item.Value()^closestpeer.Key(), item.Value()^self.Key(), closestpeer.Addr)
	if item.Value()^closestpeer.Key() > item.Value()^self.Key() {
		self.Lock()
		self.data[item.Key] = item.Data
		self.Unlock()
		*success = true
		return nil
	} else {
		log.Println(self.Key(), closestpeer.Key(), item.Value())
		var err error
		*success, err = self.Peers.Put(closestpeer, item)
		return err
	}
	return nil
}

func (self *Node) Get(key *string, data *[]byte) error {
	var found bool
	self.RLock()
	if *data, found = self.data[*key]; found {
		self.RUnlock()
		return nil
	}
	self.RUnlock()
	closestpeer := self.Peers.ClosestPeer(murmur3.Sum64([]byte(*key)) >> 16)
	log.Println(self.Key(), closestpeer.Key())
	var err error
	*data, err = self.Peers.Get(closestpeer, *key)
	return err
}

func (self *Node) AllocateTask(task *tasks.Task, result *[]byte) error {
	self.routelock.Lock()
	defer self.routelock.Unlock()
	if task.Id == 0 {
		task.Id = murmur3.Sum64([]byte(fmt.Sprintf("%v", time.Now().UnixNano()))) >> 16
	}
	//temp := time.Now()
	//self.timeQueue(temp, string(task.Id))
	task.Add(self.Addr)
	for true {
		if self.Compute && (int(atomic.LoadInt64(&self.TaskValue)+int64(task.Value)) < 10000 || len(task.Jumps) >= 48) {
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
			peer := self.Peers.GetAPeer()
			tries := 0
			for task.Visited(peer.Addr) && tries < 48 {
				time.Sleep(1 * time.Millisecond)
				peer = self.Peers.GetAPeer()
				tries++
			}
			*result, err = self.Peers.AllocateTask(peer, task)
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

func (self *Node) download_task(task *tasks.Task) string {
	output, err := os.Create(task.FileName)
	if err != nil {
		log.Println("Error while creating", task.FileName, "-", err)
		return ""
	}
	defer output.Close()

	response, err := http.Get(task.Url)
	if err != nil {
		fmt.Println("Error while downloading", task.Url, "-", err)
		return ""
	}
	defer response.Body.Close()
	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", task.Url, "-", err)
		return ""
	}

	fmt.Println(n, "bytes downloaded.")
	return task.FileName
}
