package main

import "fmt"
import "github.com/rhino1998/cluster/lib/client"

// /import "os"

func main() {
	conn := client.MustConn("108.56.251.125:2004")
	val, _ := conn.Get([]byte("test"))
	fmt.Println(string(val))
}
