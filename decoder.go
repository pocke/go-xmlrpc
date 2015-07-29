package xmlrpc

import (
	"encoding/xml"
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

func (f *fault) err() error {
	err := &FaultResponseError{}
	for _, m := range f.Members {
		switch m.Name {
		case "faultCode":
			err.Code = m.Code
		case "faultString":
			err.String = m.String
		}
	}
	return err
}

type member struct {
	Name   string `xml:"name"`
	Code   int    `xml:"value>int"`
	String string `xml:"value>string"`
}

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
		return r.Fault.err()
	}

	d.body = r.Params.Body
	return nil
}
