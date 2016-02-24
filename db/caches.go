package db

import (
	"errors"
	"github.com/muesli/cache2go"
	"github.com/rhino1998/god/common"
	"github.com/rhino1998/god/dhash"
	"sync"
)

var (
	CacheNotFound      error = errors.New("Cache not found")
	CacheAlreadyExists error = errors.New("Cache already exists")
)

type CacheLayer struct {
	sync.RWMutex
	Caches map[string]*cache2go.CacheTable
	db     *dhash.Node
}

func NewCacheLayer(node *dhash.Node) *CacheLayer {
	return &CacheLayer{Caches: make(map[string]*cache2go.CacheTable), db: node}
}

func (self *CacheLayer) AddCache(id string) error {
	self.RLock()
	if _, found := self.Caches[id]; found {
		return CacheAlreadyExists
	}
	self.RUnlock()
	self.Lock()
	self.Caches[id] = cache2go.Cache(id)
	self.Unlock()
	return nil
}

func (self *CacheLayer) Add(id string, key string, data []byte) error {
	self.RLock()
	if _, found := self.Caches[id]; !found {
		return CacheNotFound
	}
	self.RUnlock()
	self.Lock()
	self.Caches[id].Add(key, 0, data)
	self.Unlock()
	return nil
}

func (self *CacheLayer) Delete(id string, key string) error {
	self.RLock()
	if _, found := self.Caches[id]; !found {
		return CacheNotFound
	}
	self.RUnlock()
	self.Lock()
	self.Caches[id].Delete(key)
	self.Unlock()
	return nil
}

func (self *CacheLayer) Get(id string, key []byte) ([]byte, error) {
	self.RLock()
	if _, found := self.Caches[id]; !found {
		data := &common.Item{}
		err := self.db.Get(common.Item{Key: key}, data)
		return data.Value, err
	}
	self.RUnlock()
	self.Lock()
	data, err := self.Caches[id].Value(key)
	self.Unlock()
	return data.Data().([]byte), err
}

func (self *CacheLayer) push(key interface{}, item *cache2go.CacheItem) {
	self.db.Put(common.Item{Key: key.([]byte), Value: item.Data().([]byte), Sync: false})
}

func (self *CacheLayer) Commit(id string) error {
	self.Caches[id].Foreach(self.push)
	return nil

}
