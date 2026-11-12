package message

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

func TestUnexpectedGroupedInstance(t *testing.T) {
	// TS 29.274 V20.0.0 §7.7.9: an unexpected instance is skipped.
	// Delete Session Response, valid Cause, then Bearer Context instance 15.
	// The outer TLIV is intact; its content is irrelevant to this message.
	b := []byte{0x48, 37, 0, 19, 0, 0, 0, 1, 0, 0, 1, 0, 2, 0, 2, 0, 16, 0, 93, 0, 1, 15, 0xff}
	m, err := ParseDeleteSessionResponse(b)
	if err != nil {
		t.Fatalf("unexpected grouped IE prevents response processing: %v", err)
	}
	if len(m.AdditionalIEs) != 1 || m.AdditionalIEs[0].Instance() != 15 || !bytes.Equal(m.AdditionalIEs[0].Payload, []byte{0xff}) {
		t.Fatal("unexpected grouped IE not retained")
	}
	if _, err := m.AdditionalIEs[0].BearerContext(); err == nil {
		t.Fatal("malformed grouped accessor succeeded")
	}
}

func TestIntegerRoleWidths(t *testing.T) {
	// TS 29.274 V20.0.0 §8.118 Figure 8.118-1: UE Usage Type uses one
	// octet; Remaining Running Service Gap Timer uses four even for small values.
	m := NewContextResponse(1, 1, ie.NewCause(64, 0, 0, 0, nil), ie.NewIntegerNumber(5, 2), ie.NewIntegerNumber(5, 2).WithInstance(1))
	b, err := m.Marshal()
	if err != nil {
		return // Rejecting a value whose width does not fit its role is also safe.
	}
	parsed, err := ParseContextResponse(b)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.UEUsageType.Length != 1 || parsed.RemainingRunningServiceGapTimer.Length != 4 {
		t.Fatalf("Integer Number widths = %d/%d; want 1/4 for usage type/service gap timer", parsed.UEUsageType.Length, parsed.RemainingRunningServiceGapTimer.Length)
	}
}

func TestEchoRequestErroneousIE(t *testing.T) {
	// TS 29.274 V20.0.0 §7.1.2 NOTE, §7.7.7: answer Echo despite bad IEs.
	b := []byte{0x40, 1, 0, 8, 0, 0, 1, 0, 3, 0, 2, 0}
	if _, err := ParseEchoRequest(b); err != nil {
		t.Fatalf("Echo Request with erroneous IE cannot be processed: %v", err)
	}
}

func TestPiggybackFailurePreservesPrimary(t *testing.T) {
	// TS 29.274 V20.0.0 §7.7.0: a dedicated bearer error must not affect
	// default bearer establishment. Both outer message lengths are correct.
	b := []byte{0x58, 33, 0, 82, 0, 0, 0, 1, 0, 0, 1, 0,
		2, 0, 2, 0, 16, 0,
		87, 0, 9, 0, 0x8b, 0, 0, 0, 1, 192, 0, 2, 1,
		87, 0, 9, 1, 0x87, 0, 0, 0, 2, 192, 0, 2, 2,
		79, 0, 5, 0, 1, 192, 0, 2, 3,
		127, 0, 1, 0, 0,
		93, 0, 24, 0,
		2, 0, 2, 0, 16, 0,
		73, 0, 1, 0, 5,
		87, 0, 9, 0, 0x81, 0, 0, 0, 3, 192, 0, 2, 1,
		0x48, 95, 0, 12, 0, 0, 0, 1, 0, 0, 2, 0, 93, 0, 1, 0}
	m, err := ParseCreateSessionResponse(b)
	if err != nil || m == nil || m.Cause == nil || m.Cause.MustCause() != 16 {
		t.Fatalf("valid primary response lost because secondary IE is malformed: %v", err)
	}
	var secondary *IEError
	if !errors.As(m.PiggybackError, &secondary) || secondary.Type != ie.BearerContext || !secondary.Framing {
		t.Fatalf("secondary error lost: %v", m.PiggybackError)
	}
	generic, err := Parse(b)
	if err != nil || generic == nil || generic.GetHeader().PiggybackError == nil {
		t.Fatalf("generic decoder lost primary or secondary error: %v", err)
	}
	if _, err := m.Marshal(); err == nil {
		t.Fatal("encoder transmitted an invalid secondary")
	}
}

