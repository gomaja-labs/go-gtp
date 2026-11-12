package gtpv1

import (
	"context"
	"errors"
	"net"
	"os"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
)

func TestPayloadBacklogDoesNotDelayEcho(t *testing.T) {
	u := serveUDP(t, listenUDP(t, "127.0.0.1", 0), func(u *UPlaneConn) { u.DisableErrorIndication() })
	peer := listenUDP(t, "127.0.0.1", 0)
	// Sequential datagrams on one socket put the echo behind the unread payload.
	for range 100 {
		sendUDP(t, peer, u.LocalAddr(), message.NewTPDU(1, []byte{0x45}))
	}
	sendUDP(t, peer, u.LocalAddr(), message.NewEchoRequest(42))
	response, _ := readUDP(t, peer)
	if response.MessageType() != message.MsgTypeEchoResponse || response.Sequence() != 42 {
		t.Fatal(response)
	}
}

func TestPayloadDeadlineDoesNotStopEcho(t *testing.T) {
	u := serveUDP(t, listenUDP(t, "127.0.0.1", 0), nil)
	peer := listenUDP(t, "127.0.0.1", 0)
	if err := u.SetReadDeadline(time.Now().Add(10 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := u.ReadFromGTP(make([]byte, 100)); !errors.Is(err, os.ErrDeadlineExceeded) {
		t.Fatal(err)
	}
	sendUDP(t, peer, u.LocalAddr(), message.NewEchoRequest(43))
	response, _ := readUDP(t, peer)
	if response.MessageType() != message.MsgTypeEchoResponse || response.Sequence() != 43 {
		t.Fatal(response)
	}
}

func TestEndMarkerUnknownExtensionIsNotNotified(t *testing.T) {
	u := serveUDP(t, listenUDP(t, "127.0.0.1", 0), nil)
	peer := listenUDP(t, "127.0.0.3", 2152)
	raw := []byte{0x34, 0xfe, 0, 8, 0, 0, 0, 1, 0, 0, 0, 0xc7, 1, 0xaa, 0xbb, 0}
	if _, err := peer.WriteTo(raw, u.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	// The echo is an ordering barrier: any incorrect notification precedes it.
	sendUDP(t, peer, u.LocalAddr(), message.NewEchoRequest(44))
	response, _ := readUDP(t, peer)
	if response.MessageType() != message.MsgTypeEchoResponse || response.Sequence() != 44 {
		t.Fatalf("unexpected notification: %v", response)
	}
}

func TestSecondServeIsRejected(t *testing.T) {
	u := serveUDP(t, listenUDP(t, "127.0.0.1", 0), nil)
	peer := listenUDP(t, "127.0.0.1", 0)
	sendUDP(t, peer, u.LocalAddr(), message.NewEchoRequest(45))
	_, _ = readUDP(t, peer) // First serving loop has entered dispatch.
	done := make(chan error, 1)
	go func() { done <- u.ListenAndServe(context.Background()) }()
	select {
	case err := <-done:
		if err == nil || errors.Is(err, net.ErrClosed) {
			t.Fatalf("second serve: %v", err)
		}
	case <-time.After(time.Second):
		_ = u.Close()
		if err := <-done; err != nil && !errors.Is(err, net.ErrClosed) {
			t.Error(err)
		}
		t.Fatal("second serving loop started")
	}
	sendUDP(t, peer, u.LocalAddr(), message.NewEchoRequest(46))
	response, _ := readUDP(t, peer)
	if response.Sequence() != 46 {
		t.Fatal(response)
	}
}

func TestDialReservesItsServingLoop(t *testing.T) {
	peer := listenUDP(t, "127.0.0.1", 0)
	result := make(chan *UPlaneConn, 1)
	errs := make(chan error, 1)
	go func() {
		u, err := DialUPlane(context.Background(), &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, peer.LocalAddr())
		if err != nil {
			errs <- err
			return
		}
		result <- u
	}()
	request, client := readUDP(t, peer)
	sendUDP(t, peer, client, message.NewEchoResponse(request.Sequence(), ie.NewRecovery(0)))
	var u *UPlaneConn
	select {
	case u = <-result:
	case err := <-errs:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Dial did not finish")
	}
	defer func() { _ = u.Close() }()
	done := make(chan error, 1)
	go func() { done <- u.ListenAndServe(context.Background()) }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("duplicate serve accepted")
		}
	case <-time.After(time.Second):
		_ = u.Close()
		if err := <-done; err != nil && !errors.Is(err, net.ErrClosed) {
			t.Error(err)
		}
		t.Fatal("Dial connection started a second loop")
	}
	sendUDP(t, peer, client, message.NewEchoRequest(48))
	response, _ := readUDP(t, peer)
	if response.Sequence() != 48 || response.MessageType() != message.MsgTypeEchoResponse {
		t.Fatal(response)
	}
}

func TestPayloadDeadlineOnClosedConnection(t *testing.T) {
	u, err := NewUPlaneConnFromUDPConn(listenUDP(t, "127.0.0.1", 0))
	if err != nil {
		t.Fatal(err)
	}
	if err := u.Close(); err != nil {
		t.Fatal(err)
	}
	if err := u.SetReadDeadline(time.Time{}); !errors.Is(err, net.ErrClosed) {
		t.Fatalf("closed deadline: %v", err)
	}
}
