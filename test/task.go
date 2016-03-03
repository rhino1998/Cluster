package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	c := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(1)
	for ; c > 0; c-- {
		go (func() {
			for {
				time.Now()
			}
		})()
	}
	<-time.After(60 * time.Second)
	fmt.Print(time.Now())
}
