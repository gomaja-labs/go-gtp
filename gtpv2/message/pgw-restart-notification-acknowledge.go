// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// PGWRestartNotificationAcknowledge is a PGWRestartNotificationAcknowledge Header and its IEs above.
type PGWRestartNotificationAcknowledge struct {
	*Header
	Cause            *ie.IE
	PrivateExtension *ie.IE
	AdditionalIEs    []*ie.IE
}

// NewPGWRestartNotificationAcknowledge creates a new PGWRestartNotificationAcknowledge.
func NewPGWRestartNotificationAcknowledge(teid, seq uint32, ies ...*ie.IE) *PGWRestartNotificationAcknowledge {
	m := &PGWRestartNotificationAcknowledge{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypePGWRestartNotificationAcknowledge, teid, seq, nil,
		),
	}

	for _, i := range ies {
		m.addIE(i)
	}

	m.SetLength()
	return m
}

// Marshal serializes PGWRestartNotificationAcknowledge into bytes.
func (m *PGWRestartNotificationAcknowledge) Marshal() ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes PGWRestartNotificationAcknowledge into bytes.
func (m *PGWRestartNotificationAcknowledge) MarshalTo(b []byte) error {
	if m == nil || m.Header == nil {
		return ErrNilMessage
	}
	if err := m.Header.validate(MsgTypePGWRestartNotificationAcknowledge, m.MarshalLen()); err != nil {
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

// ParsePGWRestartNotificationAcknowledge decodes given bytes as PGWRestartNotificationAcknowledge.
func ParsePGWRestartNotificationAcknowledge(b []byte) (*PGWRestartNotificationAcknowledge, error) {
	m := &PGWRestartNotificationAcknowledge{}
	if err := m.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return m, err
		}
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes given bytes as PGWRestartNotificationAcknowledge.
func (m *PGWRestartNotificationAcknowledge) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	*m = PGWRestartNotificationAcknowledge{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypePGWRestartNotificationAcknowledge {
		return &UnexpectedTypeError{Expected: MsgTypePGWRestartNotificationAcknowledge, Actual: h.Type}
	}
	m.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, m.addIE)
}

// MarshalLen returns the serial length in int.
func (m *PGWRestartNotificationAcknowledge) MarshalLen() int {
	if m == nil || m.Header == nil {
		return 0
	}
	l := m.Header.MarshalLen() - len(m.Header.Payload)
	if i := m.Cause; i != nil {
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
func (m *PGWRestartNotificationAcknowledge) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.setLength(m.MarshalLen() - m.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *PGWRestartNotificationAcknowledge) MessageTypeName() string {
	return "PGW Restart Notification Acknowledge"
}

// TEID returns the TEID in uint32.
func (m *PGWRestartNotificationAcknowledge) TEID() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.teid()
}

func (m *PGWRestartNotificationAcknowledge) GetHeader() *Header {
	if m == nil {
		return nil
	}
	return m.Header
}

// TS 29.274 Table 7.9.6-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (m *PGWRestartNotificationAcknowledge) addIE(i *ie.IE) {
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
func (m *PGWRestartNotificationAcknowledge) Version() int { return m.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotificationAcknowledge) MessageType() uint8 { return m.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotificationAcknowledge) Sequence() uint32 { return m.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotificationAcknowledge) IsPiggybacking() bool {
	return m.GetHeader().IsPiggybacking()
}

// HasTEID reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotificationAcknowledge) HasTEID() bool { return m.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotificationAcknowledge) HasMessagePriority() bool {
	return m.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotificationAcknowledge) MessagePriority() uint8 {
	return m.GetHeader().MessagePriority()
}

// SetSequenceNumber updates the common header safely on a nil receiver.
func (m *PGWRestartNotificationAcknowledge) SetSequenceNumber(seq uint32) error {
	return m.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (m *PGWRestartNotificationAcknowledge) SetMessagePriority(priority uint8) error {
	return m.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (m *PGWRestartNotificationAcknowledge) SetTEID(teid uint32) { m.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (m *PGWRestartNotificationAcknowledge) SetPiggybacking(value uint8) {
	m.GetHeader().SetPiggybacking(value)
}
