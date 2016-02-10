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
	degrade    func(value interface{}) (interface{}, error)
}

func New(name string, comp func(value, check interface{}) (bool, error), value interface{}) *Req {
	return &Req{name: name, comp: comp, value: value}
}

func (self *Req) SetDegrade(degrade func(value interface{}) (interface{}, error)) {
	self.Lock()
	defer self.Unlock()
	self.degradable = true
	self.degrade = degrade
}

func (self *Req) UnSetDegrade() {
	self.Lock()
	defer self.Unlock()
	self.degradable = false
	self.degrade = nil
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

func (self *Req) Degrade() error {
	if !self.degradable {
		return nil
	}
	self.Lock()
	defer self.Unlock()
	degrade := self.degrade
	value, err := degrade(self.value)
	if err != nil {
		return err
	}
	self.value = value
	return nil
}
