// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// EchoRequest is a EchoRequest Header and its IEs above.
type EchoRequest struct {
	*Header
	PrivateExtensions []*ie.IE
	AdditionalIEs     []*ie.IE
}

// NewEchoRequest creates a new GTP.
func NewEchoRequest(seq uint16, ies ...*ie.IE) *EchoRequest {
	e := &EchoRequest{
		Header: NewHeader(0x32, MsgTypeEchoRequest, 0, seq, nil),
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

// Marshal returns the byte sequence generated from a EchoRequest.
func (e *EchoRequest) Marshal() ([]byte, error) {
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
func (e *EchoRequest) MarshalTo(b []byte) error {
	if e == nil || e.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{}
	ies = append(ies, e.PrivateExtensions...)
	ies = append(ies, e.AdditionalIEs...)
	return marshalIEs(e.Header, b, ies)
}

// ParseEchoRequest decodes a given byte sequence as a EchoRequest.
func ParseEchoRequest(b []byte) (*EchoRequest, error) {
	e := &EchoRequest{}
	if err := e.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return e, err
		}
		return nil, err
	}
	return e, nil
}

// UnmarshalBinary decodes a given byte sequence as a EchoRequest.
func (e *EchoRequest) UnmarshalBinary(b []byte) error {
	if e == nil {
		return ErrInvalidLength
	}
	*e = EchoRequest{}
	var err error
	e.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if e.Header.Type != MsgTypeEchoRequest {
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
func (e *EchoRequest) MarshalLen() int {
	if e == nil || e.Header == nil {
		return 0
	}
	ies := []*ie.IE{}
	ies = append(ies, e.PrivateExtensions...)
	ies = append(ies, e.AdditionalIEs...)
	return messageLen(e.Header, ies)
}

// SetLength sets the length in Length field.
func (e *EchoRequest) SetLength() {
	if e == nil || e.Header == nil {
		return
	}
	e.Header.Length = uint16(e.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (e *EchoRequest) MessageTypeName() string {
	return "Echo Request"
}

// TEID returns the TEID in human-readable string.
func (e *EchoRequest) TEID() uint32 {
	if e == nil || e.Header == nil {
		return 0
	}
	return e.Header.TEID
}
