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
   0.4.8
   
COMMANDS:
     ethernet, e        configures ethernet with systemd
     4g-ndis, 4g        configures 4G with systemd
     3g-ras, 3g         configures 3G 
     wifi-client, w     configures wifi client with systemd
     access-point, a    configures access point with systemd
     voice-channel, v   configures voice channel for 3g dongle
     list-interface, i  prints a json array of all interfaces
     help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --interface value  the interface
   --pid value        process id to send SIGHUP to (default: 0)
   --help, -h         show help
   --version, -v      print the version
   
```


# Network Configuration

## Ethernet

```
NAME:
   fconf ethernet - configures ethernet with systemd

USAGE:
   fconf ethernet [command options] [arguments...]

OPTIONS:
   --name value    The name of the unit file (default: "fconf-wired.network")
   --dir value     The directory in which to write the file (default: "/etc/systemd/network")
   --config value  The path to the json configuration file (default: "wired.json")
   --enable        Enables ethernet
   --disable       Disable ethernet
   --remove        Remove ethernet
```


You need to supply  the json file with the Ethernet configuration as the
first argument.


Example

	fconf e --config=fixture/wired_static.json

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
   --name value    The name of the unit file (default: "fconf-wireless.network")
   --dir value     The directory in which to write the file (default: "/etc/systemd/network")
   --config value  The path to the json configuration file (default: "wireless.json")
   --enable        Enables wifi
   --disable       Disable wifi
   --remove        Remove wifi
```

This shares the same configuration as for Ethernet, except you can add
username(ssid) and password for the wifi network.

Example

	fconf e --config=fixture/wireless.json

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
   "ssid":"HackME",
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
	"interface": "wlan0",
	"hidden": false,
	"channel": 0,
	"ssid": "voxbox",
	"passphrase": "voxbox99",
	"gateway": "192.168.12.1",
	"shared_interface": "eth0"
}
```
