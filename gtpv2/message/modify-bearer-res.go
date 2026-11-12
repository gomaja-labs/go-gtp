// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// ModifyBearerResponse is a ModifyBearerResponse Header and its IEs above.
type ModifyBearerResponse struct {
	*Header
	Cause                          *ie.IE
	MSISDN                         *ie.IE
	LinkedEBI                      *ie.IE
	APNRestriction                 *ie.IE
	PCO                            *ie.IE
	BearerContextsModified         []*ie.IE
	BearerContextsMarkedForRemoval []*ie.IE
	ChangeReportingAction          *ie.IE
	CSGInformationReportingAction  *ie.IE
	HeNBInformationReporting       *ie.IE
	ChargingGatewayName            *ie.IE
	ChargingGatewayAddress         *ie.IE
	PGWFQCSID                      *ie.IE
	SGWFQCSID                      *ie.IE
	Recovery                       *ie.IE
	SGWLDN                         *ie.IE
	PGWLDN                         *ie.IE
	IndicationFlags                *ie.IE
	PresenceReportingAreaAction    []*ie.IE
	PGWNodeLoadControlInformation  *ie.IE
	PGWAPNLoadControlInformation   []*ie.IE
	SGWNodeLoadControlInformation  *ie.IE
	PGWOverloadControlInformation  []*ie.IE
	SGWOverloadControlInformation  *ie.IE
	PDNConnectionChargingID        *ie.IE
	PrivateExtension               *ie.IE
	AdditionalIEs                  []*ie.IE
}

// NewModifyBearerResponse creates a new ModifyBearerResponse.
func NewModifyBearerResponse(teid, seq uint32, ies ...*ie.IE) *ModifyBearerResponse {
	m := &ModifyBearerResponse{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeModifyBearerResponse, teid, seq, nil,
		),
	}

	for _, i := range ies {
		m.addIE(i)
	}

	m.SetLength()
	return m
}

