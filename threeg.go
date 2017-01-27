package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

type ThreeG struct {
	IMEI           string `json:"imei"`
	IMSI           string `json:"imsi"`
	APN            string `json:"apn"`
	Dial           string `json:"dial"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	DefaultGateway bool   `json:"defaultGateway"`
}

func (c *ThreeG) Prepare() *ThreeG {
	if c.Username == "" {
		c.Username = "0"
	}
	if c.Password == "" {
		c.Password = "0"
	}
	return c
}

func (c *ThreeG) WriteTo(out io.Writer) error {
	cfgTpl := `
[Dialer Defaults]
Init2 = ATQ0 V1 E1 S0=0 &C1 &D2
Modem Type = USB Modem
; Phone = <Target Phone Number>
Phone = {{.Dial}}
ISDN = 0
; Username = <Your Login Name>, 0 if not specified
Username = {{.Username}}
; Password = <Your Password>, 0 if not specified
Password = {{.Password}}
Init1 = ATZ
; modem command port (by IMSI)
Modem = /dev/{{.IMSI}}.imsi
Baud = 115200
; set APN
Init3 = AT+CGDCONT=1,"IP","{{.APN}}"
Stupid Mode = 1
	`
	tpl, err := template.New("cfg").Parse(cfgTpl)
	if err != nil {
		return err
	}
	return tpl.Execute(out, c)
}

type ThreeGState struct {
	Enabled bool    `json:"enabled"`
	Configg *ThreeG `json:"config"`
}

func ThreegCMD(ctx *cli.Context) error {
	if ctx.IsSet(enableFlag) {
		return EnableThreeg(ctx)
	}
	if ctx.IsSet(disableFlag) {
		return DisableThreeg(ctx)
	}
	//if ctx.IsSet(removeFlag) {
	//return RemoveThreeg(ctx)
	//}
	if ctx.IsSet(configFlag) {
		return configThreegCMD(ctx)
	}
	return nil
}

func configThreegCMD(ctx *cli.Context) error {
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
	e := ThreeG{}
	err = json.Unmarshal(b, &e)
	if err != nil {
		return err
	}
	e.Prepare()
	err = checkDir(base)
	if err != nil {
		return err
	}
	if strings.Contains(name, "%s") {
		name = fmt.Sprintf(name, e.IMEI)
	}
	filename := filepath.Join(base, name)
	var buf bytes.Buffer
	err = e.WriteTo(&buf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("successful written 3G configuration to %s \n", filename)
	state := &ThreeGState{Configg: &e}
	ms, err := threeGState(e.IMEI)
	if err == nil {
		state.Enabled = ms.Enabled
	}
	b, _ = json.Marshal(state)
	setInterface(ctx, e.IMEI)
	return keepState(
		fmt.Sprintf(defaultThreeGGConfig, e.IMEI), b)
}

func threeGState(i string) (*ThreeGState, error) {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	b, err := ioutil.ReadFile(filepath.Join(dir,
		fmt.Sprintf(defaultThreeGGConfig, i)))
	if err != nil {
		return nil, err
	}
	f := &ThreeGState{}
	err = json.Unmarshal(b, f)
	if err != nil {
		return nil, err
	}
	if f.Configg == nil {
		return nil, ErrWrongStateFile
	}
	return f, nil
}

func EnableThreeg(ctx *cli.Context) error {
	if ctx.IsSet(configFlag) {
		err := configThreegCMD(ctx)
		if err != nil {
			return err
		}
	}
	i := getInterface(ctx)
	if i == "" {
		return errors.New("missing imei , you must specify imei")
	}
	e, err := threeGState(i)
	if err != nil {
		return err
	}
	service := "wvdial"
	err = restartService(service)
	if err != nil {
		return fmt.Errorf("ERROR: restarting systemd %v ", err)
	}
	err = enableService(service)
	if err != nil {
		return fmt.Errorf("ERROR: enabling systemd %v ", err)
	}
	e.Enabled = true
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	fmt.Printf("successfully enabled 3g for %s \n", i)
	return keepState(
		fmt.Sprintf(defaultThreeGGConfig, i), data)
}

func DisableThreeg(ctx *cli.Context) error {
	if ctx.IsSet(configFlag) {
		fmt.Println("WARN: config flag will be ignored when diable flag is used")
	}
	i := getInterface(ctx)
	if i == "" {
		return errors.New("missing imei, you must specify imei")
	}
	e, err := threeGState(i)
	if err != nil {
		return err
	}
	e.Enabled = false
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	service := "wvdial"
	err = stopService(service)
	if err != nil {
		return err
	}
	err = disableService(service)
	if err != nil {
		return err
	}
	fmt.Println("successfully disabled 3G for ", i)
	return keepState(
		fmt.Sprintf(defaultThreeGGConfig, i), data)
}
