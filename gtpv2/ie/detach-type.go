// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewDetachType creates a new DetachType IE.
func NewDetachType(dtype uint8) *IE {
	// TS 29.274 §8.81: validate the transmitted value.
	if dtype < 1 || dtype > 2 {
		return newErrorIE(DetachType, ErrMalformed)
	}
	return NewUint8IE(DetachType, dtype)
}

// DetachType returns DetachType in uint8 if the type of IE matches.
func (i *IE) DetachType() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != DetachType {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	v, err := i.ValueAsUint8()
	if err != nil {
		return 0, err
	}
	if v < 1 || v > 2 {
		return 0, ErrMalformed
	}
	return v, nil
}

// MustDetachType returns DetachType in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustDetachType() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.DetachType()
	return v
}
