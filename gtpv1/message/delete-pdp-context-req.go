// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// DeletePDPContextRequest is a DeletePDPContextRequest Header and its IEs above.
type DeletePDPContextRequest struct {
	*Header
	Cause               *ie.IE
	TeardownInd         *ie.IE
	NSAPI               *ie.IE
	PCO                 *ie.IE
	ULI                 *ie.IE
	MSTimeZone          *ie.IE
	ExtendedCommonFlags *ie.IE
	ULITimestamp        *ie.IE
	PrivateExtensions   []*ie.IE
	AdditionalIEs       []*ie.IE
}

// NewDeletePDPContextRequest creates a new GTPv1 DeletePDPContextRequest.
func NewDeletePDPContextRequest(teid uint32, seq uint16, ies ...*ie.IE) *DeletePDPContextRequest {
	d := &DeletePDPContextRequest{
		Header: NewHeader(0x32, MsgTypeDeletePDPContextRequest, teid, seq, nil),
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
		case ie.TeardownInd:
			if d.TeardownInd == nil {
				d.TeardownInd = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.NSAPI:
			if d.NSAPI == nil {
				d.NSAPI = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.ProtocolConfigurationOptions:
			if d.PCO == nil {
				d.PCO = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.UserLocationInformation:
			if d.ULI == nil {
				d.ULI = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.MSTimeZone:
			if d.MSTimeZone == nil {
				d.MSTimeZone = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.ExtendedCommonFlags:
			if d.ExtendedCommonFlags == nil {
				d.ExtendedCommonFlags = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.ULITimestamp:
			if d.ULITimestamp == nil {
				d.ULITimestamp = i
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
	ies := []*ie.IE{d.Cause, d.TeardownInd, d.NSAPI, d.PCO, d.ULI, d.MSTimeZone, d.ExtendedCommonFlags, d.ULITimestamp}
	ies = append(ies, d.PrivateExtensions...)
	ies = append(ies, d.AdditionalIEs...)
	return marshalIEs(d.Header, b, ies)
}

// ParseDeletePDPContextRequest decodes a given byte sequence as a DeletePDPContextRequest.
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

// UnmarshalBinary decodes a given byte sequence as a DeletePDPContextRequest.
func (d *DeletePDPContextRequest) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrInvalidLength
	}
	*d = DeletePDPContextRequest{}
	var err error
	d.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if d.Header.Type != MsgTypeDeletePDPContextRequest {
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
		case ie.TeardownInd:
			if d.TeardownInd == nil {
				d.TeardownInd = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.NSAPI:
			if d.NSAPI == nil {
				d.NSAPI = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.ProtocolConfigurationOptions:
			if d.PCO == nil {
				d.PCO = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.UserLocationInformation:
			if d.ULI == nil {
				d.ULI = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.MSTimeZone:
			if d.MSTimeZone == nil {
				d.MSTimeZone = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.ExtendedCommonFlags:
			if d.ExtendedCommonFlags == nil {
				d.ExtendedCommonFlags = i
			} else {
				d.AdditionalIEs = append(d.AdditionalIEs, i)
			}
		case ie.ULITimestamp:
			if d.ULITimestamp == nil {
				d.ULITimestamp = i
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
func (d *DeletePDPContextRequest) MarshalLen() int {
	if d == nil || d.Header == nil {
		return 0
	}
	ies := []*ie.IE{d.Cause, d.TeardownInd, d.NSAPI, d.PCO, d.ULI, d.MSTimeZone, d.ExtendedCommonFlags, d.ULITimestamp}
	ies = append(ies, d.PrivateExtensions...)
	ies = append(ies, d.AdditionalIEs...)
	return messageLen(d.Header, ies)
}

// SetLength sets the length in Length field.
func (d *DeletePDPContextRequest) SetLength() {
	if d == nil || d.Header == nil {
		return
	}
	d.Length = uint16(d.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (d *DeletePDPContextRequest) MessageTypeName() string {
	return "Delete PDP Context Request"
}

// TEID returns the TEID in human-readable string.
func (d *DeletePDPContextRequest) TEID() uint32 {
	if d == nil || d.Header == nil {
		return 0
	}
	return d.Header.TEID
}
