package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	defaultEthernetConfig = "wired.json"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1.4"
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
			},
			Action: EthernetCMD,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("fconf: %v", err)
	}
}
