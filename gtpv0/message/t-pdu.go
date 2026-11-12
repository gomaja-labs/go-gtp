// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

// TPDU represents a T-PDU type of GTPv0 message.
type TPDU struct {
	*Header
}

// NewTPDU creates a new TPDU.
func NewTPDU(seq, label uint16, tid uint64, payload []byte) *TPDU {
	t := &TPDU{
		Header: NewHeader(
			0x1e, MsgTypeTPDU, seq, label, tid, payload,
		),
	}
	t.SetLength()
	return t
}

// Marshal returns the byte sequence generated from a TPDU.
func (t *TPDU) Marshal() ([]byte, error) {
	if t == nil || t.Header == nil {
		return nil, ErrInvalidLength
	}
	b := make([]byte, t.MarshalLen())
	if err := t.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (t *TPDU) MarshalTo(b []byte) error {
	if t == nil || t.Header == nil {
		return ErrInvalidLength
	}
	if len(b) < t.MarshalLen() {
		return ErrTooShortToMarshal
	}

	t.Header.Payload = t.Payload
	t.Header.SetLength()
	return t.Header.MarshalTo(b)
}

// ParseTPDU parses a given byte sequence as a TPDU.
func ParseTPDU(b []byte) (*TPDU, error) {
	t := &TPDU{}
	if err := t.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return t, nil
}

// UnmarshalBinary parses a given byte sequence as a TPDU.
func (t *TPDU) UnmarshalBinary(b []byte) error {
	if t == nil {
		return ErrInvalidLength
	}
	*t = TPDU{}
	var err error
	t.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if t.Header.Type != MsgTypeTPDU {
		return ErrInvalidMessageType
	}
	t.Payload = t.Header.Payload

	return nil
}

// MarshalLen returns the serial length of Data.
func (t *TPDU) MarshalLen() int {
	if t == nil || t.Header == nil {
		return 0
	}
	return t.Header.MarshalLen() - len(t.Header.Payload) + len(t.Payload)
}

// SetLength sets the length in Length field.
func (t *TPDU) SetLength() {
	if t == nil || t.Header == nil {
		return
	}
	t.Header.Length = uint16(len(t.Payload))
}

// MessageTypeName returns the name of protocol.
func (t *TPDU) MessageTypeName() string {
	return "T-PDU"
}

// TID returns the TID in human-readable string.
func (t *TPDU) TID() string {
	if t == nil || t.Header == nil {
		return ""
	}
	return t.tid()
}
