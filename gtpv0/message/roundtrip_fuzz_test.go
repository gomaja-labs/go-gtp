package message_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv0/message"
)

// FuzzRoundTrip: whatever Parse accepts must Marshal without panicking, and the
// re-marshalled bytes must parse back to an identical encoding (stability).
func FuzzRoundTrip(f *testing.F) {
	f.Add([]byte{0x48, 0x20, 0x00, 0x08, 0, 0, 0, 0, 0, 0, 1, 0})
	f.Fuzz(func(t *testing.T, b []byte) {
		m, err := message.Parse(b)
		if err != nil {
			return
		}
		out := make([]byte, m.MarshalLen())
		if err := m.MarshalTo(out); err != nil {
			return
		}
		m2, err := message.Parse(out)
		if err != nil {
			t.Fatalf("re-parse of own output failed: %v\nin:  % x\nout: % x", err, b, out)
		}
		out2 := make([]byte, m2.MarshalLen())
		if err := m2.MarshalTo(out2); err != nil {
			t.Fatalf("re-marshal failed: %v", err)
		}
		if !bytes.Equal(out, out2) {
			t.Fatalf("unstable encoding\nin:   % x\nout1: % x\nout2: % x", b, out, out2)
		}
	})
}

// FuzzGenericRoundTrip requires exact replay of every structurally accepted frame.
func FuzzGenericRoundTrip(f *testing.F) {
	for _, tc := range genericStructuralCases {
		b, err := hex.DecodeString(tc.wire)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		g, err := message.ParseGeneric(b)
		if err != nil {
			return
		}
		out, err := g.Marshal()
		if err != nil {
			t.Fatalf("Marshal accepted input: %v; in %x", err, b)
		}
		if !bytes.Equal(out, b) {
			t.Fatalf("changed frame: in %x; out %x", b, out)
		}
		if g.MarshalLen() != len(out) {
			t.Fatalf("MarshalLen %d; wrote %d", g.MarshalLen(), len(out))
		}
	})
}
