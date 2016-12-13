package main

//Static holds information for static IP address.
type Static struct {

	// This is a string representation of an IP address.
	//
	// The IP adress should also contain the mask, in the format
	//xxx.xxx.xxx.xxx/xx
	IP      string `json:"ip"`
	Gateway string `json:"gateway"`
}

//Ethernet configuration settings for ethernet network.
type Ethernet struct {
	Static *Static  `json:"static"`
	DHCP   bool     `json:"dhcp"`
	DNS    []string `json:"dns-servers"`
}
