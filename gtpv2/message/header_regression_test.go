package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

func TestHeaderBoundaries(t *testing.T) {
	for _, typ := range []uint8{MsgTypeEchoRequest, MsgTypeVersionNotSupportedIndication, MsgTypeCreateSessionRequest} {
		flags, headerSize := uint8(0x48), 12
		if !requiresTEID(typ) {
			flags, headerSize = 0x40, 8
		}
		h := NewHeader(flags, typ, 0x12345678, 0xffffff, nil)
		for size := 0; size < h.MarshalLen(); size++ {
			if err := h.MarshalTo(make([]byte, size)); !errors.Is(err, ErrTooShortToMarshal) {
				t.Fatalf("short size %d: %v", size, err)
			}
		}
		for _, length := range []int{65534, 65535, 65536} {
			h.Payload = make([]byte, length-(headerSize-4))
			b, err := h.Marshal()
			if length > 65535 {
				if !errors.Is(err, ErrInvalidLength) {
					t.Fatalf("overflow: %v", err)
				}
				continue
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(b) != h.MarshalLen() || int(binary.BigEndian.Uint16(b[2:4])) != length {
				t.Fatal("maximum length changed")
			}
			parsed, err := ParseHeader(b)
			if err != nil || parsed.Sequence() != 0xffffff {
				t.Fatalf("maximum parse: %v", err)
			}
		}
	}
	h := NewHeader(0x48, 32, 0, 1, nil)
	if err := h.SetSequenceNumber(0x1000000); !errors.Is(err, ErrInvalidSequence) || h.Sequence() != 1 {
		t.Fatal("sequence overflow changed header")
	}
	if err := h.SetSequenceNumber(0xffffff); err != nil {
		t.Fatal(err)
	}
	for priority := uint8(0); priority < 16; priority++ {
		if err := h.SetMessagePriority(priority); err != nil {
			t.Fatal(err)
		}
		b, err := h.Marshal()
		if err != nil || b[0]&4 == 0 || b[11] != priority<<4 || h.MessagePriority() != priority {
			t.Fatalf("priority %d: %x %v", priority, b, err)
		}
	}
	if err := h.SetMessagePriority(16); !errors.Is(err, ErrInvalidPriority) || h.MessagePriority() != 15 {
		t.Fatal("priority overflow")
	}
}

func TestHeaderSpareBits(t *testing.T) {
	for _, flags := range []byte{0x40, 0x44, 0x48, 0x4c} {
		typ := MsgTypeCreateSessionRequest
		if flags == 0x40 || flags == 0x44 {
			typ = MsgTypeEchoRequest
		}
		h := NewHeader(flags, typ, 0, 1, nil)
		b, err := h.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		b[0] |= 3
		b[len(b)-1] = 0xaf
		parsed, err := ParseHeader(b)
		if err != nil {
			t.Fatal(err)
		}
		if flags == 0x44 && parsed.HasMessagePriority() {
			t.Fatal("spare MP bit treated as a priority flag")
		}
		canonical, err := parsed.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		wantPriority := byte(0)
		if flags == 0x4c {
			wantPriority = 0xa0
		}
		if canonical[0]&3 != 0 || canonical[len(canonical)-1] != wantPriority {
			t.Fatalf("spare bits emitted: %x", canonical)
		}
	}
}

func TestPiggybackFraming(t *testing.T) {
	for _, pair := range []struct{ first, second Message }{
		{NewCreateSessionResponse(1, 1, ie.NewCause(16, 0, 0, 0, nil)), NewCreateBearerRequest(2, 2, ie.NewEPSBearerID(5))},
		{NewCreateBearerResponse(3, 3, ie.NewCause(16, 0, 0, 0, nil)), NewModifyBearerRequest(4, 4, ie.NewRecovery(1))},
	} {
		firstOnly := codecWire(t, pair.first)
		secondOnly := codecWire(t, pair.second)
		pair.first.GetHeader().Piggybacked = pair.second
		wire := codecWire(t, pair.first)
		if len(wire) != len(firstOnly)+len(secondOnly) || pair.first.MarshalLen() != len(wire) || wire[0]&0x10 == 0 || wire[len(firstOnly)]&0x10 != 0 || int(binary.BigEndian.Uint16(wire[2:4])) != len(firstOnly)-4 {
			t.Fatal("inconsistent piggyback encoding")
		}
		dst := bytes.Repeat([]byte{0xa5}, len(wire)+8)
		if err := pair.first.MarshalTo(dst); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(dst[:len(wire)], wire) || !bytes.Equal(dst[len(wire):], bytes.Repeat([]byte{0xa5}, 8)) {
			t.Fatal("MarshalLen/write mismatch")
		}
		decoded, err := Parse(wire)
		if err != nil {
			t.Fatal(err)
		}
		if decoded.GetHeader().Piggybacked == nil || decoded.GetHeader().Piggybacked.Sequence() != pair.second.Sequence() {
			t.Fatal("missing child")
		}
		if !bytes.Equal(codecWire(t, decoded), wire) {
			t.Fatal("piggyback round trip changed bytes")
		}
		for _, invalid := range [][]byte{wire[:len(firstOnly)], wire[:len(wire)-1], append(append([]byte{}, wire...), 0)} {
			if _, err := Parse(invalid); !errors.Is(err, ErrInvalidLength) {
				t.Fatalf("invalid framing: %v", err)
			}
		}
		nested := append([]byte{}, wire...)
		nested[len(firstOnly)] |= 0x10
		primary, err := Parse(nested)
		if err != nil || primary == nil || !errors.Is(primary.GetHeader().PiggybackError, ErrInvalidPiggyback) {
			t.Fatal("nested P did not retain primary with a secondary error")
		}
		pair.second.GetHeader().Piggybacked = pair.first
		if _, err := Marshal(pair.first); !errors.Is(err, ErrInvalidPiggyback) {
			t.Fatal("cycle accepted")
		}
	}
	first := NewCreateSessionResponse(0, 1)
	first.Piggybacked = (*CreateBearerRequest)(nil)
	if _, err := Marshal(first); !errors.Is(err, ErrInvalidPiggyback) {
		t.Fatal("typed nil child accepted")
	}
	first.Piggybacked = NewEchoRequest(2)
	if _, err := Marshal(first); !errors.Is(err, ErrInvalidPiggyback) {
		t.Fatal("invalid pair accepted")
	}
}

func TestAdditionalPrivateExtensionInstances(t *testing.T) {
	for _, factory := range codecFactories {
		t.Run(factory.name, func(t *testing.T) {
			a, b := ie.NewPrivateExtension(123, []byte{1}), ie.NewPrivateExtension(456, []byte{2}).WithInstance(1)
			m := factory.make(a, b)
			wire := codecWire(t, m)
			if len(wire) != m.GetHeader().headerLen()+a.MarshalLen()+b.MarshalLen() {
				t.Fatal("different vendor instance lost")
			}
			parsed, err := Parse(wire)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(codecWire(t, parsed), wire) {
				t.Fatal("different vendor instance lost on parse")
			}
		})
	}
}
