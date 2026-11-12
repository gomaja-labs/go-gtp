package gtpv1

import (
	"bytes"
	"context"
	"errors"
	"net"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func listenUDP(t *testing.T, ip string, port int) *net.UDPConn {
	t.Helper()
	c, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	// macOS limits a UDP datagram to the send buffer (9216 octets by default);
	// the tests send datagrams up to 64 KiB.
	if err := c.SetWriteBuffer(1 << 17); err != nil {
		t.Fatal(err)
	}
	return c
}

func serveUDP(t *testing.T, c *net.UDPConn, configure func(*UPlaneConn)) *UPlaneConn {
	t.Helper()
	u, err := NewUPlaneConnFromUDPConn(c)
	if err != nil {
		t.Fatal(err)
	}
	if configure != nil {
		configure(u)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- u.ListenAndServe(ctx) }()
	t.Cleanup(func() {
		cancel()
		_ = u.Close()
		select {
		case err := <-done:
			if err != nil && !errors.Is(err, net.ErrClosed) {
				t.Error(err)
			}
		case <-time.After(time.Second):
			t.Error("serve did not exit")
		}
	})
	return u
}

func sendUDP(t *testing.T, c *net.UDPConn, to net.Addr, m message.Message) {
	t.Helper()
	b, err := message.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = c.WriteTo(b, to); err != nil {
		t.Fatal(err)
	}
}

func readUDP(t *testing.T, c *net.UDPConn) (message.Message, *net.UDPAddr) {
	t.Helper()
	if err := c.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	b := make([]byte, 65535)
	n, from, err := c.ReadFromUDP(b)
	if err != nil {
		t.Fatal(err)
	}
	m, err := message.Parse(b[:n])
	if err != nil {
		t.Fatal(err)
	}
	return m, from
}

func TestOrderedPayloadOwnership(t *testing.T) {
	sock := listenUDP(t, "127.0.0.1", 0)
	received := make(chan []byte, 16)
	u := serveUDP(t, sock, func(u *UPlaneConn) {
		u.AddHandler(message.MsgTypeTPDU, func(_ Conn, _ net.Addr, m message.Message) error { received <- m.(*message.TPDU).Payload; return nil })
	})
	sender := listenUDP(t, "127.0.0.1", 0)
	var kept [][]byte
	for j := 1; j <= 8; j++ {
		sendUDP(t, sender, u.LocalAddr(), message.NewTPDU(1, bytes.Repeat([]byte{byte(j)}, 1500)))
		select {
		case p := <-received:
			kept = append(kept, p)
		case <-time.After(time.Second):
			t.Fatal("packet not dispatched")
		}
	}
	for j, p := range kept {
		if !bytes.Equal(p, bytes.Repeat([]byte{byte(j + 1)}, 1500)) {
			t.Fatalf("packet %d overwritten or reordered", j)
		}
	}
}

func TestLargeDatagramDelivery(t *testing.T) {
	sock := listenUDP(t, "127.0.0.1", 0)
	u := serveUDP(t, sock, func(u *UPlaneConn) { u.DisableErrorIndication() })
	sender := listenUDP(t, "127.0.0.1", 0)
	want := bytes.Repeat([]byte{0xa5}, 60000)
	sendUDP(t, sender, u.LocalAddr(), message.NewTPDU(1, want))
	done := make(chan error, 1)
	go func() {
		b := make([]byte, 65535)
		n, _, teid, err := u.ReadFromGTP(b)
		if err == nil && (teid != 1 || !bytes.Equal(b[:n], want)) {
			err = errors.New("truncated or corrupted datagram")
		}
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("large datagram not delivered")
	}
}

func TestErrorIndicationWildcardDestinationAndPort(t *testing.T) {
	requireDestinationAddresses(t, "udp4")
	// TS 29.281 §§4.4.2.4, 7.3.1: notify port 2152 with the original destination address.
	sock := listenUDP(t, "0.0.0.0", 0)
	u := serveUDP(t, sock, nil)
	receiver := listenUDP(t, "127.0.0.3", 2152)
	sender := listenUDP(t, "127.0.0.3", 0)
	dst := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 2), Port: u.LocalAddr().(*net.UDPAddr).Port}
	sendUDP(t, sender, dst, message.NewTPDU(0x12345678, []byte{0x45}))
	m, from := readUDP(t, receiver)
	ind, ok := m.(*message.ErrorIndication)
	if !ok {
		t.Fatalf("got %T", m)
	}
	if got := ind.GTPUPeerAddress.MustIPAddress(); got != "127.0.0.2" {
		t.Fatalf("peer address %s", got)
	}
	if ind.TEIDDataI.MustTEID() != 0x12345678 || !from.IP.Equal(dst.IP) {
		t.Fatalf("wrong TEID or source address: %v", from)
	}
	sendUDP(t, sender, dst, message.NewTPDU(0, []byte{0x45}))
	if err := receiver.SetReadDeadline(time.Now().Add(50 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	if _, _, err := receiver.ReadFromUDP(make([]byte, 65535)); err == nil {
		t.Fatal("Error Indication sent for zero TEID")
	}
}

func TestRequiredExtensionBeforeRelay(t *testing.T) {
	sock := listenUDP(t, "127.0.0.1", 0)
	u := serveUDP(t, sock, nil)
	receiver := listenUDP(t, "127.0.0.4", 2152)
	relay := listenUDP(t, "127.0.0.1", 0)
	if err := u.RelayTo(u, 1, 2, relay.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	sendUDP(t, receiver, u.LocalAddr(), message.NewTPDUWithExtentionHeader(1, []byte{1}, message.NewExtensionHeader(0xff, []byte{1, 2}, 0)))
	m, _ := readUDP(t, receiver)
	if m.MessageType() != message.MsgTypeSupportedExtensionHeaderNotification {
		t.Fatalf("got %T", m)
	}
	if err := relay.SetReadDeadline(time.Now().Add(50 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	if _, _, err := relay.ReadFromUDP(make([]byte, 65535)); err == nil {
		t.Fatal("forwarded unknown comprehension-required extension")
	}
	sendUDP(t, receiver, u.LocalAddr(), message.NewTPDUWithExtentionHeader(1, []byte{1}, message.NewExtensionHeader(0x7f, []byte{1, 2}, 0)))
	m, _ = readUDP(t, relay)
	if m.TEID() != 2 {
		t.Fatalf("optional extension did not relay: %d", m.TEID())
	}
}

func TestEchoWildcardSourceAndSequence(t *testing.T) {
	requireDestinationAddresses(t, "udp4")
	sock := listenUDP(t, "0.0.0.0", 0)
	u := serveUDP(t, sock, nil)
	sender := listenUDP(t, "127.0.0.1", 0)
	dst := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 2), Port: u.LocalAddr().(*net.UDPAddr).Port}
	sendUDP(t, sender, dst, message.NewEchoRequest(0x1234))
	m, from := readUDP(t, sender)
	res, ok := m.(*message.EchoResponse)
	if !ok || m.Sequence() != 0x1234 || res.Recovery.MustRecovery() != 0 || !from.IP.Equal(dst.IP) {
		t.Fatalf("wrong echo response: %v from %v", m, from)
	}
}

func TestDialUPlaneCancelHandshake(t *testing.T) {
	server := listenUDP(t, "127.0.0.1", 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		u, err := DialUPlane(ctx, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, server.LocalAddr())
		if u != nil {
			_ = u.Close()
		}
		result <- err
	}()
	_, client := readUDP(t, server)
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("handshake did not cancel")
	}
	rebound, err := net.ListenUDP("udp4", client)
	if err != nil {
		t.Fatalf("handshake socket leaked: %v", err)
	}
	_ = rebound.Close()
}

func TestHandshakeCorrelation(t *testing.T) {
	server := listenUDP(t, "127.0.0.1", 0)
	impostor := listenUDP(t, "127.0.0.1", 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	type result struct {
		u   *UPlaneConn
		err error
	}
	done := make(chan result, 1)
	go func() {
		u, err := DialUPlane(ctx, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, server.LocalAddr())
		done <- result{u, err}
	}()
	request, client := readUDP(t, server)
	sendUDP(t, impostor, client, message.NewEchoResponse(request.Sequence(), ie.NewRecovery(0)))
	sendUDP(t, server, client, message.NewEchoResponse(request.Sequence()+1, ie.NewRecovery(0)))
	sendUDP(t, server, client, message.NewTPDU(1, []byte{1}))
	select {
	case r := <-done:
		if r.u != nil {
			_ = r.u.Close()
		}
		t.Fatalf("accepted unrelated handshake packet: %v", r.err)
	case <-time.After(50 * time.Millisecond):
	}
	sendUDP(t, server, client, message.NewEchoResponse(request.Sequence(), ie.NewRecovery(0)))
	select {
	case r := <-done:
		if r.err != nil {
			t.Fatal(r.err)
		}
		_ = r.u.Close()
	case <-time.After(time.Second):
		t.Fatal("matching response did not complete handshake")
	}
}

func TestIPv6WildcardDestination(t *testing.T) {
	requireDestinationAddresses(t, "udp6")
	sock, err := net.ListenUDP("udp6", &net.UDPAddr{IP: net.IPv6unspecified})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sock.Close() })
	u := serveUDP(t, sock, nil)
	receiver, err := net.ListenUDP("udp6", &net.UDPAddr{IP: net.IPv6loopback, Port: 2152})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = receiver.Close() })
	sender, err := net.ListenUDP("udp6", &net.UDPAddr{IP: net.IPv6loopback})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sender.Close() })
	dst := &net.UDPAddr{IP: net.IPv6loopback, Port: u.LocalAddr().(*net.UDPAddr).Port}
	sendUDP(t, sender, dst, message.NewTPDU(1, []byte{1}))
	m, from := readUDP(t, receiver)
	ind, ok := m.(*message.ErrorIndication)
	if !ok || ind.GTPUPeerAddress.MustIPAddress() != "::1" || !from.IP.Equal(net.IPv6loopback) {
		t.Fatalf("wrong IPv6 indication: %v from %v", m, from)
	}
}

