// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"fmt"
	"io"
)

// NewRANNASCause creates a new RANNASCause IE.
func NewRANNASCause(pType, cType uint8, cause []byte) *IE {
	v := NewRANNASCauseFields(pType, cType, cause)
	b, err := v.Marshal()
	if err != nil {
		return newErrorIE(RANNASCause, err)
	}

	return New(RANNASCause, 0x00, b)
}

// RANNASCause returns RANNASCause in RANNASCauseFields type if the type of IE matches.
func (i *IE) RANNASCause() (*RANNASCauseFields, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	switch i.Type {
	case RANNASCause:
		return ParseRANNASCauseFields(i.Payload)
	case BearerContext:
		ies, err := i.BearerContext()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve RANNASCause: %w", err)
		}

		// TS 29.274 §7.7.9: match the child type and instance together.
		for _, child := range ies {
			if child != nil && child.Type == RANNASCause && child.Instance() == 0 {
				return child.RANNASCause()
			}
		}
		return nil, ErrIENotFound
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// RANNASCauseFields is a set of fields in RANNASCause IE.
type RANNASCauseFields struct {
	ProtocolType uint8  // 4-bit
	CauseType    uint8  // 4-bit
	Cause        []byte // format depends on ProtocolType
}

// NewRANNASCauseFields creates a new RANNASCauseFields.
func NewRANNASCauseFields(pType, cType uint8, cause []byte) *RANNASCauseFields {
	return &RANNASCauseFields{
		ProtocolType: pType,
		CauseType:    cType,
		Cause:        cause,
	}
}

// Marshal serializes RANNASCauseFields.
func (f *RANNASCauseFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo serializes RANNASCauseFields.
func (f *RANNASCauseFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	n := ranNASWidth(f.ProtocolType)
	// TS 29.274 §8.103: S1AP/EMM/ESM one octet, Diameter/IKEv2 two.
	if n == 0 || len(f.Cause) != n || f.CauseType > 15 || f.ProtocolType == 1 && f.CauseType > 4 {
		return ErrMalformed
	}
	if len(b) < f.MarshalLen() {
		return ErrTooShortToMarshal
	}
	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	b[0] = ((f.ProtocolType & 0x0f) << 4) | (f.CauseType & 0x0f)

	if l < 1+len(f.Cause) {
		return io.ErrUnexpectedEOF
	}
	copy(b[1:], f.Cause)

	return nil
}

// ParseRANNASCauseFields decodes RANNASCauseFields.
func ParseRANNASCauseFields(b []byte) (*RANNASCauseFields, error) {
	f := &RANNASCauseFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return f, nil
}

// UnmarshalBinary decodes given bytes into RANNASCauseFields.
func (f *RANNASCauseFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = RANNASCauseFields{}
	l := len(b)
	if l < 2 {
		return io.ErrUnexpectedEOF
	}

	f.ProtocolType = (b[0] & 0xf0) >> 4
	f.CauseType = b[0] & 0x0f
	n := ranNASWidth(f.ProtocolType)
	if n == 0 {
		return ErrMalformed
	}
	if len(b) < 1+n {
		return io.ErrUnexpectedEOF
	}
	f.Cause = b[1 : 1+n]

	return nil
}

// MarshalLen returns the serial length of RANNASCauseFields in int.
func (f *RANNASCauseFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 1 + len(f.Cause)
}

func ranNASWidth(proto uint8) int {
	switch proto {
	case 1, 2, 3:
		return 1
	case 4, 5:
		return 2
	default:
		return 0
	}
}
