// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "fmt"

// ExtensionHeaderType definitions.
const (
	ExtHeaderTypeNoMoreExtensionHeaders                 uint8 = 0b00000000
	ExtHeaderTypeMBMSSupportIndication                  uint8 = 0b00000001
	ExtHeaderTypeMSInfoChangeReportingSupportIndication uint8 = 0b00000010
	ExtHeaderTypeLongPDCPPDUNumber                      uint8 = 0b00000011
	ExtHeaderTypeServiceClassIndicator                  uint8 = 0b00100000
	ExtHeaderTypeUDPPort                                uint8 = 0b01000000
	ExtHeaderTypeRANContainer                           uint8 = 0b10000001
	ExtHeaderTypeLongPDCPPDUNumberRequired              uint8 = 0b10000010
	ExtHeaderTypeXwRANContainer                         uint8 = 0b10000011
	ExtHeaderTypeNRRANContainer                         uint8 = 0b10000100
	ExtHeaderTypePDUSessionContainer                    uint8 = 0b10000101
	ExtHeaderTypePDCPPDUNumber                          uint8 = 0b11000000
	ExtHeaderTypeSuspendRequest                         uint8 = 0b11000001
	ExtHeaderTypeSuspendResponse                        uint8 = 0b11000010
)

// ExtensionHeader represents an Extension Header defined in §5.2, TS 29.281 and §6.1 TS 29.060.
type ExtensionHeader struct {
	Type     uint8 // this doesn't exist in the spec, but it's apparently helpful to have
	Length   uint8
	Content  []byte
	NextType uint8
}

// NewExtensionHeader creates a new ExtensionHeader.
//
// ExtensionHeader struct has its own type while it does not actually exist in the packet.
// Be sure to set an appropriate one - putting a wrong type may cause unexpected errors while
// using method depends on ExtensionHeader struct.
func NewExtensionHeader(typ uint8, content []byte, nextType uint8) *ExtensionHeader {
	eh := &ExtensionHeader{
		Type:     typ,
		Content:  content,
		NextType: nextType,
	}

	eh.SetLength()
	return eh
}

