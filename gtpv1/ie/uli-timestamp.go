// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// NewULITimestamp creates a new ULITimestamp IE. The zero time means unknown.
// Times outside NTP eras 0 and 1 are rejected; fractional seconds are discarded.
func NewULITimestamp(ts time.Time) *IE {
	// TS 29.060 §7.7.114: first four NTP timestamp octets (RFC 5905 §6).
	if ts.IsZero() {
		return newUint32ValIE(ULITimestamp, 0)
	}
	if ts.Before(time.Unix(-2208988800, 0)) || !ts.Before(time.Unix(0x200000000-2208988800, 0)) {
		return newErrorIE(ULITimestamp, ErrInvalidLength)
	}
	return newUint32ValIE(ULITimestamp, uint32((ts.Unix()+2208988800)%0x100000000))
}

// Timestamp returns the timestamp in NTP era 0 or 1 nearest the current time.
// Use TimestampNear with an explicit reference when decoding historical records.
func (i *IE) Timestamp() (time.Time, error) {
	return i.TimestampNear(time.Now())
}

// MustTimestamp returns Timestamp in time.Time, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustTimestamp() time.Time {
	v, _ := i.Timestamp()
	return v
}

// TimestampNear returns the timestamp in NTP era 0 or 1 nearest ref.
// An exact tie selects era 0. Wire zero returns the zero time (unknown).
func (i *IE) TimestampNear(ref time.Time) (time.Time, error) {
	if err := i.Err(); err != nil {
		return time.Time{}, err
	}
	if i.Type != ULITimestamp {
		return time.Time{}, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 4 {
		return time.Time{}, io.ErrUnexpectedEOF
	}
	return ntpTimestampNear(binary.BigEndian.Uint32(i.Payload), ref), nil
}

func ntpTimestampNear(seconds uint32, ref time.Time) time.Time {
	// RFC 5905 §6: zero means unknown; the era must come from external means.
	if seconds == 0 {
		return time.Time{}
	}
	eraZero := time.Unix(int64(seconds)-2208988800, 0).UTC()
	// Compare with the midpoint to avoid duration overflow for distant references.
	if ref.After(eraZero.Add(0x80000000 * time.Second)) {
		return eraZero.Add(0x100000000 * time.Second)
	}
	return eraZero
}
