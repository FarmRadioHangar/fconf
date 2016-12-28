package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

type AccessPointState struct {
	Enabled bool               `json:"enabled"`
	Configg *AccessPointConfig `json:"config"`
}

func ApCMD(ctx *cli.Context) error {
	if ctx.IsSet(enableFlag) {
		return EnableApCMD(ctx)
	}
	if ctx.IsSet(disableFlag) {
		return DisableApCMD(ctx)
	}
	if ctx.IsSet(removeFlag) {
		return RemoveApCMD(ctx)
	}
	if ctx.IsSet(configFlag) {
		return ConfigApCMD(ctx)
	}
	return nil
}

func ConfigApCMD(ctx *cli.Context) error {
	base := ctx.String("dir")
	name := ctx.String("name")
	src := ctx.String("config")
	if src == "" {
		return errors.New("fconf: missing argument")
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
	e := &AccessPointConfig{}
	err = json.Unmarshal(b, e)
	if err != nil {
		return err
	}
	err = checkDir(base)
	if err != nil {
		return err
	}
	ap := DefaultAccesPoint()
	ap.Update(e)
	filename := filepath.Join(base, name)
	var buf bytes.Buffer
	_, err = ap.WriteTo(&buf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("successful written access point configuration to %s \n", filename)
	state := &AccessPointState{Configg: ap.State()}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return keepState(defaultAccessPointConfig, data)
}

func accessPointState() (*AccessPointState, error) {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	b, err := ioutil.ReadFile(filepath.Join(dir, defaultAccessPointConfig))
	if err != nil {
		return nil, err
	}
	a := &AccessPointState{}
	err = json.Unmarshal(b, a)
	if err != nil {
		return nil, err
	}
	if a.Configg == nil {
		return nil, ErrWrongStateFile
	}
	return a, nil
}

func EnableApCMD(ctx *cli.Context) error {
	if ctx.IsSet(configFlag) {
		err := ConfigApCMD(ctx)
		if err != nil {
			return err
		}
	}
	state, err := accessPointState()
	if err != nil {
		return err
	}
	service := "create_ap"
	err = startService(service)
	if err != nil {
		return err
	}
	err = enableService(service)
	if err != nil {
		return err
	}
	state.Enabled = true
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return keepState(defaultAccessPointConfig, data)
}

func DisableApCMD(ctx *cli.Context) error {
	state, err := accessPointState()
	if err != nil {
		return err
	}
	service := "create_ap"
	err = stopService(service)
	if err != nil {
		return err
	}
	err = disableService(service)
	if err != nil {
		return err
	}
	state.Enabled = false
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return keepState(defaultAccessPointConfig, data)
}
func RemoveApCMD(ctx *cli.Context) error {
	_, err := accessPointState()
	if err != nil {
		return err
	}
	service := "create_ap"
	err = stopService(service)
	if err != nil {
		return err
	}
	err = disableService(service)
	if err != nil {
		return err
	}

	// remove the state file
	stateFile := filepath.Join(stateDir(), defaultAccessPointConfig)
	return removeFile(stateFile)
}