func TestMandatoryContentErrorPreservesMessage(t *testing.T) {
	// TS 29.274 §§7.7.7–7.7.8: retain the transaction and identify the IE
	// so the application can send Invalid length or Mandatory IE incorrect.
	for _, tc := range []struct {
		name  string
		wire  []byte
		typ   uint8
		cause uint8
	}{
		{"group", []byte{0x48, 95, 0, 13, 0, 0, 0, 7, 0, 0, 9, 0, 93, 0, 1, 0, 0xff}, ie.BearerContext, 67},
		{"nested EBI", []byte{0x48, 95, 0, 16, 0, 0, 0, 7, 0, 0, 9, 0, 93, 0, 4, 0, 73, 0, 0, 0}, ie.EPSBearerID, 67},
		{"APN", []byte{0x48, 32, 0, 14, 0, 0, 0, 7, 0, 0, 9, 0, 71, 0, 2, 0, 1, '_'}, ie.AccessPointName, 69},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m, err := Parse(tc.wire)
			var content *IEError
			if m == nil || !errors.As(err, &content) || content.Framing || content.Type != tc.typ || content.Instance != 0 || content.Cause != tc.cause {
				t.Fatalf("message=%v, error=%v", m, err)
			}
			if m.Sequence() != 9 || m.TEID() != 7 {
				t.Fatal("transaction lost on content error")
			}
			if tc.wire[1] == MsgTypeCreateBearerRequest {
				m, err := ParseCreateBearerRequest(tc.wire)
				if m == nil || m.BearerContexts == nil || !errors.As(err, &content) {
					t.Fatalf("typed decoder lost message: %v", err)
				}
			}
		})
	}
}

func TestIEFramingErrorRejectsMessage(t *testing.T) {
	b := []byte{0x48, 37, 0, 13, 0, 0, 0, 1, 0, 0, 1, 0, 93, 0, 2, 15, 0xff}
	m, err := Parse(b)
	var framing *IEError
	if m != nil || !errors.As(err, &framing) || !framing.Framing || framing.Type != ie.BearerContext || framing.Instance != 15 || !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("message=%v, error=%v", m, err)
	}
}

func TestOptionalGroupedContentErrorKeepsRawIE(t *testing.T) {
	// TS 29.274 §§7.7.1, 7.7.8: optional erroneous content cannot cost Cause.
	b := []byte{0x48, 37, 0, 24, 0, 0, 0, 1, 0, 0, 1, 0,
		2, 0, 2, 0, 16, 0, 181, 0, 1, 0, 0xff, 3, 0, 1, 0, 7}
	m, err := ParseDeleteSessionResponse(b)
	if err != nil || m == nil || m.Cause == nil || m.Recovery.MustRecovery() != 7 {
		t.Fatalf("optional group prevented decoding: %v", err)
	}
	i := m.PGWNodeLoadControlInformation
	if i == nil || !bytes.Equal(i.Payload, []byte{0xff}) {
		t.Fatalf("raw group lost: %v", i)
	}
	if _, err := i.FindByType(ie.SequenceNumber, 0); err == nil {
		t.Fatal("accessor hid malformed grouped content")
	}
	if _, err := m.Marshal(); err == nil {
		t.Fatal("encoder transmitted malformed received content")
	}
	m.PGWNodeLoadControlInformation = nil
	if _, err := m.Marshal(); err != nil {
		t.Fatalf("removing the erroneous optional IE did not repair the message: %v", err)
	}
}

func TestUnexpectedGroupedChildDoesNotCostMandatoryGroup(t *testing.T) {
	// TS 29.274 §7.7.9 and Table 7.2.4-2: ignore unknown child instance.
	b := []byte{0x48, 96, 0, 34, 0, 0, 0, 1, 0, 0, 1, 0,
		2, 0, 2, 0, 16, 0, 93, 0, 16, 0,
		2, 0, 2, 0, 16, 0, 73, 0, 1, 0, 5, 93, 0, 1, 15, 0xff}
	m, err := ParseCreateBearerResponse(b)
	if err != nil || m == nil || len(m.BearerContexts) != 1 {
		t.Fatalf("unexpected child prevented decoding: %v", err)
	}
	c := m.BearerContexts[0]
	if c.MustEPSBearerID() != 5 || c.MustCause() != 16 || len(c.ChildIEs) != 3 || c.ChildIEs[2].Instance() != 15 {
		t.Fatalf("child list lost: %v", c.ChildIEs)
	}
}

func TestEchoRequestKeepsTransactionAndErroneousIE(t *testing.T) {
	for _, tail := range [][]byte{{3, 0, 0, 0}, {3, 0, 2, 0}, {3}} {
		b := append([]byte{0x40, 1, 0, byte(4 + len(tail)), 0x12, 0x34, 0x56, 0}, tail...)
		m, err := ParseEchoRequest(b)
		if err != nil || m == nil || m.Sequence() != 0x123456 {
			t.Fatalf("Echo transaction lost: %v", err)
		}
		if !bytes.Equal(m.Header.Payload, tail) || m.Recovery == nil {
			t.Fatal("erroneous IE lost")
		}
		if _, err := m.Recovery.Recovery(); err == nil {
			t.Fatal("invalid Recovery value accepted")
		}
		if _, err := NewEchoResponse(m.Sequence(), ie.NewRecovery(8)).Marshal(); err != nil {
			t.Fatalf("cannot answer Echo: %v", err)
		}
	}
}

