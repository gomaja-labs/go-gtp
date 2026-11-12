// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// PGWRestartNotification is a PGWRestartNotification Header and its IEs above.
type PGWRestartNotification struct {
	*Header
	PGWS5S8IPAddressForControlPlaneOrPMIP *ie.IE
	SGWS11S4IPAddressForControlPlane      *ie.IE
	Cause                                 *ie.IE
	PrivateExtension                      *ie.IE
	AdditionalIEs                         []*ie.IE
}

// NewPGWRestartNotification creates a new PGWRestartNotification.
func NewPGWRestartNotification(teid, seq uint32, ies ...*ie.IE) *PGWRestartNotification {
	m := &PGWRestartNotification{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypePGWRestartNotification, teid, seq, nil,
		),
	}

	for _, i := range ies {
		m.addIE(i)
	}

	m.SetLength()
	return m
}

// Marshal serializes PGWRestartNotification into bytes.
func (m *PGWRestartNotification) Marshal() ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes PGWRestartNotification into bytes.
func (m *PGWRestartNotification) MarshalTo(b []byte) error {
	if m == nil || m.Header == nil {
		return ErrNilMessage
	}
	if err := m.Header.validate(MsgTypePGWRestartNotification, m.MarshalLen()); err != nil {
		return err
	}
	if len(b) < m.MarshalLen() {
		return ErrTooShortToMarshal
	}
	m.Header.Payload = nil
	m.Header.Payload = make([]byte, m.MarshalLen()-m.Header.MarshalLen())
	offset := 0
	if i := m.PGWS5S8IPAddressForControlPlaneOrPMIP; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "PGWS5S8IPAddressForControlPlaneOrPMIP", iePair{ie.IPAddress, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.SGWS11S4IPAddressForControlPlane; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "SGWS11S4IPAddressForControlPlane", iePair{ie.IPAddress, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
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

// ParsePGWRestartNotification decodes given bytes as PGWRestartNotification.
func ParsePGWRestartNotification(b []byte) (*PGWRestartNotification, error) {
	m := &PGWRestartNotification{}
	if err := m.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return m, err
		}
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes given bytes as PGWRestartNotification.
func (m *PGWRestartNotification) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	*m = PGWRestartNotification{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypePGWRestartNotification {
		return &UnexpectedTypeError{Expected: MsgTypePGWRestartNotification, Actual: h.Type}
	}
	m.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, m.addIE)
}

// MarshalLen returns the serial length in int.
func (m *PGWRestartNotification) MarshalLen() int {
	if m == nil || m.Header == nil {
		return 0
	}
	l := m.Header.MarshalLen() - len(m.Header.Payload)
	if i := m.PGWS5S8IPAddressForControlPlaneOrPMIP; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SGWS11S4IPAddressForControlPlane; i != nil {
		l += i.MarshalLen()
	}
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
func (m *PGWRestartNotification) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.setLength(m.MarshalLen() - m.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *PGWRestartNotification) MessageTypeName() string {
	return "PGW Restart Notification"
}

// TEID returns the TEID in uint32.
func (m *PGWRestartNotification) TEID() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.teid()
}

func (m *PGWRestartNotification) GetHeader() *Header {
	if m == nil {
		return nil
	}
	return m.Header
}

// TS 29.274 Table 7.9.5-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (m *PGWRestartNotification) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		m.Header.err = err
	}
	switch i.Type {
	case ie.IPAddress:
		switch i.Instance() {
		case 0:
			if m.PGWS5S8IPAddressForControlPlaneOrPMIP == nil {
				m.PGWS5S8IPAddressForControlPlaneOrPMIP = i
			}
		case 1:
			if m.SGWS11S4IPAddressForControlPlane == nil {
				m.SGWS11S4IPAddressForControlPlane = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
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
func (m *PGWRestartNotification) Version() int { return m.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotification) MessageType() uint8 { return m.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotification) Sequence() uint32 { return m.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotification) IsPiggybacking() bool { return m.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotification) HasTEID() bool { return m.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotification) HasMessagePriority() bool { return m.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (m *PGWRestartNotification) MessagePriority() uint8 { return m.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (m *PGWRestartNotification) SetSequenceNumber(seq uint32) error {
	return m.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (m *PGWRestartNotification) SetMessagePriority(priority uint8) error {
	return m.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (m *PGWRestartNotification) SetTEID(teid uint32) { m.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (m *PGWRestartNotification) SetPiggybacking(value uint8) { m.GetHeader().SetPiggybacking(value) }
