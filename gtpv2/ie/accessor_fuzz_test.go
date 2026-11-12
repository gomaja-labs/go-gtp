// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie_test

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
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
		// Extended and negotiated Large APCO need the NAS direction to select length widths.
		for _, direction := range []ie.PCODirection{ie.PCOMSToNetwork, ie.PCONetworkToMS} {
			_, _ = i.ExtendedProtocolConfigurationOptions(direction)
			_, _ = i.LargeAdditionalProtocolConfigurationOptions(direction)
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

// FuzzScalarConstructors checks exported constructors with arbitrary IE types,
// including grouped types which must be rejected without panicking.
func FuzzScalarConstructors(f *testing.F) {
	f.Add(ie.BearerContext, uint64(1))
	f.Add(uint8(250), uint64(0x123456789abcdef0))
	f.Fuzz(func(t *testing.T, typ uint8, value uint64) {
		for width, i := range map[int]*ie.IE{
			1: ie.NewUint8IE(typ, uint8(value)),
			2: ie.NewUint16IE(typ, uint16(value)),
			4: ie.NewUint32IE(typ, uint32(value)),
			8: ie.NewUint64IE(typ, value),
		} {
			if i.IsGrouped() {
				if !errors.Is(i.Err(), ie.ErrInvalidType) {
					t.Fatalf("grouped scalar type %d: %v", typ, i.Err())
				}
				continue
			}
			b, err := i.Marshal()
			if err != nil || len(b) != 4+width || i.MarshalLen() != len(b) {
				t.Fatalf("scalar width %d: %x, %v", width, b, err)
			}
			for n := 0; n < width; n++ {
				if b[len(b)-1-n] != byte(value>>uint(8*n)) {
					t.Fatalf("scalar width %d lost data: %x", width, b)
				}
			}
		}
	})
}
