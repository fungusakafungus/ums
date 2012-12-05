package ums

import (
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
)

var ErrNoMimeMail = errors.New("Mail is not a mime mail")
var ErrNoBoundary = errors.New("Mail contains no boundary")
var ErrNoMultipart = errors.New("Mail contains multipart/mixed type")

// Extracts the UMS userdata import result from an email message in MIME format.
// The email message is directly read from r.
func ExtractImportResult(r io.Reader) (imp *Import, err error) {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return nil, err
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
		imp, err = NewDoc(part)
		return imp, err
	}
	return nil, io.EOF
}
