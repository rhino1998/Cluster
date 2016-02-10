package forward

import (
	"fmt"
	"github.com/syncthing/syncthing/lib/upnp"
)

type Mapping struct {
	Protocols   []upnp.Protocol `json:"protocols"`
	Port        int             `json:"port"`
	Description string          `json:"description"`
}

func Forward(nat upnp.IGD, mapping Mapping) (err error) {
	for _, protocol := range mapping.Protocols {
		err := nat.AddPortMapping(protocol, mapping.Port, mapping.Port,
			fmt.Sprintf(mapping.Description, protocol, mapping.Port), 0)
		if err != nil {
			return err
		}
	}
	return
}

func ForwardAll(nat upnp.IGD, mappings map[string]Mapping) (err error) {
	for _, mapping := range mappings {
		err := Forward(nat, mapping)
		if err != nil {
			return err
		}
	}
	return
}
