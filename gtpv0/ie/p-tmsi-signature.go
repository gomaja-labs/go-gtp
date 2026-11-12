// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"

	"github.com/gomaja-labs/go-gtp/utils"
)

// NewPTMSISignature creates a new PTMSISignature IE.
func NewPTMSISignature(sig uint32) *IE {
	b, err := utils.Uint32To24(sig)
	if err != nil {
		return newErrorIE(PTMSISignature, err)
	}
	return New(PTMSISignature, b)
}

// PTMSISignature returns PTMSISignature value in uint32 if type matches.
func (i *IE) PTMSISignature() (uint32, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != PTMSISignature {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 3 {
		return 0, io.ErrUnexpectedEOF
	}

	return utils.Uint24To32(i.Payload), nil
}

// MustPTMSISignature returns PTMSISignature in uint32 if type matches.
// This should only be used if it is assured to have the value.
func (i *IE) MustPTMSISignature() uint32 {
	v, _ := i.PTMSISignature()
	return v
}
