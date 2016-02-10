package main

import (
	"cluster/peers"
	"github.com/zond/god/common"
	"log"
	"time"
)

var (
	This peers.Peer = peers.Peer{
		ExternalIP: getExternalIP.ipErrHandler(),
		LocalIP:    getLocalIP.ipErrHandler(),
	}
)

/*func GetSpecs() bench.Specs {
	return
}*/

func alive() {
	if err != nil {
		log.Printf("Alive Bump Failed: %v", err)
		return
	}
	KVStore.SubPut(common.Item{
		Key:    []byte("blank"),
		SubKey: []byte(This.ExternalIP.String()),
		Value:  This,
	})
}

func aliveLoop() {
	ticker := time.NewTicker(time.Duration(Config.Timeout) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				alive()
			}
		}
	}()
}
