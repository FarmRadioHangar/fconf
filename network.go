package main

import (
	"errors"

	"github.com/coreos/go-systemd/unit"
)

//Static holds information for static IP address.
type Static struct {

	// This is a string representation of an IP address.
	//
	// The IP adress should also contain the mask, in the format
	//xxx.xxx.xxx.xxx/xx
	IP      string `json:"ip"`
	Gateway string `json:"gateway"`
}

//Network configuration settings
type Network struct {
	Static    *Static  `json:"static"`
	DHCP      bool     `json:"dhcp"`
	DNS       []string `json:"dns-servers"`
	Interface string   `json:"interface"`
	Enabled   bool     `json:"enabled"`
}

//ToSystemdUnit transforms the Network object to systemd unit file.
func (e Network) ToSystemdUnit() ([]*unit.UnitOption, error) {
	var result []*unit.UnitOption
	// network interface
	i := &unit.UnitOption{
		Section: "Match",
		Name:    "Name",
		Value:   "eth0",
	}
	if e.Interface != "" {
		i.Value = e.Interface
	}
	result = append(result, i)
	if e.Static != nil {
		// add the IP
		result = append(result, &unit.UnitOption{
			Section: "Network",
			Name:    "Address",
			Value:   e.Static.IP,
		})

		if !e.DHCP {
			// Gateway
			result = append(result, &unit.UnitOption{
				Section: "Network",
				Name:    "Gateway",
				Value:   e.Static.Gateway,
			})
		} else {
			result = append(result, &unit.UnitOption{
				Section: "Network",
				Name:    "DHCP",
				Value:   "ipv4",
			})
		}
		// DNS
		for _, v := range e.DNS {
			result = append(result, &unit.UnitOption{
				Section: "Network",
				Name:    "DNS",
				Value:   v,
			})
		}
		return result, nil
	}

	if e.DHCP {
		result = append(result, &unit.UnitOption{
			Section: "Network",
			Name:    "DHCP",
			Value:   "ipv4",
		})
		if e.DNS != nil {
			for _, v := range e.DNS {
				result = append(result, &unit.UnitOption{
					Section: "Network",
					Name:    "DNS",
					Value:   v,
				})
			}
		}
		return result, nil
	}

	return nil, errors.New("at least either static should specifid or dhcp")
}
