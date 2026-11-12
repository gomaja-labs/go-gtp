// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestNilCodecs(t *testing.T) {
	values := []any{(*IE)(nil), (*AggregateMaximumBitRateFields)(nil), (*BearerQoSFields)(nil), (*FlowQoSFields)(nil), (*FullyQualifiedTEIDFields)(nil), (*FullyQualifiedCSIDFields)(nil), (*GUTIFields)(nil), (*PDNAddressAllocationFields)(nil), (*PagingAndServiceInformationFields)(nil), (*ProtocolConfigurationOptionsFields)(nil), (*PCOContainer)(nil), (*PCOPPP)(nil), (*PAPFields)(nil), (*CHAPFields)(nil), (*IPCPOption)(nil), (*TrafficFlowTemplate)(nil), (*TFTPacketFilter)(nil), (*TFTPFComponent)(nil), (*TFTParameter)(nil), (*UserLocationInformationFields)(nil), (*UserCSGInformationFields)(nil), (*ThrottlingFields)(nil), (*TraceReferenceFields)(nil), (*RANNASCauseFields)(nil), (*S103PDNDataForwardingInfoFields)(nil), (*S1UDataForwardingFields)(nil), (*ExtendedProtocolConfigurationOptionsFields)(nil), (*SecondaryRATUsageDataReportFields)(nil)}
	for _, value := range values {
		v := reflect.ValueOf(value)
		t.Run(v.Type().String(), func(t *testing.T) {
			for n := 0; n < v.NumMethod(); n++ {
				m := v.Method(n)
				if m.Type().NumIn() == 0 {
					codecNoPanic(t, func() { m.Call(nil) })
				}
			}
			if c, ok := value.(interface{ MarshalTo([]byte) error }); ok {
				if err := c.MarshalTo(nil); !errors.Is(err, ErrIENotFound) {
					t.Errorf("nil MarshalTo: %v", err)
				}
			}
		})
	}
}

func TestExtendedPCOWireAndBounds(t *testing.T) {
	c := &ExtendedPCOContainer{ID: 0x41, Contents: bytes.Repeat([]byte{42}, 256)}
	i := NewExtendedProtocolConfigurationOptions(0, PCONetworkToMS, c, c)
	b, err := i.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 525 || !bytes.Equal(i.Payload[:5], []byte{0x80, 0, 0x41, 1, 0}) {
		t.Fatalf("wire length/header: %d %x", len(b), i.Payload[:5])
	}
	f, err := i.ExtendedProtocolConfigurationOptions(PCONetworkToMS)
	if err != nil || len(f.ProtocolOrContainers) != 2 || !bytes.Equal(f.ProtocolOrContainers[1].Contents, c.Contents) {
		t.Fatalf("decode: %+v %v", f, err)
	}
	for n := 0; n < f.MarshalLen(); n++ {
		if err := f.MarshalTo(make([]byte, n)); err == nil {
			t.Fatalf("accepted short destination %d", n)
		}
	}
	for _, p := range [][]byte{{}, {0x80, 0}, {0x80, 0, 0x41, 0}, {0x80, 0, 0x41, 0, 1}} {
		if err := f.UnmarshalBinary(p); err == nil {
			t.Fatalf("accepted truncated ePCO: %x", p)
		}
	}
	if err := f.UnmarshalBinary([]byte{0x80}); err != nil || len(f.ProtocolOrContainers) != 0 {
		t.Fatalf("receiver reset: %v", err)
	}
	if NewExtendedProtocolConfigurationOptions(0, PCONetworkToMS, nil).Err() == nil {
		t.Fatal("nil ePCO container accepted")
	}
}

func TestAdditionalPCOWire(t *testing.T) {
	i := NewAdditionalProtocolConfigurationOptions(0, NewPCOContainer(13, []byte{1}), NewPCOContainer(13, nil))
	f, err := i.AdditionalProtocolConfigurationOptions()
	if err != nil || len(f.ProtocolOrContainers) != 2 || !bytes.Equal(i.Payload, []byte{0x80, 0, 13, 1, 1, 0, 13, 0}) {
		t.Fatalf("APCO: %+v %v", f, err)
	}
	if NewAdditionalProtocolConfigurationOptions(0, nil).Err() == nil {
		t.Fatal("nil APCO container accepted")
	}
}

