package db

import (
	"errors"
	"github.com/muesli/cache2go"
	"github.com/rhino1998/god/common"
	"github.com/rhino1998/god/dhash"
	"sync"
)

var (
	TransactionNotFound      error = errors.New("Cache not found")
	TransactionAlreadyExists error = errors.New("Cache already exists")
)

type TransactionLayer struct {
	sync.RWMutex
	transactions map[string]*cache2go.CacheTable
	DB           *dhash.Node
}

func NewTransactionLayer(node *dhash.Node) *TransactionLayer {
	return &TransactionLayer{transactions: make(map[string]*cache2go.CacheTable), DB: node}
}

func (self *TransactionLayer) NewTransaction(id string) error {
	self.RLock()
	if _, found := self.transactions[id]; found {
		return TransactionAlreadyExists
	}
	self.RUnlock()
	self.Lock()
	self.transactions[id] = cache2go.Cache(id)
	self.Unlock()
	return nil
}

func (self *TransactionLayer) Add(id string, key string, data []byte) error {
	self.RLock()
	if _, found := self.transactions[id]; !found {
		return TransactionNotFound
	}
	self.RUnlock()
	self.Lock()
	self.transactions[id].Add(key, 0, data)
	self.Unlock()
	return nil
}

func (self *TransactionLayer) Del(id string, key string) error {
	self.RLock()
	if _, found := self.transactions[id]; !found {
		return TransactionNotFound
	}
	self.RUnlock()
	self.Lock()
	self.transactions[id].Delete(key)
	self.Unlock()
	return nil
}

func (self *TransactionLayer) Get(id string, key string) ([]byte, bool) {
	self.RLock()
	if _, found := self.transactions[id]; !found {
		return self.DBGet(id, key)
	}
	self.RUnlock()
	self.Lock()
	data, err := self.transactions[id].Value([]byte(key))
	self.Unlock()
	return data.Data().([]byte), err == nil
}

func (self *TransactionLayer) DBPut(id string, key string, data []byte) {
	self.DB.SubPut(common.Item{Key: []byte(id), SubKey: []byte(key), Value: data, Sync: false})
}

func (self *TransactionLayer) DBDel(id string, key string) {
	self.DB.SubDel(common.Item{Key: []byte(id), SubKey: []byte(key), Sync: false})
}

func (self *TransactionLayer) DBGet(id string, key string) ([]byte, bool) {
	var result *common.Item
	self.DB.SubGet(common.Item{Key: []byte(id), SubKey: []byte(key), Sync: false}, result)
	return result.Value, result.Exists
}

func (self *TransactionLayer) push(id string) func(key interface{}, item *cache2go.CacheItem) {
	return func(key interface{}, item *cache2go.CacheItem) {
		self.DB.SubPut(common.Item{Key: []byte(id), SubKey: key.([]byte), Value: item.Data().([]byte), Sync: false})
	}

}

func (self *TransactionLayer) Commit(id string) error {
	self.transactions[id].Foreach(self.push(id))
	return nil

}

func (self *TransactionLayer) Discard(id string) error {
	self.transactions[id].Flush()
	return nil
}
