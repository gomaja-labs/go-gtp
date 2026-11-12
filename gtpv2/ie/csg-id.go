// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
)

// NewCSGID creates a new CSGID IE.
func NewCSGID(id uint32) *IE {
	// TS 29.274 §8.78: validate the transmitted value.
	if id > 0x7ffffff {
		return newErrorIE(CSGID, ErrMalformed)
	}
	return NewUint32IE(CSGID, id&0x7ffffff)
}

// CSGID returns CSGID in uint32 if the type of IE matches.
func (i *IE) CSGID() (uint32, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	switch i.Type {
	case CSGID:
		v, err := i.ValueAsUint32()
		return v & 0x7ffffff, err
	case UserCSGInformation:
		if len(i.Payload) < 7 {
			return 0, io.ErrUnexpectedEOF
		}
		return binary.BigEndian.Uint32(i.Payload[3:7]) & 0x7ffffff, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// MustCSGID returns CSGID in uint32, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustCSGID() uint32 {
	if i == nil {
		return 0
	}
	v, _ := i.CSGID()
	return v
}
