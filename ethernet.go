package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"
)

func EthernetCMD(ctx *cli.Context) error {
	if ctx.IsSet(enableFlag) {
		return EnableEthernet(ctx)
	}
	if ctx.IsSet(configFlag) {
		return configEthernetCMD(ctx)
	}
	return nil
}

// Enable ethernet enables ethernet network in the host machine. This relies on
// systemd to be the init system of the host machine.
//
// If the config flag is set, ethernet will be configured before  enabling it.
// Ommit the config flag if  ethernet is already configured using this tool.
func EnableEthernet(ctx *cli.Context) error {
	if ctx.IsSet(configFlag) {
		err := configEthernetCMD(ctx)
		if err != nil {
			return err
		}
	}
	e, err := ethernetState()
	if err != nil {
		return err
	}
	_, err = exec.Command("ip", "link", "set", "up", e.Interface).Output()
	if err != nil {
		return err
	}
	return restartService("systemd-networkd")
}

// gives the current state of the ethernet configuration. This will return an
// error if the system hast been configured yet.
//
// Configuration state files are written in $FCONF_CONFIGDIR directory.
func ethernetState() (*Ethernet, error) {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	b, err := ioutil.ReadFile(filepath.Join(dir, defaultEthernetConfig))
	if err != nil {
		return nil, err
	}
	e := &Ethernet{}
	err = json.Unmarshal(b, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func configEthernetCMD(ctx *cli.Context) error {
	base := ctx.String("dir")
	name := ctx.String("name")
	src := ctx.String("config")
	if src == "" {
		return errors.New("fconf: missing configuration source file")
	}
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	e := Ethernet{}
	err = json.Unmarshal(b, &e)
	if err != nil {
		return err
	}
	err = checkDir(base)
	if err != nil {
		return err
	}
	filename := filepath.Join(base, name)
	err = CreateSystemdFile(e, filename, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("successful written ethernet configuration to %s \n", filename)
	return keepState(defaultEthernetConfig, b)
}

func keepState(filename string, src []byte) error {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	return ioutil.WriteFile(filepath.Join(dir, filename), src, 0644)
}
