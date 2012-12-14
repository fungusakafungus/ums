package main

import (
	"flag"
	"net"
	"strings"

	"github.com/laziac/go-nagios/nagios"
	"github.com/nightlyone/ums"

	pop3 "github.com/bytbox/go-pop3"
)

var acct struct {
	Email    string
	Password string
	Host     string
	Port     string
}

func main() {
	flag.StringVar(&acct.Email, "email", "", "email address")
	flag.StringVar(&acct.Password, "password", "", "password")
	flag.StringVar(&acct.Host, "host", "localhost", "POP3 host")
	flag.StringVar(&acct.Port, "port", "pop3s", "POP3 via SSL port")
	flag.Parse()

	if acct.Email == "" || acct.Password == "" {
		defer nagios.Exit(nagios.UNKNOWN, "invalid or missing comandline parameters. Use -h for help")
		return
	}

	client, err := pop3.DialTLS(net.JoinHostPort(acct.Host, acct.Port))
	if err != nil {
		defer nagios.Exit(nagios.UNKNOWN, err.Error())
		return
	}

	defer client.Quit()

	if err = client.Auth(acct.Email, acct.Password); err != nil {
		defer nagios.Exit(nagios.UNKNOWN, err.Error())
		return
	}

	msgId, _, err := client.ListAll()
	if err != nil {
		defer nagios.Exit(nagios.UNKNOWN, err.Error())
		return
	}

	for _, id := range msgId {
		msg, err := client.Retr(id)
		if err != nil {
			defer nagios.Exit(nagios.UNKNOWN, err.Error())
			return
		}
		imp, err := ums.ExtractImportResult(strings.NewReader(msg))
		if err != nil {
			// Not for us
			if err == ums.ErrWrongSender {
				continue
			}
			defer nagios.Exit(nagios.UNKNOWN, err.Error())
			return
		}
		if imp.Successful() {
			defer nagios.Exit(nagios.OK, "UMS import successful")
			return
		} else {
			err := &ums.ImportError{imp}
			defer nagios.Exit(nagios.WARNING, err.Error())
			return
		}
	}

	defer nagios.Exit(nagios.OK, "No mails for us")
}
