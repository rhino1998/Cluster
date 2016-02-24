package task

import (
	"cluster/reqs"
)

type Task struct {
	Id   []byte
	Name string
	Reqs []reqs.Req
	Loc  string
}
