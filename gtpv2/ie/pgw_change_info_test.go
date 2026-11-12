package ie_test

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// pgwChangeInfoHex is a synthetic PGW Change Info IE (TS 29.274 V20.0.0
// §8.145) holding the children of Table 7.2.2-6: PGW Set FQDN (PGW FQDN,
// instance 0), two Alternative PGW-C/SMF IP Addresses (IP Address, instance
// 0), an Alternative PGW-C/SMF FQDN (PGW FQDN, instance 1) and a Group Id.
const pgwChangeInfoHex = "d6 00 49 00" +
	"d7 00 10 00 07 70 67 77 2d 73 65 74 07 65 78 61 6d 70 6c 65" + // pgw-set.example
	"4a 00 04 00 c0 00 02 63" + // 192.0.2.99
	"4a 00 10 00 20 01 0d b8 00 00 00 00 00 00 00 00 00 00 00 63" + // 2001:db8::63
	"d7 00 0d 01 04 70 67 77 32 07 65 78 61 6d 70 6c 65" + // pgw2.example
	"d8 00 04 00 01 02 03 04" // Group Id

func pgwHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(strings.ReplaceAll(s, " ", ""))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// TS 29.274 V20.0.0 §8.145: PGW Change Info is a grouped IE.
func TestPGWChangeInfoIsGrouped(t *testing.T) {
	in := pgwHex(t, pgwChangeInfoHex)
	i, err := ie.Parse(in)
	if err != nil {
		t.Fatal(err)
	}
	if !i.IsGrouped() || i.Err() != nil {
		t.Fatalf("IE %d: IsGrouped=%v, Err=%v", i.Type, i.IsGrouped(), i.Err())
	}
	want := []struct{ typ, instance uint8 }{{215, 0}, {74, 0}, {74, 0}, {215, 1}, {216, 0}}
	if len(i.ChildIEs) != len(want) {
		t.Fatalf("%d children, want %d", len(i.ChildIEs), len(want))
	}
	for k, w := range want {
		if c := i.ChildIEs[k]; c.Type != w.typ || c.Instance() != w.instance {
			t.Fatalf("child %d: type %d instance %d, want %d/%d", k, c.Type, c.Instance(), w.typ, w.instance)
		}
	}
	for k, addr := range map[int]string{1: "192.0.2.99", 2: "2001:db8::63"} {
		if got, err := i.ChildIEs[k].IPAddress(); err != nil || got != addr {
			t.Fatalf("child %d: IPAddress %q, %v; want %s", k, got, err, addr)
		}
	}
	if !bytes.Equal(i.ChildIEs[4].Payload, []byte{1, 2, 3, 4}) {
		t.Fatalf("Group Id payload % x", i.ChildIEs[4].Payload)
	}
	out, err := i.Marshal()
	if err != nil || !bytes.Equal(out, in) {
		t.Fatalf("round trip: %v\n in: % x\nout: % x", err, in, out)
	}
}

// A child crossing the end of the group is retained as a content error, like
// any other grouped IE (TS 29.274 §§7.7.1, 7.7.9), and is refused on encode.
func TestPGWChangeInfoChildFraming(t *testing.T) {
	i, err := ie.Parse(pgwHex(t, "d6 00 08 00 4a 00 05 00 c0 00 02 63"))
	if err != nil {
		t.Fatal(err)
	}
	if i.Err() == nil {
		t.Fatal("child crossing the group end was not reported")
	}
	if _, err := i.Marshal(); err == nil {
		t.Fatal("encoded a group with invalid content")
	}
}

// TS 29.274 V20.0.0 Table 8.1-1: IE type values 214 to 223.
func TestIETypeConstants214To223(t *testing.T) {
	for want, got := range []uint8{
		ie.PGWChangeInfo, ie.PGWFQDN, ie.GroupID, ie.PSCellID, ie.UPSecurityPolicy,
		ie.AlternativeIMSI, ie.NFInstanceID, ie.TimerInSeconds, ie.MDTConfigurationNR, ie.URI,
	} {
		if int(got) != 214+want {
			t.Errorf("constant %d has value %d, want %d", want, got, 214+want)
		}
	}
	if !(&ie.IE{Type: ie.PGWChangeInfo}).IsGrouped() {
		t.Error("PGW Change Info is not registered as grouped")
	}
	// Only PGW Change Info is grouped among 214 to 223 (Table 8.1-1).
	for typ := ie.PGWFQDN; typ <= ie.URI; typ++ {
		if (&ie.IE{Type: typ}).IsGrouped() {
			t.Errorf("type %d registered as grouped", typ)
		}
	}
}

// TS 29.274 V20.0.0 Table 8.1-1: IE types 214 to 223 have names.
func TestIETypeNames214To223(t *testing.T) {
	for typ, name := range map[uint8]string{
		214: "PGWChangeInfo",
		215: "PGWFQDN",
		216: "GroupID",
		217: "PSCellID",
		218: "UPSecurityPolicy",
		219: "AlternativeIMSI",
		220: "NFInstanceID",
		221: "TimerInSeconds",
		222: "MDTConfigurationNR",
		223: "URI",
	} {
		if got := ie.New(typ, 0, nil).Name(); got != name {
			t.Errorf("type %d: Name() = %q, want %q", typ, got, name)
		}
	}
}
