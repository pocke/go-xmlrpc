package xmlrpc

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

type req struct {
	Str   *string `xmlrpc:"str"`
	Int   int
	Slice []int
}

func TestEncodeRequest(t *testing.T) {
	assert := func(method string, params interface{}, expected string) {
		buf := bytes.NewBuffer([]byte{})
		e := newEncoder(buf)

		err := e.EncodeRequest(method, params)
		if err != nil {
			t.Fatal(err)
		}

		s := buf.String()
		if s != expected {
			t.Fatalf(`
Expected: %s,
but got   %s`, expected, s)
		}
	}

	assert("nya", []int{1, 2, 3}, createEncoded("nya", "<param><value><int>1</int></value></param><param><value><int>2</int></value></param><param><value><int>3</int></value></param>"))
	assert("foo", 10, createEncoded("foo", "<param><value><int>10</int></value></param>"))
	assert("foo", -10, createEncoded("foo", "<param><value><int>-10</int></value></param>"))
	assert("bar", true, createEncoded("bar", "<param><value><boolean>1</boolean></value></param>"))
	assert("hoge", false, createEncoded("hoge", "<param><value><boolean>0</boolean></value></param>"))
	assert("fuga", "piyo", createEncoded("fuga", "<param><value><string>piyo</string></value></param>"))
	assert("fuga", "<piyo>", createEncoded("fuga", "<param><value><string>&lt;piyo&gt;</string></value></param>"))
	assert("poyo", 3.141592, createEncoded("poyo", "<param><value><double>3.141592</double></value></param>"))
	n := time.Now()
	assert("poyo", n, createEncoded("poyo", "<param><value><dateTime.iso8601>"+n.Format(iso8601)+"</dateTime.iso8601></value></param>"))
	s := "hoge"
	assert("nyoa", &req{
		Str:   &s,
		Int:   10,
		Slice: []int{1, 2, 3},
	}, createEncoded("nyoa", "<param><value><struct><member><name>str</name><value><string>hoge</string></value></member><member><name>Int</name><value><int>10</int></value></member><member><name>Slice</name><value><array><data><value><int>1</int></value><value><int>2</int></value><value><int>3</int></value></data></array></value></member></struct></value></param>"))
	assert("nyoa", &req{
		Int:   10,
		Slice: []int{1, 2, 3},
	}, createEncoded("nyoa", "<param><value><struct><member><name>Int</name><value><int>10</int></value></member><member><name>Slice</name><value><array><data><value><int>1</int></value><value><int>2</int></value><value><int>3</int></value></data></array></value></member></struct></value></param>"))
}

func createEncoded(method, params string) string {
	return fmt.Sprintf(`<?xml version="1.0"?><methodCall><methodName>%s</methodName><params>%s</params></methodCall>`, method, params)
}

func TestEncodeResponse(t *testing.T) {
	assert := func(params interface{}, expected string) {
		buf := bytes.NewBuffer([]byte{})
		e := newEncoder(buf)

		err := e.EncodeResponse(params)
		if err != nil {
			t.Fatal(err)
		}

		s := buf.String()
		if s != expected {
			t.Fatalf(`
Expected: %s,
but got   %s`, expected, s)
		}
	}

	assert(10, `<?xml version="1.0"?><methodResponse><params><param><value><int>10</int></value></param></params></methodResponse>`)
}
