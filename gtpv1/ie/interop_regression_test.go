package ie_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv1"
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// https://github.com/secdev/scapy/commit/104d5ead0800; TS 29.060 V19.0.0 §7.7: TLV Length excludes three header octets.
func TestInteropV1TLVLength(t *testing.T) {
	b, err := ie.NewMSTimeZone(time.Hour, 0).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{153, 0, 2, 0x40, 0}) {
		t.Fatalf("wire = %x", b)
	}
}

// https://github.com/secdev/scapy/commit/d223f4d32eb6; TS 29.060 V19.0.0 §7.7.32.
func TestInteropIPv6GSN(t *testing.T) {
	b, err := ie.NewGSNAddress("fd01:1::1").Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 19 || b[2] != 16 {
		t.Fatalf("wire = %x", b)
	}
	p, err := ie.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	s, err := p.GSNAddress()
	if err != nil || s != "fd01:1::1" {
		t.Fatalf("address = %q, %v", s, err)
	}
}

// https://github.com/travelping/gtplib/commit/a16cc690e9c1; TS 29.060 V19.0.0 §7.7.2.
func TestInteropShortIMSIFiller(t *testing.T) {
	b, err := ie.NewIMSI("001010123456").Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{2, 0, 1, 1, 0x21, 0x43, 0x65, 0xff, 0xff}
	if !bytes.Equal(b, want) {
		t.Fatalf("short IMSI wire = %x; want %x", b, want)
	}
}

// https://github.com/travelping/gtplib/commit/f3702b2dd10c; TS 29.060 V19.0.0 §7.7.1.
func TestInteropCauseSystemFailure(t *testing.T) {
	b, err := ie.NewCause(gtpv1.ResCauseSystemFailure).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{1, 204}) {
		t.Fatalf("wire = %x", b)
	}
}

// https://github.com/wireshark/wireshark/commit/cf6a3deae44c; TS 29.060 V19.0.0 §7.7.118.
func TestInteropPMTSMIBit(t *testing.T) {
	p := ie.NewExtendedCommonFlagsII(1, 0, 0)
	if len(p.Payload) != 1 || p.Payload[0] != 4 || !p.IsPMTSMI() || p.IsDTCI() {
		t.Fatalf("flags = %x", p.Payload)
	}
}

// https://github.com/wireshark/wireshark/commit/9ff6bb28d2d0; TS 29.060 V19.0.0 §7.7.51 Figure 7.7.51.4.
func TestInteropV1ULIRAIPadding(t *testing.T) {
	b, err := ie.NewUserLocationInformationWithRAI("001", "01", 0x1234, 0x56).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{152, 0, 8, 2, 0, 0xf1, 0x10, 0x12, 0x34, 0x56, 0xff}
	if !bytes.Equal(b, want) {
		t.Fatalf("ULI RAI wire = %x; want %x", b, want)
	}
}
