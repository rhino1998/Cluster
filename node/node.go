package node

import (
	"github.com/muesli/cache2go"
	"github.com/rhino1998/cluster/info"
	"github.com/rhino1998/cluster/peer"
	"github.com/rhino1998/cluster/peers"
	"github.com/rhino1998/cluster/tasks"
	"github.com/rhino1998/cluster/util"
	"github.com/rhino1998/god/client"
	"log"
	"net/http"
	"os/exec"
	//"reflect"
	"errors"
	"fmt"
	"github.com/rhino1998/cluster/lib/godmutex"

	"github.com/rhino1998/cluster/lib/multimutex"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	queuelock    *multimutex.MultiMutex
	allocatelock *multimutex.MultiMutex
	processlock  *multimutex.MultiMutex
	routelock    *multimutex.MultiMutex
	sync.RWMutex
	DB                  *client.Conn
	DBMutex             *godmutex.RWMutex
	TaskValue           int64
	RoutedTasks         int64
	TotalTasks          int64
	TotalTasksCompleted int64
	TotalRoutedTasks    int64
	TotalRouteFailures  int64
	TotalQueueTime      int64
	TotalTaskFailures   int64
	TotalExecutionTime  int64
	TotalTaskTime       int64
	MaxTasks            int
	MaxRouted           int
	TTL                 time.Duration
	Peers               *peers.Peers
	peer.Peer
	peerCache            *cache2go.CacheTable
	LocalIP              string
	lastroutetableupdate time.Time
}

func (self *Node) incrementvalue(key []byte, amount int64, id string) {
	seed := strconv.Itoa(rand.Int())
	self.DBMutex.Lock(key, id+seed)
	var num int
	var err error
	value, found := self.DB.Get(key)
	if !found {
		self.DB.SPut(key, []byte("0"))
		num = 0
	} else {
		num, err = strconv.Atoi(string(value))
		if err != nil {
			log.Println(err)
			return
		}
	}
	self.DB.SPut(key, []byte(strconv.Itoa(num+int(amount))))
	self.DBMutex.Unlock(key, id+seed)
}

func NewNode(extip, locip string, description info.Info, kvstoreaddr string, ttl time.Duration, maxtasks int) *Node {
	clientconn := client.MustConn(kvstoreaddr)
	return &Node{
		DB:                   clientconn,
		Peers:                peers.NewPeers(ttl),
		Peer:                 *peer.ThisPeer(extip, description),
		LocalIP:              locip,
		lastroutetableupdate: time.Now(),
		TaskValue:            0,
		TTL:                  ttl,
		peerCache:            cache2go.Cache("PeerCache"),
		MaxTasks:             maxtasks,
		TotalTasks:           0,
		TotalRoutedTasks:     0,
		TotalRouteFailures:   0,
		TotalTaskFailures:    0,
		queuelock:            multimutex.NewMultiMutex(500),
		allocatelock:         multimutex.NewMultiMutex(50),
		processlock:          multimutex.NewMultiMutex(1),
		routelock:            multimutex.NewMultiMutex(maxtasks),
		DBMutex:              godmutex.NewRWMutex(clientconn, "mutex", extip),
		RoutedTasks:          0}
}

func (self *Node) GreetPeer(addr string) error {
	self.peerCache.Add(addr, 2*self.TTL, nil)
	newpeer, err := peer.NewPeer(self.Addr, addr, 1*time.Second)
	if err != nil {
		return err
	}
	self.Peers.AddPeer(newpeer)
	self.Peers.Clean(self.Addr)
	newpeernodeaddrs, err := newpeer.GetPeers(12, 1*time.Second)
	if err == nil {
		for _, newpeernodeaddr := range newpeernodeaddrs {
			if !self.Peers.Exists(newpeernodeaddr) && newpeernodeaddr != self.Addr {
				self.GreetPeer(newpeernodeaddr)
			}

		}
	}
	return nil
}

