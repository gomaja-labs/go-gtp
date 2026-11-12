// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// ModifyBearerRequest is a ModifyBearerRequest Header and its IEs above.
type ModifyBearerRequest struct {
	*Header
	MEI                                    *ie.IE
	ULI                                    *ie.IE
	ServingNetwork                         *ie.IE
	RATType                                *ie.IE
	IndicationFlags                        *ie.IE
	SenderFTEIDC                           *ie.IE
	AMBR                                   *ie.IE
	DelayDownlinkPacketNotificationRequest *ie.IE
	BearerContextsToBeModified             []*ie.IE
	BearerContextsToBeRemoved              []*ie.IE
	Recovery                               *ie.IE
	UETimeZone                             *ie.IE
	MMEFQCSID                              *ie.IE
	SGWFQCSID                              *ie.IE
	UCI                                    *ie.IE
	UELocalIPAddress                       *ie.IE
	UEUDPPort                              *ie.IE
	MMESGSNLDN                             *ie.IE
	SGWLDN                                 *ie.IE
	HeNBLocalIPAddress                     *ie.IE
	HeNBUDPPort                            *ie.IE
	MMESGSNIdentifier                      *ie.IE
	CNOperatorSelectionEntity              *ie.IE
	PresenceReportingAreaInformation       []*ie.IE
	MMESGSNOverloadControlInformation      *ie.IE
	SGWOverloadControlInformation          *ie.IE
	EPDGOverloadControlInformation         *ie.IE
	ServingPLMNRateControl                 *ie.IE
	MOExceptionDataCounter                 *ie.IE
	IMSI                                   *ie.IE
	ULIForSGW                              *ie.IE
	WLANLocationInformation                *ie.IE
	WLANLocationTimeStamp                  *ie.IE
	SecondaryRATUsageDataReport            []*ie.IE
	PrivateExtension                       *ie.IE
	AdditionalIEs                          []*ie.IE
}

// NewModifyBearerRequest creates a new ModifyBearerRequest.
func NewModifyBearerRequest(teid, seq uint32, ies ...*ie.IE) *ModifyBearerRequest {
	m := &ModifyBearerRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeModifyBearerRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		m.addIE(i)
	}

	m.SetLength()
	return m
}

