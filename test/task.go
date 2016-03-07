package main

import (
	"fmt"
	"os"
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
	<-time.After(5000 * time.Millisecond)
	fmt.Print(os.Args)
}
