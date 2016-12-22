package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/kr/pretty"
)

func TestLoadAPFromSrc(t *testing.T) {
	b, err := ioutil.ReadFile("fixture/create_ap.conf")
	if err != nil {
		t.Fatal(err)
	}
	a, err := LoadAPFromConf(b)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	_, err = a.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}
	j, err := json.MarshalIndent(a.State(), "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("fixture/create_ap.json", j, 0644)

	b, err = ioutil.ReadFile("fixture/create_ap.json")
	if err != nil {
		t.Fatal(err)
	}
	a, err = LoadFromJSON(b)
	if err != nil {
		t.Fatal(err)
	}

	pretty.Println(a)
}
