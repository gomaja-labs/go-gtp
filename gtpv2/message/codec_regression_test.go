package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

type codecFactory struct {
	name string
	make func(...*ie.IE) Message
}

var codecFactories = []codecFactory{
	{"Generic", func(xs ...*ie.IE) Message { return NewGeneric(254, 0, 1, xs...) }},
	{"ChangeNotificationRequest", func(xs ...*ie.IE) Message { return NewChangeNotificationRequest(0, 1, xs...) }},
	{"ChangeNotificationResponse", func(xs ...*ie.IE) Message { return NewChangeNotificationResponse(0, 1, xs...) }},
	{"ContextAcknowledge", func(xs ...*ie.IE) Message { return NewContextAcknowledge(0, 1, xs...) }},
	{"ContextRequest", func(xs ...*ie.IE) Message { return NewContextRequest(0, 1, xs...) }},
	{"ContextResponse", func(xs ...*ie.IE) Message { return NewContextResponse(0, 1, xs...) }},
	{"CreateBearerRequest", func(xs ...*ie.IE) Message { return NewCreateBearerRequest(0, 1, xs...) }},
	{"CreateBearerResponse", func(xs ...*ie.IE) Message { return NewCreateBearerResponse(0, 1, xs...) }},
	{"CreateSessionRequest", func(xs ...*ie.IE) Message { return NewCreateSessionRequest(0, 1, xs...) }},
	{"CreateSessionResponse", func(xs ...*ie.IE) Message { return NewCreateSessionResponse(0, 1, xs...) }},
	{"DeleteBearerCommand", func(xs ...*ie.IE) Message { return NewDeleteBearerCommand(0, 1, xs...) }},
	{"DeleteBearerFailureIndication", func(xs ...*ie.IE) Message { return NewDeleteBearerFailureIndication(0, 1, xs...) }},
	{"DeleteBearerRequest", func(xs ...*ie.IE) Message { return NewDeleteBearerRequest(0, 1, xs...) }},
	{"DeleteBearerResponse", func(xs ...*ie.IE) Message { return NewDeleteBearerResponse(0, 1, xs...) }},
	{"DeletePDNConnectionSetRequest", func(xs ...*ie.IE) Message { return NewDeletePDNConnectionSetRequest(0, 1, xs...) }},
	{"DeletePDNConnectionSetResponse", func(xs ...*ie.IE) Message { return NewDeletePDNConnectionSetResponse(0, 1, xs...) }},
	{"DeleteSessionRequest", func(xs ...*ie.IE) Message { return NewDeleteSessionRequest(0, 1, xs...) }},
	{"DeleteSessionResponse", func(xs ...*ie.IE) Message { return NewDeleteSessionResponse(0, 1, xs...) }},
	{"DetachAcknowledge", func(xs ...*ie.IE) Message { return NewDetachAcknowledge(0, 1, xs...) }},
	{"DetachNotification", func(xs ...*ie.IE) Message { return NewDetachNotification(0, 1, xs...) }},
	{"DownlinkDataNotificationAcknowledge", func(xs ...*ie.IE) Message { return NewDownlinkDataNotificationAcknowledge(0, 1, xs...) }},
	{"DownlinkDataNotificationFailureIndication", func(xs ...*ie.IE) Message { return NewDownlinkDataNotificationFailureIndication(0, 1, xs...) }},
	{"DownlinkDataNotification", func(xs ...*ie.IE) Message { return NewDownlinkDataNotification(0, 1, xs...) }},
	{"EchoRequest", func(xs ...*ie.IE) Message { return NewEchoRequest(1, xs...) }},
	{"EchoResponse", func(xs ...*ie.IE) Message { return NewEchoResponse(1, xs...) }},
	{"ModifyAccessBearersRequest", func(xs ...*ie.IE) Message { return NewModifyAccessBearersRequest(0, 1, xs...) }},
	{"ModifyAccessBearersResponse", func(xs ...*ie.IE) Message { return NewModifyAccessBearersResponse(0, 1, xs...) }},
	{"ModifyBearerCommand", func(xs ...*ie.IE) Message { return NewModifyBearerCommand(0, 1, xs...) }},
	{"ModifyBearerFailureIndication", func(xs ...*ie.IE) Message { return NewModifyBearerFailureIndication(0, 1, xs...) }},
	{"ModifyBearerRequest", func(xs ...*ie.IE) Message { return NewModifyBearerRequest(0, 1, xs...) }},
	{"ModifyBearerResponse", func(xs ...*ie.IE) Message { return NewModifyBearerResponse(0, 1, xs...) }},
	{"PGWRestartNotificationAcknowledge", func(xs ...*ie.IE) Message { return NewPGWRestartNotificationAcknowledge(0, 1, xs...) }},
	{"PGWRestartNotification", func(xs ...*ie.IE) Message { return NewPGWRestartNotification(0, 1, xs...) }},
	{"ReleaseAccessBearersRequest", func(xs ...*ie.IE) Message { return NewReleaseAccessBearersRequest(0, 1, xs...) }},
	{"ReleaseAccessBearersResponse", func(xs ...*ie.IE) Message { return NewReleaseAccessBearersResponse(0, 1, xs...) }},
	{"ResumeAcknowledge", func(xs ...*ie.IE) Message { return NewResumeAcknowledge(0, 1, xs...) }},
	{"ResumeNotification", func(xs ...*ie.IE) Message { return NewResumeNotification(0, 1, xs...) }},
	{"StopPagingIndication", func(xs ...*ie.IE) Message { return NewStopPagingIndication(0, 1, xs...) }},
	{"SuspendAcknowledge", func(xs ...*ie.IE) Message { return NewSuspendAcknowledge(0, 1, xs...) }},
	{"SuspendNotification", func(xs ...*ie.IE) Message { return NewSuspendNotification(0, 1, xs...) }},
	{"UpdateBearerRequest", func(xs ...*ie.IE) Message { return NewUpdateBearerRequest(0, 1, xs...) }},
	{"UpdateBearerResponse", func(xs ...*ie.IE) Message { return NewUpdateBearerResponse(0, 1, xs...) }},
	{"UpdatePDNConnectionSetRequest", func(xs ...*ie.IE) Message { return NewUpdatePDNConnectionSetRequest(0, 1, xs...) }},
	{"UpdatePDNConnectionSetResponse", func(xs ...*ie.IE) Message { return NewUpdatePDNConnectionSetResponse(0, 1, xs...) }},
	{"VersionNotSupportedIndication", func(xs ...*ie.IE) Message { return NewVersionNotSupportedIndication(0, 1, xs...) }},
}

