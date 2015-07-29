package xmlrpc

import (
	"strings"
	"testing"
)

const requestXML = `
<?xml version="1.0"?>
<methodCall>
  <methodName>examples.getStateName</methodName>
  <params>
    <param>
      <value><i4>41</i4></value>
    </param>
  </params>
</methodCall>
`

const responseXML = `
<?xml version="1.0"?>
<methodResponse>
  <params>
    <param>
      <value><string>South Dakota</string></value>
    </param>
  </params>
</methodResponse>
`

const faultResponseXML = `
<?xml version="1.0"?>
<methodResponse>
  <fault>
    <value>
      <struct>
        <member>
          <name>faultCode</name>
          <value><int>4</int></value>
        </member>
        <member>
          <name>faultString</name>
          <value><string>Too many parameters.</string></value>
        </member>
      </struct>
    </value>
  </fault>
</methodResponse>
`

func TestDecodeRequestHeader(t *testing.T) {
	r := strings.NewReader(requestXML)
	dec := newDecoder(r)
	m, err := dec.DecodeRequestHeader()
	if err != nil {
		t.Fatal(err)
	}
	if m != "examples.getStateName" {
		t.Fatalf("Cannot get method name. got %s", m)
	}
	if len(dec.body) == 0 {
		t.Fatal("Body should be set")
	}
}

func TestDecodeResponseHeader(t *testing.T) {
	r := strings.NewReader(responseXML)
	dec := newDecoder(r)
	err := dec.DecodeResponseHeader()
	if err != nil {
		t.Fatal(err)
	}
	if len(dec.body) == 0 {
		t.Fatal("Body should be set")
	}
}

func TestDecodeResponseHeaderWhenFault(t *testing.T) {
	r := strings.NewReader(faultResponseXML)
	dec := newDecoder(r)
	err := dec.DecodeResponseHeader()
	if err == nil {
		t.Fatal("Should be error, but got nil")
	}

	if _, ok := err.(*FaultResponseError); !ok {
		t.Fatal("Should be FaultResponseError, but not got.")
	}
}
