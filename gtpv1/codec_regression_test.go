package gtpv1

import (
	"bytes"
	"context"
	"errors"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

type testPacketConn struct {
	net.PacketConn
	read   func([]byte) (int, net.Addr, error)
	writes chan []byte
	dests  chan net.Addr
	closed chan struct{}
	once   sync.Once
	addr   net.Addr
}

func testConnection() (*UPlaneConn, *testPacketConn) {
	p := &testPacketConn{writes: make(chan []byte, 4096), dests: make(chan net.Addr, 4096), closed: make(chan struct{})}
	u := NewUPlaneConn(&net.UDPAddr{IP: net.IPv4(192, 0, 2, 1), Port: 2152})
	u.pktConn = p
	return u, p
}
func (p *testPacketConn) ReadFrom(b []byte) (int, net.Addr, error) { return p.read(b) }
func (p *testPacketConn) WriteToWithDSCPECN(b []byte, addr net.Addr, _ int) (int, error) {
	p.dests <- addr
	p.writes <- bytes.Clone(b)
	return len(b), nil
}
func TestErrorIndicationPort(t *testing.T) {
	u, p := testConnection()
	if err := u.ErrorIndication(&net.UDPAddr{IP: net.IPv4(192, 0, 2, 2), Port: 49152}, message.NewTPDU(1, []byte{1})); err != nil {
		t.Fatal(err)
	}
	got := <-p.dests
	if a, ok := got.(*net.UDPAddr); !ok || a.Port != 2152 {
		t.Fatalf("Error Indication destination=%s; want UDP port 2152", got)
	}
}
func (p *testPacketConn) Close() error            { p.once.Do(func() { close(p.closed) }); return nil }
func (p *testPacketConn) File() (*os.File, error) { return nil, errors.New("not used") }
func (p *testPacketConn) LocalAddr() net.Addr {
	if p.addr != nil {
		return p.addr
	}
	return &net.UDPAddr{IP: net.IPv4(192, 0, 2, 1), Port: 2152}
}
func TestErrorIndicationWildcardAddress(t *testing.T) {
	u, p := testConnection()
	p.addr = &net.UDPAddr{IP: net.IPv4zero, Port: 2152}
	if err := u.ErrorIndication(p.LocalAddr(), message.NewTPDU(1, []byte{1})); err != nil {
		return
	}
	b := <-p.writes
	m, err := message.ParseErrorIndication(b)
	if err != nil {
		t.Fatal(err)
	}
	got := m.GTPUPeerAddress.MustIPAddress()
	if got == "0.0.0.0" || got == "::" {
		t.Fatalf("Error Indication identifies wildcard address %s rather than the triggering destination", got)
	}
}

func TestRelayShortPacket(t *testing.T) {
	if os.Getenv("GTPV1_SHORT_RELAY_PACKET_CHILD") == "1" {
		u, p := testConnection()
		_ = u.RelayTo(u, 1, 2, p.LocalAddr())
		n := 0
		p.read = func(b []byte) (int, net.Addr, error) {
			n++
			if n == 1 {
				return 0, p.LocalAddr(), nil
			}
			<-p.closed
			return 0, nil, io.EOF
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		_ = u.serve(ctx)
		return
	}
	c := exec.Command(os.Args[0], "-test.run=^TestRelayShortPacket$", "-test.timeout=3s")
	c.Env = append(os.Environ(), "GTPV1_SHORT_RELAY_PACKET_CHILD=1")
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("short relay datagram crashed process: %s", out)
	}
}
func TestRelayMapRace(t *testing.T) {
	u, p := testConnection()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	raw, _ := message.NewEchoRequest(1).Marshal()
	count := 0
	p.read = func(b []byte) (int, net.Addr, error) {
		count++
		if count > 2000 {
			return 0, nil, io.EOF
		}
		return copy(b, raw), p.LocalAddr(), nil
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 2000; j++ {
			_ = u.RelayTo(u, uint32(j), 2, p.LocalAddr())
			_ = u.CloseRelay(uint32(j))
		}
	}()
	go func() { defer wg.Done(); _ = u.serve(ctx) }()
	wg.Wait()
}
func TestErrorIndicationFlagRace(t *testing.T) {
	u, p := testConnection()
	pdu := message.NewTPDU(1, []byte{1})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			u.DisableErrorIndication()
			u.EnableErrorIndication()
		}
	}()
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			_ = handleTPDU(u, p.LocalAddr(), pdu)
		}
	}()
	wg.Wait()
	_ = u.Close()
}
func TestFullSizeRelayPacket(t *testing.T) {
	u, p := testConnection()
	_ = u.RelayTo(u, 1, 2, p.LocalAddr())
	raw, _ := message.NewTPDU(1, make([]byte, 1500)).Marshal()
	count := 0
	p.read = func(b []byte) (int, net.Addr, error) {
		count++
		if count > 1 {
			return 0, nil, io.EOF
		}
		return copy(b, raw), p.LocalAddr(), nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := u.serve(ctx); err != nil {
		t.Fatal(err)
	}
	got := <-p.writes
	if len(got) != len(raw) {
		t.Fatalf("relay forwarded %d bytes; want %d", len(got), len(raw))
	}
}
func TestCloseLifecycle(t *testing.T) {
	u, p := testConnection()
	if err := u.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-p.closed:
	default:
		t.Error("Close did not close owned packet connection before serve")
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("second Close panicked: %v", r)
		}
	}()
	_ = u.Close()
}
func TestReadAfterClose(t *testing.T) {
	u, _ := testConnection()
	_ = u.Close()
	if _, _, _, err := u.ReadFromGTP(make([]byte, 10)); err == nil {
		t.Fatal("ReadFromGTP after Close returned success")
	}
}
func TestPacketOrdering(t *testing.T) {
	u, p := testConnection()
	firstStarted := make(chan struct{})
	release := make(chan struct{})
	got := make(chan uint16, 2)
	u.AddHandler(message.MsgTypeEchoRequest, func(_ Conn, _ net.Addr, m message.Message) error {
		if m.Sequence() == 1 {
			close(firstStarted)
			pcs := make([]uintptr, 32)
			n := runtime.Callers(0, pcs)
			frames := runtime.CallersFrames(pcs[:n])
			inline := false
			for {
				f, more := frames.Next()
				if strings.HasSuffix(f.Function, "(*UPlaneConn).serve") {
					inline = true
				}
				if !more {
					break
				}
			}
			if !inline {
				<-release
			}
		}
		got <- m.Sequence()
		if m.Sequence() == 2 {
			close(release)
		}
		return nil
	})
	count := 0
	p.read = func(b []byte) (int, net.Addr, error) {
		count++
		if count > 2 {
			return 0, nil, io.EOF
		}
		if count == 2 {
			<-firstStarted
		}
		raw, _ := message.NewEchoRequest(uint16(count)).Marshal()
		return copy(b, raw), p.LocalAddr(), nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := u.serve(ctx); err != nil {
		t.Fatal(err)
	}
	a, b := <-got, <-got
	if a != 1 || b != 2 {
		t.Fatalf("handler order=%d,%d; receive order=1,2", a, b)
	}
}
func TestZeroTEIDErrorIndication(t *testing.T) {
	u, p := testConnection()
	if err := handleTPDU(u, p.LocalAddr(), message.NewTPDU(0, []byte{1})); err != nil {
		t.Fatal(err)
	}
	select {
	case b := <-p.writes:
		t.Fatalf("sent Error Indication for TEID zero: %x", b)
	default:
	}
}
func TestUnknownRequiredExtension(t *testing.T) {
	u, p := testConnection()
	h := message.NewTPDUWithExtentionHeader(1, []byte{1}, message.NewExtensionHeader(0xff, []byte{0, 0}, 0))
	if err := u.handleMessage(p.LocalAddr(), h); err != nil {
		t.Fatal(err)
	}
	b := <-p.writes
	if b[1] != message.MsgTypeSupportedExtensionHeaderNotification {
		t.Fatalf("unknown mandatory extension produced message type %d; want 31", b[1])
	}
}
func TestErrorIndicationMissingIEs(t *testing.T) {
	u, p := testConnection()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("missing mandatory IE crashed handler: %v", r)
		}
	}()
	raw, _ := message.NewHeader(0x32, 26, 0, 1, nil).Marshal()
	m, err := message.Parse(raw)
	if err != nil {
		return
	}
	if err := u.handleMessage(p.LocalAddr(), m); err == nil {
		t.Error("accepted missing mandatory IE")
	}
}
func TestServeFailureCleanup(t *testing.T) {
	u, p := testConnection()
	p.read = func([]byte) (int, net.Addr, error) { return 0, nil, errors.New("synthetic read failure") }
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := u.serve(ctx); err == nil {
		t.Fatal("expected read failure")
	}
	select {
	case <-p.closed:
	default:
		t.Error("serve error returned without closing packet conn")
	}
}
func TestDialCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	u, err := DialUPlane(ctx, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 2152})
	if err != nil && strings.Contains(err.Error(), "operation not permitted") {
		t.Skip("sandbox prohibits UDP sockets; cancellation path cannot be exercised")
	}
	if u != nil {
		_ = u.Close()
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled Dial returned conn=%v err=%v; want context.Canceled", u, err)
	}
}
func TestRelayValidation(t *testing.T) {
	u, p := testConnection()
	_ = u.RelayTo(u, 1, 2, p.LocalAddr())
	raw, _ := message.NewTPDU(1, []byte{1, 2, 3}).Marshal()
	raw[0] = 0x50
	raw[3] = 100
	count := 0
	p.read = func(b []byte) (int, net.Addr, error) {
		count++
		if count > 1 {
			return 0, nil, io.EOF
		}
		return copy(b, raw), p.LocalAddr(), nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = u.serve(ctx)
	select {
	case got := <-p.writes:
		t.Fatalf("forwarded malformed version/length without validation: %x", got)
	case <-time.After(time.Second):
	}
}