// Marshal returns the byte sequence generated from an ExtensionHeader.
func (e *ExtensionHeader) Marshal() ([]byte, error) {
	b := make([]byte, e.MarshalLen())
	if err := e.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (e *ExtensionHeader) MarshalTo(b []byte) error {
	return e.marshalTo(b, false)
}

func (e *ExtensionHeader) marshalTo(b []byte, preserving bool) error {
	if e == nil || e.MarshalLen() > 1020 {
		return ErrInvalidLength
	}
	if preserving {
		// TS 29.281 §5.2.1: replay requires complete four-octet units.
		if len(e.Content)+2 != e.MarshalLen() {
			return ErrInvalidLength
		}
	} else if err := e.validate(); err != nil {
		return err
	}
	n := e.MarshalLen()
	if len(b) < n {
		return ErrTooShortToMarshal
	}
	// TS 29.281 §5.2.1: four-octet units; all padding is sent as zero.
	clear(b[:n])
	b[0] = uint8(n / 4)
	copy(b[1:n-1], e.Content)
	b[n-1] = e.NextType
	if preserving {
		return nil
	}
	// TS 29.281 §§5.2.2.2A, 5.2.2.3: type-specific spare bits and octets.
	switch e.Type {
	case ExtHeaderTypeLongPDCPPDUNumber, ExtHeaderTypeLongPDCPPDUNumberRequired:
		b[1] &= 3
		clear(b[4 : n-1])
	case ExtHeaderTypeServiceClassIndicator:
		b[2] = 0
	}
	return nil
}

// ParseExtensionHeader decodes given byte sequence as a ExtensionHeader.
func ParseExtensionHeader(b []byte) (*ExtensionHeader, error) {
	e := &ExtensionHeader{}
	if err := e.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return e, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in ExtensionHeader.
func (e *ExtensionHeader) UnmarshalBinary(b []byte) error {
	if e == nil {
		return ErrInvalidLength
	}
	typ := e.Type
	*e = ExtensionHeader{Type: typ}
	l := len(b)
	if l < 2 {
		return ErrTooShortToParse
	}

	e.Length = b[0]
	if e.Length == 0 {
		return ErrInvalidLength
	}
	offset := int(e.Length)*4 - 1
	if l < offset+1 {
		return ErrTooShortToParse
	}

	e.Content = b[1:offset]
	e.NextType = b[offset]

	return nil
}

// ParseMultiExtensionHeaders parses given bytes sequence as multiple ExtensionHeaders.
func ParseMultiExtensionHeaders(b []byte) ([]*ExtensionHeader, error) {
	var ehs []*ExtensionHeader
	next := ExtHeaderTypeNoMoreExtensionHeaders
	for {
		eh, err := ParseExtensionHeader(b)
		if err != nil {
			return nil, err
		}
		eh.Type = next

		ehs = append(ehs, eh)
		b = b[eh.MarshalLen():]
		next = eh.NextType
		if next == ExtHeaderTypeNoMoreExtensionHeaders {
			break
		}
		continue
	}
	return ehs, nil
}

// MarshalLen returns the serial length of ExtensionHeader.
func (e *ExtensionHeader) MarshalLen() int {
	if e == nil {
		return 0
	}
	return pad4Len(len(e.Content) + 2)
}

// SetLength sets the length calculated from the length of contents to Length field.
func (e *ExtensionHeader) SetLength() {
	if e == nil {
		return
	}
	if e.MarshalLen() > 1020 {
		e.Length = 0
		return
	}
	e.Length = uint8(e.MarshalLen() / 4)
}

// String returns an ExtensionHeader fields in human readable format.
func (e *ExtensionHeader) String() string {
	return fmt.Sprintf("{Type: %#x, Length: %d, Content: %#x, NextType: %x}",
		e.Type,
		e.Length,
		e.Content,
		e.NextType,
	)
}

// IsComprehensionRequired reports whether the comprehension of the ExtensionHeader is
// required or not.
func (e *ExtensionHeader) IsComprehensionRequired() bool {
	return e != nil && e.Type>>7 == 1
}

func pad4Len(n int) int {
	return n + ((4 - n) & 0b11)
}

func (e *ExtensionHeader) validate() error {
	if e == nil || e.MarshalLen() > 1020 || e.Type == 0 {
		return ErrInvalidLength
	}
	n := len(e.Content)
	// TS 29.281 §§5.2.2.1–5.2.2.3: known extensions have fixed wire lengths.
	switch e.Type {
	case ExtHeaderTypeUDPPort, ExtHeaderTypePDCPPDUNumber:
		if n != 2 {
			return ErrInvalidLength
		}
	case ExtHeaderTypeLongPDCPPDUNumber, ExtHeaderTypeLongPDCPPDUNumberRequired:
		if n < 3 || n > 6 {
			return ErrInvalidLength
		}
	case ExtHeaderTypeServiceClassIndicator:
		if n < 1 || n > 2 {
			return ErrInvalidLength
		}
	case ExtHeaderTypePDUSessionContainer:
		return validatePDUSession(e.Content, false)
	}
	return nil
}

func (e *ExtensionHeader) validateReceived() error {
	if e != nil && e.Type == ExtHeaderTypePDUSessionContainer {
		return validatePDUSession(e.Content, true)
	}
	return e.validate()
}

func validatePDUSession(p []byte, received bool) error {
	// TS 38.415 §§5.5.2.1, 5.5.2.2: check every field selected by the presence flags.
	if len(p) < 2 {
		return ErrInvalidLength
	}
	need := 2
	switch p[0] >> 4 {
	case 0:
		if p[1]&0x80 != 0 {
			if len(p) < 3 {
				return ErrInvalidLength
			}
			need++
			if p[2]&2 != 0 {
				need += 3
			}
			if p[2]&1 != 0 {
				need += 2
			}
		}
		if p[0]&8 != 0 {
			need += 8
		}
		if p[0]&4 != 0 {
			need += 3
		}
		if p[0]&2 != 0 {
			need += 4
		}
	case 1:
		if p[0]&8 != 0 {
			need += 24
		}
		if p[0]&4 != 0 {
			need += 4
		}
		if p[0]&2 != 0 {
			need += 4
		}
		if p[0]&1 != 0 {
			need += 3
		}
		if p[1]&0x80 != 0 {
			need += 4
		}
		if p[1]&0x40 != 0 {
			// TS 38.415 §5.5.2.2: extension flag octets precede the selected new IEs.
			var flags []byte
			for {
				if need >= len(p) {
					return ErrInvalidLength
				}
				f := p[need]
				need++
				flags = append(flags, f)
				if f&0x80 == 0 {
					break
				}
			}
			if flags[0]&1 != 0 {
				need++
			}
			if flags[0]&2 != 0 {
				need += 2
			}
			if flags[0]&4 != 0 {
				need += 2
			}
			if flags[0]&8 != 0 {
				need += 4
			}
			if flags[0]&16 != 0 {
				need += 4
			}
		}
	default:
		return ErrInvalidMessageType
	}
	// TS 38.415 §5.5.1: receivers can remove additional Future Extension
	// fields. Keep send-side layout validation strict.
	if len(p) < need || (!received && len(p)-need > 3) {
		return ErrInvalidLength
	}
	return nil
}
