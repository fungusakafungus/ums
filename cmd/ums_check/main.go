package main

import (
	"flag"
	"log"
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

var verbose bool

var delete_after bool

func info(format string, a ...interface{}) {
	if verbose {
		if !strings.HasSuffix(format, "\n") {
			format = format + "\n"
		}
		log.Printf(format, a...)
	}
}

func main() {
	// disable date/time prefixes. Use logger, if you need timestamps
	log.SetFlags(0)

	flag.StringVar(&Account.Email, "email", "", "email address")
	flag.StringVar(&Account.Password, "password", "", "password")
	flag.StringVar(&Account.Host, "host", "localhost", "POP3 host")
	flag.StringVar(&Account.Port, "port", "pop3s", "POP3 via SSL port")
	flag.BoolVar(&verbose, "verbose", false, "verbose logging on stderr")
	flag.BoolVar(&delete_after, "delete_after", false, "delete email after processing it")
	flag.Parse()

	if Account.Email == "" || Account.Password == "" {
		defer nagios.Exit(nagios.UNKNOWN, "invalid or missing comandline parameters. Use -h for help")
		return
	}

	hp := net.JoinHostPort(Account.Host, Account.Port)
	info("Dialing %v ...", hp)
	client, err := pop3.DialTLS(hp)
	if err != nil {
		defer nagios.Exit(nagios.UNKNOWN, err.Error())
		return
	}

	defer client.Quit()

	info("Authorizing as %v ...", Account.Email)
	if err = client.Auth(Account.Email, Account.Password); err != nil {
		defer nagios.Exit(nagios.UNKNOWN, err.Error())
		return
	}

	info("Retrieving messages ids ...")
	msgId, _, err := client.ListAll()
	if err != nil {
		defer nagios.Exit(nagios.UNKNOWN, err.Error())
		return
	}
	info("Got %v message ids", len(msgId))

	for i, id := range msgId {
		info("Retrieving message %v/%v id = %v", i+1, len(msgId), id)
		msg, err := client.Retr(id)
		if err != nil {
			defer nagios.Exit(nagios.UNKNOWN, err.Error())
			return
		}
		info("Retrieved %v bytes", len(msg))
		imp, err := ums.ExtractImportResult(strings.NewReader(msg))
		if err != nil {
			// Not for us
			if err == ums.ErrWrongSender {
				info("Message %v is not for UMS check", msgId)
				continue
			}
			defer nagios.Exit(nagios.UNKNOWN, err.Error())
			return
		}

		if delete_after {
			err := client.Dele(id)

			if err != nil {
				info("Error deleting processed mail %v, Error %v", id, err)
			} else {
				info("Deleted processed mail %v", id)
			}
		}

		if imp.Successful() {
			defer nagios.Exit(nagios.OK, "UMS import successful")
			return
		} else {
			err := &ums.ImportError{Import: imp}
			defer nagios.Exit(nagios.WARNING, err.Error())
			return
		}
	}

	defer nagios.Exit(nagios.OK, "No mails for us")
}
