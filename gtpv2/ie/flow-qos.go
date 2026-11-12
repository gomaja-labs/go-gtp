// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"

	"github.com/gomaja-labs/go-gtp/utils"
)

// NewFlowQoS creates a new FlowQoS IE.
func NewFlowQoS(qci uint8, umbr, dmbr, ugbr, dgbr uint64) *IE {
	v := NewFlowQoSFields(qci, umbr, dmbr, ugbr, dgbr)
	b, err := v.Marshal()
	if err != nil {
		return newErrorIE(FlowQoS, err)
	}

	return New(FlowQoS, 0x00, b)
}

// FlowQoS returns FlowQoS in FlowQoSFields type if the type of IE matches.
func (i *IE) FlowQoS() (*FlowQoSFields, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	switch i.Type {
	case FlowQoS:
		return ParseFlowQoSFields(i.Payload)
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// FlowQoSFields is a set of fields in FlowQoS IE.
type FlowQoSFields struct {
	QCI                          uint8
	MaximumBitRateForUplink      uint64 // 40 bits
	MaximumBitRateForDownlink    uint64 // 40 bits
	GuaranteedBitRateForUplink   uint64 // 40 bits
	GuaranteedBitRateForDownlink uint64 // 40 bits
}

// NewFlowQoSFields creates a new FlowQoSFields.
func NewFlowQoSFields(qci uint8, umbr, dmbr, ugbr, dgbr uint64) *FlowQoSFields {
	return &FlowQoSFields{
		QCI:                          qci,
		MaximumBitRateForUplink:      umbr,
		MaximumBitRateForDownlink:    dmbr,
		GuaranteedBitRateForUplink:   ugbr,
		GuaranteedBitRateForDownlink: dgbr,
	}
}

// Marshal serializes FlowQoSFields.
func (f *FlowQoSFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo serializes FlowQoSFields.
func (f *FlowQoSFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if len(b) < f.MarshalLen() {
		return ErrTooShortToMarshal
	}
	if len(b) < 21 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.QCI

	encoded, err := utils.Uint64To40(f.MaximumBitRateForUplink)
	if err != nil {
		return err
	}
	copy(b[1:6], encoded)
	encoded, err = utils.Uint64To40(f.MaximumBitRateForDownlink)
	if err != nil {
		return err
	}
	copy(b[6:11], encoded)
	encoded, err = utils.Uint64To40(f.GuaranteedBitRateForUplink)
	if err != nil {
		return err
	}
	copy(b[11:16], encoded)
	encoded, err = utils.Uint64To40(f.GuaranteedBitRateForDownlink)
	if err != nil {
		return err
	}
	copy(b[16:21], encoded)

	return nil
}

// ParseFlowQoSFields decodes FlowQoSFields.
func ParseFlowQoSFields(b []byte) (*FlowQoSFields, error) {
	f := &FlowQoSFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return f, nil
}

// UnmarshalBinary decodes given bytes into FlowQoSFields.
func (f *FlowQoSFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = FlowQoSFields{}
	if len(b) < 21 {
		return io.ErrUnexpectedEOF
	}

	f.QCI = b[0]

	f.MaximumBitRateForUplink = utils.Uint40To64(b[1:6])
	f.MaximumBitRateForDownlink = utils.Uint40To64(b[6:11])
	f.GuaranteedBitRateForUplink = utils.Uint40To64(b[11:16])
	f.GuaranteedBitRateForDownlink = utils.Uint40To64(b[16:21])

	return nil
}

// MarshalLen returns the serial length of FlowQoSFields in int.
func (f *FlowQoSFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 21
}
