package gtpv2

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

type testDatagram struct {
	b   []byte
	err error
}
type testPacketConn struct {
	mu     sync.Mutex
	input  chan testDatagram
	done   chan struct{}
	once   sync.Once
	writes [][]byte
	peers  []net.Addr
	hook   func([]byte) error
}

func newReviewPC() *testPacketConn {
	return &testPacketConn{input: make(chan testDatagram, 128), done: make(chan struct{})}
}
func (p *testPacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	select {
	case d := <-p.input:
		return copy(b, d.b), testPeer(), d.err
	case <-p.done:
		return 0, nil, net.ErrClosed
	}
}
func (p *testPacketConn) WriteTo(b []byte, a net.Addr) (int, error) {
	if p.hook != nil {
		if err := p.hook(b); err != nil {
			return 0, err
		}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.writes = append(p.writes, append([]byte(nil), b...))
	p.peers = append(p.peers, a)
	return len(b), nil
}
func (p *testPacketConn) Close() error { p.once.Do(func() { close(p.done) }); return nil }
func (p *testPacketConn) LocalAddr() net.Addr {
	return &net.UDPAddr{IP: net.IPv4(192, 0, 2, 1), Port: 2123}
}
func (p *testPacketConn) SetDeadline(time.Time) error      { return nil }
func (p *testPacketConn) SetReadDeadline(time.Time) error  { return nil }
func (p *testPacketConn) SetWriteDeadline(time.Time) error { return nil }
func testPeer() net.Addr                                   { return &net.UDPAddr{IP: net.IPv4(192, 0, 2, 2), Port: 43210} }
func testConnection(t testing.TB) (*Conn, *testPacketConn) {
	c := NewConn(&net.UDPAddr{}, 10, 1)
	p := newReviewPC()
	c.pktConn = p
	t.Cleanup(func() { _ = c.Close() })
	return c, p
}
func mustMarshal(t *testing.T, m message.Message) []byte {
	t.Helper()
	b, err := message.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func testRecover(t *testing.T) {
	t.Helper()
	if p := recover(); p != nil {
		t.Errorf("unexpected panic: %v", p)
	}
}
func testStart(t *testing.T, c *Conn, p *testPacketConn) chan error {
	t.Helper()
	end := make(chan error, 1)
	go func() { end <- c.Serve(context.Background()) }()
	t.Cleanup(func() {
		_ = p.Close()
		select {
		case err := <-end:
			if err != nil {
				t.Error(err)
			}
		case <-time.After(time.Second):
			t.Error("Serve did not exit")
		}
		select {
		case <-c.closed():
		default:
			_ = c.Close()
		}
	})
	return end
}
func TestConnLargeDatagram(t *testing.T) {
	c, p := testConnection(t)
	got := make(chan int, 1)
	c.AddHandler(message.MsgTypeStopPagingIndication, func(_ *Conn, _ net.Addr, m message.Message) error { got <- m.MarshalLen(); return nil })
	testStart(t, c, p)
	xs := make([]*ie.IE, 250)
	for j := range xs {
		xs[j] = ie.New(250, 0, []byte{1, 2, 3, 4})
	}
	m := message.NewStopPagingIndication(0, 1, xs...)
	b := mustMarshal(t, m)
	p.input <- testDatagram{b: b}
	select {
	case n := <-got:
		if n != len(b) {
			t.Errorf("%d-byte datagram silently truncated to %d bytes", len(b), n)
		}
	case <-time.After(time.Second):
		t.Error("large valid datagram not delivered")
	}
}
func TestConnUnboundedHandlers(t *testing.T) {
	c, p := testConnection(t)
	const n = 64
	entered := make(chan struct{}, n)
	release := make(chan struct{})
	c.AddHandler(message.MsgTypeStopPagingIndication, func(_ *Conn, _ net.Addr, _ message.Message) error { entered <- struct{}{}; <-release; return nil })
	testStart(t, c, p)
	defer close(release)
	b := mustMarshal(t, message.NewStopPagingIndication(0, 1))
	for j := 0; j < n; j++ {
		p.input <- testDatagram{b: b}
	}
	// All handlers block, so reaching n proves there is no per-peer/session serialization or worker limit at this load.
	for j := 0; j < n; j++ {
		select {
		case <-entered:
		case <-time.After(time.Second):
			return
		}
	}
	t.Errorf("all %d handlers entered concurrently while preceding handlers remained blocked", n)
}
func TestConnSequenceRollback(t *testing.T) {
	c, p := testConnection(t)
	firstEntered := make(chan struct{})
	release := make(chan struct{})
	firstDone := make(chan struct{})
	p.hook = func(b []byte) error {
		if b[6] == 1 {
			close(firstEntered)
			<-release
			return errors.New("synthetic send failure")
		}
		return nil
	}
	go func() {
		defer close(firstDone)
		_, _ = c.SendMessageTo(message.NewEchoRequest(0, ie.NewRecovery(1)), testPeer())
	}()
	<-firstEntered
	seq2, err := c.SendMessageTo(message.NewEchoRequest(0, ie.NewRecovery(1)), testPeer())
	if err != nil {
		t.Fatal(err)
	}
	close(release)
	<-firstDone
	seq3, err := c.SendMessageTo(message.NewEchoRequest(0, ie.NewRecovery(1)), testPeer())
	if err != nil {
		t.Fatal(err)
	}
	if seq2 == seq3 {
		t.Errorf("outstanding sequence reused after another sender fails: %d", seq3)
	}

}
func TestConnCommandSequence(t *testing.T) {
	c, _ := testConnection(t)
	for _, m := range []message.Message{message.NewModifyBearerCommand(0, 0), message.NewDeleteBearerCommand(0, 0)} {
		seq, err := c.SendMessageTo(m, testPeer())
		if err != nil {
			t.Fatal(err)
		}
		if seq&0x800000 == 0 {
			t.Errorf("%s sent with command sequence MSB=0 (%#x)", m.MessageTypeName(), seq)
		}
	}
}
func TestConnClose(t *testing.T) {
	t.Run("socket", func(t *testing.T) {
		c, p := testConnection(t)
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
		select {
		case <-p.done:
		default:
			t.Error("Close returned without closing underlying PacketConn when Serve is not running")
		}
	})
	t.Run("double", func(t *testing.T) { defer testRecover(t); c, _ := testConnection(t); _ = c.Close(); _ = c.Close() })
}
func TestConnServeReturnCleanup(t *testing.T) {
	c, p := testConnection(t)
	p.input <- testDatagram{err: io.EOF}
	if err := c.Serve(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	select {
	case <-p.done:
	default:
		t.Error("Serve returned but PacketConn and cancellation watcher remain open")
	}
}
func TestConnReservedTEID(t *testing.T) {
	defer testRecover(t)
	c, _ := testConnection(t)
	i := c.NewSenderFTEID("192.0.2.1", "")
	if i == nil {
		t.Fatal("no F-TEID")
	}
	teid, err := i.TEID()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetSessionByTEID(teid); err == nil {
		t.Error("unregistered reserved TEID accepted")
	}
}
func testSession(imsi string) *Session {
	return NewSession(testPeer(), &Subscriber{IMSI: imsi, Location: &Location{}})
}
func TestConnSessionRemoval(t *testing.T) {
	t.Run("reregister", func(t *testing.T) {
		c, _ := testConnection(t)
		s := testSession("001010000000001")
		c.RegisterSession(10, s)
		c.RegisterSession(11, s)
		c.RemoveSession(s)
		if _, err := c.GetSessionByTEID(10); err == nil {
			t.Error("old incoming TEID still routes to removed session")
		}
	})
	t.Run("replacement", func(t *testing.T) {
		c, _ := testConnection(t)
		a, b := testSession("001010000000001"), testSession("001010000000001")
		c.RegisterSession(10, a)
		c.RegisterSession(11, b)
		c.RemoveSession(a)
		if got, err := c.GetSessionByIMSI(b.IMSI); err != nil || got != b {
			t.Error("removing old session deleted replacement IMSI mapping")
		}
	})
}
func TestConnDownlinkGBR(t *testing.T) {
	c, _ := testConnection(t)
	s, err := c.ParseCreateSession(testPeer(), ie.NewBearerContext(ie.NewBearerQoS(1, 2, 1, 9, 100, 200, 300, 400)))
	if err != nil {
		t.Fatal(err)
	}
	if got := s.GetDefaultBearer().GBRDL; got != 400 {
		t.Errorf("downlink GBR=%d, want 400", got)
	}
}
func TestConnPartialRegistration(t *testing.T) {
	c, _ := testConnection(t)
	_, err := c.ParseCreateSession(testPeer(), ie.NewIMSI("001010000000001"), ie.NewFullyQualifiedTEID(10, 123, "192.0.2.1", ""), ie.New(ie.RATType, 0, nil))
	if err == nil {
		t.Fatal("malformed RAT accepted")
	}
	if _, err := c.GetSessionByTEID(123); err == nil {
		t.Error("ParseCreateSession failure left partially constructed session registered")
	}
}
func TestConnMultipleBearers(t *testing.T) {
	c, _ := testConnection(t)
	s, err := c.ParseCreateSession(testPeer(), ie.NewBearerContext(ie.NewEPSBearerID(5)), ie.NewBearerContext(ie.NewEPSBearerID(6)))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.LookupBearerByEBI(5); err != nil {
		t.Error("first bearer overwritten by second bearer context")
	}
	if s.BearerCount() != 2 {
		t.Errorf("bearer count=%d, want 2", s.BearerCount())
	}
}
func TestConnWaitMessageLosesUnmatched(t *testing.T) {
	s := testSession("001010000000001")
	_ = PassMessageTo(s, message.NewEchoResponse(2))
	_ = PassMessageTo(s, message.NewEchoResponse(1))
	if _, err := s.WaitMessage(1, time.Second); err != nil {
		t.Errorf("WaitMessage(1) failed instead of retaining unmatched sequence 2: %v", err)
	}
	if _, err := s.WaitMessage(2, time.Millisecond); err != nil {
		t.Errorf("unmatched response was destroyed: %v", err)
	}
}
func TestConnPeerAddressRace(t *testing.T) {
	c, _ := testConnection(t)
	s := testSession("001010000000001")
	c.RegisterSession(1, s)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			s.UpdatePeerAddr(testPeer())
		}
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			_ = s.PeerAddr()
			_, _ = c.GetSessionByTEID(1)
		}
	}()
	wg.Wait()
}
func TestConnBearerRace(t *testing.T) {
	b := NewBearer(5, "example", nil)
	s := testSession("001010000000001")
	s.SetDefaultBearer(b)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			b.SetIncomingTEID(uint32(j))
			b.SetOutgoingTEID(uint32(j))
			b.SetRemoteAddress(testPeer())
		}
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			_ = b.IncomingTEID()
			_ = b.OutgoingTEID()
			_ = b.RemoteAddress()
			_ = s.LookupEBIByTEID(uint32(j))
		}
	}()
	wg.Wait()
}
func TestConnValidationRace(t *testing.T) {
	c, _ := testConnection(t)
	m := message.NewEchoResponse(1)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			c.EnableValidation()
			c.DisableValidation()
		}
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			_ = c.handleMessage(testPeer(), m)
		}
	}()
	wg.Wait()
}
func TestConnUnknownTEIDReply(t *testing.T) {
	c, p := testConnection(t)
	_ = c.handleMessage(testPeer(), message.NewDeleteSessionRequest(123, 1, ie.NewEPSBearerID(5)))
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.writes) == 0 {
		t.Error("unknown TEID request discarded without Context Not Found response")
	}
}
func TestConnMandatoryIEs(t *testing.T) {
	c, p := testConnection(t)
	called := false
	c.AddHandler(message.MsgTypeCreateSessionRequest, func(_ *Conn, _ net.Addr, _ message.Message) error { called = true; return nil })
	_ = c.handleMessage(testPeer(), message.NewCreateSessionRequest(0, 1))
	if called {
		t.Error("request missing all mandatory IEs passed validation and reached handler")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.writes) == 0 {
		t.Error("no Mandatory IE Missing response")
	}
}
func TestConnTransactionCorrelation(t *testing.T) {
	c, _ := testConnection(t)
	calls := 0
	c.AddHandler(message.MsgTypeCreateSessionResponse, func(_ *Conn, _ net.Addr, _ message.Message) error { calls++; return nil })
	m := message.NewCreateSessionResponse(0, 123, ie.NewCause(CauseContextNotFound, 0, 0, 0, nil))
	_ = c.handleMessage(testPeer(), m)
	_ = c.handleMessage(testPeer(), m)
	if calls != 0 {
		t.Errorf("unsolicited/duplicate response reached handler %d times without any outstanding request", calls)
	}
}
func TestConnUnexpectedIEError(t *testing.T) {
	e := &UnexpectedIEError{IEType: 87}
	if got := e.Error(); !strings.Contains(got, "87") {
		t.Errorf("IE type value missing from error: %s", got)
	}
}
func TestConnSequenceRange(t *testing.T) {
	c, _ := testConnection(t)
	c.sequence = 0x7fffff
	if got := c.IncSequence(); got != 0 {
		t.Errorf("ordinary request sequence failed to wrap: %#x", got)
	}
}
func TestConnUnsupportedVersionReply(t *testing.T) {
	c, p := testConnection(t)
	response := make(chan uint8, 1)
	p.hook = func(b []byte) error { response <- b[1]; return nil }
	testStart(t, c, p)
	b := mustMarshal(t, message.NewEchoRequest(1, ie.NewRecovery(1)))
	b[0] = 0x60
	p.input <- testDatagram{b: b}
	select {
	case typ := <-response:
		if typ != message.MsgTypeVersionNotSupportedIndication {
			t.Errorf("unsupported version 3 answered with message type %d, want Version Not Supported (3)", typ)
		}
	case <-time.After(time.Second):
		t.Error("unsupported version 3 did not receive Version Not Supported")
	}
}
func TestConnMultiplePDNConnections(t *testing.T) {
	c, _ := testConnection(t)
	a, b := testSession("001010000000001"), testSession("001010000000001")
	_ = a.Activate()
	_ = b.Activate()
	c.RegisterSession(10, a)
	c.RegisterSession(11, b)
	if n := c.SessionCount(); n != 2 {
		t.Errorf("two PDN connections with same IMSI counted as %d sessions", n)
	}
	if n := len(c.Sessions()); n != 2 {
		t.Errorf("Sessions silently omits %d PDN connection(s)", 2-n)
	}
}
func TestConnEmptyBearerName(t *testing.T) {
	s := testSession("001010000000001")
	s.AddBearer("", NewBearer(6, "example", nil))
	s.RemoveBearerByEBI(6)
	if _, err := s.LookupBearerByName(""); err == nil {
		t.Error("RemoveBearerByEBI failed for accepted empty bearer name")
	}
}
func TestConnLocalAddrBeforeListen(t *testing.T) {
	defer testRecover(t)
	c := NewConn(&net.UDPAddr{}, 10, 0)
	_ = c.LocalAddr()
}

