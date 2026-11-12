package ie

import "testing"

func TestOrdinaryPCOLengthLimit(t *testing.T) {
	// TS 29.060 V19.0.0 §7.7.31 Figure 45 carries octets 3-z of
	// TS 24.008 V20.1.0 §10.5.6.3 Figure 10.5.136: at most 251 value octets.
	for _, n := range []int{248, 252} {
		i := NewProtocolConfigurationOptions(0, NewConfigurationProtocolOption(0xff00, make([]byte, n)))
		if b, err := i.Marshal(); err == nil {
			t.Errorf("PCO accepted %d value octets; maximum is 251", len(b)-3)
		}
	}
}

func TestOrdinaryPCOBoundaryAndRawValidation(t *testing.T) {
	for _, size := range []int{247, 248} {
		p := NewPCOPayload(0, NewConfigurationProtocolOption(0xff00, make([]byte, size)))
		b, err := p.Marshal()
		if size == 247 {
			if err != nil || len(b) != 251 || p.MarshalLen() != len(b) {
				t.Fatalf("maximum PCO: %d, %v", len(b), err)
			}
			if _, err := ParsePCOPayload(b); err != nil {
				t.Fatal(err)
			}
			if err := p.MarshalTo(make([]byte, 250)); err == nil {
				t.Fatal("accepted short destination")
			}
		} else {
			if err == nil {
				t.Fatal("accepted oversized fields")
			}
			// Bypass the constructor to cover received and directly populated IEs.
			raw := append([]byte{0x80, 0xff, 0, byte(size)}, make([]byte, size)...)
			if _, err := ParsePCOPayload(raw); err == nil {
				t.Fatal("parsed oversized PCO")
			}
			if _, err := New(ProtocolConfigurationOptions, raw).Marshal(); err == nil {
				t.Fatal("encoded oversized raw PCO")
			}
		}
	}
}
