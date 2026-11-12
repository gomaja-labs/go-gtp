package gtpv2

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

func udpPeer(t *testing.T) net.PacketConn {
	t.Helper()
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = pc.Close() })
	return pc
}
func udpServer(t *testing.T, configure func(*Conn)) *Conn {
	t.Helper()
	c := NewConn(&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, IFTypeS11MMEGTPC, 1)
	if err := c.Listen(context.Background()); err != nil {
		t.Fatal(err)
	}
	if configure != nil {
		configure(c)
	}
	done := make(chan error, 1)
	go func() { done <- c.Serve(context.Background()) }()
	t.Cleanup(func() {
		_ = c.Close()
		select {
		case err := <-done:
			if err != nil {
				t.Error(err)
			}
		case <-time.After(time.Second):
			t.Error("Serve failed to exit")
		}
	})
	return c
}
func writeUDP(t *testing.T, pc net.PacketConn, addr net.Addr, msg message.Message) []byte {
	t.Helper()
	raw := mustMarshal(t, msg)
	if _, err := pc.WriteTo(raw, addr); err != nil {
		t.Fatal(err)
	}
	return raw
}
func readUDP(t *testing.T, pc net.PacketConn) ([]byte, net.Addr) {
	t.Helper()
	if err := pc.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	raw := make([]byte, 65535)
	n, peer, err := pc.ReadFrom(raw)
	if err != nil {
		t.Fatal(err)
	}
	return raw[:n], peer
}
func parseUDP(t *testing.T, raw []byte) message.Message {
	t.Helper()
	m, err := message.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return m
}
func noUDP(t *testing.T, pc net.PacketConn) {
	t.Helper()
	if err := pc.SetReadDeadline(time.Now().Add(35 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	raw := make([]byte, 65535)
	n, _, err := pc.ReadFrom(raw)
	var ne net.Error
	if !errors.As(err, &ne) || !ne.Timeout() {
		t.Fatalf("expected silence, received %x (%v)", raw[:n], err)
	}
}
func TestReceiveUDPFramingAndErrors(t *testing.T) {
	handled := make(chan uint32, 64)
	c := udpServer(t, func(c *Conn) {
		c.AddHandler(message.MsgTypeStopPagingIndication, func(_ *Conn, _ net.Addr, m message.Message) error { handled <- m.Sequence(); return nil })
	})
	peer := udpPeer(t)
	large := message.NewStopPagingIndication(0, 1, ie.New(250, 0, make([]byte, 60000)))
	writeUDP(t, peer, c.LocalAddr(), large)
	writeUDP(t, peer, c.LocalAddr(), message.NewStopPagingIndication(0, 2))
	for _, want := range []uint32{1, 2} {
		select {
		case got := <-handled:
			if got != want {
				t.Fatalf("order %d want %d", got, want)
			}
		case <-time.After(time.Second):
			t.Fatal("large datagram or following message lost")
		}
	}
	for _, version := range []byte{0, 1, 3, 7} {
		raw := mustMarshal(t, message.NewEchoRequest(0x123456, ie.NewRecovery(1)))
		raw[0] = version << 5
		if _, err := peer.WriteTo(raw, c.LocalAddr()); err != nil {
			t.Fatal(err)
		}
		if version < 2 {
			noUDP(t, peer)
			continue
		}
		reply, _ := readUDP(t, peer)
		m := parseUDP(t, reply)
		if m.MessageType() != message.MsgTypeVersionNotSupportedIndication || m.Sequence() != 0x123456 {
			t.Fatalf("wrong version reply %v", m)
		}
	}
	t.Run("unknown TEID", func(t *testing.T) {
		writeUDP(t, peer, c.LocalAddr(), message.NewDeleteSessionRequest(123, 15))
		raw, _ := readUDP(t, peer)
		m := parseUDP(t, raw)
		if m.MessageType() != message.MsgTypeDeleteSessionResponse || m.TEID() != 0 || m.Sequence() != 15 || m.(*message.DeleteSessionResponse).Cause.MustCause() != CauseContextNotFound {
			t.Fatalf("wrong context reply: %v", m)
		}
	})
	t.Run("invalid request length", func(t *testing.T) {
		raw := mustMarshal(t, message.NewDeleteSessionRequest(0, 16))
		binary.BigEndian.PutUint16(raw[2:4], 50)
		if _, err := peer.WriteTo(raw, c.LocalAddr()); err != nil {
			t.Fatal(err)
		}
		reply, _ := readUDP(t, peer)
		m := parseUDP(t, reply).(*message.DeleteSessionResponse)
		if m.Sequence() != 16 || m.Cause.MustCause() != CauseInvalidLength {
			t.Fatalf("wrong length reply: %v", m)
		}
	})
	t.Run("short unknown and invalid response discarded", func(t *testing.T) {
		for _, raw := range [][]byte{{0x48, 36}, mustMarshal(t, message.NewGeneric(254, 0, 3)), {0x48, 37, 0, 50, 0, 0, 0, 0, 0, 0, 16, 0}} {
			if _, err := peer.WriteTo(raw, c.LocalAddr()); err != nil {
				t.Fatal(err)
			}
			noUDP(t, peer)
		}
	})
}

func TestDialMatchesEchoAndCancellation(t *testing.T) {
	peer := udpPeer(t)
	attacker := udpPeer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan *Conn, 1)
	errs := make(chan error, 1)
	go func() {
		c, err := Dial(ctx, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, peer.LocalAddr(), IFTypeS11MMEGTPC, 1)
		if err != nil {
			errs <- err
			return
		}
		result <- c
	}()
	raw, client := readUDP(t, peer)
	seq := parseUDP(t, raw).Sequence()
	writeUDP(t, attacker, client, message.NewEchoResponse(seq, ie.NewRecovery(1)))
	writeUDP(t, peer, client, message.NewEchoResponse(seq+1, ie.NewRecovery(1)))
	writeUDP(t, peer, client, message.NewEchoRequest(seq, ie.NewRecovery(1)))
	writeUDP(t, peer, client, message.NewEchoResponse(seq))
	select {
	case c := <-result:
		_ = c.Close()
		t.Fatal("Dial accepted invalid handshake")
	case err := <-errs:
		var missing *RequiredIEMissingError
		if !errors.As(err, &missing) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("invalid Echo response not reported")
	}
	rebound, err := net.ListenPacket(client.Network(), client.String())
	if err != nil {
		t.Fatal(err)
	}
	_ = rebound.Close()
	t.Run("valid", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		done := make(chan error, 1)
		go func() {
			c, err := Dial(ctx, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, peer.LocalAddr(), IFTypeS11MMEGTPC, 1)
			if c != nil {
				_ = c.Close()
			}
			done <- err
		}()
		raw, addr := readUDP(t, peer)
		seq := parseUDP(t, raw).Sequence()
		writeUDP(t, attacker, addr, message.NewEchoResponse(seq, ie.NewRecovery(1)))
		writeUDP(t, peer, addr, message.NewEchoResponse(seq+1, ie.NewRecovery(1)))
		select {
		case err := <-done:
			t.Fatalf("Dial completed on unrelated datagram: %v", err)
		case <-time.After(30 * time.Millisecond):
		}
		writeUDP(t, peer, addr, message.NewEchoResponse(seq, ie.NewRecovery(1)))
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(time.Second):
			t.Fatal("matching Echo not accepted")
		}
	})
	t.Run("cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		done := make(chan error, 1)
		go func() {
			c, err := Dial(ctx, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, peer.LocalAddr(), IFTypeS11MMEGTPC, 1)
			if c != nil {
				_ = c.Close()
			}
			done <- err
		}()
		_, addr := readUDP(t, peer)
		cancel()
		select {
		case err := <-done:
			if !errors.Is(err, context.Canceled) {
				t.Fatal(err)
			}
		case <-time.After(300 * time.Millisecond):
			t.Fatal("Dial ignored cancellation")
		}
		pc, err := net.ListenPacket(addr.Network(), addr.String())
		if err != nil {
			t.Fatal(err)
		}
		_ = pc.Close()
	})
}

