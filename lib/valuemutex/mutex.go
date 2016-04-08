package valuemutex

import (
	"sync/atomic"
)

type ValueMutex struct {
	locks  int32
	maxval int32
	queue  chan chan struct{}
}

func NewValueMutex(maxval int32) *ValueMutex {
	queue := make(chan chan struct{})
	return &ValueMutex{locks: 0, maxval: int32(maxval), queue: queue}
}

func (self *ValueMutex) Lock(val int32) {
	if self.locks+val > self.maxval {
		cont := make(chan struct{})
		self.queue <- cont
		select {
		case <-cont:
		}
	}
	atomic.AddInt32(&self.locks, val)
}

func (self *ValueMutex) Value() int32 {
	return atomic.LoadInt32(&self.locks)
}

func (self *ValueMutex) Unlock(val int32) {
	if self.locks <= 0 {
		panic("Locking Error, locks <= 0")
	}
	atomic.AddInt32(&self.locks, -val)
	select {
	case temp := <-self.queue:
		temp <- struct{}{}
	default:
	}
}
