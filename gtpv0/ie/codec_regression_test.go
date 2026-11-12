package ie

import (
	"bytes"
	"strings"
	"testing"
)

func TestAPNValidation(t *testing.T) {
	i := NewAccessPointName(strings.Repeat("a", 64))
	if i == nil {
		return
	}
	if _, err := i.Marshal(); err == nil {
		t.Fatal("encoded an APN label longer than 63 octets")
	}
}

func TestIMSI(t *testing.T) {
	s := "12345"
	b, err := NewIMSI(s).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	i, err := Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	got, err := i.IMSI()
	if err != nil || got != s {
		t.Fatalf("IMSI=%q err=%v; want %q", got, err, s)
	}
}
func TestIdentityDigits(t *testing.T) {
	for _, f := range []func(string) *IE{NewIMSI, NewMSISDN} {
		if i := f("1234567890123ab"); i.Err() == nil {
			t.Errorf("accepted nondecimal identity: type=%d payload=%x", i.Type, i.Payload)
		}
	}
}
func TestMSISDNOdd(t *testing.T) {
	got, err := NewMSISDN("12345").MSISDN()
	if err != nil || got != "12345" {
		t.Fatalf("MSISDN=%q err=%v; want 12345", got, err)
	}
}
func TestAPNExactEnd(t *testing.T) {
	got, err := NewAccessPointName("internet").AccessPointName()
	if err != nil || got != "internet" {
		t.Fatalf("APN=%q err=%v; want internet", got, err)
	}
}
func TestUnknownTV(t *testing.T) {
	got, err := ParseMultiIEs([]byte{30, 14, 7})
	if err == nil && len(got) > 1 {
		t.Fatalf("parsed %d IEs after unknown TV type 30", len(got))
	}
}
func TestQoSAccessors(t *testing.T) {
	i := NewQualityOfServiceProfile(2, 3, 4, 2, 18)
	delay, _ := i.QoSDelay()
	peak, _ := i.QoSPeak()
	mean, _ := i.QoSMean()
	if delay != 2 || peak != 4 || mean != 18 {
		t.Fatalf("delay=%d peak=%d mean=%d; want 2,4,18", delay, peak, mean)
	}
}
func TestIPLengths(t *testing.T) {
	if got, err := New(GSNAddress, make([]byte, 5)).GSNAddress(); err == nil {
		t.Fatalf("accepted malformed IP: %s", got)
	}
}
func TestPPPType(t *testing.T) {
	got := NewEndUserAddressPPP().Payload
	if !bytes.Equal(got, []byte{0xf0, 1}) {
		t.Fatalf("PPP payload=%x; want f001", got)
	}
}
func TestRAIAccessors(t *testing.T) {
	i := NewRouteingAreaIdentity("310", "260", 1, 2)
	mcc, _ := i.MCC()
	mnc, _ := i.MNC()
	if mcc != "310" || mnc != "260" {
		t.Fatalf("MCC=%q MNC=%q; want 310,260", mcc, mnc)
	}
}
func TestLengthOverflow(t *testing.T) {
	if _, err := New(PrivateExtension, make([]byte, 65536)).Marshal(); err == nil {
		t.Fatal("encoded 65536-byte TLV with wrapped length")
	}
}
func TestTVLengthMismatch(t *testing.T) {
	if _, err := New(Recovery, []byte{1, 2}).Marshal(); err == nil {
		t.Fatal("silently truncated a two-octet Recovery payload")
	}
}

func TestAPNSyntaxAndBoundary(t *testing.T) {
	for _, s := range []string{"-apn", "apn-", "apn_name", "a\x00b", "é.example"} {
		if NewAccessPointName(s).Err() == nil {
			t.Fatalf("accepted APN %q", s)
		}
		b := append([]byte{AccessPointName, 0, byte(len(s) + 1), byte(len(s))}, []byte(s)...)
		if _, err := Parse(b); err == nil {
			t.Fatalf("parsed APN %q", s)
		}
	}
	s := strings.Repeat("a", 63) + "." + strings.Repeat("b", 35)
	i := NewAccessPointName(s)
	if i.Err() != nil || len(i.Payload) != 100 {
		t.Fatalf("100-octet APN: %v", i.Err())
	}
	if NewAccessPointName(s+"b").Err() == nil {
		t.Fatal("accepted 101-octet APN")
	}
}
