package main

import (
	"io"

	ini "gopkg.in/ini.v1"
)

//AccessPoint  is the access point configuration
type AccessPoint struct {
	Channel         string `ini:"CHANNEL" json:"channel"`
	Gateway         string `ini:"GATEWAY" json:"gateway"`
	WPAVersion      int    `ini:"WPA_VERSION" json:"wpa_version"`
	ETCHosts        int    `ini:"ETC_HOSTS" json:"etc_hosts"`
	DHCPDNS         string `ini:"DHCP_DNS" json:"dhcp_dns"`
	NoDNS           int    `ini:"NO_DNS" json:"no_dns"`
	Hidden          int    `ini:"HIDDEN" json:"hidden"`
	MACFilter       int    `ini:"MAC_FILTER" json:"mac_filter"`
	MACFilterAccept string `ini:"MAC_FILTER_ACCEPT" json:"mac-filter_accept"`
	IsolateClients  int    `ini:"ISOLATE_CLIENTS" json:"isolate_clients"`
	ShareMethod     string `ini:"SHARE_METHOD" json:"share_method"`
	IEEE80211N      int
	IEEE80211AC     int
	HTCapAb         string  `ini:"HT_CAPAB" json:"ht_capab"`
	VHTCapAb        string  `ini:"VHT_CAPAB" json:"vht_capab"`
	Driver          string  `ini:"DRIVER" json:"driver"`
	NoVirt          int     `ini:"NO_VIRT" json:"no_virt"`
	Country         string  `ini:"COUNTRY" json:"country"`
	FreqBand        float64 `ini:"FREQ_BAND" json:"freq_band"`
	NewMACAddr      string  `ini:"NEW_MACADDR" json:"new_macaddr"`
	Daemonize       int     `ini:"DAEMONIZE" json:"daemonize"`
	NoHaveGED       int     `ini:"NO_HAVEGED" json:"no_haveged"`
	WifiIface       string  `ini:"WIFI_IFACE" json:"wifi_interface"`
	InternetIface   string  `ini:"INTERNET_IFACE" json:"internet_interface"`
	SSID            string  `ini:"SSID" json:"ssid"`
	Passptrase      string  `ini:"PASSPHRASE" json:"passphrase"`
	UsePsk          int     `ini:"USE_PSK" json:"use_psk"`
}

//LoadAPFromSrc loads access point configuration fom [byte
func LoadAPFromSrc(src []byte) (*AccessPoint, error) {
	cfg, err := ini.Load(src)
	if err != nil {
		return nil, err
	}
	a := &AccessPoint{}
	err = cfg.MapTo(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

//WriteTo writes ini representation of *AccessPoint to dst.
func (a *AccessPoint) WriteTo(dst io.Writer) (int64, error) {
	f := ini.Empty()
	err := f.ReflectFrom(a)
	if err != nil {
		return 0, err
	}
	return f.WriteTo(dst)
}
