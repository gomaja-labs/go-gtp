package message

import (
	"bytes"
	"errors"
	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
	"testing"
)

func TestEchoRequestRetainsInvalidUnexpectedIE(t *testing.T) {
	// GSM 09.60 V7.10.0 §10.1.11: unexpected TLV is skipped.
	b := []byte{0x1e, 1, 0, 4, 0, 1, 0, 0, 0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 131, 0, 1, 0xff}
	m, err := ParseEchoRequest(b)
	if err != nil {
		t.Fatalf("unexpected APN prevents Echo Request processing: %v", err)
	}
	if len(m.AdditionalIEs) != 1 || !bytes.Equal(m.AdditionalIEs[0].Payload, []byte{0xff}) {
		t.Fatal("lost erroneous unexpected IE")
	}
	if _, err := m.AdditionalIEs[0].AccessPointName(); err == nil {
		t.Fatal("invalid unexpected APN accessor succeeded")
	}
	if _, err := m.Marshal(); err == nil {
		t.Fatal("encoded invalid unexpected APN")
	}
	// Delimiting the IE is still mandatory, even when its type is unexpected.
	b = b[:len(b)-1]
	b[3]--
	if m, err := ParseEchoRequest(b); m != nil || err == nil {
		t.Fatal("accepted truncated unexpected IE")
	}
}

func TestInvalidMandatoryIEReturnsMessageAndTypedError(t *testing.T) {
	var payload []byte
	gsn := 0
	for _, i := range mandatoryCreateIEs() {
		if i.Type == ie.GSNAddress {
			if gsn == 1 {
				payload = append(payload, ie.GSNAddress, 0, 1, 0xff)
				gsn++
				continue
			}
			gsn++
		}
		b, err := i.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		payload = append(payload, b...)
	}
	b, err := NewHeader(0x1e, MsgTypeCreatePDPContextRequest, 123, 1, 1, payload).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	m, err := ParseCreatePDPContextRequest(b)
	var invalid *InvalidIEError
	if m == nil || !errors.As(err, &invalid) || invalid.Type != ie.GSNAddress || invalid.Occurrence != 1 || !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("message=%v error=%v", m, err)
	}
	generic, err := Parse(b)
	if generic == nil || !errors.As(err, &invalid) {
		t.Fatalf("Parse lost mandatory error context: %v %v", generic, err)
	}
	if _, err := m.Marshal(); err == nil {
		t.Fatal("encoded incorrect mandatory IE")
	}
}
