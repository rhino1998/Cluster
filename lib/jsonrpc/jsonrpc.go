package jsonrpc

import (
	"bytes"
	"fmt"
	"github.com/gorilla/rpc/json"
	"log"
	"net/http"
)

type Client struct {
	Addr string
}

func NewClient(addr string) *Client {
	return &Client{Addr: fmt.Sprintf("http://%v/rpc", addr)}

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
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.DecodeClientResponse(resp.Body, reply)
	log.Println(reply)
	if err != nil {
		log.Fatalf("Couldn't decode response. %s", err)
		return err
	}
	return nil
}
