// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// TS 24.008 §10.5.6.3 and TS 29.274 §§8.13, 8.94: NAS octets 3-z.
const maxPCOPayloadLen = 251

// ConfigurationProtocol definitions.
const (
	ConfigurationProtocolPPPForUseWithIPPDPTypeOrIPPDNType uint8 = 0b000
)

// NewProtocolConfigurationOptions creates a new ProtocolConfigurationOptions IE.
func NewProtocolConfigurationOptions(proto uint8, options ...*PCOContainer) *IE {
	v := NewProtocolConfigurationOptionsFields(proto, options...)
	b, err := v.Marshal()
	if err != nil {
		return newErrorIE(ProtocolConfigurationOptions, err)
	}

	if len(b) > maxPCOPayloadLen {
		return newErrorIE(ProtocolConfigurationOptions, ErrInvalidLength)
	}
	return New(ProtocolConfigurationOptions, 0x00, b)
}

// ProtocolConfigurationOptions returns ProtocolConfigurationOptions in
// ProtocolConfigurationOptionsFields type if the type of IE matches.
func (i *IE) ProtocolConfigurationOptions() (*ProtocolConfigurationOptionsFields, error) {
	if i != nil && i.Type == ProtocolConfigurationOptions && len(i.Payload) > maxPCOPayloadLen {
		return nil, ErrInvalidLength
	}
	if err := i.Err(); err != nil {
		return nil, err
	}
	switch i.Type {
	case ProtocolConfigurationOptions:
		if len(i.Payload) < 1 {
			return nil, io.ErrUnexpectedEOF
		}

		return ParseProtocolConfigurationOptionsFields(i.Payload)
	case BearerContext:
		ies, err := i.BearerContext()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve ProtocolConfigurationOptions: %w", err)
		}

		// TS 29.274 §7.7.9: match the child type and instance together.
		for _, child := range ies {
			if child != nil && child.Type == ProtocolConfigurationOptions && child.Instance() == 0 {
				return child.ProtocolConfigurationOptions()
			}
		}
		return nil, ErrIENotFound
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// MustProtocolConfigurationOptions returns ProtocolConfigurationOptions in *ProtocolConfigurationOptionsFields, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustProtocolConfigurationOptions() *ProtocolConfigurationOptionsFields {
	v, _ := i.ProtocolConfigurationOptions()
	return v
}

// ProtocolConfigurationOptionsFields is a set of fields in ProtocolConfigurationOptions IE.
type ProtocolConfigurationOptionsFields struct {
	Extension             uint8 // bit 8 of octet 1
	ConfigurationProtocol uint8 // bit 1-3 of octet 1
	ProtocolOrContainers  []*PCOContainer
}

// NewProtocolConfigurationOptionsFields creates a new ProtocolConfigurationOptionsFields.
func NewProtocolConfigurationOptionsFields(proto uint8, opts ...*PCOContainer) *ProtocolConfigurationOptionsFields {
	f := &ProtocolConfigurationOptionsFields{ConfigurationProtocol: proto}
	f.ProtocolOrContainers = append(f.ProtocolOrContainers, opts...)

	return f
}

