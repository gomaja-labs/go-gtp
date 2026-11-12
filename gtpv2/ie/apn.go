// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"strings"

	"github.com/gomaja-labs/go-gtp/utils"
)

// NewAccessPointName creates a new AccessPointName IE.
func NewAccessPointName(apn string) *IE {
	// TS 23.003 §9.1: at most 100 encoded octets; no root label.
	if err := validateAPNName(apn); err != nil {
		return newErrorIE(AccessPointName, err)
	}
	i := NewFQDNIE(AccessPointName, apn)
	if len(i.Payload) > 100 {
		return newErrorIE(AccessPointName, ErrInvalidLength)
	}
	return i
}

// AccessPointName returns AccessPointName in string if the type of IE matches.
func (i *IE) AccessPointName() (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if i.Type != AccessPointName {
		return "", &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) == 0 || len(i.Payload) > 100 {
		return "", ErrInvalidLength
	}
	// TS 23.003 §9.1: APNs have no zero-length root label or DNS compression.
	for offset := 0; offset < len(i.Payload); {
		n := int(i.Payload[offset])
		offset++
		if n == 0 || n > 63 || n > len(i.Payload)-offset {
			return "", ErrInvalidLength
		}
		offset += n
	}
	name, err := utils.DecodeFQDN(i.Payload)
	if err != nil {
		return "", err
	}
	if err := validateAPNName(name); err != nil {
		return "", err
	}
	return name, nil
}

// MustAccessPointName returns AccessPointName in string, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustAccessPointName() string {
	if i == nil {
		return ""
	}
	v, _ := i.AccessPointName()
	return v
}

func validateAPNName(name string) error {
	// TS 23.003 §9.1: ASCII letters, digits and interior hyphens only.
	for label := range strings.SplitSeq(name, ".") {
		if len(label) == 0 || len(label) > 63 {
			return ErrInvalidLength
		}
		for n := range len(label) {
			c := label[n]
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
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
