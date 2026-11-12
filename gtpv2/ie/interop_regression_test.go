// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// https://github.com/secdev/scapy/commit/18d4c30fa4ec; TS 29.274 V20.0.0 §8.2.
func TestV2IELength(t *testing.T) {
	b, err := ie.NewRecovery(13).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, []byte{3, 0, 1, 0, 13}) {
		t.Fatalf("wire = %x", b)
	}
}

// https://github.com/secdev/scapy/commit/07854abd9319; TS 29.274 V20.0.0 §8.66.
func TestFQDNLabels(t *testing.T) {
	p := ie.NewFullyQualifiedDomainName("pgw1.test.local")
	want := []byte{4, 'p', 'g', 'w', '1', 4, 't', 'e', 's', 't', 5, 'l', 'o', 'c', 'a', 'l'}
	if !bytes.Equal(p.Payload, want) {
		t.Fatalf("FQDN = %x", p.Payload)
	}
	s, err := p.FullyQualifiedDomainName()
	if err != nil || s != "pgw1.test.local" {
		t.Fatalf("FQDN = %q, %v", s, err)
	}
}

// https://github.com/travelping/gtplib/commit/d990d22bbf8f; TS 29.274 V20.0.0 §8.10.
func TestMEITBCD(t *testing.T) {
	p := ie.NewMobileEquipmentIdentity("123456789012345")
	want := []byte{0x21, 0x43, 0x65, 0x87, 9, 0x21, 0x43, 0xf5}
	if !bytes.Equal(p.Payload, want) {
		t.Fatalf("MEI = %x", p.Payload)
	}
	s, err := p.MobileEquipmentIdentity()
	if err != nil || s != "123456789012345" {
		t.Fatalf("MEI = %q, %v", s, err)
	}
}

// https://github.com/travelping/gtplib/commit/f63ff4f362d6; TS 29.274 V20.0.0 §8.12 Figure 8.12-1 (mandatory octets 5 and 6).
func TestIndicationMinimum(t *testing.T) {
	for _, p := range []*ie.IE{ie.NewIndicationFromOctets(), ie.NewIndicationFromOctets(0x80)} {
		b, err := p.Marshal()
		if err == nil && len(b) < 6 {
			t.Errorf("Indication emitted %d value octets; want at least 2", len(b)-4)
		}
	}
}

// https://github.com/travelping/gtplib/commit/602c4adaa318; TS 29.274 V20.0.0 §8.18.
func TestPLMNNibbles(t *testing.T) {
	for _, tc := range []struct {
		mnc  string
		want []byte
	}{{"260", []byte{0x13, 0, 0x62}}, {"26", []byte{0x13, 0xf0, 0x62}}} {
		p := ie.NewServingNetwork("310", tc.mnc)
		if !bytes.Equal(p.Payload, tc.want) {
			t.Fatalf("MNC %s: wire = %x", tc.mnc, p.Payload)
		}
		s, err := p.ServingNetwork()
		if err != nil || s != "310"+tc.mnc {
			t.Fatalf("PLMN = %q, %v", s, err)
		}
	}
}

// https://github.com/travelping/gtplib/commit/94f853c77071; TS 29.274 V20.0.0 §8.15 Figure 8.15-1.
func TestQoS40Bits(t *testing.T) {
	p := ie.NewBearerQoS(0, 0, 0, 9, 0x0102030405, 0x060708090a, 0x0b0c0d0e0f, 0x1011121314)
	want := []byte{0, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	if !bytes.Equal(p.Payload, want) {
		t.Fatalf("QoS = %x", p.Payload)
	}
	f, err := p.BearerQoS()
	if err != nil {
		t.Fatal(err)
	}
	if f.MaximumBitRateForUplink != 0x0102030405 || f.GuaranteedBitRateForDownlink != 0x1011121314 {
		t.Fatalf("QoS = %+v", f)
	}
}

// https://github.com/travelping/gtplib/commit/8556074f0c37; TS 24.008 V20.1.0 §10.5.6.3 (FF00 operator-specific container).
func TestPCOVendorFF00(t *testing.T) {
	p := ie.NewProtocolConfigurationOptions(0, ie.NewPCOContainer(0xff00, []byte{0x13, 1, 0x84}))
	want := []byte{0x80, 0xff, 0, 3, 0x13, 1, 0x84}
	if !bytes.Equal(p.Payload, want) {
		t.Fatalf("PCO = %x", p.Payload)
	}
	f, err := p.ProtocolConfigurationOptions()
	if err != nil {
		t.Fatal(err)
	}
	if len(f.ProtocolOrContainers) != 1 || f.ProtocolOrContainers[0].ID != 0xff00 {
		t.Fatalf("PCO = %+v", f)
	}
}

// https://github.com/travelping/gtplib/commit/f036d98d2e01; TS 24.008 V20.1.0 §10.5.6.3.
func TestPCOMixedOrder(t *testing.T) {
	raw := []byte{0x80, 0, 13, 0, 0x80, 0x21, 10, 1, 7, 0, 10, 129, 6, 192, 0, 2, 1}
	f, err := ie.ParseProtocolConfigurationOptionsFields(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.ProtocolOrContainers) != 2 || f.ProtocolOrContainers[1].ID != 0x8021 {
		t.Fatalf("PCO = %+v", f)
	}
	ppp, err := ie.ParsePCOPPP(f.ProtocolOrContainers[1].Contents)
	if err != nil || ppp.Code != 1 || ppp.Identifier != 7 {
		t.Fatalf("PPP following container = %+v, %v", ppp, err)
	}
	b, err := f.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, raw) {
		t.Fatalf("PCO = %x", b)
	}
}

