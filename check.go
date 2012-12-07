package ums

import (
	"fmt"
	"strings"
)

type ImportError struct {
	*Import
}

// did the import succeed?
func (imp *Import) Successful() bool {
	return imp.Orders[0].ErrorNo == 2 && imp.Orders[0].Inserts[0].ErrorNo == 600
}

func (impErr *ImportError) Error() string {
	imp := impErr.Import
	if imp == nil {
		panic("this is no import error, you are using me wrong")
	}

	if len(imp.Orders) == 0 {
		return "No orders in result"
	}
	if len(imp.Orders[0].Inserts) == 0 {
		return "No inserts in order result"
	}

	msg := strings.SplitN(imp.Orders[0].Inserts[0].ErrorText, "'", 2)

	// Magic value, TODO get meaning of magic value
	if imp.Orders[0].ErrorNo != 2 {
		return fmt.Sprintf("Import failed with %v, %s, wrong call", imp.Orders[0].ErrorNo, msg[0])
	}

	// Magic value, TODO get meaning of magic value
	if imp.Orders[0].Inserts[0].ErrorNo != 600 {
		return fmt.Sprintf("inserts failed with %v, %s", imp.Orders[0].Inserts[0].ErrorNo, msg[0])
	}
	panic("Unknown error, you are using me wrong")
}
