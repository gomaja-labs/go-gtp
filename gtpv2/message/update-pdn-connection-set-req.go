// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// UpdatePDNConnectionSetRequest is a UpdatePDNConnectionSetRequest Header and its IEs above.
type UpdatePDNConnectionSetRequest struct {
	*Header
	MMEFQCSID        *ie.IE
	SGWFQCSID        *ie.IE
	PrivateExtension *ie.IE
	AdditionalIEs    []*ie.IE
}

// NewUpdatePDNConnectionSetRequest creates a new UpdatePDNConnectionSetRequest.
func NewUpdatePDNConnectionSetRequest(teid, seq uint32, ies ...*ie.IE) *UpdatePDNConnectionSetRequest {
	m := &UpdatePDNConnectionSetRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeUpdatePDNConnectionSetRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		m.addIE(i)
	}

	m.SetLength()
	return m
}

// Marshal serializes UpdatePDNConnectionSetRequest into bytes.
func (m *UpdatePDNConnectionSetRequest) Marshal() ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes UpdatePDNConnectionSetRequest into bytes.
func (m *UpdatePDNConnectionSetRequest) MarshalTo(b []byte) error {
	if m == nil || m.Header == nil {
		return ErrNilMessage
	}
	if err := m.Header.validate(MsgTypeUpdatePDNConnectionSetRequest, m.MarshalLen()); err != nil {
		return err
	}
	if len(b) < m.MarshalLen() {
		return ErrTooShortToMarshal
	}
	m.Header.Payload = nil
	m.Header.Payload = make([]byte, m.MarshalLen()-m.Header.MarshalLen())
	offset := 0
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

// ParseUpdatePDNConnectionSetRequest decodes given bytes as UpdatePDNConnectionSetRequest.
func ParseUpdatePDNConnectionSetRequest(b []byte) (*UpdatePDNConnectionSetRequest, error) {
	m := &UpdatePDNConnectionSetRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return m, err
		}
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes given bytes as UpdatePDNConnectionSetRequest.
func (m *UpdatePDNConnectionSetRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	*m = UpdatePDNConnectionSetRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeUpdatePDNConnectionSetRequest {
		return &UnexpectedTypeError{Expected: MsgTypeUpdatePDNConnectionSetRequest, Actual: h.Type}
	}
	m.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, m.addIE)
}

// MarshalLen returns the serial length in int.
func (m *UpdatePDNConnectionSetRequest) MarshalLen() int {
	if m == nil || m.Header == nil {
		return 0
	}
	l := m.Header.MarshalLen() - len(m.Header.Payload)
	if i := m.MMEFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SGWFQCSID; i != nil {
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
func (m *UpdatePDNConnectionSetRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.setLength(m.MarshalLen() - m.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *UpdatePDNConnectionSetRequest) MessageTypeName() string {
	return "Update PDN Connection Set Request"
}

// TEID returns the TEID in uint32.
func (m *UpdatePDNConnectionSetRequest) TEID() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.teid()
}

func (m *UpdatePDNConnectionSetRequest) GetHeader() *Header {
	if m == nil {
		return nil
	}
	return m.Header
}

// TS 29.274 Table 7.9.3-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (m *UpdatePDNConnectionSetRequest) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		m.Header.err = err
	}
	switch i.Type {
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
func (m *UpdatePDNConnectionSetRequest) Version() int { return m.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetRequest) MessageType() uint8 { return m.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetRequest) Sequence() uint32 { return m.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetRequest) IsPiggybacking() bool { return m.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetRequest) HasTEID() bool { return m.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetRequest) HasMessagePriority() bool {
	return m.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetRequest) MessagePriority() uint8 {
	return m.GetHeader().MessagePriority()
}

// SetSequenceNumber updates the common header safely on a nil receiver.
func (m *UpdatePDNConnectionSetRequest) SetSequenceNumber(seq uint32) error {
	return m.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (m *UpdatePDNConnectionSetRequest) SetMessagePriority(priority uint8) error {
	return m.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (m *UpdatePDNConnectionSetRequest) SetTEID(teid uint32) { m.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (m *UpdatePDNConnectionSetRequest) SetPiggybacking(value uint8) {
	m.GetHeader().SetPiggybacking(value)
}
