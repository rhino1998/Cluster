package main

import (
	"cluster/forward"
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
)

type Configuration struct {
	Timeout   int                        `json:"timeout"`
	Mappings  map[string]forward.Mapping `json:"mappings"`
	Forwarded bool                       `json:"forwarded"`
	PeerSeed  string                     `json:"peerseed"`
	DHTSeed   string                     `json:"dhtseed"`
	MaxTasks  int                        `json:"maxtasks"`
	External  string                     `json:"externalip"`
}

var (
	Config Configuration = getConfig(".")
)

func getConfig(path string) Configuration {
	file, err := os.Open(path + "/conf.json")
	if err != nil {
		panic(err)
	}
	comment, err := regexp.Compile("\\/\\/.*\\n")
	conf, err := ioutil.ReadAll(file)
	conf = comment.ReplaceAll(conf, []byte(""))
	configuration := Configuration{}
	err = json.Unmarshal(conf, &configuration)
	if err != nil {
		panic(err)
	}
	return configuration
}
