// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv1/ie"

// ErrorIndication is a ErrorIndication Header and its IEs above.
type ErrorIndication struct {
	*Header
	TEIDDataI         *ie.IE
	GTPUPeerAddress   *ie.IE
	PrivateExtensions []*ie.IE
	AdditionalIEs     []*ie.IE
}

// NewErrorIndication creates a new GTPv1 ErrorIndication.
func NewErrorIndication(teid uint32, seq uint16, ies ...*ie.IE) *ErrorIndication {
	e := &ErrorIndication{
		Header: NewHeader(0x32, MsgTypeErrorIndication, teid, seq, nil),
	}

	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.TEIDDataI:
			if e.TEIDDataI == nil {
				e.TEIDDataI = i
			} else {
				e.AdditionalIEs = append(e.AdditionalIEs, i)
			}
		case ie.GSNAddress:
			if e.GTPUPeerAddress == nil {
				e.GTPUPeerAddress = i
			} else {
				e.AdditionalIEs = append(e.AdditionalIEs, i)
			}
		case ie.PrivateExtension:
			e.PrivateExtensions = append(e.PrivateExtensions, i)
		default:
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	}

	e.SetLength()
	return e
}

// Marshal returns the byte sequence generated from a ErrorIndication.
func (e *ErrorIndication) Marshal() ([]byte, error) {
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
func (e *ErrorIndication) MarshalTo(b []byte) error {
	if e == nil || e.Header == nil {
		return ErrInvalidLength
	}
	ies := []*ie.IE{e.TEIDDataI, e.GTPUPeerAddress}
	ies = append(ies, e.PrivateExtensions...)
	ies = append(ies, e.AdditionalIEs...)
	return marshalIEs(e.Header, b, ies)
}

// ParseErrorIndication decodes a given byte sequence as a ErrorIndication.
func ParseErrorIndication(b []byte) (*ErrorIndication, error) {
	e := &ErrorIndication{}
	if err := e.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return e, err
		}
		return nil, err
	}
	return e, nil
}

// UnmarshalBinary decodes a given byte sequence as a ErrorIndication.
func (e *ErrorIndication) UnmarshalBinary(b []byte) error {
	if e == nil {
		return ErrInvalidLength
	}
	*e = ErrorIndication{}
	var err error
	e.Header, err = ParseHeader(b)
	if err != nil {
		return err
	}
	if e.Header.Type != MsgTypeErrorIndication {
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
		case ie.TEIDDataI:
			if e.TEIDDataI == nil {
				e.TEIDDataI = i
			} else {
				e.AdditionalIEs = append(e.AdditionalIEs, i)
			}
		case ie.GSNAddress:
			if e.GTPUPeerAddress == nil {
				e.GTPUPeerAddress = i
			} else {
				e.AdditionalIEs = append(e.AdditionalIEs, i)
			}
		case ie.PrivateExtension:
			e.PrivateExtensions = append(e.PrivateExtensions, i)
		default:
			e.AdditionalIEs = append(e.AdditionalIEs, i)
		}
	}
	return validateMandatory(e.Header.Type, e.Header.Payload)
}

// MarshalLen returns the serial length of Data.
func (e *ErrorIndication) MarshalLen() int {
	if e == nil || e.Header == nil {
		return 0
	}
	ies := []*ie.IE{e.TEIDDataI, e.GTPUPeerAddress}
	ies = append(ies, e.PrivateExtensions...)
	ies = append(ies, e.AdditionalIEs...)
	return messageLen(e.Header, ies)
}

// SetLength sets the length in Length field.
func (e *ErrorIndication) SetLength() {
	if e == nil || e.Header == nil {
		return
	}
	e.Length = uint16(e.MarshalLen() - 8)
}

// MessageTypeName returns the name of protocol.
func (e *ErrorIndication) MessageTypeName() string {
	return "Errror Indication"
}

// TEID returns the TEID in human-readable string.
func (e *ErrorIndication) TEID() uint32 {
	if e == nil || e.Header == nil {
		return 0
	}
	return e.Header.TEID
}
