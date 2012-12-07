package ums

import (
	"os"
	"testing"
)

func TestUms(t *testing.T) {
	file, err := os.Open("testdata/result.xml")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Open works\n")
	}
	_, err = NewDoc(file)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Initial parsing works\n")
	}
}

func TestValidData(t *testing.T) {
	file, err := os.Open("testdata/result.xml")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Open works\n")
	}
	imp, err := NewDoc(file)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Initial parsing works\n")
	}
	if !imp.Successful() {
		t.Error(&ImportError{imp})
	} else {
		t.Logf("import data works\n")
	}
}