// Marshal serializes ProtocolConfigurationOptionsFields.
func (f *ProtocolConfigurationOptionsFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo serializes ProtocolConfigurationOptionsFields.
func (f *ProtocolConfigurationOptionsFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.MarshalLen() > maxPCOPayloadLen {
		return ErrInvalidLength
	}
	for _, c := range f.ProtocolOrContainers {
		if c == nil {
			return ErrMalformed
		}
	}
	if f.ConfigurationProtocol != 0 {
		return ErrMalformed
	}
	if len(b) < f.MarshalLen() {
		return ErrTooShortToMarshal
	}
	b[0] = (f.ConfigurationProtocol & 0x07) | 0x80
	offset := 1
	for _, opt := range f.ProtocolOrContainers {
		if err := opt.MarshalTo(b[offset:]); err != nil {
			return err
		}
		offset += opt.MarshalLen()
	}

	return nil
}

// ParseProtocolConfigurationOptionsFields decodes ProtocolConfigurationOptionsFields.
func ParseProtocolConfigurationOptionsFields(b []byte) (*ProtocolConfigurationOptionsFields, error) {
	f := &ProtocolConfigurationOptionsFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return f, nil
}

// UnmarshalBinary decodes given bytes into ProtocolConfigurationOptionsFields.
func (f *ProtocolConfigurationOptionsFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = ProtocolConfigurationOptionsFields{}
	if len(b) > maxPCOPayloadLen {
		return ErrInvalidLength
	}
	if len(b) < 1 {
		return ErrTooShortToParse
	}

	f.Extension = (b[0] >> 7) & 0x01
	// TS 24.008 Table 10.5.154: other values are interpreted as PPP.
	f.ConfigurationProtocol = 0

	offset := 1
	for {
		if offset >= len(b) {
			return nil
		}
		opt, err := ParsePCOContainer(b[offset:])
		if err != nil {
			return err
		}
		f.ProtocolOrContainers = append(f.ProtocolOrContainers, opt)
		offset += opt.MarshalLen()
	}
}

// MarshalLen returns the serial length of ProtocolConfigurationOptionsFields in int.
func (f *ProtocolConfigurationOptionsFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	l := 1
	for _, opt := range f.ProtocolOrContainers {
		if opt != nil {
			l += opt.MarshalLen()
		}
	}

	return l
}

// ProtocolIdentifier definitions.
//
// [Table 10.5.154/3GPP TS 24.008]
//
// At least the following protocol identifiers (as defined in RFC 3232, Description) shall be
// supported in this version of the protocol:
// - C021H (LCP);
// - C023H (PAP);
// - C223H (CHAP); and
// - 8021H (IPCP).
const (
	PCOProtocolIdentifierLCP  uint16 = 0xc021
	PCOProtocolIdentifierPAP  uint16 = 0xc023
	PCOProtocolIdentifierCHAP uint16 = 0xc223
	PCOProtocolIdentifierIPCP uint16 = 0x8021
)

// PCOContainer is either of a Configuration protocol option or Additional parameters in PCO,
// which are not distinguishable without meta information(link direction) but fortunately
// the format is the same.
type PCOContainer struct {
	ID       uint16
	Length   uint8
	Contents []byte
}

// NewPCOContainer creates a new PCOContainer.
func NewPCOContainer(pid uint16, contents []byte) *PCOContainer {
	c := &PCOContainer{
		ID:       pid,
		Length:   uint8(len(contents)),
		Contents: contents,
	}
	return c
}

// Marshal serializes PCOContainer.
func (c *PCOContainer) Marshal() ([]byte, error) {
	if c == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo serializes PCOContainer.
func (c *PCOContainer) MarshalTo(b []byte) error {
	if c == nil {
		return ErrIENotFound
	}
	if len(c.Contents) > 255 || int(c.Length) != len(c.Contents) {
		return ErrInvalidLength
	}
	if len(b) < c.MarshalLen() {
		return ErrTooShortToMarshal
	}
	binary.BigEndian.PutUint16(b[0:2], c.ID)
	b[2] = c.Length
	if c.Length != 0 {
		copy(b[3:], c.Contents)
	}

	return nil
}

// ParsePCOContainer decodes PCOContainer.
func ParsePCOContainer(b []byte) (*PCOContainer, error) {
	c := &PCOContainer{}
	if err := c.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return c, nil
}

// UnmarshalBinary decodes given bytes into PCOContainer.
func (c *PCOContainer) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrIENotFound
	}
	*c = PCOContainer{}
	l := len(b)
	if l < 3 {
		return ErrTooShortToParse
	}

	c.ID = binary.BigEndian.Uint16(b[0:2])
	c.Length = b[2]

	if l < 3+int(c.Length) {
		return ErrInvalidLength
	}
	if c.Length != 0 {
		c.Contents = make([]byte, c.Length)
		copy(c.Contents, b[3:3+int(c.Length)])
	}

	return nil
}

