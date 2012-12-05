package ums

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
)

type Import struct {
	XMLName       xml.Name `xml:"import"`
	TransactionId uint32   `xml:"transaction_id"`
	RespondTo     string   `xml:"respond_to"`
	ResponseType  string   `xml:"response_type"`
	Orders        []Order  `xml:"orders"`
}

type Error struct {
	ErrorText string `xml:"ERR_TXT"`
	ErrorNo   int32  `xml:"ERR_NO"`
}

type Order struct {
	Inserts []Insert `xml:"insert_member_csv"`
	Error
}

type Insert struct {
	Error
	Stuff string `xml:",chardata"`
}

// parses an UMS userdata import result and returns it
func NewDoc(r io.Reader) (imp *Import, err error) {
	imp = &Import{}
	err = xml.NewDecoder(bufio.NewReader(r)).Decode(imp)
	return imp, err
}

func (a *Import) Dump(w io.Writer) {
	_, err := fmt.Fprintf(w, "%+v", a)
	if err != nil {
		panic(err)
	}
}
