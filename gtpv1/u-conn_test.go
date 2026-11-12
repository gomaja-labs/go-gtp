// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv1_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/gomaja-labs/go-gtp/gtpv1"
)

type testVal struct {
	teidIn, teidOut uint32
	seq             uint16
	payload         []byte
}

// setup starts a serving UPlaneConn on 127.0.0.2:2152 and a client dialled to
// it from 127.0.0.1:2152. Both are closed when the test ends.
func setup(t *testing.T, ctx context.Context) (cliConn, srvConn *gtpv1.UPlaneConn) {
	t.Helper()
	cliAddr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 2152}
	srvAddr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 2), Port: 2152}

	srvConn = gtpv1.NewUPlaneConn(srvAddr)
	srvConn.DisableErrorIndication()
	served := make(chan error, 1)
	go func() { served <- srvConn.ListenAndServe(ctx) }()
	t.Cleanup(func() {
		_ = srvConn.Close()
		select {
		case err := <-served:
			if err != nil {
				t.Logf("ListenAndServe after Close: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Error("ListenAndServe did not return after Close")
		}
	})

	// DialUPlane's Echo handshake succeeds once the server is serving; retry
	// while it is still starting.
	deadline := time.Now().Add(5 * time.Second)
	for {
		dctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
		c, err := gtpv1.DialUPlane(dctx, cliAddr, srvAddr)
		cancel()
		if err == nil {
			cliConn = c
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("server did not answer the Echo handshake: %v", err)
		}
	}
	cliConn.DisableErrorIndication()
	t.Cleanup(func() { _ = cliConn.Close() })
	return cliConn, srvConn
}

func TestClientWrite(t *testing.T) {
	var (
		okCh  = make(chan struct{})
		errCh = make(chan error)
		buf   = make([]byte, 2048)
		tv    = &testVal{
			0x11111111, 0x22222222, 0x3333,
			[]byte{0xde, 0xad, 0xbe, 0xef},
		}
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cliConn, srvConn := setup(t, ctx)

	go func(tv *testVal) {
		n, addr, teid, err := srvConn.ReadFromGTP(buf)
		if err != nil {
			errCh <- err
			return
		}

		if diff := cmp.Diff(n, len(tv.payload)); diff != "" {
			t.Error(diff)
		}
		if diff := cmp.Diff(addr, cliConn.LocalAddr()); diff != "" {
			t.Error(diff)
		}
		if diff := cmp.Diff(teid, tv.teidOut); diff != "" {
			t.Error(diff)
		}
		if diff := cmp.Diff(buf[:n], tv.payload); diff != "" {
			t.Error(diff)
		}
		okCh <- struct{}{}
	}(tv)

	if _, err := cliConn.WriteToGTP(tv.teidOut, tv.payload, srvConn.LocalAddr()); err != nil {
		t.Fatal(err)
	}

	select {
	case <-okCh:
		return
	case err := <-errCh:
		t.Fatal(err)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out while waiting for response to come")
	}
}

func TestNewUPlaneConnFromUDPConn(t *testing.T) {
	var (
		okCh  = make(chan struct{})
		errCh = make(chan error)
		buf   = make([]byte, 2048)
		tv    = &testVal{
			0x11111111, 0x22222222, 0x3333,
			[]byte{0xde, 0xad, 0xbe, 0xef},
		}
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	// the point of the constructor: the socket can be tuned before serving.
	if err := udpConn.SetReadBuffer(1 << 20); err != nil {
		t.Fatal(err)
	}

	srvConn, err := gtpv1.NewUPlaneConnFromUDPConn(udpConn)
	if err != nil {
		t.Fatal(err)
	}
	srvConn.DisableErrorIndication()
	if diff := cmp.Diff(srvConn.LocalAddr(), udpConn.LocalAddr()); diff != "" {
		t.Error(diff)
	}

	go func() {
		if err := srvConn.ListenAndServe(ctx); err != nil {
			errCh <- err
		}
	}()

	cliConn, err := gtpv1.DialUPlane(ctx, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, srvConn.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}

	go func(tv *testVal) {
		n, _, teid, err := srvConn.ReadFromGTP(buf)
		if err != nil {
			errCh <- err
			return
		}
		if diff := cmp.Diff(teid, tv.teidOut); diff != "" {
			t.Error(diff)
		}
		if diff := cmp.Diff(buf[:n], tv.payload); diff != "" {
			t.Error(diff)
		}
		okCh <- struct{}{}
	}(tv)

	if _, err := cliConn.WriteToGTP(tv.teidOut, tv.payload, srvConn.LocalAddr()); err != nil {
		t.Fatal(err)
	}

	select {
	case <-okCh:
		return
	case err := <-errCh:
		t.Fatal(err)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out while waiting for response to come")
	}
}

func TestNewUPlaneConnFromUDPConnNil(t *testing.T) {
	if _, err := gtpv1.NewUPlaneConnFromUDPConn(nil); err == nil {
		t.Fatal("expected an error for a nil conn")
	}
}
