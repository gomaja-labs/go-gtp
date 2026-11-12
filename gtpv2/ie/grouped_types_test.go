package ie_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

func TestGroupedTypesMatchV20(t *testing.T) {
	// TS 29.274 §§8.28, 8.39, 8.111, 8.112, 8.122, 8.126, 8.139, 8.140, 8.145.
	want := map[uint8]bool{93: true, 109: true, 180: true, 181: true, 191: true, 195: true, 208: true, 209: true, 214: true}
	for n := 0; n < 256; n++ {
		i := ie.New(uint8(n), 0, nil)
		if i.IsGrouped() != want[uint8(n)] {
			t.Errorf("grouped type %d", n)
		}
	}
	names := []struct {
		typ  uint8
		name string
	}{
		{ie.PGWChangeInfo, "PGWChangeInfo"}, {ie.PGWFQDN, "PGWFQDN"}, {ie.GroupID, "GroupID"}, {ie.PSCellID, "PSCellID"}, {ie.UPSecurityPolicy, "UPSecurityPolicy"}, {ie.AlternativeIMSI, "AlternativeIMSI"}, {ie.NFInstanceID, "NFInstanceID"}, {ie.TimerInSeconds, "TimerInSeconds"}, {ie.MDTConfigurationNR, "MDTConfigurationNR"}, {ie.URI, "URI"},
	}
	for j, tc := range names {
		if tc.typ != uint8(214+j) || ie.New(tc.typ, 0, nil).Name() != tc.name {
			t.Errorf("type %d name %s", tc.typ, tc.name)
		}
	}
	i := ie.NewPGWChangeInfo(nil, ie.NewIPAddress("192.0.2.99"))
	wantBytes := []byte{214, 0, 8, 0, 74, 0, 4, 0, 192, 0, 2, 99}
	b, err := i.Marshal()
	if err != nil || !bytes.Equal(b, wantBytes) || len(i.ChildIEs) != 1 {
		t.Fatalf("PGW Change Info: %x %v", b, err)
	}
	if _, err := ie.NewPGWChangeInfo(ie.NewMSISDN("x")).Marshal(); err == nil {
		t.Fatal("invalid child construction accepted")
	}
}

func TestMalformedPGWChangeInfoDoesNotRejectMessage(t *testing.T) {
	b := []byte{0x48, 33, 0, 19, 0, 0, 0, 1, 0, 0, 1, 0, 2, 0, 2, 0, 16, 0, 214, 0, 1, 0, 255}
	m, err := message.ParseCreateSessionResponse(b)
	if err != nil {
		t.Fatal(err)
	}
	if m.Cause.MustCause() != 16 || len(m.AdditionalIEs) != 1 {
		t.Fatal("message lost fields")
	}
	bad := m.AdditionalIEs[0]
	if !bad.IsGrouped() || !errors.Is(bad.Err(), ie.ErrTooShortToParse) {
		t.Fatalf("content error: %v", bad.Err())
	}
}
