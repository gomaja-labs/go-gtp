// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
)

// UpdatePDPContextRequest is a UpdatePDPContextRequest Header and its IEs above.
type UpdatePDPContextRequest struct {
	*Header
	RAI                       *ie.IE
	QoSProfile                *ie.IE
	Recovery                  *ie.IE
	FlowLabelDataI            *ie.IE
	FlowLabelSignalling       *ie.IE
	EndUserAddress            *ie.IE
	SGSNAddressForSignalling  *ie.IE
	SGSNAddressForUserTraffic *ie.IE
	PrivateExtensions         []*ie.IE
	AdditionalIEs             []*ie.IE
}

// NewUpdatePDPContextRequest creates a new UpdatePDPContextRequest.
func NewUpdatePDPContextRequest(seq, label uint16, tid uint64, ies ...*ie.IE) *UpdatePDPContextRequest {
	u := &UpdatePDPContextRequest{
		Header: NewHeader(
			0x1e, MsgTypeUpdatePDPContextRequest, seq, label, tid, nil,
		),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.RouteingAreaIdentity:
			if u.RAI == nil {
				u.RAI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.QualityOfServiceProfile:
			if u.QoSProfile == nil {
				u.QoSProfile = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.Recovery:
			if u.Recovery == nil {
				u.Recovery = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.FlowLabelDataI:
			if u.FlowLabelDataI == nil {
				u.FlowLabelDataI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.FlowLabelSignalling:
			if u.FlowLabelSignalling == nil {
				u.FlowLabelSignalling = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.EndUserAddress:
			if u.EndUserAddress == nil {
				u.EndUserAddress = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.GSNAddress:
			if u.SGSNAddressForSignalling == nil {
				u.SGSNAddressForSignalling = i
			} else {
				u.SGSNAddressForUserTraffic = i
			}
		case ie.PrivateExtension:
			u.PrivateExtensions = append(u.PrivateExtensions, i)
		default:
			u.AdditionalIEs = append(u.AdditionalIEs, i)
		}
	}

	u.SetLength()
	return u
}

// Marshal returns the byte sequence generated from a UpdatePDPContextRequest.
func (u *UpdatePDPContextRequest) Marshal() ([]byte, error) {
	if u == nil || u.Header == nil {
		return nil, ErrInvalidLength
	}
	b := make([]byte, u.MarshalLen())
	if err := u.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (u *UpdatePDPContextRequest) MarshalTo(b []byte) error {
	if u == nil || u.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{u.RAI, u.QoSProfile, u.Recovery, u.FlowLabelDataI, u.FlowLabelSignalling, u.EndUserAddress, u.SGSNAddressForSignalling, u.SGSNAddressForUserTraffic}
	ies = append(ies, u.PrivateExtensions...)
	ies = append(ies, u.AdditionalIEs...)
	return marshalIEs(u.Header, b, ies)
}

// ParseUpdatePDPContextRequest parses a given byte sequence as a UpdatePDPContextRequest.
func ParseUpdatePDPContextRequest(b []byte) (*UpdatePDPContextRequest, error) {
	u := &UpdatePDPContextRequest{}
	if err := u.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return u, err
		}
		return nil, err
	}
	return u, nil
}

// UnmarshalBinary parses a given byte sequence as a UpdatePDPContextRequest.
func (u *UpdatePDPContextRequest) UnmarshalBinary(b []byte) error {
	if u == nil {
		return ErrInvalidLength
	}
	*u = UpdatePDPContextRequest{}
	var err error
	u.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if u.Header.Type != MsgTypeUpdatePDPContextRequest {
		return ErrInvalidMessageType
	}

	ies, err := ie.ParseMultiIEs(u.Header.Payload)
	if err != nil {
		return err
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.RouteingAreaIdentity:
			if u.RAI == nil {
				u.RAI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.QualityOfServiceProfile:
			if u.QoSProfile == nil {
				u.QoSProfile = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.Recovery:
			if u.Recovery == nil {
				u.Recovery = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.FlowLabelDataI:
			if u.FlowLabelDataI == nil {
				u.FlowLabelDataI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.FlowLabelSignalling:
			if u.FlowLabelSignalling == nil {
				u.FlowLabelSignalling = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.EndUserAddress:
			if u.EndUserAddress == nil {
				u.EndUserAddress = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.GSNAddress:
			if u.SGSNAddressForSignalling == nil {
				u.SGSNAddressForSignalling = i
			} else {
				u.SGSNAddressForUserTraffic = i
			}
		case ie.PrivateExtension:
			u.PrivateExtensions = append(u.PrivateExtensions, i)
		default:
			u.AdditionalIEs = append(u.AdditionalIEs, i)
		}
	}

	return validateMandatory(u.Header.Type, u.Header.Payload)
}

// MarshalLen returns the serial length of Data.
func (u *UpdatePDPContextRequest) MarshalLen() int {
	if u == nil || u.Header == nil {
		return 0
	}
	ies := []*ie.IE{u.RAI, u.QoSProfile, u.Recovery, u.FlowLabelDataI, u.FlowLabelSignalling, u.EndUserAddress, u.SGSNAddressForSignalling, u.SGSNAddressForUserTraffic}
	ies = append(ies, u.PrivateExtensions...)
	ies = append(ies, u.AdditionalIEs...)
	return messageLen(u.Header, ies)
}

// SetLength sets the length in Length field.
func (u *UpdatePDPContextRequest) SetLength() {
	if u == nil || u.Header == nil {
		return
	}
	u.Header.Length = uint16(u.MarshalLen() - 20)
}

// MessageTypeName returns the name of protocol.
func (u *UpdatePDPContextRequest) MessageTypeName() string {
	return "Update PDP Context Request"
}

// TID returns the TID in human-readable string.
func (u *UpdatePDPContextRequest) TID() string {
	if u == nil || u.Header == nil {
		return ""
	}
	return u.tid()
}
