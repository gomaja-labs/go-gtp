// Copyright 2019-2024 go-gtp authors. All rights reserveu.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// UpdatePDPContextResponse is a UpdatePDPContextResponse Header and its IEs above.
type UpdatePDPContextResponse struct {
	*Header
	Cause                         *ie.IE
	Recovery                      *ie.IE
	TEIDDataI                     *ie.IE
	TEIDCPlane                    *ie.IE
	ChargingID                    *ie.IE
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
	BearerControlMode             *ie.IE
	MSInfoChangeReportingAction   *ie.IE
	EvolvedARPI                   *ie.IE
	CSGInformationReportingAction *ie.IE
	APNAMBR                       *ie.IE
	PrivateExtensions             []*ie.IE
	AdditionalIEs                 []*ie.IE
}

// NewUpdatePDPContextResponse creates a new GTPv1 UpdatePDPContextResponse.
func NewUpdatePDPContextResponse(teid uint32, seq uint16, ies ...*ie.IE) *UpdatePDPContextResponse {
	u := &UpdatePDPContextResponse{
		Header: NewHeader(0x32, MsgTypeUpdatePDPContextResponse, teid, seq, nil),
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
		case ie.Recovery:
			if u.Recovery == nil {
				u.Recovery = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TEIDDataI:
			if u.TEIDDataI == nil {
				u.TEIDDataI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TEIDCPlane:
			if u.TEIDCPlane == nil {
				u.TEIDCPlane = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.ChargingID:
			if u.ChargingID == nil {
				u.ChargingID = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.ProtocolConfigurationOptions:
			if u.PCO == nil {
				u.PCO = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.GSNAddress:
			if u.GGSNAddressForCPlane == nil {
				u.GGSNAddressForCPlane = i
				continue
			}
			if u.GGSNAddressForUserTraffic == nil {
				u.GGSNAddressForUserTraffic = i
				continue
			}
			if u.AltGGSNAddressForCPlane == nil {
				u.AltGGSNAddressForCPlane = i
				continue
			}
			if u.AltGGSNAddressForUserTraffic == nil {
				u.AltGGSNAddressForUserTraffic = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.QoSProfile:
			if u.QoSProfile == nil {
				u.QoSProfile = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.ChargingGatewayAddress:
			if u.ChargingGatewayAddress == nil {
				u.ChargingGatewayAddress = i
			} else if u.AltChargingGatewayAddress == nil {
				u.AltChargingGatewayAddress = i
			}
		case ie.CommonFlags:
			if u.CommonFlags == nil {
				u.CommonFlags = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.APNRestriction:
			if u.APNRestriction == nil {
				u.APNRestriction = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.BearerControlMode:
			if u.BearerControlMode == nil {
				u.BearerControlMode = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.MSInfoChangeReportingAction:
			if u.MSInfoChangeReportingAction == nil {
				u.MSInfoChangeReportingAction = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.EvolvedAllocationRetentionPriorityI:
			if u.EvolvedARPI == nil {
				u.EvolvedARPI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.CSGInformationReportingAction:
			if u.CSGInformationReportingAction == nil {
				u.CSGInformationReportingAction = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.AggregateMaximumBitRate:
			if u.APNAMBR == nil {
				u.APNAMBR = i
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
	ies := []*ie.IE{u.Cause, u.Recovery, u.TEIDDataI, u.TEIDCPlane, u.ChargingID, u.PCO, u.GGSNAddressForCPlane, u.GGSNAddressForUserTraffic, u.AltGGSNAddressForCPlane, u.AltGGSNAddressForUserTraffic, u.QoSProfile, u.ChargingGatewayAddress, u.AltChargingGatewayAddress, u.CommonFlags, u.APNRestriction, u.BearerControlMode, u.MSInfoChangeReportingAction, u.EvolvedARPI, u.CSGInformationReportingAction, u.APNAMBR}
	ies = append(ies, u.PrivateExtensions...)
	ies = append(ies, u.AdditionalIEs...)
	return marshalIEs(u.Header, b, ies)
}

// ParseUpdatePDPContextResponse decodes a given byte sequence as a UpdatePDPContextResponse.
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

// UnmarshalBinary decodes a given byte sequence as a UpdatePDPContextResponse.
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
		case ie.Recovery:
			if u.Recovery == nil {
				u.Recovery = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TEIDDataI:
			if u.TEIDDataI == nil {
				u.TEIDDataI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TEIDCPlane:
			if u.TEIDCPlane == nil {
				u.TEIDCPlane = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.ChargingID:
			if u.ChargingID == nil {
				u.ChargingID = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.ProtocolConfigurationOptions:
			if u.PCO == nil {
				u.PCO = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.GSNAddress:
			if u.GGSNAddressForCPlane == nil {
				u.GGSNAddressForCPlane = i
				continue
			}
			if u.GGSNAddressForUserTraffic == nil {
				u.GGSNAddressForUserTraffic = i
				continue
			}
			if u.AltGGSNAddressForCPlane == nil {
				u.AltGGSNAddressForCPlane = i
				continue
			}
			if u.AltGGSNAddressForUserTraffic == nil {
				u.AltGGSNAddressForUserTraffic = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.QoSProfile:
			if u.QoSProfile == nil {
				u.QoSProfile = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.ChargingGatewayAddress:
			if u.ChargingGatewayAddress == nil {
				u.ChargingGatewayAddress = i
			} else if u.AltChargingGatewayAddress == nil {
				u.AltChargingGatewayAddress = i
			}
		case ie.CommonFlags:
			if u.CommonFlags == nil {
				u.CommonFlags = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.APNRestriction:
			if u.APNRestriction == nil {
				u.APNRestriction = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.BearerControlMode:
			if u.BearerControlMode == nil {
				u.BearerControlMode = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.MSInfoChangeReportingAction:
			if u.MSInfoChangeReportingAction == nil {
				u.MSInfoChangeReportingAction = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.EvolvedAllocationRetentionPriorityI:
			if u.EvolvedARPI == nil {
				u.EvolvedARPI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.CSGInformationReportingAction:
			if u.CSGInformationReportingAction == nil {
				u.CSGInformationReportingAction = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.AggregateMaximumBitRate:
			if u.APNAMBR == nil {
				u.APNAMBR = i
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
	ies := []*ie.IE{u.Cause, u.Recovery, u.TEIDDataI, u.TEIDCPlane, u.ChargingID, u.PCO, u.GGSNAddressForCPlane, u.GGSNAddressForUserTraffic, u.AltGGSNAddressForCPlane, u.AltGGSNAddressForUserTraffic, u.QoSProfile, u.ChargingGatewayAddress, u.AltChargingGatewayAddress, u.CommonFlags, u.APNRestriction, u.BearerControlMode, u.MSInfoChangeReportingAction, u.EvolvedARPI, u.CSGInformationReportingAction, u.APNAMBR}
	ies = append(ies, u.PrivateExtensions...)
	ies = append(ies, u.AdditionalIEs...)
	return messageLen(u.Header, ies)
}

// SetLength sets the length in Length field.
func (u *UpdatePDPContextResponse) SetLength() {
	if u == nil || u.Header == nil {
		return
	}
	u.Length = uint16(u.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (u *UpdatePDPContextResponse) MessageTypeName() string {
	return "Update PDP Context Response"
}

// TEID returns the TEID in human-readable string.
func (u *UpdatePDPContextResponse) TEID() uint32 {
	if u == nil || u.Header == nil {
		return 0
	}
	return u.Header.TEID
}