// MarshalLen returns the serial length of PCOContainer in int.
func (c *PCOContainer) MarshalLen() int {
	if c == nil {
		return 0
	}
	return 3 + len(c.Contents)
}

// NewAdditionalProtocolConfigurationOptions creates an APCO IE (TS 29.274 §8.94).
func NewAdditionalProtocolConfigurationOptions(proto uint8, options ...*PCOContainer) *IE {
	// TS 29.274 §8.94: normal APCO uses PCO contents.
	b, err := NewProtocolConfigurationOptionsFields(proto, options...).Marshal()
	if err != nil {
		return newErrorIE(AdditionalProtocolConfigurationOptions, err)
	}
	if len(b) > maxPCOPayloadLen {
		return newErrorIE(AdditionalProtocolConfigurationOptions, ErrInvalidLength)
	}
	return New(AdditionalProtocolConfigurationOptions, 0, b)
}

// AdditionalProtocolConfigurationOptions decodes an APCO IE.
func (i *IE) AdditionalProtocolConfigurationOptions() (*ProtocolConfigurationOptionsFields, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	if i.Type != AdditionalProtocolConfigurationOptions {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) > maxPCOPayloadLen {
		return nil, ErrInvalidLength
	}
	return ParseProtocolConfigurationOptionsFields(i.Payload)
}

// ExtendedPCOContainer is a container in ePCO. Its length width depends on ID and direction.
type ExtendedPCOContainer struct {
	ID       uint16
	Contents []byte
}

// ExtendedProtocolConfigurationOptionsFields contains ePCO configuration and containers.
type ExtendedProtocolConfigurationOptionsFields struct {
	Direction             PCODirection
	ConfigurationProtocol uint8
	ProtocolOrContainers  []*ExtendedPCOContainer
}

// NewExtendedProtocolConfigurationOptions creates an ePCO IE (TS 29.274 §8.128).
func NewExtendedProtocolConfigurationOptions(proto uint8, direction PCODirection, options ...*ExtendedPCOContainer) *IE {
	f := &ExtendedProtocolConfigurationOptionsFields{Direction: direction, ConfigurationProtocol: proto, ProtocolOrContainers: options}
	b, err := f.Marshal()
	if err != nil {
		return newErrorIE(ExtendedProtocolConfigurationOptions, err)
	}
	return New(ExtendedProtocolConfigurationOptions, 0, b)
}

