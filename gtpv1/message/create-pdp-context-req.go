// Copyright 2019-2023 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// CreatePDPContextRequest is a CreatePDPContextRequest Header and its IEs above.
type CreatePDPContextRequest struct {
	*Header
	IMSI                               *ie.IE
	RAI                                *ie.IE
	Recovery                           *ie.IE
	SelectionMode                      *ie.IE
	TEIDDataI                          *ie.IE
	TEIDCPlane                         *ie.IE
	NSAPI                              *ie.IE
	LinkedNSAPI                        *ie.IE
	ChargingCharacteristics            *ie.IE
	TraceReference                     *ie.IE
	TraceType                          *ie.IE
	EndUserAddress                     *ie.IE
	APN                                *ie.IE
	PCO                                *ie.IE
	SGSNAddressForSignalling           *ie.IE
	SGSNAddressForUserTraffic          *ie.IE
	MSISDN                             *ie.IE
	QoSProfile                         *ie.IE
	TFT                                *ie.IE
	TriggerID                          *ie.IE
	OMCIdentity                        *ie.IE
	CommonFlags                        *ie.IE
	APNRestriction                     *ie.IE
	RATType                            *ie.IE
	UserLocationInformation            *ie.IE
	MSTimeZone                         *ie.IE
	IMEI                               *ie.IE
	CAMELChargingInformationContainer  *ie.IE
	AdditionalTraceInfo                *ie.IE
	CorrelationID                      *ie.IE
	EvolvedARPI                        *ie.IE
	ExtendedCommonFlags                *ie.IE
	UCI                                *ie.IE
	APNAMBR                            *ie.IE
	SignallingPriorityIndication       *ie.IE
	CNOperatorSelectionEntity          *ie.IE
	MappedUEUsageType                  *ie.IE
	UPFunctionSelectionIndicationFlags *ie.IE
	PrivateExtensions                  []*ie.IE
	AdditionalIEs                      []*ie.IE
}

