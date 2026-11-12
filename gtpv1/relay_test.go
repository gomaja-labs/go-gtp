// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv1_test

import (
	"bytes"
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv1"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
)

// listenLoopback opens a UDP socket on an ephemeral loopback port that is
// closed when the test ends.
func listenLoopback(t *testing.T) *net.UDPConn {
	t.Helper()
	c, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

// serveUPlane serves a UPlaneConn on an ephemeral loopback port and stops it,
// waiting for the serving loop to return, before the test ends.
func serveUPlane(t *testing.T) *gtpv1.UPlaneConn {
	t.Helper()
	u, err := gtpv1.NewUPlaneConnFromUDPConn(listenLoopback(t))
	if err != nil {
		t.Fatal(err)
	}
	u.DisableErrorIndication()
	ctx, cancel := context.WithCancel(context.Background())
	served := make(chan error, 1)
	go func() { served <- u.ListenAndServe(ctx) }()
	t.Cleanup(func() {
		cancel()
		_ = u.Close()
		select {
		case err := <-served:
			if err != nil && !errors.Is(err, net.ErrClosed) {
				t.Errorf("ListenAndServe: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Error("ListenAndServe did not return after Close")
		}
	})
	return u
}

// TestRelay sends a T-PDU into one UPlaneConn and checks that the relay
// forwards it through another UPlaneConn with the outgoing TEID.
func TestRelay(t *testing.T) {
	left, right := serveUPlane(t), serveUPlane(t)
	sender, sink := listenLoopback(t), listenLoopback(t)

	const teidIn, teidOut = 0x11111111, 0x22222222
	if err := left.RelayTo(right, teidIn, teidOut, sink.LocalAddr()); err != nil {
		t.Fatal(err)
	}

	payload := []byte{0xde, 0xad, 0xbe, 0xef}
	raw, err := message.NewTPDU(teidIn, payload).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sender.WriteTo(raw, left.LocalAddr()); err != nil {
		t.Fatal(err)
	}

	if err := sink.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1500)
	n, from, err := sink.ReadFrom(buf)
	if err != nil {
		t.Fatalf("relayed T-PDU not received: %v", err)
	}
	if from.String() != right.LocalAddr().String() {
		t.Errorf("relayed from %s, want %s", from, right.LocalAddr())
	}
	msg, err := message.Parse(buf[:n])
	if err != nil {
		t.Fatal(err)
	}
	tpdu, ok := msg.(*message.TPDU)
	if !ok {
		t.Fatalf("relayed %s, want T-PDU", msg.MessageTypeName())
	}
	if tpdu.TEID() != teidOut {
		t.Errorf("relayed TEID %#x, want %#x", tpdu.TEID(), uint32(teidOut))
	}
	if !bytes.Equal(tpdu.Decapsulate(), payload) {
		t.Errorf("relayed payload % x, want % x", tpdu.Decapsulate(), payload)
	}
}
