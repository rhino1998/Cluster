package querymap

import (
	"sync"
)

type QueryMap struct {
	sync.RWMutex
	data map[string]map[string]interface{}
}

//Makes a new QueryMap
func New() *QueryMap {
	return &QueryMap{data: make(map[string]map[string]interface{})}
}

//Creates a row
func (self *QueryMap) Assign(key string, data map[string]interface{}) {
	self.Lock()
	defer self.Unlock()
	self.data[key] = data
}

//Deletes a row
func (self *QueryMap) Delete(key string) {
	self.Lock()
	defer self.Unlock()
	delete(self.data, key)
}

func (self *QueryMap) Get(key string) (map[string]interface{}, bool) {
	self.RLock()
	defer self.RUnlock()
	data, ok := self.data[key]
	return data, ok
}

//Assigns a value in row key column subkey
func (self *QueryMap) SubAssign(key, subkey string, data interface{}) {
	self.Lock()
	defer self.Unlock()
	if _, ok := self.data[key]; !ok {
		self.data[key] = make(map[string]interface{})
	}
	self.data[key][subkey] = data
}

//Removes a value in row key column subkey
func (self *QueryMap) SubDelete(key, subkey string) {
	self.Lock()
	defer self.Unlock()
	delete(self.data[key], subkey)
}

func (self *QueryMap) SubGet(key, subkey string) (interface{}, bool) {
	self.RLock()
	defer self.RUnlock()
	data, ok := self.data[key][subkey]
	return data, ok
}

//Returns a QueryMap containing only items for which comp returns true
func (self *QueryMap) Mask(subkey string, comp func(value, check interface{}) (bool, error), value interface{}) (*QueryMap, error) {
	self.RLock()
	defer self.RUnlock()
	newdata := make(map[string]map[string]interface{})
	for key, row := range self.data {
		if ok, err := comp(row[subkey], value); ok && err == nil {
			if _, ok := newdata[key]; !ok {
				newdata[key] = make(map[string]interface{})
			}
			newdata[key][subkey] = row[subkey]
		} else if err != nil {
			return nil, err
		}
	}
	return &QueryMap{data: newdata}, nil
}

func (self *QueryMap) Length() int {
	self.RLock()
	defer self.RUnlock()
	return len(self.data)
}

//Returns QueryMap as a map of strings to a map of strings to interfaces
func (self *QueryMap) Items() map[string]map[string]interface{} {
	self.RLock()
	defer self.RUnlock()
	newdata := make(map[string]map[string]interface{})
	for key, row := range self.data {
		newdata[key] = make(map[string]interface{})
		for subkey, data := range row {
			newdata[key][subkey] = data
		}
	}
	return newdata
}

//Returns a slice of all row keys in QueryMap
func (self *QueryMap) Keys() []string {
	keys := make([]string, 0, len(self.data))
	for k := range self.data {
		keys = append(keys, k)
	}
	return keys
}
