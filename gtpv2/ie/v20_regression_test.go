package ie_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

func TestIndicationLateFlagAccessors(t *testing.T) {
	b, _ := hex.DecodeString("4d000a00" + "0000000000000000" + "8000")
	i, err := ie.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	if !i.HasNSPUSI() || i.HasPGWRNSI() || i.HasUPIPSI() {
		t.Fatal("octets 13 and 14 flag positions")
	}

	for _, name := range []string{"NSPUSI", "PGWRNSI", "RPPCSI", "PGWCHI", "SISSME", "NSENBI", "IDFUPF", "UHRRSI", "DISREI", "PPRI", "LAPCOSI", "LTEMSAI", "SRTPI", "UPIPSI", "PMTSMI"} {
		if _, ok := reflect.TypeFor[*ie.IE]().MethodByName("Has" + name); !ok {
			t.Errorf("missing Has%s", name)
		}
	}
}

func TestTypeExtensionAccessor(t *testing.T) {
	b, _ := hex.DecodeString("fe0005000100deadbe")
	i, err := ie.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	m := reflect.ValueOf(i).MethodByName("TypeExtension")
	if !m.IsValid() {
		t.Fatal("missing TypeExtension")
	}
	got := m.Call(nil)
	if got[0].Uint() != 256 || !got[1].IsNil() {
		t.Fatalf("extension: %v", got)
	}
}

func TestV20AdditionalFlagAccessors(t *testing.T) {
	for _, tc := range []struct {
		i      *ie.IE
		method string
	}{{ie.NewNodeFeatures(0x80), "HasPSSET"}, {ie.New(ie.BearerFlags, 0, []byte{0x10}), "HasUHRREI"}} {
		m := reflect.ValueOf(tc.i).MethodByName(tc.method)
		if !m.IsValid() {
			t.Errorf("missing %s", tc.method)
			continue
		}
		if !m.Call(nil)[0].Bool() {
			t.Errorf("%s did not read its spec bit", tc.method)
		}
	}
}

func TestBearerFlagsUHRREIAndNodePSSET(t *testing.T) {
	for bit := 0; bit < 8; bit++ {
		n := ie.NewNodeFeatures(1 << bit)
		accessors := []func() bool{n.HasPRN, n.HasMABR, n.HasNTSR, n.HasCIOT, n.HasS1UN, n.HasETH, n.HasMTEDT, n.HasPSSET}
		for j, has := range accessors {
			if has() != (j == bit) {
				t.Errorf("Node Features bit %d accessor %d", bit, j)
			}
		}
	}
	for bit := 0; bit < 5; bit++ {
		values := make([]uint8, 5)
		values[4-bit] = 1
		b := ie.NewBearerFlags(values[0], values[1], values[2], values[3], values[4])
		if !bytes.Equal(b.Payload, []byte{1 << bit}) {
			t.Fatalf("Bearer Flags bit %d: %x", bit, b.Payload)
		}
		accessors := []func() bool{b.HasPPC, b.HasVB, b.HasVind, b.HasASI, b.HasUHRREI}
		for j, has := range accessors {
			if has() != (j == bit) {
				t.Errorf("Bearer Flags bit %d accessor %d", bit, j)
			}
		}
		values[4-bit] = 2
		if _, err := ie.NewBearerFlags(values[0], values[1], values[2], values[3], values[4]).Marshal(); !errors.Is(err, ie.ErrMalformed) {
			t.Fatalf("invalid bit %d: %v", bit, err)
		}
	}
	for _, i := range []*ie.IE{nil, ie.NewRecovery(255), ie.New(ie.BearerFlags, 0, nil)} {
		if i.HasUHRREI() {
			t.Fatal("invalid UHRREI receiver")
		}
	}
	for _, i := range []*ie.IE{nil, ie.NewRecovery(255), ie.New(ie.NodeFeatures, 0, nil)} {
		if i.HasPSSET() {
			t.Fatal("invalid PSSET receiver")
		}
	}
}

func TestGroupedConstructionRejectsMalformedDescendants(t *testing.T) {
	malformed, err := ie.Parse([]byte{93, 0, 1, 0, 255})
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range []*ie.IE{ie.NewGroupedIE(ie.BearerContext, malformed), ie.New(ie.BearerContext, 0, []byte{93, 0, 1, 0, 255}), ie.New(ie.BearerContext, 0, []byte{93, 0, 5, 0, 93, 0, 1, 0, 255})} {
		if !errors.Is(i.Err(), ie.ErrTooShortToParse) {
			t.Fatalf("constructor did not retain descendant error: %v", i.Err())
		}
		if _, err := i.Marshal(); !errors.Is(err, ie.ErrTooShortToParse) {
			t.Fatalf("constructed malformed descendant accepted: %v", err)
		}
	}
}

func TestSRVCCCauseTypeHasName(t *testing.T) {
	// TS 29.274 Table 8.1-1 explicitly allocates 56, within the Sv range.
	if ie.SRVCCCause != 56 {
		t.Fatalf("SRVCC Cause type: %d", ie.SRVCCCause)
	}
	if name := ie.New(ie.SRVCCCause, 0, nil).Name(); name != "SRVCCCause" {
		t.Fatalf("type 56 name: %q", name)
	}
}
