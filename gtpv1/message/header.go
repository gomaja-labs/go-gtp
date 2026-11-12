// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"encoding/binary"
	"fmt"
)

// Header is a GTPv1 common header.
type Header struct {
	err                     error
	Flags                   uint8
	Type                    uint8
	Length                  uint16
	TEID                    uint32
	SequenceNumber          uint16
	NPDUNumber              uint8
	NextExtensionHeaderType uint8
	ExtensionHeaders        []*ExtensionHeader
	Payload                 []byte
}

// NewHeader creates a new Header.
func NewHeader(flags, mtype uint8, teid uint32, seqnum uint16, payload []byte) *Header {
	h := &Header{
		Flags:                   flags,
		Type:                    mtype,
		TEID:                    teid,
		SequenceNumber:          seqnum,
		NextExtensionHeaderType: ExtHeaderTypeNoMoreExtensionHeaders,
		Payload:                 payload,
	}

	h.SetLength()
	return h
}

// NewHeaderWithNPDUNumber creates a new Header with NPDUNumber.
func NewHeaderWithNPDUNumber(flags, mtype uint8, teid uint32, seqnum uint16, npdu uint8, payload []byte) *Header {
	h := NewHeader(flags, mtype, teid, seqnum, payload)
	h.SetNPDUNumber(npdu)
	return h
}

// NewHeaderWithExtensionHeaders creates a new Header with ExtensionHeaders.
func NewHeaderWithExtensionHeaders(flags, mtype uint8, teid uint32, seqnum uint16, payload []byte, extHdrs ...*ExtensionHeader) *Header {
	h := NewHeader(flags, mtype, teid, seqnum, payload)
	h.err = h.AddExtensionHeaders(extHdrs...)
	return h
}

// NewHeaderFlags returns a Header Flag built by its components given as arguments.
func NewHeaderFlags(v, p, e, s, n int) uint8 {
	return uint8(
		((v & 0x7) << 5) | ((p & 0x1) << 4) | ((e & 0x1) << 2) | ((s & 0x1) << 1) | (n & 0x1),
	)
}

