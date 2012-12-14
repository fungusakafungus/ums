package main

import (
	"flag"
	"net"
	"strings"

	"github.com/laziac/go-nagios/nagios"
	"github.com/nightlyone/ums"

	pop3 "github.com/bytbox/go-pop3"
)

// describes one POP3 account
var Account struct {
	Email    string
	Password string
	Host     string
	Port     string
}

func main() {
	flag.StringVar(&Account.Email, "email", "", "email address")
	flag.StringVar(&Account.Password, "password", "", "password")
	flag.StringVar(&Account.Host, "host", "localhost", "POP3 host")
	flag.StringVar(&Account.Port, "port", "pop3s", "POP3 via SSL port")
	flag.Parse()

	if Account.Email == "" || Account.Password == "" {
		defer nagios.Exit(nagios.UNKNOWN, "invalid or missing comandline parameters. Use -h for help")
		return
	}

	client, err := pop3.DialTLS(net.JoinHostPort(Account.Host, Account.Port))
	if err != nil {
		defer nagios.Exit(nagios.UNKNOWN, err.Error())
		return
	}

	defer client.Quit()

	if err = client.Auth(Account.Email, Account.Password); err != nil {
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
