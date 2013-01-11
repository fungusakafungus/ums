ums
===

Tools for handling UMS mail API responses.


[![Build Status][1]][2]

[1]: https://secure.travis-ci.org/nightlyone/ums.png
[2]: http://travis-ci.org/nightlyone/ums



install
-------
Install [Go 1][3], either [from source][4] or [with a prepackaged binary][5].

Then run

	go get github.com/nightlyone/ums

To build the nagios check

	go get github.com/nightlyone/ums/cmd/check_ums
	go install github.com/nightlyone/ums/cmd/check_ums

Get usage of the nagios check

	$GOPATH/bin/check_ums -h

[3]: http://golang.org
[4]: http://golang.org/doc/install/source
[5]: http://golang.org/doc/install

LICENSE
-------
BSD

documentation
-------------
[package documentation at go.pkgdoc.org](http://go.pkgdoc.org/github.com/nightlyone/ums)

contributing
============

Contributions are welcome. Please open an issue or send me a pull request for a dedicated branch.
Make sure the git commit hooks show it works.

git commit hooks
-----------------------
enable commit hooks via

        cd .git ; rm -rf hooks; ln -s ../git-hooks hooks ; cd ..

usage of nagios plugin
----------------------
production usage, deleting all processed mails

	define command{
		command_name    check_ums
		command_line    /usr/lib/nagios/plugins/check_ums -delete_after -host=$ARG2$ -email=$ARG1$ -password=$ARG3$
		}

test usage, NOT deleting processed mails

	define command{
		command_name    check_ums_test
		command_line    /usr/lib/nagios/plugins/check_ums -host=$ARG2$ -email=$ARG1$ -password=$ARG3$
		}
