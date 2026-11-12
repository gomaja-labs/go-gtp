// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package utils provides encoding helpers for GTP and other telecom protocols.
package utils

import (
	"errors"
	"fmt"
	"strings"
)

// Errors distinguish truncated inputs, invalid sizes and malformed values.
var (
	ErrInvalidLength = errors.New("invalid length")
	ErrMalformed     = errors.New("malformed value")
	ErrOutOfRange    = errors.New("value out of range")
)

const tbcdDigits = "0123456789*#abc"

// EncodeTBCD encodes a TBCD-STRING, with the first digit in the low nibble.
func EncodeTBCD(s string) ([]byte, error) {
	// TS 24.008 §10.5.4.7: digit codes and trailing high-nibble end mark.
	b := make([]byte, (len(s)+1)/2)
	for n := range b {
		b[n] = 0xff
	}
	for n := range len(s) {
		digit := strings.IndexByte(tbcdDigits, s[n])
		if digit < 0 {
			return nil, fmt.Errorf("TBCD digit at %d: %w", n, ErrMalformed)
		}
		if n%2 == 0 {
			b[n/2] = b[n/2]&0xf0 | byte(digit)
		} else {
			b[n/2] = b[n/2]&0x0f | byte(digit)<<4
		}
	}
	return b, nil
}

// DecodeTBCD decodes a TBCD-STRING. Filler is allowed only in a trailing run,
// including multiple filler octets used to pad fixed-width fields.
func DecodeTBCD(b []byte) (string, error) {
	var out strings.Builder
	filler := false
	for _, octet := range b {
		for _, digit := range []byte{octet & 0x0f, octet >> 4} {
			if digit == 0xf {
				filler = true
				continue
			}
			if filler {
				return "", fmt.Errorf("digit after TBCD filler: %w", ErrMalformed)
			}
			out.WriteByte(tbcdDigits[digit])
		}
	}
	return out.String(), nil
}

// EncodeBCDDigits encodes decimal identity digits with TBCD nibble ordering.
func EncodeBCDDigits(s string) ([]byte, error) {
	for n := range len(s) {
		if s[n] < '0' || s[n] > '9' {
			return nil, fmt.Errorf("BCD digit at %d: %w", n, ErrMalformed)
		}
	}
	return EncodeTBCD(s)
}

// DecodeBCDDigits decodes decimal identity digits, rejecting non-decimal codes.
func DecodeBCDDigits(b []byte) (string, error) {
	s, err := DecodeTBCD(b)
	if err != nil {
		return "", err
	}
	for n := range len(s) {
		if s[n] < '0' || s[n] > '9' {
			return "", fmt.Errorf("non-decimal BCD digit: %w", ErrMalformed)
		}
	}
	return s, nil
}

