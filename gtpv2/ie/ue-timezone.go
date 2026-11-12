// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
	"math"
	"time"
)

// NewUETimeZone creates a new UETimeZone IE.
func NewUETimeZone(tz time.Duration, daylightSaving uint8) *IE {
	// TS 29.274 §8.44 and TS 24.008 §10.5.3.8: signed quarter-hours.
	if tz%(15*time.Minute) != 0 || tz < -79*15*time.Minute || tz > 79*15*time.Minute || daylightSaving > 2 {
		return newErrorIE(UETimeZone, ErrMalformed)
	}
	i := New(UETimeZone, 0x00, make([]byte, 2))
	min := tz.Minutes() / 15
	absMin := int(math.Abs(min))
	hex := byte(((absMin % 10) << 4) | (absMin / 10))
	if min < 0 {
		hex |= 0x08
	}
	i.Payload[0] = hex
	i.Payload[1] = daylightSaving & 0x03

	return i
}

// TimeZone returns TimeZone in time.Duration if the type of IE matches.
func (i *IE) TimeZone() (time.Duration, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != UETimeZone {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 1 {
		return 0, io.ErrUnexpectedEOF
	}

	unsigned := i.Payload[0] & 0xf7
	dec := int((unsigned >> 4) + (unsigned&0x0f)*10)
	if (i.Payload[0]&0x08)>>3 == 1 {
		dec *= -1
	}

	return time.Duration(dec*15) * time.Minute, nil
}

// MustTimeZone returns TimeZone in time.Duration, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustTimeZone() time.Duration {
	v, _ := i.TimeZone()
	return v
}

// DaylightSaving returns DaylightSaving in uint8 if the type of IE matches.
func (i *IE) DaylightSaving() (uint8, error) {
	if err := i.Err(); err != nil {
		return 0, err
	}
	if i.Type != UETimeZone {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	return i.Payload[1] & 3, nil
}

// MustDaylightSaving returns DaylightSaving in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustDaylightSaving() uint8 {
	if i == nil {
		return 0
	}
	v, _ := i.DaylightSaving()
	return v
}
