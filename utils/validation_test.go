// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package utils_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-gtp/utils"
)

func TestPLMNRejectsInvalidInput(t *testing.T) {
	for _, pair := range [][2]string{{"12", "34"}, {"1234", "56"}, {"12a", "34"}, {"123", "4"}, {"123", "4567"}} {
		if _, err := utils.EncodePLMN(pair[0], pair[1]); err == nil {
			t.Errorf("accepted %q", pair)
		}
	}
}

func TestPLMNTruncated(t *testing.T) {
	for n := 0; n < 3; n++ {
		if _, _, err := utils.DecodePLMN(make([]byte, n)); err == nil {
			t.Errorf("accepted %d bytes", n)
		}
	}
}

func TestECIRange(t *testing.T) {
	if _, _, err := utils.ParseECI(1 << 28); err == nil {
		t.Fatal("accepted 29-bit ECI")
	}
}

func TestTBCDDigitsAndFiller(t *testing.T) {
	for _, tc := range []struct {
		s string
		b []byte
	}{
		{"0123456789*#abc", []byte{0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe}},
		{"12", []byte{0x21}}, {"1", []byte{0xf1}}, {"", []byte{}},
	} {
		b, err := utils.EncodeTBCD(tc.s)
		if err != nil || !bytes.Equal(b, tc.b) {
			t.Fatalf("encode %q: %x %v", tc.s, b, err)
		}
		s, err := utils.DecodeTBCD(b)
		if err != nil || s != tc.s {
			t.Fatalf("decode %x: %q %v", b, s, err)
		}
	}
	for _, s := range []string{"f", "d", "A", "+1", "é"} {
		if _, err := utils.EncodeTBCD(s); !errors.Is(err, utils.ErrMalformed) {
			t.Errorf("accepted %q: %v", s, err)
		}
	}
	for _, b := range [][]byte{{0x1f}, {0xf1, 0x21}, {0xff, 0x00}} {
		if _, err := utils.DecodeTBCD(b); !errors.Is(err, utils.ErrMalformed) {
			t.Errorf("accepted %x: %v", b, err)
		}
	}
	if s, err := utils.DecodeTBCD([]byte{0x21, 0xf3, 0xff, 0xff}); err != nil || s != "123" {
		t.Fatalf("padding: %q %v", s, err)
	}
	for _, s := range []string{"*", "#", "a", "b", "c"} {
		if _, err := utils.EncodeBCDDigits(s); !errors.Is(err, utils.ErrMalformed) {
			t.Errorf("accepted identity %q", s)
		}
	}
	for n := byte(10); n < 15; n++ {
		if _, err := utils.DecodeBCDDigits([]byte{0xf0 | n}); !errors.Is(err, utils.ErrMalformed) {
			t.Errorf("accepted identity %x", n)
		}
	}
	b, err := utils.EncodeBCDDigits("0123456789")
	if err != nil {
		t.Fatal(err)
	}
	if s, err := utils.DecodeBCDDigits(b); err != nil || s != "0123456789" {
		t.Fatalf("digits %q %v", s, err)
	}
}

func TestPLMNInvalidNibbles(t *testing.T) {
	for pos := 0; pos < 6; pos++ {
		for digit := byte(10); digit < 16; digit++ {
			if pos == 3 && digit == 15 {
				continue
			}
			b := []byte{0x21, 0x63, 0x54}
			shift := uint(pos % 2 * 4)
			b[pos/2] = b[pos/2]&^(0xf<<shift) | digit<<shift
			if _, _, err := utils.DecodePLMN(b); !errors.Is(err, utils.ErrMalformed) {
				t.Errorf("accepted %x: %v", b, err)
			}
		}
	}
}

