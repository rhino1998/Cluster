package main

import (
	"github.com/rhino1998/cluster/lib/multimutex"
	"log"
	"time"
)

type Temp struct {
	queue *multimutex.MultiMutex
}

func (self *Temp) printlock(msg int, length time.Duration) {
	go func() {
		self.queue.Lock()
		log.Println("yo", msg)
		time.Sleep(length)
		self.queue.Unlock()
		log.Println("done", msg)
	}()
}

func main() {
	temp := &Temp{queue: multimutex.NewMultiMutex(5)}
	temp.printlock(1, 5*time.Second)
	temp.printlock(2, 5*time.Second)
	temp.printlock(3, 5*time.Second)
	temp.printlock(4, 6*time.Second)
	temp.printlock(5, 5*time.Second)
	temp.printlock(6, 3*time.Second)
	temp.printlock(7, 5*time.Second)
	temp.printlock(8, 2*time.Second)
	temp.printlock(9, 5*time.Second)
	temp.printlock(10, 7*time.Second)
	temp.printlock(11, 5*time.Second)
	select {}
}
