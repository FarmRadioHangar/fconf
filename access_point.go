package main

//AccessPoint  is the access point configuration
type AccessPoint struct {
	Channel         string `ini:"CHANNEL"`
	Gateway         string `ini:"GATEWAY"`
	WPAVersion      int    `ini:"WPA_VERSION"`
	ETCHosts        int    `ini:"ETC_HOSTS"`
	DHCPDNS         string `ini:"DHCP_DNS"`
	Hidden          int    `ini:"HIDDEN"`
	MACFilter       int    `ini:"MAC_FILTER"`
	MACFilterAccept string `ini:"MAC_FILTER_ACCEPT"`
	IsolateClients  int    `ini:"ISOLATE_CLIENTS"`
	ShareMethod     string `ini:"SHARE_METHOD"`
	IEEE80211N      int
	IEEE80211AC     int
	HTCapAb         string  `ini:"HT_CAPAB"`
	VHTCapAb        string  `ini:"VHT_CAPAB"`
	Driver          string  `init:"DRIVER"`
	NoVirt          int     `ini:"NO_VIRT"`
	Country         string  `nini:"COUNTRY"`
	FreqBand        float64 `ini:"FREQ_BAND"`
	NewMACAddr      string  `ini:"NEW_MACADDR"`
	Daemonize       int     `ini:"DAEMONIZE"`
	NoHaveGED       int     `ini:"NO_HAVEGED"`
	WifiIface       string  `ini:"WIFI_IFACE"`
	InternetIface   string  `ini:"INTERNET_IFACE"`
	SSID            string  `ini:"SSID"`
	Passptrase      string  `ini:"PASSPHRASE"`
	UsePsk          int     `ini:"USE_PSK"`
}