func TestConcurrentCloseUnblocksRead(t *testing.T) {
	sock := listenUDP(t, "127.0.0.1", 0)
	u := serveUDP(t, sock, nil)
	done := make(chan error, 1)
	go func() { _, _, _, err := u.ReadFromGTP(make([]byte, 1)); done <- err }()
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			if err := u.Close(); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	select {
	case err := <-done:
		if !errors.Is(err, net.ErrClosed) {
			t.Fatalf("read after close: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("reader did not unblock")
	}
}

// requireDestinationAddresses skips a test that needs the destination address
// of received datagrams where golang.org/x/net cannot report it (Windows:
// https://github.com/golang/go/issues/7175).
func requireDestinationAddresses(t *testing.T, network string) {
	t.Helper()
	ip := net.IPv4(127, 0, 0, 1)
	if network == "udp6" {
		ip = net.IPv6loopback
	}
	c, err := net.ListenUDP(network, &net.UDPAddr{IP: ip})
	if err != nil {
		t.Skipf("%s loopback unavailable: %v", network, err)
	}
	defer func() { _ = c.Close() }()
	if network == "udp6" {
		err = ipv6.NewPacketConn(c).SetControlMessage(ipv6.FlagDst, true)
	} else {
		err = ipv4.NewPacketConn(c).SetControlMessage(ipv4.FlagDst, true)
	}
	if err != nil {
		t.Skipf("destination addresses of received datagrams unavailable on %s/%s: %v (https://github.com/golang/go/issues/7175)", runtime.GOOS, runtime.GOARCH, err)
	}
}
