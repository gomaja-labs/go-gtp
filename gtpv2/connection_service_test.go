package gtpv2

import (
	"bytes"
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

func TestDDNFailureIndicationIsInitial(t *testing.T) {
	c, _ := testConnection(t)
	called := false
	c.AddHandler(message.MsgTypeDownlinkDataNotificationFailureIndication, func(_ *Conn, _ net.Addr, _ message.Message) error { called = true; return nil })
	c.processDatagram(testPeer(), mustMarshal(t, message.NewDownlinkDataNotificationFailureIndication(0, 7, ie.NewCause(64, 0, 0, 0, nil))))
	if !called {
		t.Fatal("initial indication discarded as unsolicited reply")
	}
}

func TestReplayFloodDoesNotStarveOtherPeers(t *testing.T) {
	c, _ := testConnection(t)
	for seq := uint32(1); seq <= 5000; seq++ {
		_ = c.isDuplicateOrUnmatched(testPeer(), message.NewDeleteSessionRequest(0, seq))
	}
	other := &net.UDPAddr{IP: net.IPv4(192, 0, 2, 3), Port: 2123}
	if c.isDuplicateOrUnmatched(other, message.NewDeleteSessionRequest(0, 1)) {
		t.Fatal("one peer exhausted replay capacity for another peer")
	}
}

func TestRegisteredRawMessageTypeDelivered(t *testing.T) {
	c, _ := testConnection(t)
	called := false
	c.AddHandler(254, func(_ *Conn, _ net.Addr, _ message.Message) error { called = true; return nil })
	c.processDatagram(testPeer(), mustMarshal(t, message.NewGeneric(254, 0, 2)))
	if !called {
		t.Fatal("registered raw type discarded")
	}
}

func TestBearerResourceCommandUsesCommandSequence(t *testing.T) {
	c, _ := testConnection(t)
	c.DisableValidation()
	seq, err := c.SendMessageTo(message.NewGeneric(message.MsgTypeBearerResourceCommand, 0, 0), testPeer())
	if err != nil {
		t.Fatal(err)
	}
	if seq&0x800000 == 0 {
		t.Fatal("command MSB not set")
	}
	called := false
	c.AddHandler(message.MsgTypeBearerResourceFailureIndication, func(_ *Conn, _ net.Addr, _ message.Message) error { called = true; return nil })
	c.processDatagram(testPeer(), mustMarshal(t, message.NewGeneric(message.MsgTypeBearerResourceFailureIndication, 0, seq)))
	if !called || len(c.transactions) != 0 {
		t.Fatal("failure indication did not complete command")
	}
}

func TestLocalTEIDAcceptsRelocatedPeer(t *testing.T) {
	c, _ := testConnection(t)
	c.RegisterSession(123, testSession("001010000000001"))
	peer := &net.UDPAddr{IP: net.ParseIP("2001:db8::1"), Port: 2123}
	if err := c.validate(peer, message.NewModifyBearerRequest(123, 2)); err != nil {
		t.Fatal(err)
	}
}

func TestPiggybackInitialSequencesAreAllocated(t *testing.T) {
	c, _ := testConnection(t)
	var previous uint32
	for seq := uint32(1); seq <= 2; seq++ {
		response := message.NewCreateSessionResponse(0, 0, ie.NewCause(16, 0, 0, 0, nil))
		response.Piggybacked = message.NewCreateBearerRequest(0, 0)
		if err := c.RespondTo(testPeer(), message.NewCreateSessionRequest(0, seq), response); err != nil {
			t.Fatal(err)
		}
		got := response.Piggybacked.Sequence()
		if got == 0 || got == previous {
			t.Fatalf("piggyback sequence reused: %d", got)
		}
		previous = got
	}
}

func TestFailedResponseIsNotReplayed(t *testing.T) {
	c, p := testConnection(t)
	req := message.NewDeleteSessionRequest(0, 7)
	_ = c.isDuplicateOrUnmatched(testPeer(), req)
	p.hook = func([]byte) error { return errors.New("write failed") }
	if err := c.RespondTo(testPeer(), req, message.NewDeleteSessionResponse(0, 0, ie.NewCause(16, 0, 0, 0, nil))); err == nil {
		t.Fatal("write succeeded")
	}
	p.hook = nil
	_ = c.isDuplicateOrUnmatched(testPeer(), req)
	if len(p.writes) != 0 {
		t.Fatal("failed response replayed")
	}
}

func TestSessionDropsLateReplies(t *testing.T) {
	s := NewSession(nil, nil)
	for seq := uint32(0); seq < 2*maxPendingMessages; seq++ {
		if _, err := s.WaitMessage(seq, 0); !errors.Is(err, ErrTimeout) {
			t.Fatal(err)
		}
		if err := PassMessageTo(s, message.NewEchoResponse(seq)); err != nil {
			t.Fatalf("late reply %d: %v", seq, err)
		}
	}
	if err := PassMessageTo(s, message.NewEchoResponse(5000)); err != nil {
		t.Fatalf("fresh reply after late ones: %v", err)
	}
	if msg, err := s.WaitMessage(5000, time.Second); err != nil || msg.Sequence() != 5000 {
		t.Fatalf("fresh reply not delivered: %v %v", msg, err)
	}
	s.mu.Lock()
	queued := len(s.pending)
	s.mu.Unlock()
	if queued != 0 {
		t.Fatalf("%d late replies are still queued", queued)
	}
}

func TestSessionKeepsReplyArrivingBeforeWait(t *testing.T) {
	s := NewSession(nil, nil)
	if err := PassMessageTo(s, message.NewEchoResponse(7)); err != nil {
		t.Fatal(err)
	}
	if msg, err := s.WaitMessage(7, time.Second); err != nil || msg.Sequence() != 7 {
		t.Fatalf("got %v %v", msg, err)
	}
}

func TestSessionEvictsOldestUnclaimedMessage(t *testing.T) {
	s := NewSession(nil, nil)
	for seq := uint32(1); seq <= maxPendingMessages+1; seq++ {
		if err := PassMessageTo(s, message.NewEchoResponse(seq)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := s.WaitMessage(1, 0); !errors.Is(err, ErrTimeout) {
		t.Fatalf("oldest message was not evicted: %v", err)
	}
	if msg, err := s.WaitMessage(maxPendingMessages+1, 0); err != nil || msg.Sequence() != maxPendingMessages+1 {
		t.Fatalf("newest message lost: %v %v", msg, err)
	}
}

func TestSessionExpiresUnclaimedMessages(t *testing.T) {
	old := pendingMessageTTL
	pendingMessageTTL = 10 * time.Millisecond
	t.Cleanup(func() { pendingMessageTTL = old })

	s := NewSession(nil, nil)
	if err := PassMessageTo(s, message.NewEchoResponse(1)); err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond) // let the unclaimed message expire
	if err := PassMessageTo(s, message.NewEchoResponse(2)); err != nil {
		t.Fatal(err)
	}
	if _, err := s.WaitMessage(1, 0); !errors.Is(err, ErrTimeout) {
		t.Fatalf("expired message delivered: %v", err)
	}
}

func TestVersionIndicationIgnoresSequence(t *testing.T) {
	c, _ := testConnection(t)
	if _, err := c.SendMessageTo(message.NewEchoRequest(0, ie.NewRecovery(1)), testPeer()); err != nil {
		t.Fatal(err)
	}
	if !c.matchReply(testPeer(), 777, message.MsgTypeVersionNotSupportedIndication) {
		t.Fatal("VNSI correlated by sequence")
	}
}

func TestInvalidIEGetsErrorResponse(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		wire                 []byte
		cause, typ, instance uint8
	}{
		{"framing", []byte{0x48, 36, 0, 13, 0, 0, 0, 0, 0, 0, 9, 0, 93, 0, 2, 15, 0xff}, 67, 93, 15},
		{"mandatory length", []byte{0x48, 95, 0, 13, 0, 0, 0, 0, 0, 0, 9, 0, 93, 0, 1, 0, 0xff}, 67, 93, 0},
		{"semantic", []byte{0x48, 32, 0, 14, 0, 0, 0, 0, 0, 0, 9, 0, 71, 0, 2, 0, 1, '_'}, 69, 71, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, p := testConnection(t)
			c.processDatagram(testPeer(), tc.wire)
			if len(p.writes) != 1 {
				t.Fatalf("responses: %d", len(p.writes))
			}
			m := parseUDP(t, p.writes[0])
			cause := messageIEs(m)[0]
			if m.Sequence() != 9 || cause.MustCause() != tc.cause || len(cause.Payload) != 6 || cause.Payload[2+0] != tc.typ || cause.Payload[5]&15 != tc.instance {
				t.Fatalf("wrong cause: %x", cause.Payload)
			}
		})
	}
}