func TestMandatoryExtendableIESuffixIgnored(t *testing.T) {
	// TS 29.274 §7.7.7: discard extra octets of a length-consistent IE.
	b := []byte{0x48, 38, 0, 14, 0, 0, 0, 1, 0, 0, 1, 0, 82, 0, 2, 0, 6, 0xff}
	m, err := ParseChangeNotificationRequest(b)
	if err != nil || m == nil || m.RATType.MustRATType() != 6 || !bytes.Equal(m.RATType.Payload, []byte{6, 0xff}) {
		t.Fatalf("extension octet rejected or lost: %v", err)
	}
}

func TestIntegerNumberMessageRoleWidths(t *testing.T) {
	for _, width := range []int{1, 2, 3, 4} {
		m := NewContextResponse(1, 1, ie.NewIntegerNumber(5, width))
		_, err := m.Marshal()
		if (width == 1) != (err == nil) {
			t.Fatalf("UE Usage Type width %d: %v", width, err)
		}
		m = NewContextResponse(1, 1, ie.NewIntegerNumber(5, width).WithInstance(1))
		_, err = m.Marshal()
		if (width == 4) != (err == nil) {
			t.Fatalf("Service Gap Timer width %d: %v", width, err)
		}
	}
	m := NewContextResponse(1, 1, ie.New(ie.IntegerNumber, 0, nil), ie.NewIntegerNumber(5, 4).WithInstance(1))
	b, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := ParseContextResponse(b)
	if err != nil || decoded.UEUsageType.Length != 0 || decoded.RemainingRunningServiceGapTimer.MustIntegerNumber() != 5 {
		t.Fatalf("unavailable UE Usage Type or timer lost: %v", err)
	}
}

func TestRejectionCauseMakesBearerContextOptional(t *testing.T) {
	// TS 29.274 §6.1.1: Bearer Context may be included optionally on rejection.
	b := []byte{0x48, 33, 0, 19, 0, 0, 0, 1, 0, 0, 1, 0,
		2, 0, 2, 0, 64, 0, 93, 0, 1, 0, 0xff}
	m, err := ParseCreateSessionResponse(b)
	if err != nil || m == nil || m.Cause.MustCause() != 64 || len(m.BearerContextsCreated) != 1 {
		t.Fatalf("optional group prevented rejection handling: %v", err)
	}
}

func TestDecodeReuseClearsContentAndPiggybackErrors(t *testing.T) {
	m := &CreateSessionResponse{}
	badGroup := []byte{0x48, 33, 0, 19, 0, 0, 0, 1, 0, 0, 1, 0, 2, 0, 2, 0, 16, 0, 93, 0, 1, 0, 0xff}
	if err := m.UnmarshalBinary(badGroup); err == nil {
		t.Fatal("mandatory content error missing")
	}
	good := []byte{0x48, 33, 0, 14, 0, 0, 0, 2, 0, 0, 2, 0, 2, 0, 2, 0, 64, 0}
	badPiggyback := append([]byte(nil), good...)
	badPiggyback[0] |= 0x10
	badPiggyback = append(badPiggyback, 0x48, 95, 0, 13, 0, 0, 0, 2, 0, 0, 3, 0, 93, 0, 1, 0, 0xff)
	if err := m.UnmarshalBinary(badPiggyback); err != nil || m.PiggybackError == nil {
		t.Fatalf("secondary error not retained: %v", err)
	}
	if err := m.UnmarshalBinary(good); err != nil || len(m.BearerContextsCreated) != 0 || m.Sequence() != 2 || m.PiggybackError != nil || m.Piggybacked != nil {
		t.Fatalf("reused message retained stale content: %v", err)
	}
	if _, err := m.Marshal(); err != nil {
		t.Fatalf("reused message retained an encoder error: %v", err)
	}
}

func TestUnexpectedCauseInDeletePDNConnectionSetRequest(t *testing.T) {
	// TS 29.274 Table 7.9.1-1 and §7.7.9: Cause is not a Request IE.
	b := []byte{0x48, 101, 0, 12, 0, 0, 0, 0, 0, 0, 1, 0, 2, 0, 0, 0}
	m, err := ParseDeletePDNConnectionSetRequest(b)
	if err != nil || m == nil || len(m.AdditionalIEs) != 1 {
		t.Fatalf("unexpected IE prevented decoding: %v", err)
	}
}
