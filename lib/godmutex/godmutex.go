package godmutex

import (
	"errors"
	"fmt"
	"github.com/rhino1998/god/client"
	"time"
)

var (
	MutexNotFound error = errors.New("Mutex not found")
)

type RWMutex struct {
	mutexname  string
	connection *client.Conn
	id         string
}

func NewRWMutex(connection *client.Conn, mutexname, id string) *RWMutex {
	return &RWMutex{mutexname: mutexname, connection: connection, id: id}
}

func (self *RWMutex) Lock(key []byte) error {
	keyname := []byte(fmt.Sprintf("%v:%v", string(key), self.mutexname))
	for {
		if val, found := self.connection.Get(keyname); found && len(val) == 0 {
			self.connection.SPut(keyname, []byte(self.id))
			if val, found := self.connection.Get(keyname); found && string(val) == self.id {
				return nil
			}
		} else if !found {
			return MutexNotFound
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func (self *RWMutex) Unlock(key []byte) error {
	keyname := []byte(fmt.Sprintf("%v:%v", string(key), self.mutexname))
	for {
		if val, found := self.connection.Get(keyname); found && string(val) == self.id {
			self.connection.SPut(keyname, nil)
			if val, found := self.connection.Get(keyname); found && len(val) == 0 {
				return nil
			}
		} else if !found {
			return MutexNotFound
		}
		return nil
	}

}
