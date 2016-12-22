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
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	e := &AccessPoint{}
	err = json.Unmarshal(b, e)
	if err != nil {
		return err
	}
	err = checkDir(base)
	if err != nil {
		return err
	}
	filename := filepath.Join(base, name)
	var buf bytes.Buffer
	_, err = e.WriteTo(&buf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("successful written access point configuration to %s \n", filename)
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return keepState(defaultAccessPointConfig, data)
}

func accessPointState() (*AccessPoint, error) {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	b, err := ioutil.ReadFile(filepath.Join(dir, defaultAccessPointConfig))
	if err != nil {
		return nil, err
	}
	a := &AccessPoint{}
	err = json.Unmarshal(b, a)
	if err != nil {
		return nil, err
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
	_, err := accessPointState()
	if err != nil {
		return err
	}
	service := "create_ap"
	err = startService(service)
	if err != nil {
		return err
	}
	return enableService(service)
}

func DisableApCMD(ctx *cli.Context) error {
	_, err := accessPointState()
	if err != nil {
		return err
	}
	service := "create_ap"
	err = stopService(service)
	if err != nil {
		return err
	}
	return disableService(service)
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
