package main

import (
	"encoding/json"
	"github.com/gorilla/context"
	"io"
	"net/http"
	"strconv"
	//"github.com/nu7hatch/gouuid"
	"fmt"
	"log"
)

type Status struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

type Job struct {
	ID     string            `json:"id" bson:"_id"`
	Name   string            `json:"name"`
	OS     map[string]string `json:"os"`
	TaskID string            `json:"task" bson:"task"`
}

type Server struct {
}

func NewServer(config Configuration) (*Server, error) {
	//if err != nil { return nil, err }
	return &Server{}, nil
}

func (s *Server) WithData(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*if _, present := context.GetOk(r, "db"); !present {
			mongocopy := s.mongosession.Copy()
			redisclient := s.redisclient
			config := s.config
			defer mongocopy.Close()
			context.Set(r, "mongo", mongocopy)
			context.Set(r, "redis", redisclient)
			context.Set(r, "config", config)

		}*/
		fn(w, r)
	}
}

func (s *Server) Close() {
}

func task(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var job Job
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&job)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		NewTasks.NotFoundAdd(job.ID, 0, job)
		body, err := json.Marshal(Status{ID: job.ID, State: "Recieved"})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(string(body))))
		io.WriteString(w, string(body))
	} else {
		log.Println("405")
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func startServer() {
	//loggerSetup("server.log")

	srv, err := NewServer(Config)
	errpanic(err)
	defer srv.Close()
	http.HandleFunc("/", srv.WithData(task))
	errpanic(http.ListenAndServe(fmt.Sprintf("%s:8000", localIP), context.ClearHandler(http.DefaultServeMux)))
}
