package main

import (
	"3_1/everything"
	"fmt"
	"strconv"
	"runtime"
	"github.com/zond/god/common"
	"github.com/zond/god/client"

)

func getTime(a int) string {
	b := 1 + a/3600
	a -= 3600 * (a / 3600)
	c := a/60
	s := strconv.Itoa(b)
	s2 := strconv.Itoa(c)
	return s + ":" + s2

}

func between(a1, b1, a2, b2 int) bool {
	if a1 > a2 && a1 < b2 {
		return true
	}
	if b1 > a2 && b1 < b2 {
		return true
	}
	av := a1 + b1 / 2
	if av > a2 && av < b2 {
		return true
	}
	return false
}

type info struct {
	personID []byte
	rideID []byte
}

func main() {
	
	c := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(1)
	conn := client.MustConn("")
	argsWithoutProg := os.Args[1:]
	var i info
	common.MustJSONDecode([]byte(argsWithoutProg[0]), &i)
	r := everything.RGet(i.rideID)
	fmt.Println(r.Name)
	p, _ := everything.GetG(append([]byte{0},[]byte{0}[0],[]byte{0}[0]))
	count := 0
	t := true
	fmt.Print("[")
	for i := r.CurTrip; i < len(r.Times); i++ {
		if r.Capacity - r.Capacities[i] > 2 {
			a := r.GetTimeRange(i)
			for j := 0; j < len(p.Times); j++ {
				if between(a[0], a[1], p.Times[j][0], p.Times[j][1]) {
					t = false
				}
			}
			if t {
				conn.Put([]byte(strconv.Itoa(i)),common.MustJSONEncode(p.Times[j]))
				conn.Put([]byte(strconv.Itoa(i) + "r"),common.MustJSONEncode(r.K)
				fmt.Print("{" + r.Name + ":" + " [" + getTime(p.Times[j][0]) + "," + getTime(p.Times[j][1]) + "],id:" + strconv.Itoa(i) +  "},")
				count += 1
			}
		}
		if count == 3 { break }
	}
	fmt.Print("]")


}