func TestDialContextOnlyBoundsHandshake(t *testing.T) {
	peer := udpPeer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan *Conn, 1)
	errs := make(chan error, 1)
	go func() {
		c, err := Dial(ctx, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, peer.LocalAddr(), 1, 0)
		if err != nil {
			errs <- err
			return
		}
		result <- c
	}()
	raw, addr := readUDP(t, peer)
	writeUDP(t, peer, addr, message.NewEchoResponse(parseUDP(t, raw).Sequence(), ie.NewRecovery(1)))
	var c *Conn
	select {
	case c = <-result:
	case err := <-errs:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Dial did not finish")
	}
	defer func() { _ = c.Close() }()
	cancel()
	select {
	case <-c.closed():
		t.Fatal("Dial context closed returned connection")
	case <-time.After(20 * time.Millisecond):
	}
	writeUDP(t, peer, c.LocalAddr(), message.NewEchoRequest(99, ie.NewRecovery(1)))
	raw, _ = readUDP(t, peer)
	if parseUDP(t, raw).Sequence() != 99 {
		t.Fatal("connection no longer serving")
	}
}

func TestUnmatchedPrimaryDiscardsPiggyback(t *testing.T) {
	c, _ := testConnection(t)
	c.DisableValidation()
	calls := 0
	c.AddHandler(message.MsgTypeCreateBearerRequest, func(_ *Conn, _ net.Addr, _ message.Message) error { calls++; return nil })
	response := message.NewCreateSessionResponse(0, 12)
	response.Piggybacked = message.NewCreateBearerRequest(0, 13)
	c.processDatagram(testPeer(), mustMarshal(t, response))
	if calls != 0 {
		t.Fatal("unsolicited primary activated dedicated bearer")
	}
}

