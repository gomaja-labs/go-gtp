package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
	"testing"
)

func testNoPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if p := recover(); p != nil {
			t.Errorf("unexpected panic: %v", p)
		}
	}()
	f()
}
func testMessages() []Message {
	return []Message{
		NewCreatePDPContextRequest(1, 1, mandatoryCreateIEs()...), NewCreatePDPContextResponse(1, 1, ie.NewCause(193)),
		NewUpdatePDPContextRequest(1, 1, ie.NewNSAPI(5)), NewUpdatePDPContextResponse(1, 1, ie.NewCause(193)),
		NewDeletePDPContextRequest(1, 1, ie.NewNSAPI(5)), NewDeletePDPContextResponse(1, 1, ie.NewCause(193)),
		NewErrorIndication(0, 1, ie.NewTEIDDataI(1), ie.NewGSNAddress("192.0.2.1")),
		NewSupportedExtensionHeaderNotification(0, 1, ie.NewExtensionHeaderTypeList(0x85)),
		NewVersionNotSupported(0, 1, ie.NewPrivateExtension(1, []byte{1})),
	}
}
func TestRepeatedMarshal(t *testing.T) {
	for _, m := range testMessages() {
		t.Run(m.MessageTypeName(), func(t *testing.T) {
			b, err := Marshal(m)
			if err != nil {
				t.Fatal(err)
			}
			got, err := Marshal(m)
			if err != nil || !bytes.Equal(got, b) {
				t.Fatalf("second Marshal: %v; want identical bytes", err)
			}
		})
	}
}
func TestRequiredIEs(t *testing.T) {
	for _, typ := range []uint8{2, 16, 17, 18, 19, 20, 21, 26, 31} {
		t.Run(fmt.Sprint(typ), func(t *testing.T) {
			b, _ := NewHeader(0x32, typ, 0, 1, nil).Marshal()
			if _, err := Parse(b); err == nil {
				t.Fatal("accepted message missing mandatory IEs")
			}
		})
	}
	t.Run("one_byte_IE", func(t *testing.T) {
		b, _ := NewHeader(0x32, 1, 0, 1, []byte{255}).Marshal()
		if _, err := Parse(b); err == nil {
			t.Fatal("accepted truncated single-byte IE")
		}
	})
}
func TestHeaderLengthArithmetic(t *testing.T) {
	testNoPanic(t, func() {
		b := make([]byte, 65543)
		b[0] = 0x30
		b[1] = 255
		binary.BigEndian.PutUint16(b[2:4], 65535)
		h, err := ParseHeader(b)
		if err != nil || len(h.Payload) != 65535 {
			t.Errorf("legal maximal header: payload/err=%v", err)
		}
	})
}
func TestHeaderDiscriminators(t *testing.T) {
	for _, flags := range []byte{0x10, 0x50, 0x20} {
		t.Run(fmt.Sprintf("%x", flags), func(t *testing.T) {
			b := []byte{flags, 255, 0, 0, 0, 0, 0, 1}
			if _, err := Parse(b); err == nil {
				t.Fatalf("accepted non-GTPv1/PT flags %02x", flags)
			}
		})
	}
}
func TestHeaderLengthMismatch(t *testing.T) {
	b, _ := NewHeader(0x32, 1, 0, 1, nil).Marshal()
	b = append(b, 14, 0)
	if _, err := Parse(b); err == nil {
		t.Fatal("accepted UDP bytes beyond the advertised GTP length")
	}
}
func TestExtensionAppend(t *testing.T) {
	h := NewHeaderWithExtensionHeaders(0x30, 255, 1, 0, nil, NewExtensionHeader(0x40, []byte{1, 2}, 0))
	if err := h.AddExtensionHeaders(NewExtensionHeader(0x85, []byte{0, 4}, 0)); err != nil {
		t.Fatal(err)
	}
	b, err := h.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	p, err := ParseHeader(b)
	if err != nil || len(p.ExtensionHeaders) != 2 || p.ExtensionHeaders[0].Type != 0x40 {
		t.Fatalf("appending header broke chain: parsed=%+v err=%v", p, err)
	}
}
func TestInvalidExtensionChain(t *testing.T) {
	testNoPanic(t, func() {
		h := NewHeaderWithExtensionHeaders(0x30, 255, 1, 0, nil, NewExtensionHeader(0x40, []byte{1, 2}, 0x85))
		if _, err := h.Marshal(); err == nil {
			t.Error("constructor discarded unterminated extension-chain error")
		}
	})
	testNoPanic(t, func() { _ = NewHeaderWithExtensionHeaders(0x30, 255, 1, 0, nil, nil) })
}
func TestLengthOverflow(t *testing.T) {
	if _, err := NewTPDU(1, make([]byte, 65536)).Marshal(); err == nil {
		t.Error("encoded 65536-byte payload with wrapped length")
	}
	testNoPanic(t, func() {
		if _, err := NewExtensionHeader(0x85, make([]byte, 1019), 0).Marshal(); err == nil {
			t.Error("encoded oversized extension with wrapped length")
		}
	})
}
func TestSetTEIDSpare(t *testing.T) {
	h := NewHeader(0x30, 255, 1, 0, nil)
	h.SetTEID(2)
	b, err := h.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if b[0] != 0x30 {
		t.Fatalf("SetTEID changed flags to %02x; want 30", b[0])
	}
}
func TestUnmarshalReuse(t *testing.T) {
	a, _ := NewCreatePDPContextRequest(1, 1, ie.NewNSAPI(5), ie.NewTEIDDataI(1), ie.NewQoSProfile([]byte{1, 2, 3, 4}), ie.NewGSNAddress("192.0.2.1"), ie.NewGSNAddress("192.0.2.3")).Marshal()
	b, _ := NewCreatePDPContextRequest(1, 2, ie.NewNSAPI(6), ie.NewTEIDDataI(1), ie.NewQoSProfile([]byte{1, 2, 3, 4}), ie.NewGSNAddress("192.0.2.2"), ie.NewGSNAddress("192.0.2.4")).Marshal()
	m := &CreatePDPContextRequest{}
	if err := m.UnmarshalBinary(a); err != nil {
		t.Fatal(err)
	}
	if err := m.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if m.NSAPI.MustNSAPI() != 6 || m.SGSNAddressForSignalling.MustGSNAddress() != "192.0.2.2" {
		t.Fatalf("stale NSAPI=%d address=%s", m.NSAPI.MustNSAPI(), m.SGSNAddressForSignalling.MustGSNAddress())
	}
}
func TestHeaderReuse(t *testing.T) {
	a, _ := NewHeaderWithExtensionHeaders(0x32, 255, 1, 1, nil, NewExtensionHeader(0x40, []byte{1, 2}, 0)).Marshal()
	b, _ := NewHeader(0x30, 255, 2, 0, nil).Marshal()
	h := &Header{}
	if err := h.UnmarshalBinary(a); err != nil {
		t.Fatal(err)
	}
	if err := h.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if h.SequenceNumber != 0 || len(h.ExtensionHeaders) != 0 {
		t.Fatalf("stale optional state: seq=%d extensions=%d", h.SequenceNumber, len(h.ExtensionHeaders))
	}
}
func TestRepeatedIEFirstWins(t *testing.T) {
	b, _ := NewHeader(0x32, 17, 1, 1, []byte{1, 193, 1, 128}).Marshal()
	m, err := ParseCreatePDPContextResponse(b)
	if err != nil {
		t.Fatal(err)
	}
	if m.Cause.MustCause() != 193 {
		t.Fatalf("duplicate Cause overwrote first: got %d", m.Cause.MustCause())
	}
}
func TestIEOrdering(t *testing.T) {
	b, err := NewGeneric(99, 0, 1, ie.NewRecovery(1), ie.NewCause(193)).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if b[12] != 1 {
		t.Fatalf("encoded IE type %d before type 1", b[12])
	}
}
func TestNamedMessageIEOrdering(t *testing.T) {
	cases := []Message{
		NewCreatePDPContextResponse(1, 1, ie.NewCause(128), ie.New(ie.ChargingGatewayAddress, []byte{192, 0, 2, 1}), ie.NewCommonFlags(0, 0, 0, 0, 0, 0, 0, 0)),
		NewUpdatePDPContextResponse(1, 1, ie.NewCause(128), ie.New(ie.BearerControlMode, []byte{1}), ie.New(ie.MSInfoChangeReportingAction, []byte{1})),
		NewUpdatePDPContextRequest(1, 1, ie.NewNSAPI(5), ie.NewIMEISV("1234567890123456"), ie.New(ie.CNOperatorSelectionEntity, []byte{0})),
	}
	for _, m := range cases {
		t.Run(m.MessageTypeName(), func(t *testing.T) {
			b, err := Marshal(m)
			if err != nil {
				t.Fatal(err)
			}
			h, err := ParseHeader(b)
			if err != nil {
				t.Fatal(err)
			}
			xs, err := ie.ParseMultiIEs(h.Payload)
			if err != nil {
				t.Fatal(err)
			}
			for j := 1; j < len(xs); j++ {
				if xs[j].Type < xs[j-1].Type {
					t.Fatalf("IE type %d encoded after %d", xs[j].Type, xs[j-1].Type)
				}
			}
		})
	}
}
func TestZeroNextExtension(t *testing.T) {
	b := []byte{0x34, 255, 0, 4, 0, 0, 0, 1, 0, 0, 0, 0}
	h, err := ParseHeader(b)
	if err != nil {
		t.Fatalf("E=1, next=0, no extensions rejected: %v", err)
	}
	if len(h.ExtensionHeaders) != 0 {
		t.Fatal("invented extensions")
	}
}
func TestNotificationSequenceFlag(t *testing.T) {
	b, err := NewSupportedExtensionHeaderNotification(0, 1, ie.NewExtensionHeaderTypeList(0x85)).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if b[0]&2 == 0 {
		t.Fatalf("Supported Extension Headers Notification S flag is zero: %02x", b[0])
	}
}
func TestMessageTypeConstants(t *testing.T) {
	if MsgTypeIdentificationRequest != 48 || MsgTypeSGSNContextAcknowledge != 52 {
		t.Fatalf("Identification=%d SGSN Ack=%d; want 48,52", MsgTypeIdentificationRequest, MsgTypeSGSNContextAcknowledge)
	}
}
func TestUnusedOctets(t *testing.T) {
	h := NewHeader(0x32, 1, 0, 1, nil)
	b := bytes.Repeat([]byte{0xff}, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		t.Fatal(err)
	}
	if b[10] != 0 || b[11] != 0 {
		t.Errorf("unused optional octets=%x; want 0000", b[10:12])
	}
	e := NewExtensionHeader(0x20, []byte{1}, 0)
	b = bytes.Repeat([]byte{0xff}, e.MarshalLen())
	if err := e.MarshalTo(b); err != nil {
		t.Fatal(err)
	}
	if b[2] != 0 {
		t.Errorf("extension spare octet=%02x; want 00", b[2])
	}
}
func TestExtensionSemanticLength(t *testing.T) {
	for _, typ := range []uint8{0x40, 0xc0, 0x03, 0x20, 0x85} {
		t.Run(fmt.Sprintf("%x", typ), func(t *testing.T) {
			// QMP=1 in a two-octet DL PDU Session frame advertises an absent timestamp.
			content := []byte{0x08, 1}
			if typ != 0x85 {
				content = make([]byte, 10)
			}
			h := NewHeaderWithExtensionHeaders(0x30, 255, 1, 0, nil, NewExtensionHeader(typ, content, 0))
			b, err := h.Marshal()
			if err == nil {
				_, err = Parse(b)
			}
			if err == nil {
				t.Fatalf("accepted invalid extension layout for type %02x", typ)
			}
		})
	}
}

