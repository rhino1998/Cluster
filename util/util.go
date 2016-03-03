package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

func GetLocalIP() (ip net.IP, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				if !ip.IsLoopback() && ip.String() != "0.0.0.0" {
					return ip, nil
				}
			case *net.IPAddr:
				ip = v.IP
				if !ip.IsLoopback() && ip.String() != "0.0.0.0" {
					return ip, nil
				}
			}
		}
	}
	return nil, err
}

//replace this eventually
func GetExternalIP() (ip net.IP, err error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ip, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(string(ip)[:len(string(ip))-1]), nil
}

func ByteSliceEq(a, b []byte) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func NewUUID() string {
	uuid := make([]byte, 16)
	io.ReadFull(rand.Reader, uuid)
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
