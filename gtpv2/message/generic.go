// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"fmt"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// Generic is a Generic formed GTP Header and its IEs above.
type Generic struct {
	*Header
	IEs []*ie.IE
}

// NewGeneric creates a new Generic.
func NewGeneric(msgType uint8, teid, seq uint32, ie ...*ie.IE) *Generic {
	teidFlag := 0
	if requiresTEID(msgType) {
		teidFlag = 1
	}
	g := &Generic{
		Header: NewHeader(
			NewHeaderFlags(2, 0, teidFlag),
			msgType,
			teid, seq,
			nil,
		),
		IEs: ie,
	}
	g.SetLength()
	return g
}

// NewGenericWithoutTEID creates a new Generic.
func NewGenericWithoutTEID(msgType uint8, teid, seq uint32, ie ...*ie.IE) *Generic {
	g := &Generic{
		Header: NewHeader(
			NewHeaderFlags(2, 0x00, 0),
			msgType, teid, seq, nil,
		),
		IEs: ie,
	}

	g.SetLength()
	return g
}

// Marshal returns the byte sequence generated from a Generic.
func (g *Generic) Marshal() ([]byte, error) {
	b := make([]byte, g.MarshalLen())
	if err := g.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (g *Generic) MarshalTo(b []byte) error {
	if g == nil || g.Header == nil {
		return ErrNilMessage
	}
	if err := g.Header.validate(g.Type, g.MarshalLen()); err != nil {
		return err
	}
	if len(b) < g.MarshalLen() {
		return ErrTooShortToMarshal
	}
	if g.Header.Payload != nil {
		g.Header.Payload = nil
	}
	g.Header.Payload = make([]byte, g.MarshalLen()-g.Header.MarshalLen())

	offset := 0
	for _, ie := range g.IEs {
		if ie == nil {
			continue
		}
		if err := ie.MarshalTo(g.Header.Payload[offset:]); err != nil {
			return err
		}
		offset += ie.MarshalLen()
	}

	g.Header.SetLength()
	return g.Header.MarshalTo(b)
}

// ParseGeneric decodes a given byte sequence as a Generic.
func ParseGeneric(b []byte) (*Generic, error) {
	g := &Generic{}
	if err := g.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return g, nil
}

// UnmarshalBinary decodes a given byte sequence as a Generic. A piggybacked
// message is decoded as a Generic as well, so Marshal reproduces the datagram
// whatever the IE order of either message (TS 29.274 §§5.5.1, 8.1).
func (g *Generic) UnmarshalBinary(b []byte) error {
	if g == nil {
		return ErrNilMessage
	}
	*g = Generic{}
	h := &Header{}
	if err := h.unmarshal(b, parseGenericMessage); err != nil {
		return err
	}
	g.Header = h
	if len(g.Header.Payload) == 0 {
		return nil
	}

	var err error
	g.IEs, err = ie.ParseMultiIEs(g.Header.Payload)
	if err != nil {
		return err
	}
	return nil
}

// parseGenericMessage decodes a piggybacked message for Generic. On failure it
// returns a nil Message rather than a nil *Generic, so Piggybacked stays nil.
func parseGenericMessage(b []byte) (Message, error) {
	g, err := ParseGeneric(b)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// MarshalLen returns the serial length of Data.
func (g *Generic) MarshalLen() int {
	if g == nil || g.Header == nil {
		return 0
	}
	l := g.Header.MarshalLen() - len(g.Header.Payload)

	for _, ie := range g.IEs {
		if ie == nil {
			continue
		}
		l += ie.MarshalLen()
	}

	return l
}

// SetLength sets the length in Length field.
func (g *Generic) SetLength() {
	if g != nil && g.Header != nil {
		g.Header.setLength(g.MarshalLen() - g.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (g *Generic) MessageTypeName() string {
	return fmt.Sprintf("Unknown (%d)", g.GetHeader().MessageType())
}

// TEID returns the TEID in uint32.
func (g *Generic) TEID() uint32 {
	if g == nil {
		return 0
	}
	return g.Header.teid()
}

// AddIE add IEs to Generic type of GTPv2 message and update Length field.
func (g *Generic) AddIE(ie ...*ie.IE) {
	g.IEs = append(g.IEs, ie...)
	g.SetLength()
}

func (g *Generic) GetHeader() *Header {
	if g == nil {
		return nil
	}
	return g.Header
}

// Version reads the common header, returning zero for a nil receiver.
func (g *Generic) Version() int { return g.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (g *Generic) MessageType() uint8 { return g.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (g *Generic) Sequence() uint32 { return g.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (g *Generic) IsPiggybacking() bool { return g.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (g *Generic) HasTEID() bool { return g.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (g *Generic) HasMessagePriority() bool { return g.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (g *Generic) MessagePriority() uint8 { return g.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (g *Generic) SetSequenceNumber(seq uint32) error { return g.GetHeader().SetSequenceNumber(seq) }

// SetMessagePriority updates the common header safely on a nil receiver.
func (g *Generic) SetMessagePriority(priority uint8) error {
	return g.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (g *Generic) SetTEID(teid uint32) { g.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (g *Generic) SetPiggybacking(value uint8) { g.GetHeader().SetPiggybacking(value) }
