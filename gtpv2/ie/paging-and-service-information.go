// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
)

// NewPagingAndServiceInformation creates a new PagingAndServiceInformation IE.
func NewPagingAndServiceInformation(ebi, flags, ppi uint8) *IE {
	// TS 29.274 §8.117: reject fields that do not fit.
	if ebi > 15 || flags > 1 || ppi > 127 {
		return newErrorIE(PagingAndServiceInformation, ErrMalformed)
	}
	v := NewPagingAndServiceInformationFields(ebi, flags, ppi)
	b, err := v.Marshal()
	if err != nil {
		return newErrorIE(PagingAndServiceInformation, err)
	}

	return New(PagingAndServiceInformation, 0x00, b)
}

// PagingAndServiceInformation returns PagingAndServiceInformation in PagingAndServiceInformationFields type if the type of IE matches.
func (i *IE) PagingAndServiceInformation() (*PagingAndServiceInformationFields, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	switch i.Type {
	case PagingAndServiceInformation:
		return ParsePagingAndServiceInformationFields(i.Payload)
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// PagingAndServiceInformationFields is a set of fields in PagingAndServiceInformation IE.
type PagingAndServiceInformationFields struct {
	err                    error
	EPSBearerID            uint8 // 4-bit
	Flags                  uint8 // 1st bit
	PagingPolicyIndication uint8 // 7-bit
}

// NewPagingAndServiceInformationFields creates a new PagingAndServiceInformationFields.
func NewPagingAndServiceInformationFields(ebi, flags, ppi uint8) *PagingAndServiceInformationFields {
	if ebi > 15 || flags > 1 || ppi > 127 {
		return &PagingAndServiceInformationFields{err: ErrMalformed}
	}
	return &PagingAndServiceInformationFields{
		EPSBearerID:            ebi & 0x0f,
		Flags:                  flags, // no mask, may be used in the future
		PagingPolicyIndication: ppi & 0x7f,
	}
}

// Marshal serializes PagingAndServiceInformationFields.
func (f *PagingAndServiceInformationFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo serializes PagingAndServiceInformationFields.
func (f *PagingAndServiceInformationFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.err != nil {
		return f.err
	}
	// TS 29.274 §8.117: validate exported fields before encoding.
	if f.EPSBearerID > 15 || f.Flags > 1 || f.PagingPolicyIndication > 127 {
		return ErrMalformed
	}
	if len(b) < f.MarshalLen() {
		return ErrTooShortToMarshal
	}
	l := len(b)
	if l < 2 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.EPSBearerID & 0x0f
	b[1] = f.Flags
	if has1stBit(f.Flags) {
		if l < 3 {
			return io.ErrUnexpectedEOF
		}
		b[2] = f.PagingPolicyIndication & 0x7f
	}

	return nil
}

// ParsePagingAndServiceInformationFields decodes PagingAndServiceInformationFields.
func ParsePagingAndServiceInformationFields(b []byte) (*PagingAndServiceInformationFields, error) {
	f := &PagingAndServiceInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return f, nil
}

// UnmarshalBinary decodes given bytes into PagingAndServiceInformationFields.
func (f *PagingAndServiceInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = PagingAndServiceInformationFields{}
	l := len(b)
	if l < 2 {
		return io.ErrUnexpectedEOF
	}

	f.EPSBearerID = b[0] & 15 & 0x0f
	f.Flags = b[1] & 1
	if has1stBit(f.Flags) {
		if l < 3 {
			return io.ErrUnexpectedEOF
		}
		f.PagingPolicyIndication = b[2] & 0x7f
	}

	return nil
}

// MarshalLen returns the serial length of PagingAndServiceInformationFields in int.
func (f *PagingAndServiceInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	l := 2
	if has1stBit(f.Flags) {
		l++
	}
	return l
}

// PagingPolicyIndication returns PagingPolicyIndication in uint8 if the type of IE matches.
func (i *IE) PagingPolicyIndication() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	switch i.Type {
	case PagingAndServiceInformation:
		f, err := ParsePagingAndServiceInformationFields(i.Payload)
		if err != nil {
			return 0, err
		}

		if !has1stBit(f.Flags) {
			return 0, ErrIEValueNotFound
		}

		return f.PagingPolicyIndication, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// MustPagingPolicyIndication returns PagingPolicyIndication in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustPagingPolicyIndication() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.PagingPolicyIndication()
	return v
}