// Marshal serializes ModifyBearerRequest into bytes.
func (m *ModifyBearerRequest) Marshal() ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ModifyBearerRequest into bytes.
func (m *ModifyBearerRequest) MarshalTo(b []byte) error {
	if m == nil || m.Header == nil {
		return ErrNilMessage
	}
	if err := m.Header.validate(MsgTypeModifyBearerRequest, m.MarshalLen()); err != nil {
		return err
	}
	if len(b) < m.MarshalLen() {
		return ErrTooShortToMarshal
	}
	m.Header.Payload = nil
	m.Header.Payload = make([]byte, m.MarshalLen()-m.Header.MarshalLen())
	offset := 0
	if i := m.MEI; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "MEI", iePair{ie.MobileEquipmentIdentity, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.ULI; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "ULI", iePair{ie.UserLocationInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.ServingNetwork; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "ServingNetwork", iePair{ie.ServingNetwork, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.RATType; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "RATType", iePair{ie.RATType, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.IndicationFlags; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "IndicationFlags", iePair{ie.Indication, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.SenderFTEIDC; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "SenderFTEIDC", iePair{ie.FullyQualifiedTEID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.AMBR; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "AMBR", iePair{ie.AggregateMaximumBitRate, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.DelayDownlinkPacketNotificationRequest; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "DelayDownlinkPacketNotificationRequest", iePair{ie.DelayValue, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.BearerContextsToBeModified {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "BearerContextsToBeModified", iePair{ie.BearerContext, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.BearerContextsToBeRemoved {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "BearerContextsToBeRemoved", iePair{ie.BearerContext, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.Recovery; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "Recovery", iePair{ie.Recovery, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.UETimeZone; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "UETimeZone", iePair{ie.UETimeZone, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.MMEFQCSID; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "MMEFQCSID", iePair{ie.FullyQualifiedCSID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.SGWFQCSID; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "SGWFQCSID", iePair{ie.FullyQualifiedCSID, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.UCI; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "UCI", iePair{ie.UserCSGInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.UELocalIPAddress; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "UELocalIPAddress", iePair{ie.IPAddress, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.UEUDPPort; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "UEUDPPort", iePair{ie.PortNumber, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.MMESGSNLDN; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "MMESGSNLDN", iePair{ie.LocalDistinguishedName, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.SGWLDN; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "SGWLDN", iePair{ie.LocalDistinguishedName, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.HeNBLocalIPAddress; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "HeNBLocalIPAddress", iePair{ie.IPAddress, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.HeNBUDPPort; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "HeNBUDPPort", iePair{ie.PortNumber, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.MMESGSNIdentifier; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "MMESGSNIdentifier", iePair{ie.IPAddress, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.CNOperatorSelectionEntity; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "CNOperatorSelectionEntity", iePair{ie.CNOperatorSelectionEntity, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.PresenceReportingAreaInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "PresenceReportingAreaInformation", iePair{ie.PresenceReportingAreaInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.MMESGSNOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "MMESGSNOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.SGWOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "SGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.EPDGOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "EPDGOverloadControlInformation", iePair{ie.OverloadControlInformation, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.ServingPLMNRateControl; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "ServingPLMNRateControl", iePair{ie.ServingPLMNRateControl, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.MOExceptionDataCounter; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "MOExceptionDataCounter", iePair{ie.Counter, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.IMSI; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "IMSI", iePair{ie.IMSI, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.ULIForSGW; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "ULIForSGW", iePair{ie.UserLocationInformation, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.WLANLocationInformation; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "WLANLocationInformation", iePair{ie.TWANIdentifier, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.WLANLocationTimeStamp; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "WLANLocationTimeStamp", iePair{ie.TWANIdentifierTimestamp, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.SecondaryRATUsageDataReport {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "SecondaryRATUsageDataReport", iePair{ie.SecondaryRATUsageDataReport, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.PrivateExtension; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "PrivateExtension", iePair{ie.PrivateExtension, 255}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.AdditionalIEs {
		if i == nil {
			continue
		}
		if err := i.MarshalTo(m.Payload[offset:]); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	m.Header.SetLength()
	return m.Header.MarshalTo(b)
}

// ParseModifyBearerRequest decodes given bytes as ModifyBearerRequest.
func ParseModifyBearerRequest(b []byte) (*ModifyBearerRequest, error) {
	c := &ModifyBearerRequest{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as ModifyBearerRequest.
func (m *ModifyBearerRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	*m = ModifyBearerRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeModifyBearerRequest {
		return &UnexpectedTypeError{Expected: MsgTypeModifyBearerRequest, Actual: h.Type}
	}
	m.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, m.addIE)
}

// MarshalLen returns the serial length in int.
func (m *ModifyBearerRequest) MarshalLen() int {
	if m == nil || m.Header == nil {
		return 0
	}
	l := m.Header.MarshalLen() - len(m.Header.Payload)
	if i := m.MEI; i != nil {
		l += i.MarshalLen()
	}
	if i := m.ULI; i != nil {
		l += i.MarshalLen()
	}
	if i := m.ServingNetwork; i != nil {
		l += i.MarshalLen()
	}
	if i := m.RATType; i != nil {
		l += i.MarshalLen()
	}
	if i := m.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SenderFTEIDC; i != nil {
		l += i.MarshalLen()
	}
	if i := m.AMBR; i != nil {
		l += i.MarshalLen()
	}
	if i := m.DelayDownlinkPacketNotificationRequest; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range m.BearerContextsToBeModified {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range m.BearerContextsToBeRemoved {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := m.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := m.UETimeZone; i != nil {
		l += i.MarshalLen()
	}
	if i := m.MMEFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := m.UCI; i != nil {
		l += i.MarshalLen()
	}
	if i := m.UELocalIPAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := m.UEUDPPort; i != nil {
		l += i.MarshalLen()
	}
	if i := m.MMESGSNLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SGWLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := m.HeNBLocalIPAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := m.HeNBUDPPort; i != nil {
		l += i.MarshalLen()
	}
	if i := m.MMESGSNIdentifier; i != nil {
		l += i.MarshalLen()
	}
	if i := m.CNOperatorSelectionEntity; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range m.PresenceReportingAreaInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := m.MMESGSNOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := m.EPDGOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := m.ServingPLMNRateControl; i != nil {
		l += i.MarshalLen()
	}
	if i := m.MOExceptionDataCounter; i != nil {
		l += i.MarshalLen()
	}
	if i := m.IMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := m.ULIForSGW; i != nil {
		l += i.MarshalLen()
	}
	if i := m.WLANLocationInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := m.WLANLocationTimeStamp; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range m.SecondaryRATUsageDataReport {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := m.PrivateExtension; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range m.AdditionalIEs {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	return l
}

// SetLength sets the length in Length field.
func (m *ModifyBearerRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.setLength(m.MarshalLen() - m.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *ModifyBearerRequest) MessageTypeName() string {
	return "Modify Bearer Request"
}

// TEID returns the TEID in uint32.
func (m *ModifyBearerRequest) TEID() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.teid()
}

func (m *ModifyBearerRequest) GetHeader() *Header {
	if m == nil {
		return nil
	}
	return m.Header
}

// TS 29.274 Table 7.2.7-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (m *ModifyBearerRequest) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		m.Header.err = err
	}
	switch i.Type {
	case ie.MobileEquipmentIdentity:
		switch i.Instance() {
		case 0:
			if m.MEI == nil {
				m.MEI = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.UserLocationInformation:
		switch i.Instance() {
		case 0:
			if m.ULI == nil {
				m.ULI = i
			}
		case 1:
			if m.ULIForSGW == nil {
				m.ULIForSGW = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.ServingNetwork:
		switch i.Instance() {
		case 0:
			if m.ServingNetwork == nil {
				m.ServingNetwork = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.RATType:
		switch i.Instance() {
		case 0:
			if m.RATType == nil {
				m.RATType = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.Indication:
		switch i.Instance() {
		case 0:
			if m.IndicationFlags == nil {
				m.IndicationFlags = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.FullyQualifiedTEID:
		switch i.Instance() {
		case 0:
			if m.SenderFTEIDC == nil {
				m.SenderFTEIDC = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.AggregateMaximumBitRate:
		switch i.Instance() {
		case 0:
			if m.AMBR == nil {
				m.AMBR = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.DelayValue:
		switch i.Instance() {
		case 0:
			if m.DelayDownlinkPacketNotificationRequest == nil {
				m.DelayDownlinkPacketNotificationRequest = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.BearerContext:
		switch i.Instance() {
		case 0:
			m.BearerContextsToBeModified = append(m.BearerContextsToBeModified, i)
		case 1:
			m.BearerContextsToBeRemoved = append(m.BearerContextsToBeRemoved, i)
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.Recovery:
		switch i.Instance() {
		case 0:
			if m.Recovery == nil {
				m.Recovery = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.UETimeZone:
		switch i.Instance() {
		case 0:
			if m.UETimeZone == nil {
				m.UETimeZone = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.FullyQualifiedCSID:
		switch i.Instance() {
		case 0:
			if m.MMEFQCSID == nil {
				m.MMEFQCSID = i
			}
		case 1:
			if m.SGWFQCSID == nil {
				m.SGWFQCSID = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.UserCSGInformation:
		switch i.Instance() {
		case 0:
			if m.UCI == nil {
				m.UCI = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.IPAddress:
		switch i.Instance() {
		case 1:
			if m.UELocalIPAddress == nil {
				m.UELocalIPAddress = i
			}
		case 0:
			if m.HeNBLocalIPAddress == nil {
				m.HeNBLocalIPAddress = i
			}
		case 2:
			if m.MMESGSNIdentifier == nil {
				m.MMESGSNIdentifier = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.PortNumber:
		switch i.Instance() {
		case 1:
			if m.UEUDPPort == nil {
				m.UEUDPPort = i
			}
		case 0:
			if m.HeNBUDPPort == nil {
				m.HeNBUDPPort = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.LocalDistinguishedName:
		switch i.Instance() {
		case 0:
			if m.MMESGSNLDN == nil {
				m.MMESGSNLDN = i
			}
		case 1:
			if m.SGWLDN == nil {
				m.SGWLDN = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.CNOperatorSelectionEntity:
		switch i.Instance() {
		case 0:
			if m.CNOperatorSelectionEntity == nil {
				m.CNOperatorSelectionEntity = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.PresenceReportingAreaInformation:
		switch i.Instance() {
		case 0:
			m.PresenceReportingAreaInformation = append(m.PresenceReportingAreaInformation, i)
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			if m.MMESGSNOverloadControlInformation == nil {
				m.MMESGSNOverloadControlInformation = i
			}
		case 1:
			if m.SGWOverloadControlInformation == nil {
				m.SGWOverloadControlInformation = i
			}
		case 2:
			if m.EPDGOverloadControlInformation == nil {
				m.EPDGOverloadControlInformation = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.ServingPLMNRateControl:
		switch i.Instance() {
		case 0:
			if m.ServingPLMNRateControl == nil {
				m.ServingPLMNRateControl = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.Counter:
		switch i.Instance() {
		case 0:
			if m.MOExceptionDataCounter == nil {
				m.MOExceptionDataCounter = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.IMSI:
		switch i.Instance() {
		case 0:
			if m.IMSI == nil {
				m.IMSI = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.TWANIdentifier:
		switch i.Instance() {
		case 0:
			if m.WLANLocationInformation == nil {
				m.WLANLocationInformation = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.TWANIdentifierTimestamp:
		switch i.Instance() {
		case 0:
			if m.WLANLocationTimeStamp == nil {
				m.WLANLocationTimeStamp = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.SecondaryRATUsageDataReport:
		switch i.Instance() {
		case 0:
			m.SecondaryRATUsageDataReport = append(m.SecondaryRATUsageDataReport, i)
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.PrivateExtension:
		if m.PrivateExtension == nil {
			m.PrivateExtension = i
		} else if m.PrivateExtension.Instance() != i.Instance() {
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	default:
		m.AdditionalIEs = append(m.AdditionalIEs, i)
	}
}

// Version reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerRequest) Version() int { return m.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerRequest) MessageType() uint8 { return m.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerRequest) Sequence() uint32 { return m.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerRequest) IsPiggybacking() bool { return m.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerRequest) HasTEID() bool { return m.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerRequest) HasMessagePriority() bool { return m.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerRequest) MessagePriority() uint8 { return m.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (m *ModifyBearerRequest) SetSequenceNumber(seq uint32) error {
	return m.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (m *ModifyBearerRequest) SetMessagePriority(priority uint8) error {
	return m.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (m *ModifyBearerRequest) SetTEID(teid uint32) { m.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (m *ModifyBearerRequest) SetPiggybacking(value uint8) { m.GetHeader().SetPiggybacking(value) }
