// Copyright 2019-2024 go-gtp authors. All rights reserveu.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// UpdatePDPContextRequest is a UpdatePDPContextRequest Header and its IEs above.
type UpdatePDPContextRequest struct {
	*Header
	IMSI                                 *ie.IE
	RAI                                  *ie.IE
	Recovery                             *ie.IE
	TEIDDataI                            *ie.IE
	TEIDCPlane                           *ie.IE
	NSAPI                                *ie.IE
	TraceReference                       *ie.IE
	TraceType                            *ie.IE
	PCO                                  *ie.IE
	SGSNAddressForCPlane                 *ie.IE
	SGSNAddressForUserTraffic            *ie.IE
	AlternativeSGSNAddressForCPlane      *ie.IE
	AlternativeSGSNAddressForUserTraffic *ie.IE
	QoSProfile                           *ie.IE
	TFT                                  *ie.IE
	TriggerID                            *ie.IE
	OMCIdentity                          *ie.IE
	CommonFlags                          *ie.IE
	RATType                              *ie.IE
	ULI                                  *ie.IE
	MSTimeZone                           *ie.IE
	AdditionalTraceInfo                  *ie.IE
	DirectTunnelFlags                    *ie.IE
	EvolvedARPI                          *ie.IE
	ExtendedCommonFlags                  *ie.IE
	UCI                                  *ie.IE
	APNAMBR                              *ie.IE
	SignallingPriorityIndication         *ie.IE
	CNOperatorSelectionEntity            *ie.IE
	IMEI                                 *ie.IE
	PrivateExtensions                    []*ie.IE
	AdditionalIEs                        []*ie.IE
}

