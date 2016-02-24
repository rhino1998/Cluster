package jsonrpc

import (
	"bytes"
	"github.com/gorilla/rpc/json"
	"net/http"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr}

}

func (self *Client) Call(service string, args interface{}, reply interface{}) error {
	message, err := json.EncodeClientRequest(service, args)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", self.addr, bytes.NewBuffer(message))
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
	err = json.DecodeClientResponse(resp.Body, &reply)
	if err != nil {
		return err
	}
	return nil
}
