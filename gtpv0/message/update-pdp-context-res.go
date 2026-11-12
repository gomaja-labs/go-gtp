// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
)

// UpdatePDPContextResponse is a UpdatePDPContextResponse Header and its IEs above.
type UpdatePDPContextResponse struct {
	*Header
	Cause                     *ie.IE
	QoSProfile                *ie.IE
	Recovery                  *ie.IE
	FlowLabelDataI            *ie.IE
	FlowLabelSignalling       *ie.IE
	ChargingID                *ie.IE
	EndUserAddress            *ie.IE
	GGSNAddressForSignalling  *ie.IE
	GGSNAddressForUserTraffic *ie.IE
	ChargingGatewayAddress    *ie.IE
	PrivateExtensions         []*ie.IE
	AdditionalIEs             []*ie.IE
}

// NewUpdatePDPContextResponse creates a new UpdatePDPContextResponse.
func NewUpdatePDPContextResponse(seq, label uint16, tid uint64, ies ...*ie.IE) *UpdatePDPContextResponse {
	u := &UpdatePDPContextResponse{
		Header: NewHeader(
			0x1e, MsgTypeUpdatePDPContextResponse, seq, label, tid, nil,
		),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Cause:
			if u.Cause == nil {
				u.Cause = i
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
		case ie.ChargingID:
			if u.ChargingID == nil {
				u.ChargingID = i
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
			if u.GGSNAddressForSignalling == nil {
				u.GGSNAddressForSignalling = i
			} else {
				u.GGSNAddressForUserTraffic = i
			}
		case ie.ChargingGatewayAddress:
			if u.ChargingGatewayAddress == nil {
				u.ChargingGatewayAddress = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
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

// Marshal returns the byte sequence generated from a UpdatePDPContextResponse.
func (u *UpdatePDPContextResponse) Marshal() ([]byte, error) {
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
func (u *UpdatePDPContextResponse) MarshalTo(b []byte) error {
	if u == nil || u.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{u.Cause, u.QoSProfile, u.Recovery, u.FlowLabelDataI, u.FlowLabelSignalling, u.ChargingID, u.EndUserAddress, u.GGSNAddressForSignalling, u.GGSNAddressForUserTraffic, u.ChargingGatewayAddress}
	ies = append(ies, u.PrivateExtensions...)
	ies = append(ies, u.AdditionalIEs...)
	return marshalIEs(u.Header, b, ies)
}

// ParseUpdatePDPContextResponse parses a given byte sequence as a UpdatePDPContextResponse.
func ParseUpdatePDPContextResponse(b []byte) (*UpdatePDPContextResponse, error) {
	u := &UpdatePDPContextResponse{}
	if err := u.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return u, err
		}
		return nil, err
	}
	return u, nil
}

// UnmarshalBinary parses a given byte sequence as a UpdatePDPContextResponse.
func (u *UpdatePDPContextResponse) UnmarshalBinary(b []byte) error {
	if u == nil {
		return ErrInvalidLength
	}
	*u = UpdatePDPContextResponse{}
	var err error
	u.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if u.Header.Type != MsgTypeUpdatePDPContextResponse {
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
		case ie.Cause:
			if u.Cause == nil {
				u.Cause = i
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
		case ie.ChargingID:
			if u.ChargingID == nil {
				u.ChargingID = i
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
			if u.GGSNAddressForSignalling == nil {
				u.GGSNAddressForSignalling = i
			} else {
				u.GGSNAddressForUserTraffic = i
			}
		case ie.ChargingGatewayAddress:
			if u.ChargingGatewayAddress == nil {
				u.ChargingGatewayAddress = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
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
func (u *UpdatePDPContextResponse) MarshalLen() int {
	if u == nil || u.Header == nil {
		return 0
	}
	ies := []*ie.IE{u.Cause, u.QoSProfile, u.Recovery, u.FlowLabelDataI, u.FlowLabelSignalling, u.ChargingID, u.EndUserAddress, u.GGSNAddressForSignalling, u.GGSNAddressForUserTraffic, u.ChargingGatewayAddress}
	ies = append(ies, u.PrivateExtensions...)
	ies = append(ies, u.AdditionalIEs...)
	return messageLen(u.Header, ies)
}

// SetLength sets the length in Length field.
func (u *UpdatePDPContextResponse) SetLength() {
	if u == nil || u.Header == nil {
		return
	}
	u.Header.Length = uint16(u.MarshalLen() - 20)
}

// MessageTypeName returns the name of protocol.
func (u *UpdatePDPContextResponse) MessageTypeName() string {
	return "Update PDP Context Response"
}

// TID returns the TID in human-readable string.
func (u *UpdatePDPContextResponse) TID() string {
	if u == nil || u.Header == nil {
		return ""
	}
	return u.tid()
}
