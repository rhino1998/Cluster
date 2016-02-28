package main

import (
	"encoding/json"
	"github.com/rhino1998/cluster/tasks"
	"io"
	"log"
	"net/http"
	"strconv"
)

type dbgetresp struct {
	data  string `json:"data"`
	found bool   `json:"found"`
}

type dbget struct {
	id        string `json:"id"`
	key       string `json:"key"`
	skipcache bool   `json:"skipcache"`
}
type dbput struct {
	id        string                      `json:"id"`
	key       string                      `json:"key"`
	skipcache bool                        `json:"skipcache""`
	data      map[interface{}]interface{} `json:"data"`
}
type dbdel struct {
	id        string `json:"id"`
	key       string `json:"key"`
	skipcache bool   `json:"skipcache""`
}
type commit struct {
	id string `json:"id"`
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

func api_db_put(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var put dbput
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&put)
		if err != nil {
			log.Println("405 dbset", r.Body)
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
		body, err := json.Marshal(put.data)
		if err != nil {
			log.Println("405 dbset", r.Body)
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
		if put.skipcache {
			This.DB.DBPut(put.id, put.key, body)
		} else {
			This.DB.Add(put.id, put.key, body)
		}
	}
}

func api_db_get(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var get dbget
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&get)
		if err != nil {
			log.Println("405 dbput", r.Body)
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
		var data []byte
		var found bool
		if get.skipcache {
			data, found = This.DB.DBGet(get.id, get.key)

		} else {
			data, found = This.DB.Get(get.id, get.key)
		}
		body, err := json.Marshal(dbgetresp{data: string(data), found: found})
		if err != nil {
			log.Println("405 dbset", r.Body)
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
		io.WriteString(w, string(body))

	}
}

func api_db_del(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var del dbdel
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&del)
		if err != nil {
			log.Println("405 dbset", r.Body)
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
		if del.skipcache {
			This.DB.Del(del.id, del.key)
		} else {
			This.DB.DBDel(del.id, del.key)
		}
	}
}

func api_task(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var task tasks.Task
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&task)
		if err != nil {
			log.Println("405 task", r.Body)
			http.Error(w, "", http.StatusNotAcceptable)
		}
		This.NewTask(task)
	}
}
