// Copyright 2019-2024 go-gtp authors. All rights reservev.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

// ModifyAccessBearersResponse is a ModifyAccessBearersResponse Header and its IEs above.
type ModifyAccessBearersResponse struct {
	*Header
	Cause                          *ie.IE
	BearerContextsModified         []*ie.IE
	BearerContextsMarkedForRemoval []*ie.IE
	Recovery                       *ie.IE
	IndicationFlags                *ie.IE
	SGWNodeLoadControlInformation  *ie.IE
	SGWOverloadControlInformation  *ie.IE
	PrivateExtension               *ie.IE
	AdditionalIEs                  []*ie.IE
}

// NewModifyAccessBearersResponse creates a new ModifyAccessBearersResponse.
func NewModifyAccessBearersResponse(teid, seq uint32, ies ...*ie.IE) *ModifyAccessBearersResponse {
	m := &ModifyAccessBearersResponse{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeModifyAccessBearersResponse, teid, seq, nil,
		),
	}

	for _, i := range ies {
		m.addIE(i)
	}

	m.SetLength()
	return m
}

// Marshal serializes ModifyAccessBearersResponse into bytes.
func (m *ModifyAccessBearersResponse) Marshal() ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ModifyAccessBearersResponse into bytes.
func (m *ModifyAccessBearersResponse) MarshalTo(b []byte) error {
	if m == nil || m.Header == nil {
		return ErrNilMessage
	}
	if err := m.Header.validate(MsgTypeModifyAccessBearersResponse, m.MarshalLen()); err != nil {
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
	if i := m.Recovery; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "Recovery", iePair{ie.Recovery, 0}); err != nil {
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
	if i := m.SGWNodeLoadControlInformation; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "SGWNodeLoadControlInformation", iePair{ie.LoadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.SGWOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "SGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
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

// ParseModifyAccessBearersResponse decodes given bytes as ModifyAccessBearersResponse.
func ParseModifyAccessBearersResponse(b []byte) (*ModifyAccessBearersResponse, error) {
	m := &ModifyAccessBearersResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return m, err
		}
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes given bytes as ModifyAccessBearersResponse.
func (m *ModifyAccessBearersResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	*m = ModifyAccessBearersResponse{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeModifyAccessBearersResponse {
		return &UnexpectedTypeError{Expected: MsgTypeModifyAccessBearersResponse, Actual: h.Type}
	}
	m.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, m.addIE)
}

// MarshalLen returns the serial length in int.
func (m *ModifyAccessBearersResponse) MarshalLen() int {
	if m == nil || m.Header == nil {
		return 0
	}
	l := m.Header.MarshalLen() - len(m.Header.Payload)
	if i := m.Cause; i != nil {
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
	if i := m.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := m.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SGWNodeLoadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SGWOverloadControlInformation; i != nil {
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
func (m *ModifyAccessBearersResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.setLength(m.MarshalLen() - m.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *ModifyAccessBearersResponse) MessageTypeName() string {
	return "Modify Access Bearers Response"
}

// TEID returns the TEID in uint32.
func (m *ModifyAccessBearersResponse) TEID() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.teid()
}

func (m *ModifyAccessBearersResponse) GetHeader() *Header {
	if m == nil {
		return nil
	}
	return m.Header
}

// TS 29.274 Table 7.2.25-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (m *ModifyAccessBearersResponse) addIE(i *ie.IE) {
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
	case ie.BearerContext:
		switch i.Instance() {
		case 0:
			m.BearerContextsModified = append(m.BearerContextsModified, i)
		case 1:
			m.BearerContextsMarkedForRemoval = append(m.BearerContextsMarkedForRemoval, i)
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
	case ie.Indication:
		switch i.Instance() {
		case 0:
			if m.IndicationFlags == nil {
				m.IndicationFlags = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.LoadControlInformation:
		switch i.Instance() {
		case 0:
			if m.SGWNodeLoadControlInformation == nil {
				m.SGWNodeLoadControlInformation = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			if m.SGWOverloadControlInformation == nil {
				m.SGWOverloadControlInformation = i
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
func (m *ModifyAccessBearersResponse) Version() int { return m.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersResponse) MessageType() uint8 { return m.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersResponse) Sequence() uint32 { return m.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersResponse) IsPiggybacking() bool { return m.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersResponse) HasTEID() bool { return m.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersResponse) HasMessagePriority() bool {
	return m.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (m *ModifyAccessBearersResponse) MessagePriority() uint8 { return m.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (m *ModifyAccessBearersResponse) SetSequenceNumber(seq uint32) error {
	return m.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (m *ModifyAccessBearersResponse) SetMessagePriority(priority uint8) error {
	return m.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (m *ModifyAccessBearersResponse) SetTEID(teid uint32) { m.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (m *ModifyAccessBearersResponse) SetPiggybacking(value uint8) {
	m.GetHeader().SetPiggybacking(value)
}
