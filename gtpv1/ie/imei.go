// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"

	"github.com/gomaja-labs/go-gtp/utils"
)

// NewIMEISV creates a new IMEISV IE.
func NewIMEISV(imei string) *IE {
	// TS 29.060 §7.7.53: 15-digit IMEI or 16-digit IMEISV.
	if len(imei) != 15 && len(imei) != 16 {
		return newErrorIE(IMEISV, ErrInvalidLength)
	}
	i, err := utils.EncodeBCDDigits(imei)
	if err != nil {
		return newErrorIE(IMEISV, err)
	}
	return New(IMEISV, i)
}

// IMEISV returns IMEISV value if type matches.
func (i *IE) IMEISV() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if i.Type != IMEISV {
		return "", &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) == 0 {
		return "", io.ErrUnexpectedEOF
	}
	return utils.DecodeBCDDigits(i.Payload)
}

// MustIMEISV returns IMEISV in string if type matches.
// This should only be used if it is assured to have the value.
func (i *IE) MustIMEISV() string {
	v, _ := i.IMEISV()
	return v
}
