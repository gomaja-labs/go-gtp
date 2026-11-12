package message_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

var genericDatagrams = []string{
	"5821001b1111111100010000020002001000570009018722222222c0000201485f002a33333333000101005d0012004900010000570009018544444444c0000202e6000300aabbcc4900010005",
	"4821001a1111111100010000020002001000d60008004a000400c0000263",
	"4821000d00000001000001005d000100ff",
	"4b21000d00000001000001ff030001f007",
	"582100080000000100000100485f00080000000200000200",
	"582100080000000100000100482200080000000200000200",
	"582100080000000100000100485f00090000000200000200ff",
}

func TestGenericRetainsDatagramBytes(t *testing.T) {
	for _, s := range genericDatagrams {
		t.Run(s, func(t *testing.T) {
			in, err := hex.DecodeString(s)
			if err != nil {
				t.Fatal(err)
			}
			g, err := message.ParseGeneric(in)
			if err != nil {
				t.Fatal(err)
			}
			if g.Piggybacked != nil {
				if _, ok := g.Piggybacked.(*message.Generic); !ok {
					t.Fatalf("secondary type %T", g.Piggybacked)
				}
			}
			out, err := g.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(in, out) {
				t.Fatalf("datagram changed\nin: %x\nout:%x", in, out)
			}
			if len(out) != g.MarshalLen() {
				t.Fatal("wrong MarshalLen")
			}
			if err := g.MarshalTo(make([]byte, len(out)-1)); err == nil {
				t.Fatal("short destination accepted")
			}
		})
	}
}

func TestPGWChangeInfoDecodesChildren(t *testing.T) {
	in, _ := hex.DecodeString(genericDatagrams[1])
	g, err := message.ParseGeneric(in)
	if err != nil {
		t.Fatal(err)
	}
	pci := g.IEs[1]
	if !pci.IsGrouped() || len(pci.ChildIEs) != 1 {
		t.Fatalf("grouped=%v children=%d", pci.IsGrouped(), len(pci.ChildIEs))
	}
}

func FuzzGenericRoundTrip(f *testing.F) {
	for _, s := range genericDatagrams {
		b, _ := hex.DecodeString(s)
		f.Add(b)
	}
	for typ := 0; typ < 256; typ++ {
		flags, length := byte(0x48), byte(8)
		if typ >= 1 && typ <= 3 {
			flags, length = 0x40, 4
		}
		b := make([]byte, int(length)+4)
		b[0], b[1], b[3] = flags, byte(typ), length
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		g, err := message.ParseGeneric(b)
		if err != nil {
			return
		}
		out, err := g.Marshal()
		if err != nil {
			t.Fatalf("accepted datagram cannot marshal: %v", err)
		}
		if !bytes.Equal(b, out) {
			t.Fatalf("datagram changed\nin: %x\nout:%x", b, out)
		}
	})
}

func TestGenericAndTypedPiggybackErrors(t *testing.T) {
	for _, s := range []string{
		"582100080000000100000100582200080000000200000200",
		"582100080000000100000100482200080000000200000200",
		"582100080000000100000100285f00080000000200000200",
		"582100080000000100000100405f000400000200",
	} {
		b, _ := hex.DecodeString(s)
		h, err := message.ParseHeader(b)
		if err != nil {
			t.Fatal(err)
		}
		typed, err := message.Parse(b)
		if err != nil {
			t.Fatal(err)
		}
		generic, err := message.ParseGeneric(b)
		if err != nil {
			t.Fatal(err)
		}
		if h.PiggybackError == nil || typed.GetHeader().PiggybackError == nil || generic.PiggybackError == nil {
			t.Fatalf("secondary error missing: %s", s)
		}
		for _, err := range []error{h.PiggybackError, typed.GetHeader().PiggybackError, generic.PiggybackError} {
			switch {
			case errors.Is(h.PiggybackError, message.ErrInvalidPiggyback):
				if !errors.Is(err, message.ErrInvalidPiggyback) {
					t.Fatal(err)
				}
			case errors.Is(h.PiggybackError, message.ErrInvalidTEIDFlag):
				if !errors.Is(err, message.ErrInvalidTEIDFlag) {
					t.Fatal(err)
				}
			default:
				var version *message.UnsupportedVersionError
				if !errors.As(err, &version) || version.Version != 1 {
					t.Fatalf("secondary version error: %v", err)
				}
			}
		}
		if generic.Piggybacked != nil {
			t.Fatalf("invalid secondary became %T", generic.Piggybacked)
		}
		if _, err := message.Marshal(typed); err == nil {
			t.Fatal("typed encoder replayed invalid secondary")
		}
		out, err := generic.Marshal()
		if err != nil || !bytes.Equal(out, b) {
			t.Fatalf("generic replay: %x %v", out, err)
		}
		if err := generic.UnmarshalBinary([]byte{0x40, 1, 0, 4, 0, 0, 2, 0}); err != nil || generic.PiggybackError != nil {
			t.Fatalf("reuse retained error: %v", err)
		}
		out, err = generic.Marshal()
		if err != nil || len(out) != 8 || out[0] != 0x40 {
			t.Fatalf("reuse retained secondary/spares: %x %v", out, err)
		}
	}
}

