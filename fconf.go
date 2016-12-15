package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/coreos/go-systemd/unit"
	"github.com/urfave/cli"
)

const (
	networkBase       = "/etc/systemd/network"
	ethernetService   = "wired.service"
	wirelessService   = "wireless.service"
	accessPointConfig = "/etc/create_ap.conf"
)

//Ethernet is the ehternet configuration.
type Ethernet struct {
	Network
}

//ToSystemdUnit implement UnitFile interface
func (e Ethernet) ToSystemdUnit() ([]*unit.UnitOption, error) {
	if e.Interface == "" {
		e.Interface = "eth0"
	}
	return e.Network.ToSystemdUnit()
}

//Wifi is the wifi configuration.
type Wifi struct {
	Network
	Username string `json:"username"`
	Password string `json:"password"`
}

//UnitFile is an interface for systemd uni file
type UnitFile interface {
	ToSystemdUnit() ([]*unit.UnitOption, error)
}

//ToSystemdUnit implement UnitFile interface
func (w Wifi) ToSystemdUnit() ([]*unit.UnitOption, error) {
	if w.Interface == "" {
		w.Interface = "wlan0"
	}
	return w.Network.ToSystemdUnit()
}

//CreateSystemdFile creates a file that has systemd unit file content.
func CreateSystemdFile(u UnitFile, filename string, mode os.FileMode, out ...io.Writer) error {
	x, err := u.ToSystemdUnit()
	if err != nil {
		return err
	}
	r := unit.Serialize(x)
	if len(out) > 0 {
		_, err := io.Copy(out[0], r)
		return err
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = io.Copy(f, r)
	return err
}

func ethernet(ctx *cli.Context) error {
	base := ctx.String("dir")
	name := ctx.String("name")
	src := ctx.Args().First()
	restart := ctx.BoolT("restart")
	if src == "" {
		return errors.New("fconf: missing argument")
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
	if restart {
		return restartService("systemd-networkd")
	}
	return nil
}

// Checks if the directory exists. If the directory doesnt exist, this function
// will create the directory with permission 0755.
//
// The directory created will recursively create subdirectory. It will behave
// something like mkdir -p /dir/subdir.
func checkDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 07755)
		if err != nil {
			return err
		}
	}
	return nil
}

func wifiClient(ctx *cli.Context) error {
	base := ctx.String("dir")
	name := ctx.String("name")
	src := ctx.Args().First()
	restart := ctx.BoolT("restart")
	connect := ctx.BoolT("connect")
	if src == "" {
		return errors.New("fconf: missing argument")
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
	if connect {
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
		err = enableService("wpa_supplicant@wlan0")
		if err != nil {
			return err
		}
		err = restartService("wpa_supplicant@wlan0")
		if err != nil {
			return err
		}
	}
	if restart {
		return restartService("systemd-networkd")
	}
	return nil
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

func restartService(name string) error {
	fmt.Print("restarting ", name, "...")
	_, err := exec.Command("systemctl", "restart", name).Output()
	fmt.Println("done")
	return err
}
func enableService(name string) error {
	fmt.Print("enabling ", name, "...")
	_, err := exec.Command("systemctl", "enable", name).Output()
	fmt.Println("done")
	return err
}

func accessPoint(ctx *cli.Context) error {
	base := ctx.String("dir")
	name := ctx.String("name")
	src := ctx.Args().First()
	restart := ctx.BoolT("restart")
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
	if restart {
		return restartService("create_ap")
	}
	return nil
}