func TestReliableDeliveryUDP(t *testing.T) {
	failures := make(chan error, 4)
	handled := make(chan uint32, 4)
	c := udpServer(t, func(c *Conn) {
		if err := c.SetRetransmission(25*time.Millisecond, 2); err != nil {
			t.Fatal(err)
		}
		c.SetErrorHandler(func(err error) { failures <- err })
		c.AddHandler(message.MsgTypeDeleteSessionResponse, func(_ *Conn, _ net.Addr, m message.Message) error { handled <- m.Sequence(); return nil })
	})
	peer := udpPeer(t)
	attacker := udpPeer(t)
	seq, err := c.SendMessageTo(message.NewDeleteSessionRequest(0, 0), peer.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}
	original, _ := readUDP(t, peer)
	retried, _ := readUDP(t, peer)
	if !bytes.Equal(original, retried) {
		t.Fatal("retry changed bytes")
	}
	writeUDP(t, attacker, c.LocalAddr(), message.NewDeleteSessionResponse(0, seq, ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil)))
	writeUDP(t, peer, c.LocalAddr(), message.NewModifyBearerResponse(0, seq, ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil)))
	retried, _ = readUDP(t, peer)
	if !bytes.Equal(original, retried) {
		t.Fatal("retry changed bytes")
	}
	select {
	case err := <-failures:
		if !errors.Is(err, ErrTimeout) {
			t.Fatal(err)
		}
		var tx *TransactionError
		if !errors.As(err, &tx) || tx.Sequence != seq {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("missing transfer failure")
	}
	seq, err = c.SendMessageTo(message.NewDeleteSessionRequest(0, 0), peer.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}
	_, _ = readUDP(t, peer)
	reply := message.NewDeleteSessionResponse(0, seq, ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil))
	writeUDP(t, peer, c.LocalAddr(), reply)
	writeUDP(t, peer, c.LocalAddr(), reply)
	select {
	case got := <-handled:
		if got != seq {
			t.Fatal(got)
		}
	case <-time.After(time.Second):
		t.Fatal("correlated reply lost")
	}
	select {
	case <-handled:
		t.Fatal("duplicate reply delivered")
	case <-time.After(60 * time.Millisecond):
	}
	noUDP(t, peer)
}

