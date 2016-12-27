package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/coreos/go-systemd/unit"
)

const (
	networkBase       = "/etc/systemd/network"
	apConfigBase      = "/etc/"
	ethernetService   = "fconf-wired.network"
	fourgService      = "fconf-4g.network"
	wirelessService   = "fconf-wireless.network"
	apConfigFile      = "create_ap.conf"
	accessPointConfig = "create_ap.conf"
	enableFlag        = "enable"
	disableFlag       = "disable"
	removeFlag        = "remove"
	configFlag        = "config"
	fconfConfigDir    = "/etc/fconf"
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
	Username string `json:"ssid"`
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

func restartService(name string) error {
	return systemdCMD("restart", name)
}

func startService(name string) error {
	return systemdCMD("start", name)
}

func enableService(name string) error {
	return systemdCMD("enable", name)
}
func disableService(name string) error {
	return systemdCMD("disable", name)
}

func stopService(name string) error {
	return systemdCMD("stop", name)
}

func systemdCMD(name, service string) error {
	fmt.Printf("%s %s ...", name, service)
	_, err := exec.Command("systemctl", name, service).Output()
	if err != nil {
		fmt.Println("done with error")
		return err
	}
	fmt.Println("done without error")
	return nil
}

func ReadFromStdin() ([]byte, error) {
	r := bufio.NewReader(os.Stdin)
	return r.ReadBytes('\n')
}
