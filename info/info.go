package info

import (
	"github.com/rhino1998/cluster/bench"
)

type Info struct {
	Compute bool `json:"compute"`
	IntAddr string
	bench.Specs
}
