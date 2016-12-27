package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/coreos/go-systemd/unit"
	"github.com/urfave/cli"
)

type FourGState struct {
	Enabled bool   `json:"enabled"`
	Configg *FourG `json:"config"`
}

type FourG struct {
	Network
}

//ToSystemdUnit implement UnitFile interface
func (f FourG) ToSystemdUnit() ([]*unit.UnitOption, error) {
	if f.Interface == "" {
		f.Interface = "eth1"
	}
	return f.Network.ToSystemdUnit()
}

func fourGState() (*FourGState, error) {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	b, err := ioutil.ReadFile(filepath.Join(dir, defaultFougGConfig))
	if err != nil {
		return nil, err
	}
	f := &FourGState{}
	err = json.Unmarshal(b, f)
	if err != nil {
		return nil, err
	}
	if f.Configg.Interface == "" {
		f.Configg.Interface = "eth1"
	}
	return f, nil
}

func FourgCMD(ctx *cli.Context) error {
	if ctx.IsSet(enableFlag) {
		return EnableFourg(ctx)
	}
	if ctx.IsSet(disableFlag) {
		return DisableFourg(ctx)
	}
	if ctx.IsSet(removeFlag) {
		return RemoveFourg(ctx)
	}
	if ctx.IsSet(configFlag) {
		return configFourgCMD(ctx)
	}
	return nil
}
func RemoveFourg(ctx *cli.Context) error {
	err := DisableFourg(ctx)
	if err != nil {
		return err
	}

	// removestate file
	stateFile := filepath.Join(stateDir(), defaultFougGConfig)
	err = removeFile(stateFile)
	if err != nil {
		return err
	}
	// remove systemd file
	unit := filepath.Join(networkBase, fourgService)
	err = removeFile(unit)
	if err != nil {
		return err
	}

	// reload systemd-networkd
	return restartService("systemd-networkd")
}
func EnableFourg(ctx *cli.Context) error {
	if ctx.IsSet(configFlag) {
		err := configFourgCMD(ctx)
		if err != nil {
			return err
		}
	}
	e, err := fourGState()
	if err != nil {
		return err
	}
	_, err = exec.Command("ip", "link", "set", "up", e.Configg.Interface).Output()
	if err != nil {
		return err
	}
	err = restartService("systemd-networkd")
	if err != nil {
		return err
	}
	//e.Enabled = true
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return keepState(defaultFougGConfig, data)
}

func DisableFourg(ctx *cli.Context) error {
	e, err := fourGState()
	if err != nil {
		return err
	}
	_, err = exec.Command("ip", "link", "set", "down", "dev", e.Configg.Interface).Output()
	if err != nil {
		return err
	}
	fmt.Println("successfully disabled 4G")
	//e.Enabled = false
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return keepState(defaultFougGConfig, data)
}

func configFourgCMD(ctx *cli.Context) error {
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
	e := FourG{}
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
	fmt.Printf("successful written 4G configuration to %s \n", filename)
	state := &FourGState{Configg: &e}
	b, _ = json.Marshal(state)
	return keepState(defaultFougGConfig, b)
}
