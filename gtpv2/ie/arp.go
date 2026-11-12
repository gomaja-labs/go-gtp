// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewAllocationRetensionPriority creates a new AllocationRetensionPriority IE.
func NewAllocationRetensionPriority(pci, pl, pvi uint8) *IE {
	// TS 29.274 §8.86: reject fields that do not fit.
	if pci > 1 || pl > 15 || pvi > 1 {
		return newErrorIE(AllocationRetensionPriority, ErrMalformed)
	}
	i := New(AllocationRetensionPriority, 0x00, make([]byte, 1))
	i.Payload[0] |= (pci << 6 & 0x40) | (pl << 2 & 0x3c) | (pvi & 0x01)
	return i
}

// AllocationRetensionPriority returns AllocationRetensionPriority in uint8 if the type of IE matches.
func (i *IE) AllocationRetensionPriority() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != AllocationRetensionPriority {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	return i.ValueAsUint8()
}

// HasPVI reports whether an IE has PVI bit.
func (i *IE) HasPVI() bool {
	if i == nil {
		return false
	}
	v, err := i.AllocationRetensionPriority()
	if err != nil {
		return false
	}

	return has1stBit(v)
}

// HasPCI reports whether an IE has PCI bit.
func (i *IE) HasPCI() bool {
	if i == nil {
		return false
	}
	v, err := i.AllocationRetensionPriority()
	if err != nil {
		return false
	}

	return has7thBit(v)
}

// PriorityLevel returns PriorityLevel in uint8 if the type of IE matches.
func (i *IE) PriorityLevel() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	switch i.Type {
	case AllocationRetensionPriority:
		v, err := i.AllocationRetensionPriority()
		if err != nil {
			return 0, err
		}

		return (v & 0x3c) >> 2, nil
	case BearerQoS:
		v, err := i.BearerQoS()
		if err != nil {
			return 0, err
		}

		return (v.ARP & 0x3c) >> 2, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
