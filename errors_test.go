package xmlrpc

import "testing"

func TestFaultResponseError(t *testing.T) {
	err := &FaultResponseError{
		Code:   1,
		String: "hoge",
	}
	err.Error()
}
