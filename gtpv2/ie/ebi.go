// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewEPSBearerID creates a new EPSBearerID IE.
func NewEPSBearerID(ebi uint8) *IE {
	// TS 29.274 §8.8: validate the transmitted value.
	if ebi > 15 {
		return newErrorIE(EPSBearerID, ErrMalformed)
	}
	return NewUint8IE(EPSBearerID, ebi&0x0f)
}

// EPSBearerID returns EPSBearerID if the type of IE matches.
func (i *IE) EPSBearerID() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	switch i.Type {
	case EPSBearerID:
		v, err := i.ValueAsUint8()
		return v & 0x0f, err
	case BearerContext:
		ies, err := i.BearerContext()
		if err != nil {
			return 0, err
		}

		// TS 29.274 §7.7.9: match the child type and instance together.
		for _, child := range ies {
			if child != nil && child.Type == EPSBearerID && child.Instance() == 0 {
				return child.EPSBearerID()
			}
		}
		return 0, ErrIENotFound
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// MustEPSBearerID returns EPSBearerID in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustEPSBearerID() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.EPSBearerID()
	return v
}