func TestGenericRetainsNestedIESpares(t *testing.T) {
	b := []byte{0x4f, 33, 0, 18, 0, 0, 0, 1, 0, 0, 1, 0xab, 93, 0, 6, 0xe1, 2, 0, 2, 0xf0, 16, 0}
	g, err := message.ParseGeneric(b)
	if err != nil {
		t.Fatal(err)
	}
	out, err := g.Marshal()
	if err != nil || !bytes.Equal(out, b) {
		t.Fatalf("nested spares changed: %x %v", out, err)
	}
	g.IEs[0].ChildIEs[0].Payload[0] = 17
	out, err = g.Marshal()
	if err != nil || out[len(out)-2] != 17 {
		t.Fatalf("edits ignored: %x %v", out, err)
	}
}

func TestGenericCanDiscardUndecodableSecondary(t *testing.T) {
	b, _ := hex.DecodeString("582100080000000100000100482200080000000200000200")
	g, err := message.ParseGeneric(b)
	if err != nil {
		t.Fatal(err)
	}
	g.PiggybackError = nil
	if _, err := g.Marshal(); !errors.Is(err, message.ErrInvalidPiggyback) {
		t.Fatalf("cleared secondary error without clearing P: %v", err)
	}
	g.Piggybacked = nil
	g.SetPiggybacking(0)
	out, err := g.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 12 || out[0]&0x10 != 0 {
		t.Fatalf("discarded secondary restored: %x", out)
	}
	if _, err := message.ParseGeneric(out); err != nil {
		t.Fatal(err)
	}
}

func TestTypedPiggybackRetainsMessageWithContentError(t *testing.T) {
	// TS 29.274 §7.7.0: the secondary is processed independently.
	wire, err := hex.DecodeString("582100080000000100000100485f000c000000020000020049000000")
	if err != nil {
		t.Fatal(err)
	}
	m, err := message.Parse(wire)
	if err != nil {
		t.Fatal(err)
	}
	h := m.GetHeader()
	var content *message.IEError
	if !errors.As(h.PiggybackError, &content) {
		t.Fatalf("secondary error = %v", h.PiggybackError)
	}
	secondary, ok := h.Piggybacked.(*message.CreateBearerRequest)
	if !ok || secondary == nil || secondary.LinkedEBI == nil {
		t.Fatalf("lost secondary: %T", h.Piggybacked)
	}
	if !bytes.Equal(secondary.LinkedEBI.Payload, []byte{}) {
		t.Fatal("lost erroneous IE contents")
	}
}

func TestNewGenericUsesStrictIEEncoding(t *testing.T) {
	i, err := ie.Parse([]byte{93, 0, 1, 0xf0, 255})
	if err != nil {
		t.Fatal(err)
	}
	for _, g := range []*message.Generic{
		message.NewGeneric(message.MsgTypeCreateSessionResponse, 1, 1, i),
		message.NewGenericWithoutTEID(message.MsgTypeEchoRequest, 0, 1, i),
	} {
		if _, err := g.Marshal(); !errors.Is(err, ie.ErrTooShortToParse) {
			t.Fatalf("constructed Generic Marshal = %v", err)
		}
	}
}

func TestDecodedGenericAppendsLosslessIEs(t *testing.T) {
	g, err := message.ParseGeneric([]byte{0x40, 1, 0, 4, 0, 0, 1, 0})
	if err != nil {
		t.Fatal(err)
	}
	i, err := ie.Parse([]byte{93, 0, 1, 0xf2, 255})
	if err != nil {
		t.Fatal(err)
	}
	g.IEs = append(g.IEs, i)
	out, err := g.Marshal()
	want := []byte{0x40, 1, 0, 9, 0, 0, 1, 0, 93, 0, 1, 0xf2, 255}
	if err != nil || !bytes.Equal(out, want) {
		t.Fatalf("appended IE: %x, %v", out, err)
	}
}
