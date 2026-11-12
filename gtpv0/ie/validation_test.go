// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"errors"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-gtp/utils"
)

func TestNilIEMarshal(t *testing.T) {
	var i *IE
	if _, err := i.Marshal(); !errors.Is(err, ErrIENotFound) {
		t.Fatalf("nil Marshal: %v", err)
	}
	if err := i.MarshalTo(nil); !errors.Is(err, ErrIENotFound) {
		t.Fatal(err)
	}
	if i.MarshalLen() != 0 {
		t.Fatal("nil length")
	}
	if !errors.Is(i.Err(), ErrIENotFound) {
		t.Fatal(i.Err())
	}
	if err := i.UnmarshalBinary(nil); !errors.Is(err, ErrIENotFound) {
		t.Fatal(err)
	}
	i.SetLength()
	if i.Name() != "Undefined" || i.String() != "nil" {
		t.Fatal("nil diagnostic")
	}
	if _, err := i.IMSI(); !errors.Is(err, ErrIENotFound) {
		t.Fatal(err)
	}
	if i.MustIMSI() != "" {
		t.Fatal("nil MustIMSI")
	}
}

func TestConstructionError(t *testing.T) {
	i := NewIMSI("12x")
	if i == nil {
		t.Fatal("constructor discarded error")
	}
	if !errors.Is(i.Err(), utils.ErrMalformed) {
		t.Fatal(i.Err())
	}
	if !strings.Contains(i.Err().Error(), "IE type 2") && !strings.Contains(i.Err().Error(), "IE type 1:") {
		t.Fatal(i.Err())
	}
	if _, err := i.Marshal(); !errors.Is(err, utils.ErrMalformed) {
		t.Fatal(err)
	}
	if err := i.MarshalTo(make([]byte, 50)); !errors.Is(err, utils.ErrMalformed) {
		t.Fatal(err)
	}
	if _, err := i.IMSI(); !errors.Is(err, utils.ErrMalformed) {
		t.Fatal(err)
	}
	if i.MarshalLen() != 0 {
		t.Fatal("construction error length must be zero")
	}
	// Reuse must erase the construction error and all prior payload state.
	good, err := NewIMSI("123456789012345").Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := i.UnmarshalBinary(good); err != nil {
		t.Fatal(err)
	}
	if i.Err() != nil || i.MustIMSI() != "123456789012345" {
		t.Fatal(i.Err())
	}
}

func TestShortMarshalDestination(t *testing.T) {
	i := NewIMSI("123456789012345")
	for n := 0; n < i.MarshalLen(); n++ {
		if err := i.MarshalTo(make([]byte, n)); !errors.Is(err, ErrTooShortToMarshal) {
			t.Fatalf("length %d: %v", n, err)
		}
	}
	b, err := i.Marshal()
	if err != nil || len(b) != i.MarshalLen() {
		t.Fatalf("MarshalLen: %d %v", len(b), err)
	}
}

func TestIdentityFillerValidation(t *testing.T) {
	i := NewIMSI("123456789012345")
	i.Payload = []byte{0xf1, 0x32, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	if _, err := i.IMSI(); !errors.Is(err, utils.ErrMalformed) {
		t.Fatal(err)
	}
	i.Payload = []byte{0xfa, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	if _, err := i.IMSI(); !errors.Is(err, utils.ErrMalformed) {
		t.Fatal(err)
	}
	i.Payload = []byte{0x21, 0x43, 0x65, 0x87, 0x09, 0x21, 0x43, 0xff}
	if got, err := i.IMSI(); err != nil || got != "12345678901234" {
		t.Fatalf("even identity: %q %v", got, err)
	}
	m := NewMSISDN("123")
	if got, err := m.MSISDN(); err != nil || got != "123" {
		t.Fatalf("MSISDN filler: %q %v", got, err)
	}
	m = NewMSISDN("12a")
	if m == nil || !errors.Is(m.Err(), utils.ErrMalformed) {
		t.Fatal("invalid MSISDN lost")
	}
}

func TestSignatureOverflow(t *testing.T) {
	i := NewPTMSISignature(1 << 24)
	if i == nil || !errors.Is(i.Err(), utils.ErrOutOfRange) {
		t.Fatalf("%v", i)
	}
	if _, err := i.PTMSISignature(); !errors.Is(err, utils.ErrOutOfRange) {
		t.Fatal(err)
	}
}

func TestErrorTypes(t *testing.T) {
	cause := &InvalidTypeError{Type: 99}
	i := newErrorIE(1, cause)
	_, err := i.Marshal()
	var got *InvalidTypeError
	if !errors.As(err, &got) || got.Type != 99 {
		t.Fatalf("typed cause lost: %v", err)
	}
	i = New(PrivateExtension, make([]byte, 65536))
	if !errors.Is(i.Err(), ErrInvalidLength) || i.MarshalLen() != 0 {
		t.Fatalf("oversized IE: %v", i.Err())
	}
}

func TestPLMNCallSites(t *testing.T) {
	for _, mnc := range []string{"45", "456"} {
		i := NewRouteingAreaIdentity("123", mnc, 1, 2)
		if got, err := i.MCC(); err != nil || got != "123" {
			t.Fatalf("MCC %q %v", got, err)
		}
		if got, err := i.MNC(); err != nil || got != mnc {
			t.Fatalf("MNC %q %v", got, err)
		}
		i.Payload[0] = 0xaa
		if _, err := i.MCC(); !errors.Is(err, utils.ErrMalformed) {
			t.Fatal(err)
		}
	}
	i := NewRouteingAreaIdentity("12a", "45", 1, 2)
	if i == nil || !errors.Is(i.Err(), utils.ErrMalformed) {
		t.Fatal("invalid PLMN lost")
	}
}
