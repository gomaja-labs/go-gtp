// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewServiceIndicator creates a new ServiceIndicator IE.
func NewServiceIndicator(ind uint8) *IE {
	// TS 29.274 §8.80: validate the transmitted value.
	if ind < 1 || ind > 2 {
		return newErrorIE(ServiceIndicator, ErrMalformed)
	}
	return NewUint8IE(ServiceIndicator, ind)
}

// ServiceIndicator returns ServiceIndicator in uint8 if the type of IE matches.
func (i *IE) ServiceIndicator() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != ServiceIndicator {
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

// MustServiceIndicator returns ServiceIndicator in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustServiceIndicator() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.ServiceIndicator()
	return v
}
