package gtpv2

import "testing"

func FuzzReceiveDatagram(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0x48, 36, 0, 50, 0, 0, 0, 0, 0, 0, 1, 0})
	f.Add([]byte{0x60, 1, 0, 4, 0, 0, 1, 0})
	f.Add([]byte{0x40, 1, 0, 9, 0, 0, 1, 0, 3, 0, 1, 0, 1})
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 65535 {
			t.Skip()
		}
		c, _ := testConnection(t)
		c.processDatagram(testPeer(), raw)
	})
}
