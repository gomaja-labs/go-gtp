// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"

	"github.com/gomaja-labs/go-gtp/utils"
)

// NewMobileEquipmentIdentity creates a new MobileEquipmentIdentity IE.
func NewMobileEquipmentIdentity(mei string) *IE {
	m, err := utils.EncodeBCDDigits(mei)
	if err != nil {
		return newErrorIE(MobileEquipmentIdentity, err)
	}
	return New(MobileEquipmentIdentity, 0x00, m)
}

// MobileEquipmentIdentity returns MobileEquipmentIdentity in string if the
// type of IE matches.
func (i *IE) MobileEquipmentIdentity() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if i.Type != MobileEquipmentIdentity {
		return "", &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 1 {
		return "", io.ErrUnexpectedEOF
	}

	return utils.DecodeBCDDigits(i.Payload)
}

// MustMobileEquipmentIdentity returns MobileEquipmentIdentity in string, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustMobileEquipmentIdentity() string {
	v, _ := i.MobileEquipmentIdentity()
	return v
}
