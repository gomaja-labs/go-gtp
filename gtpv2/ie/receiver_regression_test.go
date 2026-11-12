package ie

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestS1UForwardingTEID(t *testing.T) {
	// TS 29.274 V20.0.0 §8.26 Figure 8.26-1: EBI, address length,
	// address, then TEID. This vector is independent of the encoder.
	i := New(S1UDataForwarding, 0, []byte{5, 4, 192, 0, 2, 1, 0x12, 0x34, 0x56, 0x78})
	got, err := i.TEID()
	if err != nil || got != 0x12345678 {
		t.Fatalf("TEID = %#x, %v; want 0x12345678", got, err)
	}
}

func TestScalarGroupedDataLoss(t *testing.T) {
	for _, i := range []*IE{NewUint32IE(BearerContext, 0x03000000), NewUint64IE(BearerContext, 0x0300000003000000)} {
		if i.Err() != nil {
			continue
		}
		b, err := i.Marshal()
		if err != nil {
			continue
		}
		if !bytes.Equal(b[4:], i.Payload) {
			t.Errorf("scalar constructor payload %x marshals as %x", i.Payload, b[4:])
		}
	}
}

func TestInvalidCSIDHex(t *testing.T) {
	i := NewFullyQualifiedCSID("01020304zz", 1)
	if i.Err() == nil {
		b, err := i.Marshal()
		t.Fatalf("invalid hex silently accepted: wire=%x error=%v", b, err)
	}
}

func TestPagingFieldOverflow(t *testing.T) {
	// TS 29.274 V20.0.0 §8.117 Figure 8.117-1: EBI has four bits.
	f := &PagingAndServiceInformationFields{EPSBearerID: 0x15, Flags: 1, PagingPolicyIndication: 1}
	if b, err := f.Marshal(); err == nil {
		t.Fatalf("out-of-range EBI 21 silently encoded as %d: %x", b[0], b)
	}
}

func TestTFTDeleteSpareBits(t *testing.T) {
	// TS 24.008 V20.1.0 §10.5 and §10.5.6.12 Figure 10.5.144a:
	// delete one packet filter, ID 5, with the spare high nibble set.
	f, err := ParseTrafficFlowTemplate([]byte{0xa1, 0xf5})
	if err != nil {
		t.Fatalf("spare bits reject TFT delete: %v", err)
	}
	if len(f.PacketFilterIdentifiers) != 1 || f.PacketFilterIdentifiers[0] != 5 {
		t.Fatalf("spare bits leaked into identifiers: %v", f.PacketFilterIdentifiers)
	}
}

func TestGroupedAccessorsIgnoreUnknownInstances(t *testing.T) {
	// TS 29.274 V20.0.0 §7.7.9, Table 7.2.4-2: Cause and EBI use instance 0.
	i := NewBearerContext(NewCause(64, 0, 0, 0, nil).WithInstance(1),
		NewCause(16, 0, 0, 0, nil), NewEPSBearerID(6).WithInstance(1), NewEPSBearerID(5))
	cause, err := i.Cause()
	if err != nil || cause != 16 {
		t.Errorf("Cause = %d, %v; want instance-0 value 16", cause, err)
	}
	ebi, err := i.EPSBearerID()
	if err != nil || ebi != 5 {
		t.Errorf("EBI = %d, %v; want instance-0 value 5", ebi, err)
	}
}

func TestErrorIEFlagHelpers(t *testing.T) {
	// The package contract says Has helpers return zero values on errors.
	i := NewFullyQualifiedTEID(0, 1, "192.0.2.1", "2001:db8::1").WithInstance(16)
	if i.Err() == nil {
		t.Fatal("fixture must carry an error")
	}
	if i.HasIPv4() || i.HasIPv6() {
		t.Error("HasIPv4/HasIPv6 report true on an error-carrying IE")
	}
	c := NewCause(16, 1, 1, 1, nil).WithInstance(16)
	if c.IsRemoteCause() || c.IsBearerContextIEError() || c.IsPDNConnectionIEError() {
		t.Error("Cause boolean helpers ignore the stored construction error")
	}
}

func TestAPNUsesAPNSyntax(t *testing.T) {
	// TS 23.003 V20.1.0 §9.1 explicitly restricts APN labels and
	// prohibits the terminating zero-length root label.
	if b, err := NewAccessPointName("bad_apn").Marshal(); err == nil {
		t.Errorf("invalid APN label emitted as %x", b)
	}
	i := New(AccessPointName, 0, []byte{1, 'a', 0})
	if s, err := i.AccessPointName(); err == nil {
		t.Errorf("root-terminated APN accepted as %q", s)
	}
}

