package gtpv1_test

import (
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv1"
)

func FuzzDecapsulate(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		if _, _, err := gtpv1.Decapsulate(b); err != nil {
			t.Skip()
		}
	})
}