// Uint24To32 reads the first three octets in big-endian order, like encoding/binary.
// It panics if b is shorter than three octets; callers must check input lengths.
func Uint24To32(b []byte) uint32 {
	_ = b[2]
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// Uint32To24 encodes a big-endian 24-bit integer, rejecting overflow.
func Uint32To24(n uint32) ([]byte, error) {
	if n > 0xffffff {
		return nil, ErrOutOfRange
	}
	return []byte{byte(n >> 16), byte(n >> 8), byte(n)}, nil
}

// Uint40To64 reads the first five octets in big-endian order, like encoding/binary.
// It panics if b is shorter than five octets; callers must check input lengths.
func Uint40To64(b []byte) uint64 {
	_ = b[4]
	return uint64(b[0])<<32 | uint64(b[1])<<24 | uint64(b[2])<<16 | uint64(b[3])<<8 | uint64(b[4])
}

// Uint64To40 encodes a big-endian 40-bit integer, rejecting overflow.
func Uint64To40(n uint64) ([]byte, error) {
	if n > 0xffffffffff {
		return nil, ErrOutOfRange
	}
	return []byte{byte(n >> 32), byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}, nil
}

// EncodePLMN encodes a three-digit MCC and a two- or three-digit MNC.
func EncodePLMN(mcc, mnc string) ([]byte, error) {
	// TS 24.008 §10.5.1.13 and TS 23.003 §2.2: three octets, decimal digits.
	if len(mcc) != 3 || (len(mnc) != 2 && len(mnc) != 3) {
		return nil, ErrInvalidLength
	}
	c, err := EncodeBCDDigits(mcc)
	if err != nil {
		return nil, err
	}
	n, err := EncodeBCDDigits(mnc)
	if err != nil {
		return nil, err
	}
	third := byte(0xf)
	if len(mnc) == 3 {
		third = n[1] & 0xf
	}
	return []byte{c[0], c[1]&0xf | third<<4, n[0]}, nil
}

// DecodePLMN decodes the first three octets of a PLMN, preserving MNC width.
func DecodePLMN(b []byte) (mcc, mnc string, err error) {
	if len(b) < 3 {
		return "", "", ErrInvalidLength
	}
	digits := []byte{b[0] & 0xf, b[0] >> 4, b[1] & 0xf, b[2] & 0xf, b[2] >> 4}
	for _, d := range digits {
		if d > 9 {
			return "", "", ErrMalformed
		}
	}
	for n := range digits {
		digits[n] += '0'
	}
	mcc, mnc = string(digits[:3]), string(digits[3:])
	if d := b[1] >> 4; d != 0xf {
		if d > 9 {
			return "", "", ErrMalformed
		}
		mnc += string('0' + d)
	}
	return mcc, mnc, nil
}

// ParseECI splits a 28-bit ECI into a macro eNodeB ID and Cell ID.
func ParseECI(eci uint32) (enbID uint32, cellID uint8, err error) {
	// TS 23.003 §19.6: the E-UTRAN Cell Identity has 28 bits.
	if eci > 0xfffffff {
		return 0, 0, ErrOutOfRange
	}
	return eci >> 8, uint8(eci), nil
}

// EncodeFQDN encodes labels without a root label. A single trailing dot is accepted.
func EncodeFQDN(name string) ([]byte, error) {
	// RFC 1035 §§2.3.4, 3.1; RFC 2181 §11: labels are 1..63 octets and the
	// whole name, including the root label's zero length octet, is at most
	// 255 octets, so the encoding without the root label is at most 254.
	name = strings.TrimSuffix(name, ".")
	if len(name)+2 > 255 {
		return nil, ErrInvalidLength
	}
	b := make([]byte, 0, len(name)+1)
	for label := range strings.SplitSeq(name, ".") {
		if len(label) == 0 || len(label) > 63 {
			return nil, ErrInvalidLength
		}
		b = append(b, byte(len(label)))
		b = append(b, label...)
	}
	return b, nil
}

// DecodeFQDN decodes labels to dotted form, tolerating a single terminal root label.
func DecodeFQDN(b []byte) (string, error) {
	// RFC 1035 §§2.3.4, 3.1: no compression in this label-only representation.
	if len(b) == 0 {
		return "", ErrInvalidLength
	}
	if len(b) > 255 {
		return "", ErrInvalidLength
	}
	var labels []string
	for offset := 0; offset < len(b); {
		l := int(b[offset])
		offset++
		if l == 0 && offset == len(b) && offset > 1 {
			return strings.Join(labels, "."), nil
		}
		// RFC 1035 §3.1: at most 255 octets including the root label.
		if offset+l > 254 {
			return "", ErrInvalidLength
		}
		if l == 0 || l > 63 {
			return "", ErrMalformed
		}
		if l > len(b)-offset {
			return "", ErrInvalidLength
		}
		label := string(b[offset : offset+l])
		if strings.ContainsRune(label, '.') {
			return "", ErrMalformed
		}
		labels = append(labels, label)
		offset += l
	}
	return strings.Join(labels, "."), nil
}
