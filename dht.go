package main

import (
	"fmt"
	"github.com/rhino1998/god/dhash"
)

var (
	KVStore *dhash.Node
)

func initDHT() {
	fmt.Println(This.Addr)
	KVStore = dhash.NewNodeDir(fmt.Sprintf("%v:%v", "0.0.0.0", Config.Mappings["DHT"].Port), fmt.Sprintf("%v:%v", This.Addr, Config.Mappings["DHT"].Port), "")
	KVStore.Start()
	if Config.DHTSeed != "" {
		KVStore.MustJoin(Config.DHTSeed)
	}
}
