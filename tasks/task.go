package tasks

import (
//"github.com/rhino1998/cluster/reqs"
)

type Task struct {
	Id []byte `json:"id"`
	//Jumps map[string]int `json:"jumps"`
	Name string `json:"name"`
	//Reqs     //[]reqs.Req
	FileName string `json:"filename"`
	Loc      string `json:"loc"`
}
