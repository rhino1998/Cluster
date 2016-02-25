package tasks

import (
	"github.com/rhino1998/cluster/reqs"
)

type Task struct {
	Id       []byte
	Jumps    map[string]int
	Name     string
	Reqs     []reqs.Req
	FileName string
	Loc      string
}
