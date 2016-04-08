package tasks

import (
	//"github.com/rhino1998/cluster/reqs"
	"fmt"
	"github.com/spaolacci/murmur3"
	"time"
)

type Task struct {
	Id    uint64         `json:"id"`
	Jumps map[string]int `json:"jumps"`
	Name  string         `json:"name"`
	Args  []string       `json:"args"`
	//Reqs     //[]reqs.Req
	Checksum []byte `json:"Checksum"`
	FileName string `json:"filename"`
	Url      string `json:"url"`
	Value    int    `json:"value"`
}

func NewTask(name, url, filename string, args []string, value int) Task {
	return Task{Id: murmur3.Sum64([]byte(fmt.Sprintf("%v", time.Now().UnixNano()))) >> 16, Jumps: make(map[string]int), FileName: filename, Name: name, Url: url, Value: value, Args: args}
}

func (self Task) Add(addr string) {
	self.Jumps[addr] = len(self.Jumps)
}

func (self Task) Visited(addr string) bool {
	_, found := self.Jumps[addr]
	return found
}