func codecIE(typ, ins uint8) *ie.IE {
	i := ie.New(typ, ins, []byte{typ, ins, 1, 2})
	// Supply meaningful values for roles whose mandatory content is validated.
	switch typ {
	case ie.Cause:
		i = ie.NewCause(16, 0, 0, 0, nil)
	case ie.RATType:
		i = ie.NewRATType(6)
	case ie.FullyQualifiedTEID:
		i = ie.NewFullyQualifiedTEID(0, 1, "192.0.2.1", "")
	case ie.AccessPointName:
		i = ie.NewAccessPointName("internet")
	case ie.AggregateMaximumBitRate:
		i = ie.NewAggregateMaximumBitRate(1, 2)
	case ie.IMSI:
		i = ie.NewIMSI("001010123456789")
	}
	i.SetInstance(ins)
	if i.IsGrouped() {
		return ie.NewGroupedIE(typ, ie.NewEPSBearerID(5)).WithInstance(ins)
	}
	return i
}
func codecRoleIE(typ, ins uint8, field string) *ie.IE {
	if typ == ie.IntegerNumber {
		width := 2
		switch field {
		case "UEUsageType":
			width = 1
		case "RemainingRunningServiceGapTimer":
			width = 4
		}
		return ie.NewIntegerNumber(5, width).WithInstance(ins)
	}
	return codecIE(typ, ins)
}
func codecWire(t *testing.T, m Message) []byte {
	t.Helper()
	b, err := Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func codecPacket(t *testing.T, typ uint8, xs ...*ie.IE) []byte {
	t.Helper()
	return codecWire(t, NewGeneric(typ, 0, 1, xs...))
}
func codecNoPanic(t *testing.T) {
	t.Helper()
	if p := recover(); p != nil {
		t.Errorf("unexpected panic: %v", p)
	}
}
func codecFields(m Message) map[string][]byte {
	v := reflect.ValueOf(m).Elem()
	out := map[string][]byte{}
	for j := 1; j < v.NumField(); j++ {
		f := v.Field(j)
		name := v.Type().Field(j).Name
		switch x := f.Interface().(type) {
		case *ie.IE:
			if x != nil {
				out[name], _ = x.Marshal()
			}
		case []*ie.IE:
			for _, i := range x {
				if i != nil {
					b, _ := i.Marshal()
					out[name] = append(out[name], b...)
				}
			}
		}
	}
	return out
}

// Every named field reachable through a constructor is populated independently,
// then all fields are marshalled together and compared after Parse.
func TestAllFieldsRoundTrip(t *testing.T) {
	for _, f := range codecFactories {
		t.Run(f.name, func(t *testing.T) {
			defer codecNoPanic(t)
			chosen := map[string]*ie.IE{}
			for typ := 1; typ < 256; typ++ {
				for ins := 0; ins < 16; ins++ {
					i := codecIE(uint8(typ), uint8(ins))
					m := f.make(i)
					for field := range codecFields(m) {
						if _, ok := chosen[field]; !ok {
							chosen[field] = codecRoleIE(i.Type, i.Instance(), field)
						}
					}
				}
			}
			m := f.make()
			v := reflect.ValueOf(m).Elem()
			for field, i := range chosen {
				dst := v.FieldByName(field)
				if dst.Type() == reflect.TypeOf(i) {
					dst.Set(reflect.ValueOf(i))
				} else {
					dst.Set(reflect.ValueOf([]*ie.IE{i}))
				}
			}
			b := codecWire(t, m)
			got, err := Parse(b)
			if err != nil {
				t.Fatal(err)
			}
			wantFields, gotFields := codecFields(m), codecFields(got)
			for name, want := range wantFields {
				if !bytes.Equal(want, gotFields[name]) {
					t.Errorf("field %s lost or misrouted: want %x got %x", name, want, gotFields[name])
				}
			}
		})
	}
}
func TestHeaderFraming(t *testing.T) {
	for _, tc := range []struct {
		name     string
		declared uint16
		size     int
	}{{"truncated", 100, 12}, {"trailing", 8, 16}} {
		t.Run(tc.name, func(t *testing.T) {
			b := make([]byte, tc.size)
			b[0] = 0x48
			b[1] = 33
			binary.BigEndian.PutUint16(b[2:4], tc.declared)
			if _, err := Parse(b); err == nil {
				t.Errorf("accepted Length=%d in %d-byte datagram", tc.declared, len(b))
			}
		})
	}
	t.Run("maximum", func(t *testing.T) {
		b := make([]byte, 65539)
		b[0] = 0x48
		binary.BigEndian.PutUint16(b[2:4], 65535)
		if _, err := ParseHeader(b); err != nil {
			t.Errorf("valid maximum header length rejected: %v", err)
		}
	})
}
func TestEightByteHeader(t *testing.T) {
	b := []byte{0x40, 3, 0, 4, 0, 0, 1, 0}
	if _, err := Parse(b); err != nil {
		t.Errorf("valid eight-byte Version Not Supported header rejected: %v", err)
	}
}
func TestVersion(t *testing.T) {
	for _, v := range []byte{0, 1, 3, 7} {
		t.Run(fmt.Sprint(v), func(t *testing.T) {
			b := codecPacket(t, 32)
			b[0] = (b[0] & 0x1f) | (v << 5)
			m, err := Parse(b)
			if err == nil && m.Version() == 2 {
				t.Errorf("wire version %d accepted and reported as version 2", v)
			}
		})
	}
}
func TestVersionNotSupportedHeader(t *testing.T) {
	b := codecWire(t, NewVersionNotSupportedIndication(0, 1))
	if len(b) != 8 || b[0]&8 != 0 {
		t.Errorf("Version Not Supported encoded T=%d length=%d; want T=0 length=8", (b[0]>>3)&1, len(b))
	}
}
func TestMessagePriority(t *testing.T) {
	h := NewHeader(0x48, 32, 0, 1, nil)
	if err := h.SetMessagePriority(9); err != nil {
		t.Fatal(err)
	}
	b, err := h.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if b[11] != 0x90 || h.MessagePriority() != 9 {
		t.Errorf("priority=9 encoded octet=%#02x accessor=%d", b[11], h.MessagePriority())
	}
}
func TestPiggyback(t *testing.T) {
	a := codecWire(t, NewCreateSessionResponse(0, 1))
	b := codecWire(t, NewCreateBearerRequest(0, 2))
	a[0] |= 0x10
	input := append(a, b...)
	m, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if got := codecWire(t, m); !bytes.Equal(got, input) {
		t.Errorf("piggybacked Create Bearer Request lost: input=%d remarshal=%d bytes", len(input), len(got))
	}
}
func TestShortMarshalBuffer(t *testing.T) {
	for _, f := range codecFactories {
		t.Run(f.name, func(t *testing.T) {
			defer codecNoPanic(t)
			m := f.make()
			if err := m.MarshalTo(make([]byte, m.MarshalLen()-1)); err == nil {
				t.Error("short MarshalTo buffer accepted")
			}
		})
	}
}
func TestLengthOverflow(t *testing.T) {
	m := NewGeneric(254, 0, 1, ie.New(250, 0, make([]byte, 65528)))
	if b, err := Marshal(m); err == nil {
		t.Errorf("oversized message marshalled: %d bytes with Length=%d", len(b), binary.BigEndian.Uint16(b[2:4]))
	}
}
func TestModifyAccessSenderFTEID(t *testing.T) {
	i := ie.NewFullyQualifiedTEID(10, 123, "192.0.2.1", "")
	m, err := Parse(codecPacket(t, 211, i))
	if err != nil {
		t.Fatal(err)
	}
	if m.(*ModifyAccessBearersRequest).SenderFTEIDC == nil {
		t.Error("F-TEID instance 0 not assigned to SenderFTEIDC")
	}
}
func TestModifyBearerAddressMapping(t *testing.T) {
	for _, tc := range []struct {
		typ, ins uint8
		field    string
	}{{ie.IPAddress, 0, "HeNBLocalIPAddress"}, {ie.IPAddress, 1, "UELocalIPAddress"}, {ie.IPAddress, 2, "MMESGSNIdentifier"}, {ie.PortNumber, 0, "HeNBUDPPort"}, {ie.PortNumber, 1, "UEUDPPort"}} {
		t.Run(tc.field, func(t *testing.T) {
			i := codecIE(tc.typ, tc.ins)
			m, err := Parse(codecPacket(t, 34, i))
			if err != nil {
				t.Fatal(err)
			}
			if reflect.ValueOf(m).Elem().FieldByName(tc.field).IsNil() {
				t.Errorf("type=%d instance=%d not mapped to %s", tc.typ, tc.ins, tc.field)
			}
		})
	}
}
func TestCreateSessionLDN(t *testing.T) {
	m, err := Parse(codecPacket(t, 33, ie.NewLocalDistinguishedName("sgw").WithInstance(0)))
	if err != nil {
		t.Fatal(err)
	}
	if m.(*CreateSessionResponse).SGWLDN == nil {
		t.Error("SGW LDN instance 0 mapped to PGWLDN")
	}
}
func TestModifyBearerControlMapping(t *testing.T) {
	for _, tc := range []struct {
		typ, ins uint8
		field    string
	}{{ie.LoadControlInformation, 2, "SGWNodeLoadControlInformation"}, {ie.OverloadControlInformation, 0, "PGWOverloadControlInformation"}, {ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"}} {
		t.Run(tc.field, func(t *testing.T) {
			m, err := Parse(codecPacket(t, 35, codecIE(tc.typ, tc.ins)))
			if err != nil {
				t.Fatal(err)
			}
			if reflect.ValueOf(m).Elem().FieldByName(tc.field).IsNil() {
				t.Errorf("type=%d instance=%d not mapped to %s", tc.typ, tc.ins, tc.field)
			}
		})
	}
}
func TestDeleteSessionControlMapping(t *testing.T) {
	for _, tc := range []struct {
		typ, ins uint8
		field    string
	}{{ie.LoadControlInformation, 0, "PGWNodeLoadControlInformation"}, {ie.LoadControlInformation, 1, "PGWAPNLoadControlInformation"}, {ie.LoadControlInformation, 2, "SGWNodeLoadControlInformation"}, {ie.OverloadControlInformation, 0, "PGWOverloadControlInformation"}, {ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"}} {
		t.Run(tc.field, func(t *testing.T) {
			i := codecIE(tc.typ, tc.ins)
			m, err := Parse(codecPacket(t, 37, i))
			if err != nil {
				t.Fatal(err)
			}
			got := codecFields(m)[tc.field]
			want, _ := i.Marshal()
			if !bytes.Equal(want, got) {
				t.Errorf("type=%d instance=%d not mapped to %s", tc.typ, tc.ins, tc.field)
			}
		})
	}
}
func TestContextOmittedFields(t *testing.T) {
	for _, i := range []*ie.IE{codecIE(ie.MonitoringEventExtensionInformation, 0), codecIE(ie.AdditionalRRMPolicyIndex, 0), codecIE(ie.AdditionalRRMPolicyIndex, 1)} {
		t.Run(fmt.Sprint(i.Type, i.Instance()), func(t *testing.T) {
			m := NewContextResponse(0, 1, i)
			b := codecWire(t, m)
			want, err := i.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if len(b) != 12+i.MarshalLen() || !bytes.Equal(b[12:], want) {
				t.Errorf("Context Response silently drops IE type=%d instance=%d", i.Type, i.Instance())
			}
		})
	}
}
func TestContextMM(t *testing.T) {
	i := codecIE(ie.MMContextEPSSecurityContextQuadrupletsAndQuintuplets, 0)
	m, err := Parse(codecPacket(t, 131, i))
	if err != nil {
		t.Fatal(err)
	}
	if m.(*ContextResponse).UEMMContext == nil {
		t.Error("UE MM Context is nil despite MM Context IE")
	}
}
func TestRepeatedIEs(t *testing.T) {
	for _, tc := range []struct{ typ, iet, ins uint8 }{{33, ie.BearerContext, 1}, {36, ie.SecondaryRATUsageDataReport, 0}, {131, ie.MonitoringEventInformation, 0}, {176, ie.EPSBearerID, 0}, {99, ie.BearerContext, 0}} {
		t.Run(fmt.Sprint(tc.typ, tc.iet), func(t *testing.T) {
			a, b := codecIE(tc.iet, tc.ins), codecIE(tc.iet, tc.ins)
			input := codecPacket(t, tc.typ, a, b)
			m, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			if got := codecWire(t, m); len(got) != len(input) {
				t.Errorf("repeated IE overwritten: input=%d remarshal=%d bytes", len(input), len(got))
			}
		})
	}
}
func TestUnknownInstances(t *testing.T) {
	for _, f := range codecFactories {
		if f.name == "VersionNotSupportedIndication" || f.name == "Generic" {
			continue
		}
		t.Run(f.name, func(t *testing.T) {
			var bad []string
			for typ := 1; typ < 255; typ++ {
				if typ == int(ie.PrivateExtension) {
					continue
				}
				i := codecIE(uint8(typ), 15)
				m := f.make(i)
				fields := codecFields(m)
				for name := range fields {
					if name != "AdditionalIEs" {
						bad = append(bad, fmt.Sprintf("%s(type=%d)", name, typ))
					}
				}
				if len(fields) == 0 {
					bad = append(bad, fmt.Sprintf("dropped(type=%d)", typ))
				}
			}
			if len(bad) > 0 {
				t.Errorf("instance 15 not retained as AdditionalIEs: %v", bad)
			}
		})
	}
}
func TestUnknownSwitchDrops(t *testing.T) {
	for _, tc := range []struct{ typ, iet uint8 }{{96, ie.FullyQualifiedCSID}, {96, ie.TWANIdentifier}, {96, ie.TWANIdentifierTimestamp}, {96, ie.OverloadControlInformation}, {96, ie.PortNumber}, {96, ie.IPAddress}, {98, ie.IPAddress}, {64, ie.OverloadControlInformation}, {65, ie.OverloadControlInformation}, {66, ie.OverloadControlInformation}, {67, ie.OverloadControlInformation}, {211, ie.BearerContext}, {212, ie.BearerContext}, {164, 250}, {165, 250}} {
		t.Run(fmt.Sprint(tc.typ, tc.iet), func(t *testing.T) {
			input := codecPacket(t, tc.typ, codecIE(tc.iet, 15))
			m, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			if len(codecWire(t, m)) != len(input) {
				t.Error("unexpected IE silently lost instead of retained")
			}
		})
	}
}
func TestPresenceReportingInformation(t *testing.T) {
	for _, typ := range []uint8{96, 98} {
		t.Run(fmt.Sprint(typ), func(t *testing.T) {
			i := codecIE(ie.PresenceReportingAreaInformation, 0)
			m, err := Parse(codecPacket(t, typ, i))
			if err != nil {
				t.Fatal(err)
			}
			for field := range codecFields(m) {
				if field != "AdditionalIEs" {
					return
				}
			}
			t.Error("Presence Reporting Area Information not assigned to named presence field")
		})
	}
}
func TestDecodeReuse(t *testing.T) {
	for _, f := range codecFactories {
		if f.name == "VersionNotSupportedIndication" || f.name == "ResumeAcknowledge" || f.name == "ResumeNotification" {
			continue
		}
		t.Run(f.name, func(t *testing.T) {
			m := f.make(codecIE(250, 0))
			b := codecWire(t, f.make())
			if len(b) == 8 {
				b = codecWire(t, f.make(ie.NewRecovery(1)))
			}
			if err := m.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
			if len(codecFields(m)["AdditionalIEs"]) > 0 {
				t.Error("UnmarshalBinary retained stale IE from preceding message")
			}
		})
	}
}
func TestTrailingIEByte(t *testing.T) {
	for _, f := range codecFactories {
		if f.name == "VersionNotSupportedIndication" {
			continue
		}
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			v := reflect.ValueOf(m).Elem().FieldByName("Header").Interface().(*Header)
			v.Payload = []byte{1}
			v.SetLength()
			b, err := v.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			err = m.UnmarshalBinary(b)
			if f.name == "EchoRequest" {
				// TS 29.274 §7.1.2 NOTE: malformed IEs do not prevent Echo Response.
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil {
				t.Error("incomplete one-byte IE accepted")
			}
		})
	}
}
func TestNilSliceEntries(t *testing.T) {
	for _, f := range codecFactories {
		m := f.make()
		v := reflect.ValueOf(m).Elem()
		for j := 1; j < v.NumField(); j++ {
			field := v.Field(j)
			if field.Type() != reflect.TypeOf([]*ie.IE{}) {
				continue
			}
			t.Run(f.name+"/"+v.Type().Field(j).Name, func(t *testing.T) {
				defer codecNoPanic(t)
				fresh := f.make()
				reflect.ValueOf(fresh).Elem().Field(j).Set(reflect.ValueOf([]*ie.IE{nil}))
				if _, err := Marshal(fresh); err != nil {
					t.Fatalf("nil IE must be skipped: %v", err)
				}
			})
		}
	}
}

func TestPrivateExtensionOverwrite(t *testing.T) {
	a := ie.NewPrivateExtension(123, []byte{1})
	b := ie.New(250, 0, []byte{2})
	m := NewCreateSessionRequest(0, 1, a, b)
	raw := codecWire(t, m)
	if raw[12] != ie.PrivateExtension {
		t.Errorf("AdditionalIEs overwrote PrivateExtension at offset 12: got type %d", raw[12])
		return
	}
	got, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.(*CreateSessionRequest).PrivateExtension == nil {
		t.Error("AdditionalIEs overwrote serialized PrivateExtension")
	}
}
func TestSequenceOverflow(t *testing.T) {
	h := NewHeader(0x48, 32, 0, 0x1000001, nil)
	if b, err := h.Marshal(); err == nil {
		t.Errorf("sequence %#x silently encoded as %#x", h.Sequence(), uint32(b[8])<<16|uint32(b[9])<<8|uint32(b[10]))
	}
}
func TestHeaderTEIDPolicy(t *testing.T) {
	for _, tc := range []struct {
		typ   uint8
		flags uint8
	}{{1, 0x48}, {2, 0x48}, {3, 0x48}, {32, 0x40}} {
		t.Run(fmt.Sprint(tc.typ), func(t *testing.T) {
			h := NewHeader(tc.flags, tc.typ, 0, 1, []byte{250, 0, 0, 0})
			if _, err := h.Marshal(); !errors.Is(err, ErrInvalidTEIDFlag) {
				t.Fatalf("invalid T accepted by encoder: %v", err)
			}
			b := make([]byte, h.headerLen())
			b[0], b[1] = tc.flags, tc.typ
			binary.BigEndian.PutUint16(b[2:4], uint16(len(b)-4))
			if _, err := Parse(b); !errors.Is(err, ErrInvalidTEIDFlag) {
				t.Errorf("message type %d accepted invalid T=%d: %v", tc.typ, (tc.flags>>3)&1, err)
			}
		})
	}
}

func TestTypedDecoderMessageType(t *testing.T) {
	for _, f := range codecFactories {
		if f.name == "Generic" {
			continue
		}
		t.Run(f.name, func(t *testing.T) {
			m := f.make()
			b := codecPacket(t, 254, ie.NewRecovery(1))
			if err := m.UnmarshalBinary(b); err == nil {
				t.Errorf("%s decoder accepted message type 254", f.name)
			}
		})
	}
}
func TestDuplicateSingleton(t *testing.T) {
	for _, typ := range []uint8{33, 35, 37, 96, 98, 100, 131, 132, 177, 163, 165} {
		t.Run(fmt.Sprint(typ), func(t *testing.T) {
			a := ie.NewCause(16, 0, 0, 0, nil)
			b := ie.NewCause(64, 0, 0, 0, nil)
			m, err := Parse(codecPacket(t, typ, a, b))
			if err != nil {
				t.Fatal(err)
			}
			got := reflect.ValueOf(m).Elem().FieldByName("Cause").Interface().(*ie.IE)
			if got == nil || len(got.Payload) == 0 || got.Payload[0] != 16 {
				t.Error("later duplicate Cause overwrote the first occurrence")
			}
		})
	}
}

func TestGenericNil(t *testing.T) {
	t.Run("constructor", func(t *testing.T) { defer codecNoPanic(t); _ = NewGeneric(254, 0, 1, nil) })
	t.Run("AddIE", func(t *testing.T) { defer codecNoPanic(t); NewGeneric(254, 0, 1).AddIE(nil) })
}

func TestLegalAdditionalRoles(t *testing.T) {
	charging := ie.New(ie.FullyQualifiedDomainName, 0, []byte{2, 'c', 'g'})
	node := ie.New(ie.FullyQualifiedDomainName, 1, []byte{2, 'p', 'g'})
	m, err := Parse(codecPacket(t, 33, charging, node))
	if err != nil {
		t.Fatal(err)
	}
	got := m.(*CreateSessionResponse)
	if !bytes.Equal(got.ChargingGatewayName.Payload, charging.Payload) {
		t.Error("legal PGW node name (FQDN instance 1) overwrote Charging Gateway Name (instance 0)")
	}
	if len(got.AdditionalIEs) != 1 {
		t.Error("legal PGW node name not retained in AdditionalIEs")
	}
}