// Marshal returns the byte sequence generated from a Header.
func (h *Header) Marshal() ([]byte, error) {
	b := make([]byte, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (h *Header) MarshalTo(b []byte) error {
	return h.marshalTo(b, nil)
}

// marshalTo preserves received spare and unused fields only for Generic replay.
func (h *Header) marshalTo(b, received []byte) error {
	if h == nil {
		return ErrInvalidLength
	}
	if h.err != nil {
		return h.err
	}
	// TS 29.281 §§5.1, 5.2.1: representable length and a complete extension chain.
	if h.Flags>>5 != 1 || h.Flags&0x10 == 0 {
		return ErrInvalidMessageType
	}
	if h.MarshalLen()-8 > 65535 {
		return ErrInvalidLength
	}
	if err := h.validateChain(received); err != nil {
		return err
	}
	n := h.MarshalLen()
	if len(b) < n {
		return ErrTooShortToMarshal
	}
	clear(b[:n])
	b[0] = h.Flags & 0xf7
	if received != nil {
		b[0] = h.Flags
	}
	b[1] = h.Type
	binary.BigEndian.PutUint16(b[2:4], uint16(n-8))
	binary.BigEndian.PutUint32(b[4:8], h.TEID)
	off := 8
	// TS 29.281 §5.1: any E/S/PN flag includes all four optional octets; unused fields are zero.
	if h.Flags&7 != 0 {
		if len(received) >= 12 {
			copy(b[8:12], received[8:12])
		}
		if h.HasSequence() {
			binary.BigEndian.PutUint16(b[off:off+2], h.SequenceNumber)
		}
		if h.HasNPDUNumber() {
			b[off+2] = h.NPDUNumber
		}
		if h.HasExtensionHeader() {
			b[off+3] = h.NextExtensionHeaderType
		}
		off += 4
	}
	if h.HasExtensionHeader() {
		for _, eh := range h.ExtensionHeaders {
			if err := eh.marshalTo(b[off:], received != nil); err != nil {
				return err
			}
			off += eh.MarshalLen()
		}
	}
	copy(b[off:n], h.Payload)
	return nil
}

// ParseHeader decodes given byte sequence as a GTPv1 header.
func ParseHeader(b []byte) (*Header, error) {
	h := &Header{}
	if err := h.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return h, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in GTPv1 header.
func (h *Header) UnmarshalBinary(b []byte) error {
	return h.unmarshalBinary(b, true)
}

func (h *Header) unmarshalBinary(b []byte, contents bool) error {
	if h == nil {
		return ErrInvalidLength
	}
	*h = Header{}
	// TS 29.281 §5.1: version 1, PT=1 and Length excludes the first eight octets.
	if len(b) < 8 {
		return ErrTooShortToParse
	}
	if b[0]>>5 != 1 || b[0]&0x10 == 0 {
		return ErrInvalidMessageType
	}
	h.Flags = b[0]
	h.Type = b[1]
	h.Length = binary.BigEndian.Uint16(b[2:4])
	h.TEID = binary.BigEndian.Uint32(b[4:8])
	end := 8 + int(h.Length)
	if end > len(b) {
		return ErrTooShortToParse
	}
	if end != len(b) {
		return ErrInvalidLength
	}
	off := 8
	if h.Flags&7 != 0 {
		if end-off < 4 {
			return ErrTooShortToParse
		}
		if h.HasSequence() {
			h.SequenceNumber = binary.BigEndian.Uint16(b[off : off+2])
		}
		if h.HasNPDUNumber() {
			h.NPDUNumber = b[off+2]
		}
		if h.HasExtensionHeader() {
			h.NextExtensionHeaderType = b[off+3]
		}
		off += 4
	}
	// TS 29.281 §5.2.1: next type zero terminates the chain, including an empty chain.
	next := h.NextExtensionHeaderType
	for next != 0 {
		eh := &ExtensionHeader{Type: next}
		if err := eh.UnmarshalBinary(b[off:end]); err != nil {
			return err
		}
		if contents {
			if err := eh.validateReceived(); err != nil {
				return err
			}
		}
		h.ExtensionHeaders = append(h.ExtensionHeaders, eh)
		off += eh.MarshalLen()
		next = eh.NextType
	}
	if off < end {
		h.Payload = b[off:end]
	}
	return nil
}

// SetTEID sets the always-present TEID (TS 29.281 §5.1).
func (h *Header) SetTEID(teid uint32) {
	if h == nil {
		return
	}
	h.TEID = teid
}

// Sequence returns SequenceNumber in uint16.
func (h *Header) Sequence() uint16 {
	if h == nil {
		return 0
	}
	return h.SequenceNumber
}

// HasSequence reports whether a Header has SequenceNumber by checking the flag.
func (h *Header) HasSequence() bool {
	if h == nil {
		return false
	}
	return ((int(h.Flags) >> 1) & 0x1) == 1
}

// SetSequenceNumber sets the SequenceNumber in Header.
func (h *Header) SetSequenceNumber(seq uint16) {
	if h == nil {
		return
	}
	h.Flags |= 0x02
	h.SequenceNumber = seq
	h.SetLength()
}

// WithSequenceNumber returns the Header with SequenceNumber added.
func (h *Header) WithSequenceNumber(seq uint16) *Header {
	h.SetSequenceNumber(seq)
	return h
}

// HasNPDUNumber reports whether a Header has N-PDU Number by checking the flag.
func (h *Header) HasNPDUNumber() bool {
	if h == nil {
		return false
	}
	return (int(h.Flags) & 0x1) == 1
}

// SetNPDUNumber sets the NPDUNumber in Header.
func (h *Header) SetNPDUNumber(npdu uint8) {
	if h == nil {
		return
	}
	h.Flags |= 0x01
	h.NPDUNumber = npdu
	h.SetLength()
}

// WithNPDUNumber returns the Header with NPDUNumber added.
func (h *Header) WithNPDUNumber(npdu uint8) *Header {
	h.SetNPDUNumber(npdu)
	return h
}

// HasExtensionHeader reports whether a Header has extension header by checking the flag.
func (h *Header) HasExtensionHeader() bool {
	if h == nil {
		return false
	}
	return ((int(h.Flags) >> 2) & 0x1) == 1
}

// SetNextExtensionHeaderType sets the ExtensionHeaderType in Header.
func (h *Header) SetNextExtensionHeaderType(exhType uint8) {
	if h == nil {
		return
	}
	h.Flags |= 0x04
	h.NextExtensionHeaderType = exhType
}

// WithExtensionHeaders returns the Header with ExtensionHeaders added.
func (h *Header) WithExtensionHeaders(extHdrs ...*ExtensionHeader) *Header {
	h.err = h.AddExtensionHeaders(extHdrs...)
	return h
}

// MarshalLen returns the serial length of Header.
func (h *Header) MarshalLen() int {
	if h == nil {
		return 0
	}
	l := len(h.Payload) + 8
	if h.HasSequence() || h.HasNPDUNumber() || h.HasExtensionHeader() {
		l += 4
	}

	if h.HasExtensionHeader() {
		for _, eh := range h.ExtensionHeaders {
			l += eh.MarshalLen()
		}
	}

	return l
}

// SetLength sets the length in Length field.
func (h *Header) SetLength() {
	if h == nil {
		return
	}
	h.Length = uint16(h.MarshalLen() - 8)
}

// Version returns GTP version in int.
func (h *Header) Version() int {
	return 1
}

// MessageType returns the type of message.
func (h *Header) MessageType() uint8 {
	if h == nil {
		return 0
	}
	return h.Type
}

// String returns the GTPv1 header values in human readable format.
func (h *Header) String() string {
	if h == nil {
		return "nil"
	}
	return fmt.Sprintf(
		"{Flags: %#x, Type: %#x, Length: %d, TEID: %#08x, SequenceNumber: %#04x, NPDUNumber: %#x, "+
			"ExtensionHeaderType: %#x, ExtensionHeaders: %v, Payload: %#v}",
		h.Flags,
		h.Type,
		h.Length,
		h.TEID,
		h.SequenceNumber,
		h.NPDUNumber,
		h.NextExtensionHeaderType,
		h.ExtensionHeaders,
		h.Payload,
	)
}

// AddExtensionHeaders adds ExtensionHeader(s) to Header.
//
// This function validates if the next extension header type matches the actual one for safety.
// To create arbitrary(possibly malformed) Header, access ExtensionHeaders field on your own.
func (h *Header) AddExtensionHeaders(extHdrs ...*ExtensionHeader) error {
	if h == nil {
		return ErrInvalidLength
	}
	if len(extHdrs) == 0 {
		return nil
	}
	// TS 29.281 §5.2.1: validate before linking the former tail to the new first type.
	for j, eh := range extHdrs {
		if eh == nil || eh.Type == 0 {
			return ErrInvalidLength
		}
		if err := eh.validate(); err != nil {
			return err
		}
		next := uint8(0)
		if j+1 < len(extHdrs) {
			if extHdrs[j+1] == nil {
				return ErrInvalidLength
			}
			next = extHdrs[j+1].Type
		}
		if eh.NextType != next {
			return ErrInvalidLength
		}
	}
	if err := h.validateChain(nil); err != nil {
		return err
	}
	combined := append([]*ExtensionHeader(nil), h.ExtensionHeaders...)
	if len(combined) > 0 {
		tail := *combined[len(combined)-1]
		tail.NextType = extHdrs[0].Type
		combined[len(combined)-1] = &tail
	}
	combined = append(combined, extHdrs...)
	h.Flags |= 4
	h.NextExtensionHeaderType = combined[0].Type
	h.ExtensionHeaders = combined
	h.SetLength()
	return nil
}

func (h *Header) validateChain(received []byte) error {
	next := h.NextExtensionHeaderType
	if !h.HasExtensionHeader() {
		if len(h.ExtensionHeaders) > 0 {
			return ErrInvalidLength
		}
		return nil
	}
	for _, eh := range h.ExtensionHeaders {
		if eh == nil || eh.Type == 0 || eh.Type != next {
			return ErrInvalidLength
		}
		if eh.MarshalLen() > 1020 {
			return ErrInvalidLength
		}
		if received == nil {
			if err := eh.validate(); err != nil {
				return err
			}
		}
		if received != nil && len(eh.Content)+2 != eh.MarshalLen() {
			return ErrInvalidLength
		}
		next = eh.NextType
	}
	if next != 0 {
		return ErrInvalidLength
	}
	return nil
}

// GetHeader returns the common header, including transport extension headers.
func (h *Header) GetHeader() *Header { return h }
