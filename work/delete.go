package main


import (
	
	"runtime"
	"3_1/everything"
)

func main() {
	//c := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(1)
	p, _ := everything.GetG(append([]byte{0},[]byte{0}[0],[]byte{0}[0]))
	p = everything.Group{K: append([]byte{0},[]byte{0}[0],[]byte{0}[0])}
	p.Save()
}