// NewUpdatePDPContextRequest creates a new GTPv1 UpdatePDPContextRequest.
func NewUpdatePDPContextRequest(teid uint32, seq uint16, ies ...*ie.IE) *UpdatePDPContextRequest {
	u := &UpdatePDPContextRequest{
		Header: NewHeader(0x32, MsgTypeUpdatePDPContextRequest, teid, seq, nil),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.IMSI:
			if u.IMSI == nil {
				u.IMSI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.RouteingAreaIdentity:
			if u.RAI == nil {
				u.RAI = i
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
		case ie.NSAPI:
			if u.NSAPI == nil {
				u.NSAPI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TraceReference:
			if u.TraceReference == nil {
				u.TraceReference = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TraceType:
			if u.TraceType == nil {
				u.TraceType = i
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
			if u.SGSNAddressForCPlane == nil {
				u.SGSNAddressForCPlane = i
				continue
			}
			if u.SGSNAddressForUserTraffic == nil {
				u.SGSNAddressForUserTraffic = i
				continue
			}
			if u.AlternativeSGSNAddressForCPlane == nil {
				u.AlternativeSGSNAddressForCPlane = i
				continue
			}
			if u.AlternativeSGSNAddressForUserTraffic == nil {
				u.AlternativeSGSNAddressForUserTraffic = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.QoSProfile:
			if u.QoSProfile == nil {
				u.QoSProfile = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TrafficFlowTemplate:
			if u.TFT == nil {
				u.TFT = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TriggerID:
			if u.TriggerID == nil {
				u.TriggerID = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.OMCIdentity:
			if u.OMCIdentity == nil {
				u.OMCIdentity = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.CommonFlags:
			if u.CommonFlags == nil {
				u.CommonFlags = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.RATType:
			if u.RATType == nil {
				u.RATType = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.UserLocationInformation:
			if u.ULI == nil {
				u.ULI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.MSTimeZone:
			if u.MSTimeZone == nil {
				u.MSTimeZone = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.AdditionalTraceInfo:
			if u.AdditionalTraceInfo == nil {
				u.AdditionalTraceInfo = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.DirectTunnelFlags:
			if u.DirectTunnelFlags == nil {
				u.DirectTunnelFlags = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.EvolvedAllocationRetentionPriorityI:
			if u.EvolvedARPI == nil {
				u.EvolvedARPI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.ExtendedCommonFlags:
			if u.ExtendedCommonFlags == nil {
				u.ExtendedCommonFlags = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.UserCSGInformation:
			if u.UCI == nil {
				u.UCI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.AggregateMaximumBitRate:
			if u.APNAMBR == nil {
				u.APNAMBR = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.SignallingPriorityIndication:
			if u.SignallingPriorityIndication == nil {
				u.SignallingPriorityIndication = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.CNOperatorSelectionEntity:
			if u.CNOperatorSelectionEntity == nil {
				u.CNOperatorSelectionEntity = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.IMEISV:
			if u.IMEI == nil {
				u.IMEI = i
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
	ies := []*ie.IE{u.IMSI, u.RAI, u.Recovery, u.TEIDDataI, u.TEIDCPlane, u.NSAPI, u.TraceReference, u.TraceType, u.PCO, u.SGSNAddressForCPlane, u.SGSNAddressForUserTraffic, u.AlternativeSGSNAddressForCPlane, u.AlternativeSGSNAddressForUserTraffic, u.QoSProfile, u.TFT, u.TriggerID, u.OMCIdentity, u.CommonFlags, u.RATType, u.ULI, u.MSTimeZone, u.AdditionalTraceInfo, u.DirectTunnelFlags, u.EvolvedARPI, u.ExtendedCommonFlags, u.UCI, u.APNAMBR, u.SignallingPriorityIndication, u.CNOperatorSelectionEntity, u.IMEI}
	ies = append(ies, u.PrivateExtensions...)
	ies = append(ies, u.AdditionalIEs...)
	return marshalIEs(u.Header, b, ies)
}

// ParseUpdatePDPContextRequest decodes a given byte sequence as a UpdatePDPContextRequest.
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

// UnmarshalBinary decodes a given byte sequence as a UpdatePDPContextRequest.
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
		case ie.IMSI:
			if u.IMSI == nil {
				u.IMSI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.RouteingAreaIdentity:
			if u.RAI == nil {
				u.RAI = i
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
		case ie.NSAPI:
			if u.NSAPI == nil {
				u.NSAPI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TraceReference:
			if u.TraceReference == nil {
				u.TraceReference = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TraceType:
			if u.TraceType == nil {
				u.TraceType = i
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
			if u.SGSNAddressForCPlane == nil {
				u.SGSNAddressForCPlane = i
				continue
			}
			if u.SGSNAddressForUserTraffic == nil {
				u.SGSNAddressForUserTraffic = i
				continue
			}
			if u.AlternativeSGSNAddressForCPlane == nil {
				u.AlternativeSGSNAddressForCPlane = i
				continue
			}
			if u.AlternativeSGSNAddressForUserTraffic == nil {
				u.AlternativeSGSNAddressForUserTraffic = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.QoSProfile:
			if u.QoSProfile == nil {
				u.QoSProfile = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TrafficFlowTemplate:
			if u.TFT == nil {
				u.TFT = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.TriggerID:
			if u.TriggerID == nil {
				u.TriggerID = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.OMCIdentity:
			if u.OMCIdentity == nil {
				u.OMCIdentity = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.CommonFlags:
			if u.CommonFlags == nil {
				u.CommonFlags = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.RATType:
			if u.RATType == nil {
				u.RATType = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.UserLocationInformation:
			if u.ULI == nil {
				u.ULI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.MSTimeZone:
			if u.MSTimeZone == nil {
				u.MSTimeZone = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.AdditionalTraceInfo:
			if u.AdditionalTraceInfo == nil {
				u.AdditionalTraceInfo = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.DirectTunnelFlags:
			if u.DirectTunnelFlags == nil {
				u.DirectTunnelFlags = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.EvolvedAllocationRetentionPriorityI:
			if u.EvolvedARPI == nil {
				u.EvolvedARPI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.ExtendedCommonFlags:
			if u.ExtendedCommonFlags == nil {
				u.ExtendedCommonFlags = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.UserCSGInformation:
			if u.UCI == nil {
				u.UCI = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.AggregateMaximumBitRate:
			if u.APNAMBR == nil {
				u.APNAMBR = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.SignallingPriorityIndication:
			if u.SignallingPriorityIndication == nil {
				u.SignallingPriorityIndication = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.CNOperatorSelectionEntity:
			if u.CNOperatorSelectionEntity == nil {
				u.CNOperatorSelectionEntity = i
			} else {
				u.AdditionalIEs = append(u.AdditionalIEs, i)
			}
		case ie.IMEISV:
			if u.IMEI == nil {
				u.IMEI = i
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
func (u *UpdatePDPContextRequest) MarshalLen() int {
	if u == nil || u.Header == nil {
		return 0
	}
	ies := []*ie.IE{u.IMSI, u.RAI, u.Recovery, u.TEIDDataI, u.TEIDCPlane, u.NSAPI, u.TraceReference, u.TraceType, u.PCO, u.SGSNAddressForCPlane, u.SGSNAddressForUserTraffic, u.AlternativeSGSNAddressForCPlane, u.AlternativeSGSNAddressForUserTraffic, u.QoSProfile, u.TFT, u.TriggerID, u.OMCIdentity, u.CommonFlags, u.RATType, u.ULI, u.MSTimeZone, u.AdditionalTraceInfo, u.DirectTunnelFlags, u.EvolvedARPI, u.ExtendedCommonFlags, u.UCI, u.APNAMBR, u.SignallingPriorityIndication, u.CNOperatorSelectionEntity, u.IMEI}
	ies = append(ies, u.PrivateExtensions...)
	ies = append(ies, u.AdditionalIEs...)
	return messageLen(u.Header, ies)
}

// SetLength sets the length in Length field.
func (u *UpdatePDPContextRequest) SetLength() {
	if u == nil || u.Header == nil {
		return
	}
	u.Length = uint16(u.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (u *UpdatePDPContextRequest) MessageTypeName() string {
	return "Update PDP Context Request"
}

// TEID returns the TEID in human-readable string.
func (u *UpdatePDPContextRequest) TEID() uint32 {
	if u == nil || u.Header == nil {
		return 0
	}
	return u.Header.TEID
}
