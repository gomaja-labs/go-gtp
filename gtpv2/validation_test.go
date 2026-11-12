package gtpv2

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

func TestMandatoryIEMessageTables(t *testing.T) {
	// TS 29.274 V20.0.0 Tables 7.1.1-1 through 7.9.6-1, for modeled messages.
	// Numeric type/instance pairs make this independent of the implementation table.
	expected := map[uint8][][2]uint8{
		1: {{3, 0}}, 2: {{3, 0}}, 3: {},
		32: {{82, 0}, {87, 0}, {71, 0}, {93, 0}}, 33: {{2, 0}, {93, 0}},
		34: {}, 35: {{2, 0}}, 36: {}, 37: {{2, 0}}, 38: {{82, 0}}, 39: {{2, 0}},
		64: {{72, 0}, {93, 0}}, 65: {{2, 0}}, 66: {{93, 0}}, 67: {{2, 0}, {93, 0}},
		68: {{73, 0}, {100, 0}, {85, 0}}, 69: {{2, 0}, {73, 0}, {100, 0}},
		70: {{2, 0}}, 73: {}, 95: {{73, 0}, {93, 0}}, 96: {{2, 0}, {93, 0}},
		97: {{93, 0}, {72, 0}}, 98: {{2, 0}, {93, 0}}, 99: {}, 100: {{2, 0}},
		101: {}, 102: {{2, 0}}, 130: {}, 131: {{2, 0}}, 132: {{2, 0}},
		149: {{2, 0}}, 150: {{2, 0}}, 162: {}, 163: {{2, 0}}, 164: {{1, 0}}, 165: {{2, 0}},
		170: {}, 171: {{2, 0}}, 176: {}, 177: {{2, 0}}, 179: {}, 180: {{2, 0}}, 200: {}, 201: {{2, 0}}, 211: {}, 212: {{2, 0}},
	}
	// Check actual wire type values for the restoration messages explicitly.
	expected[message.MsgTypePGWRestartNotification] = [][2]uint8{{74, 0}, {74, 1}}
	expected[message.MsgTypePGWRestartNotificationAcknowledge] = [][2]uint8{{2, 0}}
	if len(expected) != len(mandatoryIEs) {
		t.Fatalf("modeled messages: snapshot=%d table=%d", len(expected), len(mandatoryIEs))
	}
	for typ, wants := range expected {
		got, ok := mandatoryIEs[typ]
		if !ok || len(got) != len(wants) {
			t.Errorf("type %d: %v want %v", typ, got, wants)
			continue
		}
		for j, want := range wants {
			if got[j].typ != want[0] || got[j].instance != want[1] {
				t.Errorf("type %d pair %d: %v want %v", typ, j, got[j], want)
			}
		}
	}
}

func TestMandatoryRequestReplyAndOffendingIE(t *testing.T) {
	peer := udpPeer(t)
	c := udpServer(t, nil)
	req := message.NewCreateSessionRequest(0, 27)
	writeUDP(t, peer, c.LocalAddr(), req)
	raw, _ := readUDP(t, peer)
	rsp := parseUDP(t, raw).(*message.CreateSessionResponse)
	if rsp.Sequence() != 27 || rsp.Cause.MustCause() != CauseMandatoryIEMissing || len(rsp.Cause.Payload) != 6 || rsp.Cause.Payload[2] != ie.RATType || rsp.Cause.Payload[5] != 0 {
		t.Fatalf("missing IE reply %x", raw)
	}
	// Nested required EBI must be checked in every bearer, including later contexts.
	req = message.NewCreateSessionRequest(0, 28, ie.NewRATType(6), ie.NewFullyQualifiedTEID(10, 1, "192.0.2.1", ""), ie.NewAccessPointName("internet"),
		ie.NewBearerContext(ie.NewEPSBearerID(5), ie.NewBearerQoS(0, 1, 0, 9, 0, 0, 0, 0)), ie.NewBearerContext(ie.NewBearerQoS(0, 1, 0, 9, 0, 0, 0, 0)))
	writeUDP(t, peer, c.LocalAddr(), req)
	raw, _ = readUDP(t, peer)
	rsp = parseUDP(t, raw).(*message.CreateSessionResponse)
	if rsp.Cause.MustCause() != CauseMandatoryIEMissing || rsp.Cause.Payload[2] != ie.EPSBearerID {
		t.Fatalf("nested missing IE reply %x", raw)
	}
}

