// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewIntegerNumber creates an Integer Number with an explicit width in octets.
// TS 29.274 §8.118: the containing field determines the width, not the value.
func NewIntegerNumber(num uint32, width int) *IE {
	if width < 1 || width > 4 {
		return newErrorIE(IntegerNumber, ErrInvalidLength)
	}
	if width < 4 && num >= 1<<uint(width*8) {
		return newErrorIE(IntegerNumber, ErrMalformed)
	}
	b := make([]byte, width)
	for n := width - 1; n >= 0; n-- {
		b[n] = byte(num)
		num >>= 8
	}
	return New(IntegerNumber, 0, b)
}

// IntegerNumber returns IntegerNumber in uint32 if the type of IE matches.
func (i *IE) IntegerNumber() (uint32, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	switch i.Type {
	case IntegerNumber:
		if len(i.Payload) < 1 || len(i.Payload) > 4 {
			return 0, ErrInvalidLength
		}
		var v uint32
		for _, b := range i.Payload {
			v = v<<8 | uint32(b)
		}
		return v, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// MustIntegerNumber returns IntegerNumber in uint32, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustIntegerNumber() uint32 {
	if i == nil {
		return 0
	}
	v, _ := i.IntegerNumber()
	return v
}
