package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	addr string
}

type getreturn struct {
	Key    string `json:"Key"`
	Value  string `json:"Value"`
	Exists bool   `json:"Exists"`
}

func MustConn(addr string) {
	return &Client{addr: addr}
}

func (self *Client) Get(url string, key []byte) ([]byte, bool) {
	var jsonStr = []byte(fmt.Sprintf(`{"Key":"%v"}`, base64.StdEncoding.EncodeToString(key)))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%v/rpc/DHash.Get", url), bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var dat getreturn
	json.Unmarshal(body, &dat)
	val, _ := base64.StdEncoding.DecodeString(dat.Value)
	return val, dat.Exists
}

func (self *Client) Put(url string, key, val []byte) {
	var jsonStr = []byte(fmt.Sprintf(`{"key":"%v", "Sync": false, "Value":"%v"}`, base64.StdEncoding.EncodeToString(key), base64.StdEncoding.EncodeToString(val)))
	fmt.Println(string(jsonStr))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%v/rpc/DHash.Put", url), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Length", string(len(string(jsonStr))))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	resp.Body.Close()
}
