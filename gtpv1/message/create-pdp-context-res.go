// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// CreatePDPContextResponse is a CreatePDPContextResponse Header and its IEs above.
type CreatePDPContextResponse struct {
	*Header
	Cause                         *ie.IE
	ReorderingRequired            *ie.IE
	Recovery                      *ie.IE
	TEIDDataI                     *ie.IE
	TEIDCPlane                    *ie.IE
	NSAPI                         *ie.IE
	ChargingID                    *ie.IE
	EndUserAddress                *ie.IE
	PCO                           *ie.IE
	GGSNAddressForCPlane          *ie.IE
	GGSNAddressForUserTraffic     *ie.IE
	AltGGSNAddressForCPlane       *ie.IE
	AltGGSNAddressForUserTraffic  *ie.IE
	QoSProfile                    *ie.IE
	ChargingGatewayAddress        *ie.IE
	AltChargingGatewayAddress     *ie.IE
	CommonFlags                   *ie.IE
	APNRestriction                *ie.IE
	MSInfoChangeReportingAction   *ie.IE
	BearerControlMode             *ie.IE
	EvolvedARPI                   *ie.IE
	ExtendedCommonFlags           *ie.IE
	CSGInformationReportingAction *ie.IE
	APNAMBR                       *ie.IE
	GGSNBackOffTime               *ie.IE
	ExtendedCommonFlagsII         *ie.IE
	PrivateExtensions             []*ie.IE
	AdditionalIEs                 []*ie.IE
}

