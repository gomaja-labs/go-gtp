package gtpv2

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

// DDN Failure Indication is an Initial message without a reply (TS 29.274 §4.2.5),
// so a node must hand it to the application.
func TestConnDDNFailureIndicationDelivered(t *testing.T) {
	got := make(chan struct{}, 1)
	c := udpServer(t, func(c *Conn) {
		c.AddHandler(message.MsgTypeDownlinkDataNotificationFailureIndication, func(*Conn, net.Addr, message.Message) error {
			got <- struct{}{}
			return nil
		})
	})
	peer := udpPeer(t)
	writeUDP(t, peer, c.LocalAddr(), message.NewDownlinkDataNotificationFailureIndication(0, 0x55, ie.NewCause(CauseUnableToPageUE, 0, 0, 0, nil)))
	select {
	case <-got:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Downlink Data Notification Failure Indication (initial message) was discarded as an unmatched reply")
	}
}

// AddHandler documents raw uint8 types. A handler for an unmodeled type
// (Bearer Resource Command, 68) should run.

// A valid TEID-C used by a new MME (different IP, e.g. MME relocation or the
// IPv4->IPv6 switch of §4.2.2.1 EXAMPLE 2) must not get Context not found.
func TestConnTEIDFromNewPeer(t *testing.T) {
	got := make(chan struct{}, 1)
	c := udpServer(t, func(c *Conn) {
		c.AddHandler(message.MsgTypeModifyBearerRequest, func(*Conn, net.Addr, message.Message) error { got <- struct{}{}; return nil })
	})
	old := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 2123}
	sess := NewSession(old, &Subscriber{IMSI: "001010000000001"})
	c.RegisterSession(0x1234, sess)
	newMME, err := net.ListenPacket("udp", "127.0.0.2:0")
	if err != nil {
		t.Skip(err)
	}
	defer func() { _ = newMME.Close() }()
	writeUDP(t, newMME, c.LocalAddr(), message.NewModifyBearerRequest(0x1234, 7))
	select {
	case <-got:
		return
	case <-time.After(300 * time.Millisecond):
	}
	raw, _ := readUDP(t, newMME)
	m := parseUDP(t, raw)
	t.Fatalf("MBR with a registered TEID from a new peer IP rejected: got type %d %s", m.MessageType(), message.Prettify(m))
}

// The replay cache is shared by every peer and admits at most 4096 requests per
// T3*(N3+2) window; once full, every new request from every peer is dropped.
func TestConnReplayCacheStarvation(t *testing.T) {
	c := udpServer(t, nil)
	attacker := udpPeer(t)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 2048)
		_ = attacker.SetReadDeadline(time.Now().Add(3 * time.Second))
		for {
			if _, _, err := attacker.ReadFrom(buf); err != nil {
				return
			}
		}
	}()
	for i := 1; i <= 4096; i++ {
		raw := mustMarshal(t, message.NewEchoRequest(uint32(i), ie.NewRecovery(1)))
		if _, err := attacker.WriteTo(raw, c.LocalAddr()); err != nil {
			t.Fatal(err)
		}
		if i%256 == 0 {
			time.Sleep(5 * time.Millisecond)
		}
	}
	time.Sleep(200 * time.Millisecond)
	victim := udpPeer(t)
	writeUDP(t, victim, c.LocalAddr(), message.NewEchoRequest(1, ie.NewRecovery(1)))
	if err := victim.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 2048)
	if _, _, err := victim.ReadFrom(buf); err != nil {
		t.Fatalf("another peer's Echo Request went unanswered after 4096 requests from one peer: %v", err)
	}
	wg.Wait()
}

// A request sent with WriteTo (not SendMessageTo) has no transaction, so its
// response is silently discarded before any handler runs.

// A response that arrives after WaitMessage timed out stays queued forever;
// after 1000 such messages every PassMessageTo blocks the serving goroutine.
func TestConnWaitMessageLateLeak(t *testing.T) {
	s := NewSession(&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 2123}, &Subscriber{IMSI: "001010000000001"})
	for i := 0; i < 1000; i++ {
		if _, err := s.WaitMessage(uint32(i), time.Microsecond); err == nil {
			t.Fatal("unexpected message")
		}
		if err := PassMessageTo(s, message.NewEchoResponse(uint32(i), ie.NewRecovery(1))); err != nil {
			t.Fatal(err)
		}
	}
	start := time.Now()
	err := PassMessageTo(s, message.NewEchoResponse(5000, ie.NewRecovery(1)))
	if err != nil {
		t.Fatalf("after 1000 late responses, PassMessageTo blocked %v and failed: %v", time.Since(start).Round(time.Millisecond), err)
	}
}

