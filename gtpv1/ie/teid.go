// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// NewTEIDDataI creates a new TEIDDataI IE.
func NewTEIDDataI(teid uint32) *IE {
	return newUint32ValIE(TEIDDataI, teid)
}

// NewTEIDCPlane creates a new TEID C-Plane IE.
func NewTEIDCPlane(teid uint32) *IE {
	return newUint32ValIE(TEIDCPlane, teid)
}

// NewTEIDDataII creates a new TEIDDataII IE.
func NewTEIDDataII(nsapi uint8, teid uint32) *IE {
	// TS 29.060 §7.7.15: one NSAPI octet followed by four TEID octets.
	if nsapi > 15 {
		return newErrorIE(TEIDDataII, ErrMalformed)
	}
	p := make([]byte, 5)
	p[0] = nsapi
	binary.BigEndian.PutUint32(p[1:], teid)
	return New(TEIDDataII, p)
}

// TEID returns TEID value if type matches.
func (i *IE) TEID() (uint32, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if len(i.Payload) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	switch i.Type {
	case TEIDDataII:
		if len(i.Payload) != 5 {
			return 0, ErrInvalidLength
		}
		return binary.BigEndian.Uint32(i.Payload[1:]), nil
	case TEIDCPlane, TEIDDataI:
		return binary.BigEndian.Uint32(i.Payload), nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// MustTEID returns TEID in uint32 if type matches.
// This should only be used if it is assured to have the value.
func (i *IE) MustTEID() uint32 {
	v, _ := i.TEID()
	return v
}
