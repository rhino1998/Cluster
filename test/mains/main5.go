package main

import (
	"3_1/everything"
	"3_1/nm"
	"3_1/update"
	"fmt"
	"github.com/rhino1998/cluster/lib/client"
	"github.com/zond/god/common"
	"time"
)

func main() {
	conn := client.MustConn("108.56.251.125:2004")
	conn.Clear()
	g := everything.Group{K: append([]byte{0}, []byte{0}[0], []byte{0}[0])}
	conn.Put(g.K, common.MustJSONEncode(g))
	rides, _, _, _, _, _ := nm.MainR()

	T := 0
	//panic(0)
	fmt.Println("hi")
	for T < 36000*5 {
		fmt.Println(T)
		t := time.Now()
		T += 60
		conn.Put([]byte{83, 83, 83}, common.MustJSONEncode(T))
		for i := 0; i < len(rides); i++ {
			fmt.Println(i)
			update.Update(rides[i], T)
		}
		fmt.Println("jhfds")

		for time.Since(t)/time.Second < 3 {
			//fmt.Println(time.Since(t))
			//fmt.Println(time.Since(t)/time.Second)
			time.Sleep(500)
		}

	}

}
