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

func (self *RWMutex) Lock(key []byte, id string) error {
	keyid := fmt.Sprintf("%v:%v", string(key), id)
	keyname := []byte(fmt.Sprintf("%v:%v", string(key), self.mutexname))
	for {
		if val, _ := self.connection.Get(keyname); len(val) == 0 || string(val) == keyid {
			self.connection.SPut(keyname, []byte(keyid))
			if val, _ := self.connection.Get(keyname); string(val) == keyid {
				return nil
			}
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func (self *RWMutex) Unlock(key []byte, id string) error {
	keyid := fmt.Sprintf("%v:%v", string(key), id)
	keyname := []byte(fmt.Sprintf("%v:%v", string(key), self.mutexname))
	for {
		if val, found := self.connection.Get(keyname); found && string(val) == keyid {
			self.connection.Del(keyname)
			if _, found := self.connection.Get(keyname); !found {
				return nil
			}
		} else if !found {
			return MutexNotFound
		}
		return nil
	}

}
