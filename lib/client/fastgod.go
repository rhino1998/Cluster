package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/rhino1998/god/common"
	"io/ioutil"
	"net/http"
)

type Conn struct {
	addr string
}

type getreturn struct {
	Key    string `json:"Key"`
	Value  string `json:"Value"`
	Exists bool   `json:"Exists"`
}

func MustConn(addr string) *Conn {
	return &Conn{addr: addr}
}

func (self *Conn) Get(key []byte) ([]byte, bool) {
	var jsonStr = []byte(fmt.Sprintf(`{"Key":"%v"}`, string(key)))
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v/rpc/DHash.Get", self.addr), bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var dat getreturn
		json.Unmarshal(body, &dat)
		val, _ := base64.StdEncoding.DecodeString(dat.Value)
		return val, dat.Exists
	}
	return nil, false
}

func (self *Conn) Put(key, val []byte) {
	var jsonStr = []byte(fmt.Sprintf(`{"Key":"%v", "Sync": false, "Value":"%v"}`, string(key), string(val)))
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v/rpc/DHash.Put", self.addr), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Length", string(len(string(jsonStr))))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	resp.Body.Close()
}

func (self *Conn) SubPut(key, subkey, val []byte) {
	var jsonStr = []byte(fmt.Sprintf(`{"Key":"%v", "SubKey":"%v", "Sync": false, "Value":"%v"}`, string(key), string(subkey), string(val)))
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v/rpc/DHash.SubPut", self.addr), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Length", string(len(string(jsonStr))))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	resp.Body.Close()
}

func (self *Conn) Del(key []byte) {
	var jsonStr = []byte(fmt.Sprintf(`{"Key":"%v", "Sync": false}`, string(key)))
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v/rpc/DHash.Del", self.addr), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Length", string(len(string(jsonStr))))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	resp.Body.Close()
}

func (self *Conn) SubDel(key, subkey []byte) {
	var jsonStr = []byte(fmt.Sprintf(`{"Key":"%v", "SubKey":"%v", "Sync": false}`, string(key), string(subkey)))
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v/rpc/DHash.SubDel", self.addr), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Length", string(len(string(jsonStr))))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	resp.Body.Close()
}

func (self *Conn) Slice(key, min, max []byte, mininc, maxinc bool) []common.Item {
	var jsonStr = []byte(fmt.Sprintf(`{"Key":"%v", "Min":"%v", "Max":"%v", "Sync": false, "MinInc": %v, "MaxInc": %v}`, string(key), string(min), string(max), mininc, maxinc))
	fmt.Println(string(jsonStr))
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v/rpc/DHash.Slice", self.addr), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Length", string(len(string(jsonStr))))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var dat []common.Item
		json.Unmarshal(body, &dat)
		return dat
	}
	return nil
}

func (self *Conn) Clear() {
	var jsonStr = []byte(fmt.Sprintf(`{}`))
	fmt.Println(string(jsonStr))
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v/rpc/DHash.Clear", self.addr), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Length", string(len(string(jsonStr))))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	resp.Body.Close()
}
