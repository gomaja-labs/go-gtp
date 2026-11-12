package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
	"testing"
)

func TestTPDURetainsPDUSessionFutureExtension(t *testing.T) {
	// TS 38.415 V19.1.0 §5.5.1 Figure 5.5.1-1 and §5.5.2.2
	// Figure 5.5.2.2-1: remove unknown Future Extension fields.
	// UL/QFI 9, New IE Flags bit 5, four future value octets, three pad octets.
	// The extension length (3 four-octet units) and GTP length are consistent.
	b := []byte{0x34, 255, 0, 17, 0, 0, 0, 1, 0, 0, 0, 0x85,
		3, 0x10, 0x49, 0x20, 0xde, 0xad, 0xbe, 0xef, 0, 0, 0, 0, 0x45}
	m, err := ParseTPDU(b)
	if err != nil {
		t.Fatalf("future PDU Session Container extension rejects G-PDU: %v", err)
	}
	if len(m.Payload) != 1 || m.Payload[0] != 0x45 {
		t.Fatalf("inner payload lost: %x", m.Payload)
	}
	if !bytes.Equal(m.ExtensionHeaders[0].Content, b[13:23]) {
		t.Fatal("lost future extension content")
	}
	if _, err := m.Marshal(); err == nil {
		t.Fatal("send-side validation accepted unknown future layout")
	}
	// Selecting known fields still requires their complete values.
	b[15] = 0x1f
	if _, err := ParseTPDU(b); err == nil {
		t.Fatal("future data hid truncated known fields")
	}
}

func TestEchoRequestRetainsInvalidUnexpectedIE(t *testing.T) {
	// TS 29.060 V19.0.0 §11.1.11: ignore an unexpected, framed TLV.
	// Echo Request followed by an unexpected APN containing an invalid label.
	b := []byte{0x32, 1, 0, 8, 0, 0, 0, 0, 0, 1, 0, 0, 131, 0, 1, 0xff}
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

func TestCreatePDPContextMixedAndNonIPAddresses(t *testing.T) {
	// TS 29.060 V19.0.0 §7.7.27 Figures 38/39: dual stack permits a
	// single static family, and Non-IP uses ETSI number 2 without an address.
	for _, p := range [][]byte{
		{0xf1, 0x8d, 192, 0, 2, 9},
		{0xf1, 0x8d, 0x20, 1, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9},
		{0xf0, 2},
	} {
		t.Run(string([]byte{'n', byte('0' + len(p)/10), byte('0' + len(p)%10)}), func(t *testing.T) {
			b := []byte{0x32, 16, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0,
				2, 0, 1, 1, 0, 0, 0, 0, 0xf1, 15, 0,
				16, 0, 0, 0, 1, 17, 0, 0, 0, 2, 20, 5, 128, 0, byte(len(p))}
			b = append(b, p...)
			b = append(b, 131, 0, 9, 8, 'i', 'n', 't', 'e', 'r', 'n', 'e', 't',
				133, 0, 4, 192, 0, 2, 1, 133, 0, 4, 192, 0, 2, 1,
				134, 0, 3, 0x91, 0x21, 0xf3, 135, 0, 4, 0, 0, 0, 0)
			binary.BigEndian.PutUint16(b[2:4], uint16(len(b)-8))
			m, err := ParseCreatePDPContextRequest(b)
			if err != nil {
				t.Fatalf("valid End User Address %x rejects Create PDP Context: %v", p, err)
			}
			if got, err := m.EndUserAddress.EndUserAddress(); err != nil || !bytes.Equal(got, p) {
				t.Fatalf("End User Address accessor = %x, %v", got, err)
			}
			if _, err := m.Marshal(); err != nil {
				t.Fatalf("could not send valid End User Address: %v", err)
			}
		})
	}
}

func TestDeletePDPContextRetainsInvalidOptionalPCO(t *testing.T) {
	// TS 29.060 Table 11, §§11.1.8, 11.1.13: discard erroneous optional PCO.
	b := []byte{0x32, 20, 0, 10, 0, 0, 0, 1, 0, 1, 0, 0, 20, 5, 132, 0, 1, 0x80}
	// A length-consistent PCO containing a truncated protocol/container header.
	b = append(b, 0xff)
	b[3], b[16] = 11, 2
	m, err := ParseDeletePDPContextRequest(b)
	if err != nil {
		t.Fatalf("optional PCO prevents Delete PDP Context processing: %v", err)
	}

	if m.NSAPI.MustNSAPI() != 5 || m.PCO == nil || !bytes.Equal(m.PCO.Payload, []byte{0x80, 0xff}) {
		t.Fatal("lost mandatory or erroneous optional IE")
	}
	if _, err := m.PCO.ProtocolConfigurationOptions(); err == nil {
		t.Fatal("invalid PCO accessor succeeded")
	}
	if _, err := m.Marshal(); err == nil {
		t.Fatal("encoded invalid PCO")
	}
	good, err := NewDeletePDPContextRequest(1, 2, ie.NewNSAPI(5)).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := m.UnmarshalBinary(good); err != nil || m.PCO != nil {
		t.Fatalf("reuse retained invalid PCO: %v", err)
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
	b, err := NewHeader(0x32, MsgTypeCreatePDPContextRequest, 1, 123, payload).Marshal()
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
