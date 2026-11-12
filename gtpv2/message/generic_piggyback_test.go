package message

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// Synthetic datagrams for the two pairs allowed by TS 29.274 §5.5.1. The
// piggybacked messages list their IEs out of table order, which §8.1 allows.
const (
	// Create Session Response (P=1): Cause 16, PGW S5/S8 F-TEID (instance 1).
	csrespPiggyHex = "58 21 00 1b 11 11 11 11 00 01 00 00" +
		"02 00 02 00 10 00" +
		"57 00 09 01 87 22 22 22 22 c0 00 02 01"
	// Create Bearer Request: Bearer Context, unknown IE 230, then the LBI.
	cbreqUnorderedHex = "48 5f 00 2a 33 33 33 33 00 01 01 00" +
		"5d 00 12 00 49 00 01 00 00 57 00 09 01 85 44 44 44 44 c0 00 02 02" +
		"e6 00 03 00 aa bb cc" +
		"49 00 01 00 05"
	// Create Bearer Response (P=1): Cause 16, Bearer Context {EBI 5, Cause 16}.
	cbrespPiggyHex = "58 60 00 1d 55 55 55 55 00 01 02 00" +
		"02 00 02 00 10 00" +
		"5d 00 0b 00 49 00 01 00 05 02 00 02 00 10 00"
	// Modify Bearer Request: unknown IE 230, Bearer Context {EBI 5,
	// S1-U eNodeB F-TEID}, then Recovery.
	mbreqUnorderedHex = "48 22 00 28 66 66 66 66 00 01 03 00" +
		"e6 00 01 00 01" +
		"5d 00 12 00 49 00 01 00 05 57 00 09 00 80 77 77 77 77 c0 00 02 05" +
		"03 00 01 00 09"
)

func piggyHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(strings.ReplaceAll(s, " ", ""))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func piggyWire(t *testing.T, first, second string) (wire []byte, firstLen int) {
	t.Helper()
	f := piggyHex(t, first)
	return append(f, piggyHex(t, second)...), len(f)
}

// TS 29.274 §5.5.1: ParseGeneric keeps both messages of a datagram as Generic,
// so the datagram re-encodes unchanged whatever the IE order of either message.
func TestGenericPiggybackRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name          string
		first, second string
		secondType    uint8
	}{
		{"CreateSessionResponse+CreateBearerRequest", csrespPiggyHex, cbreqUnorderedHex, MsgTypeCreateBearerRequest},
		{"CreateBearerResponse+ModifyBearerRequest", cbrespPiggyHex, mbreqUnorderedHex, MsgTypeModifyBearerRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in, _ := piggyWire(t, tc.first, tc.second)
			g, err := ParseGeneric(in)
			if err != nil {
				t.Fatal(err)
			}
			p, ok := g.Header.Piggybacked.(*Generic)
			if !ok || g.Header.PiggybackError != nil || p.MessageType() != tc.secondType || len(p.IEs) != 3 {
				t.Fatalf("piggybacked message decoded as %T, error %v", g.Header.Piggybacked, g.Header.PiggybackError)
			}
			if out := codecWire(t, g); !bytes.Equal(out, in) {
				t.Fatalf("datagram not re-encoded byte-exact\n in: % x\nout: % x", in, out)
			}
			// Parse keeps decoding the piggybacked message with the typed codec.
			m, err := Parse(in)
			if err != nil {
				t.Fatal(err)
			}
			if _, generic := m.GetHeader().Piggybacked.(*Generic); generic || m.GetHeader().Piggybacked == nil {
				t.Fatalf("Parse returned piggybacked %T", m.GetHeader().Piggybacked)
			}
		})
	}
}

