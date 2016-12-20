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

func WifiClientCMD(ctx *cli.Context) error {
	if ctx.IsSet(enableFlag) {
		return EnableWifiClient(ctx)
	}
	if ctx.IsSet(disableFlag) {
		return DisableWifi(ctx)
	}
	if ctx.IsSet(removeFlag) {
		return RemoveWifi(ctx)
	}
	if ctx.IsSet(configFlag) {
		return configWifiClient(ctx)
	}
	return nil
}

//EnableWifiClient enables wifi client. If the config flag is set, wifi is
//configured before being enabled.
func EnableWifiClient(ctx *cli.Context) error {
	if ctx.IsSet(configFlag) {
		err := configWifiClient(ctx)
		if err != nil {
			return err
		}
	}
	w, err := wifiClientState()
	if err != nil {
		return err
	}
	if w.Interface == "" {
		w.Interface = "wlan0"
	}
	service := "wpa_supplicant@" + w.Interface
	err = restartService(service)
	if err != nil {
		return err
	}
	err = enableService(service)
	if err != nil {
		return err
	}
	return restartService("systemd-networkd")

}

func wifiClientState() (*Wifi, error) {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	b, err := ioutil.ReadFile(filepath.Join(dir, defaultWifiClientConfig))
	if err != nil {
		return nil, err
	}
	w := &Wifi{}
	err = json.Unmarshal(b, w)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func configWifiClient(ctx *cli.Context) error {
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
	e := Wifi{}
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
	fmt.Printf("successful written wifi configuration to %s \n", filename)
	path := "/etc/wpa_supplicant/"
	err = checkDir(path)
	if err != nil {
		return err
	}
	if e.Interface == "" {
		e.Interface = "wlan0"
	}
	cname := "wpa_supplicant-" + e.Interface + ".conf"
	s, err := wifiConfig(e.Username, e.Password)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(path, cname), []byte(s), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("successful written wifi connection  configuration to %s \n", filepath.Join(path, cname))
	return keepState(defaultWifiClientConfig, b)
}

func wifiConfig(username, password string) (string, error) {
	cmd := "/usr/bin/wpa_passphrase"
	firstLine := "ctrl_interface=/run/wpa_supplicant_fconf"
	o, err := exec.Command(cmd, username, password).Output()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s \n \n%s\n", firstLine, string(o)), nil
}

func DisableWifi(ctx *cli.Context) error {
	w, err := wifiClientState()
	if err != nil {
		return err
	}
	if w.Interface == "" {
		w.Interface = "wlan0"
	}

	service := "wpa_supplicant@" + w.Interface
	err = disableService(service)
	if err != nil {
		return err
	}
	err = stopService(service)
	if err != nil {
		return err
	}
	return restartService("systemd-networkd")
}

func RemoveWifi(ctx *cli.Context) error {
	w, err := wifiClientState()
	if err != nil {
		return err
	}
	if w.Interface == "" {
		w.Interface = "wlan0"
	}

	// remove systemd file
	unit := filepath.Join(networkBase, wirelessService)
	err = removeFile(unit)
	if err != nil {
		return err
	}
	path := "/etc/wpa_supplicant/"
	cname := "wpa_supplicant-" + w.Interface + ".conf"

	// remove client connection
	err = removeFile(filepath.Join(path, cname))
	if err != nil {
		return err
	}
	err = DisableWifi(ctx)
	if err != nil {
		return err
	}

	// remove the state file
	stateFile := filepath.Join(stateDir(), defaultWifiClientConfig)
	return removeFile(stateFile)
}
