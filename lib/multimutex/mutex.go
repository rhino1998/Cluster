package multimutex

import (
	"sync/atomic"
)

type MultiMutex struct {
	locks    int32
	maxlocks int32
	queue    chan chan bool
}

func NewMultiMutex(maxlocks int) *MultiMutex {
	queue := make(chan chan bool)
	return &MultiMutex{locks: 0, maxlocks: int32(maxlocks), queue: queue}
}

func (self *MultiMutex) Lock() {
	if self.locks+1 > self.maxlocks {
		cont := make(chan bool)
		self.queue <- cont
		select {
		case <-cont:
		}
	}
	atomic.AddInt32(&self.locks, 1)
}

func (self *MultiMutex) Unlock() {
	if self.locks <= 0 {
		panic("Locking Error, locks <= 0")
	}
	atomic.AddInt32(&self.locks, -1)
	select {
	case temp := <-self.queue:
		temp <- true
	default:
	}
}
