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

var ErrWrongStateFile = errors.New("fconf: wrong state file")

type WifiState struct {
	Enabled bool  `json:"enabled"`
	Configg *Wifi `json:"config"`
}

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
	var i string
	if ctx.IsSet("interface") {
		i = ctx.String("interface")
	} else {
		i = ctx.Args().First()
	}
	if i == "" {
		return errors.New("missing interface, you must specify interface")
	}
	w, err := wifiClientState(i)
	if err != nil {
		return err
	}
	if w.Configg.Interface == "" {
		w.Configg.Interface = "wlan0"
	}
	service := "wpa_supplicant@" + w.Configg.Interface
	err = restartService(service)
	if err != nil {
		return err
	}
	err = enableService(service)
	if err != nil {
		return err
	}
	err = restartService("systemd-networkd")
	if err != nil {
		return err
	}
	//w.Enabled = true
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}
	return keepState(
		fmt.Sprintf(defaultWifiClientConfig, w.Configg.Interface), data)

}

func wifiClientState(i string) (*WifiState, error) {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	b, err := ioutil.ReadFile(filepath.Join(dir,
		fmt.Sprintf(defaultWifiClientConfig, i)))
	if err != nil {
		return nil, err
	}
	w := &WifiState{}
	err = json.Unmarshal(b, w)
	if err != nil {
		return nil, err
	}
	if w.Configg == nil {
		return nil, ErrWrongStateFile
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
	e := Wifi{}
	err = json.Unmarshal(b, &e)
	if err != nil {
		return err
	}
	if e.Interface == "" {
		e.Interface = "wlan0"
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
	fmt.Printf("successful written wifi configuration to %s \n", filename)
	path := "/etc/wpa_supplicant/"
	err = checkDir(path)
	if err != nil {
		return err
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
	state := &WifiState{Configg: &e}
	b, _ = json.Marshal(state)
	fmt.Printf("successful written wifi connection  configuration to %s \n", filepath.Join(path, cname))
	ctx.Set("interface", e.Interface)
	return keepState(
		fmt.Sprintf(defaultWifiClientConfig, e.Interface), b)
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
	var i string
	if ctx.IsSet("interface") {
		i = ctx.String("interface")
	} else {
		i = ctx.Args().First()
	}
	if i == "" {
		return errors.New("missing interface, you must specify interface")
	}
	w, err := wifiClientState(i)
	if err != nil {
		return err
	}
	service := "wpa_supplicant@" + w.Configg.Interface
	err = disableService(service)
	if err != nil {
		return err
	}
	err = stopService(service)
	if err != nil {
		return err
	}
	err = restartService("systemd-networkd")
	if err != nil {
		return err
	}
	//w.Enabled = false
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}
	return keepState(
		fmt.Sprintf(defaultWifiClientConfig, i), data)
}

func RemoveWifi(ctx *cli.Context) error {
	var i string
	if ctx.IsSet("interface") {
		i = ctx.String("interface")
	} else {
		i = ctx.Args().First()
	}
	if i == "" {
		return errors.New("missing interface, you must specify interface")
	}
	w, err := wifiClientState(i)
	if err != nil {
		return err
	}

	// remove systemd file
	unit := filepath.Join(networkBase, wirelessService)
	err = removeFile(unit)
	if err != nil {
		return err
	}
	path := "/etc/wpa_supplicant/"
	cname := "wpa_supplicant-" + w.Configg.Interface + ".conf"

	// remove client connection
	err = removeFile(filepath.Join(path, cname))
	if err != nil {
		return err
	}
	err = DisableWifi(ctx)
	if err != nil {
		return err
	}

	// Remove any interface settings
	err = FlushInterface(w.Configg.Interface)
	if err != nil {
		return err
	}

	// remove the state file
	stateFile := filepath.Join(stateDir(),
		fmt.Sprintf(defaultWifiClientConfig, i))
	return removeFile(stateFile)
}
