package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

type VoiceChannel struct {
	IMEI   string `json:"imei"`
	IMSI   string `json:"imsi"`
	Name   string `json:"name"`
	Number string `json:"number"`
	RxGain int    `json:"rx-gain"`
	TxGain int    `json:"tx-gain"`
	Label  string `json:"label"`
}

type VoiceState struct {
	Enabled bool          `json:"enabled"`
	Config  *VoiceChannel `json:"config"`
}

func VoiceChannelCMD(ctx *cli.Context) error {
	if ctx.IsSet(enableFlag) {
		return EnableVoiceChan(ctx)
	}
	if ctx.IsSet(disableFlag) {
		return DisableVoiceChan(ctx)
	}
	if ctx.IsSet(removeFlag) {
		return RemoveVoiceChan(ctx)
	}
	if ctx.IsSet(configFlag) {
		return configVoice(ctx)
	}
	return nil
}

func configVoice(ctx *cli.Context) error {
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
	e := VoiceChannel{}
	err = json.Unmarshal(b, &e)
	if err != nil {
		return err
	}
	state := &VoiceState{Config: &e}
	ws, err := voiceChanState(e.IMEI)
	if err == nil {
		state.Enabled = ws.Enabled
	}
	b, _ = json.Marshal(state)
	setInterface(ctx, e.IMEI)
	return keepState(
		fmt.Sprintf(defaultVoiceChanConfig, e.IMEI), b)
}

func voiceChanState(i string) (*VoiceState, error) {
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	b, err := ioutil.ReadFile(filepath.Join(dir,
		fmt.Sprintf(defaultVoiceChanConfig, i)))
	if err != nil {
		return nil, err
	}
	w := &VoiceState{}
	err = json.Unmarshal(b, w)
	if err != nil {
		return nil, err
	}
	if w.Config == nil {
		return nil, ErrWrongStateFile
	}
	return w, nil
}

func EnableVoiceChan(ctx *cli.Context) error {
	if ctx.IsSet(configFlag) {
		err := configVoice(ctx)
		if err != nil {
			return err
		}
	}
	i := getInterface(ctx)
	if i == "" {
		return errors.New("missing imei , you must specify imei")
	}
	w, err := voiceChanState(i)
	if err != nil {
		return err
	}
	w.Enabled = true
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}
	return keepState(
		fmt.Sprintf(defaultVoiceChanConfig, w.Config.IMEI), data)

}

func DisableVoiceChan(ctx *cli.Context) error {
	i := getInterface(ctx)
	if i == "" {
		return errors.New("missing imei , you must specify imei")
	}
	w, err := voiceChanState(i)
	if err != nil {
		return err
	}
	w.Enabled = false
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}
	return keepState(
		fmt.Sprintf(defaultVoiceChanConfig, w.Config.IMEI), data)

}

func RemoveVoiceChan(ctx *cli.Context) error {
	i := getInterface(ctx)
	if i == "" {
		return errors.New("missing imei , you must specify imei")
	}
	w, err := voiceChanState(i)
	if err != nil {
		return err
	}
	dir := os.Getenv("FCONF_CONFIGDIR")
	if dir == "" {
		dir = fconfConfigDir
	}
	err = checkDir(dir)
	if err != nil {
		return err
	}
	name := filepath.Join(dir, fmt.Sprintf(defaultVoiceChanConfig, w.Config.IMEI))
	return removeFile(name)
}
