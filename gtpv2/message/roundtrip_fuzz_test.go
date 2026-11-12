package message_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

// FuzzRoundTrip: decoding and encoding must never panic. When a tolerated
// received message also meets the stricter send rules, encoding must stabilize.
func FuzzRoundTrip(f *testing.F) {
	f.Add([]byte{0x48, 0x20, 0x00, 0x08, 0, 0, 0, 0, 0, 0, 1, 0})
	for typ := 0; typ < 256; typ++ {
		flags, size := byte(0x48), 12
		if typ >= 1 && typ <= 3 {
			flags, size = 0x40, 8
		}
		b := make([]byte, size)
		b[0], b[1] = flags, byte(typ)
		binary.BigEndian.PutUint16(b[2:4], uint16(size-4))
		f.Add(b)
		b = append(b, 250, 0, 1, 0, 0x42)
		binary.BigEndian.PutUint16(b[2:4], uint16(len(b)-4))
		f.Add(b)
	}
	f.Add([]byte{0x58, 33, 0, 8, 0, 0, 0, 0, 0, 0, 1, 0, 0x48, 95, 0, 8, 0, 0, 0, 0, 0, 0, 2, 0})
	// Create Session Response with a PGW Change Info (TS 29.274 §8.145) holding
	// an IP Address child, valid and with the child crossing the group end.
	f.Add([]byte{0x48, 33, 0, 26, 0, 0, 0, 1, 0, 0, 3, 0, 2, 0, 2, 0, 16, 0, 0xd6, 0, 8, 0, 74, 0, 4, 0, 192, 0, 2, 99})
	f.Add([]byte{0x48, 33, 0, 26, 0, 0, 0, 1, 0, 0, 3, 0, 2, 0, 2, 0, 16, 0, 0xd6, 0, 8, 0, 74, 0, 5, 0, 192, 0, 2, 99})
	f.Fuzz(func(t *testing.T, b []byte) {
		m, err := message.Parse(b)
		if err != nil {
			return
		}
		out := make([]byte, m.MarshalLen())
		if err := m.MarshalTo(out); err != nil {
			// TS 29.274 §§7.1.2, 7.7.1, 7.7.9 permit receiving erroneous
			// IEs which may not be transmitted. Unit tests check these policies.
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
