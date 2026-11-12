// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"fmt"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// Generic holds a header and structurally framed IEs in received order.
// It does not enforce mandatory presence, IE contents, or IE order.
// Decoded Generics encode every IE, including appended IEs, with MarshalToLossless.
// NewGeneric uses normal IE and header encoding checks.
type Generic struct {
	*Header
	IEs            []*ie.IE
	receivedHeader []byte
}

// NewGeneric creates a new GTPv1 Generic.
func NewGeneric(msgType uint8, teid uint32, seq uint16, ie ...*ie.IE) *Generic {
	g := &Generic{
		Header: NewHeader(0x32, msgType, teid, seq, nil),
		IEs:    ie,
	}

	g.SetLength()
	return g
}

// Marshal returns the byte sequence generated from a Generic.
func (g *Generic) Marshal() ([]byte, error) {
	if g == nil || g.Header == nil {
		return nil, ErrInvalidLength
	}
	b := make([]byte, g.MarshalLen())
	if err := g.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (g *Generic) MarshalTo(b []byte) error {
	if g == nil || g.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{}
	ies = append(ies, g.IEs...)
	return marshalGenericIEs(g.Header, b, ies, g.receivedHeader)
}

// ParseGeneric decodes a given byte sequence as a Generic.
func ParseGeneric(b []byte) (*Generic, error) {
	g := &Generic{}
	if err := g.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return g, nil
}

// UnmarshalBinary decodes a given byte sequence as a Generic.
func (g *Generic) UnmarshalBinary(b []byte) error {
	if g == nil {
		return ErrInvalidLength
	}
	*g = Generic{}
	var err error
	g.Header = &Header{}
	err = g.Header.unmarshalBinary(b, false)
	if err != nil {
		return err
	}

	g.IEs, err = ie.ParseMultiIEs(g.Header.Payload)
	if err != nil {
		return err
	}
	g.receivedHeader = append([]byte(nil), b[:len(b)-len(g.Header.Payload)]...)
	return nil
}

// MarshalLen returns the serial length of Data.
func (g *Generic) MarshalLen() int {
	if g == nil || g.Header == nil {
		return 0
	}
	ies := []*ie.IE{}
	ies = append(ies, g.IEs...)
	return messageLen(g.Header, ies)
}

// SetLength sets the length in Length field.
func (g *Generic) SetLength() {
	if g == nil || g.Header == nil {
		return
	}
	g.Header.Length = uint16(g.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (g *Generic) MessageTypeName() string {
	if g == nil || g.Header == nil {
		return "Unknown"
	}
	return fmt.Sprintf("Unknown (%d)", g.Type)
}

// TEID returns the TEID in human-readable string.
func (g *Generic) TEID() uint32 {
	if g == nil || g.Header == nil {
		return 0
	}
	return g.Header.TEID
}

// AddIE appends IEs and updates Length.
func (g *Generic) AddIE(ie ...*ie.IE) {
	g.IEs = append(g.IEs, ie...)
	g.SetLength()
}
