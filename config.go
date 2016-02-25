package main

import (
	"encoding/json"
	"github.com/rhino1998/cluster/forward"
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
	Compute   bool                       `json:"compute"`
}

var (
	Config Configuration = getConfig("./conf.json")
)

func getConfig(path string) Configuration {
	file, err := os.Open(path)
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
