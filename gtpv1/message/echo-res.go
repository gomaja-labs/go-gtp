// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// EchoResponse is a EchoResponse Header and its IEs above.
type EchoResponse struct {
	*Header
	Recovery          *ie.IE
	PrivateExtensions []*ie.IE
	AdditionalIEs     []*ie.IE
}

// NewEchoResponse creates a new GTPv1 EchoResponse.
func NewEchoResponse(seq uint16, ies ...*ie.IE) *EchoResponse {
	e := &EchoResponse{
		Header: NewHeader(0x32, MsgTypeEchoResponse, 0, seq, nil),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Recovery:
			if e.Recovery == nil {
				e.Recovery = i
			} else {
				e.AdditionalIEs = append(e.AdditionalIEs, i)
			}
		case ie.PrivateExtension:
			e.PrivateExtensions = append(e.PrivateExtensions, i)
		default:
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	}

	e.SetLength()
	return e
}

// Marshal returns the byte sequence generated from a EchoResponse.
func (e *EchoResponse) Marshal() ([]byte, error) {
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
func (e *EchoResponse) MarshalTo(b []byte) error {
	if e == nil || e.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{e.Recovery}
	ies = append(ies, e.PrivateExtensions...)
	ies = append(ies, e.AdditionalIEs...)
	return marshalIEs(e.Header, b, ies)
}

// ParseEchoResponse decodes a given byte sequence as a EchoResponse.
func ParseEchoResponse(b []byte) (*EchoResponse, error) {
	e := &EchoResponse{}
	if err := e.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return e, err
		}
		return nil, err
	}
	return e, nil
}

// UnmarshalBinary decodes a given byte sequence as a EchoResponse.
func (e *EchoResponse) UnmarshalBinary(b []byte) error {
	if e == nil {
		return ErrInvalidLength
	}
	*e = EchoResponse{}
	var err error
	e.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if e.Header.Type != MsgTypeEchoResponse {
		return ErrInvalidMessageType
	}

	decodedIEs, err := ie.ParseMultiIEs(e.Header.Payload)
	if err != nil {
		return err
	}

	for _, i := range decodedIEs {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Recovery:
			if e.Recovery == nil {
				e.Recovery = i
			} else {
				e.AdditionalIEs = append(e.AdditionalIEs, i)
			}
		case ie.PrivateExtension:
			e.PrivateExtensions = append(e.PrivateExtensions, i)
		default:
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	}
	return validateMandatory(e.Header.Type, e.Header.Payload)
}

// MarshalLen returns the serial length of Data.
func (e *EchoResponse) MarshalLen() int {
	if e == nil || e.Header == nil {
		return 0
	}
	ies := []*ie.IE{e.Recovery}
	ies = append(ies, e.PrivateExtensions...)
	ies = append(ies, e.AdditionalIEs...)
	return messageLen(e.Header, ies)
}

// SetLength sets the length in Length field.
func (e *EchoResponse) SetLength() {
	if e == nil || e.Header == nil {
		return
	}
	e.Header.Length = uint16(e.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (e *EchoResponse) MessageTypeName() string {
	return "Echo Response"
}

// TEID returns the TEID in human-readable string.
func (e *EchoResponse) TEID() uint32 {
	if e == nil || e.Header == nil {
		return 0
	}
	return e.Header.TEID
}
