package message_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
)

// https://github.com/secdev/scapy/commit/7bf257e6eb1a; TS 29.281 V19.2.0 §5.2.2.1.
func TestInteropUDPPortExtensionLength(t *testing.T) {
	b, err := message.NewExtensionHeader(0x40, []byte{0x08, 0x68}, 0).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{1, 8, 0x68, 0}) {
		t.Fatalf("extension = %x", b)
	}
}

// https://github.com/travelping/gtplib/commit/0c9e2bba47ef; TS 29.281 V19.2.0 §5.2.1.
func TestInteropExtensionChain(t *testing.T) {
	raw := []byte{0x34, 255, 0, 13, 0, 0, 0, 1, 0, 0, 0, 0x40, 1, 8, 0x68, 0xc0, 1, 0x12, 0x34, 0, 0x45}
	h, err := message.ParseHeader(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(h.ExtensionHeaders) != 2 || h.ExtensionHeaders[1].Type != 0xc0 || !bytes.Equal(h.Payload, []byte{0x45}) {
		t.Fatalf("header = %+v", h)
	}
	b, err := h.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, raw) {
		t.Fatalf("wire = %x", b)
	}
}

// https://github.com/wireshark/wireshark/commit/7557909a8395; TS 29.281 V19.2.0 §5.1.
func TestInteropTEIDWidth(t *testing.T) {
	h, err := message.ParseHeader([]byte{0x30, 255, 0, 1, 0x12, 0x34, 0x56, 0x78, 0x45})
	if err != nil {
		t.Fatal(err)
	}
	if h.TEID != 0x12345678 {
		t.Fatalf("TEID = %#x", h.TEID)
	}
}

// https://github.com/open5gs/open5gs/commit/3b6da3ddb2d0; TS 29.060 V19.0.0 §8.2.
func TestInteropControlSequence(t *testing.T) {
	for _, m := range []interface{ Marshal() ([]byte, error) }{message.NewEchoResponse(0x1234, ie.NewRecovery(1)), message.NewVersionNotSupported(0, 0x1234)} {
		b, err := m.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if len(b) < 12 || b[0]&2 == 0 || b[8] != 0x12 || b[9] != 0x34 {
			t.Fatalf("control wire = %x", b)
		}
	}
}

// https://github.com/edgecomllc/eupf/commit/24488216f22c; TS 29.281 V19.2.0 §5.1.
func TestInteropSequenceOnlyGPDU(t *testing.T) {
	h, err := message.ParseHeader([]byte{0x32, 255, 0, 5, 0, 0, 0, 1, 0x12, 0x34, 0, 0, 0x45})
	if err != nil {
		t.Fatal(err)
	}
	if h.SequenceNumber != 0x1234 || !bytes.Equal(h.Payload, []byte{0x45}) {
		t.Fatalf("header = %+v", h)
	}
}

// https://github.com/edgecomllc/eupf/commit/19593cf57282; TS 29.281 V19.2.0 §§5.1, 7.2.1.
func TestInteropEchoSequenceLength(t *testing.T) {
	b, err := message.NewEchoRequest(0x1234).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{0x32, 1, 0, 4, 0, 0, 0, 0, 0x12, 0x34, 0, 0}) {
		t.Fatalf("Echo wire = %x", b)
	}
}
