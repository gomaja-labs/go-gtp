// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"fmt"

	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
)

// EchoRequest is a EchoRequest Header and its AdditionalIEs above.
type EchoRequest struct {
	*Header
	PrivateExtensions []*ie.IE
	AdditionalIEs     []*ie.IE
}

// NewEchoRequest creates a new EchoRequest.
func NewEchoRequest(seq, label uint16, tid uint64, ies ...*ie.IE) *EchoRequest {
	e := &EchoRequest{
		Header: NewHeader(
			0x1e, MsgTypeEchoRequest, seq, label, tid, nil,
		),
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

// ParseEchoRequest parses a given byte sequence as a EchoRequest.
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

// UnmarshalBinary parses a given byte sequence as a EchoRequest.
func (e *EchoRequest) UnmarshalBinary(b []byte) error {
	if e == nil {
		return ErrInvalidLength
	}
	*e = EchoRequest{}
	var err error
	e.Header, err = ParseHeader(b)
	if err != nil {
		return fmt.Errorf("failed to Parse Header: %w", err)
	}

	ies, err := ie.ParseMultiIEs(e.Header.Payload)
	if err != nil {
		return err
	}
	if e.Header.Type != MsgTypeEchoRequest {
		return ErrInvalidMessageType
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
	e.Header.Length = uint16(e.MarshalLen() - 20)
}

// MessageTypeName returns the name of protocol.
func (e *EchoRequest) MessageTypeName() string {
	return "Echo Request"
}

// TID returns the TID in human-readable string.
func (e *EchoRequest) TID() string {
	if e == nil || e.Header == nil {
		return ""
	}
	return e.tid()
}