func TestRawSendCorrelatesReplies(t *testing.T) {
	c, _ := testConnection(t)
	raw := mustMarshal(t, message.NewDeleteSessionRequest(0, 456))
	if err := c.SendRawTo(raw, testPeer()); err != nil {
		t.Fatal(err)
	}
	calls := 0
	c.AddHandler(message.MsgTypeDeleteSessionResponse, func(_ *Conn, _ net.Addr, _ message.Message) error { calls++; return nil })
	reply := mustMarshal(t, message.NewDeleteSessionResponse(0, 456, ie.NewCause(16, 0, 0, 0, nil)))
	c.processDatagram(testPeer(), reply)
	c.processDatagram(testPeer(), reply)
	if calls != 1 || len(c.transactions) != 0 {
		t.Fatal("raw request reply was lost or duplicated")
	}
}

func TestAnsweredReplayCacheContinuesAtCapacity(t *testing.T) {
	c, _ := testConnection(t)
	for seq := uint32(1); seq <= 5000; seq++ {
		req := message.NewDeleteSessionRequest(0, seq)
		if c.isDuplicateOrUnmatched(testPeer(), req) {
			t.Fatalf("answered requests stopped admission at %d", seq)
		}
		if err := c.RespondTo(testPeer(), req, message.NewDeleteSessionResponse(0, 0, ie.NewCause(16, 0, 0, 0, nil))); err != nil {
			t.Fatal(err)
		}
	}
	if len(c.received) > maxReplayRequestsPerPeer {
		t.Fatal("peer bound exceeded")
	}
}

