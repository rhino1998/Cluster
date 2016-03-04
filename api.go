package main

import (
	"encoding/json"
	"github.com/rhino1998/cluster/tasks"
	"io"
	"log"
	"net/http"
	"strconv"
)

type newtask struct {
	Id    []byte         `json:"id"`
	Jumps map[string]int `json:"jumps"`
	Name  string         `json:"name"`
	//Reqs     //[]reqs.Req
	FileName string `json:"filename"`
	Loc      string `json:"loc"`
	Value    int    `json:"value"`
	Args     string `json:"args"`
}

func api_peers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		body, _ := json.Marshal(This.Peers.Items())
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		io.WriteString(w, string(body))
	} else {
		log.Println("405 peers")
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func api_task(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var task newtask
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&task)
		if err != nil {
			log.Println("405 task", r.Body)
			http.Error(w, "", http.StatusNotAcceptable)
		}
		This.NewTask(tasks.NewTask(task.Name, task.Loc, task.FileName, task.Args, task.Value))
	}
}
