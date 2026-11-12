package ie_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

func TestTypeExtensionBoundsAndOpaqueContent(t *testing.T) {
	for _, ext := range []uint16{0, 1, 255, 256, 65535} {
		i := ie.NewWithTypeExtension(ext, 15, []byte{0xde, 0xad, 0xbe})
		if ext < 256 {
			if _, err := i.Marshal(); !errors.Is(err, ie.ErrMalformed) {
				t.Errorf("constructor extension %d: %v", ext, err)
			}
		} else {
			want := []byte{254, 0, 5, 15, byte(ext >> 8), byte(ext), 0xde, 0xad, 0xbe}
			b, err := i.Marshal()
			if err != nil || !bytes.Equal(b, want) {
				t.Fatalf("extension %d: %x %v", ext, b, err)
			}
			if got, err := i.TypeExtension(); err != nil || got != ext {
				t.Errorf("accessor %d %v", got, err)
			}
			for n := 0; n < i.MarshalLen(); n++ {
				if err := i.MarshalTo(make([]byte, n)); !errors.Is(err, ie.ErrTooShortToMarshal) {
					t.Fatalf("short destination %d: %v", n, err)
				}
			}
		}
		raw := ie.New(254, 0, []byte{byte(ext >> 8), byte(ext)})
		b, err := raw.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := ie.Parse(b)
		if err != nil {
			t.Fatal(err)
		}
		got, err := decoded.TypeExtension()
		if ext < 256 {
			if !errors.Is(err, ie.ErrMalformed) {
				t.Errorf("received extension %d: %v", ext, err)
			}
		} else if err != nil || got != ext {
			t.Errorf("received extension %d: %d %v", ext, got, err)
		}
	}
	for _, i := range []*ie.IE{ie.New(254, 0, nil), ie.New(254, 0, []byte{1})} {
		if _, err := i.TypeExtension(); !errors.Is(err, ie.ErrInvalidLength) {
			t.Fatalf("short extension: %v", err)
		}
	}
	if _, err := (*ie.IE)(nil).TypeExtension(); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	if _, err := ie.NewRecovery(1).TypeExtension(); !errors.Is(err, ie.ErrInvalidType) {
		t.Fatal(err)
	}
	if _, err := ie.NewWithTypeExtension(256, 16, nil).Marshal(); !errors.Is(err, ie.ErrMalformed) {
		t.Fatal(err)
	}
	for _, n := range []int{65533, 65534} {
		i := ie.NewWithTypeExtension(256, 0, make([]byte, n))
		b, err := i.Marshal()
		if n == 65533 {
			if err != nil || len(b) != 65539 {
				t.Fatalf("maximum extension content: %d %v", len(b), err)
			}
		} else if !errors.Is(err, ie.ErrInvalidLength) {
			t.Fatal(err)
		}
	}
}

func TestUnsupportedExtensionStaysInAdditionalIEs(t *testing.T) {
	// TS 29.274 Table 8.1-1 NOTE 2: unsupported extensions remain opaque.
	for _, payload := range [][]byte{nil, {0}, {0, 255}, {1, 0, 0xde, 0xad, 0xbe}, {255, 255}} {
		raw := ie.New(254, 0, payload)
		g := message.NewGeneric(message.MsgTypeEchoRequest, 0, 1, ie.NewRecovery(7), raw, raw)
		b, err := g.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		m, err := message.ParseEchoRequest(b)
		if err != nil {
			t.Fatal(err)
		}
		if len(m.AdditionalIEs) != 2 || !bytes.Equal(m.AdditionalIEs[0].Payload, payload) {
			t.Fatalf("extended IEs lost: %+v", m.AdditionalIEs)
		}
		decoded, err := message.ParseGeneric(b)
		if err != nil {
			t.Fatal(err)
		}
		out, err := decoded.Marshal()
		if err != nil || !bytes.Equal(out, b) {
			t.Fatalf("opaque extension changed: %x %v", out, err)
		}
	}
}