func TestPCOLengthLimit(t *testing.T) {
	// TS 29.274 V20.0.0 §8.13 Figure 8.13-1 and §8.94 Figure 8.94-1
	// carry NAS octets 3-z. TS 24.008 V20.1.0 §10.5.6.3 Figure 10.5.136
	// limits the complete NAS IE to 253 octets, leaving 251 value octets.
	for _, i := range []*IE{
		NewProtocolConfigurationOptions(0, NewPCOContainer(0xff00, make([]byte, 248))),
		NewAdditionalProtocolConfigurationOptions(0, NewPCOContainer(0xff00, make([]byte, 248))),
	} {
		if b, err := i.Marshal(); err == nil {
			t.Errorf("IE type %d accepted %d PCO value octets; maximum is 251", i.Type, len(b)-4)
		}
	}
}

func TestScalarGroupedConstructor(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("NewUint16IE panicked instead of returning an error IE: %v", p)
		}
	}()
	i := NewUint16IE(BearerContext, 1)
	if i.Err() == nil {
		t.Fatal("scalar constructor accepted a grouped IE type")
	}
}

func TestNilGroupedChildren(t *testing.T) {
	for _, name := range []string{"Cause", "CauseFlags", "BearerFlags", "BearerQoS", "ChargingID", "EPSBearerID", "RANNASCause", "ProtocolConfigurationOptions", "BearerTFT"} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if p := recover(); p != nil {
					t.Fatalf("nil child caused panic: %v", p)
				}
			}()
			i := &IE{Type: BearerContext, ChildIEs: []*IE{nil}}
			var err error
			switch name {
			case "Cause":
				_, err = i.Cause()
			case "CauseFlags":
				_, err = i.CauseFlags()
			case "BearerFlags":
				_, err = i.BearerFlags()
			case "BearerQoS":
				_, err = i.BearerQoS()
			case "ChargingID":
				_, err = i.ChargingID()
			case "EPSBearerID":
				_, err = i.EPSBearerID()
			case "RANNASCause":
				_, err = i.RANNASCause()
			case "ProtocolConfigurationOptions":
				_, err = i.ProtocolConfigurationOptions()
			case "BearerTFT":
				_, err = i.BearerTFT()
			}
			if !errors.Is(err, ErrMalformed) {
				t.Fatalf("nil child error: %v", err)
			}
		})
	}
}

func TestScalarConstructorsRejectGroupedTypes(t *testing.T) {
	for _, typ := range []uint8{BearerContext, PDNConnection, OverloadControlInformation, LoadControlInformation, RemoteUEContext, SCEFPDNConnection, V2XContext, PC5QoSParameters} {
		for _, i := range []*IE{NewUint8IE(typ, 1), NewUint16IE(typ, 1), NewUint32IE(typ, 1), NewUint64IE(typ, 1)} {
			if !errors.Is(i.Err(), ErrInvalidType) {
				t.Fatalf("grouped scalar type %d: %v", typ, i.Err())
			}
			if _, err := i.Marshal(); !errors.Is(err, ErrInvalidType) {
				t.Fatalf("grouped scalar marshaled: %v", err)
			}
		}
	}
	for _, tc := range []struct {
		i    *IE
		want []byte
	}{
		{NewUint8IE(250, 0x12), []byte{0x12}},
		{NewUint16IE(250, 0x1234), []byte{0x12, 0x34}},
		{NewUint32IE(250, 0x12345678), []byte{0x12, 0x34, 0x56, 0x78}},
		{NewUint64IE(250, 0x123456789abcdef0), []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}},
	} {
		b, err := tc.i.Marshal()
		if err != nil || !bytes.Equal(b[4:], tc.want) || len(b) != tc.i.MarshalLen() {
			t.Fatalf("scalar bytes = %x, %v; want %x", b, err, tc.want)
		}
	}
}

func TestS1UForwardingTEIDAddressLengths(t *testing.T) {
	for _, address := range [][]byte{{192, 0, 2, 1}, {0x20, 1, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}} {
		payload := append([]byte{0xf5, byte(len(address))}, address...)
		payload = append(payload, 0x12, 0x34, 0x56, 0x78)
		for n := 0; n < len(payload); n++ {
			if _, err := New(S1UDataForwarding, 0, payload[:n]).TEID(); err == nil {
				t.Fatalf("truncated %d-byte address payload of length %d accepted", len(address), n)
			}
		}
		if got, err := New(S1UDataForwarding, 0, payload).TEID(); err != nil || got != 0x12345678 {
			t.Fatalf("TEID = %#x, %v", got, err)
		}
	}
}

