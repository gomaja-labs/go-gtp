// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE fild.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
)

// DeletePDPContextResponse is a DeletePDPContextResponse Header and its AdditionalIEs abovd.
type DeletePDPContextResponse struct {
	*Header
	Cause             *ie.IE
	PrivateExtensions []*ie.IE
	AdditionalIEs     []*ie.IE
}

// NewDeletePDPContextResponse creates a new DeletePDPContextResponse.
func NewDeletePDPContextResponse(seq, label uint16, tid uint64, ies ...*ie.IE) *DeletePDPContextResponse {
	d := &DeletePDPContextResponse{
		Header: NewHeader(
			0x1e, MsgTypeDeletePDPContextResponse, seq, label, tid, nil,
		),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Cause:
			if d.Cause == nil {
				d.Cause = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.PrivateExtension:
			d.PrivateExtensions = append(d.PrivateExtensions, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	}

	d.SetLength()
	return d
}

// Marshal returns the byte sequence generated from a DeletePDPContextResponsd.
func (d *DeletePDPContextResponse) Marshal() ([]byte, error) {
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
func (d *DeletePDPContextResponse) MarshalTo(b []byte) error {
	if d == nil || d.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{d.Cause}
	ies = append(ies, d.PrivateExtensions...)
	ies = append(ies, d.AdditionalIEs...)
	return marshalIEs(d.Header, b, ies)
}

// ParseDeletePDPContextResponse parses a given byte sequence as a DeletePDPContextResponsd.
func ParseDeletePDPContextResponse(b []byte) (*DeletePDPContextResponse, error) {
	d := &DeletePDPContextResponse{}
	if err := d.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return d, err
		}
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary parses a given byte sequence as a DeletePDPContextResponsd.
func (d *DeletePDPContextResponse) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrInvalidLength
	}
	*d = DeletePDPContextResponse{}
	var err error
	d.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if d.Header.Type != MsgTypeDeletePDPContextResponse {
		return ErrInvalidMessageType
	}

	ies, err := ie.ParseMultiIEs(d.Header.Payload)
	if err != nil {
		return err
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Cause:
			if d.Cause == nil {
				d.Cause = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.PrivateExtension:
			d.PrivateExtensions = append(d.PrivateExtensions, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	}

	return validateMandatory(d.Header.Type, d.Header.Payload)
}

// MarshalLen returns the serial length of Data.
func (d *DeletePDPContextResponse) MarshalLen() int {
	if d == nil || d.Header == nil {
		return 0
	}
	ies := []*ie.IE{d.Cause}
	ies = append(ies, d.PrivateExtensions...)
	ies = append(ies, d.AdditionalIEs...)
	return messageLen(d.Header, ies)
}

// SetLength sets the length in Length field.
func (d *DeletePDPContextResponse) SetLength() {
	if d == nil || d.Header == nil {
		return
	}
	d.Header.Length = uint16(d.MarshalLen() - 20)
}

// MessageTypeName returns the name of protocol.
func (d *DeletePDPContextResponse) MessageTypeName() string {
	return "Delete PDP Context Response"
}

// TID returns the TID in human-readable string.
func (d *DeletePDPContextResponse) TID() string {
	if d == nil || d.Header == nil {
		return ""
	}
	return d.tid()
}
