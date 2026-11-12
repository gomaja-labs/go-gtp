package message_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
)

var genericStructuralCases = []struct{ name, wire string }{
	{"missing_mandatory", "3210000f00000000000100000200010121436587f91405"},
	{"out_of_sequence", "3233000800000001000400000e070180"},
	{"invalid_optional", "320100080000000000010000830001ff"},
	{"empty_gsn", "3202000900000000000100000e07850000"},
	{"unknown_tlv", "320100070000000000010000fa0000"},
	{"header_spares", "3a010004000000000001fefe"},
	{"ie_spares", "32010006000000000001000014f5"},
	{"extension_spares", "3401000800000000000000200101ff00"},
	{"extension_invalid_content", "34010008000000000000008501080100"},
}

func TestGenericPreservesStructuralMessages(t *testing.T) {
	for _, tc := range genericStructuralCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := hex.DecodeString(tc.wire)
			if err != nil {
				t.Fatal(err)
			}
			g, err := message.ParseGeneric(b)
			if err != nil {
				t.Fatal(err)
			}
			out, err := g.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b, out) {
				t.Fatalf("round trip: in %x, out %x", b, out)
			}
			if g.MarshalLen() != len(b) {
				t.Fatalf("MarshalLen = %d, want %d", g.MarshalLen(), len(b))
			}
			dst := bytes.Repeat([]byte{0xcc}, len(b)+1)
			if err := g.MarshalTo(dst); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(dst[:len(b)], b) || dst[len(b)] != 0xcc {
				t.Fatalf("MarshalTo = %x", dst)
			}
			if err := g.MarshalTo(dst[:len(b)-1]); !errors.Is(err, message.ErrTooShortToMarshal) {
				t.Fatalf("short destination: %v", err)
			}
			fresh := message.NewGeneric(message.MsgTypeEchoRequest, 0, 1)
			if err := g.UnmarshalBinary([]byte{0x32, 0x1, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0}); err != nil {
				t.Fatal(err)
			}
			got, err := g.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			want, err := fresh.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) || len(g.IEs) != 0 {
				t.Fatalf("reuse retained state: %x", got)
			}
		})
	}
}

func TestGenericRejectsConstructionErrors(t *testing.T) {
	for _, i := range []*ie.IE{ie.NewIMSI("invalid"), ie.NewGSNAddress("invalid"), ie.New(ie.GSNAddress, []byte{1})} {
		g := message.NewGeneric(message.MsgTypeEchoRequest, 0, 1, i)
		if _, err := g.Marshal(); err == nil {
			t.Fatal("encoded invalid constructed IE")
		}
	}
}

func TestTypedDecodeEnforcesIEOrder(t *testing.T) {
	cases := []struct {
		name    string
		typ     byte
		payload []byte
		missing bool
	}{
		{"echo_request", message.MsgTypeEchoRequest, []byte{14, 7, 1, 128}, false},
		{"echo_response", message.MsgTypeEchoResponse, []byte{14, 7, 1, 128}, false},
		{"create_response", message.MsgTypeCreatePDPContextResponse, []byte{14, 7, 1, 128}, false},
		{"update_response", message.MsgTypeUpdatePDPContextResponse, []byte{14, 7, 1, 128}, false},
		{"delete_response", message.MsgTypeDeletePDPContextResponse, []byte{14, 7, 1, 128}, false},
		{"delete_request", message.MsgTypeDeletePDPContextRequest, []byte{20, 5, 14, 7}, false},
		{"missing_before_invalid_and_sequence", message.MsgTypeCreatePDPContextRequest, []byte{133, 0, 0, 133, 0, 0, 14, 7}, true},
		{"missing_before_sequence", message.MsgTypeCreatePDPContextRequest, []byte{14, 7, 1, 128}, true},
		{"missing_echo_recovery", message.MsgTypeEchoResponse, []byte{255, 0, 2, 0, 1, 1, 128}, true},
		{"end_marker", message.MsgTypeEndMarker, []byte{14, 7, 1, 128}, false},
		{"error_indication", message.MsgTypeErrorIndication, []byte{133, 0, 4, 192, 0, 2, 1, 16, 0, 0, 0, 1}, false},
		{"notification", message.MsgTypeSupportedExtensionHeaderNotification, []byte{255, 0, 2, 0, 1, 141, 1, 0x40}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := message.NewHeader(0x32, tc.typ, 0, 1, tc.payload)
			b, err := h.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			m, err := message.Parse(b)
			want := message.ErrIEOutOfSequence
			if tc.missing {
				want = message.ErrMissingIE
			}
			if m != nil || !errors.Is(err, want) {
				t.Fatalf("Parse = %T, %v; want nil, %v", m, err, want)
			}
			// Each typed parser and direct UnmarshalBinary must enforce the same rule.
			typed := newTypedMessage(tc.typ)
			if err := typed.UnmarshalBinary(b); !errors.Is(err, want) {
				t.Fatalf("UnmarshalBinary = %v, want %v", err, want)
			}
		})
	}
}

func TestTypedEncodeSortsStableRepetitions(t *testing.T) {
	m := message.NewEchoRequest(1, ie.NewRecovery(7), ie.NewCause(128), ie.NewRecovery(8))
	b, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	h, err := message.ParseHeader(b)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{1, 128, 14, 7, 14, 8}
	if !bytes.Equal(h.Payload, want) {
		t.Fatalf("IEs = %x, want %x", h.Payload, want)
	}
	if _, err := message.Parse(b); err != nil {
		t.Fatalf("repeated type rejected: %v", err)
	}
}

