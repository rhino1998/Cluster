package main

import (
	"3_1/everything"
	"github.com/rhino1998/lib/client"
	"github.com/zond/god/client"
	"os"
	"runtime"
	"strings"
	//"strconv"
)

type info struct {
	personID []byte
	rideID   []byte
	times    []byte
}

func main() {
	//c := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(1)
	conn := client.MustConn("")
	argsWithoutProg := os.Args[1:]
	var i info
	common.MustJSONDecode([]byte(argsWithoutProg[0]), &i)
	s := strings.Split(string(i.times), "id:")[1]
	s = strings.Split(s, "}")[0]
	var time []int
	v, _ := conn.Get([]byte(s))
	common.MustJSONDecode(v, &time)
	p, _ := everything.GetG(append([]byte{0}, []byte{0}[0], []byte{0}[0]))
	p.Times = append(p.Times, time)

	var ri []byte
	v, _ = conn.Get([]byte(s + "r"))
	common.MustJSONDecode(v, &ri)
	p.AppointRide = append(p.AppointRide, ri)

	p.Save()
}
