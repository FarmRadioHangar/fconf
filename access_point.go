package main

import (
	"io"

	ini "gopkg.in/ini.v1"
)

//AccessPoint  is the access point configuration
type AccessPoint struct {
	Channel         string `ini:"CHANNEL"`
	Gateway         string `ini:"GATEWAY"`
	WPAVersion      int    `ini:"WPA_VERSION"`
	ETCHosts        int    `ini:"ETC_HOSTS"`
	DHCPDNS         string `ini:"DHCP_DNS"`
	NoDNS           int    `ini:"NO_DNS"`
	Hidden          int    `ini:"HIDDEN"`
	MACFilter       int    `ini:"MAC_FILTER"`
	MACFilterAccept string `ini:"MAC_FILTER_ACCEPT"`
	IsolateClients  int    `ini:"ISOLATE_CLIENTS"`
	ShareMethod     string `ini:"SHARE_METHOD"`
	IEEE80211N      int
	IEEE80211AC     int
	HTCapAb         string  `ini:"HT_CAPAB"`
	VHTCapAb        string  `ini:"VHT_CAPAB"`
	Driver          string  `ini:"DRIVER"`
	NoVirt          int     `ini:"NO_VIRT"`
	Country         string  `ini:"COUNTRY"`
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
