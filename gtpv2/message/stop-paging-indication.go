// Copyright 2019-2024 go-gtp authors. All rights reservev.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

// StopPagingIndication is a StopPagingIndication Header and its IEs above.
type StopPagingIndication struct {
	*Header
	IMSI             *ie.IE
	PrivateExtension *ie.IE
	AdditionalIEs    []*ie.IE
}

// NewStopPagingIndication creates a new StopPagingIndication.
func NewStopPagingIndication(teid, seq uint32, ies ...*ie.IE) *StopPagingIndication {
	s := &StopPagingIndication{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeStopPagingIndication, teid, seq, nil,
		),
	}

	for _, i := range ies {
		s.addIE(i)
	}

	s.SetLength()
	return s
}

// Marshal serializes StopPagingIndication into bytes.
func (s *StopPagingIndication) Marshal() ([]byte, error) {
	b := make([]byte, s.MarshalLen())
	if err := s.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes StopPagingIndication into bytes.
func (s *StopPagingIndication) MarshalTo(b []byte) error {
	if s == nil || s.Header == nil {
		return ErrNilMessage
	}
	if err := s.Header.validate(MsgTypeStopPagingIndication, s.MarshalLen()); err != nil {
		return err
	}
	if len(b) < s.MarshalLen() {
		return ErrTooShortToMarshal
	}
	s.Header.Payload = nil
	s.Header.Payload = make([]byte, s.MarshalLen()-s.Header.MarshalLen())
	offset := 0
	if i := s.IMSI; i != nil {
		if err := marshalNamedIE(i, s.Payload[offset:], "IMSI", iePair{ie.IMSI, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := s.PrivateExtension; i != nil {
		if err := marshalNamedIE(i, s.Payload[offset:], "PrivateExtension", iePair{ie.PrivateExtension, 255}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range s.AdditionalIEs {
		if i == nil {
			continue
		}
		if err := i.MarshalTo(s.Payload[offset:]); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	s.Header.SetLength()
	return s.Header.MarshalTo(b)
}

// ParseStopPagingIndication decodes given bytes as StopPagingIndication.
func ParseStopPagingIndication(b []byte) (*StopPagingIndication, error) {
	s := &StopPagingIndication{}
	if err := s.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return s, err
		}
		return nil, err
	}
	return s, nil
}

// UnmarshalBinary decodes given bytes as StopPagingIndication.
func (s *StopPagingIndication) UnmarshalBinary(b []byte) error {
	if s == nil {
		return ErrNilMessage
	}
	*s = StopPagingIndication{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeStopPagingIndication {
		return &UnexpectedTypeError{Expected: MsgTypeStopPagingIndication, Actual: h.Type}
	}
	s.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, s.addIE)
}

// MarshalLen returns the serial length in int.
func (s *StopPagingIndication) MarshalLen() int {
	if s == nil || s.Header == nil {
		return 0
	}
	l := s.Header.MarshalLen() - len(s.Header.Payload)
	if i := s.IMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := s.PrivateExtension; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range s.AdditionalIEs {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	return l
}

// SetLength sets the length in Length field.
func (s *StopPagingIndication) SetLength() {
	if s != nil && s.Header != nil {
		s.Header.setLength(s.MarshalLen() - s.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (s *StopPagingIndication) MessageTypeName() string {
	return "Stop Paging Indication"
}

// TEID returns the TEID in uint32.
func (s *StopPagingIndication) TEID() uint32 {
	if s == nil {
		return 0
	}
	return s.Header.teid()
}

func (s *StopPagingIndication) GetHeader() *Header {
	if s == nil {
		return nil
	}
	return s.Header
}

// TS 29.274 Table 7.2.23-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (s *StopPagingIndication) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		s.Header.err = err
	}
	switch i.Type {
	case ie.IMSI:
		switch i.Instance() {
		case 0:
			if s.IMSI == nil {
				s.IMSI = i
			}
		default:
			s.AdditionalIEs = append(s.AdditionalIEs, i)
		}
	case ie.PrivateExtension:
		if s.PrivateExtension == nil {
			s.PrivateExtension = i
		} else if s.PrivateExtension.Instance() != i.Instance() {
			s.AdditionalIEs = append(s.AdditionalIEs, i)
		}
	default:
		s.AdditionalIEs = append(s.AdditionalIEs, i)
	}
}

// Version reads the common header, returning zero for a nil receiver.
func (s *StopPagingIndication) Version() int { return s.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (s *StopPagingIndication) MessageType() uint8 { return s.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (s *StopPagingIndication) Sequence() uint32 { return s.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (s *StopPagingIndication) IsPiggybacking() bool { return s.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (s *StopPagingIndication) HasTEID() bool { return s.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (s *StopPagingIndication) HasMessagePriority() bool { return s.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (s *StopPagingIndication) MessagePriority() uint8 { return s.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (s *StopPagingIndication) SetSequenceNumber(seq uint32) error {
	return s.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (s *StopPagingIndication) SetMessagePriority(priority uint8) error {
	return s.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (s *StopPagingIndication) SetTEID(teid uint32) { s.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (s *StopPagingIndication) SetPiggybacking(value uint8) { s.GetHeader().SetPiggybacking(value) }
