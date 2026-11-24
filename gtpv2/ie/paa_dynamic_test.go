package ie

import (
	"bytes"
	"testing"
)

// TS 29.274 V20.0.0 §7.2.1 Table 7.2.1-1 permits the all-zero dynamic
// IPv6 request, unlike the assigned /64 prefix described in §8.14.
func TestPAADynamicIPv6(t *testing.T) {
	for _, typ := range []uint8{2, 3} {
		t.Run(string(rune('0'+typ)), func(t *testing.T) {
			size := 18
			if typ == 3 {
				size = 22
			}
			wire := make([]byte, size)
			wire[0] = typ
			decoded, err := ParsePDNAddressAllocationFields(wire)
			if err != nil {
				t.Fatal(err)
			}
			got, err := decoded.Marshal()
			if err != nil || !bytes.Equal(got, wire) {
				t.Fatalf("round trip: %x, %v", got, err)
			}
			wire[17] = 1
			if _, err := ParsePDNAddressAllocationFields(wire); err == nil {
				t.Fatal("zero prefix with nonzero address accepted")
			}
		})
	}
	for _, value := range []*IE{NewPDNAddressAllocationIPv6("::", 0), NewPDNAddressAllocationDual("0.0.0.0", "::", 0), NewPDNAddressAllocationDual("192.0.2.1", "::", 0)} {
		if err := value.Err(); err != nil {
			t.Fatal(err)
		}
		if _, err := ParsePDNAddressAllocationFields(value.Payload); err != nil {
			t.Fatal(err)
		}
	}
}
