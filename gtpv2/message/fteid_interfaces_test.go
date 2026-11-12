package message_test

import (
	"encoding/hex"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2"
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

func TestFTEIDN26AndN19mbInterfaces(t *testing.T) {
	for _, tc := range []struct{ value, want uint8 }{{gtpv2.IFTypeN26AMFGTPC, 40}, {gtpv2.IFTypeN19mbUPFGTPU, 41}} {
		if tc.value != tc.want {
			t.Fatalf("interface %d want %d", tc.value, tc.want)
		}
		b, err := ie.NewFullyQualifiedTEID(tc.value, 1, "192.0.2.1", "").Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if b[4] != 0x80|tc.want {
			t.Fatalf("interface flags %x", b[4])
		}
		decoded, err := ie.Parse(b)
		if err != nil {
			t.Fatal(err)
		}
		v, err := decoded.InterfaceType()
		if err != nil || v != tc.want {
			t.Fatalf("interface %d %v", v, err)
		}
	}
	b, _ := hex.DecodeString("57000900a900000001c0000201")
	i, err := ie.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	v, err := i.InterfaceType()
	if err != nil || v != gtpv2.IFTypeN19mbUPFGTPU {
		t.Fatalf("interface %d %v", v, err)
	}
}
