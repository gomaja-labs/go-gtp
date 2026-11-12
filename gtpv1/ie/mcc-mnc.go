// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"

	"github.com/gomaja-labs/go-gtp/utils"
)

// MCC returns MCC value if type matches.
func (i *IE) MCC() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	switch i.Type {
	case RouteingAreaIdentity:
		if len(i.Payload) < 3 {
			return "", io.ErrUnexpectedEOF
		}
		mcc, _, err := utils.DecodePLMN(i.Payload[0:3])
		return mcc, err
	case UserLocationInformation:
		if len(i.Payload) < 4 {
			return "", io.ErrUnexpectedEOF
		}
		mcc, _, err := utils.DecodePLMN(i.Payload[1:4])
		return mcc, err
	default:
		return "", &InvalidTypeError{Type: i.Type}
	}
}

// MustMCC returns MCC in string if type matches.
// This should only be used if it is assured to have the value.
func (i *IE) MustMCC() string {
	v, _ := i.MCC()
	return v
}

// MNC returns MNC value if type matches.
func (i *IE) MNC() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	switch i.Type {
	case RouteingAreaIdentity:
		if len(i.Payload) < 3 {
			return "", io.ErrUnexpectedEOF
		}
		_, mnc, err := utils.DecodePLMN(i.Payload[0:3])
		return mnc, err
	case UserLocationInformation:
		if len(i.Payload) < 4 {
			return "", io.ErrUnexpectedEOF
		}
		_, mnc, err := utils.DecodePLMN(i.Payload[1:4])
		return mnc, err
	default:
		return "", &InvalidTypeError{Type: i.Type}
	}
}

// MustMNC returns MNC in string if type matches.
// This should only be used if it is assured to have the value.
func (i *IE) MustMNC() string {
	v, _ := i.MNC()
	return v
}