func TestFQCSIDFieldsPropagateHexError(t *testing.T) {
	for _, name := range []string{"01020304zz", "01020304f"} {
		f := NewFullyQualifiedCSIDFields(name, 1)
		if err := f.MarshalTo(make([]byte, f.MarshalLen())); err == nil {
			t.Fatalf("bad hex %q accepted by fields marshaler", name)
		}
	}
	for _, name := range []string{"192.0.2.1", "2001:db8::1", "01020304"} {
		i := NewFullyQualifiedCSID(name, 1)
		b, err := i.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		f, err := ParseFullyQualifiedCSIDFields(b[4:])
		if err != nil || len(f.CSIDs) != 1 || f.CSIDs[0] != 1 {
			t.Fatalf("valid node %q: %v", name, err)
		}
	}
}

func TestPagingFieldsValidateBeforeWriting(t *testing.T) {
	for _, f := range []*PagingAndServiceInformationFields{
		{EPSBearerID: 16}, {EPSBearerID: 5, Flags: 2}, {EPSBearerID: 5, Flags: 1, PagingPolicyIndication: 128},
	} {
		b := []byte{0xaa, 0xaa, 0xaa}
		if err := f.MarshalTo(b); !errors.Is(err, ErrMalformed) || !bytes.Equal(b, []byte{0xaa, 0xaa, 0xaa}) {
			t.Fatalf("out-of-range fields: %x, %v", b, err)
		}
	}
	f := &PagingAndServiceInformationFields{EPSBearerID: 15, Flags: 1, PagingPolicyIndication: 127}
	for n := 0; n < f.MarshalLen(); n++ {
		if err := f.MarshalTo(make([]byte, n)); !errors.Is(err, ErrTooShortToMarshal) {
			t.Fatalf("short destination %d: %v", n, err)
		}
	}
	if b, err := f.Marshal(); err != nil || !bytes.Equal(b, []byte{15, 1, 127}) || len(b) != f.MarshalLen() {
		t.Fatalf("valid fields = %x, %v", b, err)
	}
}

func TestTFTDeleteMasksWithoutAliasingInput(t *testing.T) {
	input := []byte{0xa2, 0xf5, 0x86}
	original := append([]byte(nil), input...)
	f, err := ParseTrafficFlowTemplate(input)
	if err != nil || !bytes.Equal(input, original) || !bytes.Equal(f.PacketFilterIdentifiers, []byte{5, 6}) {
		t.Fatalf("delete IDs = %v; input=%x; error=%v", f, input, err)
	}
	f.PacketFilterIdentifiers[0] = 7
	if !bytes.Equal(input, original) {
		t.Fatal("decoded IDs alias input")
	}
	out, err := f.Marshal()
	if err != nil || !bytes.Equal(out, []byte{0xa2, 7, 6}) {
		t.Fatalf("spare bits sent: %x, %v", out, err)
	}
}

func TestPCOValueLengthBoundaries(t *testing.T) {
	// TS 24.008 §10.5.6.3: 253 NAS octets minus IEI and Length = 251.
	for _, size := range []int{250, 251, 252, 253, 256} {
		options := NewPCOContainer(0xff00, make([]byte, size-4))
		for _, i := range []*IE{NewProtocolConfigurationOptions(0, options), NewAdditionalProtocolConfigurationOptions(0, options)} {
			b, err := i.Marshal()
			if size <= 251 {
				if err != nil || len(b) != 4+size {
					t.Fatalf("valid %d-octet PCO: %v", size, err)
				}
			} else if !errors.Is(err, ErrInvalidLength) {
				t.Fatalf("%d-octet PCO error: %v", size, err)
			}
		}
		f := NewProtocolConfigurationOptionsFields(0, options)
		b, err := f.Marshal()
		if (size <= 251) != (err == nil) {
			t.Fatalf("fields size %d: %v", size, err)
		}
		if size > 251 {
			// Form a length-consistent raw received value independently of the encoder.
			b = append([]byte{0x80, 0xff, 0, byte(size - 4)}, make([]byte, size-4)...)
		}
		_, err = New(ProtocolConfigurationOptions, 0, b).ProtocolConfigurationOptions()
		if (size <= 251) != (err == nil) {
			t.Fatalf("PCO accessor size %d: %v", size, err)
		}
		_, err = New(AdditionalProtocolConfigurationOptions, 0, b).AdditionalProtocolConfigurationOptions()
		if (size <= 251) != (err == nil) {
			t.Fatalf("APCO accessor size %d: %v", size, err)
		}
	}
}

