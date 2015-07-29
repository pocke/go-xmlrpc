package xmlrpc

import "fmt"

type FaultResponseError struct {
	Code   int
	String string
}

func (e *FaultResponseError) Error() string {
	return fmt.Sprintf("Code: %d, String: %s", e.Code, e.String)
}

var _ error = &FaultResponseError{}
