// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"

	"github.com/gomaja-labs/go-gtp/utils"
)

// NewServingNetwork creates a ServingNetwork IE.
func NewServingNetwork(mcc, mnc string) *IE {
	encoded, err := utils.EncodePLMN(mcc, mnc)
	if err != nil {
		return newErrorIE(ServingNetwork, err)
	}

	return New(ServingNetwork, 0x00, encoded)
}

// ServingNetwork returns ServingNetwork(MCC and MNC) in string if the type of IE matches.
func (i *IE) ServingNetwork() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if i.Type != ServingNetwork {
		return "", &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 3 {
		return "", io.ErrUnexpectedEOF
	}

	mcc, mnc, err := utils.DecodePLMN(i.Payload)
	if err != nil {
		return "", err
	}

	return mcc + mnc, nil
}

// MustServingNetwork returns ServingNetwork in string, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustServingNetwork() string {
	if i == nil {
		return ""
	}
	v, _ := i.ServingNetwork()
	return v
}
