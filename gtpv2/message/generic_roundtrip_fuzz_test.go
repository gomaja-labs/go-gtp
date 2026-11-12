package message_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

// canonicalGTPv2 returns b with the spare bits cleared exactly as the encoder
// writes them (TS 29.274 §§5.1, 5.3, 5.4, 8.2.1): octet 1 bits 2-1, the MP flag
// without a TEID, the spare octet (bits 4-1 of octet 12 when MP is set), and
// bits 8-5 of octet 4 of every IE, grouped children included. It walks a
// piggybacked message the same way. It returns nil if b is not well framed.
func canonicalGTPv2(b []byte) []byte {
	out := append([]byte(nil), b...)
	for off := 0; off < len(out); {
		m := out[off:]
		if len(m) < 8 {
			return nil
		}
		end := 4 + int(binary.BigEndian.Uint16(m[2:4]))
		hdr := 8
		if m[0]&0x08 != 0 {
			hdr = 12
		}
		if end < hdr || end > len(m) {
			return nil
		}
		m[0] &^= 0x03
		if m[0]&0x08 == 0 {
			m[0] &^= 0x04
		}
		if m[0]&0x04 != 0 {
			m[hdr-1] &= 0xf0
		} else {
			m[hdr-1] = 0
		}
		if !canonicalIEs(m[hdr:end], 0) {
			return nil
		}
		next := m[0]&0x10 != 0
		off += end
		if !next && off != len(out) {
			return nil
		}
	}
	return out
}

func canonicalIEs(b []byte, depth int) bool {
	for len(b) > 0 {
		if len(b) < 4 || depth > 8 {
			return false
		}
		n := 4 + int(binary.BigEndian.Uint16(b[1:3]))
		if n > len(b) {
			return false
		}
		b[3] &= 0x0f
		if (&ie.IE{Type: b[0]}).IsGrouped() && !canonicalIEs(b[4:n], depth+1) {
			return false
		}
		b = b[n:]
	}
	return true
}

func fuzzHex(s string) []byte {
	b, err := hex.DecodeString(strings.ReplaceAll(s, " ", ""))
	if err != nil {
		panic(err)
	}
	return b
}

// FuzzGenericRoundTrip: whatever ParseGeneric accepts and Marshal encodes must
// come back as the received datagram, except for the spare bits the encoder
// clears. This covers IE order, unknown and repeated IEs, grouped children and
// piggybacked messages (TS 29.274 §§5.5.1, 8.1, 8.2.1).
func FuzzGenericRoundTrip(f *testing.F) {
	for _, s := range []string{
		// Echo Request with Recovery.
		"40 01 00 09 00 00 2a 00 03 00 01 00 05",
		// Version Not Supported Indication.
		"40 03 00 04 00 00 2b 00",
		// Unknown message type 250 with unknown, type-254, Private Extension and
		// a Bearer Context holding an unknown child.
		"48 fa 00 51 99 99 99 99 00 00 2d 00 03 00 01 00 07 e6 00 03 03 aa bb cc" +
			"fe 00 05 00 01 00 de ad be ff 00 04 00 00 01 ca fe" +
			"5d 00 28 01 49 00 01 00 06 e7 00 02 05 01 02" +
			"57 00 19 02 c4 55 55 55 55 c0 00 02 03 20 01 0d b8 00 00 00 00 00 00 00 00 00 00 00 03",
		// Spare bits set in octet 1, octet 12 and an IE octet 4.
		"4b 20 00 14 00 00 00 00 00 00 32 0f 01 00 08 f0 00 01 01 21 43 65 87 f9",
		// Message priority 5.
		"4c 20 00 14 00 00 00 00 00 00 31 50 01 00 08 00 00 01 01 21 43 65 87 f9",
		// Create Session Response + Create Bearer Request, IEs out of table order.
		"58 21 00 1b 11 11 11 11 00 01 00 00 02 00 02 00 10 00 57 00 09 01 87 22 22 22 22 c0 00 02 01" +
			"48 5f 00 2a 33 33 33 33 00 01 01 00 5d 00 12 00 49 00 01 00 00 57 00 09 01 85 44 44 44 44 c0 00 02 02" +
			"e6 00 03 00 aa bb cc 49 00 01 00 05",
		// Create Bearer Response + Modify Bearer Request, IEs out of table order.
		"58 60 00 1d 55 55 55 55 00 01 02 00 02 00 02 00 10 00 5d 00 0b 00 49 00 01 00 05 02 00 02 00 10 00" +
			"48 22 00 28 66 66 66 66 00 01 03 00 e6 00 01 00 01" +
			"5d 00 12 00 49 00 01 00 05 57 00 09 00 80 77 77 77 77 c0 00 02 05 03 00 01 00 09",
	} {
		f.Add(fuzzHex(s))
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		g, err := message.ParseGeneric(b)
		if err != nil {
			return
		}
		out, err := g.Marshal()
		if err != nil {
			// Tolerated receive errors (invalid grouped content, a failed
			// piggybacked message) are refused on encode by design.
			return
		}
		want := canonicalGTPv2(b)
		if want == nil {
			t.Fatalf("ParseGeneric accepted a datagram the framing walk rejects\nin: % x", b)
		}
		if !bytes.Equal(out, want) {
			t.Fatalf("Generic round trip changed the datagram\n in: % x\nout: % x\nwant % x", b, out, want)
		}
	})
}
