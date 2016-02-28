package jsonrpc

import (
	"bytes"
	"fmt"
	"github.com/gorilla/rpc/json"
	"io/ioutil"
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
	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}
	// Restore the io.ReadCloser to its original state
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// Use the content
	bodyString := string(bodyBytes)
	log.Println(bodyString)
	defer resp.Body.Close()
	err = json.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		log.Println("Couldn't decode response. %s", err)
		return err
	}
	return nil
}
