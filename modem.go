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

func fourGState(i string) (*FourGState, error) {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	b, err := ioutil.ReadFile(filepath.Join(dir,
		fmt.Sprintf(defaultFougGConfig, i)))
	if err != nil {
		return nil, err
	}
	f := &FourGState{}
	err = json.Unmarshal(b, f)
	if err != nil {
		return nil, err
	}
	if f.Configg == nil {
		return nil, ErrWrongStateFile
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

func getInterface(ctx *cli.Context) string {
	var i string
	if ctx.GlobalIsSet("interface") {
		i = ctx.GlobalString("interface")
	} else {
		i = ctx.Args().First()
	}
	return i
}
func setInterface(ctx *cli.Context, i string) {
	ctx.GlobalSet("interface", i)
}
func RemoveFourg(ctx *cli.Context) error {
	i := getInterface(ctx)
	if i == "" {
		return errors.New("missing interface, you must specify interface")
	}
	f, err := fourGState(i)
	if err != nil {
		return err
	}
	if f.Enabled {
		err = DisableFourg(ctx)
		if err != nil {
			return err
		}
	}
	// removestate file
	stateFile := filepath.Join(stateDir(),
		fmt.Sprintf(defaultFougGConfig, i))
	err = removeFile(stateFile)
	if err != nil {
		return err
	}
	// remove systemd file
	unit := filepath.Join(networkBase,
		fmt.Sprintf(fourgService, i))
	err = removeFile(unit)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	err = FlushInterface(f.Configg.Interface)
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
	i := getInterface(ctx)
	if i == "" {
		return errors.New("missing interface, you must specify interface")
	}
	e, err := fourGState(i)
	if err != nil {
		return err
	}
	unit := filepath.Join(networkBase,
		fmt.Sprintf(fourgService, e.Configg.Interface))
	_, err = os.Stat(unit)
	if os.IsNotExist(err) {
		err = CreateSystemdFile(e.Configg, unit, 0644)
		if err != nil {
			return err
		}
	}
	_, err = exec.Command("ip", "link", "set", "up", e.Configg.Interface).Output()
	if err != nil {
		return fmt.Errorf("ERROR: runnin ip link set up %s %v",
			e.Configg.Interface, err,
		)
	}
	err = restartService("systemd-networkd")
	if err != nil {
		return fmt.Errorf("ERROR: restarting systemd %v ", err)
	}
	e.Enabled = true
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return keepState(
		fmt.Sprintf(defaultFougGConfig, i), data)
}

func DisableFourg(ctx *cli.Context) error {
	if ctx.IsSet(configFlag) {
		fmt.Println("WARN: config flag will be ignored when diable flag is used")
	}
	i := getInterface(ctx)
	if i == "" {
		return errors.New("missing interface, you must specify interface")
	}
	e, err := fourGState(i)
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
		fmt.Sprintf(fourgService, e.Configg.Interface))
	err = removeFile(unit)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	fmt.Println("successfully disabled 4G")
	e.Enabled = false
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	err = restartService("systemd-networkd")
	if err != nil {
		return err
	}
	return keepState(
		fmt.Sprintf(defaultFougGConfig, i), data)
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
	if e.Interface == "" {
		e.Interface = "eth1"
	}
	err = checkDir(base)
	if err != nil {
		return err
	}
	if strings.Contains(name, "%s") {
		name = fmt.Sprintf(name, e.Interface)
	}
	filename := filepath.Join(base, name)
	err = CreateSystemdFile(e, filename, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("successful written 4G configuration to %s \n", filename)
	state := &FourGState{Configg: &e}
	ms, err := fourGState(e.Interface)
	if err == nil {
		state.Enabled = ms.Enabled
	}
	ctx.GlobalSet("interface", e.Interface)
	b, _ = json.Marshal(state)
	return keepState(
		fmt.Sprintf(defaultFougGConfig, e.Interface), b)
}
