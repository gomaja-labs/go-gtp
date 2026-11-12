// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// EchoRequest is a EchoRequest Header and its IEs above.
type EchoRequest struct {
	*Header
	Recovery            *ie.IE
	SendingNodeFeatures *ie.IE
	PrivateExtension    *ie.IE
	AdditionalIEs       []*ie.IE
}

// NewEchoRequest creates a new EchoRequest.
func NewEchoRequest(seq uint32, ies ...*ie.IE) *EchoRequest {
	e := &EchoRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 0),
			MsgTypeEchoRequest, 0, seq, nil,
		),
	}

	for _, i := range ies {
		e.addIE(i)
	}

	e.SetLength()
	return e
}

// Marshal returns the byte sequence generated from a EchoRequest.
func (e *EchoRequest) Marshal() ([]byte, error) {
	b := make([]byte, e.MarshalLen())
	if err := e.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (e *EchoRequest) MarshalTo(b []byte) error {
	if e == nil || e.Header == nil {
		return ErrNilMessage
	}
	if err := e.Header.validate(MsgTypeEchoRequest, e.MarshalLen()); err != nil {
		return err
	}
	if len(b) < e.MarshalLen() {
		return ErrTooShortToMarshal
	}
	e.Header.Payload = nil
	e.Header.Payload = make([]byte, e.MarshalLen()-e.Header.MarshalLen())
	offset := 0
	if i := e.Recovery; i != nil {
		if err := marshalNamedIE(i, e.Payload[offset:], "Recovery", iePair{ie.Recovery, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := e.SendingNodeFeatures; i != nil {
		if err := marshalNamedIE(i, e.Payload[offset:], "SendingNodeFeatures", iePair{ie.NodeFeatures, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := e.PrivateExtension; i != nil {
		if err := marshalNamedIE(i, e.Payload[offset:], "PrivateExtension", iePair{ie.PrivateExtension, 255}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range e.AdditionalIEs {
		if i == nil {
			continue
		}
		if err := i.MarshalTo(e.Payload[offset:]); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	e.Header.SetLength()
	return e.Header.MarshalTo(b)
}

// ParseEchoRequest decodes a given byte sequence as a EchoRequest.
func ParseEchoRequest(b []byte) (*EchoRequest, error) {
	e := &EchoRequest{}
	if err := e.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return e, err
		}
		return nil, err
	}
	return e, nil
}

// UnmarshalBinary decodes a given byte sequence as a EchoRequest.
func (e *EchoRequest) UnmarshalBinary(b []byte) error {
	if e == nil {
		return ErrNilMessage
	}
	*e = EchoRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeEchoRequest {
		return &UnexpectedTypeError{Expected: MsgTypeEchoRequest, Actual: h.Type}
	}
	e.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, e.addIE)
}

// MarshalLen returns the serial length of Data.
func (e *EchoRequest) MarshalLen() int {
	if e == nil || e.Header == nil {
		return 0
	}
	l := e.Header.MarshalLen() - len(e.Header.Payload)
	if i := e.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := e.SendingNodeFeatures; i != nil {
		l += i.MarshalLen()
	}
	if i := e.PrivateExtension; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range e.AdditionalIEs {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	return l
}

// SetLength sets the length in Length field.
func (e *EchoRequest) SetLength() {
	if e != nil && e.Header != nil {
		e.Header.setLength(e.MarshalLen() - e.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (e *EchoRequest) MessageTypeName() string {
	return "Echo Request"
}

// TEID returns the TEID in uint32.
func (e *EchoRequest) TEID() uint32 {
	if e == nil {
		return 0
	}
	return e.Header.teid()
}

func (e *EchoRequest) GetHeader() *Header {
	if e == nil {
		return nil
	}
	return e.Header
}

// TS 29.274 Table 7.1.1-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (e *EchoRequest) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		e.Header.err = err
	}
	switch i.Type {
	case ie.Recovery:
		switch i.Instance() {
		case 0:
			if e.Recovery == nil {
				e.Recovery = i
			}
		default:
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	case ie.NodeFeatures:
		switch i.Instance() {
		case 0:
			if e.SendingNodeFeatures == nil {
				e.SendingNodeFeatures = i
			}
		default:
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	case ie.PrivateExtension:
		if e.PrivateExtension == nil {
			e.PrivateExtension = i
		} else if e.PrivateExtension.Instance() != i.Instance() {
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	default:
		e.AdditionalIEs = append(e.AdditionalIEs, i)
	}
}

// Version reads the common header, returning zero for a nil receiver.
func (e *EchoRequest) Version() int { return e.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (e *EchoRequest) MessageType() uint8 { return e.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (e *EchoRequest) Sequence() uint32 { return e.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (e *EchoRequest) IsPiggybacking() bool { return e.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (e *EchoRequest) HasTEID() bool { return e.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (e *EchoRequest) HasMessagePriority() bool { return e.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (e *EchoRequest) MessagePriority() uint8 { return e.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (e *EchoRequest) SetSequenceNumber(seq uint32) error {
	return e.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (e *EchoRequest) SetMessagePriority(priority uint8) error {
	return e.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (e *EchoRequest) SetTEID(teid uint32) { e.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (e *EchoRequest) SetPiggybacking(value uint8) { e.GetHeader().SetPiggybacking(value) }
