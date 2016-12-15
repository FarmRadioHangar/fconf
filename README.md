# fconf
voxbox configuration

This provide command line application that helps configure voxbox through json
files.

# Installation

You can install precompiled binaries

Or if you have Go installed.
	go get github.com/FarmRadioHangar/fconf

## Usage
```

NAME:
   fconf - fessbox configuration manager

USAGE:
   fconf [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     ethernet, e       configures ethernet with systemd
     wifi-client, w    configures wifi client with systemd
     access-point, ap  configures access point with systemd
     help, h           Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```


# Network Configuration

## Ethernet
```
NAME:
   fconf ethernet - configures ethernet with systemd

USAGE:
   fconf ethernet [command options] [arguments...]

OPTIONS:
   --name value  The name of the unit file (default: "wired.service")
   --dir value   The directory in which to write the file (default: "/etc/systemd/network")
   --restart     restarts the network service
```

You need to supply  the json file with the Ethernet configuration as the
first argument.

Example

	fconf e fixture/wired_static.json

This is a sample content of the json configuration file for Ethernet

```json
{
   "static":{
      "ip":"192.168.1.8/24",
      "gateway":"192.168.1.1"
   },
   "dhcp":true,
   "dns-servers":[
      "8.8.8.8",
      "8.8.8.8"
   ],
   "interface":"eth0"
}
```
That sample configures ethernet for both dhcp and  static.


## Wireless

### Wifi client

```
NAME:
   fconf wifi-client - configures wifi client with systemd

USAGE:
   fconf wifi-client [command options] [arguments...]

OPTIONS:
   --name value  The name of the unit file (default: "wireless.service")
   --dir value   The directory in which to write the file (default: "/etc/systemd/network")
   --restart     restarts the network service
   --connect     generates and starts service for wifi connection
```

This shares the same configuration as for Ethernet, except you can add
username(ssid) and password for the wifi network.

Example

	fconf e fixture/wireless.json

This is a sample content of the json configuration file for wifi client

```json
{
   "static":{
      "ip":"192.168.1.8/24",
      "gateway":"192.168.1.1"
   },
   "dhcp":true,
   "dns-servers":[
      "8.8.8.8",
      "8.8.8.8"
   ],
   "interface":"wlan0",
   "username":"HackME",
   "password":"mypassworld"
}
```

Note that the interface is changed to `wlan0` . You can omit the interface and
the default wifi network of `wlan0` will be used.


## Access Point

```
NAME:
   fconf access-point - configures access point with systemd

USAGE:
   fconf access-point [command options] [arguments...]

OPTIONS:
   --name value  The name of the configuration file (default: "create_ap.conf")
   --dir value   The directory in which to write the file (default: "/etc/")
   --restart     restarts the access point service
```

Example

	fconf ap  fixture/create_ap.json


This is a sample content of the json configuration file for access point

```json
{
	"channel": "default",
	"gateway": "10.0.0.1",
	"wpa_version": 2,
	"etc_hosts": 0,
	"dhcp_dns": "gateway",
	"no_dns": 0,
	"hidden": 0,
	"mac_filter": 0,
	"mac-filter_accept": "/etc/hostapd/hostapd.accept",
	"isolate_clients": 0,
	"share_method": "nat",
	"IEEE80211N": 0,
	"IEEE80211AC": 0,
	"ht_capab": "[HT40+]",
	"vht_capab": "",
	"driver": "nl80211",
	"no_virt": 0,
	"country": "",
	"freq_band": 2.4,
	"new_macaddr": "",
	"daemonize": 0,
	"no_haveged": 0,
	"wifi_interface": "wlan0",
	"internet_interface": "eth0",
	"ssid": "MyAccessPoint",
	"passphrase": "12345678",
	"use_psk": 0
}
```
