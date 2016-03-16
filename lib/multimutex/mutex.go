package multimutex

import ()

type MultiMutex struct {
	tokens chan struct{}
}

func NewMultiMutex(maxlocks int) *MultiMutex {
	return &MultiMutex{tokens: make(chan struct{}, maxlocks)}
}

func (self *MultiMutex) Lock() {
	self.tokens <- struct{}{}
}

func (self *MultiMutex) Unlock() {
	select {
	case <-self.tokens:
	default:
		panic("Unlocked when not locked")
	}
}
