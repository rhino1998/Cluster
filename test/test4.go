package main

import (
	//"fmt"
	"github.com/rhino1998/cluster/common"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3002")
	client := rpc.NewClient(conn)
	var wg sync.WaitGroup
	if err != nil {
		panic(err)
	}
	var puttimer int64 = 0
	var reply []byte
	now := time.Now()
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(i int) {
			temp := time.Now()
			item := &common.Item{Key: strconv.Itoa(i), Data: []byte("randrandrd")}
			err = client.Call("Node.Put", item, nil)
			if err != nil {
				log.Fatal(err)
			}
			atomic.AddInt64(&puttimer, int64(time.Since(temp)))
			wg.Done()
		}(i)
	}
	wg.Wait()
	var gettimer int64 = 0
	var wg2 sync.WaitGroup
	for i := 0; i < 100000; i++ {
		wg2.Add(1)
		go func(i int) {
			temp := time.Now()
			stri := strconv.Itoa(i)
			err = client.Call("Node.Get", &stri, &reply)
			if err != nil || string(reply) != "randrandrd" {
				log.Fatal(err, string(reply))
			}
			log.Println(string(reply))
			atomic.AddInt64(&gettimer, int64(time.Since(temp)))
			wg2.Done()
		}(i)
	}
	log.Println("alloc", float64(time.Since(now).Nanoseconds()/1000000))
	wg2.Wait()
	log.Println(puttimer/1000000/20, gettimer/1000000/20)
}