// Dial ties the Conn lifetime to the dial context. A Dial with a timeout
// context closes the returned Conn when the caller cancels it.
func TestConnDialContextLifetime(t *testing.T) {
	srv := udpServer(t, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	c, err := Dial(ctx, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, srv.LocalAddr(), IFTypeS11MMEGTPC, 0)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	time.Sleep(50 * time.Millisecond)
	if _, err := c.EchoRequest(srv.LocalAddr()); err != nil {
		t.Fatalf("Conn returned by Dial unusable after the dial context was cancelled: %v", err)
	}
	_ = c.Close()
}

// TS 29.274 §5.3: the VNSI Sequence Number "may be set to any number and shall
// be ignored by the receiver". A VNSI with another sequence is dropped here.
func TestConnVNSISequenceIgnored(t *testing.T) {
	got := make(chan struct{}, 1)
	c := udpServer(t, func(c *Conn) {
		c.AddHandler(message.MsgTypeVersionNotSupportedIndication, func(*Conn, net.Addr, message.Message) error { got <- struct{}{}; return nil })
	})
	peer := udpPeer(t)
	if _, err := c.EchoRequest(peer.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	readUDP(t, peer)
	writeUDP(t, peer, c.LocalAddr(), message.NewVersionNotSupportedIndication(0, 0))
	select {
	case <-got:
	case <-time.After(300 * time.Millisecond):
		t.Fatal("VNSI with a sequence other than the request's was discarded as an unmatched reply")
	}
}

// TS 29.274 §7.7.7: a Request with an IE whose Length exceeds the message must
// be answered with Cause "Invalid length" and the offending IE.
func TestConnIELengthErrorAnswered(t *testing.T) {
	c := udpServer(t, nil)
	peer := udpPeer(t)
	raw := mustMarshal(t, message.NewModifyBearerRequest(0, 3, ie.NewRATType(6)))
	raw[13] = 0x20 // RAT Type IE length 32 > remaining
	if _, err := peer.WriteTo(raw, c.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	if err := peer.SetReadDeadline(time.Now().Add(300 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 2048)
	if _, _, err := peer.ReadFrom(buf); err != nil {
		t.Fatalf("request with an invalid-length IE silently discarded: %v", err)
	}
}

// RespondTo leaves the piggybacked initial message's sequence to the caller
// (TS 29.274 §7.6 says the sender assigns it). With the constructor default,
// a second attach to the same peer fails, yet its unsent response is cached.
func TestConnPiggybackSequenceCollision(t *testing.T) {
	c := udpServer(t, nil)
	peer := udpPeer(t)
	for i, seq := range []uint32{100, 101} {
		csr := message.NewCreateSessionRequest(0, seq)
		rsp := message.NewCreateSessionResponse(0x10, 0, ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil), ie.NewBearerContext(ie.NewEPSBearerID(5), ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil)))
		rsp.Piggybacked = message.NewCreateBearerRequest(0x10, 0)
		if err := c.RespondTo(peer.LocalAddr(), csr, rsp); err != nil {
			t.Fatalf("attach %d: RespondTo with piggybacked CBR failed: %v", i+1, err)
		}
	}
}

func TestConnRawTypeHandlerValidMessage(t *testing.T) {
	got := make(chan struct{}, 1)
	c := udpServer(t, func(c *Conn) {
		c.AddHandler(68, func(*Conn, net.Addr, message.Message) error { got <- struct{}{}; return nil })
	})
	peer := udpPeer(t)
	// Table 7.2.5-1 mandatory IEs: Linked EBI, PTI, TAD ("delete existing TFT", no filters).
	writeUDP(t, peer, c.LocalAddr(), message.NewGeneric(68, 0, 0x800001,
		ie.NewEPSBearerID(5), ie.NewProcedureTransactionID(1), ie.New(ie.TrafficAggregateDescription, 0, []byte{0x40})))
	select {
	case <-got:
	case <-time.After(300 * time.Millisecond):
		t.Fatal("handler for raw type 68 not invoked for a valid Bearer Resource Command")
	}
}

func TestConnResponseToSendRawToDelivered(t *testing.T) {
	got := make(chan struct{}, 1)
	c := udpServer(t, func(c *Conn) {
		c.AddHandler(message.MsgTypeEchoResponse, func(*Conn, net.Addr, message.Message) error { got <- struct{}{}; return nil })
	})
	peer := udpPeer(t)
	if err := c.SendRawTo(mustMarshal(t, message.NewEchoRequest(9, ie.NewRecovery(1))), peer.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	raw, _ := readUDP(t, peer)
	req := parseUDP(t, raw)
	writeUDP(t, peer, c.LocalAddr(), message.NewEchoResponse(req.Sequence(), ie.NewRecovery(1)))
	select {
	case <-got:
	case <-time.After(300 * time.Millisecond):
		t.Fatal("response to a request sent with SendRawTo never reached its handler")
	}
}
