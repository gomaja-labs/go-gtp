package ie_test

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// FuzzAccessors parses an IE and calls every exported zero-argument method on
// it. Any panic is a defect: accessors must return an error on bad payloads.
func FuzzAccessors(f *testing.F) {
	f.Add([]byte{0x01, 0x00, 0x01, 0x00, 0x00})
	f.Fuzz(func(t *testing.T, b []byte) {
		i, err := ie.Parse(b)
		if err != nil {
			return
		}
		v := reflect.ValueOf(i)
		typ := v.Type()
		for m := 0; m < typ.NumMethod(); m++ {
			meth := typ.Method(m)
			if meth.Type.NumIn() != 1 { // receiver only
				continue
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("%s() panicked on % x: %v\n%s", meth.Name, b, r, debug.Stack())
					}
				}()
				_ = fmt.Sprint(v.Method(m).Call(nil))
			}()
		}
	})
}

func TestNilAccessors(t *testing.T) {
	var i *ie.IE
	v := reflect.ValueOf(i)
	for m := 0; m < v.NumMethod(); m++ {
		method := v.Type().Method(m)
		if method.Type.NumIn() != 1 {
			continue
		}
		t.Run(method.Name, func(t *testing.T) {
			out := v.Method(m).Call(nil)
			if len(out) > 1 && method.Type.Out(len(out)-1).Name() == "error" && out[len(out)-1].IsNil() {
				t.Fatal("nil accessor returned success")
			}
		})
	}
}
