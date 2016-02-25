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

func (self *TransactionLayer) Delete(id string, key string) error {
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

func (self *TransactionLayer) Get(id string, key []byte) ([]byte, error) {
	self.RLock()
	if _, found := self.transactions[id]; !found {
		data := &common.Item{}
		err := self.DB.Get(common.Item{Key: key}, data)
		return data.Value, err
	}
	self.RUnlock()
	self.Lock()
	data, err := self.transactions[id].Value(key)
	self.Unlock()
	return data.Data().([]byte), err
}

func (self *TransactionLayer) push(key interface{}, item *cache2go.CacheItem) {
	self.DB.Put(common.Item{Key: key.([]byte), Value: item.Data().([]byte), Sync: false})
}

func (self *TransactionLayer) Commit(id string) error {
	self.transactions[id].Foreach(self.push)
	return nil

}
