package main

import (
	"fmt"
	"github.com/rhino1998/cluster/tasks"
	"log"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:3000")
	var wg sync.WaitGroup
	if err != nil {
		panic(err)
	}
	var reply []byte
	var a int64 = 0
	now := time.Now()
	for i := 0; i < 20000; i++ {
		wg.Add(1)
		//time.Sleep(20 * time.Millisecond)
		go func() {
			task := tasks.NewTask("FLOOP", "http://localhost:8080/task2.exe", "task2.exe", []string{"hey"}, 3400)
			err = client.Call("Node.AllocateTask", &task, &reply)
			if err != nil || string(reply) != "yo\n" {
				log.Fatal(err, string(reply))
			}
			wg.Done()
			atomic.AddInt64(&a, 1)
			log.Println(a)
		}()
	}
	log.Println("alloc", float64(time.Since(now).Nanoseconds()/1000000))
	wg.Wait()
	fmt.Print(string(reply), float64(time.Since(now).Nanoseconds()/1000000)/20000)
}
