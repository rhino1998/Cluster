package jsonrpc

import (
	"bytes"
	"fmt"
	"github.com/gorilla/rpc/json"
	//"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var client *http.Client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 0,
		Proxy:               http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   500 * time.Millisecond,
			KeepAlive: 0,
			LocalAddr: nil,
		}).Dial,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 500 * time.Millisecond,
	},
}

type Client struct {
	client *http.Client
	addr   string
}

func NewClient(laddr, raddr string) *Client {
	laddress, _ := net.ResolveTCPAddr("tcp", laddr)
	laddress.IP = net.ParseIP("localhost")
	return &Client{
		addr:   fmt.Sprintf("http://%v/rpc", raddr),
		client: client,
	}
}

func (self *Client) Call(service string, args interface{}, reply interface{}) error {
	message, err := json.EncodeClientRequest(service, args)
	//log.Println(string(message))
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", self.addr, bytes.NewBuffer(message))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		return err
	}
	/*var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}
	// Restore the io.ReadCloser to its original state
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// Use the content
	bodyString := string(bodyBytes)
	log.Println(bodyString)*/
	defer resp.Body.Close()
	err = json.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		log.Println("Couldn't decode response. %s", err)
		return err
	}
	return nil
}