func TestReplayBoundIncludesChangingPortsAndPeers(t *testing.T) {
	c, _ := testConnection(t)
	for port := 1; port <= 5000; port++ {
		peer := &net.UDPAddr{IP: net.IPv4(192, 0, 2, 2), Port: port}
		_ = c.isDuplicateOrUnmatched(peer, message.NewDeleteSessionRequest(0, 1))
	}
	if len(c.received) > maxReplayRequestsPerPeer {
		t.Fatal("source port changes bypassed peer bound")
	}
	for n := 1; n <= 200; n++ {
		peer := &net.UDPAddr{IP: net.IPv4(198, 51, 100, byte(n)), Port: 2123}
		if c.isDuplicateOrUnmatched(peer, message.NewDeleteSessionRequest(0, 1)) {
			t.Fatal("new peer starved")
		}
	}
	if len(c.replayPeers) > maxReplayPeers || len(c.received) > maxReplayPeers*maxReplayRequestsPerPeer {
		t.Fatal("total replay bound exceeded")
	}
}

func TestReplayWindowIndependentOfLocalRetries(t *testing.T) {
	c, _ := testConnection(t)
	if err := c.SetRetransmission(time.Millisecond, 0); err != nil {
		t.Fatal(err)
	}
	req := message.NewDeleteSessionRequest(0, 1)
	_ = c.isDuplicateOrUnmatched(testPeer(), req)
	c.mu.Lock()
	c.replayPeerLocked(testPeer(), time.Now().Add(time.Second))
	retained := c.received[requestKey{testPeer().String(), 1, req.MessageType()}] != nil
	c.mu.Unlock()
	if !retained {
		t.Fatal("replay retention ended at local retry window")
	}
	if err := c.SetReplayWindow(time.Hour); err != nil {
		t.Fatal(err)
	}
	if err := c.SetReplayWindow(0); err == nil {
		t.Fatal("zero replay window accepted")
	}
}

func TestBadSecondaryPreservesPrimaryAndAnswersIEError(t *testing.T) {
	c, p := testConnection(t)
	c.DisableValidation()
	seq, err := c.SendMessageTo(message.NewCreateSessionRequest(0, 0), testPeer())
	if err != nil {
		t.Fatal(err)
	}
	called := false
	c.AddHandler(message.MsgTypeCreateSessionResponse, func(_ *Conn, _ net.Addr, _ message.Message) error { called = true; return nil })
	response := message.NewCreateSessionResponse(0, seq, ie.NewCause(16, 0, 0, 0, nil))
	primary := mustMarshal(t, response)
	primary[0] |= 0x10
	bad := []byte{0x48, 95, 0, 13, 0, 0, 0, 0, 0, 0, 9, 0, 93, 0, 2, 3, 0xff}
	c.processDatagram(testPeer(), append(primary, bad...))
	if !called {
		t.Fatal("secondary failure discarded primary")
	}
	if len(p.writes) != 2 {
		t.Fatalf("secondary error not answered: %d", len(p.writes))
	}
	reply := parseUDP(t, p.writes[1]).(*message.CreateBearerResponse)
	if reply.Sequence() != 9 || reply.Cause.MustCause() != 67 || reply.Cause.Payload[2] != 93 || reply.Cause.Payload[5]&15 != 3 {
		t.Fatalf("wrong secondary cause: %x", reply.Cause.Payload)
	}
}

func TestDuplicatePrimaryStillReplaysDedicatedResponse(t *testing.T) {
	c, p := testConnection(t)
	c.DisableValidation()
	seq, err := c.SendMessageTo(message.NewCreateSessionRequest(0, 0), testPeer())
	if err != nil {
		t.Fatal(err)
	}
	primaryCalls, secondaryCalls := 0, 0
	c.AddHandler(message.MsgTypeCreateSessionResponse, func(_ *Conn, _ net.Addr, _ message.Message) error { primaryCalls++; return nil })
	c.AddHandler(message.MsgTypeCreateBearerRequest, func(c *Conn, peer net.Addr, m message.Message) error {
		secondaryCalls++
		return c.RespondTo(peer, m, message.NewCreateBearerResponse(0, 0, ie.NewCause(16, 0, 0, 0, nil)))
	})
	response := message.NewCreateSessionResponse(0, seq, ie.NewCause(16, 0, 0, 0, nil))
	response.Piggybacked = message.NewCreateBearerRequest(0, 77)
	raw := mustMarshal(t, response)
	c.processDatagram(testPeer(), raw)
	c.processDatagram(testPeer(), raw)
	if primaryCalls != 1 || secondaryCalls != 1 || len(p.writes) != 3 || !bytes.Equal(p.writes[1], p.writes[2]) {
		t.Fatal("dedicated response not replayed independently")
	}
}