func TestConnListenFailure(t *testing.T) {
	c, p := testConnection(t)
	c.laddr = &net.UnixAddr{Net: "review-invalid-network", Name: "unused"}
	if err := c.Listen(context.Background()); err == nil {
		t.Fatal("invalid network accepted")
	}
	if c.pktConn != p {
		t.Error("failed Listen discarded the existing PacketConn without closing it")
	}
}

func TestConnInitialSourcePortChange(t *testing.T) {
	c, _ := testConnection(t)
	s := testSession("001010000000001")
	c.RegisterSession(123, s)
	samePeerNewPort := &net.UDPAddr{IP: net.IPv4(192, 0, 2, 2), Port: 54321}
	if err := c.validate(samePeerNewPort, message.NewModifyBearerRequest(123, 2)); err != nil {
		t.Errorf("valid peer Initial message from a newly allocated source port rejected: %v", err)
	}
}
func TestConnReplySourcePort(t *testing.T) {
	c, p := testConnection(t)
	req := message.NewEchoRequest(77, ie.NewRecovery(1))
	if err := c.handleMessage(testPeer(), req); err != nil {
		t.Fatal(err)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.writes) != 1 || p.peers[0].String() != testPeer().String() {
		t.Fatal("Echo response not sent to request source IP and port")
	}
	m, err := message.Parse(p.writes[0])
	if err != nil {
		t.Fatal(err)
	}
	if m.Sequence() != 77 {
		t.Errorf("Echo reply sequence=%d, want 77", m.Sequence())
	}
}
