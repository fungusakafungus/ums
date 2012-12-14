package ums

import (
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

var XMLInterfaceSender = "<xmlinterface@fagms.net>"
var ErrWrongSender = errors.New("Mail is not from API sender " + XMLInterfaceSender)
var ErrNoMimeMail = errors.New("Mail is not a mime mail")
var ErrNoBoundary = errors.New("Mail contains no boundary")
var ErrNoMultipart = errors.New("Mail contains multipart/mixed type")
var ErrNoResult = errors.New("Mail contains no result.xml")

// Reader which strips \n (linefeed) characters from a stream
type lineFeedFilter struct {
	r io.Reader
}

// Read any bytes, which are no linefeeds
func (r *lineFeedFilter) Read(p []byte) (n int, err error) {
	line := make([]byte, len(p), cap(p))
	for {
		remain := len(p)
		if remain < len(line) {
			line = line[:remain]
		}
		ln, err := r.r.Read(line)
		if ln > 0 {
			copied := 0
			scanned := 0
			for i, b := range line[:ln] {
				if b == '\n' {
					slice := line[scanned:i]
					copy(p, slice)
					p = p[len(slice):]
					copied += len(slice)
					scanned += len(slice) + 1
				}
			}
			slice := line[scanned:ln]
			copy(p, slice)
			p = p[len(slice):]
			copied += len(slice)
			n += copied
		}
		if ln == 0 || err != nil {
			return n, err
		}
	}
	return n, nil
}

// Extracts the UMS userdata import result from an email message in MIME format.
// The email message is directly read from r.
func ExtractImportResult(r io.Reader) (imp *Import, err error) {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return nil, err
	}

	if from := msg.Header.Get("From"); !strings.Contains(from, XMLInterfaceSender) {
		return nil, ErrWrongSender
	}
	if msg.Header.Get("MIME-Version") != "1.0" {
		return nil, ErrNoMimeMail
	}
	ct := msg.Header.Get("Content-Type")
	media, params, err := mime.ParseMediaType(ct)
	boundary, ok := params["boundary"]
	if !ok {
		return nil, ErrNoBoundary
	}
	if media != "multipart/mixed" {
		return nil, ErrNoMultipart
	}

	parts := multipart.NewReader(msg.Body, boundary)
	part, err := parts.NextPart()
	for ; err == nil; part, err = parts.NextPart() {
		media, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			return nil, err
		}
		if media != "text/xml" || part.FileName() != "result.xml" {
			continue
		}

		switch part.Header.Get("Content-Transfer-Encoding") {
		case "base64":
			imp, err = NewDoc(base64.NewDecoder(base64.StdEncoding, &lineFeedFilter{r: part}))
		default:
			imp, err = NewDoc(part)
		}

		return imp, err
	}
	return nil, ErrNoResult
}
