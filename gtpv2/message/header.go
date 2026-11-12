// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"encoding/binary"
	"fmt"

	"github.com/gomaja-labs/go-gtp/utils"
)

const (
	fixedHeaderSize  = 4
	seqSpareSize     = 4
	teidSize         = 4
	noTEIDHeaderSize = fixedHeaderSize + seqSpareSize
	teidHeaderSize   = noTEIDHeaderSize + teidSize
)

// Header is a GTPv2 common header.
type Header struct {
	Flags          uint8
	Type           uint8
	Length         uint16
	TEID           uint32
	SequenceNumber uint32
	Spare          uint8
	Payload        []byte
	// Piggybacked is the optional second message in the datagram. Marshal sets
	// P from this field; Length describes only the first message (TS 29.274 §5.5.1).
	// Process this message after the enclosing message (TS 29.274 §7.7.0).
	Piggybacked Message
	// PiggybackError reports a secondary-message decode failure. The primary
	// message remains usable (TS 29.274 §7.7.0). Typed encoders reject this
	// state; Generic retains the secondary bytes for lossless replay.
	PiggybackError error
	err            error
	// Generic decoding retains received spare bits and an undecodable secondary.
	replay       bool
	rawFlags     uint8
	rawSpare     uint8
	piggybackRaw []byte
}

// NewHeader creates a new Header. Invalid values are rejected by MarshalTo.
func NewHeader(flags, mtype uint8, teid, seqnum uint32, data []byte) *Header {
	h := &Header{Flags: flags, Type: mtype, TEID: teid, SequenceNumber: seqnum, Payload: data}
	if !h.HasTEID() {
		h.TEID = 0
	}
	h.SetLength()
	return h
}

// NewHeaderFlags returns a Header Flag built by its components given as arguments.
func NewHeaderFlags(v, p, t int) uint8 {
	return uint8(((v & 7) << 5) | ((p & 1) << 4) | ((t & 1) << 3))
}

// GetHeader returns the common header, including the optional piggybacked message.
func (h *Header) GetHeader() *Header { return h }

func (h *Header) headerLen() int {
	if h.HasTEID() {
		return teidHeaderSize
	}
	return noTEIDHeaderSize
}

func (h *Header) piggybackLen() int {
	if h != nil && len(h.piggybackRaw) != 0 && h.PiggybackError != nil {
		return len(h.piggybackRaw)
	}
	if h == nil || h.Piggybacked == nil {
		return 0
	}
	// Reject chains and cycles in validate before asking the child for its size.
	if p := h.Piggybacked.GetHeader(); p == nil || p == h || p.Piggybacked != nil {
		return 0
	}
	return h.Piggybacked.MarshalLen()
}

func requiresTEID(typ uint8) bool {
	return typ != MsgTypeEchoRequest && typ != MsgTypeEchoResponse && typ != MsgTypeVersionNotSupportedIndication
}

func validPiggyback(first, second uint8) bool {
	// TS 29.274 §5.5.1: only these two EPC message pairs may share a datagram.
	return first == MsgTypeCreateSessionResponse && second == MsgTypeCreateBearerRequest ||
		first == MsgTypeCreateBearerResponse && second == MsgTypeModifyBearerRequest
}

func (h *Header) validate(expected uint8, size int) error {
	if h == nil {
		return ErrNilMessage
	}
	if h.PiggybackError != nil && len(h.piggybackRaw) == 0 {
		return h.PiggybackError
	}
	if h.err != nil {
		return h.err
	}
	if h.Version() != 2 {
		return &UnsupportedVersionError{Version: uint8(h.Version())}
	}
	if h.Type != expected {
		return &UnexpectedTypeError{Expected: expected, Actual: h.Type}
	}
	// TS 29.274 §§5.3, 5.4, 5.5.1: only Echo and VNSI omit TEID.
	if h.HasTEID() != requiresTEID(h.Type) {
		return ErrInvalidTEIDFlag
	}
	if h.SequenceNumber > 0xffffff {
		return ErrInvalidSequence
	}
	if h.Piggybacked != nil {
		p := h.Piggybacked.GetHeader()
		if p == nil || p == h || p.Piggybacked != nil || !validPiggyback(h.Type, p.Type) {
			return ErrInvalidPiggyback
		}
	} else if h.IsPiggybacking() && (len(h.piggybackRaw) == 0 || h.PiggybackError == nil) {
		return ErrInvalidPiggyback
	}
	if size-h.piggybackLen()-fixedHeaderSize > 0xffff {
		return ErrInvalidLength
	}
	return nil
}