// NewCreatePDPContextRequest creates a new GTPv1 CreatePDPContextRequest.
func NewCreatePDPContextRequest(teid uint32, seq uint16, ies ...*ie.IE) *CreatePDPContextRequest {
	c := &CreatePDPContextRequest{
		Header: NewHeader(0x32, MsgTypeCreatePDPContextRequest, teid, seq, nil),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.IMSI:
			if c.IMSI == nil {
				c.IMSI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.RouteingAreaIdentity:
			if c.RAI == nil {
				c.RAI = i
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
				continue
			}
			if c.LinkedNSAPI == nil {
				c.LinkedNSAPI = i
				continue
			}
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		case ie.ChargingCharacteristics:
			if c.ChargingCharacteristics == nil {
				c.ChargingCharacteristics = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.TraceReference:
			if c.TraceReference == nil {
				c.TraceReference = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.TraceType:
			if c.TraceType == nil {
				c.TraceType = i
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
				continue
			}
			if c.SGSNAddressForUserTraffic == nil {
				c.SGSNAddressForUserTraffic = i
				continue
			}
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		case ie.MSISDN:
			if c.MSISDN == nil {
				c.MSISDN = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.QoSProfile:
			if c.QoSProfile == nil {
				c.QoSProfile = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.TrafficFlowTemplate:
			if c.TFT == nil {
				c.TFT = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.TriggerID:
			if c.TriggerID == nil {
				c.TriggerID = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.OMCIdentity:
			if c.OMCIdentity == nil {
				c.OMCIdentity = i
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
		case ie.RATType:
			if c.RATType == nil {
				c.RATType = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.UserLocationInformation:
			if c.UserLocationInformation == nil {
				c.UserLocationInformation = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.MSTimeZone:
			if c.MSTimeZone == nil {
				c.MSTimeZone = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.IMEISV:
			if c.IMEI == nil {
				c.IMEI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.CAMELChargingInformationContainer:
			if c.CAMELChargingInformationContainer == nil {
				c.CAMELChargingInformationContainer = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.AdditionalTraceInfo:
			if c.AdditionalTraceInfo == nil {
				c.AdditionalTraceInfo = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.CorrelationID:
			if c.CorrelationID == nil {
				c.CorrelationID = i
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
		case ie.UserCSGInformation:
			if c.UCI == nil {
				c.UCI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.AggregateMaximumBitRate:
			if c.APNAMBR == nil {
				c.APNAMBR = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.SignallingPriorityIndication:
			if c.SignallingPriorityIndication == nil {
				c.SignallingPriorityIndication = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.CNOperatorSelectionEntity:
			if c.CNOperatorSelectionEntity == nil {
				c.CNOperatorSelectionEntity = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.MappedUEUsageType:
			if c.MappedUEUsageType == nil {
				c.MappedUEUsageType = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.UPFunctionSelectionIndicationFlags:
			if c.UPFunctionSelectionIndicationFlags == nil {
				c.UPFunctionSelectionIndicationFlags = i
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
	ies := []*ie.IE{c.IMSI, c.RAI, c.Recovery, c.SelectionMode, c.TEIDDataI, c.TEIDCPlane, c.NSAPI, c.LinkedNSAPI, c.ChargingCharacteristics, c.TraceReference, c.TraceType, c.EndUserAddress, c.APN, c.PCO, c.SGSNAddressForSignalling, c.SGSNAddressForUserTraffic, c.MSISDN, c.QoSProfile, c.TFT, c.TriggerID, c.OMCIdentity, c.CommonFlags, c.APNRestriction, c.RATType, c.UserLocationInformation, c.MSTimeZone, c.IMEI, c.CAMELChargingInformationContainer, c.AdditionalTraceInfo, c.CorrelationID, c.EvolvedARPI, c.ExtendedCommonFlags, c.UCI, c.APNAMBR, c.SignallingPriorityIndication, c.CNOperatorSelectionEntity, c.MappedUEUsageType, c.UPFunctionSelectionIndicationFlags}
	ies = append(ies, c.PrivateExtensions...)
	ies = append(ies, c.AdditionalIEs...)
	return marshalIEs(c.Header, b, ies)
}

// ParseCreatePDPContextRequest decodes a given byte sequence as a CreatePDPContextRequest.
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

// UnmarshalBinary decodes a given byte sequence as a CreatePDPContextRequest.
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
		case ie.IMSI:
			if c.IMSI == nil {
				c.IMSI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.RouteingAreaIdentity:
			if c.RAI == nil {
				c.RAI = i
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
				continue
			}
			if c.LinkedNSAPI == nil {
				c.LinkedNSAPI = i
				continue
			}
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		case ie.ChargingCharacteristics:
			if c.ChargingCharacteristics == nil {
				c.ChargingCharacteristics = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.TraceReference:
			if c.TraceReference == nil {
				c.TraceReference = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.TraceType:
			if c.TraceType == nil {
				c.TraceType = i
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
				continue
			}
			if c.SGSNAddressForUserTraffic == nil {
				c.SGSNAddressForUserTraffic = i
				continue
			}
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		case ie.MSISDN:
			if c.MSISDN == nil {
				c.MSISDN = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.QoSProfile:
			if c.QoSProfile == nil {
				c.QoSProfile = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.TrafficFlowTemplate:
			if c.TFT == nil {
				c.TFT = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.TriggerID:
			if c.TriggerID == nil {
				c.TriggerID = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.OMCIdentity:
			if c.OMCIdentity == nil {
				c.OMCIdentity = i
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
		case ie.RATType:
			if c.RATType == nil {
				c.RATType = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.UserLocationInformation:
			if c.UserLocationInformation == nil {
				c.UserLocationInformation = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.MSTimeZone:
			if c.MSTimeZone == nil {
				c.MSTimeZone = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.IMEISV:
			if c.IMEI == nil {
				c.IMEI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.CAMELChargingInformationContainer:
			if c.CAMELChargingInformationContainer == nil {
				c.CAMELChargingInformationContainer = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.AdditionalTraceInfo:
			if c.AdditionalTraceInfo == nil {
				c.AdditionalTraceInfo = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.CorrelationID:
			if c.CorrelationID == nil {
				c.CorrelationID = i
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
		case ie.UserCSGInformation:
			if c.UCI == nil {
				c.UCI = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.AggregateMaximumBitRate:
			if c.APNAMBR == nil {
				c.APNAMBR = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.SignallingPriorityIndication:
			if c.SignallingPriorityIndication == nil {
				c.SignallingPriorityIndication = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.CNOperatorSelectionEntity:
			if c.CNOperatorSelectionEntity == nil {
				c.CNOperatorSelectionEntity = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.MappedUEUsageType:
			if c.MappedUEUsageType == nil {
				c.MappedUEUsageType = i
			} else {
				c.AdditionalIEs = append(c.AdditionalIEs, i)
			}
		case ie.UPFunctionSelectionIndicationFlags:
			if c.UPFunctionSelectionIndicationFlags == nil {
				c.UPFunctionSelectionIndicationFlags = i
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
	ies := []*ie.IE{c.IMSI, c.RAI, c.Recovery, c.SelectionMode, c.TEIDDataI, c.TEIDCPlane, c.NSAPI, c.LinkedNSAPI, c.ChargingCharacteristics, c.TraceReference, c.TraceType, c.EndUserAddress, c.APN, c.PCO, c.SGSNAddressForSignalling, c.SGSNAddressForUserTraffic, c.MSISDN, c.QoSProfile, c.TFT, c.TriggerID, c.OMCIdentity, c.CommonFlags, c.APNRestriction, c.RATType, c.UserLocationInformation, c.MSTimeZone, c.IMEI, c.CAMELChargingInformationContainer, c.AdditionalTraceInfo, c.CorrelationID, c.EvolvedARPI, c.ExtendedCommonFlags, c.UCI, c.APNAMBR, c.SignallingPriorityIndication, c.CNOperatorSelectionEntity, c.MappedUEUsageType, c.UPFunctionSelectionIndicationFlags}
	ies = append(ies, c.PrivateExtensions...)
	ies = append(ies, c.AdditionalIEs...)
	return messageLen(c.Header, ies)
}

// SetLength sets the length in Length field.
func (c *CreatePDPContextRequest) SetLength() {
	if c == nil || c.Header == nil {
		return
	}
	c.Length = uint16(c.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (c *CreatePDPContextRequest) MessageTypeName() string {
	return "Create PDP Context Request"
}

// TEID returns the TEID in human-readable string.
func (c *CreatePDPContextRequest) TEID() uint32 {
	if c == nil || c.Header == nil {
		return 0
	}
	return c.Header.TEID
}
