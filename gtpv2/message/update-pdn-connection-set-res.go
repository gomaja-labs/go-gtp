// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// UpdatePDNConnectionSetResponse is a UpdatePDNConnectionSetResponse Header and its IEs above.
type UpdatePDNConnectionSetResponse struct {
	*Header
	Cause            *ie.IE
	PGWFQCSID        *ie.IE
	Recovery         *ie.IE
	PrivateExtension *ie.IE
	AdditionalIEs    []*ie.IE
}

// NewUpdatePDNConnectionSetResponse creates a new UpdatePDNConnectionSetResponse.
func NewUpdatePDNConnectionSetResponse(teid, seq uint32, ies ...*ie.IE) *UpdatePDNConnectionSetResponse {
	m := &UpdatePDNConnectionSetResponse{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeUpdatePDNConnectionSetResponse, teid, seq, nil,
		),
	}

	for _, i := range ies {
		m.addIE(i)
	}

	m.SetLength()
	return m
}

// Marshal returns the byte sequence generated from a UpdatePDNConnectionSetResponse.
func (m *UpdatePDNConnectionSetResponse) Marshal() ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (m *UpdatePDNConnectionSetResponse) MarshalTo(b []byte) error {
	if m == nil || m.Header == nil {
		return ErrNilMessage
	}
	if err := m.Header.validate(MsgTypeUpdatePDNConnectionSetResponse, m.MarshalLen()); err != nil {
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
	if i := m.PGWFQCSID; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "PGWFQCSID", iePair{ie.FullyQualifiedCSID, 0}); err != nil {
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

// ParseUpdatePDNConnectionSetResponse decodes a given byte sequence as a UpdatePDNConnectionSetResponse.
func ParseUpdatePDNConnectionSetResponse(b []byte) (*UpdatePDNConnectionSetResponse, error) {
	m := &UpdatePDNConnectionSetResponse{}
	if err := m.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return m, err
		}
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes given byte sequence as UpdatePDNConnectionSetResponse.
func (m *UpdatePDNConnectionSetResponse) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	*m = UpdatePDNConnectionSetResponse{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeUpdatePDNConnectionSetResponse {
		return &UnexpectedTypeError{Expected: MsgTypeUpdatePDNConnectionSetResponse, Actual: h.Type}
	}
	m.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, m.addIE)
}

// MarshalLen returns the serial length of Data.
func (m *UpdatePDNConnectionSetResponse) MarshalLen() int {
	if m == nil || m.Header == nil {
		return 0
	}
	l := m.Header.MarshalLen() - len(m.Header.Payload)
	if i := m.Cause; i != nil {
		l += i.MarshalLen()
	}
	if i := m.PGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := m.Recovery; i != nil {
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
func (m *UpdatePDNConnectionSetResponse) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.setLength(m.MarshalLen() - m.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *UpdatePDNConnectionSetResponse) MessageTypeName() string {
	return "Update PDN Connection Set Response"
}

// TEID returns the TEID in uint32.
func (m *UpdatePDNConnectionSetResponse) TEID() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.teid()
}

func (m *UpdatePDNConnectionSetResponse) GetHeader() *Header {
	if m == nil {
		return nil
	}
	return m.Header
}

// TS 29.274 Table 7.9.4-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (m *UpdatePDNConnectionSetResponse) addIE(i *ie.IE) {
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
	case ie.FullyQualifiedCSID:
		switch i.Instance() {
		case 0:
			if m.PGWFQCSID == nil {
				m.PGWFQCSID = i
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
func (m *UpdatePDNConnectionSetResponse) Version() int { return m.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetResponse) MessageType() uint8 { return m.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetResponse) Sequence() uint32 { return m.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetResponse) IsPiggybacking() bool { return m.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetResponse) HasTEID() bool { return m.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetResponse) HasMessagePriority() bool {
	return m.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (m *UpdatePDNConnectionSetResponse) MessagePriority() uint8 {
	return m.GetHeader().MessagePriority()
}

// SetSequenceNumber updates the common header safely on a nil receiver.
func (m *UpdatePDNConnectionSetResponse) SetSequenceNumber(seq uint32) error {
	return m.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (m *UpdatePDNConnectionSetResponse) SetMessagePriority(priority uint8) error {
	return m.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (m *UpdatePDNConnectionSetResponse) SetTEID(teid uint32) { m.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (m *UpdatePDNConnectionSetResponse) SetPiggybacking(value uint8) {
	m.GetHeader().SetPiggybacking(value)
}
