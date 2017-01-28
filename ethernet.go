package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

type EthernetState struct {
	Enabled bool      `json:"enabled"`
	Configg *Ethernet `json:"config"`
}

func EthernetCMD(ctx *cli.Context) error {
	if ctx.IsSet(enableFlag) {
		return EnableEthernet(ctx)
	}
	if ctx.IsSet(disableFlag) {
		return DisableEthernet(ctx)
	}
	if ctx.IsSet(removeFlag) {
		return RemoveEthernet(ctx)
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
	i := getInterface(ctx)
	e, err := ethernetState(i)
	if err != nil {
		return err
	}
	unit := filepath.Join(networkBase,
		fmt.Sprintf(ethernetService, e.Configg.Interface))
	_, err = os.Stat(unit)
	if os.IsNotExist(err) {
		err = CreateSystemdFile(e.Configg, unit, 0644)
		if err != nil {
			return err
		}
	}
	_, err = exec.Command("ip", "link", "set", "up", e.Configg.Interface).Output()
	if err != nil {
		return err
	}
	err = restartService("systemd-networkd")
	if err != nil {
		return err
	}
	e.Enabled = true
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return keepState(
		fmt.Sprintf(defaultEthernetConfig, i), data)
}

// gives the current state of the ethernet configuration. This will return an
// error if the system hast been configured yet.
//
// Configuration state files are written in $FCONF_CONFIGDIR directory.
func ethernetState(i string) (*EthernetState, error) {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	b, err := ioutil.ReadFile(filepath.Join(dir,
		fmt.Sprintf(defaultEthernetConfig, i)))
	if err != nil {
		return nil, err
	}
	e := &EthernetState{}
	err = json.Unmarshal(b, e)
	if err != nil {
		return nil, err
	}
	if e.Configg == nil {
		return nil, ErrWrongStateFile
	}
	if e.Configg.Interface == "" {
		e.Configg.Interface = "eth0"
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
	var b []byte
	var err error
	if src == "stdin" {
		b, err = ReadFromStdin()
		if err != nil {
			return err
		}
	} else {
		b, err = ioutil.ReadFile(src)
		if err != nil {
			return err
		}
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
	if e.Interface == "" {
		e.Interface = "eth0"
	}
	if strings.Contains(name, "%s") {
		name = fmt.Sprintf(name, e.Interface)
	}
	filename := filepath.Join(base, name)
	err = CreateSystemdFile(e, filename, 0644)
	if err != nil {
		return err
	}
	setInterface(ctx, e.Interface)
	fmt.Printf("successful written ethernet configuration to %s \n", filename)
	state := &EthernetState{Configg: &e}
	es, err := ethernetState(e.Interface)
	if err == nil {
		state.Enabled = es.Enabled
	}
	b, _ = json.Marshal(state)
	return keepState(
		fmt.Sprintf(defaultEthernetConfig, e.Interface), b)
}

func keepState(filename string, src []byte) error {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	err := checkDir(dir)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dir, filename), src, 0644)
}

//DisableEthernet disables ethernet temporaly.
func DisableEthernet(ctx *cli.Context) error {
	if ctx.IsSet(configFlag) {
		fmt.Println("WARN: config flag will be ignored when diable flag is used")
	}
	i := getInterface(ctx)
	if i == "" {
		return errors.New("missing interface, you must specify interface")
	}
	e, err := ethernetState(i)
	if err != nil {
		return err
	}
	_, err = exec.Command("ip", "addr", "flush", "dev", e.Configg.Interface).Output()
	if err != nil {
		return fmt.Errorf("ERROR: running ip addr flush dev %s %v",
			e.Configg.Interface, err,
		)
	}
	unit := filepath.Join(networkBase,
		fmt.Sprintf(ethernetService, e.Configg.Interface))
	err = removeFile(unit)
	if err != nil {
		return err
	}
	err = restartService("systemd-networkd")
	if err != nil {
		return err
	}
	fmt.Println("successfully disabled ethernet")
	e.Enabled = false
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return keepState(
		fmt.Sprintf(defaultEthernetConfig, i), data)
}

//RemoveEthernet removes ethernet service.
func RemoveEthernet(ctx *cli.Context) error {
	i := getInterface(ctx)
	if i == "" {
		return errors.New("missing interface, you must specify interface")
	}
	e, err := ethernetState(i)
	if err != nil {
		return err
	}
	if e.Enabled {
		err = DisableEthernet(ctx)
		if err != nil {
			return err
		}
	}
	// removestate file
	stateFile := filepath.Join(stateDir(),
		fmt.Sprintf(defaultEthernetConfig, i))
	err = removeFile(stateFile)
	if err != nil {
		return err
	}
	// remove systemd file
	unit := filepath.Join(networkBase, ethernetService)
	err = removeFile(unit)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	// Flush settings for the interface
	err = FlushInterface(e.Configg.Interface)
	if err != nil {
		return err
	}

	// reload systemd-networkd
	return restartService("systemd-networkd")
}

func removeFile(name string) error {
	fmt.Printf("removing %s ...", name)
	err := os.Remove(name)
	if err != nil {
		fmt.Println(" error")
		return err
	}
	fmt.Println(" done without error")
	return nil
}

func stateDir() string {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	return dir
}