// NewCreatePDPContextResponse creates a new GTPv1 CreatePDPContextResponse.
func NewCreatePDPContextResponse(teid uint32, seq uint16, ies ...*ie.IE) *CreatePDPContextResponse {
	c := &CreatePDPContextResponse{
		Header: NewHeader(0x32, MsgTypeCreatePDPContextResponse, teid, seq, nil),
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
		case ie.TEIDDataI:
			if c.TEIDDataI == nil {
				c.TEIDDataI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.TEIDCPlane:
			if c.TEIDCPlane == nil {
				c.TEIDCPlane = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.NSAPI:
			if c.NSAPI == nil {
				c.NSAPI = i
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
			if c.GGSNAddressForCPlane == nil {
				c.GGSNAddressForCPlane = i
				continue
			}
			if c.GGSNAddressForUserTraffic == nil {
				c.GGSNAddressForUserTraffic = i
				continue
			}
			if c.AltGGSNAddressForCPlane == nil {
				c.AltGGSNAddressForCPlane = i
				continue
			}
			if c.AltGGSNAddressForUserTraffic == nil {
				c.AltGGSNAddressForUserTraffic = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.QoSProfile:
			if c.QoSProfile == nil {
				c.QoSProfile = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ChargingGatewayAddress:
			if c.ChargingGatewayAddress == nil {
				c.ChargingGatewayAddress = i
			} else if c.AltChargingGatewayAddress == nil {
				c.AltChargingGatewayAddress = i
			}
		case ie.CommonFlags:
			if c.CommonFlags == nil {
				c.CommonFlags = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.APNRestriction:
			if c.APNRestriction == nil {
				c.APNRestriction = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.MSInfoChangeReportingAction:
			if c.MSInfoChangeReportingAction == nil {
				c.MSInfoChangeReportingAction = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.BearerControlMode:
			if c.BearerControlMode == nil {
				c.BearerControlMode = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.EvolvedAllocationRetentionPriorityI:
			if c.EvolvedARPI == nil {
				c.EvolvedARPI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ExtendedCommonFlags:
			if c.ExtendedCommonFlags == nil {
				c.ExtendedCommonFlags = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.CSGInformationReportingAction:
			if c.CSGInformationReportingAction == nil {
				c.CSGInformationReportingAction = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.AggregateMaximumBitRate:
			if c.APNAMBR == nil {
				c.APNAMBR = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.GGSNBackOffTime:
			if c.GGSNBackOffTime == nil {
				c.GGSNBackOffTime = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ExtendedCommonFlagsII:
			if c.ExtendedCommonFlagsII == nil {
				c.ExtendedCommonFlagsII = i
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
	ies := []*ie.IE{c.Cause, c.ReorderingRequired, c.Recovery, c.TEIDDataI, c.TEIDCPlane, c.NSAPI, c.ChargingID, c.EndUserAddress, c.PCO, c.GGSNAddressForCPlane, c.GGSNAddressForUserTraffic, c.AltGGSNAddressForCPlane, c.AltGGSNAddressForUserTraffic, c.QoSProfile, c.ChargingGatewayAddress, c.AltChargingGatewayAddress, c.CommonFlags, c.APNRestriction, c.MSInfoChangeReportingAction, c.BearerControlMode, c.EvolvedARPI, c.ExtendedCommonFlags, c.CSGInformationReportingAction, c.APNAMBR, c.GGSNBackOffTime, c.ExtendedCommonFlagsII}
	ies = append(ies, c.PrivateExtensions...)
	ies = append(ies, c.AdditionalIEs...)
	return marshalIEs(c.Header, b, ies)
}

// ParseCreatePDPContextResponse decodes a given byte sequence as a CreatePDPContextResponse.
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

// UnmarshalBinary decodes a given byte sequence as a CreatePDPContextResponse.
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
		case ie.TEIDDataI:
			if c.TEIDDataI == nil {
				c.TEIDDataI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.TEIDCPlane:
			if c.TEIDCPlane == nil {
				c.TEIDCPlane = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.NSAPI:
			if c.NSAPI == nil {
				c.NSAPI = i
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
			if c.GGSNAddressForCPlane == nil {
				c.GGSNAddressForCPlane = i
				continue
			}
			if c.GGSNAddressForUserTraffic == nil {
				c.GGSNAddressForUserTraffic = i
				continue
			}
			if c.AltGGSNAddressForCPlane == nil {
				c.AltGGSNAddressForCPlane = i
				continue
			}
			if c.AltGGSNAddressForUserTraffic == nil {
				c.AltGGSNAddressForUserTraffic = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.QoSProfile:
			if c.QoSProfile == nil {
				c.QoSProfile = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ChargingGatewayAddress:
			if c.ChargingGatewayAddress == nil {
				c.ChargingGatewayAddress = i
				continue
			}
			if c.AltChargingGatewayAddress == nil {
				c.AltChargingGatewayAddress = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.CommonFlags:
			if c.CommonFlags == nil {
				c.CommonFlags = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.APNRestriction:
			if c.APNRestriction == nil {
				c.APNRestriction = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.MSInfoChangeReportingAction:
			if c.MSInfoChangeReportingAction == nil {
				c.MSInfoChangeReportingAction = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.BearerControlMode:
			if c.BearerControlMode == nil {
				c.BearerControlMode = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.EvolvedAllocationRetentionPriorityI:
			if c.EvolvedARPI == nil {
				c.EvolvedARPI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ExtendedCommonFlags:
			if c.ExtendedCommonFlags == nil {
				c.ExtendedCommonFlags = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.CSGInformationReportingAction:
			if c.CSGInformationReportingAction == nil {
				c.CSGInformationReportingAction = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.AggregateMaximumBitRate:
			if c.APNAMBR == nil {
				c.APNAMBR = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.GGSNBackOffTime:
			if c.GGSNBackOffTime == nil {
				c.GGSNBackOffTime = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.ExtendedCommonFlagsII:
			if c.ExtendedCommonFlagsII == nil {
				c.ExtendedCommonFlagsII = i
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
	ies := []*ie.IE{c.Cause, c.ReorderingRequired, c.Recovery, c.TEIDDataI, c.TEIDCPlane, c.NSAPI, c.ChargingID, c.EndUserAddress, c.PCO, c.GGSNAddressForCPlane, c.GGSNAddressForUserTraffic, c.AltGGSNAddressForCPlane, c.AltGGSNAddressForUserTraffic, c.QoSProfile, c.ChargingGatewayAddress, c.AltChargingGatewayAddress, c.CommonFlags, c.APNRestriction, c.MSInfoChangeReportingAction, c.BearerControlMode, c.EvolvedARPI, c.ExtendedCommonFlags, c.CSGInformationReportingAction, c.APNAMBR, c.GGSNBackOffTime, c.ExtendedCommonFlagsII}
	ies = append(ies, c.PrivateExtensions...)
	ies = append(ies, c.AdditionalIEs...)
	return messageLen(c.Header, ies)
}

// SetLength sets the length in Length field.
func (c *CreatePDPContextResponse) SetLength() {
	if c == nil || c.Header == nil {
		return
	}
	c.Length = uint16(c.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (c *CreatePDPContextResponse) MessageTypeName() string {
	return "Create PDP Context Response"
}

// TEID returns the TEID in human-readable string.
func (c *CreatePDPContextResponse) TEID() uint32 {
	if c == nil || c.Header == nil {
		return 0
	}
	return c.Header.TEID
}