func TestFQDNStrict(t *testing.T) {
	// 254 octets without the root label, 255 with it: the RFC 1035 §3.1 maximum.
	longest := strings.Repeat("a", 63) + "." + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 61)
	for _, name := range []string{"a", "example.net", strings.Repeat("a", 63), longest, "a\x00"} {
		b, err := utils.EncodeFQDN(name + ".")
		if err != nil {
			t.Fatalf("encode %q: %v", name, err)
		}
		for _, wire := range [][]byte{b, append(append([]byte(nil), b...), 0)} {
			got, err := utils.DecodeFQDN(wire)
			if err != nil || got != name {
				t.Fatalf("decode %x: %q %v", wire, got, err)
			}
		}
	}
	for _, name := range []string{"", ".", "a..b", ".a", "a..", strings.Repeat("a", 64), longest + "e"} {
		if _, err := utils.EncodeFQDN(name); err == nil {
			t.Errorf("accepted name %q", name)
		}
	}
	for _, b := range [][]byte{nil, {0}, {2, 'a'}, {1, 'a', 0, 0}, {1, 'a', 0, 1, 'b'}, {1, '.'}, {64}, {0xc0, 0x00}} {
		if _, err := utils.DecodeFQDN(b); err == nil {
			t.Errorf("accepted encoding %x", b)
		}
	}
	b, err := utils.EncodeFQDN(longest)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 254 {
		t.Fatalf("longest name encodes to %d octets, want 254", len(b))
	}
	// labels encodes the given labels without a root label.
	labels := func(ls ...string) []byte {
		var out []byte
		for _, l := range ls {
			out = append(append(out, byte(len(l))), l...)
		}
		return out
	}
	a63, b63, c63 := strings.Repeat("a", 63), strings.Repeat("b", 63), strings.Repeat("c", 63)
	for _, oversized := range [][]byte{
		labels(a63, b63, c63, strings.Repeat("d", 62)),            // 255 octets, 256 with the root label
		append(labels(a63, b63, c63, strings.Repeat("d", 62)), 0), // the same with the root label
		append(append([]byte(nil), b...), 1, 'a'),                 // 256 octets
	} {
		if _, err := utils.DecodeFQDN(oversized); err == nil {
			t.Errorf("accepted oversized %d-octet name", len(oversized))
		}
	}
	if _, err := utils.EncodeFQDN(longest + "d"); err == nil {
		t.Error("accepted a name of 255 octets without the root label")
	}
}

func TestFixedWidthIntegers(t *testing.T) {
	if _, err := utils.Uint32To24(1 << 24); !errors.Is(err, utils.ErrOutOfRange) {
		t.Fatal(err)
	}
	if _, err := utils.Uint64To40(1 << 40); !errors.Is(err, utils.ErrOutOfRange) {
		t.Fatal(err)
	}
	if got := utils.Uint24To32([]byte{1, 2, 3, 4}); got != 0x010203 {
		t.Fatalf("%x", got)
	}
	if got := utils.Uint40To64([]byte{1, 2, 3, 4, 5, 6}); got != 0x0102030405 {
		t.Fatalf("%x", got)
	}
	for _, read := range []func(){func() { utils.Uint24To32([]byte{1, 2}) }, func() { utils.Uint40To64([]byte{1, 2, 3, 4}) }} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("short programmer-supplied slice did not panic")
				}
			}()
			read()
		}()
	}
	enb, cell, err := utils.ParseECI(0xfffffff)
	if err != nil || enb != 0xfffff || cell != 255 {
		t.Fatalf("%x %x %v", enb, cell, err)
	}
}

func FuzzDecodeHelpers(f *testing.F) {
	for _, b := range [][]byte{nil, {0}, {0xff}, {1, 'a'}, {1, 'a', 0}, {0x21, 0xf3, 0x54}, {0x1f}} {
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		_, _ = utils.DecodeTBCD(b)
		_, _ = utils.DecodeBCDDigits(b)
		_, _, _ = utils.DecodePLMN(b)
		_, _ = utils.DecodeFQDN(b)
	})
}