func TestMandatoryResponseNotifiesAndStopsRetry(t *testing.T) {
	failures := make(chan error, 2)
	called := make(chan struct{}, 1)
	c := udpServer(t, func(c *Conn) {
		c.SetErrorHandler(func(err error) { failures <- err })
		if err := c.SetRetransmission(20*time.Millisecond, 1); err != nil {
			t.Fatal(err)
		}
		c.AddHandler(message.MsgTypeCreateSessionResponse, func(_ *Conn, _ net.Addr, _ message.Message) error { called <- struct{}{}; return nil })
	})
	peer := udpPeer(t)
	seq, err := c.SendMessageTo(message.NewCreateSessionRequest(0, 0), peer.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}
	_, _ = readUDP(t, peer)
	writeUDP(t, peer, c.LocalAddr(), message.NewCreateSessionResponse(0, seq))
	select {
	case err := <-failures:
		var missing *RequiredIEMissingError
		if !errors.As(err, &missing) || missing.Type != ie.Cause {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("missing application error")
	}
	noUDP(t, peer)
	select {
	case <-called:
		t.Fatal("invalid response delivered")
	default:
	}
	select {
	case err := <-failures:
		t.Fatalf("unexpected retry failure: %v", err)
	default:
	}
}

func TestPiggybackReceiveOrder(t *testing.T) {
	order := make(chan string, 2)
	c := udpServer(t, func(c *Conn) {
		c.DisableValidation()
		c.AddHandler(message.MsgTypeCreateSessionResponse, func(_ *Conn, _ net.Addr, _ message.Message) error { order <- "response"; return nil })
		c.AddHandler(message.MsgTypeCreateBearerRequest, func(_ *Conn, _ net.Addr, _ message.Message) error { order <- "request"; return nil })
	})
	peer := udpPeer(t)
	seq, err := c.SendMessageTo(message.NewCreateSessionRequest(0, 0), peer.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}
	_, _ = readUDP(t, peer)
	rsp := message.NewCreateSessionResponse(0, seq)
	rsp.Piggybacked = message.NewCreateBearerRequest(0, 987)
	writeUDP(t, peer, c.LocalAddr(), rsp)
	for _, want := range []string{"response", "request"} {
		select {
		case got := <-order:
			if got != want {
				t.Fatalf("%s before %s", got, want)
			}
		case <-time.After(time.Second):
			t.Fatal("piggyback lost")
		}
	}
}

func TestCommandTriggeredExchange(t *testing.T) {
	c, p := testConnection(t)
	defer func() { _ = c.Close() }()
	c.DisableValidation()
	seq, err := c.SendMessageTo(message.NewDeleteBearerCommand(0, 0), testPeer())
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	c.AddHandler(message.MsgTypeDeleteBearerRequest, func(c *Conn, peer net.Addr, req message.Message) error {
		calls++
		return c.RespondTo(peer, req, message.NewDeleteBearerResponse(0, 0, ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil)))
	})
	req := message.NewDeleteBearerRequest(0, seq, ie.NewEPSBearerID(5))
	if err = c.handleMessage(testPeer(), req); err != nil {
		t.Fatal(err)
	}
	if err = c.handleMessage(testPeer(), req); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatal(calls)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, raw := range p.writes {
		m := parseUDP(t, raw)
		if m.Sequence() != seq || m.Sequence()&0x800000 == 0 {
			t.Fatal("command sequence lost")
		}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.transactions) != 0 {
		t.Fatal("command not correlated")
	}
}

func TestGroupedMandatoryAndUnexpectedIEs(t *testing.T) {
	c, _ := testConnection(t)
	c.SetErrorHandler(func(error) {})
	tests := []struct {
		msg     message.Message
		missing uint8
	}{
		{message.NewContextResponse(0, 1, ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil), ie.New(ie.PDNConnection, 0, nil)), ie.AccessPointName},
		{message.NewModifyBearerResponse(0, 1, ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil), ie.New(ie.OverloadControlInformation, 0, nil)), ie.SequenceNumber},
	}
	for _, tc := range tests {
		err := c.validateMandatory(testPeer(), tc.msg)
		var missing *RequiredIEMissingError
		if !errors.As(err, &missing) || missing.Type != tc.missing {
			t.Fatalf("%s: %v", tc.msg.MessageTypeName(), err)
		}
	}
	// A newer optional grouped IE, undefined in the Stop Paging message table,
	// must not be interpreted as a defined group with missing children (§7.7.9).
	if err := c.validateMandatory(testPeer(), message.NewStopPagingIndication(0, 2, ie.New(ie.OverloadControlInformation, 0, nil))); err != nil {
		t.Fatal(err)
	}
	if err := c.validateMandatory(testPeer(), message.NewCreateSessionResponse(0, 3, ie.NewCause(CauseContextNotFound, 0, 0, 0, nil))); err != nil {
		t.Fatalf("rejection requires only Cause: %v", err)
	}
}

func TestPiggybackReliableDelivery(t *testing.T) {
	c, p := testConnection(t)
	if err := c.SetRetransmission(10*time.Millisecond, 1); err != nil {
		t.Fatal(err)
	}
	failures := make(chan error, 1)
	c.SetErrorHandler(func(err error) { failures <- err })
	response := message.NewCreateSessionResponse(0, 1, ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil))
	response.Piggybacked = message.NewCreateBearerRequest(0, 10)
	if err := c.RespondTo(testPeer(), message.NewCreateSessionRequest(0, 1), response); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-failures:
		var tx *TransactionError
		if !errors.As(err, &tx) || tx.Sequence != 10 {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("no piggyback transfer failure")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.writes) != 2 || !bytes.Equal(p.writes[0], p.writes[1]) {
		t.Fatal("piggyback retry did not resend complete original datagram")
	}
}

