package message

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// pgwChangeInfoIEHex is a synthetic PGW Change Info (TS 29.274 V20.0.0 §8.145)
// with the children of Table 7.2.2-6.
const pgwChangeInfoIEHex = "d6 00 49 00" +
	"d7 00 10 00 07 70 67 77 2d 73 65 74 07 65 78 61 6d 70 6c 65" + // PGW Set FQDN pgw-set.example
	"4a 00 04 00 c0 00 02 63" + // Alternative PGW-C/SMF IP Address 192.0.2.99
	"4a 00 10 00 20 01 0d b8 00 00 00 00 00 00 00 00 00 00 00 63" + // and 2001:db8::63
	"d7 00 0d 01 04 70 67 77 32 07 65 78 61 6d 70 6c 65" + // Alternative PGW-C/SMF FQDN pgw2.example
	"d8 00 04 00 01 02 03 04" // Group Id

// pgwChangeInfoMessage builds a message of type typ (T=1) holding Cause 16,
// an F-TEID and the PGW Change Info IE.
func pgwChangeInfoMessage(t *testing.T, typ uint8) []byte {
	t.Helper()
	ies, err := hex.DecodeString(strings.ReplaceAll("02 00 02 00 10 00"+
		"57 00 09 01 87 22 22 22 22 c0 00 02 01"+
		pgwChangeInfoIEHex, " ", ""))
	if err != nil {
		t.Fatal(err)
	}
	b := []byte{0x48, typ, 0, 0, 0x11, 0x11, 0x11, 0x11, 0x00, 0x01, 0x00, 0x00}
	b = append(b, ies...)
	binary.BigEndian.PutUint16(b[2:4], uint16(len(b)-4))
	return b
}

// PGW Change Info is carried on S5/S8 in Create Session Response (Table
// 7.2.2-6), Create Bearer Request (Table 7.2.3-5) and Modify Bearer Response
// (Table 7.2.8-6). Its children decode, and each message re-encodes unchanged.
func TestPGWChangeInfoInS5S8Messages(t *testing.T) {
	for _, typ := range []uint8{MsgTypeCreateSessionResponse, MsgTypeCreateBearerRequest, MsgTypeModifyBearerResponse} {
		in := pgwChangeInfoMessage(t, typ)
		g, err := ParseGeneric(in)
		if err != nil {
			t.Fatalf("type %d: %v", typ, err)
		}
		pci := g.IEs[2]
		if pci.Type != 214 || !pci.IsGrouped() || len(pci.ChildIEs) != 5 || pci.Err() != nil {
			t.Fatalf("type %d: IE %d grouped=%v children=%d err=%v", typ, pci.Type, pci.IsGrouped(), len(pci.ChildIEs), pci.Err())
		}
		if addr, err := pci.ChildIEs[2].IPAddress(); err != nil || addr != "2001:db8::63" {
			t.Fatalf("type %d: child IP Address %q, %v", typ, addr, err)
		}
		if out := codecWire(t, g); !bytes.Equal(out, in) {
			t.Fatalf("type %d: Generic round trip changed\n in: % x\nout: % x", typ, in, out)
		}
	}

	// The typed Create Session Response keeps the IE, grouped, among AdditionalIEs.
	in := pgwChangeInfoMessage(t, MsgTypeCreateSessionResponse)
	m, err := ParseCreateSessionResponse(in)
	if err != nil {
		t.Fatal(err)
	}
	var pci *ie.IE
	for _, a := range m.AdditionalIEs {
		if a.Type == 214 {
			pci = a
		}
	}
	if pci == nil || len(pci.ChildIEs) != 5 {
		t.Fatalf("typed Create Session Response: PGW Change Info %v", pci)
	}
	if out := codecWire(t, m); !bytes.Equal(out, in) {
		t.Fatalf("typed round trip changed\n in: % x\nout: % x", in, out)
	}
}
