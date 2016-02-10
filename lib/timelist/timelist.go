package timelist

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Item struct {
	val   interface{}
	index time.Time
}

func (self Item) String() string {
	return fmt.Sprintf("%v:%v", self.index, self.val)
}

func (self *Item) Value() interface{} {
	return self.val
}

func (self *Item) Time() time.Time {
	return self.index.UTC()
}

type TimeList struct {
	sync.RWMutex
	vals  []Item
	start time.Time
	end   time.Time
}

func New() *TimeList {
	t := new(TimeList)
	t.start = time.Now().UTC()
	t.end = time.Now().UTC()
	return t
}

func (self *TimeList) Length() int {
	self.RLock()
	defer self.RUnlock()
	return len(self.vals)
}

func (self *TimeList) First() *Item {
	self.RLock()
	defer self.RUnlock()
	return &self.vals[0]
}

func (self *TimeList) Last() *Item {
	self.RLock()
	defer self.RUnlock()
	return &self.vals[len(self.vals)]
}

func (self *TimeList) Start() time.Time {
	self.RLock()
	defer self.RUnlock()
	return self.start.UTC()
}

func (self *TimeList) End() time.Time {
	self.RLock()
	defer self.RUnlock()
	return self.end.UTC()
}

//Retuns contents of ItemList as []Item
func (self *TimeList) Items() []Item {
	self.RLock()
	defer self.RUnlock()
	return self.vals
}

func (self *TimeList) frontSearch(index time.Time) int {
	nextIndex := len(self.vals) - 1
	for i, item := range self.vals {
		if index.Equal(item.index) {
			return i
		}
		if index.Before(item.index) {
			return i
		}
		nextIndex = i
	}
	return nextIndex
}

func (self *TimeList) backSearch(index time.Time) int {
	prevIndex := 0
	for i := len(self.vals) - 1; i >= 0; i-- {
		item := self.vals[i]
		if index.Equal(item.index) {
			return i
		}
		if index.After(item.index) {
			return i
		}
		prevIndex = i
	}
	return prevIndex
}

func (self *TimeList) search(index time.Time) int {
	if index.After(self.start) && index.Before(self.end) {
		if index.Sub(self.start) < self.end.Sub(index) {
			return self.backSearch(index) + 1
		}
		return self.frontSearch(index)
	}
	return -1
}

//Append item to end of TimeList
func (self *TimeList) Append(value interface{}) {
	self.Lock()
	defer self.Unlock()
	if len(self.vals) == cap(self.vals) {
		self.grow()
	}
	self.end = time.Now().UTC()
	newItem := Item{val: value, index: self.end}
	self.vals = append(self.vals, newItem)
	return
}

//Inserts value at given time
func (self *TimeList) Insert(value interface{}, index time.Time) {
	newItem := Item{val: value, index: index.UTC()}
	self.Lock()
	defer self.Unlock()
	if len(self.vals) == cap(self.vals) {
		self.grow()
	}
	if index.After(self.end) || index.Equal(self.end) {
		self.end = index
		self.vals = append(self.vals, newItem)

	} else if index.Before(self.start) || index.Equal(self.end) {
		self.start = index
		self.vals = append([]Item{newItem}, self.vals...)
	} else {
		i := self.search(index)
		self.vals = append(self.vals, Item{})
		copy(self.vals[i+1:], self.vals[i:])
		self.vals[i] = newItem
	}
	return
}

func (self *TimeList) PopAfter(index time.Time) (*TimeList, error) {
	self.Lock()
	defer self.Unlock()
	i := self.search(index)
	if i == -1 {
		return nil, errors.New(fmt.Sprintf("%v is outside range", index))
	}
	temp := make([]Item, len(self.vals)-i)
	copy(temp, self.vals[i:])
	self.end = self.vals[i].index
	self.vals = self.vals[:i]
	return &TimeList{vals: temp, start: temp[0].Time(), end: temp[len(temp)-1].Time()}, nil
}

func (self *TimeList) PopBefore(index time.Time) (*TimeList, error) {
	self.Lock()
	defer self.Unlock()
	i := self.search(index)
	if i == -1 {
		return nil, errors.New(fmt.Sprintf("%v is outside range", index))
	}
	temp := make([]Item, i)
	copy(temp, self.vals[:i])
	self.start = self.vals[i].index
	self.vals = self.vals[i:]
	return &TimeList{vals: temp, start: temp[0].Time(), end: temp[len(temp)-1].Time()}, nil
}

func (self *TimeList) RemoveAfter(index time.Time) {
	self.Lock()
	defer self.Unlock()
	i := self.search(index)
	if i == -1 {
		return
	}
	self.end = self.vals[i].index
	self.vals = self.vals[:i]
}

func (self *TimeList) RemoveBefore(index time.Time) {
	self.Lock()
	defer self.Unlock()
	i := self.search(index)
	if i == -1 {
		return
	}
	self.start = self.vals[i].index
	self.vals = self.vals[i:]
}

//Retuns a new ItemList containing items after given time
func (self *TimeList) After(index time.Time) *TimeList {
	self.RLock()
	defer self.RUnlock()
	i := self.search(index)
	if i == -1 {
		return nil
	}
	return &TimeList{vals: self.vals[i:], start: self.vals[i].index, end: self.end}
}

//Retuns a new ItemList containing items before given time
func (self *TimeList) Before(index time.Time) *TimeList {
	self.RLock()
	defer self.RUnlock()
	i := self.search(index)
	if i == -1 {
		return nil
	}
	return &TimeList{vals: self.vals[:i], start: self.start, end: self.vals[i].index}
}

//Retuns athe first value at or after given time
func (self *TimeList) FirstAfter(index time.Time) *Item {
	self.RLock()
	defer self.RUnlock()
	i := self.search(index)
	if i == -1 {
		return nil
	}
	return &self.vals[i]
}

//Retuns athe first value at or before given time
func (self *TimeList) FirstBefore(index time.Time) *Item {
	self.RLock()
	defer self.RUnlock()
	i := self.search(index)
	if i == -1 || i == 0 {
		return nil
	}
	return &self.vals[i-1]
}

func (self *TimeList) grow() {
	t := make([]Item, len(self.vals), (cap(self.vals)+1)*2)
	copy(t, self.vals)
	self.vals = t
	return
}
