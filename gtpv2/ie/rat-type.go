// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewRATType creates a new RATType IE.
func NewRATType(rat uint8) *IE {
	// TS 29.274 §8.17: validate the transmitted value.
	if rat < 1 || rat > 22 {
		return newErrorIE(RATType, ErrMalformed)
	}
	return NewUint8IE(RATType, rat)
}

// RATType returns RATType in uint8 if the type of IE matches.
func (i *IE) RATType() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	switch i.Type {
	case RATType:
		v, err := i.ValueAsUint8()
		if err != nil {
			return 0, err
		}
		if v < 1 || v > 22 {
			return 0, ErrMalformed
		}
		return v, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// MustRATType returns RATType in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustRATType() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.RATType()
	return v
}