func TestEveryModeledMandatoryIEIsEnforced(t *testing.T) {
	c, _ := testConnection(t)
	c.SetErrorHandler(func(error) {})
	for typ, required := range mandatoryIEs {
		t.Run(fmt.Sprintf("message-%d", typ), func(t *testing.T) {
			ies := make([]*ie.IE, 0, len(required))
			for _, req := range required {
				var i *ie.IE
				switch req.typ {
				case ie.Cause:
					i = ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil)
				case ie.BearerContext:
					i = ie.NewBearerContext(ie.NewEPSBearerID(5), ie.NewBearerQoS(0, 1, 0, 9, 0, 0, 0, 0), ie.New(ie.BearerTFT, 0, nil), ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil))
				default:
					i = ie.New(req.typ, req.instance, nil)
				}
				ies = append(ies, i.WithInstance(req.instance))
			}
			if err := c.validateMandatory(testPeer(), message.NewGeneric(typ, 0, 1, ies...)); err != nil {
				t.Fatalf("complete mandatory set rejected: %v", err)
			}
			for j, req := range required {
				omitted := append([]*ie.IE(nil), ies[:j]...)
				omitted = append(omitted, ies[j+1:]...)
				err := c.validateMandatory(testPeer(), message.NewGeneric(typ, 0, uint32(j+2), omitted...))
				var missing *RequiredIEMissingError
				if !errors.As(err, &missing) || missing.Type != req.typ || missing.Instance != req.instance {
					t.Fatalf("omitted type %d instance %d: %v", req.typ, req.instance, err)
				}
			}
		})
	}
}

func TestInvalidPiggybackOverallLength(t *testing.T) {
	failures := make(chan error, 2)
	c := udpServer(t, func(c *Conn) { c.SetErrorHandler(func(err error) { failures <- err }) })
	peer := udpPeer(t)
	seq, err := c.SendMessageTo(message.NewCreateSessionRequest(0, 0), peer.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}
	_, _ = readUDP(t, peer)
	rsp := message.NewCreateSessionResponse(0, seq, ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil))
	rsp.Piggybacked = message.NewCreateBearerRequest(0, 321)
	raw := append(mustMarshal(t, rsp), 0)
	if _, err = peer.WriteTo(raw, c.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	raw, _ = readUDP(t, peer)
	reply := parseUDP(t, raw).(*message.CreateBearerResponse)
	if reply.Sequence() != 321 || reply.Cause.MustCause() != CauseInvalidOverallLengthOfTheTriggeredResponseMessageAndAPiggybackedInitialMessage {
		t.Fatalf("invalid overall length reply: %x", raw)
	}
	select {
	case err := <-failures:
		if !errors.Is(err, message.ErrInvalidLength) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("invalid triggered response not reported")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.transactions) != 0 {
		t.Fatal("invalid triggered response did not stop retransmissions")
	}
}

func TestDuplicateCachePreservesUnansweredAtCapacity(t *testing.T) {
	c, _ := testConnection(t)
	original := message.NewDeleteSessionRequest(0, 1)
	for seq := uint32(1); seq <= maxReplayRequestsPerPeer; seq++ {
		if c.isDuplicateOrUnmatched(testPeer(), message.NewDeleteSessionRequest(0, seq)) {
			t.Fatal("cache filled early")
		}
	}
	if !c.isDuplicateOrUnmatched(testPeer(), message.NewDeleteSessionRequest(0, 10000)) {
		t.Fatal("unanswered entry evicted")
	}
	if !c.isDuplicateOrUnmatched(testPeer(), original) {
		t.Fatal("unanswered request delivered twice")
	}
}

func TestReplayCapacityDoesNotCompleteCommand(t *testing.T) {
	c, _ := testConnection(t)
	seq, err := c.SendMessageTo(message.NewDeleteBearerCommand(0, 0), testPeer())
	if err != nil {
		t.Fatal(err)
	}
	for j := uint32(0); j < maxReplayRequestsPerPeer; j++ {
		if c.isDuplicateOrUnmatched(testPeer(), message.NewDeleteSessionRequest(0, j)) {
			t.Fatal("cache filled early")
		}
	}
	req := message.NewDeleteBearerRequest(0, seq, ie.NewEPSBearerID(5))
	if !c.isDuplicateOrUnmatched(testPeer(), req) {
		t.Fatal("request admitted without replay capacity")
	}
	c.mu.Lock()
	pending := len(c.transactions)
	c.mu.Unlock()
	if pending != 1 {
		t.Fatal("command completed before its triggered request could be admitted")
	}
}
