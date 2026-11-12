// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"encoding/binary"
	"fmt"
)

// Header is a GTPv0 common header.
type Header struct {
	Flags          uint8
	Type           uint8
	Length         uint16
	SequenceNumber uint16
	SndcpNumber    uint8
	FlowLabel      uint16
	TID            uint64
	Payload        []byte
}

// NewHeader creates a new Header.
func NewHeader(flags, mtype uint8, seq, label uint16, tid uint64, payload []byte) *Header {
	h := &Header{
		Flags:          flags,
		Type:           mtype,
		SequenceNumber: seq,
		FlowLabel:      label,
		SndcpNumber:    0xff,
		TID:            tid,
		Payload:        payload,
	}
	h.SetLength()

	return h
}

// HeaderFlags returns a Header Flag built by its components given as arguments.
func HeaderFlags(v, p, s int) uint8 {
	return uint8(
		((v & 0x7) << 5) | ((p & 0x1) << 4) | (s & 0x1) | 0x0e,
	)
}

// Marshal returns the byte sequence generated from an IE instance.
func (h *Header) Marshal() ([]byte, error) {
	b := make([]byte, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (h *Header) MarshalTo(b []byte) error {
	if h == nil || len(h.Payload) > 65535 {
		return ErrInvalidLength
	}
	if h.Flags>>5 != 0 || h.Flags&0x10 == 0 {
		return ErrInvalidLength
	}
	if len(b) < h.MarshalLen() {
		return ErrTooShortToMarshal
	}

	b[0] = h.Flags | 0x0e
	b[1] = h.Type
	binary.BigEndian.PutUint16(b[2:4], uint16(len(h.Payload)))
	binary.BigEndian.PutUint16(b[4:6], h.SequenceNumber)
	binary.BigEndian.PutUint16(b[6:8], h.FlowLabel)
	binary.BigEndian.PutUint32(b[8:12], uint32(int(h.SndcpNumber)<<24|0xffffff))
	binary.BigEndian.PutUint64(b[12:20], h.TID)
	// two bytes of padding before payload.
	copy(b[20:h.MarshalLen()], h.Payload)
	return nil
}

// ParseHeader Parses given byte sequence as a GTPv1 header.
func ParseHeader(b []byte) (*Header, error) {
	h := &Header{}
	if err := h.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return h, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in GTPv1 header.
func (h *Header) UnmarshalBinary(b []byte) error {
	if h == nil {
		return ErrInvalidLength
	}
	*h = Header{}
	// GSM 09.60 §6 Figure 2: fixed twenty-octet header; SNDCP occupies octet 9.
	if len(b) < 20 {
		return ErrTooShortToParse
	}
	if b[0]>>5 != 0 || b[0]&0x10 == 0 {
		return ErrInvalidLength
	}
	h.Flags = b[0]
	h.Type = b[1]
	h.Length = binary.BigEndian.Uint16(b[2:4])
	if 20+int(h.Length) != len(b) {
		return ErrInvalidLength
	}
	h.SequenceNumber = binary.BigEndian.Uint16(b[4:6])
	h.FlowLabel = binary.BigEndian.Uint16(b[6:8])
	h.SndcpNumber = b[8]
	h.TID = binary.BigEndian.Uint64(b[12:20])
	if h.Length > 0 {
		h.Payload = b[20 : 20+int(h.Length)]
	}
	return nil
}

// MarshalLen returns the serial length of Header.
func (h *Header) MarshalLen() int {
	if h == nil {
		return 0
	}
	return 20 + len(h.Payload)
}

// SetLength sets the length in Length field.
func (h *Header) SetLength() {
	if h == nil {
		return
	}
	h.Length = uint16(len(h.Payload))
}

// String returns the GTPv1 header values in human readable format.
func (h *Header) String() string {
	if h == nil {
		return "nil"
	}
	return fmt.Sprintf("{Flags: %#x, Type: %#x, Length: %d, SequenceNumber: %#04x, FlowLabel: %#04x, SndcpNumber: %#02x, TID: %#016x, Payload: %#v}",
		h.Flags,
		h.Type,
		h.Length,
		h.SequenceNumber,
		h.FlowLabel,
		h.SndcpNumber,
		h.TID,
		h.Payload,
	)
}

// tid returns the tid in human-readable string.
func (h *Header) tid() string {
	if h == nil {
		return ""
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, h.TID)

	// GSM 09.60 §7.2: preserve the NSAPI nibble alongside the IMSI digits.
	for n := range b {
		b[n] = b[n]>>4 | b[n]<<4
	}
	return fmt.Sprintf("%x", b)
}

// Version returns the GTP version.
func (h *Header) Version() int {
	return 0
}

// MessageType returns the type of message.
func (h *Header) MessageType() uint8 {
	if h == nil {
		return 0
	}
	return h.Type
}
