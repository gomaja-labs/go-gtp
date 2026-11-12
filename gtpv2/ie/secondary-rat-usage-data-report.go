// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"io"
	"time"
)

// SecondaryRATUsageDataReportFields contains the usage counters or an NG-RAN transfer.
// Timestamps hold the first four NTP timestamp octets, as specified by TS 29.274 §8.132.
// They remain raw on decoding; use StartTimeNear and EndTimeNear with the same
// reference time to resolve their eras. Zero represents unknown time.
type SecondaryRATUsageDataReportFields struct {
	Flags                               uint8
	SecondaryRATType                    uint8
	EPSBearerID                         uint8
	StartTimestamp                      uint32
	EndTimestamp                        uint32
	UsageDataDL                         uint64
	UsageDataUL                         uint64
	SecondaryRATDataUsageReportTransfer []byte
}

// StartTimeNear returns the recording start in NTP era 0 or 1 nearest ref.
// An exact tie selects era 0. Wire zero returns the zero time (unknown).
func (f *SecondaryRATUsageDataReportFields) StartTimeNear(ref time.Time) (time.Time, error) {
	if f == nil {
		return time.Time{}, ErrIENotFound
	}
	// TS 29.274 §8.132: start/end use the NTP seconds format (RFC 5905 §6).
	return ntpTimestampNear(f.StartTimestamp, ref), nil
}

// EndTimeNear returns the recording end in NTP era 0 or 1 nearest ref.
// An exact tie selects era 0. Wire zero returns the zero time (unknown).
func (f *SecondaryRATUsageDataReportFields) EndTimeNear(ref time.Time) (time.Time, error) {
	if f == nil {
		return time.Time{}, ErrIENotFound
	}
	return ntpTimestampNear(f.EndTimestamp, ref), nil
}

// NewSecondaryRATUsageDataReport creates a Secondary RAT Usage Data Report IE.
func NewSecondaryRATUsageDataReport(f *SecondaryRATUsageDataReportFields) *IE {
	b, err := f.Marshal()
	if err != nil {
		return newErrorIE(SecondaryRATUsageDataReport, err)
	}
	return New(SecondaryRATUsageDataReport, 0, b)
}

// SecondaryRATUsageDataReport decodes a Secondary RAT Usage Data Report IE.
func (i *IE) SecondaryRATUsageDataReport() (*SecondaryRATUsageDataReportFields, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	if i.Type != SecondaryRATUsageDataReport {
		return nil, &InvalidTypeError{Type: i.Type}
	}
	f := &SecondaryRATUsageDataReportFields{}
	if err := f.UnmarshalBinary(i.Payload); err != nil {
		return nil, err
	}
	return f, nil
}

// Marshal serializes the report fields.
func (f *SecondaryRATUsageDataReportFields) Marshal() ([]byte, error) {
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes the report fields into b.
func (f *SecondaryRATUsageDataReportFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if len(b) < f.MarshalLen() {
		return ErrTooShortToMarshal
	}
	if f.Flags > 7 || f.SecondaryRATType > 1 || f.EPSBearerID > 15 {
		return ErrMalformed
	}
	// TS 29.274 §8.132: SRUDN requires IRPGW and zero legacy usage fields.
	if f.Flags&4 != 0 {
		if f.Flags&1 == 0 || f.StartTimestamp != 0 || f.EndTimestamp != 0 || f.UsageDataDL != 0 || f.UsageDataUL != 0 || len(f.SecondaryRATDataUsageReportTransfer) > 255 {
			return ErrMalformed
		}
	} else if len(f.SecondaryRATDataUsageReportTransfer) != 0 {
		return ErrMalformed
	}
	b[0], b[1], b[2] = f.Flags, f.SecondaryRATType, f.EPSBearerID
	binary.BigEndian.PutUint32(b[3:7], f.StartTimestamp)
	binary.BigEndian.PutUint32(b[7:11], f.EndTimestamp)
	binary.BigEndian.PutUint64(b[11:19], f.UsageDataDL)
	binary.BigEndian.PutUint64(b[19:27], f.UsageDataUL)
	if f.Flags&4 != 0 {
		b[27] = uint8(len(f.SecondaryRATDataUsageReportTransfer))
		copy(b[28:], f.SecondaryRATDataUsageReportTransfer)
	}
	return nil
}

// UnmarshalBinary replaces the report with decoded fields.
func (f *SecondaryRATUsageDataReportFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = SecondaryRATUsageDataReportFields{}
	if len(b) < 27 {
		return io.ErrUnexpectedEOF
	}
	f.Flags, f.SecondaryRATType, f.EPSBearerID = b[0]&7, b[1], b[2]&15
	if f.Flags&4 != 0 {
		if len(b) < 28 || int(b[27]) > len(b)-28 {
			return io.ErrUnexpectedEOF
		}
		f.SecondaryRATDataUsageReportTransfer = append([]byte(nil), b[28:28+int(b[27])]...)
	} else {
		f.StartTimestamp = binary.BigEndian.Uint32(b[3:7])
		f.EndTimestamp = binary.BigEndian.Uint32(b[7:11])
		f.UsageDataDL = binary.BigEndian.Uint64(b[11:19])
		f.UsageDataUL = binary.BigEndian.Uint64(b[19:27])
	}
	return nil
}

// MarshalLen returns the serialized report length.
func (f *SecondaryRATUsageDataReportFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	if f.Flags&4 != 0 {
		return 28 + len(f.SecondaryRATDataUsageReportTransfer)
	}
	return 27
}