func TestReplayExpiryOrderedAcrossWindowChanges(t *testing.T) {
	c, _ := testConnection(t)
	if err := c.SetReplayWindow(time.Hour); err != nil {
		t.Fatal(err)
	}
	_ = c.isDuplicateOrUnmatched(testPeer(), message.NewDeleteSessionRequest(0, 1))
	if err := c.SetReplayWindow(time.Second); err != nil {
		t.Fatal(err)
	}
	_ = c.isDuplicateOrUnmatched(testPeer(), message.NewDeleteSessionRequest(0, 2))
	c.mu.Lock()
	c.replayPeerLocked(testPeer(), time.Now().Add(2*time.Second))
	first := c.received[requestKey{testPeer().String(), 1, message.MsgTypeDeleteSessionRequest}]
	second := c.received[requestKey{testPeer().String(), 2, message.MsgTypeDeleteSessionRequest}]
	c.mu.Unlock()
	if first == nil || second != nil {
		t.Fatal("replay expiration was not ordered")
	}
}

func TestInvalidIEReplyStopsRetriesAndReportsError(t *testing.T) {
	c, p := testConnection(t)
	if err := c.SendRawTo(mustMarshal(t, message.NewDeleteSessionRequest(0, 9)), testPeer()); err != nil {
		t.Fatal(err)
	}
	var reported error
	c.SetErrorHandler(func(err error) { reported = err })
	c.processDatagram(testPeer(), []byte{0x48, 37, 0, 12, 0, 0, 0, 0, 0, 0, 9, 0, 2, 0, 0, 0})
	var tx *TransactionError
	var content *message.IEError
	if !errors.As(reported, &tx) || !errors.As(reported, &content) || content.Type != ie.Cause || tx.Sequence != 9 || len(c.transactions) != 0 || len(p.writes) != 1 {
		t.Fatalf("reply error: %v", reported)
	}
}

func TestFailedPrimaryDiscardsDedicatedRequest(t *testing.T) {
	for _, scenario := range []string{"rejection", "handler error", "invalid TEID"} {
		t.Run(scenario, func(t *testing.T) {
			c, _ := testConnection(t)
			c.DisableValidation()
			seq, err := c.SendMessageTo(message.NewCreateSessionRequest(0, 0), testPeer())
			if err != nil {
				t.Fatal(err)
			}
			called := false
			c.SetErrorHandler(func(error) {})
			c.AddHandler(message.MsgTypeCreateBearerRequest, func(_ *Conn, _ net.Addr, _ message.Message) error { called = true; return nil })
			c.AddHandler(message.MsgTypeCreateSessionResponse, func(_ *Conn, _ net.Addr, _ message.Message) error {
				if scenario == "handler error" {
					return errors.New("establishment failed")
				}
				return nil
			})
			cause := uint8(16)
			teid := uint32(0)
			if scenario == "rejection" {
				cause = 64
			}
			if scenario == "invalid TEID" {
				c.EnableValidation()
				teid = 999
			}
			response := message.NewCreateSessionResponse(teid, seq, ie.NewCause(cause, 0, 0, 0, nil))
			response.Piggybacked = message.NewCreateBearerRequest(0, 88)
			c.processDatagram(testPeer(), mustMarshal(t, response))
			if called {
				t.Fatal("dedicated request handled after primary failure")
			}
		})
	}
}

func TestRawSendRejectsInvalidEncoding(t *testing.T) {
	c, _ := testConnection(t)
	for _, raw := range [][]byte{{}, {0x48, 36}, mustMarshal(t, message.NewGeneric(message.MsgTypeBearerResourceCommand, 0, 1))} {
		if err := c.SendRawTo(raw, testPeer()); err == nil {
			t.Fatalf("invalid raw request accepted: %x", raw)
		}
	}
	if len(c.transactions) != 0 {
		t.Fatal("invalid raw request registered a transaction")
	}
}