// https://github.com/wireshark/wireshark/commit/37e871a75653; TS 29.274 V20.0.0 §8.21.3.
func TestRAIPadding(t *testing.T) {
	raw := []byte{4, 0, 0xf1, 0x10, 0x12, 0x34, 0x56, 0xff}
	f, err := ie.ParseUserLocationInformationFields(raw)
	if err != nil {
		t.Fatal(err)
	}
	if f.RAI.RAC != 0x56 {
		t.Errorf("decoded RAC = %#x; want 0x56 (padding excluded)", f.RAI.RAC)
	}
	f.RAI.RAC = 0x56
	b, err := f.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, raw) {
		t.Errorf("RAI wire = %x; want %x", b, raw)
	}
}

// https://github.com/open5gs/open5gs/commit/b21dc65741da; TS 29.274 V20.0.0 §8.1.
func TestUnknownTLV(t *testing.T) {
	raw := []byte{250, 0, 1, 0, 0xaa, 3, 0, 1, 0, 7}
	ps, err := ie.ParseMultiIEs(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 2 || ps[1].Type != ie.Recovery || ps[1].Payload[0] != 7 {
		t.Fatalf("IEs = %+v", ps)
	}
	// Unknown children must not suppress known siblings in a grouped IE.
	grouped := append([]byte{93, 0, 10, 0}, raw...)
	p, err := ie.Parse(grouped)
	if err != nil || len(p.ChildIEs) != 2 {
		t.Fatalf("group = %+v, %v", p, err)
	}
}

// https://github.com/open5gs/open5gs/commit/df92296dbf60; TS 29.274 V20.0.0 §8.3.
func TestEmptyIMSI(t *testing.T) {
	p, err := ie.Parse([]byte{1, 0, 0, 0})
	if err != nil {
		return
	}
	if s, err := p.IMSI(); err == nil {
		t.Fatalf("empty IMSI accepted as %q", s)
	}
}

// https://github.com/wireshark/wireshark/commit/6180d0fdaa5e; TS 29.274 V20.0.0 §8.6.
func TestLongAPNLabel(t *testing.T) {
	s := strings.Repeat("a", 32) + ".example"
	p := ie.NewAccessPointName(s)
	if p.Payload[0] != 32 {
		t.Fatalf("APN = %x", p.Payload)
	}
	got, err := p.AccessPointName()
	if err != nil || got != s {
		t.Fatalf("APN = %q, %v", got, err)
	}
}

// https://github.com/wireshark/wireshark/commit/83cb6ff625a4; TS 29.274 V20.0.0 §8.67.
func TestPrivateType255(t *testing.T) {
	p, err := ie.Parse([]byte{255, 0, 3, 0, 0x12, 0x34, 0xaa})
	if err != nil {
		t.Fatal(err)
	}
	id, err := p.EnterpriseID()
	if err != nil || id != 0x1234 {
		t.Fatalf("enterprise = %x, %v", id, err)
	}
	v, err := p.PrivateExtension()
	if err != nil || !bytes.Equal(v, []byte{0xaa}) {
		t.Fatalf("private = %x, %v", v, err)
	}
}

// https://github.com/wireshark/wireshark/commit/f59262b94c10; TS 29.274 V20.0.0 §§8.25–8.26.
func TestForwardingAddresses(t *testing.T) {
	for _, s := range []string{"192.0.2.1", "2001:db8::1"} {
		f, err := ie.NewS103PDNDataForwardingInfo(s, 0x12345678, 5).S103PDNDataForwardingInfo()
		if err != nil || f.HSGWAddressForForwarding.String() != s || f.GREKey != 0x12345678 {
			t.Fatalf("S103 = %+v, %v", f, err)
		}
		g, err := ie.NewS1UDataForwarding(5, s, 0x12345678).S1UDataForwarding()
		if err != nil || g.ServingGWAddress.String() != s || g.ServingGWS1UTEID != 0x12345678 {
			t.Fatalf("S1U = %+v, %v", g, err)
		}
	}
}

// https://github.com/wireshark/wireshark/commit/50018d8994b3; TS 29.274 V20.0.0 §8.15.
func TestQoSPVIBit(t *testing.T) {
	p := ie.NewBearerQoS(1, 10, 1, 9, 0, 0, 0, 0)
	if p.Payload[0] != 0x69 {
		t.Fatalf("ARP = %#x; want 0x69", p.Payload[0])
	}
	f, err := p.BearerQoS()
	if err != nil {
		t.Fatal(err)
	}
	if f.ARP&1 != 1 || (f.ARP>>2)&15 != 10 {
		t.Fatalf("ARP = %#x", f.ARP)
	}
}

// https://github.com/wireshark/wireshark/commit/6180d0fdaa5e; TS 23.003 V20.1.0 §9.1: APN maximum is 100 octets.
func TestAPNMaximumLength(t *testing.T) {
	p := ie.NewAccessPointName(strings.Repeat("a", 50) + "." + strings.Repeat("b", 50))
	if p == nil {
		return
	}
	b, err := p.Marshal()
	if err == nil {
		t.Fatalf("accepted APN encoding of %d octets; maximum is 100", len(b)-4)
	}
}