// Marshal serializes ModifyBearerResponse into bytes.
func (m *ModifyBearerResponse) Marshal() ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ModifyBearerResponse into bytes.
func (m *ModifyBearerResponse) MarshalTo(b []byte) error {
	if m == nil || m.Header == nil {
		return ErrNilMessage
	}
	if err := m.Header.validate(MsgTypeModifyBearerResponse, m.MarshalLen()); err != nil {
		return err
	}
	if len(b) < m.MarshalLen() {
		return ErrTooShortToMarshal
	}
	m.Header.Payload = nil
	m.Header.Payload = make([]byte, m.MarshalLen()-m.Header.MarshalLen())
	offset := 0
	if i := m.Cause; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "Cause", iePair{ie.Cause, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.MSISDN; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "MSISDN", iePair{ie.MSISDN, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.LinkedEBI; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "LinkedEBI", iePair{ie.EPSBearerID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.APNRestriction; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "APNRestriction", iePair{ie.APNRestriction, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.PCO; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "PCO", iePair{ie.ProtocolConfigurationOptions, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.BearerContextsModified {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "BearerContextsModified", iePair{ie.BearerContext, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.BearerContextsMarkedForRemoval {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "BearerContextsMarkedForRemoval", iePair{ie.BearerContext, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.ChangeReportingAction; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "ChangeReportingAction", iePair{ie.ChangeReportingAction, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.CSGInformationReportingAction; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "CSGInformationReportingAction", iePair{ie.CSGInformationReportingAction, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.HeNBInformationReporting; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "HeNBInformationReporting", iePair{ie.HeNBInformationReporting, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.ChargingGatewayName; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "ChargingGatewayName", iePair{ie.FullyQualifiedDomainName, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.ChargingGatewayAddress; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "ChargingGatewayAddress", iePair{ie.IPAddress, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.PGWFQCSID; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "PGWFQCSID", iePair{ie.FullyQualifiedCSID, 0}); err != nil {
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
	if i := m.Recovery; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "Recovery", iePair{ie.Recovery, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.SGWLDN; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "SGWLDN", iePair{ie.LocalDistinguishedName, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.PGWLDN; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "PGWLDN", iePair{ie.LocalDistinguishedName, 1}); err != nil {
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
	for _, i := range m.PresenceReportingAreaAction {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "PresenceReportingAreaAction", iePair{ie.PresenceReportingAreaAction, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.PGWNodeLoadControlInformation; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "PGWNodeLoadControlInformation", iePair{ie.LoadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.PGWAPNLoadControlInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "PGWAPNLoadControlInformation", iePair{ie.LoadControlInformation, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.SGWNodeLoadControlInformation; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "SGWNodeLoadControlInformation", iePair{ie.LoadControlInformation, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.PGWOverloadControlInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "PGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
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
	if i := m.PDNConnectionChargingID; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "PDNConnectionChargingID", iePair{ie.ChargingID, 0}); err != nil {
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

// ParseModifyBearerResponse decodes given bytes as ModifyBearerResponse.
func ParseModifyBearerResponse(b []byte) (*ModifyBearerResponse, error) {
	m := &ModifyBearerResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return m, err
		}
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes given bytes as ModifyBearerResponse.
func (m *ModifyBearerResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	*m = ModifyBearerResponse{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeModifyBearerResponse {
		return &UnexpectedTypeError{Expected: MsgTypeModifyBearerResponse, Actual: h.Type}
	}
	m.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, m.addIE)
}

// MarshalLen returns the serial length in int.
func (m *ModifyBearerResponse) MarshalLen() int {
	if m == nil || m.Header == nil {
		return 0
	}
	l := m.Header.MarshalLen() - len(m.Header.Payload)
	if i := m.Cause; i != nil {
		l += i.MarshalLen()
	}
	if i := m.MSISDN; i != nil {
		l += i.MarshalLen()
	}
	if i := m.LinkedEBI; i != nil {
		l += i.MarshalLen()
	}
	if i := m.APNRestriction; i != nil {
		l += i.MarshalLen()
	}
	if i := m.PCO; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range m.BearerContextsModified {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range m.BearerContextsMarkedForRemoval {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := m.ChangeReportingAction; i != nil {
		l += i.MarshalLen()
	}
	if i := m.CSGInformationReportingAction; i != nil {
		l += i.MarshalLen()
	}
	if i := m.HeNBInformationReporting; i != nil {
		l += i.MarshalLen()
	}
	if i := m.ChargingGatewayName; i != nil {
		l += i.MarshalLen()
	}
	if i := m.ChargingGatewayAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := m.PGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := m.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SGWLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := m.PGWLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := m.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range m.PresenceReportingAreaAction {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := m.PGWNodeLoadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range m.PGWAPNLoadControlInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := m.SGWNodeLoadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range m.PGWOverloadControlInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := m.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := m.PDNConnectionChargingID; i != nil {
		l += i.MarshalLen()
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
func (m *ModifyBearerResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.setLength(m.MarshalLen() - m.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *ModifyBearerResponse) MessageTypeName() string {
	return "Modify Bearer Response"
}

// TEID returns the TEID in uint32.
func (m *ModifyBearerResponse) TEID() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.teid()
}

func (m *ModifyBearerResponse) GetHeader() *Header {
	if m == nil {
		return nil
	}
	return m.Header
}

// TS 29.274 Table 7.2.8-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (m *ModifyBearerResponse) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		m.Header.err = err
	}
	switch i.Type {
	case ie.Cause:
		switch i.Instance() {
		case 0:
			if m.Cause == nil {
				m.Cause = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.MSISDN:
		switch i.Instance() {
		case 0:
			if m.MSISDN == nil {
				m.MSISDN = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.EPSBearerID:
		switch i.Instance() {
		case 0:
			if m.LinkedEBI == nil {
				m.LinkedEBI = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.APNRestriction:
		switch i.Instance() {
		case 0:
			if m.APNRestriction == nil {
				m.APNRestriction = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.ProtocolConfigurationOptions:
		switch i.Instance() {
		case 0:
			if m.PCO == nil {
				m.PCO = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.BearerContext:
		switch i.Instance() {
		case 0:
			m.BearerContextsModified = append(m.BearerContextsModified, i)
		case 1:
			m.BearerContextsMarkedForRemoval = append(m.BearerContextsMarkedForRemoval, i)
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.ChangeReportingAction:
		switch i.Instance() {
		case 0:
			if m.ChangeReportingAction == nil {
				m.ChangeReportingAction = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.CSGInformationReportingAction:
		switch i.Instance() {
		case 0:
			if m.CSGInformationReportingAction == nil {
				m.CSGInformationReportingAction = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.HeNBInformationReporting:
		switch i.Instance() {
		case 0:
			if m.HeNBInformationReporting == nil {
				m.HeNBInformationReporting = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.FullyQualifiedDomainName:
		switch i.Instance() {
		case 0:
			if m.ChargingGatewayName == nil {
				m.ChargingGatewayName = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.IPAddress:
		switch i.Instance() {
		case 0:
			if m.ChargingGatewayAddress == nil {
				m.ChargingGatewayAddress = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.FullyQualifiedCSID:
		switch i.Instance() {
		case 0:
			if m.PGWFQCSID == nil {
				m.PGWFQCSID = i
			}
		case 1:
			if m.SGWFQCSID == nil {
				m.SGWFQCSID = i
			}
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
	case ie.LocalDistinguishedName:
		switch i.Instance() {
		case 0:
			if m.SGWLDN == nil {
				m.SGWLDN = i
			}
		case 1:
			if m.PGWLDN == nil {
				m.PGWLDN = i
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
	case ie.PresenceReportingAreaAction:
		switch i.Instance() {
		case 0:
			m.PresenceReportingAreaAction = append(m.PresenceReportingAreaAction, i)
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.LoadControlInformation:
		switch i.Instance() {
		case 0:
			if m.PGWNodeLoadControlInformation == nil {
				m.PGWNodeLoadControlInformation = i
			}
		case 1:
			m.PGWAPNLoadControlInformation = append(m.PGWAPNLoadControlInformation, i)
		case 2:
			if m.SGWNodeLoadControlInformation == nil {
				m.SGWNodeLoadControlInformation = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			m.PGWOverloadControlInformation = append(m.PGWOverloadControlInformation, i)
		case 1:
			if m.SGWOverloadControlInformation == nil {
				m.SGWOverloadControlInformation = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.ChargingID:
		switch i.Instance() {
		case 0:
			if m.PDNConnectionChargingID == nil {
				m.PDNConnectionChargingID = i
			}
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
func (m *ModifyBearerResponse) Version() int { return m.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerResponse) MessageType() uint8 { return m.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerResponse) Sequence() uint32 { return m.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerResponse) IsPiggybacking() bool { return m.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerResponse) HasTEID() bool { return m.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerResponse) HasMessagePriority() bool { return m.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerResponse) MessagePriority() uint8 { return m.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (m *ModifyBearerResponse) SetSequenceNumber(seq uint32) error {
	return m.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (m *ModifyBearerResponse) SetMessagePriority(priority uint8) error {
	return m.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (m *ModifyBearerResponse) SetTEID(teid uint32) { m.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (m *ModifyBearerResponse) SetPiggybacking(value uint8) { m.GetHeader().SetPiggybacking(value) }
