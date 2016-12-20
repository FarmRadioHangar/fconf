package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	defaultEthernetConfig   = "wired.json"
	defaultWifiClientConfig = "wireless.json"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.2.0"
	app.Name = "fconf"
	app.Usage = "fessbox configuration manager"
	app.Commands = []cli.Command{
		{
			Name:    "ethernet",
			Aliases: []string{"e"},
			Usage:   "configures ethernet with systemd",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "The name of the unit file",
					Value: ethernetService,
				},
				cli.StringFlag{
					Name:  "dir",
					Usage: "The directory in which to write the file",
					Value: networkBase,
				},
				cli.StringFlag{
					Name:  "config",
					Usage: "The path to the json configuration file",
					Value: defaultEthernetConfig,
				},
				cli.BoolFlag{
					Name:  "enable",
					Usage: "Enables ethernet",
				},
				cli.BoolFlag{
					Name:  "disable",
					Usage: "Disable ethernet",
				},
				cli.BoolFlag{
					Name:  "remove",
					Usage: "Remove ethernet",
				},
			},
			Action: EthernetCMD,
		},
		{
			Name:    "wifi-client",
			Aliases: []string{"w"},
			Usage:   "configures wifi client with systemd",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "The name of the unit file",
					Value: wirelessService,
				},
				cli.StringFlag{
					Name:  "dir",
					Usage: "The directory in which to write the file",
					Value: networkBase,
				},
				cli.StringFlag{
					Name:  "config",
					Usage: "The path to the json configuration file",
					Value: defaultWifiClientConfig,
				},
				cli.BoolFlag{
					Name:  "enable",
					Usage: "Enables wifi",
				},
				cli.BoolFlag{
					Name:  "disable",
					Usage: "Disable wifi",
				},
				cli.BoolFlag{
					Name:  "remove",
					Usage: "Remove wifi",
				},
			},
			Action: WifiClientCMD,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("fconf: %v", err)
	}
}
