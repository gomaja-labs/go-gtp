// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// SupportedExtensionHeaderNotification is a SupportedExtensionHeaderNotification Header and its IEs above.
type SupportedExtensionHeaderNotification struct {
	*Header
	ExtensionHeaderTypeList *ie.IE
	AdditionalIEs           []*ie.IE
}

// NewSupportedExtensionHeaderNotification creates a new GTPv1 SupportedExtensionHeaderNotification.
func NewSupportedExtensionHeaderNotification(teid uint32, seq uint16, ies ...*ie.IE) *SupportedExtensionHeaderNotification {
	e := &SupportedExtensionHeaderNotification{
		Header: NewHeader(0x32, MsgTypeSupportedExtensionHeaderNotification, teid, seq, nil),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.ExtensionHeaderTypeList:
			if e.ExtensionHeaderTypeList == nil {
				e.ExtensionHeaderTypeList = i
			} else {
				e.AdditionalIEs = append(e.AdditionalIEs, i)
			}
		default:
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	}

	e.SetLength()
	return e
}

// Marshal returns the byte sequence generated from a SupportedExtensionHeaderNotification.
func (e *SupportedExtensionHeaderNotification) Marshal() ([]byte, error) {
	if e == nil || e.Header == nil {
		return nil, ErrInvalidLength
	}
	b := make([]byte, e.MarshalLen())
	if err := e.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (e *SupportedExtensionHeaderNotification) MarshalTo(b []byte) error {
	if e == nil || e.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{e.ExtensionHeaderTypeList}
	ies = append(ies, e.AdditionalIEs...)
	return marshalIEs(e.Header, b, ies)
}

// ParseSupportedExtensionHeaderNotification decodes a given byte sequence as a SupportedExtensionHeaderNotification.
func ParseSupportedExtensionHeaderNotification(b []byte) (*SupportedExtensionHeaderNotification, error) {
	e := &SupportedExtensionHeaderNotification{}
	if err := e.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return e, err
		}
		return nil, err
	}
	return e, nil
}

// UnmarshalBinary decodes a given byte sequence as a SupportedExtensionHeaderNotification.
func (e *SupportedExtensionHeaderNotification) UnmarshalBinary(b []byte) error {
	if e == nil {
		return ErrInvalidLength
	}
	*e = SupportedExtensionHeaderNotification{}
	var err error
	e.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if e.Header.Type != MsgTypeSupportedExtensionHeaderNotification {
		return ErrInvalidMessageType
	}

	ies, err := ie.ParseMultiIEs(e.Header.Payload)
	if err != nil {
		return err
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.ExtensionHeaderTypeList:
			if e.ExtensionHeaderTypeList == nil {
				e.ExtensionHeaderTypeList = i
			} else {
				e.AdditionalIEs = append(e.AdditionalIEs, i)
			}
		default:
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	}
	return validateMandatory(e.Header.Type, e.Header.Payload)
}

// MarshalLen returns the serial length of Data.
func (e *SupportedExtensionHeaderNotification) MarshalLen() int {
	if e == nil || e.Header == nil {
		return 0
	}
	ies := []*ie.IE{e.ExtensionHeaderTypeList}
	ies = append(ies, e.AdditionalIEs...)
	return messageLen(e.Header, ies)
}

// SetLength sets the length in Length field.
func (e *SupportedExtensionHeaderNotification) SetLength() {
	if e == nil || e.Header == nil {
		return
	}
	e.Length = uint16(e.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (e *SupportedExtensionHeaderNotification) MessageTypeName() string {
	return "Supported Extension Header Notification"
}

// TEID returns the TEID in human-readable string.
func (e *SupportedExtensionHeaderNotification) TEID() uint32 {
	if e == nil || e.Header == nil {
		return 0
	}
	return e.Header.TEID
}
