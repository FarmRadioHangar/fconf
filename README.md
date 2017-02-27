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
   0.4.9
   
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