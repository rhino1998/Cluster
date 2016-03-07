package jsonrpc

import (
	"bytes"
	"fmt"
	"github.com/gorilla/rpc/json"
	"net"
	//"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	Addr   string
}

func NewClient(addr string) *Client {
	return &Client{
		Addr: fmt.Sprintf("http://%v/rpc", addr),
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				Dial: (&net.Dialer{
					Timeout:   500 * time.Millisecond,
					KeepAlive: 60 * time.Second,
				}).Dial,
				TLSHandshakeTimeout:   5 * time.Second,
				ExpectContinueTimeout: 500 * time.Millisecond,
			},
		},
	}
}

func (self *Client) Call(service string, args interface{}, reply interface{}) error {
	message, err := json.EncodeClientRequest(service, args)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", self.Addr, bytes.NewBuffer(message))
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
