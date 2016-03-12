package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type getreturn struct {
	Key    string `json:"Key"`
	Value  string `json:"Value"`
	Exists bool   `json:"Exists"`
}

func get(url string, key []byte) ([]byte, bool) {
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

func put(url string, key, val []byte) {
	var jsonStr = []byte(fmt.Sprintf(`{"key":"%v", "Sync": false, "Value":"%v"}`, base64.StdEncoding.EncodeToString(key), base64.StdEncoding.EncodeToString(val)))
	fmt.Println(string(jsonStr))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%v/rpc/DHash.Put", url), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Length", string(len(string(jsonStr))))
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	resp.Body.Close()
}

// /import "os"

func main() {
	//url := "http://108.56.251.125:2004"

	//put(url, []byte("tes"), []byte(`[{"Ride":["ew","2erw"]},{"Ride":[1,2]},{"Ride":[1,2ddff]}]`))
	val, _ := get(url, []byte("tes"))
	fmt.Println(string(val))
}
