package ums

import (
	"strings"
	"testing"
)

func TestErrorPanicOnNil(t *testing.T) {
	var seen_panic bool

	provoke := func() {
		defer func() {
			if r := recover(); r != nil {
				if msg, ok := r.(string); ok && msg == "this is no import error, you are using me wrong" {
					seen_panic = true
				}
			}
		}()
		test := &ImportError{}
		_ = test.Error()

	}
	provoke()

	if !seen_panic {
		t.Errorf("Sorry no panic received on invalid usage")
	} else {
		t.Log("Caught panic on nil")
	}
}

func TestErrorPanicOnNoError(t *testing.T) {
	var seen_panic bool

	provoke := func() {
		defer func() {
			if r := recover(); r != nil {
				if msg, ok := r.(string); ok && msg == "Unknown error, you are using me wrong" {
					seen_panic = true
				}
			}
		}()
		data := strings.NewReader(goodImport)
		imp, err := NewDoc(data)
		if err != nil {
			t.Fatal(err)
		}

		test := &ImportError{Import: imp}
		_ = test.Error()
	}
	provoke()

	if !seen_panic {
		t.Errorf("Sorry no panic received on error, that is no error")
	} else {
		t.Log("Caught panic on no error")
	}
}

func TestBadImportError(t *testing.T) {
	data := strings.NewReader(badImport)
	imp, err := NewDoc(data)
	if err != nil {
		t.Fatal(err)
	}
	if imp.Successful() {
		t.Error("badImport is seen wrongly considered a success")
	}
}

func TestGoodImportNoError(t *testing.T) {
	data := strings.NewReader(goodImport)
	imp, err := NewDoc(data)
	if err != nil {
		t.Error(err)
	}
	if !imp.Successful() {
		t.Error("goodImport is wrongly considered a failure")
	}
}

var badImport = `<?xml version="1.0" encoding="UTF-8" standalone="no" ?><!DOCTYPE import SYSTEM "https://xml.fagms.net/XMLInterface.dtd"><import>
<transaction_id>12345678</transaction_id>
<respond_to>recipient@example.com</respond_to>
<response_type>SHORT</response_type>
<language>GERMAN</language>
<login>
 <net>4711</net>
 <uid>Admin</uid>
 <pwd>***</pwd>
</login>
<orders>

 <insert_member_csv><ERR_NO>600</ERR_NO><ERR_TXT>Your Operation was executed with the following results: '
line   1 Member already exists. Updated successfully.
line 2319 Insert Failed. Wrong Number of Values or empty Line.
'</ERR_TXT></insert_member_csv>

<ERR_NO>3</ERR_NO><ERR_TXT>Errors occured processing your Request, See above for Details.</ERR_TXT></orders>
</import>
`

var goodImport = `<?xml version="1.0" encoding="UTF-8" standalone="no" ?><!DOCTYPE import SYSTEM "http://admin-beta.fagms.net/XMLInterface.dtd"><import>
<transaction_id>12345678</transaction_id>
<respond_to>recipient@example.com</respond_to>
<response_type>SHORT</response_type>
<language>GERMAN</language>
<login>
 <net>4711</net>
 <uid>Admin</uid>
 <pwd>***</pwd>
</login>
<orders>

 <insert_member_csv><ERR_NO>600</ERR_NO><ERR_TXT>Your Operation was executed with the following results: '
line   1 Member already exists. Updated successfully.
line 33940 Inserted successfully.
'</ERR_TXT></insert_member_csv>

<ERR_NO>2</ERR_NO><ERR_TXT>No Errors occured processing your Request!</ERR_TXT></orders>
</import>

`
