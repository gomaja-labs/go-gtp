// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"fmt"
	"io"
)

// NewCause creates a new Cause IE.
func NewCause(cause uint8, pce, bce, cs uint8, offendingIE *IE) *IE {
	// TS 29.274 §8.4: reject fields that do not fit.
	if pce > 1 || bce > 1 || cs > 1 {
		return newErrorIE(Cause, ErrMalformed)
	}
	i := New(Cause, 0x00, make([]byte, 2))
	i.Payload[0] = cause
	i.Payload[1] = ((pce << 2) & 0x04) | ((bce << 1) & 0x02) | cs&0x01

	if offendingIE != nil {
		// TS 29.274 §8.4: retain Type and Instance; Length and spare are zero.
		i.Payload = append(i.Payload, []byte{offendingIE.Type, 0x00, 0x00, offendingIE.Instance()}...)
		i.SetLength()
	}
	return i
}

// Cause returns Cause in uint8 if the type of IE matches.
func (i *IE) Cause() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	switch i.Type {
	case Cause:
		if len(i.Payload) < 2 {
			return 0, ErrInvalidLength
		}
		return i.ValueAsUint8()
	case BearerContext:
		ies, err := i.BearerContext()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve Cause: %w", err)
		}

		// TS 29.274 §7.7.9: match the child type and instance together.
		for _, child := range ies {
			if child != nil && child.Type == Cause && child.Instance() == 0 {
				return child.Cause()
			}
		}
		return 0, ErrIENotFound
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// MustCause returns Cause in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustCause() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.Cause()
	return v
}

// CauseFlags returns CauseFlags in uint8 if the type of IE matches.
func (i *IE) CauseFlags() (uint8, error) {
	// TS 29.274 §8.4: CS, BCE, PCE occupy bits 1, 2, 3.
	if err := i.Err(); err != nil {
		return 0, err
	}
	switch i.Type {
	case Cause:
		if len(i.Payload) < 2 {
			return 0, io.ErrUnexpectedEOF
		}

		return i.Payload[1] & 7, nil
	case BearerContext:
		ies, err := i.BearerContext()
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve Cause: %w", err)
		}

		// TS 29.274 §7.7.9: match the child type and instance together.
		for _, child := range ies {
			if child != nil && child.Type == Cause && child.Instance() == 0 {
				return child.CauseFlags()
			}
		}
		return 0, ErrIENotFound
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// MustCauseFlags returns CauseFlags in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustCauseFlags() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.CauseFlags()
	return v
}

// HasCS reports whether an IE has CS bit.
func (i *IE) HasCS() bool {
	if i == nil {
		return false
	}
	v, err := i.CauseFlags()
	if err != nil {
		return false
	}

	return has1stBit(v)
}

// HasBCE reports whether an IE has BCE bit.
func (i *IE) HasBCE() bool {
	if i == nil {
		return false
	}
	v, err := i.CauseFlags()
	if err != nil {
		return false
	}

	return has2ndBit(v)
}

// HasPCE reports whether an IE has PCE bit.
func (i *IE) HasPCE() bool {
	if i == nil {
		return false
	}
	v, err := i.CauseFlags()
	if err != nil {
		return false
	}

	return has3rdBit(v)
}

// IsRemoteCause returns IsRemoteCause in bool if the type of IE matches.
func (i *IE) IsRemoteCause() bool {
	if i.Err() != nil {
		return false
	}
	if i.Type != Cause {
		return false
	}

	if len(i.Payload) < 2 {
		return false
	}

	if i.Payload[1]&0x01 == 1 {
		return true
	}
	return false
}

// IsBearerContextIEError returns IsBearerContextIEError in bool if the type of IE matches.
func (i *IE) IsBearerContextIEError() bool {
	if i.Err() != nil {
		return false
	}
	if i.Type != Cause {
		return false
	}

	if len(i.Payload) < 2 {
		return false
	}

	if i.Payload[1]>>1&0x01 == 1 {
		return true
	}
	return false
}

// IsPDNConnectionIEError returns IsPDNConnectionIEError in bool if the type of IE matches.
func (i *IE) IsPDNConnectionIEError() bool {
	if i.Err() != nil {
		return false
	}
	if i.Type != Cause {
		return false
	}

	if len(i.Payload) < 2 {
		return false
	}

	if i.Payload[1]>>2&0x01 == 1 {
		return true
	}
	return false
}

// OffendingIE returns OffendingIE in *IE if the type of IE matches.
//
// Note that the returned IE has no payload (cf. §8.4, TS29.274).
func (i *IE) OffendingIE() (*IE, error) {

	if err := i.Err(); err != nil {
		return nil, err
	}
	if i.Type != Cause {
		return nil, &InvalidTypeError{Type: i.Type}
	}

	if len(i.Payload) < 6 {
		return nil, io.ErrUnexpectedEOF
	}

	return Parse(i.Payload[2:6])
}

// MustOffendingIE returns OffendingIE in *IE, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustOffendingIE() *IE {
	v, _ := i.OffendingIE()
	return v
}
