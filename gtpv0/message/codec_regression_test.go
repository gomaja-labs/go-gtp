package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
	"testing"
)

func TestRepeatedMarshal(t *testing.T) {
	for _, m := range []Message{NewCreatePDPContextRequest(1, 1, 1, mandatoryCreateIEs()...), NewCreatePDPContextResponse(1, 1, 1, ie.NewCause(193)), NewDeletePDPContextResponse(1, 1, 1, ie.NewCause(193))} {
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
	b, _ := NewHeader(0x1e, 2, 1, 0, 0, nil).Marshal()
	if _, err := Parse(b); err == nil {
		t.Fatal("accepted Echo Response without Recovery")
	}
}
func TestHeaderLengthArithmetic(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Errorf("unexpected panic: %v", p)
		}
	}()
	b := make([]byte, 65555)
	b[0] = 0x1e
	b[1] = 255
	binary.BigEndian.PutUint16(b[2:4], 65535)
	if _, err := ParseHeader(b); err != nil {
		t.Fatal(err)
	}
}
func TestSNDCP(t *testing.T) {
	h := NewHeader(0x1f, 255, 1, 1, 1, nil)
	h.SndcpNumber = 0x42
	b, _ := h.Marshal()
	got, err := ParseHeader(b)
	if err != nil || got.SndcpNumber != 0x42 {
		t.Fatalf("SNDCP=%02x err=%v; want 42", got.SndcpNumber, err)
	}
}
func TestHeaderLengthValidation(t *testing.T) {
	b := make([]byte, 20)
	b[0] = 0x1e
	b[1] = 255
	b[3] = 10
	if _, err := ParseHeader(b); err == nil {
		t.Fatal("accepted declared payload length 10 with zero received bytes")
	}
}
func TestLengthOverflow(t *testing.T) {
	if _, err := NewTPDU(1, 1, 1, make([]byte, 65536)).Marshal(); err == nil {
		t.Fatal("encoded 65536-byte payload with wrapped length")
	}
}
func TestUnmarshalReuse(t *testing.T) {
	a, _ := NewEchoRequest(1, 0, 0, ie.NewPrivateExtension(1, []byte{1})).Marshal()
	b, _ := NewEchoRequest(2, 0, 0).Marshal()
	m := &EchoRequest{}
	if err := m.UnmarshalBinary(a); err != nil {
		t.Fatal(err)
	}
	if err := m.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if m.PrivateExtensions != nil {
		t.Fatal("second packet retained first packet's PrivateExtension")
	}
}
func TestRepeatedIEFirstWins(t *testing.T) {
	b, _ := NewHeader(0x1e, 17, 1, 0, 1, []byte{1, 193, 1, 128}).Marshal()
	m, err := ParseCreatePDPContextResponse(b)
	if err != nil {
		t.Fatal(err)
	}
	if m.Cause.MustCause() != 193 {
		t.Fatalf("duplicate Cause overwrote first: got %d", m.Cause.MustCause())
	}
}
func TestIEOrdering(t *testing.T) {
	b, err := NewGeneric(99, 1, 0, 0, ie.NewRecovery(1), ie.NewCause(193)).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if b[20] != 1 {
		t.Fatalf("encoded IE type %d before type 1", b[20])
	}
}
func TestDecapsulateNonTPDU(t *testing.T) {
	b, _ := NewEchoRequest(1, 0, 0).Marshal()
	if _, err := Decapsulate(b); err == nil {
		t.Fatal("Decapsulate reported success for Echo Request")
	}
}

func mandatoryCreateIEs() []*ie.IE {
	return []*ie.IE{ie.NewQualityOfServiceProfile(1, 1, 1, 1, 1), ie.NewSelectionMode(0), ie.NewFlowLabelDataI(1), ie.NewFlowLabelSignalling(1), ie.NewEndUserAddress("192.0.2.1"), ie.NewAccessPointName("internet"), ie.NewGSNAddress("192.0.2.1"), ie.NewGSNAddress("192.0.2.2"), ie.NewMSISDN("12345")}
}

func TestEachMandatoryCreateIE(t *testing.T) {
	all := mandatoryCreateIEs()
	for omitted := range all {
		ies := append([]*ie.IE{}, all[:omitted]...)
		ies = append(ies, all[omitted+1:]...)
		if _, err := NewCreatePDPContextRequest(1, 1, 1, ies...).Marshal(); err == nil {
			t.Fatalf("encoded create without mandatory index %d", omitted)
		}
		var payload []byte
		for _, i := range ies {
			b, err := i.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			payload = append(payload, b...)
		}
		b, err := NewHeader(0x1e, 16, 1, 1, 1, payload).Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Parse(b); !errors.Is(err, ErrMissingIE) {
			t.Fatalf("omitted %d: %v", omitted, err)
		}
	}
}

func TestPrivateExtensions(t *testing.T) {
	m := NewEchoRequest(1, 0, 0, ie.NewPrivateExtension(1, []byte{1}), ie.NewPrivateExtension(2, []byte{2}))
	b, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseEchoRequest(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.PrivateExtensions) != 2 {
		t.Fatal("lost private extensions")
	}
	again, err := parsed.Marshal()
	if err != nil || !bytes.Equal(again, b) {
		t.Fatalf("round trip: %x %v", again, err)
	}
}