func TestAPNLabelValidation(t *testing.T) {
	for _, name := range []string{"", "a.", ".a", "a..b", "-a", "a-", "bad_apn", "a b", "é", strings.Repeat("a", 64)} {
		if _, err := NewAccessPointName(name).Marshal(); err == nil {
			t.Errorf("bad APN %q emitted", name)
		}
	}
	for _, b := range [][]byte{nil, {0}, {1, 'a', 0}, {2, 'a'}, {0xc0, 0}, {1, '_'}, {2, '-', 'a'}} {
		if _, err := New(AccessPointName, 0, b).AccessPointName(); err == nil {
			t.Errorf("bad APN %x decoded", b)
		}
	}
	for _, name := range []string{"internet", "1a-b.APN", strings.Repeat("a", 63), strings.Repeat("a", 63) + "." + strings.Repeat("b", 35)} {
		i := NewAccessPointName(name)
		b, err := i.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := Parse(b)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := decoded.AccessPointName(); err != nil || got != name {
			t.Fatalf("APN = %q, %v; want %q", got, err, name)
		}
	}
}

func TestIntegerNumberExplicitWidth(t *testing.T) {
	for _, width := range []int{1, 2, 3, 4} {
		i := NewIntegerNumber(5, width)
		b, err := i.Marshal()
		if err != nil || len(b) != 4+width || b[len(b)-1] != 5 || i.MustIntegerNumber() != 5 {
			t.Fatalf("width %d = %x, %v", width, b, err)
		}
		if width < 4 {
			if _, err := NewIntegerNumber(1<<uint(width*8), width).Marshal(); !errors.Is(err, ErrMalformed) {
				t.Fatalf("overflow in %d octets: %v", width, err)
			}
		}
	}
	for _, width := range []int{-1, 0, 5} {
		if _, err := NewIntegerNumber(0, width).Marshal(); !errors.Is(err, ErrInvalidLength) {
			t.Fatalf("bad width %d: %v", width, err)
		}
	}
}

func TestGroupedAccessorsUseInstanceZero(t *testing.T) {
	for _, tc := range []struct {
		name  string
		child *IE
	}{
		{"Cause", NewCause(16, 0, 0, 1, nil)},
		{"CauseFlags", NewCause(16, 1, 0, 1, nil)},
		{"BearerFlags", NewBearerFlags(1, 0, 1, 0)},
		{"BearerQoS", NewBearerQoS(1, 2, 1, 9, 1, 2, 3, 4)},
		{"ChargingID", NewChargingID(0x12345678)},
		{"EPSBearerID", NewEPSBearerID(5)},
		{"RANNASCause", NewRANNASCause(1, 0, []byte{1})},
		{"ProtocolConfigurationOptions", NewProtocolConfigurationOptions(0)},
		{"BearerTFT", NewBearerTFTDeleteExistingTFT()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.child.Err(); err != nil {
				t.Fatal(err)
			}
			good := reflect.ValueOf(tc.child).MethodByName(tc.name).Call(nil)
			if !good[1].IsNil() {
				t.Fatal(good[1].Interface())
			}
			group := NewBearerContext(New(tc.child.Type, 1, nil), tc.child)
			got := reflect.ValueOf(group).MethodByName(tc.name).Call(nil)
			if !got[1].IsNil() || !reflect.DeepEqual(got[0].Interface(), good[0].Interface()) {
				t.Fatalf("wrong child selected: %v", got)
			}
		})
	}
}

func TestGroupedDecodeRetainsErrorAndResetsOnReuse(t *testing.T) {
	var i IE
	if err := i.UnmarshalBinary([]byte{BearerContext, 0, 1, 15, 0xff}); err != nil {
		t.Fatal(err)
	}
	if _, err := i.BearerContext(); err == nil {
		t.Fatal("content error lost")
	}
	i.SetLength()
	if !bytes.Equal(i.Payload, []byte{0xff}) || i.Length != 1 || i.Instance() != 15 {
		t.Fatal("raw erroneous content lost")
	}
	if err := i.UnmarshalBinary([]byte{Recovery, 0, 1, 0, 7}); err != nil || i.Err() != nil || len(i.ChildIEs) != 0 || i.MustRecovery() != 7 {
		t.Fatalf("reused IE retained stale grouped content: %v", err)
	}
}
