// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"github.com/gomaja-labs/go-gtp/utils"
	"strings"
)

// NewAccessPointName creates a validated APN.
func NewAccessPointName(apn string) *IE {
	// TS 23.003 §9.1: at most 100 encoded octets, no empty or root label.
	if len(apn) == 0 || strings.HasSuffix(apn, ".") {
		return newErrorIE(AccessPointName, ErrMalformed)
	}
	if err := validateAPNName(apn); err != nil {
		return newErrorIE(AccessPointName, err)
	}
	b, err := utils.EncodeFQDN(apn)
	if err != nil {
		return newErrorIE(AccessPointName, err)
	}
	if len(b) > 100 {
		return newErrorIE(AccessPointName, ErrInvalidLength)
	}
	return New(AccessPointName, b)
}

// AccessPointName returns the decoded APN.
func (i *IE) AccessPointName() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if i.Type != AccessPointName {
		return "", &InvalidTypeError{Type: i.Type}
	}
	return utils.DecodeFQDN(i.Payload)
}

// MustAccessPointName returns the APN, or an empty string on error.
func (i *IE) MustAccessPointName() string { v, _ := i.AccessPointName(); return v }

func validateAPNName(name string) error {
	// TS 23.003 §9.1: ASCII letters, digits and interior hyphens only.
	for label := range strings.SplitSeq(name, ".") {
		if len(label) == 0 || len(label) > 63 {
			return ErrInvalidLength
		}
		for n := range len(label) {
			c := label[n]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
				continue
			}
			if c == '-' && n > 0 && n < len(label)-1 {
				continue
			}
			return ErrMalformed
		}
	}
	return nil
}
