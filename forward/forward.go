package forward

import (
	"fmt"
	"github.com/syncthing/syncthing/lib/upnp"
	"log"
)

type Mapping struct {
	Protocols   []upnp.Protocol `json:"protocols"`
	Ports       []int           `json:"port"`
	Description string          `json:"description"`
}

func Forward(nat upnp.IGD, mapping Mapping) (err error) {
	for _, protocol := range mapping.Protocols {
		for _, port := range mapping.Ports {
			err := nat.AddPortMapping(protocol, port, port,
				fmt.Sprintf(mapping.Description, protocol, port), 0)
			log.Println("success")
			if err != nil {
				return err
			}
		}
	}
	return
}