// Marshal returns the byte sequence generated from a Header instance.
func (h *Header) Marshal() ([]byte, error) {
	if h == nil {
		return nil, ErrNilMessage
	}
	b := make([]byte, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo writes the whole datagram, returning an error for invalid fields or
// a short destination. Spare bits are zero (TS 29.274 §§5.3–5.5.1), except
// when replaying a header received through ParseGeneric.
func (h *Header) MarshalTo(b []byte) error {
	if h == nil {
		return ErrNilMessage
	}
	if err := h.validate(h.Type, h.MarshalLen()); err != nil {
		return err
	}
	if len(b) < h.MarshalLen() {
		return ErrTooShortToMarshal
	}
	h.SetLength()
	flags := h.Flags & 0xec
	if !h.HasTEID() {
		flags &^= 4
	}
	if h.Piggybacked != nil {
		flags |= 0x10
	}
	flags |= h.rawFlags
	if len(h.piggybackRaw) != 0 && h.PiggybackError != nil {
		flags |= 0x10
	}
	b[0], b[1] = flags, h.Type
	binary.BigEndian.PutUint16(b[2:4], h.Length)
	offset := fixedHeaderSize
	if h.HasTEID() {
		binary.BigEndian.PutUint32(b[offset:offset+4], h.TEID)
		offset += teidSize
	}
	b[offset], b[offset+1], b[offset+2] = byte(h.SequenceNumber>>16), byte(h.SequenceNumber>>8), byte(h.SequenceNumber)
	b[offset+3] = 0
	if h.HasTEID() && h.HasMessagePriority() {
		b[offset+3] = h.Spare & 0xf0
	}
	b[offset+3] |= h.rawSpare
	end := h.headerLen() + len(h.Payload)
	copy(b[offset+4:end], h.Payload)
	if len(h.piggybackRaw) != 0 && h.PiggybackError != nil {
		copy(b[end:], h.piggybackRaw)
		return nil
	}
	if h.Piggybacked != nil {
		return h.Piggybacked.MarshalTo(b[end:])
	}
	return nil
}

// ParseHeader decodes a complete datagram, retaining its optional second message.
func ParseHeader(b []byte) (*Header, error) {
	h := &Header{}
	if err := h.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return h, nil
}

// UnmarshalBinary resets the header and decodes a complete datagram.
// LengthError carries the safely decoded transaction fields needed by §7.7.3;
// UnsupportedVersionError supports the version-specific handling in §7.7.2.
func (h *Header) UnmarshalBinary(b []byte) error {
	return h.unmarshalBinary(b, Parse, false)
}

// TS 29.274 §§5.5.1, 7.7.0: both codecs share framing and secondary errors.
func (h *Header) unmarshalBinary(b []byte, decode func([]byte) (Message, error), replay bool) error {
	if h == nil {
		return ErrNilMessage
	}
	*h = Header{replay: replay}
	if len(b) < 1 {
		return ErrTooShortToParse
	}
	h.Flags = b[0]
	if h.Version() != 2 {
		return &UnsupportedVersionError{Version: uint8(h.Version())}
	}
	if len(b) < h.headerLen() {
		return ErrTooShortToParse
	}
	h.Type, h.Length = b[1], binary.BigEndian.Uint16(b[2:4])
	offset := fixedHeaderSize
	if h.HasTEID() {
		h.TEID = binary.BigEndian.Uint32(b[4:8])
		offset += teidSize
	}
	h.SequenceNumber = utils.Uint24To32(b[offset : offset+3])
	if h.HasTEID() && h.HasMessagePriority() {
		h.Spare = b[offset+3] & 0xf0
	}
	if replay {
		h.rawFlags = h.Flags & 3
		if !h.HasTEID() {
			h.rawFlags |= h.Flags & 4
		}
		h.rawSpare = b[offset+3]
		if h.HasMessagePriority() {
			h.rawSpare &= 0x0f
		}
	}
	end := int(h.Length) + fixedHeaderSize
	lengthErr := &LengthError{MessageType: h.Type, TEID: h.TEID, SequenceNumber: h.SequenceNumber, Declared: end, Actual: len(b)}
	// TS 29.274 §7.7.3: distinguish missing header bytes from inconsistent Length.
	if end < h.headerLen() || end > len(b) {
		return lengthErr
	}
	if h.HasTEID() != requiresTEID(h.Type) {
		return ErrInvalidTEIDFlag
	}
	if h.IsPiggybacking() {
		if len(b)-end < noTEIDHeaderSize {
			return lengthErr
		}
		if h.Type != MsgTypeCreateSessionResponse && h.Type != MsgTypeCreateBearerResponse {
			return ErrInvalidPiggyback
		}
		// TS 29.274 §7.7.3: first check the overall datagram framing.
		secondaryEnd := end + fixedHeaderSize + int(binary.BigEndian.Uint16(b[end+2:end+4]))
		if secondaryEnd != len(b) {
			return lengthErr
		}
		// TS 29.274 §7.7.0: a bad secondary message cannot cost the primary.
		if b[end]&0x10 != 0 || !validPiggyback(h.Type, b[end+1]) {
			h.PiggybackError = ErrInvalidPiggyback
		} else {
			h.Piggybacked, h.PiggybackError = decode(b[end:])
		}
		if replay && h.PiggybackError != nil {
			h.piggybackRaw = append([]byte(nil), b[end:]...)
		}
	} else if end != len(b) {
		return lengthErr
	}
	h.Payload = b[h.headerLen():end]
	return nil
}

// MarshalLen returns the complete datagram size, including any piggybacked message.
func (h *Header) MarshalLen() int {
	if h == nil {
		return 0
	}
	return h.headerLen() + len(h.Payload) + h.piggybackLen()
}

func (h *Header) setLength(length int) {
	if h == nil {
		return
	}
	if length < 0 || length > 0xffff {
		h.Length = 0
		return
	}
	h.Length = uint16(length)
}

// SetLength updates Length for this message, excluding its piggybacked message.
// MarshalTo rejects lengths that cannot be represented in sixteen bits.
func (h *Header) SetLength() {
	if h != nil {
		h.setLength(h.headerLen() + len(h.Payload) - fixedHeaderSize)
	}
}

// String returns the GTPv2 header values in human readable format.
func (h *Header) String() string {
	if h == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{Flags: %#x, Type: %d, Length: %d, TEID: %#x, SequenceNumber: %#x, Spare: %d, Payload: %#v}", h.Flags, h.Type, h.Length, h.TEID, h.SequenceNumber, h.Spare, h.Payload)
}

// IsPiggybacking reports whether the message has a piggybacked message or P flag.
func (h *Header) IsPiggybacking() bool {
	return h != nil && (h.Piggybacked != nil || h.Flags&0x10 != 0)
}

// SetPiggybacking sets P. Marshal requires a corresponding Piggybacked message.
func (h *Header) SetPiggybacking(val uint8) {
	if h != nil {
		h.Flags = (h.Flags &^ 0x10) | ((val & 1) << 4)
	}
}

// HasTEID reports whether T is set.
func (h *Header) HasTEID() bool { return h != nil && h.Flags&8 != 0 }

func (h *Header) teid() uint32 {
	if !h.HasTEID() {
		return 0
	}
	return h.TEID
}

// SetTEID sets T and TEID; Marshal validates the message-specific T policy.
func (h *Header) SetTEID(teid uint32) {
	if h != nil {
		h.Flags |= 8
		h.TEID = teid
	}
}

// Sequence returns the sequence number, or zero for a nil header.
func (h *Header) Sequence() uint32 {
	if h == nil {
		return 0
	}
	return h.SequenceNumber
}

// SetSequenceNumber rejects values outside the 24-bit field (TS 29.274 §5.5.1).
func (h *Header) SetSequenceNumber(seq uint32) error {
	if h == nil {
		return ErrNilMessage
	}
	if seq > 0xffffff {
		return ErrInvalidSequence
	}
	h.SequenceNumber = seq
	return nil
}

// HasMessagePriority reports whether MP is set.
func (h *Header) HasMessagePriority() bool { return h.HasTEID() && h.Flags&4 != 0 }

// SetMessagePriority accepts the priority value 0–15 (TS 29.274 §5.5.1).
func (h *Header) SetMessagePriority(mp uint8) error {
	if h == nil {
		return ErrNilMessage
	}
	if mp > 15 {
		return ErrInvalidPriority
	}
	if !h.HasTEID() {
		return ErrInvalidTEIDFlag
	}
	h.Flags |= 4
	h.Spare = mp << 4
	return nil
}

// MessagePriority returns the four-bit value, or zero when MP is absent.
func (h *Header) MessagePriority() uint8 {
	if !h.HasMessagePriority() {
		return 0
	}
	return h.Spare >> 4
}

// Version returns the version in the header flags, or zero for a nil header.
func (h *Header) Version() int {
	if h == nil {
		return 0
	}
	return int(h.Flags >> 5)
}

// MessageType returns the wire message type, or zero for a nil header.
func (h *Header) MessageType() uint8 {
	if h == nil {
		return 0
	}
	return h.Type
}