// ExtendedProtocolConfigurationOptions decodes an ePCO IE.
func (i *IE) ExtendedProtocolConfigurationOptions(direction PCODirection) (*ExtendedProtocolConfigurationOptionsFields, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	if i.Type != ExtendedProtocolConfigurationOptions {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	f := &ExtendedProtocolConfigurationOptionsFields{Direction: direction}
	if err := f.UnmarshalBinary(i.Payload); err != nil {
		return nil, err
	}
	return f, nil
}

// Marshal serializes ePCO fields.
func (f *ExtendedProtocolConfigurationOptionsFields) Marshal() ([]byte, error) {
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ePCO fields into b.
func (f *ExtendedProtocolConfigurationOptionsFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if len(b) < f.MarshalLen() {
		return ErrTooShortToMarshal
	}
	if f.ConfigurationProtocol != 0 || f.Direction > PCONetworkToMS || f.MarshalLen() > 65535 {
		return ErrMalformed
	}
	b[0] = 0x80 | f.ConfigurationProtocol
	offset := 1
	// TS 24.008 §10.5.6.3A reuses §10.5.6.3, Figure 10.5.136.
	for _, c := range f.ProtocolOrContainers {
		if c == nil {
			return ErrMalformed
		}
		width := extendedPCOLengthWidth(c.ID, f.Direction)
		if len(c.Contents) > 65535 || width == 1 && len(c.Contents) > 255 {
			return ErrInvalidLength
		}
		binary.BigEndian.PutUint16(b[offset:offset+2], c.ID)
		if width == 1 {
			b[offset+2] = uint8(len(c.Contents))
		} else {
			binary.BigEndian.PutUint16(b[offset+2:offset+4], uint16(len(c.Contents)))
		}
		copy(b[offset+2+width:], c.Contents)
		offset += 2 + width + len(c.Contents)
	}

	return nil
}

// UnmarshalBinary replaces the ePCO fields.
func (f *ExtendedProtocolConfigurationOptionsFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	direction := f.Direction
	*f = ExtendedProtocolConfigurationOptionsFields{Direction: direction}
	if direction > PCONetworkToMS {
		return ErrMalformed
	}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	// TS 24.008 Table 10.5.154: all configuration protocol values mean PPP.
	f.ConfigurationProtocol = 0
	for offset := 1; offset < len(b); {
		if len(b)-offset < 3 {
			return io.ErrUnexpectedEOF
		}
		id := binary.BigEndian.Uint16(b[offset : offset+2])
		width := extendedPCOLengthWidth(id, direction)
		if len(b)-offset < 2+width {
			return io.ErrUnexpectedEOF
		}
		n := int(b[offset+2])
		if width == 2 {
			n = int(binary.BigEndian.Uint16(b[offset+2 : offset+4]))
		}
		header := 2 + width
		if n > len(b)-offset-header {
			return io.ErrUnexpectedEOF
		}
		f.ProtocolOrContainers = append(f.ProtocolOrContainers, &ExtendedPCOContainer{ID: id, Contents: append([]byte(nil), b[offset+header:offset+header+n]...)})
		offset += header + n
	}

	return nil
}

// MarshalLen returns the serialized ePCO length.
func (f *ExtendedProtocolConfigurationOptionsFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	for _, c := range f.ProtocolOrContainers {
		if c != nil {
			n += 2 + extendedPCOLengthWidth(c.ID, f.Direction) + len(c.Contents)
		}
	}
	return n
}

// PCODirection identifies the origin of the contained NAS configuration.
type PCODirection uint8

const (
	PCOMSToNetwork PCODirection = iota
	PCONetworkToMS
)

// TS 24.008 §10.5.6.3, Figure 10.5.136 note: only these IDs use two octets.
func extendedPCOLengthWidth(id uint16, direction PCODirection) int {
	switch id {
	case 0x41, 0x51, 0x56, 0x61:
		return 2
	case 0x23, 0x24, 0x30, 0x31, 0x32:
		if direction == PCONetworkToMS {
			return 2
		}
	}
	return 1
}

// NewLargeAdditionalProtocolConfigurationOptions creates APCO after Large APCO negotiation.
func NewLargeAdditionalProtocolConfigurationOptions(proto uint8, direction PCODirection, options ...*ExtendedPCOContainer) *IE {
	// TS 29.274 §8.94: use ePCO contents when LAPCOSI has been negotiated.
	i := NewExtendedProtocolConfigurationOptions(proto, direction, options...)
	i.Type = AdditionalProtocolConfigurationOptions
	return i
}

// LargeAdditionalProtocolConfigurationOptions decodes negotiated Large APCO contents.
func (i *IE) LargeAdditionalProtocolConfigurationOptions(direction PCODirection) (*ExtendedProtocolConfigurationOptionsFields, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	if i.Type != AdditionalProtocolConfigurationOptions {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	f := &ExtendedProtocolConfigurationOptionsFields{Direction: direction}
	if err := f.UnmarshalBinary(i.Payload); err != nil {
		return nil, err
	}
	return f, nil
}
