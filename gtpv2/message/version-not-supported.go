// Copyright 2019-2024 go-gtp authors. All rights reservev.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// VersionNotSupportedIndication is a VersionNotSupportedIndication Header and its IEs above.
type VersionNotSupportedIndication struct {
	*Header
	AdditionalIEs []*ie.IE
}

// NewVersionNotSupportedIndication creates a new VersionNotSupportedIndication.
func NewVersionNotSupportedIndication(teid, seq uint32, ie ...*ie.IE) *VersionNotSupportedIndication {
	v := &VersionNotSupportedIndication{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 0),
			MsgTypeVersionNotSupportedIndication, teid, seq, nil,
		),
	}

	for _, i := range ie {
		if i == nil {
			continue
		}
		v.AdditionalIEs = append(v.AdditionalIEs, i)
	}

	v.SetLength()
	return v
}

// Marshal serializes VersionNotSupportedIndication into bytes.
func (v *VersionNotSupportedIndication) Marshal() ([]byte, error) {
	b := make([]byte, v.MarshalLen())
	if err := v.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes VersionNotSupportedIndication into bytes.
func (v *VersionNotSupportedIndication) MarshalTo(b []byte) error {
	if v == nil || v.Header == nil {
		return ErrNilMessage
	}
	if err := v.Header.validate(MsgTypeVersionNotSupportedIndication, v.MarshalLen()); err != nil {
		return err
	}
	if len(b) < v.MarshalLen() {
		return ErrTooShortToMarshal
	}
	v.Header.Payload = nil
	v.Header.Payload = make([]byte, v.MarshalLen()-v.Header.MarshalLen())
	offset := 0
	for _, i := range v.AdditionalIEs {
		if i == nil {
			continue
		}
		if err := i.MarshalTo(v.Payload[offset:]); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	v.Header.SetLength()
	return v.Header.MarshalTo(b)
}

// ParseVersionNotSupportedIndication decodes given bytes as VersionNotSupportedIndication.
func ParseVersionNotSupportedIndication(b []byte) (*VersionNotSupportedIndication, error) {
	v := &VersionNotSupportedIndication{}
	if err := v.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return v, err
		}
		return nil, err
	}
	return v, nil
}

// UnmarshalBinary decodes given bytes as VersionNotSupportedIndication.
func (v *VersionNotSupportedIndication) UnmarshalBinary(b []byte) error {
	if v == nil {
		return ErrNilMessage
	}
	*v = VersionNotSupportedIndication{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeVersionNotSupportedIndication {
		return &UnexpectedTypeError{Expected: MsgTypeVersionNotSupportedIndication, Actual: h.Type}
	}
	v.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, v.addIE)
}

// MarshalLen returns the serial length in int.
func (v *VersionNotSupportedIndication) MarshalLen() int {
	if v == nil || v.Header == nil {
		return 0
	}
	l := v.Header.MarshalLen() - len(v.Header.Payload)
	for _, i := range v.AdditionalIEs {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	return l
}

// SetLength sets the length in Length fielv.
func (v *VersionNotSupportedIndication) SetLength() {
	if v != nil && v.Header != nil {
		v.Header.setLength(v.MarshalLen() - v.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (v *VersionNotSupportedIndication) MessageTypeName() string {
	return "Version Not Supported Indication"
}

// TEID returns the TEID in uint32.
func (v *VersionNotSupportedIndication) TEID() uint32 {
	if v == nil {
		return 0
	}
	return v.Header.teid()
}

func (v *VersionNotSupportedIndication) GetHeader() *Header {
	if v == nil {
		return nil
	}
	return v.Header
}

// TS 29.274 §7: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (v *VersionNotSupportedIndication) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		v.Header.err = err
	}
	v.AdditionalIEs = append(v.AdditionalIEs, i)
}

// Version reads the common header, returning zero for a nil receiver.
func (v *VersionNotSupportedIndication) Version() int { return v.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (v *VersionNotSupportedIndication) MessageType() uint8 { return v.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (v *VersionNotSupportedIndication) Sequence() uint32 { return v.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (v *VersionNotSupportedIndication) IsPiggybacking() bool { return v.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (v *VersionNotSupportedIndication) HasTEID() bool { return v.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (v *VersionNotSupportedIndication) HasMessagePriority() bool {
	return v.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (v *VersionNotSupportedIndication) MessagePriority() uint8 {
	return v.GetHeader().MessagePriority()
}

// SetSequenceNumber updates the common header safely on a nil receiver.
func (v *VersionNotSupportedIndication) SetSequenceNumber(seq uint32) error {
	return v.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (v *VersionNotSupportedIndication) SetMessagePriority(priority uint8) error {
	return v.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (v *VersionNotSupportedIndication) SetTEID(teid uint32) { v.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (v *VersionNotSupportedIndication) SetPiggybacking(value uint8) {
	v.GetHeader().SetPiggybacking(value)
}
