package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	defaultEthernetConfig    = "ethernet@%s.json"
	defaultWifiClientConfig  = "wifi-client@%s.json"
	defaultAccessPointConfig = "access-point@%s.json"
	defaultFougGConfig       = "4g-ndis@%s.json"
	defaultThreeGGConfig     = "3g-ras@%s.json"
	defaultVoiceChanConfig   = "voice-channel@%s.json"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.4.7"
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
			Name:    "4g-ndis",
			Aliases: []string{"4g"},
			Usage:   "configures 4G with systemd",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "The name of the unit file",
					Value: fourgService,
				},
				cli.StringFlag{
					Name:  "dir",
					Usage: "The directory in which to write the file",
					Value: networkBase,
				},
				cli.StringFlag{
					Name:  "config",
					Usage: "The path to the json configuration file",
					Value: defaultFougGConfig,
				},
				cli.BoolFlag{
					Name:  "enable",
					Usage: "Enables 4G",
				},
				cli.BoolFlag{
					Name:  "disable",
					Usage: "Disable 4G",
				},
				cli.BoolFlag{
					Name:  "remove",
					Usage: "Remove 4G",
				},
			},
			Action: FourgCMD,
		},
		{
			Name:    "3g-ras",
			Aliases: []string{"3g"},
			Usage:   "configures 3G ",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "The name of the config file",
					Value: threeGService,
				},
				cli.StringFlag{
					Name:  "dir",
					Usage: "The directory in which to write the file",
					Value: apConfigBase,
				},
				cli.StringFlag{
					Name:  "config",
					Usage: "The path to the json configuration file",
				},
				cli.BoolFlag{
					Name:  "enable",
					Usage: "Enables 3G",
				},
				cli.BoolFlag{
					Name:  "disable",
					Usage: "Disable 3G",
				},
				cli.BoolFlag{
					Name:  "remove",
					Usage: "Remove 3G",
				},
			},
			Action: ThreegCMD,
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
		{
			Name:    "access-point",
			Aliases: []string{"a"},
			Usage:   "configures access point with systemd",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "The name of the unit file",
					Value: apConfigFile,
				},
				cli.StringFlag{
					Name:  "dir",
					Usage: "The directory in which to write the file",
					Value: apConfigBase,
				},
				cli.StringFlag{
					Name:  "config",
					Usage: "The path to the json configuration file",
					Value: defaultAccessPointConfig,
				},
				cli.BoolFlag{
					Name:  "enable",
					Usage: "Enables access point",
				},
				cli.BoolFlag{
					Name:  "disable",
					Usage: "Disable access point",
				},
				cli.BoolFlag{
					Name:  "remove",
					Usage: "Remove access point",
				},
			},
			Action: ApCMD,
		},
		{
			Name:    "voice-channel",
			Aliases: []string{"v"},
			Usage:   "configures voice channel for 3g dongle",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "The path to the json configuration file",
				},
				cli.BoolFlag{
					Name:  "enable",
					Usage: "Enables access point",
				},
				cli.BoolFlag{
					Name:  "disable",
					Usage: "Disable access point",
				},
				cli.BoolFlag{
					Name:  "remove",
					Usage: "Remove access point",
				},
			},
			Action: VoiceChannelCMD,
		},
		{
			Name:    "list-interface",
			Aliases: []string{"i"},
			Usage:   "prints a json array of all interfaces",
			Action:  ListInterface,
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "interface",
			Usage: "the interface",
		},
		cli.IntFlag{
			Name:  "pid",
			Usage: "process id to send SIGHUP to",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("fconf: %v", err)
	}
}