// Rewriting the TEID and sequence number of both messages changes only those
// header fields on the wire.
func TestGenericPiggybackHeaderRewrite(t *testing.T) {
	in, firstLen := piggyWire(t, csrespPiggyHex, cbreqUnorderedHex)
	g, err := ParseGeneric(in)
	if err != nil {
		t.Fatal(err)
	}
	g.SetTEID(0x0a0b0c0d)
	if err := g.SetSequenceNumber(0x000200); err != nil {
		t.Fatal(err)
	}
	p := g.Header.Piggybacked
	if p == nil {
		t.Fatal("no piggybacked message")
	}
	p.SetTEID(0x01020304)
	if err := p.SetSequenceNumber(0xfffffe); err != nil {
		t.Fatal(err)
	}
	want := append([]byte(nil), in...)
	copy(want[4:11], piggyHex(t, "0a 0b 0c 0d 00 02 00"))
	copy(want[firstLen+4:firstLen+11], piggyHex(t, "01 02 03 04 ff ff fe"))
	if out := codecWire(t, g); !bytes.Equal(out, want) {
		t.Fatalf("rewrite changed more than the headers\nwant % x\n got % x", want, out)
	}
}

// TS 29.274 §7.7.3: ParseGeneric applies the same framing rules as Parse.
func TestGenericPiggybackFraming(t *testing.T) {
	wire, firstLen := piggyWire(t, csrespPiggyHex, cbreqUnorderedHex)
	clone := func() []byte { return append([]byte(nil), wire...) }

	// Truncations keep their capacity equal to their length, so a read past
	// the received octets panics instead of seeing the rest of the datagram.
	prefix := func(n int) []byte { return append([]byte(nil), wire[:n]...)[:n:n] }
	longerSecond := clone()
	longerSecond[firstLen+3]++ // second Length claims one octet more than present
	shorterSecond := clone()
	shorterSecond[firstLen+3]-- // second Length leaves one octet unaccounted for
	cases := map[string][]byte{
		"trailing message shorter than a header": prefix(firstLen + 7),
		"truncated trailing message":             prefix(len(wire) - 1),
		"octet after the trailing message":       append(clone(), 0),
		"trailing Length beyond the datagram":    longerSecond,
		"trailing Length short of the datagram":  shorterSecond,
	}
	for n := 0; n < 4; n++ { // P set, but fewer octets follow than a Length field needs
		cases[fmt.Sprintf("%d octets after a P=1 message", n)] = prefix(firstLen + n)
	}
	for name, b := range cases {
		if _, err := ParseGeneric(b); !errors.Is(err, ErrInvalidLength) {
			t.Errorf("%s: got %v, want %v", name, err, ErrInvalidLength)
		}
	}

	// P set on a message that cannot carry a piggybacked message.
	notAllowed := clone()
	notAllowed[1] = MsgTypeCreateSessionRequest
	if _, err := ParseGeneric(notAllowed); !errors.Is(err, ErrInvalidPiggyback) {
		t.Errorf("P on Create Session Request: got %v", err)
	}

	// TS 29.274 §7.7.0: a bad secondary message does not cost the primary.
	nested := clone()
	nested[firstLen] |= 0x10 // the piggybacked message sets P itself
	wrongPair := clone()
	wrongPair[firstLen+1] = MsgTypeDeleteSessionRequest
	badIE := clone()
	badIE[firstLen+12+2] = 0x40 // Bearer Context Length runs past the message
	for name, tc := range map[string]struct {
		b    []byte
		want error
	}{
		"nested piggyback":               {nested, ErrInvalidPiggyback},
		"pair not allowed by §5.5.1":     {wrongPair, ErrInvalidPiggyback},
		"malformed IE in second message": {badIE, nil},
	} {
		g, err := ParseGeneric(tc.b)
		if err != nil || g == nil || g.MessageType() != MsgTypeCreateSessionResponse || len(g.IEs) != 2 {
			t.Fatalf("%s: primary not retained: %v", name, err)
		}
		if g.Header.Piggybacked != nil || g.Header.PiggybackError == nil {
			t.Fatalf("%s: got piggybacked %T, error %v", name, g.Header.Piggybacked, g.Header.PiggybackError)
		}
		if tc.want != nil && !errors.Is(g.Header.PiggybackError, tc.want) {
			t.Fatalf("%s: got %v, want %v", name, g.Header.PiggybackError, tc.want)
		}
		if _, err := g.Marshal(); err == nil {
			t.Fatalf("%s: re-encoded a datagram whose piggybacked message failed to decode", name)
		}
	}
}
