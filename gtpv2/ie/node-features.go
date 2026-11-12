// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewNodeFeatures creates a new NodeFeatures IE.
func NewNodeFeatures(flags uint8) *IE {
	return NewUint8IE(NodeFeatures, flags)
}

// NodeFeatures returns NodeFeatures in uint8 if the type of IE matches.
func (i *IE) NodeFeatures() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != NodeFeatures {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	return i.ValueAsUint8()
}

// MustNodeFeatures returns NodeFeatures in uint8 if the type of IE matches.
func (i *IE) MustNodeFeatures() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.NodeFeatures()
	return v
}

// HasPRN reports whether an IE has PRN bit.
func (i *IE) HasPRN() bool {
	if i == nil {
		return false
	}
	v, err := i.NodeFeatures()
	if err != nil {
		return false
	}

	return has1stBit(v)
}

// HasMABR reports whether an IE has MABR bit.
func (i *IE) HasMABR() bool {
	if i == nil {
		return false
	}
	v, err := i.NodeFeatures()
	if err != nil {
		return false
	}

	return has2ndBit(v)
}

// HasNTSR reports whether an IE has NTSR bit.
func (i *IE) HasNTSR() bool {
	if i == nil {
		return false
	}
	v, err := i.NodeFeatures()
	if err != nil {
		return false
	}

	return has3rdBit(v)
}

// HasCIOT reports whether an IE has CIOT bit.
func (i *IE) HasCIOT() bool {
	if i == nil {
		return false
	}
	v, err := i.NodeFeatures()
	if err != nil {
		return false
	}

	return has4thBit(v)
}

// HasS1UN reports whether an IE has S1UN bit.
func (i *IE) HasS1UN() bool {
	if i == nil {
		return false
	}
	v, err := i.NodeFeatures()
	if err != nil {
		return false
	}

	return has5thBit(v)
}

// HasETH reports whether an IE has ETH bit.
func (i *IE) HasETH() bool {
	if i == nil {
		return false
	}
	v, err := i.NodeFeatures()
	if err != nil {
		return false
	}

	return has6thBit(v)
}

// HasMTEDT reports whether an IE has MTEDT bit.
func (i *IE) HasMTEDT() bool {
	if i == nil {
		return false
	}
	v, err := i.NodeFeatures()
	if err != nil {
		return false
	}

	return has7thBit(v)
}

// HasPSSET reports PGW-C/SMF Set support (TS 29.274 Table 8.83-1, octet 5 bit 8).
func (i *IE) HasPSSET() bool {
	v, err := i.NodeFeatures()
	return err == nil && v&0x80 != 0
}
