package ie_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

func TestLosslessMarshalDoesNotValidateContents(t *testing.T) {
	i := ie.New(ie.AccessPointName, []byte{0xff})
	if i.Err() == nil {
		t.Fatal("expected invalid APN")
	}
	out := bytes.Repeat([]byte{0xcc}, i.MarshalLen()+1)
	if err := i.MarshalToLossless(out); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, []byte{131, 0, 1, 0xff, 0xcc}) {
		t.Fatalf("lossless = %x", out)
	}
	for n := 0; n < i.MarshalLen(); n++ {
		if err := i.MarshalToLossless(make([]byte, n)); !errors.Is(err, ie.ErrTooShortToMarshal) {
			t.Fatalf("short %d: %v", n, err)
		}
	}
	if err := (*ie.IE)(nil).MarshalToLossless(nil); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	if err := ie.NewIMSI("invalid").MarshalToLossless(nil); err == nil {
		t.Fatal("ignored construction error")
	}
	i.Length++
	if err := i.MarshalToLossless(out); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("length mismatch: %v", err)
	}
	i.Payload = make([]byte, 65536)
	i.Length = 0
	if err := i.MarshalToLossless(make([]byte, i.MarshalLen())); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("length overflow: %v", err)
	}
}

func TestLosslessMarshalRetainsPayloadSpares(t *testing.T) {
	wire := []byte{20, 0xf5}
	i, err := ie.Parse(wire)
	if err != nil {
		t.Fatal(err)
	}
	out := make([]byte, i.MarshalLen())
	if err := i.MarshalToLossless(out); err != nil || !bytes.Equal(out, wire) {
		t.Fatalf("spares: %x %v", out, err)
	}
	strict, err := i.Marshal()
	if err != nil || !bytes.Equal(strict, []byte{20, 5}) {
		t.Fatalf("strict: %x %v", strict, err)
	}
}

func TestLosslessMarshalRejectsExtensionListOverflow(t *testing.T) {
	i := ie.New(ie.ExtensionHeaderTypeList, make([]byte, 256))
	if err := i.MarshalToLossless(make([]byte, i.MarshalLen())); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("one-octet length overflow: %v", err)
	}
}
