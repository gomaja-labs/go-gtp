// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewLocalDistinguishedName creates a new LocalDistinguishedName IE.
func NewLocalDistinguishedName(name string) *IE {
	if !validLDN(name) {
		return newErrorIE(LocalDistinguishedName, ErrMalformed)
	}
	return NewStringIE(LocalDistinguishedName, name)
}

// LocalDistinguishedName returns LocalDistinguishedName in string if the type of IE matches.
func (i *IE) LocalDistinguishedName() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if i.Type != LocalDistinguishedName {
		return "", &InvalidTypeError{Type: i.Type}
	}
	if !validLDN(string(i.Payload)) {
		return "", ErrMalformed
	}
	return i.ValueAsString()
}

// MustLocalDistinguishedName returns LocalDistinguishedName in string, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustLocalDistinguishedName() string {
	if i == nil {
		return ""
	}
	v, _ := i.LocalDistinguishedName()
	return v
}

// TS 29.274 §8.82: 1..400 ASCII octets.
func validLDN(s string) bool {
	if len(s) < 1 || len(s) > 400 {
		return false
	}
	for _, b := range []byte(s) {
		if b > 127 {
			return false
		}
	}
	return true
}
