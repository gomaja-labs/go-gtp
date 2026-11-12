// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// EndMarker is a EndMarker Header and its IEs above.
type EndMarker struct {
	*Header
	PrivateExtensions []*ie.IE
	AdditionalIEs     []*ie.IE
}

// NewEndMarker creates a new GTP.
func NewEndMarker(ies ...*ie.IE) *EndMarker {
	e := &EndMarker{
		Header: NewHeader(0x30, MsgTypeEndMarker, 0, 0, nil),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.PrivateExtension:
			e.PrivateExtensions = append(e.PrivateExtensions, i)
		default:
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	}

	e.SetLength()
	return e
}

// Marshal returns the byte sequence generated from a EndMarker.
func (e *EndMarker) Marshal() ([]byte, error) {
	if e == nil || e.Header == nil {
		return nil, ErrInvalidLength
	}
	b := make([]byte, e.MarshalLen())
	if err := e.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (e *EndMarker) MarshalTo(b []byte) error {
	if e == nil || e.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{}
	ies = append(ies, e.PrivateExtensions...)
	ies = append(ies, e.AdditionalIEs...)
	return marshalIEs(e.Header, b, ies)
}

// ParseEndMarker decodes a given byte sequence as a EndMarker.
func ParseEndMarker(b []byte) (*EndMarker, error) {
	e := &EndMarker{}
	if err := e.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return e, err
		}
		return nil, err
	}
	return e, nil
}

// UnmarshalBinary decodes a given byte sequence as a EndMarker.
func (e *EndMarker) UnmarshalBinary(b []byte) error {
	if e == nil {
		return ErrInvalidLength
	}
	*e = EndMarker{}
	var err error
	e.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if e.Header.Type != MsgTypeEndMarker {
		return ErrInvalidMessageType
	}

	ies, err := ie.ParseMultiIEs(e.Header.Payload)
	if err != nil {
		return err
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.PrivateExtension:
			e.PrivateExtensions = append(e.PrivateExtensions, i)
		default:
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	}

	return validateMandatory(e.Header.Type, e.Header.Payload)
}

// MarshalLen returns the serial length of Data.
func (e *EndMarker) MarshalLen() int {
	if e == nil || e.Header == nil {
		return 0
	}
	ies := []*ie.IE{}
	ies = append(ies, e.PrivateExtensions...)
	ies = append(ies, e.AdditionalIEs...)
	return messageLen(e.Header, ies)
}

// SetLength sets the length in Length field.
func (e *EndMarker) SetLength() {
	if e == nil || e.Header == nil {
		return
	}
	e.Header.Length = uint16(e.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (e *EndMarker) MessageTypeName() string {
	return "End Marker"
}

// TEID returns the TEID in human-readable string.
func (e *EndMarker) TEID() uint32 {
	if e == nil || e.Header == nil {
		return 0
	}
	return e.Header.TEID
}
