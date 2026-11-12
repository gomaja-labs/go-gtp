// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewRFSPIndex creates a new RFSPIndex IE.
func NewRFSPIndex(idx uint16) *IE {
	// TS 29.274 §8.77: validate the transmitted value.
	if idx == 0 || idx > 256 {
		return newErrorIE(RFSPIndex, ErrMalformed)
	}
	return NewUint16IE(RFSPIndex, idx)
}

// RFSPIndex returns RFSPIndex in uint16 if the type of IE matches.
func (i *IE) RFSPIndex() (uint16, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != RFSPIndex {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	v, err := i.ValueAsUint16()
	if err != nil {
		return 0, err
	}
	if v < 1 || v > 256 {
		return 0, ErrMalformed
	}
	return v, nil
}

// MustRFSPIndex returns RFSPIndex in uint16, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustRFSPIndex() uint16 {
	if i == nil {
		return 0
	}
	v, _ := i.RFSPIndex()
	return v
}
