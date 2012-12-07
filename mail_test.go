package ums

import (
	"os"
	"strings"
	"testing"
)

func TestUmsMail(t *testing.T) {
	file, err := os.Open("testdata/result.msg")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Open works\n")
	}

	imp, err := ExtractImportResult(file)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Initial parsing works\n")
	}
	if imp.Sucessful() {
		t.Errorf("expected import failure here")
	}
}
