package main

import (
	//"crypto/sha1"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/rhino1998/cluster/util"
)

func read_int64(data []byte, big bool) (ret uint64) {
	if big {
		for i, b := range data {
			ret |= uint64(b) << uint((i)*8)
		}
	} else {
		l := len(data)
		for i, b := range data {
			ret |= uint64(b) << uint((l-i-1)*8)
		}
	}

	return ret
}

func comp(key []byte, addr string) (ret uint64) {
	ipstring, portstring, err := net.SplitHostPort(addr)
	if err != nil {
		return ^uint64(0)
	}
	port, err := strconv.Atoi(portstring)
	if err != nil {
		return ^uint64(0)
	}
	ip := net.ParseIP(ipstring)
	if ip == nil {
		return ^uint64(0)
	}
	return read_int64(key[:6], false) ^ (read_int64([]byte(ip)[12:], true)<<16 | uint64(port))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	addrs := make(chan string, 20)
	for i := 0; i < 20; i++ {
		addrs <- fmt.Sprintf("%v.%v.%v.%v:%v", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(65535))
	}
	fmt.Println(len(addrs))
	dat := util.PadKey([]byte{254}, 6)
	best := ^uint64(0)
	bestaddr := <-addrs
	addrs <- bestaddr
	for i := 0; i < len(addrs); i++ {
		addr := <-addrs
		addrs <- addr

		val := comp(dat[:6], addr)
		if val < best {
			bestaddr = addr
			best = val
		}
	}
	fmt.Println(best, bestaddr)

}
