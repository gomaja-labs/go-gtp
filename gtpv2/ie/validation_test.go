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
	i.Payload = []byte{0xf1, 0x32}
	if _, err := i.IMSI(); !errors.Is(err, utils.ErrMalformed) {
		t.Fatal(err)
	}
	i.Payload = []byte{0xfa}
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

func TestCommonValueAccessors(t *testing.T) {
	var missing *IE
	for _, i := range []*IE{missing, newErrorIE(1, ErrMalformed)} {
		want := ErrMalformed
		if i == nil {
			want = ErrIENotFound
		}
		_, err := i.ValueAsUint8()
		if !errors.Is(err, want) {
			t.Fatal(err)
		}
		_, err = i.ValueAsUint16()
		if !errors.Is(err, want) {
			t.Fatal(err)
		}
		_, err = i.ValueAsUint32()
		if !errors.Is(err, want) {
			t.Fatal(err)
		}
		_, err = i.ValueAsUint64()
		if !errors.Is(err, want) {
			t.Fatal(err)
		}
		_, err = i.ValueAsString()
		if !errors.Is(err, want) {
			t.Fatal(err)
		}
		_, err = i.ValueAsFQDN()
		if !errors.Is(err, want) {
			t.Fatal(err)
		}
	}
	missing.SetInstance(1)
	if missing.WithInstance(1) != nil || missing.Instance() != 0 {
		t.Fatal("nil instance methods")
	}
	i := NewFQDNIE(FullyQualifiedDomainName, "a..b").WithInstance(3)
	if _, err := i.ValueAsFQDN(); !errors.Is(err, utils.ErrInvalidLength) {
		t.Fatal(err)
	}
	i = NewFQDNIE(FullyQualifiedDomainName, "a.example.")
	if got, err := i.ValueAsFQDN(); err != nil || got != "a.example" {
		t.Fatalf("%q %v", got, err)
	}
	i.Payload = []byte{2, 'a'}
	if _, err := i.ValueAsFQDN(); !errors.Is(err, utils.ErrInvalidLength) {
		t.Fatal(err)
	}
}

func TestQoSOverflowAndRates(t *testing.T) {
	constructors := []func(uint64, uint64, uint64, uint64) *IE{
		func(a, b, c, d uint64) *IE { return NewBearerQoS(0, 0, 0, 1, a, b, c, d) },
		func(a, b, c, d uint64) *IE { return NewFlowQoS(1, a, b, c, d) },
	}
	for _, newIE := range constructors {
		for n := 0; n < 4; n++ {
			rates := [4]uint64{}
			rates[n] = 1 << 40
			i := newIE(rates[0], rates[1], rates[2], rates[3])
			if i == nil || !errors.Is(i.Err(), utils.ErrOutOfRange) {
				t.Fatalf("rate %d not rejected", n)
			}
			if _, err := i.Marshal(); !errors.Is(err, utils.ErrOutOfRange) {
				t.Fatal(err)
			}
		}
		i := newIE(0x0102030405, 2, 3, 4)
		rates := []func() (uint64, error){i.MBRForUplink, i.MBRForDownlink, i.GBRForUplink, i.GBRForDownlink}
		for n, accessor := range rates {
			want := []uint64{0x0102030405, 2, 3, 4}[n]
			if got, err := accessor(); err != nil || got != want {
				t.Fatalf("rate %d: %x %v", n, got, err)
			}
		}
	}
	i := NewTraceReference("123", "45", 1<<24)
	if i == nil || !errors.Is(i.Err(), utils.ErrOutOfRange) {
		t.Fatal("trace overflow")
	}
	c := NewTFTPFComponentFlowLabel(1 << 24)
	if _, err := c.Marshal(); !errors.Is(err, utils.ErrOutOfRange) {
		t.Fatal(err)
	}
	if _, err := c.FlowLabel(); !errors.Is(err, utils.ErrOutOfRange) {
		t.Fatal(err)
	}
}

func TestGUTIPLMN(t *testing.T) {
	i := NewGUTI("123", "456", 1, 2, 3)
	f, err := ParseGUTIFields(i.Payload)
	if err != nil || f.MCC != "123" || f.MNC != "456" {
		t.Fatalf("%v %v", f, err)
	}
}

func TestGroupedCommonMarshal(t *testing.T) {
	i := NewGroupedIE(BearerContext, NewIMSI("123456789012345"))
	for n := 0; n < i.MarshalLen(); n++ {
		if err := i.MarshalTo(make([]byte, n)); !errors.Is(err, ErrTooShortToMarshal) {
			t.Fatalf("length %d: %v", n, err)
		}
	}
	i.Payload = nil
	i.SetLength()
	b, err := i.Marshal()
	if err != nil || len(b) != i.MarshalLen() {
		t.Fatalf("%d %v", len(b), err)
	}
	i.ChildIEs = append(i.ChildIEs, NewIMSI("x"))
	i.SetLength()
	if _, err := i.Marshal(); !errors.Is(err, utils.ErrMalformed) {
		t.Fatal(err)
	}
	// Reusing a grouped IE for a scalar IE clears children and spare instance bits.
	if err := i.UnmarshalBinary([]byte{Recovery, 0, 1, 0xf2, 1}); err != nil {
		t.Fatal(err)
	}
	if len(i.ChildIEs) != 0 || i.Instance() != 2 {
		t.Fatal("retained grouped state")
	}
	b, err = i.Marshal()
	if err != nil || b[3] != 2 {
		t.Fatalf("spare bits: %x %v", b, err)
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
	i = New(PrivateExtension, 0, make([]byte, 65536))
	if !errors.Is(i.Err(), ErrInvalidLength) || i.MarshalLen() != 0 {
		t.Fatalf("oversized IE: %v", i.Err())
	}
}

func TestSignatureTruncatedPayload(t *testing.T) {
	for n := 0; n < 3; n++ {
		i := New(PTMSISignature, 0, make([]byte, n))
		if _, err := i.PTMSISignature(); !errors.Is(err, ErrTooShortToParse) {
			t.Fatalf("length %d: %v", n, err)
		}
	}
}
