package main


import (
	
	"runtime"
	"3_1/everything"
	"strconv"
	"fmt"
)

func getTime(a int) string {
	b := 1 + a/3600
	a -= 3600 * (a / 3600)
	c := a/60
	s := strconv.Itoa(b)
	s2 := strconv.Itoa(c)
	return s + ":" + s2

}

func main() {
	//c := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(1)

	p, _ := everything.GetG(append([]byte{0},[]byte{0}[0],[]byte{0}[0]))
	fmt.Print("[")
	for i := 0; i < len(p.Times); i++ {
		r := everything.RGet(p.AppointRide[i])
		a1 := getTime(p.Times[i][0])
		b1 := getTime(p.Times[i][1])
		s1 := "{" + r.Name + ": [" 
		s2 := a1 + "," 
		s3 := b1 + "]},"
		fmt.Print(s1 + s2 + s3)
	}
	fmt.Print("]")
}

//"{" + r.Name + ": [" + a1 + "," b1 + "]},"