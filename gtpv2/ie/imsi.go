// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"

	"github.com/gomaja-labs/go-gtp/utils"
)

// NewIMSI creates a new IMSI IE.
func NewIMSI(imsi string) *IE {
	// TS 29.274 §8.3: identity digit count.
	if len(imsi) == 0 || len(imsi) > 15 {
		return newErrorIE(IMSI, ErrInvalidLength)
	}
	i, err := utils.EncodeBCDDigits(imsi)
	if err != nil {
		return newErrorIE(IMSI, err)
	}

	return New(IMSI, 0x00, i)
}

// IMSI returns IMSI in string if the type of IE matches.
func (i *IE) IMSI() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if i.Type != IMSI {
		return "", &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 1 {
		return "", io.ErrUnexpectedEOF
	}

	v, err := utils.DecodeBCDDigits(i.Payload)
	if err != nil {
		return "", err
	}
	if len(v) == 0 || len(v) > 15 {
		return "", ErrInvalidLength
	}
	// TS 29.274 §§8.3, 8.10: only one trailing high-nibble filler.
	if (len(v)+1)/2 != len(i.Payload) {
		return "", ErrMalformed
	}
	return v, nil
}

// MustIMSI returns IMSI in string, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustIMSI() string {
	if i == nil {
		return ""
	}
	v, _ := i.IMSI()
	return v
}
