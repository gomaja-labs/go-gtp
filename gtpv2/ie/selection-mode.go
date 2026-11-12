// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewSelectionMode creates a new SelectionMode IE.
func NewSelectionMode(mode uint8) *IE {
	// TS 29.274 §8.58: validate the transmitted value.
	if mode > 2 {
		return newErrorIE(SelectionMode, ErrMalformed)
	}
	return NewUint8IE(SelectionMode, mode)
}

// SelectionMode returns SelectionMode value if the type of IE matches.
func (i *IE) SelectionMode() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != SelectionMode {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	v, err := i.ValueAsUint8()
	v &= 3
	if v == 3 {
		v = 2
	} // TS 29.274 §8.58.
	return v, err
}

// MustSelectionMode returns SelectionMode in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustSelectionMode() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.SelectionMode()
	return v
}
