package reqs

import (
	"sync"
)

type Req struct {
	sync.RWMutex
	degradable bool
	value      interface{}
	comp       func(value, check interface{}) (bool, error)
	name       string
}

func New(name string, comp func(value, check interface{}) (bool, error), value interface{}) *Req {
	return &Req{name: name, comp: comp, value: value}
}

func (self *Req) Name() string {
	self.RLock()
	defer self.RUnlock()
	return self.name
}

func (self *Req) Value() interface{} {
	self.RLock()
	defer self.RUnlock()
	return self.value
}

func (self *Req) Comp(value, check interface{}) (bool, error) {
	self.RLock()
	defer self.RUnlock()
	return self.comp(value, check)
}