func TestSecondaryRATUsageWireAndBounds(t *testing.T) {
	f := &SecondaryRATUsageDataReportFields{Flags: 3, SecondaryRATType: 1, EPSBearerID: 5, StartTimestamp: 0x01020304, EndTimestamp: 0x05060708, UsageDataDL: 0x0102030405060708, UsageDataUL: 0x1112131415161718}
	i := NewSecondaryRATUsageDataReport(f)
	want := []byte{3, 1, 5, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18}
	if !bytes.Equal(i.Payload, want) {
		t.Fatalf("wire: %x", i.Payload)
	}
	got, err := i.SecondaryRATUsageDataReport()
	if err != nil || !reflect.DeepEqual(got, f) {
		t.Fatalf("decode: %+v %v", got, err)
	}
	for n := 0; n < 27; n++ {
		if err := got.UnmarshalBinary(want[:n]); err == nil {
			t.Fatalf("accepted %d payload octets", n)
		}
		if err := f.MarshalTo(make([]byte, n)); err == nil {
			t.Fatalf("accepted %d destination octets", n)
		}
	}
	f = &SecondaryRATUsageDataReportFields{Flags: 5, EPSBearerID: 5, SecondaryRATDataUsageReportTransfer: []byte{1, 2, 3}}
	i = NewSecondaryRATUsageDataReport(f)
	if len(i.Payload) != 31 || i.Payload[27] != 3 {
		t.Fatalf("NG-RAN wire: %x", i.Payload)
	}
	got, err = i.SecondaryRATUsageDataReport()
	if err != nil || !reflect.DeepEqual(got, f) {
		t.Fatalf("NG-RAN decode: %+v %v", got, err)
	}
	for n := 27; n < 31; n++ {
		if err := got.UnmarshalBinary(i.Payload[:n]); err == nil {
			t.Fatalf("accepted truncated transfer %d", n)
		}
	}
	f.Flags = 4
	if NewSecondaryRATUsageDataReport(f).Err() == nil {
		t.Fatal("missing IRPGW accepted")
	}
	f.Flags = 5
	f.UsageDataDL = 1
	if NewSecondaryRATUsageDataReport(f).Err() == nil {
		t.Fatal("legacy usage in NG-RAN accepted")
	}
}

func TestScalarBoundaries(t *testing.T) {
	for _, n := range []uint16{1, 255, 256} {
		v, err := NewRFSPIndex(n).RFSPIndex()
		if err != nil || v != n {
			t.Fatalf("RFSP %d: %d %v", n, v, err)
		}
	}
	for _, n := range []uint16{0, 257, 65535} {
		if NewRFSPIndex(n).Err() == nil {
			t.Fatalf("RFSP accepted %d", n)
		}
	}
	for _, d := range []time.Duration{0, 2 * time.Second, 31 * time.Minute, 40 * time.Minute, 310 * time.Hour, time.Duration(1<<63 - 1)} {
		v, err := NewEPCTimer(d).EPCTimer()
		if err != nil || v != d {
			t.Fatalf("timer %v: %v %v", d, v, err)
		}
	}
	for _, i := range []*IE{NewBearerFlags(2, 0, 0, 0), NewAllocationRetensionPriority(2, 0, 0), NewCause(69, 2, 0, 0, nil), NewMBMSFlags(0, 2), NewUserCSGInformation("123", "45", 1, 0, 0, 2), NewPagingAndServiceInformation(16, 0, 0), NewUETimeZone(time.Second, 0), NewRATType(23)} {
		if i.Err() == nil {
			t.Errorf("invalid type %d input accepted", i.Type)
		}
	}
	for _, n := range []uint8{10, 11, 22} {
		if NewRATType(n).Err() != nil {
			t.Errorf("current RAT %d rejected", n)
		}
	}
}

func TestStrictAddressAndTFTInputs(t *testing.T) {
	for _, i := range []*IE{NewFullyQualifiedTEID(10, 1, "invalid", "2001:db8::1"), NewFullyQualifiedTEID(10, 1, "192.0.2.1", "invalid"), NewFullyQualifiedTEIDNetIP(10, 1, []byte{1}, nil), NewPDNAddressAllocation("invalid")} {
		if i.Err() == nil {
			t.Errorf("invalid address for IE %d accepted", i.Type)
		}
	}
	for _, b := range [][]byte{{0x4f}, {0x0f}, {0x10}, {0x70}, {0xc0}} {
		if _, err := ParseTrafficFlowTemplate(b); err == nil {
			t.Errorf("invalid TFT header accepted: %x", b)
		}
	}
	if _, err := NewTFTPFComponentFlowLabel(1 << 20).Marshal(); err == nil {
		t.Fatal("flow label overflow accepted")
	}
}

func TestShortMacroENodeBForm(t *testing.T) {
	// TS 29.274 §8.21.8: short ID is 18 bits; SMeNB is bit 8.
	raw := []byte{0x80, 0x21, 0xf3, 0x54, 0x81, 0x23, 0x45}
	f, err := ParseUserLocationInformationFields(raw)
	if err != nil || !f.EMENBI.Short || f.EMENBI.EMENBI != 0x12345 {
		t.Fatalf("short macro identity: %+v %v", f, err)
	}
	b, err := f.Marshal()
	if err != nil || !bytes.Equal(b, raw) {
		t.Fatalf("short macro encoding: %x %v", b, err)
	}
	f.EMENBI.EMENBI = 1 << 18
	if _, err := f.Marshal(); err == nil {
		t.Fatal("short macro ID overflow accepted")
	}
}

