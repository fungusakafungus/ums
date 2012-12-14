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

// verbosity of operation
var verbose bool

// small wrapper to abstract out verbosity and autoappend \n
func info(format string, a ...interface{}) {
	if verbose {
		if !strings.HasSuffix(format, "\n") {
			format = format + "\n"
		}
		log.Printf(format, a...)
	}
}

// whether we delete an email after processing it
var delete_after bool

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

	if !CheckArguments() {
		state, message := nagios.UNKNOWN, "invalid or missing comandline parameters. Use -h for help"
		nagios.Exit(state, message)
	}

	// The real plugin
	state, message := ProcessMails()
	nagios.Exit(state, message)
}

// Any argument checking goes here
func CheckArguments() bool {
	return Account.Email != "" && Account.Password != ""
}

// Process mails at configured POP3 account
// and return a nagios state, including message reflecting what happened
func ProcessMails() (state nagios.Status, message string) {

	hp := net.JoinHostPort(Account.Host, Account.Port)
	info("Dialing %v ...", hp)
	client, err := pop3.DialTLS(hp)
	if err != nil {
		return nagios.UNKNOWN, err.Error()
	}

	// Send QUIT, so DELEted mails get expunged from the mailbox
	defer client.Quit()

	info("Authorizing as %v ...", Account.Email)
	if err = client.Auth(Account.Email, Account.Password); err != nil {
		return nagios.UNKNOWN, err.Error()
	}

	info("Retrieving messages ids ...")
	msgId, _, err := client.ListAll()
	if err != nil {
		return nagios.UNKNOWN, err.Error()
	}
	info("Got %v message ids", len(msgId))

	for i, id := range msgId {
		info("Retrieving message %v/%v id = %v", i+1, len(msgId), id)
		msg, err := client.Retr(id)
		if err != nil {
			return nagios.UNKNOWN, err.Error()
		}
		info("Retrieved %v bytes", len(msg))
		imp, err := ums.ExtractImportResult(strings.NewReader(msg))
		if err != nil {
			// Not for us
			if err == ums.ErrWrongSender {
				info("Message %v is not for UMS check", msgId)
				continue
			}
			return nagios.UNKNOWN, err.Error()
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
			return nagios.OK, "UMS import successful"
		} else {
			err := &ums.ImportError{Import: imp}
			return nagios.WARNING, err.Error()
		}
	}

	return nagios.OK, "No mails for us"
}