func mandatoryCreateIEs() []*ie.IE {
	return []*ie.IE{ie.NewTEIDDataI(1), ie.NewNSAPI(5), ie.NewGSNAddress("192.0.2.1"), ie.NewGSNAddress("192.0.2.2"), ie.NewQoSProfile([]byte{1, 2, 3, 4})}
}

func TestEachMandatoryCreateIE(t *testing.T) {
	all := mandatoryCreateIEs()
	for omitted := range all {
		ies := append([]*ie.IE{}, all[:omitted]...)
		ies = append(ies, all[omitted+1:]...)
		if _, err := NewCreatePDPContextRequest(1, 1, ies...).Marshal(); err == nil {
			t.Fatalf("encoded create request without mandatory index %d", omitted)
		}
		var payload []byte
		for _, i := range ies {
			b, err := i.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			payload = append(payload, b...)
		}
		b, err := NewHeader(0x32, 16, 1, 1, payload).Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Parse(b); !errors.Is(err, ErrMissingIE) {
			t.Fatalf("omitted %d: %v", omitted, err)
		}
	}
}

func TestC038_PrivateExtensionsAndUnknownIEs(t *testing.T) {
	m := NewEchoRequest(1, ie.NewPrivateExtension(1, []byte{1}), ie.NewPrivateExtension(2, []byte{2}), ie.New(240, []byte{3}))
	b, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseEchoRequest(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.PrivateExtensions) != 2 || len(parsed.AdditionalIEs) != 1 {
		t.Fatalf("lost repeated or unknown IEs: %+v", parsed)
	}
	again, err := parsed.Marshal()
	if err != nil || !bytes.Equal(again, b) {
		t.Fatalf("round trip: %x %v", again, err)
	}
	empty, err := NewEchoRequest(2).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := parsed.UnmarshalBinary(empty); err != nil {
		t.Fatal(err)
	}
	if len(parsed.PrivateExtensions) != 0 || len(parsed.AdditionalIEs) != 0 {
		t.Fatal("retained prior IEs")
	}
}

func TestParsedCreateRepeatedMarshal(t *testing.T) {
	b, err := NewCreatePDPContextRequest(1, 1, mandatoryCreateIEs()...).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	m, err := Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	for range 3 {
		out, err := Marshal(m)
		if err != nil || len(out) != m.MarshalLen() || !bytes.Equal(out, b) {
			t.Fatalf("parsed create marshal: %x %v", out, err)
		}
	}
}
