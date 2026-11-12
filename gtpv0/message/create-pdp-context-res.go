// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
)

// CreatePDPContextResponse is a CreatePDPContextResponse Header and its IEs above.
type CreatePDPContextResponse struct {
	*Header
	Cause                     *ie.IE
	QoSProfile                *ie.IE
	ReorderingRequired        *ie.IE
	Recovery                  *ie.IE
	FlowLabelDataI            *ie.IE
	FlowLabelSignalling       *ie.IE
	ChargingID                *ie.IE
	EndUserAddress            *ie.IE
	PCO                       *ie.IE
	GGSNAddressForSignalling  *ie.IE
	GGSNAddressForUserTraffic *ie.IE
	ChargingGatewayAddress    *ie.IE
	PrivateExtensions         []*ie.IE
	AdditionalIEs             []*ie.IE
}

// NewCreatePDPContextResponse creates a new CreatePDPContextResponse.
func NewCreatePDPContextResponse(seq, label uint16, tid uint64, ies ...*ie.IE) *CreatePDPContextResponse {
	c := &CreatePDPContextResponse{
		Header: NewHeader(
			0x1e, MsgTypeCreatePDPContextResponse, seq, label, tid, nil,
		),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Cause:
			if c.Cause == nil {
				c.Cause = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.QualityOfServiceProfile:
			if c.QoSProfile == nil {
				c.QoSProfile = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ReorderingRequired:
			if c.ReorderingRequired == nil {
				c.ReorderingRequired = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.Recovery:
			if c.Recovery == nil {
				c.Recovery = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.FlowLabelDataI:
			if c.FlowLabelDataI == nil {
				c.FlowLabelDataI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.FlowLabelSignalling:
			if c.FlowLabelSignalling == nil {
				c.FlowLabelSignalling = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ChargingID:
			if c.ChargingID == nil {
				c.ChargingID = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.EndUserAddress:
			if c.EndUserAddress == nil {
				c.EndUserAddress = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ProtocolConfigurationOptions:
			if c.PCO == nil {
				c.PCO = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.GSNAddress:
			if c.GGSNAddressForSignalling == nil {
				c.GGSNAddressForSignalling = i
			} else {
				c.GGSNAddressForUserTraffic = i
			}
		case ie.ChargingGatewayAddress:
			if c.ChargingGatewayAddress == nil {
				c.ChargingGatewayAddress = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.PrivateExtension:
			c.PrivateExtensions = append(c.PrivateExtensions, i)
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	}

	c.SetLength()
	return c
}

// Marshal returns the byte sequence generated from a CreatePDPContextResponse.
func (c *CreatePDPContextResponse) Marshal() ([]byte, error) {
	if c == nil || c.Header == nil {
		return nil, ErrInvalidLength
	}
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (c *CreatePDPContextResponse) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{c.Cause, c.QoSProfile, c.ReorderingRequired, c.Recovery, c.FlowLabelDataI, c.FlowLabelSignalling, c.ChargingID, c.EndUserAddress, c.PCO, c.GGSNAddressForSignalling, c.GGSNAddressForUserTraffic, c.ChargingGatewayAddress}
	ies = append(ies, c.PrivateExtensions...)
	ies = append(ies, c.AdditionalIEs...)
	return marshalIEs(c.Header, b, ies)
}

// ParseCreatePDPContextResponse parses a given byte sequence as a CreatePDPContextResponse.
func ParseCreatePDPContextResponse(b []byte) (*CreatePDPContextResponse, error) {
	c := &CreatePDPContextResponse{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary parses a given byte sequence as a CreatePDPContextResponse.
func (c *CreatePDPContextResponse) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrInvalidLength
	}
	*c = CreatePDPContextResponse{}
	var err error
	c.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if c.Header.Type != MsgTypeCreatePDPContextResponse {
		return ErrInvalidMessageType
	}

	ies, err := ie.ParseMultiIEs(c.Header.Payload)
	if err != nil {
		return err
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.Cause:
			if c.Cause == nil {
				c.Cause = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.QualityOfServiceProfile:
			if c.QoSProfile == nil {
				c.QoSProfile = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ReorderingRequired:
			if c.ReorderingRequired == nil {
				c.ReorderingRequired = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.Recovery:
			if c.Recovery == nil {
				c.Recovery = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.FlowLabelDataI:
			if c.FlowLabelDataI == nil {
				c.FlowLabelDataI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.FlowLabelSignalling:
			if c.FlowLabelSignalling == nil {
				c.FlowLabelSignalling = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ChargingID:
			if c.ChargingID == nil {
				c.ChargingID = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.EndUserAddress:
			if c.EndUserAddress == nil {
				c.EndUserAddress = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ProtocolConfigurationOptions:
			if c.PCO == nil {
				c.PCO = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.GSNAddress:
			if c.GGSNAddressForSignalling == nil {
				c.GGSNAddressForSignalling = i
			} else {
				c.GGSNAddressForUserTraffic = i
			}
		case ie.ChargingGatewayAddress:
			if c.ChargingGatewayAddress == nil {
				c.ChargingGatewayAddress = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.PrivateExtension:
			c.PrivateExtensions = append(c.PrivateExtensions, i)
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	}

	return validateMandatory(c.Header.Type, c.Header.Payload)
}

// MarshalLen returns the serial length of Data.
func (c *CreatePDPContextResponse) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	ies := []*ie.IE{c.Cause, c.QoSProfile, c.ReorderingRequired, c.Recovery, c.FlowLabelDataI, c.FlowLabelSignalling, c.ChargingID, c.EndUserAddress, c.PCO, c.GGSNAddressForSignalling, c.GGSNAddressForUserTraffic, c.ChargingGatewayAddress}
	ies = append(ies, c.PrivateExtensions...)
	ies = append(ies, c.AdditionalIEs...)
	return messageLen(c.Header, ies)
}

// SetLength sets the length in Length field.
func (c *CreatePDPContextResponse) SetLength() {
	if c == nil || c.Header == nil {
		return
	}
	c.Header.Length = uint16(c.MarshalLen() - 20)
}

// MessageTypeName returns the name of protocol.
func (c *CreatePDPContextResponse) MessageTypeName() string {
	return "Create PDP Context Response"
}

// TID returns the TID in human-readable string.
func (c *CreatePDPContextResponse) TID() string {
	if c == nil || c.Header == nil {
		return ""
	}
	return c.tid()
}