func TestDuplicateRequestReplaysExactResponse(t *testing.T) {
	calls := make(chan struct{}, 2)
	c := udpServer(t, func(c *Conn) {
		c.AddHandler(message.MsgTypeDeleteSessionRequest, func(c *Conn, peer net.Addr, m message.Message) error {
			calls <- struct{}{}
			return c.RespondTo(peer, m, message.NewDeleteSessionResponse(0, 0, ie.NewCause(CauseRequestAccepted, 0, 0, 0, nil)))
		})
	})
	peer := udpPeer(t)
	req := message.NewDeleteSessionRequest(0, 55)
	writeUDP(t, peer, c.LocalAddr(), req)
	first, _ := readUDP(t, peer)
	writeUDP(t, peer, c.LocalAddr(), req)
	second, _ := readUDP(t, peer)
	if !bytes.Equal(first, second) {
		t.Fatal("replayed response differs")
	}
	if len(calls) != 1 {
		t.Fatalf("handler called %d times", len(calls))
	}
}

func TestWaitMessageConcurrentTransactions(t *testing.T) {
	s := NewSession(nil, nil)
	var wg sync.WaitGroup
	for seq := uint32(0); seq < 50; seq++ {
		wg.Go(func() {
			msg, err := s.WaitMessage(seq, time.Second)
			if err != nil || msg.Sequence() != seq {
				t.Errorf("sequence %d: %v %v", seq, msg, err)
			}
		})
	}
	for seq := uint32(49); ; seq-- {
		if err := PassMessageTo(s, message.NewEchoResponse(seq)); err != nil {
			t.Fatal(err)
		}
		if seq == 0 {
			break
		}
	}
	wg.Wait()
}

