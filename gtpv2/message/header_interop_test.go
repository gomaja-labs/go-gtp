package message_test

import (
	"encoding/binary"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

// https://github.com/secdev/scapy/commit/0260f9705a85; TS 29.274 V20.0.0 §5.5.1.
func TestV2HeaderLength(t *testing.T) {
	for _, flags := range []uint8{0x40, 0x48} {
		h := message.NewHeader(flags, map[uint8]uint8{0x40: 1, 0x48: 32}[flags], 0x12345678, 0xabcdef, []byte{3, 0, 1, 0, 1})
		b, err := h.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if int(binary.BigEndian.Uint16(b[2:4])) != len(b)-4 {
			t.Fatalf("header = %x", b)
		}
		p, err := message.ParseHeader(b)
		if err != nil || p.SequenceNumber != 0xabcdef {
			t.Fatalf("parse = %+v, %v", p, err)
		}
	}
}
