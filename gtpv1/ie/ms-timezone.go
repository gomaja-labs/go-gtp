// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
	"math"
	"time"
)

// Timezone adjustment definitions.
const (
	TimeAdjustNoDaylightSaving uint8 = iota
	TimeAdjustDaylightSavingOneHour
	TimeAdjustDaylightSavingTwoHour
)

// NewMSTimeZone creates a new MSTimeZone IE.
func NewMSTimeZone(tz time.Duration, daylightSaving uint8) *IE {
	// TS 29.060 §7.7.52: integral quarter-hours and daylight values 0..2.
	if tz%(15*time.Minute) != 0 || tz > 79*15*time.Minute || tz < -79*15*time.Minute || daylightSaving > 2 {
		return newErrorIE(MSTimeZone, ErrMalformed)
	}
	i := New(MSTimeZone, make([]byte, 2))
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
	if i.Type != MSTimeZone {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) == 0 {
		return 0, io.ErrUnexpectedEOF
	}

	unsigned := i.Payload[0] & 0xf7
	dec := int((unsigned >> 4) + (unsigned&0x0f)*10)
	if (i.Payload[0]&0x08)>>3 == 1 {
		dec *= -1
	}

	return time.Duration(dec*15) * time.Minute, nil
}

// MustTimeZone returns TimeZone in time.Duration if type matches.
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
	if i.Type != MSTimeZone {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	v := i.Payload[1] & 3
	if v == 3 {
		return 0, ErrMalformed
	}
	return v, nil
}

// MustDaylightSaving returns DaylightSaving in uint8 if type matches.
// This should only be used if it is assured to have the value.
func (i *IE) MustDaylightSaving() uint8 {
	v, _ := i.DaylightSaving()
	return v
}
