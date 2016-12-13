package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/coreos/go-systemd/unit"
)

func TestEthernet_ToSystemdUnit(t *testing.T) {
	e := &Ethernet{
		Static: &Static{
			IP:      "192.168.1.8/24",
			Gateway: "192.168.1.1",
		},
		DNS:       []string{"8.8.8.8", "8.8.8.8"},
		Interface: "eth0",
	}
	u, err := e.ToSystemdUnit()
	if err != nil {
		t.Fatal(err)
	}
	r := unit.Serialize(u)
	o, _ := ioutil.ReadAll(r)
	exp, err := ioutil.ReadFile("fixture/staticonly.service")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(exp, o) {
		t.Errorf("expected \n %s \n Got \n %s", string(exp), string(o))
	}

	// static and dhcp
	e.DHCP = true
	u, err = e.ToSystemdUnit()
	if err != nil {
		t.Fatal(err)
	}
	r = unit.Serialize(u)
	o, _ = ioutil.ReadAll(r)
	exp, err = ioutil.ReadFile("fixture/staticanddhcp.service")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(exp, o) {
		t.Errorf("expected \n %s \n Got \n %s", string(exp), string(o))
	}

	// static only
	e.Static = nil
	u, err = e.ToSystemdUnit()
	if err != nil {
		t.Fatal(err)
	}
	r = unit.Serialize(u)
	o, _ = ioutil.ReadAll(r)
	exp, err = ioutil.ReadFile("fixture/dhcponly.service")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(exp, o) {
		t.Errorf("expected \n %s \n Got \n %s", string(exp), string(o))
	}

}