func (self *Node) GetPeers(r *http.Request, x *int, peerList *[]string) error {
	temp := self.Peers.FirstX(*x)
	peers := make([]string, 0, cap(temp))
	if len(temp) > 0 {
		for _, peernode := range temp {
			if peernode != nil {
				peers = append(peers, peernode.Addr)
			}
		}
	}
	*peerList = peers
	return nil
}

func (self *Node) Ping(start time.Time, end *time.Time) error {
	*end = time.Now()
	return nil
}

func (self *Node) Greet(r *http.Request, remaddr *string, desciption *info.Info) error {
	*desciption = info.Info{Compute: self.Compute, Specs: self.Specs}
	if !self.peerCache.Exists(*remaddr) && !self.Peers.Exists(*remaddr) {
		self.GreetPeer(*remaddr)
	}
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

	log.Printf("Processing %v %v %v", string(task.Id), self.TaskValue, task.Value)

	//executes actual task
	return exec.Command(fmt.Sprintf("%v", task.FileName), task.Args...).Output()
}

func (self *Node) NewTask(task tasks.Task) error {
	temp := time.Now()
	log.Printf("Added %v to queue", string(task.Id))
	//self.incrementTotalTasks(string(task.Id))
	task.Add(self.Addr)
	go func() {
		//defer self.decrementTotalTasks(string(task.Id))
		//defer self.incrementTotalTasksCompleted(string(task.Id))
		//defer self.timeTask(time.Now(), string(task.Id))
		for true {
			peernode, err := self.Peers.GetAPeer()
			if err == nil {
				result, err := peernode.AllocateTask(&task)
				if err == nil {
					elapsed := time.Now().Sub(temp).Nanoseconds() / 1000000
					log.Println("Allocated", err, string(result), string(task.Id), elapsed)
					return
				}
				log.Println(err)
				log.Println(string(task.Id))
			}
		}
	}()
	return nil
}

func (self *Node) AllocateTask(r *http.Request, task *tasks.Task, result *[]byte) error {
	if len(task.Id) == 0 {
		task.Id = []byte(util.NewUUID())
	}
	//temp := time.Now()
	self.queuelock.Lock()
	defer self.queuelock.Unlock()
	//self.timeQueue(temp, string(task.Id))
	if len(task.Jumps) == 0 {
		//self.incrementTotalTasks(string(task.Id))
		//defer self.timeTask(time.Now(), string(task.Id))
		//defer self.decrementTotalTasks(string(task.Id))
		//defer self.incrementTotalTasksCompleted(string(task.Id))
	}
	task.Add(self.Addr)
	fails := 0
	for true {
		self.allocatelock.Lock()
		if self.Compute && (int(self.TaskValue+int64(task.Value)) < 10000 || len(task.Jumps) >= self.Peers.Length()-1) {
			//Updates current processing value to reflect task queue
			atomic.AddInt64(&self.TaskValue, int64(task.Value))
			//Processes Task
			self.allocatelock.Unlock()
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
			self.routelock.Lock()
			atomic.AddInt64(&self.RoutedTasks, 1)
			defer atomic.AddInt64(&self.RoutedTasks, -1)
			peernode, err := self.Peers.GetAPeer()
			if err == nil && !task.Visited(peernode.Addr) {
				log.Printf("Allocated %v from %v to %v", string(task.Id), self.Addr, peernode.Addr)
				atomic.AddInt64(&self.RoutedTasks, 1)
				self.allocatelock.Unlock()
				*result, err = peernode.AllocateTask(task)
				atomic.AddInt64(&self.RoutedTasks, -1)
				if err == nil {
					//self.incrementTotalRoutedTasks(string(task.Id))
					self.routelock.Unlock()
					log.Printf("Recieved %v from %v", string(task.Id), peernode.Addr)
					return nil
				}
				//self.incrementTotalRouteFailures(string(task.Id))
			} else {
				fails++
				if fails > 10 {
					self.routelock.Unlock()
					return errors.New("Failed too many times")
				}
			}
			self.routelock.Unlock()
		}
	}
	*result = []byte("Error")
	return nil
}
