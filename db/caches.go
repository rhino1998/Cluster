package db

import (
	"errors"
	"github.com/muesli/cache2go"
	"github.com/zond/god/dhash"
	"sync"
)

var (
	CacheNotFound      error = errors.New("Cache not found")
	CacheAlreadyExists error = errors.New("Cache already exists")
)

type CacheLayer struct {
	sync.RWMutex
	Caches map[string]*cache2go.CacheTable
	DB     *dhash.Node
}

func NewCacheLayer(node *dhash.Node) *CacheLayer {
	return &CacheLayer{Caches: make(map[string]*cache2go.CacheTable), DB: node}
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

func (self *CacheLayer) AddCacheItem(id string, key string, data []byte) error {
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

func (self *CacheLayer) DeleteCacheItem(id string, key string) error {
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

func (self *CacheLayer) GetCacheItem(id string, key string) ([]byte, error) {
	self.RLock()
	if _, found := self.Caches[id]; !found {
		return nil, CacheNotFound
	}
	self.RUnlock()
	self.Lock()
	data, err := self.Caches[id].Value(key)
	self.Unlock()
	return data.Data().([]byte), err
}
