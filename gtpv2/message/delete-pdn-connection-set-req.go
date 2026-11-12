// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// DeletePDNConnectionSetRequest is a DeletePDNConnectionSetRequest Header and its IEs above.
type DeletePDNConnectionSetRequest struct {
	*Header
	MMEFQCSID        []*ie.IE
	SGWFQCSID        []*ie.IE
	PGWFQCSID        []*ie.IE
	EPDGFQCSID       []*ie.IE
	TWANFQCSID       []*ie.IE
	PrivateExtension *ie.IE
	AdditionalIEs    []*ie.IE
}

// NewDeletePDNConnectionSetRequest creates a new DeletePDNConnectionSetRequest.
func NewDeletePDNConnectionSetRequest(teid, seq uint32, ies ...*ie.IE) *DeletePDNConnectionSetRequest {
	m := &DeletePDNConnectionSetRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeDeletePDNConnectionSetRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		m.addIE(i)
	}

	m.SetLength()
	return m
}

// Marshal serializes DeletePDNConnectionSetRequest into bytes.
func (m *DeletePDNConnectionSetRequest) Marshal() ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes DeletePDNConnectionSetRequest into bytes.
func (m *DeletePDNConnectionSetRequest) MarshalTo(b []byte) error {
	if m == nil || m.Header == nil {
		return ErrNilMessage
	}
	if err := m.Header.validate(MsgTypeDeletePDNConnectionSetRequest, m.MarshalLen()); err != nil {
		return err
	}
	if len(b) < m.MarshalLen() {
		return ErrTooShortToMarshal
	}
	m.Header.Payload = nil
	m.Header.Payload = make([]byte, m.MarshalLen()-m.Header.MarshalLen())
	offset := 0
	for _, i := range m.MMEFQCSID {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "MMEFQCSID", iePair{ie.FullyQualifiedCSID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.SGWFQCSID {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "SGWFQCSID", iePair{ie.FullyQualifiedCSID, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.PGWFQCSID {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "PGWFQCSID", iePair{ie.FullyQualifiedCSID, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.EPDGFQCSID {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "EPDGFQCSID", iePair{ie.FullyQualifiedCSID, 3}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range m.TWANFQCSID {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, m.Payload[offset:], "TWANFQCSID", iePair{ie.FullyQualifiedCSID, 4}); err != nil {
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

// ParseDeletePDNConnectionSetRequest decodes given bytes as DeletePDNConnectionSetRequest.
func ParseDeletePDNConnectionSetRequest(b []byte) (*DeletePDNConnectionSetRequest, error) {
	m := &DeletePDNConnectionSetRequest{}
	if err := m.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return m, err
		}
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes given bytes as DeletePDNConnectionSetRequest.
func (m *DeletePDNConnectionSetRequest) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	*m = DeletePDNConnectionSetRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeDeletePDNConnectionSetRequest {
		return &UnexpectedTypeError{Expected: MsgTypeDeletePDNConnectionSetRequest, Actual: h.Type}
	}
	m.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, m.addIE)
}

// MarshalLen returns the serial length in int.
func (m *DeletePDNConnectionSetRequest) MarshalLen() int {
	if m == nil || m.Header == nil {
		return 0
	}
	l := m.Header.MarshalLen() - len(m.Header.Payload)
	for _, i := range m.MMEFQCSID {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range m.SGWFQCSID {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range m.PGWFQCSID {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range m.EPDGFQCSID {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range m.TWANFQCSID {
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
func (m *DeletePDNConnectionSetRequest) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.setLength(m.MarshalLen() - m.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *DeletePDNConnectionSetRequest) MessageTypeName() string {
	return "Delete PDN Connection Set Request"
}

// TEID returns the TEID in uint32.
func (m *DeletePDNConnectionSetRequest) TEID() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.teid()
}

func (m *DeletePDNConnectionSetRequest) GetHeader() *Header {
	if m == nil {
		return nil
	}
	return m.Header
}

// TS 29.274 Table 7.9.1-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (m *DeletePDNConnectionSetRequest) addIE(i *ie.IE) {
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
			m.MMEFQCSID = append(m.MMEFQCSID, i)
		case 1:
			m.SGWFQCSID = append(m.SGWFQCSID, i)
		case 2:
			m.PGWFQCSID = append(m.PGWFQCSID, i)
		case 3:
			m.EPDGFQCSID = append(m.EPDGFQCSID, i)
		case 4:
			m.TWANFQCSID = append(m.TWANFQCSID, i)
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
func (m *DeletePDNConnectionSetRequest) Version() int { return m.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (m *DeletePDNConnectionSetRequest) MessageType() uint8 { return m.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (m *DeletePDNConnectionSetRequest) Sequence() uint32 { return m.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (m *DeletePDNConnectionSetRequest) IsPiggybacking() bool { return m.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (m *DeletePDNConnectionSetRequest) HasTEID() bool { return m.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (m *DeletePDNConnectionSetRequest) HasMessagePriority() bool {
	return m.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (m *DeletePDNConnectionSetRequest) MessagePriority() uint8 {
	return m.GetHeader().MessagePriority()
}

// SetSequenceNumber updates the common header safely on a nil receiver.
func (m *DeletePDNConnectionSetRequest) SetSequenceNumber(seq uint32) error {
	return m.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (m *DeletePDNConnectionSetRequest) SetMessagePriority(priority uint8) error {
	return m.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (m *DeletePDNConnectionSetRequest) SetTEID(teid uint32) { m.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (m *DeletePDNConnectionSetRequest) SetPiggybacking(value uint8) {
	m.GetHeader().SetPiggybacking(value)
}
