// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"errors"
	"io"
	"math"
	"time"
)

// NewEPCTimer creates a new EPCTimer IE.
func NewEPCTimer(duration time.Duration) *IE {
	raw, err := encodeTimer(duration)
	if err != nil {
		return newErrorIE(EPCTimer, err)
	}
	return NewEPCTimerRaw(raw)
}

// NewEPCTimerRaw creates a new EPCTimer IE from a uint8 value.
func NewEPCTimerRaw(duration uint8) *IE {
	return NewUint8IE(EPCTimer, duration)
}

// EPCTimer returns EPCTimer in time.Duration if the type of IE matches.
//
// The returned value is in time.Duration. To get the raw value as uint8,
// use ValueAsUint8 or access Payload field directly instead.
func (i *IE) EPCTimer() (time.Duration, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	return i.Timer()
}

// MustEPCTimer returns EPCTimer in time.Duration, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustEPCTimer() time.Duration {
	v, _ := i.EPCTimer()
	return v
}

// Timer returns Timer in time.Duration if the type of IE matches.
func (i *IE) Timer() (time.Duration, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if len(i.Payload) < 1 {
		return 0, io.ErrUnexpectedEOF
	}

	switch i.Type {
	case EPCTimer:
		var d time.Duration
		switch (i.Payload[0] & 0xe0) >> 5 {
		case 0x07:
			d = time.Duration(math.MaxInt64)
		case 0x04:
			d = time.Duration(i.Payload[0]&0x1f) * 10 * time.Hour
		case 0x03:
			d = time.Duration(i.Payload[0]&0x1f) * time.Hour
		case 0x02:
			d = time.Duration(i.Payload[0]&0x1f) * 10 * time.Minute
		case 0x01:
			d = time.Duration(i.Payload[0]&0x1f) * time.Minute
		case 0x00:
			d = time.Duration(i.Payload[0]&0x1f) * 2 * time.Second
		default:
			d = time.Duration(i.Payload[0]&0x1f) * time.Minute
		}
		return d, nil
	case OverloadControlInformation:
		return 0, errors.New("not implemented")
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// TS 29.274 §§8.85, 8.87: a three-bit unit and five-bit count.
func encodeTimer(d time.Duration) (uint8, error) {
	if d == time.Duration(math.MaxInt64) {
		return 0xe0, nil
	}
	if d < 0 {
		return 0, ErrMalformed
	}
	units := []time.Duration{2 * time.Second, time.Minute, 10 * time.Minute, time.Hour, 10 * time.Hour}
	for n := len(units) - 1; n >= 0; n-- {
		if d%units[n] == 0 && d/units[n] <= 31 {
			return uint8(n<<5) | uint8(d/units[n]), nil
		}
	}
	return 0, ErrMalformed
}
