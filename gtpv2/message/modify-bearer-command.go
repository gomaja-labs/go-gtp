// Copyright 2019-2024 go-gtp authors. All rights reservev.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

// ModifyBearerCommand is a ModifyBearerCommand Header and its IEs above.
type ModifyBearerCommand struct {
	*Header
	APNAMBR                            *ie.IE
	BearerContext                      *ie.IE
	MMESGSNOverloadControlInformation  *ie.IE
	SGWOverloadControlInformation      *ie.IE
	TWANePDGOverloadControlInformation *ie.IE
	SenderFTEIDC                       *ie.IE
	PrivateExtension                   *ie.IE
	AdditionalIEs                      []*ie.IE
}

// NewModifyBearerCommand creates a new ModifyBearerCommand.
func NewModifyBearerCommand(teid, seq uint32, ies ...*ie.IE) *ModifyBearerCommand {
	m := &ModifyBearerCommand{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeModifyBearerCommand, teid, seq, nil,
		),
	}

	for _, i := range ies {
		m.addIE(i)
	}

	m.SetLength()
	return m
}

// Marshal serializes ModifyBearerCommand into bytes.
func (m *ModifyBearerCommand) Marshal() ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ModifyBearerCommand into bytes.
func (m *ModifyBearerCommand) MarshalTo(b []byte) error {
	if m == nil || m.Header == nil {
		return ErrNilMessage
	}
	if err := m.Header.validate(MsgTypeModifyBearerCommand, m.MarshalLen()); err != nil {
		return err
	}
	if len(b) < m.MarshalLen() {
		return ErrTooShortToMarshal
	}
	m.Header.Payload = nil
	m.Header.Payload = make([]byte, m.MarshalLen()-m.Header.MarshalLen())
	offset := 0
	if i := m.APNAMBR; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "APNAMBR", iePair{ie.AggregateMaximumBitRate, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := m.BearerContext; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "BearerContext", iePair{ie.BearerContext, 0}); err != nil {
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
	if i := m.TWANePDGOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, m.Payload[offset:], "TWANePDGOverloadControlInformation", iePair{ie.OverloadControlInformation, 2}); err != nil {
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

// ParseModifyBearerCommand decodes given bytes as ModifyBearerCommand.
func ParseModifyBearerCommand(b []byte) (*ModifyBearerCommand, error) {
	m := &ModifyBearerCommand{}
	if err := m.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return m, err
		}
		return nil, err
	}
	return m, nil
}

// UnmarshalBinary decodes given bytes as ModifyBearerCommand.
func (m *ModifyBearerCommand) UnmarshalBinary(b []byte) error {
	if m == nil {
		return ErrNilMessage
	}
	*m = ModifyBearerCommand{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeModifyBearerCommand {
		return &UnexpectedTypeError{Expected: MsgTypeModifyBearerCommand, Actual: h.Type}
	}
	m.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, m.addIE)
}

// MarshalLen returns the serial length in int.
func (m *ModifyBearerCommand) MarshalLen() int {
	if m == nil || m.Header == nil {
		return 0
	}
	l := m.Header.MarshalLen() - len(m.Header.Payload)
	if i := m.APNAMBR; i != nil {
		l += i.MarshalLen()
	}
	if i := m.BearerContext; i != nil {
		l += i.MarshalLen()
	}
	if i := m.MMESGSNOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := m.TWANePDGOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := m.SenderFTEIDC; i != nil {
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
func (m *ModifyBearerCommand) SetLength() {
	if m != nil && m.Header != nil {
		m.Header.setLength(m.MarshalLen() - m.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (m *ModifyBearerCommand) MessageTypeName() string {
	return "Modify Bearer Command"
}

// TEID returns the TEID in uint32.
func (m *ModifyBearerCommand) TEID() uint32 {
	if m == nil {
		return 0
	}
	return m.Header.teid()
}

func (m *ModifyBearerCommand) GetHeader() *Header {
	if m == nil {
		return nil
	}
	return m.Header
}

// TS 29.274 Table 7.2.14.1-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (m *ModifyBearerCommand) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		m.Header.err = err
	}
	switch i.Type {
	case ie.AggregateMaximumBitRate:
		switch i.Instance() {
		case 0:
			if m.APNAMBR == nil {
				m.APNAMBR = i
			}
		default:
			m.AdditionalIEs = append(m.AdditionalIEs, i)
		}
	case ie.BearerContext:
		switch i.Instance() {
		case 0:
			if m.BearerContext == nil {
				m.BearerContext = i
			}
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
			if m.TWANePDGOverloadControlInformation == nil {
				m.TWANePDGOverloadControlInformation = i
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
func (m *ModifyBearerCommand) Version() int { return m.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerCommand) MessageType() uint8 { return m.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerCommand) Sequence() uint32 { return m.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerCommand) IsPiggybacking() bool { return m.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerCommand) HasTEID() bool { return m.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerCommand) HasMessagePriority() bool { return m.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (m *ModifyBearerCommand) MessagePriority() uint8 { return m.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (m *ModifyBearerCommand) SetSequenceNumber(seq uint32) error {
	return m.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (m *ModifyBearerCommand) SetMessagePriority(priority uint8) error {
	return m.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (m *ModifyBearerCommand) SetTEID(teid uint32) { m.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (m *ModifyBearerCommand) SetPiggybacking(value uint8) { m.GetHeader().SetPiggybacking(value) }
