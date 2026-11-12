package ie

import (
	"encoding/binary"
	"errors"
	"testing"
	"time"
)

func TestTimestampUnknown(t *testing.T) {
	got, err := newUint32ValIE(ULITimestamp, 0).Timestamp()
	if err != nil || !got.IsZero() {
		t.Fatalf("zero timestamp = %v, %v; want unknown", got, err)
	}
	if err := NewULITimestamp(time.Time{}).Err(); err != nil {
		t.Fatalf("encode unknown: %v", err)
	}
}

func TestTimestampNear(t *testing.T) {
	ref2026 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	ref2040 := time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name      string
		wire      uint32
		ref, want time.Time
	}{
		{"1968", 0x7fe6c600, ref2026, time.Date(1968, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"2026", 0xed003780, ref2026, ref2026},
		{"era0_end", 0xffffffff, ref2040, time.Date(2036, 2, 7, 6, 28, 15, 0, time.UTC)},
		{"era1_start_unknown", 0, ref2040, time.Time{}},
		{"era1_first_second", 1, ref2026, time.Date(2036, 2, 7, 6, 28, 17, 0, time.UTC)},
		{"2040_near_2026", 0x0754fd00, ref2026, ref2040},
		{"2040_near_2040", 0x0754fd00, ref2040, ref2040},
		{"same_wire_era0", 0x0754fd00, time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC), ref2040.Add(-0x100000000 * time.Second)},
		{"unknown", 0, ref2026, time.Time{}},
		{"distant_past", 1, time.Time{}, time.Unix(1-2208988800, 0)},
		{"distant_future", 1, time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC), time.Unix(0x100000001-2208988800, 0)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			i := newUint32ValIE(ULITimestamp, tc.wire)
			got, err := i.TimestampNear(tc.ref)
			if err != nil || !got.Equal(tc.want) {
				t.Fatalf("got %v, %v; want %v", got, err, tc.want)
			}
		})
	}
	eraZero := time.Unix(1-2208988800, 0)
	midpoint := eraZero.Add(0x80000000 * time.Second)
	for _, tc := range []struct{ ref, want time.Time }{
		{midpoint.Add(-time.Nanosecond), eraZero},
		{midpoint, eraZero},
		{midpoint.Add(time.Nanosecond), eraZero.Add(0x100000000 * time.Second)},
	} {
		got, err := newUint32ValIE(ULITimestamp, 1).TimestampNear(tc.ref)
		if err != nil || !got.Equal(tc.want) {
			t.Fatalf("midpoint %v: %v, %v; want %v", tc.ref, got, err, tc.want)
		}
	}
}

func TestTimestampEncoding(t *testing.T) {
	start := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	boundary := time.Date(2036, 2, 7, 6, 28, 16, 0, time.UTC)
	end := start.Add(0x200000000 * time.Second)
	for _, tc := range []struct {
		name    string
		date    time.Time
		wire    uint32
		invalid bool
	}{
		{"unknown", time.Time{}, 0, false},
		{"epoch", start, 0, false},
		{"before_epoch", start.Add(-time.Nanosecond), 0, true},
		{"1968", time.Date(1968, 1, 1, 0, 0, 0, 0, time.UTC), 0x7fe6c600, false},
		{"2026", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), 0xed003780, false},
		{"era0_end", boundary.Add(-time.Second), 0xffffffff, false},
		{"era1_start", boundary, 0, false},
		{"2040", time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC), 0x0754fd00, false},
		{"fraction", boundary.Add(2*time.Second + time.Nanosecond), 2, false},
		{"era1_end", end.Add(-time.Nanosecond), 0xffffffff, false},
		{"era2_start", end, 0, true},
		{"far_future", time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC), 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			i := NewULITimestamp(tc.date)
			_, err := i.Marshal()
			if tc.invalid {
				if !errors.Is(err, ErrInvalidLength) {
					t.Fatalf("invalid time marshal: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(i.Payload) != 4 || binary.BigEndian.Uint32(i.Payload) != tc.wire {
				t.Fatalf("wire = %x; want %08x", i.Payload, tc.wire)
			}
			b := make([]byte, i.MarshalLen())
			if err := i.MarshalTo(b); err != nil {
				t.Fatal(err)
			}
			if err := i.MarshalTo(b[:len(b)-1]); err == nil {
				t.Fatal("short buffer accepted")
			}
		})
	}
}

func TestTimestampErrors(t *testing.T) {
	ref := time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC)
	var nilIE *IE
	if _, err := nilIE.TimestampNear(ref); !errors.Is(err, ErrIENotFound) {
		t.Fatalf("nil: %v", err)
	}
	if _, err := nilIE.Timestamp(); !errors.Is(err, ErrIENotFound) {
		t.Fatalf("nil default: %v", err)
	}
	if !nilIE.MustTimestamp().IsZero() {
		t.Fatal("nil MustTimestamp")
	}
	for n := 0; n < 4; n++ {
		i := newUint32ValIE(ULITimestamp, 1)
		i.Payload = i.Payload[:n]
		if _, err := i.TimestampNear(ref); !errors.Is(err, ErrInvalidLength) {
			t.Fatalf("length %d: %v", n, err)
		}
	}
	i := newUint32ValIE(ULITimestamp, 1)
	i.Type = 0
	var typeErr *InvalidTypeError
	if _, err := i.TimestampNear(ref); !errors.As(err, &typeErr) {
		t.Fatalf("type: %v", err)
	}
}
