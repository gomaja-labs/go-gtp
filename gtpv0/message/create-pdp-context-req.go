// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
)

// CreatePDPContextRequest is a CreatePDPContextRequest Header and its IEs above.
type CreatePDPContextRequest struct {
	*Header
	RAI                       *ie.IE
	QoSProfile                *ie.IE
	Recovery                  *ie.IE
	SelectionMode             *ie.IE
	FlowLabelDataI            *ie.IE
	FlowLabelSignalling       *ie.IE
	EndUserAddress            *ie.IE
	APN                       *ie.IE
	PCO                       *ie.IE
	SGSNAddressForSignalling  *ie.IE
	SGSNAddressForUserTraffic *ie.IE
	MSISDN                    *ie.IE
	PrivateExtensions         []*ie.IE
	AdditionalIEs             []*ie.IE
}

// NewCreatePDPContextRequest creates a new CreatePDPContextRequest.
func NewCreatePDPContextRequest(seq, label uint16, tid uint64, ies ...*ie.IE) *CreatePDPContextRequest {
	c := &CreatePDPContextRequest{
		Header: NewHeader(
			0x1e, MsgTypeCreatePDPContextRequest, seq, label, tid, nil,
		),
	}

	// Optional IEs and Private Extensions, or any arbitrary IE.
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.RouteingAreaIdentity:
			if c.RAI == nil {
				c.RAI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.QualityOfServiceProfile:
			if c.QoSProfile == nil {
				c.QoSProfile = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.Recovery:
			if c.Recovery == nil {
				c.Recovery = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.SelectionMode:
			if c.SelectionMode == nil {
				c.SelectionMode = i
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
		case ie.EndUserAddress:
			if c.EndUserAddress == nil {
				c.EndUserAddress = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.AccessPointName:
			if c.APN == nil {
				c.APN = i
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
			if c.SGSNAddressForSignalling == nil {
				c.SGSNAddressForSignalling = i
			} else {
				c.SGSNAddressForUserTraffic = i
			}
		case ie.MSISDN:
			if c.MSISDN == nil {
				c.MSISDN = i
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

// Marshal returns the byte sequence generated from a CreatePDPContextRequest.
func (c *CreatePDPContextRequest) Marshal() ([]byte, error) {
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
func (c *CreatePDPContextRequest) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{c.RAI, c.QoSProfile, c.Recovery, c.SelectionMode, c.FlowLabelDataI, c.FlowLabelSignalling, c.EndUserAddress, c.APN, c.PCO, c.SGSNAddressForSignalling, c.SGSNAddressForUserTraffic, c.MSISDN}
	ies = append(ies, c.PrivateExtensions...)
	ies = append(ies, c.AdditionalIEs...)
	return marshalIEs(c.Header, b, ies)
}

// ParseCreatePDPContextRequest parses a given byte sequence as a CreatePDPContextRequest.
func ParseCreatePDPContextRequest(b []byte) (*CreatePDPContextRequest, error) {
	c := &CreatePDPContextRequest{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary parses a given byte sequence as a CreatePDPContextRequest.
func (c *CreatePDPContextRequest) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrInvalidLength
	}
	*c = CreatePDPContextRequest{}
	var err error
	c.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if c.Header.Type != MsgTypeCreatePDPContextRequest {
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
		case ie.RouteingAreaIdentity:
			if c.RAI == nil {
				c.RAI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.QualityOfServiceProfile:
			if c.QoSProfile == nil {
				c.QoSProfile = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.Recovery:
			if c.Recovery == nil {
				c.Recovery = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.SelectionMode:
			if c.SelectionMode == nil {
				c.SelectionMode = i
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
		case ie.EndUserAddress:
			if c.EndUserAddress == nil {
				c.EndUserAddress = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.AccessPointName:
			if c.APN == nil {
				c.APN = i
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
			if c.SGSNAddressForSignalling == nil {
				c.SGSNAddressForSignalling = i
			} else {
				c.SGSNAddressForUserTraffic = i
			}
		case ie.MSISDN:
			if c.MSISDN == nil {
				c.MSISDN = i
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
func (c *CreatePDPContextRequest) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	ies := []*ie.IE{c.RAI, c.QoSProfile, c.Recovery, c.SelectionMode, c.FlowLabelDataI, c.FlowLabelSignalling, c.EndUserAddress, c.APN, c.PCO, c.SGSNAddressForSignalling, c.SGSNAddressForUserTraffic, c.MSISDN}
	ies = append(ies, c.PrivateExtensions...)
	ies = append(ies, c.AdditionalIEs...)
	return messageLen(c.Header, ies)
}

// SetLength sets the length in Length field.
func (c *CreatePDPContextRequest) SetLength() {
	if c == nil || c.Header == nil {
		return
	}
	c.Header.Length = uint16(c.MarshalLen() - 20)
}

// MessageTypeName returns the name of protocol.
func (c *CreatePDPContextRequest) MessageTypeName() string {
	return "Create PDP Context Request"
}

// TID returns the TID in human-readable string.
func (c *CreatePDPContextRequest) TID() string {
	if c == nil || c.Header == nil {
		return ""
	}
	return c.tid()
}
