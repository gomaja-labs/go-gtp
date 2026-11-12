// Copyright 2019-2024 go-gtp authors. All rights reservev.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

// ModifyAccessBearersRequest is a ModifyAccessBearersRequest Header and its IEs above.
type ModifyAccessBearersRequest struct {
	*Header
	IndicationFlags                        *ie.IE
	SenderFTEIDC                           *ie.IE
	DelayDownlinkPacketNotificationRequest *ie.IE
	BearerContextsToBeModified             []*ie.IE
	BearerContextsToBeRemoved              []*ie.IE
	Recovery                               *ie.IE
	SecondaryRATUsageDataReport            []*ie.IE
	PrivateExtension                       *ie.IE
	AdditionalIEs                          []*ie.IE
}

// NewModifyAccessBearersRequest creates a new ModifyAccessBearersRequest.
func NewModifyAccessBearersRequest(teid, seq uint32, ies ...*ie.IE) *ModifyAccessBearersRequest {
	m := &ModifyAccessBearersRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeModifyAccessBearersRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		m.addIE(i)
	}

	m.SetLength()
	return m
}

// Marshal serializes ModifyAccessBearersRequest into bytes.
func (m *ModifyAccessBearersRequest) Marshal() ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ModifyAccessBearersRequest into bytes.
func (m *ModifyAccessBearersRequest) MarshalTo(b []byte) error {
	if m == nil || m.Header == nil {
		return ErrNilMessage
	}
	if err := m.Header.validate(MsgTypeModifyAccessBearersRequest, m.MarshalLen()); err != nil {
		return err
	}
	if len(b) < m.MarshalLen() {
		return ErrTooShortToMarshal
	}
	m.Header.Payload = nil
	m.Header.Payload = make([]byte, m.MarshalLen()-m.Header.MarshalLen())
	offset := 0
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

// ParseModifyAccessBearersRequest decodes given bytes as ModifyAccessBearersRequest.
func ParseModifyAccessBearersRequest(b []byte) (*ModifyAccessBearersRequest, error) {
	m := &ModifyAccessBearersRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return m, err
		}
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes given bytes as ModifyAccessBearersRequest.
func (m *ModifyAccessBearersRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	*m = ModifyAccessBearersRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeModifyAccessBearersRequest {
		return &UnexpectedTypeError{Expected: MsgTypeModifyAccessBearersRequest, Actual: h.Type}
	}
	m.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, m.addIE)
}

// MarshalLen returns the serial length in int.
func (m *ModifyAccessBearersRequest) MarshalLen() int {
	if m == nil || m.Header == nil {
		return 0
	}
	l := m.Header.MarshalLen() - len(m.Header.Payload)
	if i := m.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SenderFTEIDC; i != nil {
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
func (m *ModifyAccessBearersRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.setLength(m.MarshalLen() - m.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *ModifyAccessBearersRequest) MessageTypeName() string {
	return "Modify Access Bearers Request"
}

// TEID returns the TEID in uint32.
func (m *ModifyAccessBearersRequest) TEID() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.teid()
}

func (m *ModifyAccessBearersRequest) GetHeader() *Header {
	if m == nil {
		return nil
	}
	return m.Header
}

// TS 29.274 Table 7.2.24-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (m *ModifyAccessBearersRequest) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		m.Header.err = err
	}
	switch i.Type {
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
func (m *ModifyAccessBearersRequest) Version() int { return m.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersRequest) MessageType() uint8 { return m.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersRequest) Sequence() uint32 { return m.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersRequest) IsPiggybacking() bool { return m.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersRequest) HasTEID() bool { return m.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersRequest) HasMessagePriority() bool {
	return m.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersRequest) MessagePriority() uint8 { return m.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (m *ModifyAccessBearersRequest) SetSequenceNumber(seq uint32) error {
	return m.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (m *ModifyAccessBearersRequest) SetMessagePriority(priority uint8) error {
	return m.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (m *ModifyAccessBearersRequest) SetTEID(teid uint32) { m.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (m *ModifyAccessBearersRequest) SetPiggybacking(value uint8) {
	m.GetHeader().SetPiggybacking(value)
}
