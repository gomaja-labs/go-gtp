// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"

	"github.com/gomaja-labs/go-gtp/utils"
)

// NewMSISDN creates a new MSISDN IE.
func NewMSISDN(msisdn string) *IE {
	if len(msisdn) == 0 {
		return newErrorIE(MSISDN, ErrInvalidLength)
	}
	m, err := utils.EncodeBCDDigits(msisdn)
	if err != nil {
		return newErrorIE(MSISDN, err)
	}
	return New(MSISDN, 0x00, m)
}

// MSISDN returns MSISDN in string if the type of IE matches.
func (i *IE) MSISDN() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if i.Type != MSISDN {
		return "", &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 1 {
		return "", io.ErrUnexpectedEOF
	}

	v, err := utils.DecodeBCDDigits(i.Payload)
	if err != nil {
		return "", err
	}
	// TS 29.274 §8.11: only a single final high nibble may be filler.
	if len(v) == 0 || (len(v)+1)/2 != len(i.Payload) {
		return "", ErrMalformed
	}
	return v, nil
}

// MustMSISDN returns MSISDN in string, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustMSISDN() string {
	if i == nil {
		return ""
	}
	v, _ := i.MSISDN()
	return v
}
