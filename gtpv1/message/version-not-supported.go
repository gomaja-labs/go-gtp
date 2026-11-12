// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// VersionNotSupported is a VersionNotSupported Header and its IEs above.
type VersionNotSupported struct {
	*Header
	AdditionalIEs []*ie.IE
}

// NewVersionNotSupported creates a new GTPv1 VersionNotSupported.
func NewVersionNotSupported(teid uint32, seq uint16, ie ...*ie.IE) *VersionNotSupported {
	v := &VersionNotSupported{
		Header: NewHeader(0x32, MsgTypeVersionNotSupported, teid, seq, nil),
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

// Marshal returns the byte sequence generated from a VersionNotSupported.
func (v *VersionNotSupported) Marshal() ([]byte, error) {
	if v == nil || v.Header == nil {
		return nil, ErrInvalidLength
	}
	b := make([]byte, v.MarshalLen())
	if err := v.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (v *VersionNotSupported) MarshalTo(b []byte) error {
	if v == nil || v.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{}
	ies = append(ies, v.AdditionalIEs...)
	return marshalIEs(v.Header, b, ies)
}

// ParseVersionNotSupported decodes a given byte sequence as a VersionNotSupported.
func ParseVersionNotSupported(b []byte) (*VersionNotSupported, error) {
	v := &VersionNotSupported{}
	if err := v.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return v, nil
}

// UnmarshalBinary decodes a given byte sequence as a VersionNotSupported.
func (v *VersionNotSupported) UnmarshalBinary(b []byte) error {
	if v == nil {
		return ErrInvalidLength
	}
	*v = VersionNotSupported{}
	var err error
	v.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if v.Header.Type != MsgTypeVersionNotSupported {
		return ErrInvalidMessageType
	}

	ie, err := ie.ParseMultiIEs(v.Header.Payload)
	if err != nil {
		return err
	}

	for _, i := range ie {
		if i == nil {
			continue
		}
		v.AdditionalIEs = append(v.AdditionalIEs, i)
	}
	return validateMandatory(v.Header.Type, v.Header.Payload)
}

// MarshalLen returns the serial length of Data.
func (v *VersionNotSupported) MarshalLen() int {
	if v == nil || v.Header == nil {
		return 0
	}
	ies := []*ie.IE{}
	ies = append(ies, v.AdditionalIEs...)
	return messageLen(v.Header, ies)
}

// SetLength sets the length in Length field.
func (v *VersionNotSupported) SetLength() {
	if v == nil || v.Header == nil {
		return
	}
	v.Length = uint16(v.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (v *VersionNotSupported) MessageTypeName() string {
	return "Version Not Supported"
}

// TEID returns the TEID in human-readable string.
func (v *VersionNotSupported) TEID() uint32 {
	if v == nil || v.Header == nil {
		return 0
	}
	return v.Header.TEID
}
