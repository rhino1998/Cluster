package bench

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Specs struct {
	Threads int `json:"threads"`
	RAM     int `json:"RAM"`
}

func LoadSpecs(path string) (*Specs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	specs := Specs{}
	err = json.Unmarshal(data, &specs)
	if err != nil {
		return nil, err
	}
	return &specs, nil
}
