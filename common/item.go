package common

import (
	"github.com/spaolacci/murmur3"
)

type Item struct {
	Key  string `json:"key"`
	Data []byte `json:"data"`
}

func (self *Item) Value() uint64 {
	return murmur3.Sum64([]byte(self.Key)) >> 16
}
