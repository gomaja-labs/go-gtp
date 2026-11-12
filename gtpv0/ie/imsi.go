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
	// TS 29.060 §7.7.2; GSM 09.60 §7.9.2: eight octets, unused digits are F.
	if len(imsi) == 0 || len(imsi) > 15 {
		return newErrorIE(IMSI, ErrInvalidLength)
	}
	i, err := utils.EncodeBCDDigits(imsi)
	if err != nil {
		return newErrorIE(IMSI, err)
	}
	p := make([]byte, 8)
	for n := range p {
		p[n] = 0xff
	}
	copy(p, i)
	return New(IMSI, p)
}

// IMSI returns IMSI value in string if type matches.
func (i *IE) IMSI() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if i.Type != IMSI {
		return "", &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) == 0 {
		return "", io.ErrUnexpectedEOF
	}

	return utils.DecodeBCDDigits(i.Payload)
}

// MustIMSI returns IMSI in string if type matches.
// This should only be used if it is assured to have the value.
func (i *IE) MustIMSI() string {
	v, _ := i.IMSI()
	return v
}
