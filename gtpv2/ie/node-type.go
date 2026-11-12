// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewNodeType creates a new NodeType IE.
func NewNodeType(nodeType uint8) *IE {
	// TS 29.274 §8.65: validate the transmitted value.
	if nodeType > 1 {
		return newErrorIE(NodeType, ErrMalformed)
	}
	return NewUint8IE(NodeType, nodeType)
}

// NodeType returns NodeType in uint8 if the type of IE matches.
func (i *IE) NodeType() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != NodeType {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	v, err := i.ValueAsUint8()
	if err != nil {
		return 0, err
	}
	if v > 1 {
		return 0, ErrMalformed
	}
	return v, nil
}

// MustNodeType returns NodeType in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustNodeType() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.NodeType()
	return v
}
