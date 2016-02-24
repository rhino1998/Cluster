package main

import (
	"github.com/muesli/cache2go"
	"log"
	"sync/atomic"
	"time"
)

var (
	Tasks        int64                = 0
	NewTasks     *cache2go.CacheTable = cache2go.Cache("New")
	ManagedTasks *cache2go.CacheTable = cache2go.Cache("Monitor")
)

func Allocate(item *cache2go.CacheItem) {
	log.Println("Allocated", item.Key())
	log.Println(NewTasks.Count())
	log.Println(atomic.AddInt64(&Tasks, 1))
	ManagedTasks.NotFoundAdd(item.Key(), time.Duration(Config.Timeout)*time.Nanosecond*1e6, item.Data())
	//go func(){
	time.Sleep(3 * time.Second)
	log.Println("nonblock")
	//}()
}