func TestGenericRejectsFramingErrors(t *testing.T) {
	for _, payload := range [][]byte{{7, 0}, {14}, {250, 0}, {250, 0, 2, 0}} {
		h := message.NewHeader(0x32, 1, 0, 1, payload)
		b, err := h.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if g, err := message.ParseGeneric(b); g != nil || err == nil {
			t.Fatalf("accepted malformed %x: %v", payload, err)
		}
	}
}

func newTypedMessage(typ byte) message.Message {
	switch typ {
	case message.MsgTypeEchoRequest:
		return &message.EchoRequest{}
	case message.MsgTypeEchoResponse:
		return &message.EchoResponse{}
	case message.MsgTypeCreatePDPContextRequest:
		return &message.CreatePDPContextRequest{}
	case message.MsgTypeCreatePDPContextResponse:
		return &message.CreatePDPContextResponse{}
	case message.MsgTypeUpdatePDPContextResponse:
		return &message.UpdatePDPContextResponse{}
	case message.MsgTypeDeletePDPContextRequest:
		return &message.DeletePDPContextRequest{}
	case message.MsgTypeDeletePDPContextResponse:
		return &message.DeletePDPContextResponse{}
	case message.MsgTypeEndMarker:
		return &message.EndMarker{}
	case message.MsgTypeErrorIndication:
		return &message.ErrorIndication{}
	case message.MsgTypeSupportedExtensionHeaderNotification:
		return &message.SupportedExtensionHeaderNotification{}
	default:
		panic("unhandled test message type")
	}
}

func TestGenericReceivedContentErrorsRemainObservable(t *testing.T) {
	b, err := hex.DecodeString("320100080000000000010000830001ff")
	if err != nil {
		t.Fatal(err)
	}
	g, err := message.ParseGeneric(b)
	if err != nil {
		t.Fatal(err)
	}
	if g.IEs[0].Err() == nil {
		t.Fatal("lost received content error")
	}
	if _, err := g.IEs[0].Marshal(); err == nil {
		t.Fatal("ordinary IE encoding accepted invalid contents")
	}
	// A failed constructor must remain an error when added to a decoded Generic.
	g.AddIE(ie.NewIMSI("invalid"))
	if _, err := g.Marshal(); err == nil {
		t.Fatal("decoded Generic ignored construction error")
	}
}

func TestGenericEncodesEditedAndAppendedContentsLosslessly(t *testing.T) {
	b, err := hex.DecodeString("320100080000000000010000830001ff")
	if err != nil {
		t.Fatal(err)
	}
	g, err := message.ParseGeneric(b)
	if err != nil {
		t.Fatal(err)
	}
	g.IEs[0].Payload[0] = 0xfe
	g.AddIE(ie.New(ie.AccessPointName, []byte{0xfd}))
	out, err := g.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := append(append([]byte(nil), b...), 131, 0, 1, 0xfd)
	want[3] += 4
	want[len(b)-1] = 0xfe
	if !bytes.Equal(out, want) {
		t.Fatalf("lossless edit/append = %x, want %x", out, want)
	}
}

func TestGenericExtensionEditsAreStructural(t *testing.T) {
	for _, size := range []int{2, 10, 1018, 0, 3, 1019} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			b, err := hex.DecodeString("3401000800000000000000200101ff00")
			if err != nil {
				t.Fatal(err)
			}
			g, err := message.ParseGeneric(b)
			if err != nil {
				t.Fatal(err)
			}
			g.ExtensionHeaders[0].Content = bytes.Repeat([]byte{0xff}, size)
			out, err := g.Marshal()
			if size == 0 || size == 3 || size == 1019 {
				if !errors.Is(err, message.ErrInvalidLength) {
					t.Fatalf("structural error = %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if out[12] != byte((size+2)/4) || !bytes.Equal(out[13:len(out)-1], g.ExtensionHeaders[0].Content) || out[len(out)-1] != 0 {
				t.Fatalf("extension bytes changed: %x", out)
			}
			fresh := message.NewGeneric(message.MsgTypeEchoRequest, 0, 1)
			fresh.Flags, fresh.NextExtensionHeaderType = g.Flags, g.NextExtensionHeaderType
			fresh.ExtensionHeaders = g.ExtensionHeaders
			strict, err := fresh.Marshal()
			if size == 2 {
				if err != nil || strict[14] != 0 {
					t.Fatalf("constructed Generic spare = %x, %v", strict, err)
				}
			} else if err == nil {
				t.Fatal("constructed Generic accepted invalid Service Class Indicator")
			}
		})
	}
}

// TS 29.060 §11.1.10 applies the IE order rule to every message, including the
// message types that Parse decodes as Generic; ParseGeneric stays structural.
func TestParseEnforcesIEOrderForGenericTypes(t *testing.T) {
	bad, _ := hex.DecodeString("3233000800000001000400000e070180")
	if m, err := message.Parse(bad); m != nil || !errors.Is(err, message.ErrIEOutOfSequence) {
		t.Fatalf("Parse = %T, %v; want nil, %v", m, err, message.ErrIEOutOfSequence)
	}
	g, err := message.ParseGeneric(bad)
	if err != nil {
		t.Fatalf("ParseGeneric: %v", err)
	}
	if out, err := g.Marshal(); err != nil || !bytes.Equal(out, bad) {
		t.Fatalf("replay = % x, %v; want % x", out, err, bad)
	}
	good, _ := hex.DecodeString("32330008000000010004000001800e07")
	m, err := message.Parse(good)
	if err != nil {
		t.Fatalf("Parse(in order): %v", err)
	}
	if _, ok := m.(*message.Generic); !ok {
		t.Fatalf("Parse(in order) = %T, want *message.Generic", m)
	}
}
