package main

import (
	"fmt"
	"github.com/zond/god/dhash"
)

var (
	KVStore *dhash.Node = dhash.NewNodeDir(fmt.Sprintf("%v:%v", "0.0.0.0", Config.Mappings["DHT"].Port), fmt.Sprintf("%v:%v", This.ExternalIP, Config.Mappings["DHT"].Port), "")
)

func initDHT() {
	KVStore.Start()
	if Config.DHTSeed != "" {
		KVStore.MustJoin(Config.DHTSeed)
	}
}
