package main

import (
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

func TestSampleSubscriberEquipmentIdentities(t *testing.T) {
	cfg, err := loadConfig("enb.yml")
	if err != nil {
		t.Fatal(err)
	}
	for _, sub := range cfg.Subscribers {
		// TS 23.003 §6.2.2: TAC(8), SNR(6), SVN(2) form a 16-digit IMEISV.
		if len(sub.IMEISV) != 16 {
			t.Errorf("sample IMEISV has %d digits, want 16", len(sub.IMEISV))
		}
		i := ie.NewMobileEquipmentIdentity(sub.IMEISV)
		if err := i.MarshalTo(make([]byte, i.MarshalLen())); err != nil {
			t.Errorf("sample IMEISV cannot be sent: %v", err)
		}
	}
}
