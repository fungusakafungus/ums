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

	dec, err := ExtractImportResult(file)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Initial parsing works\n")
	}
	if dec.Orders[0].ErrorNo != 3 {
		t.Errorf("Expected: <orders><ERR_NO>3</ERR_NO></orders>, "+
			"got <orders><ERR_NO>%v</ERR_NO></orders>\n", dec.Orders[0].ErrorNo)
	} else {
		t.Logf("Total err: %s,%v\n", dec.Orders[0].ErrorText, dec.Orders[0].ErrorNo)
	}
	if dec.Orders[0].Inserts[0].ErrorNo != 600 {
		t.Errorf("Expected: <orders><insert_member_csv><ERR_NO>600</ERR_NO></insert_member_csv></orders>, "+
			"got <orders><insert_member_csv><ERR_NO>%v</ERR_NO></insert_member_csv></orders>\n", dec.Orders[0].Inserts[0].ErrorNo)
	} else {
		msg := strings.SplitN(dec.Orders[0].Inserts[0].ErrorText, "'", 2)
		t.Logf("Insert err: %s,%v\n", msg[0], dec.Orders[0].Inserts[0].ErrorNo)
	}
}
