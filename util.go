package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func getLocalIP() (ip net.IP, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				if !ip.IsLoopback() && ip.String() != "0.0.0.0" {
					return ip
				}
			case *net.IPAddr:
				ip = v.IP
				if !ip.IsLoopback() && ip.String() != "0.0.0.0" {
					return ip
				}
			}
		}
	}
	return nil
}

//replace this eventually
func getExternalIP() (ip net.IP, err error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	ip, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return net.ParseIP(string(ip)[:len(string(ip))-1])
}
