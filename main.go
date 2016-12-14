package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
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
				cli.BoolTFlag{
					Name:  "restart",
					Usage: "restarts the network service",
				},
			},
			Action: ethernet,
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
				cli.BoolTFlag{
					Name:  "restart",
					Usage: "restarts the network service",
				},
			},
			Action: wifiClient,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("fconf: %v", err)
	}
}
