package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

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
	NoVirt          int     `ini:"NO_VIRT" json:"no_virt,omitempty"`
	Country         string  `ini:"COUNTRY" json:"country"`
	FreqBand        float64 `ini:"FREQ_BAND" json:"freq_band"`
	NewMACAddr      string  `ini:"NEW_MACADDR" json:"new_macaddr"`
	Daemonize       int     `ini:"DAEMONIZE" json:"daemonize"`
	NoHaveGED       int     `ini:"NO_HAVEGED" json:"no_haveged"`
	WifiIface       string  `ini:"WIFI_IFACE" json:"wifi_interface"`
	InternetIface   string  `ini:"INTERNET_IFACE" json:"internet_interface"`
	SSID            string  `ini:"SSID" json:"ssid"`
	Passphrase      string  `ini:"PASSPHRASE" json:"passphrase"`
	UsePsk          int     `ini:"USE_PSK" json:"use_psk"`
}

//LoadAPFromSrc loads access point configuration fom [byte
func LoadAPFromConf(src []byte) (*AccessPoint, error) {
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

func LoadFromJSON(src []byte) (*AccessPoint, error) {
	a := DefaultAccesPoint()
	ap := &AccessPointConfig{}
	err := json.Unmarshal(src, ap)
	if err != nil {
		return nil, err
	}
	a.Update(ap)
	return a, nil
}

//WriteTo writes ini representation of *AccessPoint to dst.
func (a *AccessPoint) WriteTo(dst io.Writer) (int64, error) {
	f := ini.Empty()
	err := f.ReflectFrom(a)
	if err != nil {
		return 0, err
	}
	s := f.Sections()
	for _, sec := range s {
		names := sec.KeyStrings()
		for _, n := range names {
			v := sec.Key(n).String()
			fmt.Fprintf(dst, "%s=%s\n", n, v)
		}
	}
	return 0, nil
}

type AccessPointConfig struct {
	Interface      string `json:"interface"`
	Hidden         bool   `json:"hidden"`
	Channel        int    `json:"channel"`
	SSID           string `json:"ssid"`
	Passphrase     string `json:"passphrase"`
	Gateway        string `json:"gateway"`
	ShareInterfaec string `json:"shared_interface"`
	Enabled        bool   `json:"enabled"`
}

func (a *AccessPoint) Update(ap *AccessPointConfig) {
	if ap.ShareInterfaec != "" {
		a.ShareMethod = "nat"
		a.InternetIface = ap.ShareInterfaec
	} else {
		a.ShareMethod = "none"
		a.InternetIface = ""
	}
	if ap.SSID != "" {
		a.SSID = ap.SSID
		a.Passphrase = ap.Passphrase
	}
	if ap.Gateway != "" {
		a.Gateway = ap.Gateway
	}
	if ap.Channel > 0 {
		a.Channel = fmt.Sprint(ap.Channel)
	}
	if ap.Hidden {
		a.Hidden = 1
	} else {
		a.Hidden = 0
	}
	if ap.Interface != "" {
		a.WifiIface = ap.Interface
	}
}

func DefaultAccesPoint() *AccessPoint {
	var txt = `
{
	"channel": "default",
	"gateway": "192.168.12.1",
	"wpa_version": 2,
	"etc_hosts": 0,
	"dhcp_dns": "gateway",
	"no_dns": 0,
	"hidden": 0,
	"mac_filter": 0,
	"mac-filter_accept": "/etc/hostapd/hostapd.accept",
	"isolate_clients": 1,
	"share_method": "nat",
	"IEEE80211N": 0,
	"IEEE80211AC": 0,
	"ht_capab": "[HT40+]",
	"vht_capab": "",
	"driver": "nl80211",
	"no_virt": 1,
	"country": "",
	"freq_band": 2.4,
	"new_macaddr": "",
	"daemonize": 0,
	"no_haveged": 0,
	"wifi_interface": "wlan0",
	"internet_interface": "eth0",
	"ssid": "voxbox",
	"passphrase": "voxbox99",
	"use_psk": 0
}
	`
	a := &AccessPoint{}
	_ = json.Unmarshal([]byte(txt), a)
	return a
}

func (a *AccessPoint) State() *AccessPointConfig {
	ap := &AccessPointConfig{
		SSID:           a.SSID,
		Passphrase:     a.Passphrase,
		Gateway:        a.Gateway,
		Interface:      a.WifiIface,
		ShareInterfaec: a.InternetIface,
	}
	if a.Hidden == 1 {
		ap.Hidden = true
	}
	if a.Channel != "" && a.Channel != "default" {
		i, err := strconv.Atoi(a.Channel)
		if err != nil {
			log.Fatal(err)
		}
		ap.Channel = i
	}
	return ap
}
