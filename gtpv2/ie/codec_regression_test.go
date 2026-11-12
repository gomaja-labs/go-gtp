// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"bytes"
	"encoding/binary"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"go/types"
	"net"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func codecNoPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	f()
}

func TestGroupLengthAndView(t *testing.T) {
	c := NewEPSBearerID(5)
	g := NewBearerContext(c)
	c.Payload = append(c.Payload, 0)
	c.SetLength()
	g.SetLength()
	b, err := g.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if int(binary.BigEndian.Uint16(b[1:3])) != len(b)-4 {
		t.Errorf("group header length=%d, actual=%d", g.Length, len(b)-4)
	}
	c.Payload[0] = 6
	got, err := g.EPSBearerID()
	if err != nil || got != 6 {
		t.Errorf("group accessor uses stale child: got %d, %v; want 6", got, err)
	}
}
func TestGroupShortMarshal(t *testing.T) {
	g := NewBearerContext(NewEPSBearerID(5))
	if err := g.MarshalTo(make([]byte, 4)); err == nil {
		t.Error("group MarshalTo accepted buffer that cannot hold its child")
	}
	for _, c := range []struct {
		name string
		f    func() error
	}{
		{"FQCSID", func() error { return NewFullyQualifiedCSIDFields("1.2.3.4", 1).MarshalTo(make([]byte, 5)) }},
		{"PCOContainer", func() error { return NewPCOContainer(13, []byte{42}).MarshalTo(make([]byte, 3)) }},
		{"PPP", func() error { return NewPCOPPP(1, 1, []byte{1, 2}).MarshalTo(make([]byte, 5)) }},
		{"CHAP", func() error { return NewCHAPFields([]byte{1}, "abc").MarshalTo(make([]byte, 2)) }},
		{"IPCP", func() error { return NewIPCPOption(3, []byte{1, 2, 3, 4}).MarshalTo(make([]byte, 3)) }},
	} {
		if err := c.f(); err == nil {
			t.Errorf("%s MarshalTo silently truncated output", c.name)
		}
	}
}
func TestHeaderSpare(t *testing.T) {
	i, err := Parse([]byte{Recovery, 0, 1, 0xf1, 7})
	if err != nil {
		t.Fatal(err)
	}
	b, err := i.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if b[3] != 1 {
		t.Errorf("transmitted spare header bits: %#x", b[3])
	}
}
func TestLengthOverflow(t *testing.T) {
	i := New(PrivateExtension, 0, make([]byte, 65536))
	b, err := i.Marshal()
	if err == nil {
		t.Errorf("encoded %d-byte payload with uint16 length %d", len(b)-4, i.Length)
	}
}
func TestCauseFlagMappings(t *testing.T) {
	c := NewCause(69, 0, 0, 1, nil)
	if !c.HasCS() || c.HasBCE() || !c.IsRemoteCause() || c.IsPDNConnectionIEError() {
		t.Errorf("CS bit misidentified: CS=%v BCE=%v remote=%v PDN=%v", c.HasCS(), c.HasBCE(), c.IsRemoteCause(), c.IsPDNConnectionIEError())
	}
	c = NewCause(69, 1, 0, 0, nil)
	if !c.IsPDNConnectionIEError() || c.IsRemoteCause() {
		t.Errorf("PCE bit misidentified: PDN=%v remote=%v", c.IsPDNConnectionIEError(), c.IsRemoteCause())
	}
}
func TestCauseGroupedFlags(t *testing.T) {
	g := NewBearerContext(NewCause(69, 0, 0, 1, nil))
	got, err := g.CauseFlags()
	if err != nil || got != 1 {
		t.Errorf("group CauseFlags=%d,%v; want 1", got, err)
	}
}
func TestOffendingInstance(t *testing.T) {
	c := NewCause(69, 0, 0, 0, NewEPSBearerID(5).WithInstance(7))
	got, err := c.OffendingIE()
	if err != nil {
		t.Fatal(err)
	}
	if got.Instance() != 7 {
		t.Errorf("offending instance=%d; want 7", got.Instance())
	}
}
func TestGroupClassifierRace(t *testing.T) {
	old := isGroupedFun
	defer SetIsGroupedFun(old)
	f := func(t uint8) bool { return t == BearerContext }
	SetIsGroupedFun(f)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for n := 0; n < 10000; n++ {
			SetIsGroupedFun(f)
		}
	}()
	go func() {
		defer wg.Done()
		i := &IE{Type: BearerContext}
		for n := 0; n < 10000; n++ {
			i.IsGrouped()
		}
	}()
	wg.Wait()
}
func TestBounds(t *testing.T) {
	for _, c := range []struct {
		name string
		f    func()
	}{
		{"PCOPPP length below header", func() { _, _ = ParsePCOPPP([]byte{1, 1, 0, 0, 0}) }},
		{"IPCP length below header", func() { _, _ = ParseIPCPOption([]byte{3, 0, 0}) }},
		{"PAP missing password length", func() { _, _ = ParsePAPFields([]byte{1, 'a'}) }},
		{"PCO MarshalTo empty", func() { _ = NewProtocolConfigurationOptionsFields(0).MarshalTo(nil) }},
		{"PCO container MarshalTo empty", func() { _ = NewPCOContainer(13, nil).MarshalTo(nil) }},
		{"PAP MarshalTo missing password length", func() { _ = NewPAPFields("a", "").MarshalTo(make([]byte, 2)) }},
		{"S103 MarshalTo missing count", func() { _ = NewS103PDNDataForwardingInfoFields(net.ParseIP("192.0.2.1"), 1).MarshalTo(make([]byte, 9)) }},
		{"PCO nil option", func() { _ = NewProtocolConfigurationOptions(0, nil) }},
		{"TFT delete MarshalTo short", func() {
			_ = NewTrafficFlowTemplate(TFTOpDeletePacketFiltersFromExistingTFT, nil, []uint8{1}, nil).MarshalTo(make([]byte, 1))
		}},
	} {
		t.Run(c.name, func(t *testing.T) { codecNoPanic(t, c.f) })
	}
}
func TestReceiverReset(t *testing.T) {
	t.Run("PCO", func(t *testing.T) {
		f := &ProtocolConfigurationOptionsFields{}
		for n := 0; n < 2; n++ {
			if err := f.UnmarshalBinary([]byte{0x80, 0, 13, 0}); err != nil {
				t.Fatal(err)
			}
		}
		if len(f.ProtocolOrContainers) != 1 {
			t.Errorf("reused receiver accumulated %d containers", len(f.ProtocolOrContainers))
		}
	})
	t.Run("PCOContainer", func(t *testing.T) {
		f := &PCOContainer{}
		if err := f.UnmarshalBinary([]byte{0, 13, 1, 42}); err != nil {
			t.Fatal(err)
		}
		if err := f.UnmarshalBinary([]byte{0, 13, 0}); err != nil {
			t.Fatal(err)
		}
		if len(f.Contents) != 0 {
			t.Errorf("zero-length decode retained %x", f.Contents)
		}
	})
	t.Run("TFT", func(t *testing.T) {
		f := &TrafficFlowTemplate{}
		if err := f.UnmarshalBinary([]byte{0xd0, 3, 1, 1}); err != nil {
			t.Fatal(err)
		}
		if err := f.UnmarshalBinary([]byte{0x40}); err != nil {
			t.Fatal(err)
		}
		if len(f.Parameters) != 0 {
			t.Errorf("no-parameter decode retained %d parameters", len(f.Parameters))
		}
	})
	t.Run("FTEID", func(t *testing.T) {
		f := NewFullyQualifiedTEIDFields(10, 1, net.ParseIP("192.0.2.1"), nil)
		b, _ := NewFullyQualifiedTEIDFields(10, 2, nil, net.ParseIP("2001:db8::1")).Marshal()
		if err := f.UnmarshalBinary(b); err != nil {
			t.Fatal(err)
		}
		if len(f.IPv4Address) != 0 {
			t.Errorf("IPv6-only decode retained IPv4 %v", f.IPv4Address)
		}
	})
	t.Run("ULI", func(t *testing.T) {
		f := NewUserLocationInformationFields(NewCGI("123", "45", 1, 2), nil, nil, nil, nil, nil, nil, nil)
		if err := f.UnmarshalBinary([]byte{0}); err != nil {
			t.Fatal(err)
		}
		if f.CGI != nil {
			t.Error("flagless decode retained CGI")
		}
	})
	t.Run("FQCSID", func(t *testing.T) {
		f := &FullyQualifiedCSIDFields{}
		for n := 0; n < 2; n++ {
			if err := f.UnmarshalBinary([]byte{1, 1, 2, 3, 4, 0, 1}); err != nil {
				t.Fatal(err)
			}
		}
		if len(f.CSIDs) != 1 {
			t.Errorf("reused receiver accumulated %d CSIDs", len(f.CSIDs))
		}
	})
	t.Run("IE", func(t *testing.T) {
		f := NewBearerContext(NewEPSBearerID(5))
		if err := f.UnmarshalBinary([]byte{Recovery, 0, 1, 0, 1}); err != nil {
			t.Fatal(err)
		}
		if len(f.ChildIEs) != 0 {
			t.Error("non-grouped IE retained old children")
		}
	})
}
func TestFlowUplink(t *testing.T) {
	got, err := NewFlowQoS(9, 0x0102030405, 2, 3, 4).MBRForUplink()
	if err != nil || got != 0x0102030405 {
		t.Errorf("FlowQoS uplink=%x,%v; want 0102030405", got, err)
	}
}
func TestPagingFlags(t *testing.T) {
	got, err := NewPagingAndServiceInformation(6, 1, 42).PagingPolicyIndication()
	if err != nil || got != 42 {
		t.Errorf("PPI=%d,%v; want 42", got, err)
	}
	if _, err := NewPagingAndServiceInformation(5, 0, 0).PagingAndServiceInformation(); err != nil {
		t.Errorf("valid no-PPI IE rejected: %v", err)
	}
}
func TestTFTParameterFlag(t *testing.T) {
	f := NewTrafficFlowTemplate(TFTOpNoTFTOperation, nil, nil, []*TFTParameter{NewTFTParameter(1, []byte{1}), NewTFTParameter(2, []byte{0, 0, 0, 1})})
	b, err := f.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if b[0]&16 == 0 {
		t.Errorf("E=0 despite two parameters: %x", b)
	}
}
func TestTFTIPv6Offsets(t *testing.T) {
	t.Run("mask", func(t *testing.T) {
		f := NewTFTPFComponentIPv6RemoteAddress(net.ParseIP("2001:db8::1"), net.CIDRMask(64, 128))
		got, err := f.IPv6RemoteAddress()
		if err != nil || !bytes.Equal(got.Mask, net.CIDRMask(64, 128)) {
			t.Errorf("IPv6 mask=%v,%v; want 16-byte /64", got, err)
		}
	})
	t.Run("remote prefix", func(t *testing.T) {
		codecNoPanic(t, func() {
			f := NewTFTPFComponentIPv6RemoteAddressPrefixLength(net.ParseIP("2001:db8::1"), 64)
			_, _ = f.IPv6RemoteAddressPrefixLength()
		})
	})
	t.Run("local prefix", func(t *testing.T) {
		codecNoPanic(t, func() {
			f := NewTFTPFComponentIPv6LocalAddressPrefixLength(net.ParseIP("2001:db8::1"), 64)
			_, _ = f.IPv6LocalAddressPrefixLength()
		})
	})
}
func TestPCOTruncatedContainer(t *testing.T) {
	if _, err := ParsePCOContainer([]byte{0, 13, 1}); err == nil {
		t.Error("PCO container accepted missing declared content")
	}
	if _, err := ParseProtocolConfigurationOptionsFields([]byte{0x80, 0, 13, 1}); err == nil {
		t.Error("PCO accepted truncated container")
	}
}
func TestTimerUnitDecoding(t *testing.T) {
	for _, raw := range []byte{0xa2, 0xc2} {
		got, err := NewEPCTimerRaw(raw).Timer()
		if err != nil || got != 2*time.Minute {
			t.Errorf("EPC %x=%v,%v; want 2m", raw, got, err)
		}
		f, err := ParseThrottlingFields([]byte{raw, 50})
		if err != nil || f.DelayValue != 2*time.Minute {
			t.Errorf("Throttling %x=%v,%v; want 2m", raw, f, err)
		}
	}

}
func TestTimerConstructorLoss(t *testing.T) {
	for _, i := range []*IE{NewEPCTimer(32 * time.Minute), NewThrottling(32*time.Minute, 50)} {
		if i.Err() == nil {
			t.Error("unrepresentable timer accepted")
		}
	}
}
func TestRFSPWidth(t *testing.T) {
	if got := NewRFSPIndex(1).Payload; !bytes.Equal(got, []byte{0, 1}) {
		t.Errorf("RFSP wire=%x; want 0001", got)
	}
	got, err := New(RFSPIndex, 0, []byte{0, 1}).RFSPIndex()
	if err != nil || got != 1 {
		t.Errorf("RFSP decode=%v,%v; want 1", got, err)
	}
}
func TestDelayUnit(t *testing.T) {
	got, err := NewDelayValueRaw(10).DelayValue()
	if err != nil || got != 500*time.Millisecond {
		t.Errorf("DelayValue(10)=%v,%v; want 500ms", got, err)
	}
}
func TestBooleanMasks(t *testing.T) {
	b := NewBearerFlags(0, 1, 1, 1, 1)
	for _, c := range []struct {
		name string
		got  bool
	}{{"ASI", b.ActivityStatusIndicator()}, {"Vind", b.VSRVCC()}, {"VB", b.VoiceBearer()}, {"LMRI", NewMBMSFlags(1, 0).LocalMBMSBearerContextRelease()}} {
		if !c.got {
			t.Errorf("%s set but accessor false", c.name)
		}
	}
}
func TestIntegerWidths(t *testing.T) {
	for _, c := range []struct {
		p    []byte
		want uint32
	}{{[]byte{42}, 42}, {[]byte{0, 1, 0, 0}, 65536}} {
		got, err := New(IntegerNumber, 0, c.p).IntegerNumber()
		if err != nil || got != c.want {
			t.Errorf("IntegerNumber(%x)=%v,%v; want %d", c.p, got, err, c.want)
		}
	}
}
func TestIndicationTenOctets(t *testing.T) {
	p := make([]byte, 10)
	p[9] = 1
	got := NewIndicationFromOctets(p...)
	if !bytes.Equal(got.Payload, p) {
		t.Errorf("Indication payload=%x; want %x", got.Payload, p)
	}
}
func TestMinimumPayloads(t *testing.T) {
	for _, p := range [][]byte{{1}, {1, 2}} {
		if got, err := New(PTMSISignature, 0, p).PTMSISignature(); err == nil {
			t.Errorf("short P-TMSI signature %x accepted as %x", p, got)
		}
	}
	if _, err := New(Cause, 0, []byte{69}).Cause(); err == nil {
		t.Error("Cause accepted missing mandatory flags octet")
	}
	if NewIndicationFromOctets(1).Err() == nil {
		t.Error("Indication constructor accepted less than 2 mandatory octets")
	}
}
func TestSpareMasks(t *testing.T) {
	for _, c := range []struct {
		name string
		get  func() (uint8, error)
		want uint8
	}{
		{"EBI", New(EPSBearerID, 0, []byte{0xf5}).EPSBearerID, 5},
		{"PDN", New(PDNType, 0, []byte{0xf9}).PDNType, 1},
		{"PAA PDN", New(PDNAddressAllocation, 0, []byte{0xf9, 192, 0, 2, 1}).PDNType, 1},
		{"CMI", New(CSGMembershipIndication, 0, []byte{0xff}).CMI, 1},
		{"UCI CMI", New(UserCSGInformation, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0x41}).CMI, 1},
		{"selection", New(SelectionMode, 0, []byte{0xfd}).SelectionMode, 1},
		{"DST", New(UETimeZone, 0, []byte{0, 0xfd}).DaylightSaving, 1},
	} {
		got, err := c.get()
		if err != nil || got != c.want {
			t.Errorf("%s=%v,%v; want %d", c.name, got, err, c.want)
		}
	}
	got, err := New(CSGID, 0, []byte{0xf8, 0, 0, 1}).CSGID()
	if err != nil || got != 1 {
		t.Errorf("CSGID=%x,%v; want 1", got, err)
	}
	for _, c := range []struct {
		name string
		p    []byte
		get  func(*UserLocationInformationFields) uint32
		want uint32
	}{
		{"ECGI", []byte{0x10, 0x21, 0xf3, 0x54, 0xf1, 0x23, 0x45, 0x67}, func(f *UserLocationInformationFields) uint32 { return f.ECGI.ECI }, 0x1234567},
		{"MENBI", []byte{0x40, 0x21, 0xf3, 0x54, 0xf1, 0x23, 0x45}, func(f *UserLocationInformationFields) uint32 { return f.MENBI.MENBI }, 0x12345},
		{"EMENBI", []byte{0x80, 0x21, 0xf3, 0x54, 0x80, 0, 1}, func(f *UserLocationInformationFields) uint32 { return f.EMENBI.EMENBI }, 1},
	} {
		f, err := ParseUserLocationInformationFields(c.p)
		if err != nil {
			t.Error(err)
			continue
		}
		if got := c.get(f); got != c.want {
			t.Errorf("%s=%x; want %x", c.name, got, c.want)
		}
	}
}
func TestSelectionReserved(t *testing.T) {
	got, err := New(SelectionMode, 0, []byte{3}).SelectionMode()
	if err != nil || got != 2 {
		t.Errorf("selection3=%d,%v; want2", got, err)
	}
}
func TestEmptyPrivateExtension(t *testing.T) {
	got, err := NewPrivateExtension(1, nil).PrivateExtension()
	if err != nil || len(got) != 0 {
		t.Errorf("empty proprietary value=%x,%v; want empty success", got, err)
	}
}
func TestEvenIMSI(t *testing.T) {
	want := "00101012345678"
	got, err := NewIMSI(want).IMSI()
	if err != nil || got != want {
		t.Errorf("IMSI=%q,%v; want %q", got, err, want)
	}
}
func TestDigitValidation(t *testing.T) {
	for _, c := range []struct {
		name string
		i    *IE
	}{{"hex IMSI", NewIMSI("0010101234567a")}, {"long IMSI", NewIMSI("0010101234567890")}, {"hex MSISDN", NewMSISDN("123a")}, {"hex MEI", NewMobileEquipmentIdentity("12345678901234a")}, {"short MEI", NewMobileEquipmentIdentity("1")}} {
		if c.i.Err() == nil {
			t.Errorf("%s accepted as %x", c.name, c.i.Payload)
		}
	}
	if _, err := New(IMSI, 0, []byte{0xfa}).IMSI(); err == nil {
		t.Error("invalid decimal IMSI nibble accepted")
	}
}
func TestGUTIOffset(t *testing.T) {
	codecNoPanic(t, func() { _, _ = ParseGUTIFields([]byte{0x21, 0xf3, 0x54}) })
	got, err := NewGUTI("123", "45", 0x1111, 0x22, 0x33333333).GUTI()
	if err != nil || got.MCC != "123" || got.MNC != "45" {
		t.Errorf("GUTI=%+v,%v; want MCC123 MNC45", got, err)
	}
}
func TestPLMNValidation(t *testing.T) {
	if _, err := New(ServingNetwork, 0, []byte{0xba, 0xfc, 0x54}).ServingNetwork(); err == nil {
		t.Error("nondecimal PLMN decoded without error")
	}

	for _, c := range []struct {
		name string
		f    func() *IE
	}{{"short MCC", func() *IE { return NewServingNetwork("1", "123") }}, {"long MCC", func() *IE { return NewServingNetwork("1234", "45") }}, {"hex MCC", func() *IE { return NewServingNetwork("abc", "45") }}} {
		t.Run(c.name, func(t *testing.T) {
			codecNoPanic(t, func() {
				if got := c.f(); got.Err() == nil {
					t.Errorf("invalid PLMN accepted: %x", got.Payload)
				}
			})
		})
	}
}
func TestGUTISwallowedError(t *testing.T) {
	if got, err := NewGUTIFields("abcZ", "45", 1, 2, 3).Marshal(); err == nil {
		t.Errorf("invalid PLMN encoded as %x without error", got)
	}
}
func TestRAIWidth(t *testing.T) {
	i := NewUserLocationInformation(nil, nil, NewRAI("123", "45", 0x1234, 0x56), nil, nil, nil, nil, nil)
	if got := i.Payload[6:8]; !bytes.Equal(got, []byte{0x56, 0xff}) {
		t.Errorf("RAC wire=%x; want 56ff", got)
	}
	f, err := ParseUserLocationInformationFields([]byte{4, 0x21, 0xf3, 0x54, 0x12, 0x34, 0x56, 0xff})
	if err != nil || f.RAI.RAC != 0x56 {
		t.Errorf("RAI decode=%+v,%v; want RAC56", f, err)
	}
}
func TestULIExclusiveENB(t *testing.T) {
	if NewUserLocationInformation(nil, nil, nil, nil, nil, nil, NewMENBI("123", "45", 1), NewEMENBI("123", "45", 1)).Err() == nil {
		t.Error("mutually exclusive MENBI and EMENBI accepted")
	}
}
func TestFTEIDAddressFlags(t *testing.T) {
	if got := NewFullyQualifiedTEID(0x80, 1, "", "2001:db8::1"); got.Err() == nil {
		t.Errorf("out-of-range interface type forged V4 flag: %x", got.Payload)
	}
	if NewFullyQualifiedTEID(10, 1, "", "").Err() == nil {
		t.Error("F-TEID without any address accepted")
	}
	if _, err := ParseFullyQualifiedTEIDFields([]byte{10, 0, 0, 0, 1}); err == nil {
		t.Error("F-TEID without V4/V6 flag decoded successfully")
	}
}
func TestFQCSIDCount(t *testing.T) {
	if _, err := ParseFullyQualifiedCSIDFields([]byte{2, 1, 2, 3, 4, 0, 1}); err == nil {
		t.Error("count2 with only 1 CSID accepted")
	}
	f, err := ParseFullyQualifiedCSIDFields([]byte{1, 1, 2, 3, 4, 0, 1, 0xaa, 0xbb})
	if err != nil || len(f.CSIDs) != 1 {
		t.Errorf("extension treated as CSID: %+v,%v", f, err)
	}
}
func TestFQCSIDNodeValidation(t *testing.T) {
	codecNoPanic(t, func() {
		if NewFullyQualifiedCSID("not-an-ip", 1).Err() == nil {
			t.Error("invalid node accepted")
		}
	})
	if got := NewFullyQualifiedCSID("12", 1); got.Err() == nil {
		t.Errorf("type2 node with 1 instead of 4 bytes accepted: %x", got.Payload)
	}
}
func TestTimestampEpoch(t *testing.T) {
	want := time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := NewULITimestamp(want).TimestampNear(want)
	if err != nil || !got.Equal(want) {
		t.Errorf("timestamp=%v,%v; want %v", got, err, want)
	}
}
func TestNameValidation(t *testing.T) {
	for _, typ := range []uint8{AccessPointName, FullyQualifiedDomainName} {
		if got, err := New(typ, 0, []byte{1, 'a', 3, 'b'}).ValueAsFQDN(); err == nil {
			t.Errorf("type%d truncated label decoded as %q without error", typ, got)
		}
	}
	for _, s := range []string{"a.", "a..b", strings.Repeat("a", 64), strings.Repeat("a.", 51) + "a"} {
		if got := NewAccessPointName(s); got.Err() == nil {
			t.Errorf("invalid APN %q accepted", s)
		}
	}
}
func TestIPLength(t *testing.T) {
	for _, c := range []struct {
		name string
		i    *IE
	}{{"S103", NewS103PDNDataForwardingInfo("invalid", 1, 5)}, {"S1U", NewS1UDataForwarding(5, "invalid", 1)}} {
		if c.i.Err() == nil {
			t.Errorf("%s constructor emitted zero-length address", c.name)
		}
	}
	if _, e := ParseS1UDataForwardingFields([]byte{5, 0, 0, 0, 0, 1}); e == nil {
		t.Error("S1U decoder accepted zero-length address")
	}
	if _, e := ParseS103PDNDataForwardingInfoFields([]byte{0, 0, 0, 0, 1, 0}); e == nil {
		t.Error("S103 decoder accepted zero-length address")
	}

	if got, err := New(IPAddress, 0, []byte{1, 2, 3, 4, 5}).IPAddress(); err == nil {
		t.Errorf("5-byte IP accepted as %q", got)
	}
	if NewIPAddress("invalid").Err() == nil {
		t.Error("invalid IP constructor accepted")
	}
}
func TestPAAValidation(t *testing.T) {
	if _, err := ParsePDNAddressAllocationFields([]byte{0}); err == nil {
		t.Error("reserved PDN type0 accepted")
	}
	if NewPDNAddressAllocationIPv6("2001:db8::1", 129).Err() == nil {
		t.Error("IPv6 prefix129 accepted")
	}
}
func TestNameMap(t *testing.T) {
	for _, c := range []struct {
		typ  uint8
		want string
	}{{PLMNID, "PLMNID"}, {TargetIdentification, "TargetIdentification"}} {
		if got := New(c.typ, 0, nil).Name(); got != c.want {
			t.Errorf("IE%d Name=%q; want %q", c.typ, got, c.want)
		}
	}
}
func TestThrottleFactor(t *testing.T) {
	f, err := ParseThrottlingFields([]byte{1, 255})
	if err != nil || f.Factor != 0 {
		t.Errorf("reserved factor=%v,%v; want 0", f, err)
	}
}
func TestS103Count(t *testing.T) {
	i := New(S103PDNDataForwardingInfo, 0, []byte{4, 127, 0, 0, 1, 0, 0, 0, 1, 2, 5})
	if got, err := i.EBIs(); err == nil {
		t.Errorf("truncated EBI list accepted: %v", got)
	}
}
func TestLDNValidation(t *testing.T) {
	for _, s := range []string{"", strings.Repeat("a", 401), "é"} {
		if i := NewLocalDistinguishedName(s); i != nil {
			if got, err := i.LocalDistinguishedName(); err == nil {
				t.Errorf("invalid LDN (%d bytes) accepted: %q", len(s), got)
			}
		}
	}
}
func TestIndicationBinary(t *testing.T) {
	if NewIndicationFromBitSequence(strings.Repeat("2", 72)).Err() == nil {
		t.Error("nonbinary indication sequence accepted as all ones")
	}
}
func TestEnumDomains(t *testing.T) {
	for _, c := range []struct {
		name string
		i    *IE
	}{{"PDN0", NewPDNType(0)}, {"Selection3", NewSelectionMode(3)}, {"NodeType2", NewNodeType(2)}, {"RAT0", NewRATType(0)}, {"Service0", NewServiceIndicator(0)}, {"Detach0", NewDetachType(0)}, {"CMI2", NewCSGMembershipIndication(2)}, {"RFSP0", NewRFSPIndex(0)}} {
		if c.i.Err() == nil {
			t.Errorf("forbidden %s emitted: %x", c.name, c.i.Payload)
		}
	}
}
func TestUCIConditionalCMI(t *testing.T) {
	i := NewUserCSGInformation("123", "45", 1, 0, 0, 1)
	if i.Payload[7]&1 != 0 {
		t.Error("Closed-mode UCI transmits CMI=1")
	}
}
func TestRANNASCauseWidth(t *testing.T) {
	if _, err := ParseRANNASCauseFields([]byte{0x40, 1}); err == nil {
		t.Error("Diameter cause accepted with 1 instead of 2 octets")
	}
	f, err := ParseRANNASCauseFields([]byte{0x10, 1, 0xaa, 0xbb})
	if err != nil || len(f.Cause) != 1 {
		t.Errorf("S1AP extension treated as Cause: %v,%v", f, err)
	}
	if NewRANNASCause(4, 0, nil).Err() == nil {
		t.Error("constructor emits Diameter cause without mandatory value")
	}
}
func TestTFTValidation(t *testing.T) {
	for _, c := range []struct {
		name string
		p    []byte
	}{{"reserved opcode", []byte{0xe0}}, {"empty create", []byte{0x20}}, {"empty filter", []byte{0x21, 0x31, 0, 0}}, {"reserved component", []byte{0x21, 0x31, 0, 1, 0xff}}} {
		if got, err := ParseTrafficFlowTemplate(c.p); err == nil {
			t.Errorf("%s accepted: %+v", c.name, got)
		}
	}
	f := NewTrafficFlowTemplate(TFTOpCreateNewTFT, []*TFTPacketFilter{NewTFTPacketFilter(3, 1, 1, NewTFTPFComponent(0xff, nil))}, nil, nil)
	if b, err := f.Marshal(); err == nil {
		t.Errorf("reserved component encoded: %x", b)
	}
}
func TestNarrowLengthOverflow(t *testing.T) {
	for _, c := range []struct {
		name string
		f    func() ([]byte, error)
	}{
		{"PCO", NewPCOContainer(13, make([]byte, 256)).Marshal},
		{"PPP", NewPCOPPP(1, 1, make([]byte, 65532)).Marshal},
		{"PAP", NewPAPFields(strings.Repeat("a", 256), "p").Marshal},
		{"CHAP", NewCHAPFields(make([]byte, 256), "n").Marshal},
		{"IPCP", NewIPCPOption(3, make([]byte, 254)).Marshal},
		{"TFT parameter", NewTFTParameter(1, make([]byte, 256)).Marshal},
		{"TFT count", NewTrafficFlowTemplate(TFTOpDeletePacketFiltersFromExistingTFT, nil, make([]uint8, 16), nil).Marshal},
		{"FQCSID count", NewFullyQualifiedCSIDFields("192.0.2.1", make([]uint16, 16)...).Marshal},
	} {
		if b, err := c.f(); err == nil {
			t.Errorf("%s accepted overflowing length/count (wire %d bytes)", c.name, len(b))
		}
	}
}
func TestNarrowValueOverflow(t *testing.T) {
	for _, c := range []struct {
		name string
		i    *IE
	}{
		{"Delay", NewDelayValue(12800 * time.Millisecond)},
		{"PTMSI", NewPTMSISignature(0x1000001)},
		{"CSG", NewCSGID(0x8000001)},
		{"EBI", NewEPSBearerID(0x15)},
		{"TraceID", NewTraceReference("123", "45", 0x1000001)},
		{"QoS", NewFlowQoS(9, 1<<40, 0, 0, 0)},
	} {
		if c.i.Err() == nil {
			t.Errorf("%s silently narrowed out-of-range input: %x", c.name, c.i.Payload)
		}
	}
}
func TestTFTAndS103SpareMasks(t *testing.T) {
	for _, c := range []struct {
		name string
		get  func() (uint32, error)
		want uint32
	}{
		{"flow label", func() (uint32, error) { return NewTFTPFComponent(PFCompFlowLabel, []byte{0xf1, 2, 3}).FlowLabel() }, 0x10203},
		{"CTAG VID", func() (uint32, error) {
			v, e := NewTFTPFComponent(PFCompDot1QCTAGVID, []byte{0xf1, 2}).Dot1QCTAGVID()
			return uint32(v), e
		}, 0x102},
		{"STAG VID", func() (uint32, error) {
			v, e := NewTFTPFComponent(PFCompDot1QSTAGVID, []byte{0xf1, 2}).Dot1QSTAGVID()
			return uint32(v), e
		}, 0x102},
		{"CTAG PCPDEI", func() (uint32, error) {
			v, e := NewTFTPFComponent(PFCompDot1QCTAGPCPDEI, []byte{0xf1}).Dot1QCTAGPCPDEI()
			return uint32(v), e
		}, 1},
		{"STAG PCPDEI", func() (uint32, error) {
			v, e := NewTFTPFComponent(PFCompDot1QSTAGPCPDEI, []byte{0xf1}).Dot1QSTAGPCPDEI()
			return uint32(v), e
		}, 1},
	} {
		got, err := c.get()
		if err != nil || got != c.want {
			t.Errorf("%s=%x,%v; want %x", c.name, got, err, c.want)
		}
	}
	i := New(S103PDNDataForwardingInfo, 0, []byte{4, 127, 0, 0, 1, 0, 0, 0, 1, 1, 0xf5})
	v, err := i.EBIs()
	if err != nil || len(v) != 1 || v[0] != 5 {
		t.Errorf("S103 EBI=%v,%v; want [5]", v, err)
	}
	f, err := i.S103PDNDataForwardingInfo()
	if err != nil || len(f.EPSBearerIDs) != 1 || f.EPSBearerIDs[0] != 5 {
		t.Errorf("S103 fields EBI=%+v,%v; want [5]", f, err)
	}
}
func TestEmptyPPPConstructor(t *testing.T) {
	i := NewPCOPPP(2, 1, nil)
	if b, err := i.Marshal(); err != nil || len(b) != 4 {
		t.Errorf("empty PPP constructor cannot marshal: %x,%v", b, err)
	}
}
func TestOriginalLeads(t *testing.T) {
	t.Run("ARP and BearerQoS", func(t *testing.T) {
		for _, i := range []*IE{NewAllocationRetensionPriority(1, 9, 1), NewBearerQoS(1, 9, 1, 9, 1, 2, 3, 4)} {
			v, e := i.PriorityLevel()
			if e != nil || v != 9 {
				t.Fatalf("priority=%d,%v", v, e)
			}
		}
		b := NewBearerQoS(1, 9, 1, 9, 1, 2, 3, 4)
		got, e := b.BearerQoS()
		if e != nil || got.ARP != 0x65 || got.QCI != 9 || got.MaximumBitRateForUplink != 1 {
			t.Fatalf("BearerQoS=%+v,%v", got, e)
		}
	})
	t.Run("PAA", func(t *testing.T) {
		i := NewPDNAddressAllocationDual("192.0.2.1", "2001:db8::1", 64)
		b, e := i.Marshal()
		if e != nil || len(b) != 26 {
			t.Fatalf("PAA wire length=%d,%v", len(b), e)
		}
		f, e := ParsePDNAddressAllocationFields(i.Payload)
		if e != nil || !f.IPv4Address.Equal(net.ParseIP("192.0.2.1")) || !f.IPv6Address.Equal(net.ParseIP("2001:db8::1")) {
			t.Fatalf("PAA=%+v,%v", f, e)
		}
	})
	t.Run("PCO copying", func(t *testing.T) {
		b := []byte{0, 13, 1, 42}
		c, e := ParsePCOContainer(b)
		if e != nil {
			t.Fatal(e)
		}
		b[3] = 0
		if !bytes.Equal(c.Contents, []byte{42}) {
			t.Fatalf("PCO contents=%x", c.Contents)
		}
	})
	t.Run("same type and instance", func(t *testing.T) {
		g := NewBearerContext(New(EPSBearerID, 0, []byte{5}), New(EPSBearerID, 0, []byte{6}))
		b, e := g.Marshal()
		if e != nil {
			t.Fatal(e)
		}
		p, e := Parse(b)
		if e != nil || len(p.ChildIEs) != 2 {
			t.Fatalf("duplicate children lost: %+v,%v", p, e)
		}
		v, e := p.BearerContext()
		if e != nil || len(v) != 2 {
			t.Fatalf("duplicate getter lost: %+v,%v", v, e)
		}
	})
}
func TestCommonCodecCoverage(t *testing.T) {
	for _, name := range []string{"ExtendedProtocolConfigurationOptions", "AdditionalProtocolConfigurationOptions", "SecondaryRATUsageDataReport"} {
		if _, ok := reflect.TypeOf(&IE{}).MethodByName(name); !ok {
			t.Errorf("no typed %s decoder", name)
		}
	}
}

