// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
	"time"
)

// NewThrottling creates a new Timer IE.
func NewThrottling(delay time.Duration, factor uint8) *IE {
	v := NewThrottlingFields(delay, factor)
	b, err := v.Marshal()
	if err != nil {
		return newErrorIE(Throttling, err)
	}

	return New(Throttling, 0x00, b)
}

// Throttling returns Throttling in time.Duration if the type of IE matches.
func (i *IE) Throttling() (*ThrottlingFields, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	switch i.Type {
	case Throttling:
		return ParseThrottlingFields(i.Payload)
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// ThrottlingFields is a set of fields in Throttling IE.
type ThrottlingFields struct {
	DelayUnit  uint8
	DelayValue time.Duration
	Factor     uint8
}

// NewThrottlingFields creates a new ThrottlingFields.
func NewThrottlingFields(delay time.Duration, factor uint8) *ThrottlingFields {
	raw, _ := encodeTimer(delay)
	return &ThrottlingFields{DelayUnit: raw >> 5, DelayValue: delay, Factor: factor}
}

// Marshal serializes ThrottlingFields.
func (f *ThrottlingFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo serializes ThrottlingFields.
func (f *ThrottlingFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if len(b) < f.MarshalLen() {
		return ErrTooShortToMarshal
	}
	l := len(b)
	if l < 2 {
		return io.ErrUnexpectedEOF
	}

	// TS 29.274 §8.85: reserved factors shall not be transmitted.
	if f.Factor > 100 {
		return ErrMalformed
	}
	raw, err := encodeTimer(f.DelayValue)
	if err != nil {
		return err
	}
	b[0] = raw
	b[1] = f.Factor

	return nil
}

// ParseThrottlingFields decodes ThrottlingFields.
func ParseThrottlingFields(b []byte) (*ThrottlingFields, error) {
	f := &ThrottlingFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return f, nil
}

// UnmarshalBinary decodes given bytes into ThrottlingFields.
func (f *ThrottlingFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = ThrottlingFields{}
	l := len(b)
	if l < 2 {
		return io.ErrUnexpectedEOF
	}

	f.DelayUnit = (b[0] & 0xe0) >> 5

	switch f.DelayUnit {
	case 0x04:
		f.DelayValue = time.Duration(b[0]&0x1f) * (10 * time.Hour)
	case 0x03:
		f.DelayValue = time.Duration(b[0]&0x1f) * time.Hour
	case 0x02:
		f.DelayValue = time.Duration(b[0]&0x1f) * (10 * time.Minute)
	case 0x01, 0x05, 0x06:
		f.DelayValue = time.Duration(b[0]&0x1f) * time.Minute
	case 0x00:
		f.DelayValue = time.Duration(b[0]&0x1f) * (2 * time.Second)
	case 0x07:
		f.DelayValue = time.Duration(1<<63 - 1)
	}

	f.Factor = b[1]
	if f.Factor > 100 {
		f.Factor = 0
	}

	return nil
}

// MarshalLen returns the serial length of ThrottlingFields in int.
func (f *ThrottlingFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 2
}
