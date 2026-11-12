// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv1

import (
	"net"
	"runtime"
	"testing"
	"time"
)

func TestWriteToWithDSCPECN(t *testing.T) {
	for _, tc := range []struct {
		name string
		ip   net.IP
	}{
		{"IPv4", net.IPv4(127, 0, 0, 1)},
		{"IPv6", net.IPv6loopback},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dst, err := net.ListenUDP("udp", &net.UDPAddr{IP: tc.ip})
			if err != nil {
				t.Skipf("loopback not available: %v", err)
			}
			defer func() { _ = dst.Close() }()

			src, err := net.ListenUDP("udp", &net.UDPAddr{IP: tc.ip})
			if err != nil {
				t.Fatal(err)
			}
			pkt, err := newPktConnFromUDPConn(src, tc.ip)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = pkt.Close() }()

			dscp, ok := pkt.(interface{ DSCPECN() (int, error) })
			if !ok {
				t.Fatalf("%T does not expose DSCPECN", pkt)
			}

			if _, err := dscp.DSCPECN(); err != nil {
				// DSCP/ECN 0 needs no socket option, so plain writes still work.
				if _, err := pkt.WriteToWithDSCPECN([]byte{0}, dst.LocalAddr(), 0); err != nil {
					t.Fatalf("plain write: %v", err)
				}
				t.Skipf("DSCP/ECN unsupported on %s/%s: %v (https://github.com/golang/go/issues/7174)", runtime.GOOS, runtime.GOARCH, err)
			}

			// 0 matches the socket default (fast path), 0xb8 (EF) must be
			// applied for the write only and restored afterwards.
			for _, v := range []int{0, 0xb8, 0} {
				payload := []byte{byte(v)}
				if _, err := pkt.WriteToWithDSCPECN(payload, dst.LocalAddr(), v); err != nil {
					t.Fatalf("write with DSCP/ECN %#x: %v", v, err)
				}

				got, err := dscp.DSCPECN()
				if err != nil {
					t.Fatal(err)
				}
				if got != 0 {
					t.Errorf("DSCP/ECN after write with %#x = %#x, want 0 restored", v, got)
				}

				buf := make([]byte, 16)
				if err := dst.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
					t.Fatal(err)
				}
				n, err := dst.Read(buf)
				if err != nil {
					t.Fatal(err)
				}
				if n != 1 || buf[0] != byte(v) {
					t.Errorf("received %x, want %x", buf[:n], payload)
				}
			}
		})
	}
}

func BenchmarkWriteToWithDSCPECN(b *testing.B) {
	dst, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = dst.Close() }()
	go func() {
		buf := make([]byte, 2048)
		for {
			if _, err := dst.Read(buf); err != nil {
				return
			}
		}
	}()

	src, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		b.Fatal(err)
	}
	pkt, err := newPktConnFromUDPConn(src, net.IPv4(127, 0, 0, 1))
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = pkt.Close() }()

	payload := make([]byte, 1280)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := pkt.WriteToWithDSCPECN(payload, dst.LocalAddr(), 0); err != nil {
			b.Fatal(err)
		}
	}
}
