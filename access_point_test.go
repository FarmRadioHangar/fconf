package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestLoadAPFromSrc(t *testing.T) {
	b, err := ioutil.ReadFile("fixture/create_ap.conf")
	if err != nil {
		t.Fatal(err)
	}
	a, err := LoadAPFromSrc(b)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	_, err = a.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}
}
