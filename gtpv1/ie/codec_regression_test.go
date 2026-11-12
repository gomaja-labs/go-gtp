package ie

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func testNoPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if p := recover(); p != nil {
			t.Errorf("unexpected panic: %v", p)
		}
	}()
	f()
}

func TestTEIDDataII(t *testing.T) {
	// TS 29.060 7.7.15 Figure 23, Table 37: NSAPI followed by four TEID octets.
	i, err := Parse([]byte{18, 5, 0x11, 0x22, 0x33, 0x44})
	if err != nil {
		t.Fatal(err)
	}
	v, err := i.TEID()
	if err != nil || v != 0x11223344 || i.MarshalLen() != 6 {
		t.Fatalf("TEID=%08x len=%d err=%v; want 11223344 and 6", v, i.MarshalLen(), err)
	}
}
func TestAUTN(t *testing.T) {
	want := bytes.Repeat([]byte{0xa5}, 16)
	i := NewAuthenticationQuintuplet(make([]byte, 16), []byte{1, 2, 3, 4}, make([]byte, 16), make([]byte, 16), want)
	got, err := i.AUTN()
	if err != nil || !bytes.Equal(got, want) {
		t.Errorf("AUTN=%x err=%v; want %x", got, err, want)
	}
	testNoPanic(t, func() { _, _ = New(AuthenticationQuintuplet, make([]byte, 50)).AUTN() })
}
func TestRAIThreeDigitMNC(t *testing.T) {
	i := NewRouteingAreaIdentity("310", "260", 1, 2)
	if !bytes.Equal(i.Payload[:3], []byte{0x13, 0x00, 0x62}) {
		t.Fatalf("PLMN=%x; want 130062", i.Payload[:3])
	}
}
func TestIMSI(t *testing.T) {
	for _, s := range []string{"", "1234567890123456", "12345678901234567"} {
		if i := NewIMSI(s); i.Err() == nil {
			t.Errorf("accepted invalid IMSI digit count %d", len(s))
		}
	}
	for _, s := range []string{"12345678901234", "12345"} {
		t.Run(s, func(t *testing.T) {
			i := NewIMSI(s)
			raw, err := i.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			p, err := Parse(raw)
			if err != nil {
				t.Fatal(err)
			}
			got, err := p.IMSI()
			if err != nil || got != s {
				t.Fatalf("IMSI=%q err=%v; want %q", got, err, s)
			}
		})
	}
}
func TestTimeZoneValidation(t *testing.T) {
	if i := NewMSTimeZone(17*time.Minute, 0); i.Err() == nil {
		got, _ := i.TimeZone()
		t.Errorf("unrepresentable 17-minute time zone silently changed to %s", got)
	}
	if got, err := New(MSTimeZone, []byte{0, 3}).DaylightSaving(); err == nil {
		t.Errorf("accepted reserved daylight-saving value %d", got)
	}
}
func TestIdentityDigits(t *testing.T) {
	for _, f := range []func(string) *IE{NewIMSI, NewMSISDN, NewIMEISV} {
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
func TestAPNValidation(t *testing.T) {
	for _, s := range []string{"", "a.", "a..b", strings.Repeat("a", 64), strings.Repeat("a.", 50) + "a"} {
		t.Run(fmt.Sprintf("%d_%q", len(s), s), func(t *testing.T) {
			i := NewAccessPointName(s)
			if i == nil {
				return
			}
			if got, err := i.AccessPointName(); err == nil {
				t.Fatalf("accepted invalid APN %q as %q", s, got)
			}
		})
	}
}
func TestPCOBounds(t *testing.T) {
	t.Run("truncated", func(t *testing.T) {
		if p, err := ParsePCOPayload([]byte{0x80, 0, 13, 1}); err == nil {
			t.Fatalf("accepted truncated PCO: %+v", p)
		}
	})
	t.Run("option_buffer", func(t *testing.T) {
		testNoPanic(t, func() {
			if err := NewConfigurationProtocolOption(13, []byte{1}).MarshalTo(nil); err == nil {
				t.Error("accepted short buffer")
			}
		})
	})
	t.Run("payload_buffer", func(t *testing.T) {
		testNoPanic(t, func() {
			if err := NewPCOPayload(0).MarshalTo(nil); err == nil {
				t.Error("accepted short buffer")
			}
		})
	})
}
func TestUnknownTV(t *testing.T) {
	got, err := ParseMultiIEs([]byte{30, 14, 7})
	if err == nil && len(got) > 1 {
		t.Fatalf("parsed %d IEs after unknown TV type 30", len(got))
	}
}
func TestTimestampBeforeUnixEpoch(t *testing.T) {
	want := time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := NewULITimestamp(want).TimestampNear(want)
	if err != nil || !got.Equal(want) {
		t.Fatalf("timestamp=%v err=%v; want %v", got, err, want)
	}
}
func TestLengthOverflow(t *testing.T) {
	t.Run("extension_list", func(t *testing.T) {
		i := NewExtensionHeaderTypeList(make([]byte, 256)...)
		b, err := i.Marshal()
		if err == nil && int(b[1]) != len(i.Payload) {
			t.Fatalf("encoded 256 types with length %d", b[1])
		}
	})
	t.Run("tlv", func(t *testing.T) {
		i := New(PrivateExtension, make([]byte, 65536))
		_, err := i.Marshal()
		if err == nil {
			t.Fatal("encoded 65536-byte TLV with wrapped length")
		}
	})
	t.Run("pco_option", func(t *testing.T) {
		_, err := NewConfigurationProtocolOption(13, make([]byte, 256)).Marshal()
		if err == nil {
			t.Fatal("encoded 256-byte PCO option with wrapped length")
		}
	})
}
func TestIPLengths(t *testing.T) {
	for _, typ := range []uint8{GSNAddress, EndUserAddress} {
		t.Run(fmt.Sprint(typ), func(t *testing.T) {
			i := New(typ, []byte{0xf1, 0x21, 1, 2, 3})
			if got, err := i.IPAddress(); err == nil {
				t.Fatalf("accepted malformed IP: %s", got)
			}
		})
	}
}
func TestPPPType(t *testing.T) {
	got := NewEndUserAddressPPP().Payload
	if !bytes.Equal(got, []byte{0xf0, 1}) {
		t.Fatalf("PPP payload=%x; want f001", got)
	}
}
func TestULIRAIWidth(t *testing.T) {
	i := NewUserLocationInformationWithRAI("310", "260", 1, 2)
	if len(i.Payload) != 8 {
		t.Fatalf("ULI RAI payload length=%d; want 8 (two-octet RAC)", len(i.Payload))
	}
}
func TestSpareBitsIgnored(t *testing.T) {
	t.Run("EUA", func(t *testing.T) {
		i := New(EndUserAddress, []byte{1, 0x21, 192, 0, 2, 1})
		if got, err := i.IPAddress(); err != nil || got != "192.0.2.1" {
			t.Fatalf("spare bits changed address result: %q, %v", got, err)
		}
	})
	t.Run("daylight", func(t *testing.T) {
		got, err := New(MSTimeZone, []byte{0, 0xfd}).DaylightSaving()
		if err != nil || got != 1 {
			t.Fatalf("daylight=%d err=%v; want 1 ignoring spare bits", got, err)
		}
	})
	got, err := New(NSAPI, []byte{0xf5}).NSAPI()
	if err != nil || got != 5 {
		t.Fatalf("NSAPI=%d err=%v; want 5 ignoring spare bits", got, err)
	}
}
func TestTVLengthMismatch(t *testing.T) {
	if _, err := New(Recovery, []byte{1, 2}).Marshal(); err == nil {
		t.Fatal("silently truncated a two-octet Recovery payload")
	}
}
func TestTypedPayloadValidation(t *testing.T) {
	for _, typ := range []uint8{EndUserAddress, GSNAddress, QoSProfile, IMEISV, UserLocationInformation, MSTimeZone, PrivateExtension, AuthenticationQuintuplet} {
		t.Run(fmt.Sprint(typ), func(t *testing.T) {
			if _, err := Parse([]byte{typ, 0, 0}); err == nil {
				t.Fatal("accepted empty structured IE payload")
			}
		})
	}
}
func TestPCOReuse(t *testing.T) {
	p := NewPCOPayload(0)
	if err := p.UnmarshalBinary([]byte{0x80, 0, 13, 0}); err != nil {
		t.Fatal(err)
	}
	if err := p.UnmarshalBinary([]byte{0x80}); err != nil {
		t.Fatal(err)
	}
	if len(p.ConfigurationProtocolOptions) != 0 {
		t.Fatal("second PCO retained first packet's options")
	}
	c := NewConfigurationProtocolOption(13, []byte{1})
	if err := c.UnmarshalBinary([]byte{0, 13, 0}); err != nil {
		t.Fatal(err)
	}
	if len(c.Contents) != 0 {
		t.Error("zero-length option retained old contents")
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