func TestAdditionalCauseConstants(t *testing.T) {
	for _, c := range []struct {
		name  string
		value uint8
	}{
		{"CauseEPSTo5GSMobility", 15}, {"Cause5GCNotAllowed", 129}, {"CausePGWMismatchWithNetworkSliceSubscribedByUE", 130}, {"CauseRejectionDueToPagingRestriction", 131},
	} {
		if got := codecParentConstant(t, c.name); got != c.value {
			t.Errorf("%s=%d, want %d", c.name, got, c.value)
		}
	}
}

func TestClosedModeIgnoresMembership(t *testing.T) {
	i := New(UserCSGInformation, 0, []byte{0x21, 0xf3, 0x54, 0, 0, 0, 1, 1})
	v, err := i.CMI()
	if err != nil || v != 0 {
		t.Fatalf("Closed-mode CMI: %d %v", v, err)
	}
	f, err := i.UserCSGInformation()
	if err != nil || f.Flags&1 != 0 {
		t.Fatalf("Closed-mode flags: %+v %v", f, err)
	}
}

func TestExtendedPCOContainerLengthWidths(t *testing.T) {
	for _, c := range []struct {
		direction PCODirection
		id        uint16
		want      []byte
	}{
		{PCONetworkToMS, 0xff00, []byte{0x80, 0xff, 0, 1, 42}},
		{PCOMSToNetwork, 0x23, []byte{0x80, 0, 0x23, 1, 42}},
		{PCONetworkToMS, 0x23, []byte{0x80, 0, 0x23, 0, 1, 42}},
		{PCOMSToNetwork, 0x41, []byte{0x80, 0, 0x41, 0, 1, 42}},
	} {
		i := NewExtendedProtocolConfigurationOptions(0, c.direction, &ExtendedPCOContainer{ID: c.id, Contents: []byte{42}})
		if !bytes.Equal(i.Payload, c.want) {
			t.Fatalf("ID %x direction %d: %x", c.id, c.direction, i.Payload)
		}
		f, err := i.ExtendedProtocolConfigurationOptions(c.direction)
		if err != nil || len(f.ProtocolOrContainers) != 1 || !bytes.Equal(f.ProtocolOrContainers[0].Contents, []byte{42}) {
			t.Fatalf("decode: %+v %v", f, err)
		}
	}
	if NewExtendedProtocolConfigurationOptions(0, PCONetworkToMS, &ExtendedPCOContainer{ID: 0xff00, Contents: make([]byte, 256)}).Err() == nil {
		t.Fatal("one-octet length overflow accepted")
	}
	large := NewLargeAdditionalProtocolConfigurationOptions(0, PCONetworkToMS, &ExtendedPCOContainer{ID: 0x41, Contents: make([]byte, 256)})
	f, err := large.LargeAdditionalProtocolConfigurationOptions(PCONetworkToMS)
	if err != nil || len(f.ProtocolOrContainers) != 1 || len(f.ProtocolOrContainers[0].Contents) != 256 {
		t.Fatalf("large APCO: %+v %v", f, err)
	}
}

func TestPAAFixedIPv6Prefix(t *testing.T) {
	for _, prefix := range []uint8{0, 63, 65, 128, 129} {
		if NewPDNAddressAllocationIPv6("2001:db8::1", prefix).Err() == nil {
			t.Errorf("PAA prefix %d accepted", prefix)
		}
		raw := make([]byte, 18)
		raw[0] = 2
		raw[1] = prefix
		if _, err := ParsePDNAddressAllocationFields(raw); err == nil {
			t.Errorf("received PAA prefix %d accepted", prefix)
		}
	}
	if NewPDNAddressAllocationIPv6("invalid", 64).Err() == nil {
		t.Fatal("invalid IPv6 address accepted")
	}
	f, err := ParsePDNAddressAllocationFields(NewPDNAddressAllocation("2001:db8::1").Payload)
	if err != nil || f.IPv6PrefixLength != 64 {
		t.Fatalf("default IPv6 prefix: %+v %v", f, err)
	}
}

func TestIdentitySurplusFillers(t *testing.T) {
	for _, get := range []func() (string, error){
		New(IMSI, 0, []byte{0x21, 0xff}).IMSI,
		New(MSISDN, 0, []byte{0x21, 0xff}).MSISDN,
		New(MobileEquipmentIdentity, 0, []byte{0x21, 0x43, 0x65, 0x87, 0x09, 0x21, 0x43, 0xf5, 0xff}).MobileEquipmentIdentity,
	} {
		if _, err := get(); err == nil {
			t.Fatal("surplus identity filler accepted")
		}
	}
	if NewMSISDN("").Err() == nil {
		t.Fatal("empty MSISDN accepted")
	}
	if _, err := New(PDNType, 0, []byte{0xf8}).PDNType(); err == nil {
		t.Fatal("reserved PDN type accepted")
	}
}