func codecParentConstant(t *testing.T, name string) uint8 {
	t.Helper()
	fs := token.NewFileSet()
	f, e := parser.ParseFile(fs, "../constants.go", nil, 0)
	if e != nil {
		t.Fatal(e)
	}
	cfg := types.Config{}
	p, e := cfg.Check("gtpv2", fs, []*ast.File{f}, nil)
	if e != nil {
		t.Fatal(e)
	}
	c, ok := p.Scope().Lookup(name).(*types.Const)
	if !ok {
		t.Fatalf("missing constant %s", name)
	}
	n, ok := constant.Uint64Val(c.Val())
	if !ok || n > 255 {
		t.Fatalf("invalid constant %s", name)
	}
	return uint8(n)
}
func TestNodeTypeConstants(t *testing.T) {
	for _, c := range []struct {
		name string
		want byte
	}{{"NodeTypeMME", 0}, {"NodeTypeSGSN", 1}} {
		i := NewNodeType(codecParentConstant(t, c.name))
		if i.Payload[0] != c.want {
			t.Errorf("%s emits %d; want %d", c.name, i.Payload[0], c.want)
		}
	}
}
func TestCMIContextConstants(t *testing.T) {
	i := NewCSGMembershipIndication(codecParentConstant(t, "CMICSG"))
	if i.Payload[0] != 0 {
		t.Errorf("CMICSG emits %d (non-member); want 0 (member)", i.Payload[0])
	}
}
