// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewPDNType creates a new PDNType IE.
func NewPDNType(pdn uint8) *IE {
	// TS 29.274 §8.34: validate the transmitted value.
	if pdn < 1 || pdn > 5 {
		return newErrorIE(PDNType, ErrMalformed)
	}
	return NewUint8IE(PDNType, pdn)
}

// PDNType returns the PDNType value in uint8 if the type of IE matches.
func (i *IE) PDNType() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	switch i.Type {
	case PDNType, PDNAddressAllocation:
		v, err := i.ValueAsUint8()
		if err != nil {
			return 0, err
		}
		v &= 7
		if v < 1 || v > 5 {
			return 0, ErrMalformed
		}
		return v, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// MustPDNType returns PDNType in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustPDNType() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.PDNType()
	return v
}