func TestSessionPublicationAndBearerIdentity(t *testing.T) {
	c, pc := testConnection(t)
	pc.hook = func([]byte) error { return errors.New("send failure") }
	_, _, err := c.CreateSession(testPeer(), ie.NewIMSI("001010000000001"), ie.NewFullyQualifiedTEID(10, 100, "192.0.2.1", ""))
	if err == nil {
		t.Fatal("send unexpectedly succeeded")
	}
	if len(c.Sessions()) != 0 {
		t.Fatal("failed CreateSession published a session")
	}
	a := testSession("001010000000001")
	b := testSession(a.IMSI)
	c.RegisterSession(10, a)
	c.RegisterSession(11, b)
	if _, err = c.GetSessionByIMSI(a.IMSI); !errors.Is(err, ErrAmbiguousIMSI) {
		t.Fatalf("ambiguous lookup: %v", err)
	}
	if len(c.GetSessionsByIMSI(a.IMSI)) != 2 {
		t.Fatal("PDN connection lost")
	}
	c.RegisterSession(12, a)
	if _, err = c.GetSessionByTEID(10); err == nil {
		t.Fatal("remapped TEID remains live")
	}
	c.RemoveSessionByIMSI(a.IMSI)
	if len(c.Sessions()) != 0 {
		t.Fatal("PDN connection not removed")
	}
	unknown := ie.New(250, 0, []byte{1, 2})
	s, err := c.ParseCreateSession(testPeer(), ie.NewBearerContext(ie.NewEPSBearerID(5), ie.NewFullyQualifiedTEID(1, 101, "192.0.2.1", ""), unknown),
		ie.NewBearerContext(ie.NewEPSBearerID(6), ie.NewFullyQualifiedTEID(1, 102, "192.0.2.2", "")), ie.NewAccessPointName("internet"), unknown)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []uint8{5, 6} {
		br, err := s.LookupBearerByEBI(id)
		if err != nil {
			t.Fatal(err)
		}
		teid, err := br.GetTEID(1)
		if err != nil || teid != uint32(96+id) || br.APN != "internet" {
			t.Fatalf("bearer %d: TEID %d APN %s err %v", id, teid, br.APN, err)
		}
	}
	if len(s.AdditionalIEs) != 1 || len(s.GetDefaultBearer().AdditionalIEs) != 1 {
		t.Fatal("unhandled IEs lost")
	}
}

func TestNilAccessorsAndLifecycle(t *testing.T) {
	var c *Conn
	var s *Session
	var b *Bearer
	if c.LocalAddr() != nil || c.SessionCount() != 0 || c.BearerCount() != 0 || c.Sessions() != nil {
		t.Fatal("nil Conn accessors")
	}
	if _, err := c.GetSessionByTEID(1); err == nil {
		t.Fatal("nil session lookup")
	}
	if s.PeerAddr() != nil || s.IsActive() || s.GetDefaultBearer() != nil || s.BearerCount() != 0 {
		t.Fatal("nil Session accessors")
	}
	if _, err := s.GetTEID(1); err == nil {
		t.Fatal("nil TEID lookup")
	}
	if _, err := s.WaitMessage(1, time.Millisecond); err == nil {
		t.Fatal("nil waiter")
	}
	if b.RemoteAddress() != nil || b.IncomingTEID() != 0 || b.OutgoingTEID() != 0 {
		t.Fatal("nil Bearer accessors")
	}
	if _, err := b.GetTEID(1); err == nil {
		t.Fatal("nil bearer TEID")
	}
	n := NewConn(&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, 1, 0)
	if err := n.Close(); err != nil {
		t.Fatal(err)
	}
	if err := n.Listen(context.Background()); !errors.Is(err, net.ErrClosed) {
		t.Fatal(err)
	}
	if _, _, err := n.ReadFrom(nil); !errors.Is(err, net.ErrClosed) {
		t.Fatal(err)
	}
}

func TestRetainedMessageOwnsDatagram(t *testing.T) {
	received := make(chan message.Message, 2)
	c := udpServer(t, func(c *Conn) {
		c.AddHandler(message.MsgTypeStopPagingIndication, func(_ *Conn, _ net.Addr, m message.Message) error { received <- m; return nil })
	})
	peer := udpPeer(t)
	writeUDP(t, peer, c.LocalAddr(), message.NewStopPagingIndication(0, 1, ie.New(250, 0, []byte{1})))
	writeUDP(t, peer, c.LocalAddr(), message.NewStopPagingIndication(0, 2, ie.New(250, 0, []byte{2})))
	var messages []message.Message
	for range 2 {
		select {
		case m := <-received:
			messages = append(messages, m)
		case <-time.After(time.Second):
			t.Fatal("message lost")
		}
	}
	first := messages[0].(*message.StopPagingIndication)
	if first.AdditionalIEs[0].Payload[0] != 1 {
		t.Fatalf("retained message overwritten by later datagram: %x", first.AdditionalIEs[0].Payload)
	}
}
