// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// NewExtendedCommonFlagsII creates a new ExtendedCommonFlagsII IE.
//
// Note: each flag should be set in 1 or 0.
func NewExtendedCommonFlagsII(pmtsmi, dtci, pnsi int) *IE {
	if pmtsmi < 0 || pmtsmi > 1 || dtci < 0 || dtci > 1 || pnsi < 0 || pnsi > 1 {
		return newErrorIE(ExtendedCommonFlagsII, ErrMalformed)
	}
	return New(
		ExtendedCommonFlagsII,
		[]byte{uint8(
			pmtsmi<<2 | dtci<<1 | pnsi,
		)},
	)
}

// ExtendedCommonFlagsII returns ExtendedCommonFlagsII value if type matches.
func (i *IE) ExtendedCommonFlagsII() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != ExtendedCommonFlagsII {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) == 0 {
		return 0, io.ErrUnexpectedEOF
	}

	return i.Payload[0], nil
}

// MustExtendedCommonFlagsII returns ExtendedCommonFlagsII in uint8 if type matches.
// This should only be used if it is assured to have the value.
func (i *IE) MustExtendedCommonFlagsII() uint8 {
	v, _ := i.ExtendedCommonFlagsII()
	return v
}

// IsPMTSMI checks if PMTSMI flag exists in ExtendedCommonFlagsII.
func (i *IE) IsPMTSMI() bool {
	if i.Err() != nil {
		return false
	}
	return ((i.MustExtendedCommonFlagsII() >> 2) & 0x01) != 0
}

// IsDTCI checks if DTCI flag exists in ExtendedCommonFlagsII.
func (i *IE) IsDTCI() bool {
	if i.Err() != nil {
		return false
	}
	return ((i.MustExtendedCommonFlagsII() >> 1) & 0x01) != 0
}

// IsPNSI checks if PNSI flag exists in ExtendedCommonFlagsII.
func (i *IE) IsPNSI() bool {
	if i.Err() != nil {
		return false
	}
	return (i.MustExtendedCommonFlagsII() & 0x01) != 0
}
