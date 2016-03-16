package main

import (
	"fmt"
	//"github.com/rhino1998/cluster/tasks"
	"crypto/sha1"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	laddr, err := net.ResolveTCPAddr("tcp", "localhost")
	raddr, err := net.ResolveTCPAddr("tcp", "localhost:2002")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", laddr, raddr)
	client := rpc.NewClient(conn)
	var wg sync.WaitGroup
	if err != nil {
		panic(err)
	}
	var reply bool
	rand.Seed(time.Now().UnixNano())
	key := randSeq(16)
	data := sha1.Sum([]byte(key))
	log.Println(data, key)
	//task := tasks.NewTask("FLOOP", "./", "task2.exe", []string{"hey"}, 3400)
	now := time.Now()
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			err = client.Call("Node.Put", data[:], &reply)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	log.Println("alloc", float64(time.Since(now).Nanoseconds()/1000000))
	wg.Wait()
	fmt.Print(reply, float64(time.Since(now).Nanoseconds()/1000000))
}
