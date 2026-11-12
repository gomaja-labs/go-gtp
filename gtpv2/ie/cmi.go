// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// NewCSGMembershipIndication creates a new CSGMembershipIndication IE.
func NewCSGMembershipIndication(cmi uint8) *IE {
	// TS 29.274 §8.79: validate the transmitted value.
	if cmi > 1 {
		return newErrorIE(CSGMembershipIndication, ErrMalformed)
	}
	return NewUint8IE(CSGMembershipIndication, cmi)
}

// CMI returns CMI in uint8 if the type of IE matches.
func (i *IE) CMI() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	switch i.Type {
	case CSGMembershipIndication:
		v, err := i.ValueAsUint8()
		return v & 1, err
	case UserCSGInformation:
		if len(i.Payload) < 8 {
			return 0, io.ErrUnexpectedEOF
		}
		// TS 29.274 §8.75: ignore CMI outside Hybrid Mode.
		if i.Payload[7]>>6 != 1 {
			return 0, nil
		}
		return i.Payload[7] & 1, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// MustCMI returns CMI in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustCMI() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.CMI()
	return v
}
