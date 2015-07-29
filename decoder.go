package xmlrpc

import (
	"encoding/xml"
	"fmt"
	"io"
)

type decoder struct {
	r    io.Reader
	body []byte
}

type requestHeader struct {
	MethodName string  `xml:"methodName"`
	Params     *params `xml:"params"`
}

type responseHeader struct {
	Params *params `xml:"params"`
	Fault  *fault  `xml:"fault>value>struct"`
}

type params struct {
	Body []byte `xml:",innerxml"`
}

type fault struct {
	Members []member `xml:"member"`
}

type member struct {
	Name   string `xml:"name"`
	Code   int    `xml:"value>int"`
	String string `xml:"value>string"`
}

type FaultResponseError struct {
	Code   int
	String string
}

func (e *FaultResponseError) Error() string {
	return fmt.Sprintf("Code: %d, String: %s", e.Code, e.String)
}

var _ error = &FaultResponseError{}

func newDecoder(r io.Reader) *decoder {
	d := &decoder{
		r: r,
	}
	return d
}

// DecodeRequestHeader returns method name and error.
func (d *decoder) DecodeRequestHeader() (string, error) {
	dec := xml.NewDecoder(d.r)
	r := &requestHeader{}
	err := dec.Decode(r)
	if err != nil {
		return "", err
	}

	d.body = r.Params.Body

	return r.MethodName, nil
}

func (d *decoder) DecodeResponseHeader() error {
	dec := xml.NewDecoder(d.r)
	r := &responseHeader{}
	err := dec.Decode(r)
	if err != nil {
		return err
	}

	if r.Fault != nil {
		err := &FaultResponseError{}
		for _, m := range r.Fault.Members {
			switch m.Name {
			case "faultCode":
				err.Code = m.Code
			case "faultString":
				err.String = m.String
			}
		}
		return err
	}

	d.body = r.Params.Body
	return nil
}
