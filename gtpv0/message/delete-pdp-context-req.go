// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE fild.

package message

import (
	"fmt"

	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
)

// DeletePDPContextRequest is a DeletePDPContextRequest Header and its AdditionalIEs abovd.
type DeletePDPContextRequest struct {
	*Header
	PrivateExtensions []*ie.IE
	AdditionalIEs     []*ie.IE
}

// NewDeletePDPContextRequest creates a new DeletePDPContextRequest.
func NewDeletePDPContextRequest(seq, label uint16, tid uint64, ies ...*ie.IE) *DeletePDPContextRequest {
	d := &DeletePDPContextRequest{
		Header: NewHeader(
			0x1e, MsgTypeDeletePDPContextRequest, seq, label, tid, nil,
		),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.PrivateExtension:
			d.PrivateExtensions = append(d.PrivateExtensions, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	}

	d.SetLength()
	return d
}

// Marshal returns the byte sequence generated from a DeletePDPContextRequest.
func (d *DeletePDPContextRequest) Marshal() ([]byte, error) {
	if d == nil || d.Header == nil {
		return nil, ErrInvalidLength
	}
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (d *DeletePDPContextRequest) MarshalTo(b []byte) error {
	if d == nil || d.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{}
	ies = append(ies, d.PrivateExtensions...)
	ies = append(ies, d.AdditionalIEs...)
	return marshalIEs(d.Header, b, ies)
}

// ParseDeletePDPContextRequest parses a given byte sequence as a DeletePDPContextRequest.
func ParseDeletePDPContextRequest(b []byte) (*DeletePDPContextRequest, error) {
	d := &DeletePDPContextRequest{}
	if err := d.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return d, err
		}
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary parses a given byte sequence as a DeletePDPContextRequest.
func (d *DeletePDPContextRequest) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrInvalidLength
	}
	*d = DeletePDPContextRequest{}
	var err error
	d.Header, err = ParseHeader(b)
	if err != nil {
		return fmt.Errorf("failed to Parse Header: %w", err)
	}

	ies, err := ie.ParseMultiIEs(d.Header.Payload)
	if err != nil {
		return err
	}
	if d.Header.Type != MsgTypeDeletePDPContextRequest {
		return ErrInvalidMessageType
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.PrivateExtension:
			d.PrivateExtensions = append(d.PrivateExtensions, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	}

	return validateMandatory(d.Header.Type, d.Header.Payload)
}

// MarshalLen returns the serial length of Data.
func (d *DeletePDPContextRequest) MarshalLen() int {
	if d == nil || d.Header == nil {
		return 0
	}
	ies := []*ie.IE{}
	ies = append(ies, d.PrivateExtensions...)
	ies = append(ies, d.AdditionalIEs...)
	return messageLen(d.Header, ies)
}

// SetLength sets the length in Length field.
func (d *DeletePDPContextRequest) SetLength() {
	if d == nil || d.Header == nil {
		return
	}
	d.Header.Length = uint16(d.MarshalLen() - 20)
}

// MessageTypeName returns the name of protocol.
func (d *DeletePDPContextRequest) MessageTypeName() string {
	return "Delete PDP Context Request"
}

// TID returns the TID in human-readable string.
func (d *DeletePDPContextRequest) TID() string {
	if d == nil || d.Header == nil {
		return ""
	}
	return d.tid()
}
