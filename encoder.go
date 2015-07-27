package xmlrpc

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"time"
)

const iso8601 = "20060102T15:04:05"

type encoder struct {
	w   io.Writer
	err error
}

func newEncoder(w io.Writer) *encoder {
	return &encoder{
		w: w,
	}
}

func (e *encoder) Write(b []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}
	i, err := e.w.Write(b)
	e.err = err
	return i, err
}

var _ io.Writer = &encoder{}

func (e *encoder) EncodeRequest(method string, params interface{}) error {
	fmt.Fprint(e, `<?xml version="1.0"?>`)
	fmt.Fprint(e, "<methodCall>")
	fmt.Fprintf(e, "<methodName>%s</methodName>", method)
	fmt.Fprint(e, "<params>")

	e.encodeParams(params)

	fmt.Fprint(e, "</params>")
	fmt.Fprint(e, "</methodCall>")
	return e.err
}

func (e *encoder) EncodeResponse(params interface{}) error {
	fmt.Fprint(e, `<?xml version="1.0"?>`)
	fmt.Fprint(e, "<methodResponse>")
	fmt.Fprint(e, "<params>")

	e.encodeParams(params)

	fmt.Fprint(e, "</params>")
	fmt.Fprint(e, "</methodResponse>")
	return e.err
}

func (e *encoder) encodeParams(params interface{}) {
	val := reflect.ValueOf(params)
	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			e.encodeParam(val.Index(i))
		}
	} else {
		e.encodeParam(val)
	}
}

func (e *encoder) encodeParam(val reflect.Value) {
	if e.isNil(val) {
		return
	}

	fmt.Fprint(e, "<param>")
	e.encodeValue(val)
	fmt.Fprint(e, "</param>")
}

// TODO: base64
func (e *encoder) encodeValue(val reflect.Value) {
	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	fmt.Fprint(e, "<value>")

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(e, "<int>%d</int>", val.Int())
	case reflect.Bool:
		var n int
		if val.Bool() {
			n = 1
		} else {
			n = 0
		}
		fmt.Fprintf(e, "<boolean>%d</boolean>", n)
	case reflect.String:
		fmt.Fprint(e, "<string>")
		xml.Escape(e, []byte(val.String()))
		fmt.Fprint(e, "</string>")
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(e, "<double>%f</double>", val.Float())
	case reflect.Struct:
		if t, isTime := val.Interface().(time.Time); isTime {
			fmt.Fprintf(e, "<dateTime.iso8601>%s</dateTime.iso8601>", t.Format(iso8601))
		} else {
			e.encodeStruct(val)
		}
	case reflect.Slice:
		e.encodeSlice(val)
	}

	fmt.Fprint(e, "</value>")
}

func (e *encoder) encodeStruct(val reflect.Value) {
	fmt.Fprint(e, "<struct>")

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		v := val.FieldByName(f.Name)
		if e.isNil(v) {
			continue
		}

		fmt.Fprintf(e, "<member>")

		name := f.Tag.Get("xmlrpc")
		if name == "" {
			name = f.Name
		}
		fmt.Fprintf(e, "<name>%s</name>", name)
		e.encodeValue(v)

		fmt.Fprintf(e, "</member>")
	}

	fmt.Fprint(e, "</struct>")
}

func (e *encoder) encodeSlice(val reflect.Value) {
	fmt.Fprint(e, "<array><data>")

	for i := 0; i < val.Len(); i++ {
		v := val.Index(i)
		if e.isNil(v) {
			continue
		}
		e.encodeValue(v)
	}

	fmt.Fprint(e, "</data></array>")
}

func (_ *encoder) isNil(val reflect.Value) bool {
	return (val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface) && val.IsNil()
}
