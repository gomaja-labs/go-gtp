// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// NewNSAPI creates a new NSAPI IE.
func NewNSAPI(nsapi uint8) *IE {
	if nsapi > 15 {
		return newErrorIE(NSAPI, ErrMalformed)
	}
	return newUint8ValIE(NSAPI, nsapi)
}

// NSAPI returns NSAPI value if type matches.
func (i *IE) NSAPI() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != NSAPI && i.Type != TEIDDataII {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) == 0 {
		return 0, io.ErrUnexpectedEOF
	}

	// TS 29.060 §7.7.17: ignore the spare high nibble.
	return i.Payload[0] & 0x0f, nil
}

// MustNSAPI returns NSAPI in uint8 if type matches.
// This should only be used if it is assured to have the value.
func (i *IE) MustNSAPI() uint8 {
	v, _ := i.NSAPI()
	return v
}
