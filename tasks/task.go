package tasks

import (
	//"github.com/rhino1998/cluster/reqs"
	"github.com/rhino1998/cluster/util"
)

type Task struct {
	Id    []byte         `json:"id"`
	Jumps map[string]int `json:"jumps"`
	Name  string         `json:"name"`
	//Reqs     //[]reqs.Req
	FileName string `json:"filename"`
	Loc      string `json:"loc"`
	Value    int    `json:"value"`
}

func NewTask(name, loc, filename string, value int) Task {
	return Task{Id: []byte(util.NewUUID()), Jumps: make(map[string]int), FileName: filename, Name: name, Loc: loc}
}
