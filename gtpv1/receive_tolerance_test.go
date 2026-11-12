package gtpv1

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv1/message"
)

func TestRelayExtensionComprehensionRules(t *testing.T) {
	// TS 29.281 V19.2.0 §5.2.1 Figure 5.2.1-2: relays are intermediate
	// nodes: forward 10xxxxxx; remove 01xxxxxx and retain the payload.
	for _, typ := range []byte{0x11, 0x81, 0x7f, 0x40, 0xcf} {
		t.Run(string([]byte{'x', "0123456789abcdef"[typ>>4], "0123456789abcdef"[typ&15]}), func(t *testing.T) {
			u, p := testConnection()
			if err := u.RelayTo(u, 1, 2, p.LocalAddr()); err != nil {
				t.Fatal(err)
			}
			raw := []byte{0x34, 255, 0, 9, 0, 0, 0, 1, 0, 0, 0, typ, 1, 1, 2, 0, 0x45}
			read := false
			p.read = func(b []byte) (int, net.Addr, error) {
				if read {
					return 0, nil, io.EOF
				}
				read = true
				return copy(b, raw), p.LocalAddr(), nil
			}
			if err := u.serve(context.Background()); err != nil {
				t.Fatal(err)
			}
			select {
			case b := <-p.writes:
				m, err := message.Parse(b)
				if err != nil {
					t.Fatal(err)
				}
				if typ == 0xcf {
					if m.MessageType() != message.MsgTypeSupportedExtensionHeaderNotification {
						t.Fatal("unknown 11 extension was not notified")
					}
					return
				}
				if m.MessageType() != 255 || m.TEID() != 2 {
					t.Fatalf("intermediate node sent message type %d TEID %d instead of forwarding G-PDU", m.MessageType(), m.TEID())
				}
				if typ&0xc0 == 0x40 && len(m.GetHeader().ExtensionHeaders) != 0 {
					t.Fatal("intermediate node forwarded 01xxxxxx extension instead of removing it")
				}
				if !bytes.Equal(m.GetHeader().Payload, []byte{0x45}) {
					t.Fatal("relay lost inner payload")
				}
				if typ&0xc0 != 0x40 && (len(m.GetHeader().ExtensionHeaders) != 1 || m.GetHeader().ExtensionHeaders[0].Type != typ) {
					t.Fatal("relay lost required forwarded extension")
				}
			default:
				t.Fatal("relay produced no output")
			}
		})
	}
}

func TestInvalidFTEIDDoesNotReserve(t *testing.T) {
	u := NewUPlaneConn(&net.UDPAddr{})
	i := u.NewFTEID(0, "not-an-ip", "")
	if i == nil || i.Err() == nil {
		t.Fatal("expected a construction error")
	}
	n := 0
	u.iteiMap.syncMap.Range(func(_, _ any) bool { n++; return true })
	if n != 0 {
		t.Fatalf("failed F-TEID construction leaked %d TEID reservation", n)
	}
}

func TestReadFromGTPExpiredDeadline(t *testing.T) {
	u, err := NewUPlaneConnFromUDPConn(listenUDP(t, "127.0.0.1", 0))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = u.Close() })
	// A queued packet avoids any dependence on scheduling or timeout duration.
	// SetReadDeadline promises to affect future Read calls, even with queued data.
	u.tpduCh = make(chan *tpduSet, 1)
	u.tpduCh <- &tpduSet{teid: 1, payload: []byte{0x45}}
	if err := u.SetReadDeadline(time.Unix(1, 0)); err != nil {
		t.Fatal(err)
	}
	n, _, _, err := u.ReadFromGTP(make([]byte, 1))
	if !errors.Is(err, os.ErrDeadlineExceeded) {
		t.Fatalf("expired ReadFromGTP deadline returned n=%d, err=%v; want deadline exceeded", n, err)
	}
}

func TestRelayRebuildsExtensionChain(t *testing.T) {
	for _, flags := range []uint8{0x34, 0x35, 0x36, 0x37} {
		for _, types := range [][]uint8{{0x7f}, {0x11, 0x7f, 0x86}, {0x7f, 0x11, 0x40}, {0x11, 0x86}} {
			var extensions []*message.ExtensionHeader
			var retained []uint8
			for n, typ := range types {
				next := uint8(0)
				if n+1 < len(types) {
					next = types[n+1]
				}
				extensions = append(extensions, message.NewExtensionHeader(typ, []byte{0x12, 0x34}, next))
				if typ&0xc0 != 0x40 {
					retained = append(retained, typ)
				}
			}
			h := message.NewHeaderWithExtensionHeaders(flags, 255, 1, 123, []byte{0x45, 0x67}, extensions...)
			h.NPDUNumber = 9
			raw, err := h.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := message.ParseHeader(raw)
			if err != nil {
				t.Fatal(err)
			}
			forwarded := relayExtensions(raw, parsed)
			got, err := message.ParseHeader(forwarded)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got.Payload, h.Payload) || got.SequenceNumber != parsed.SequenceNumber || got.NPDUNumber != parsed.NPDUNumber || len(got.ExtensionHeaders) != len(retained) {
				t.Fatalf("relay changed data: %x -> %x", raw, forwarded)
			}
			for n, e := range got.ExtensionHeaders {
				if e.Type != retained[n] || !bytes.Equal(e.Content, []byte{0x12, 0x34}) {
					t.Fatalf("relay changed extension %d: %+v", n, e)
				}
			}
			if len(retained) == 0 && got.HasExtensionHeader() {
				t.Fatal("retained E flag after removing all extensions")
			}
		}
	}
}

func TestReadFromGTPDeadlineChanges(t *testing.T) {
	u, err := NewUPlaneConnFromUDPConn(listenUDP(t, "127.0.0.1", 0))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = u.Close() })
	if err := u.SetReadDeadline(time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	results := make(chan error, 2)
	for range 2 {
		go func() { _, _, _, err := u.ReadFromGTP(make([]byte, 1)); results <- err }()
	}
	if err := u.SetDeadline(time.Unix(1, 0)); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		select {
		case err := <-results:
			if !errors.Is(err, os.ErrDeadlineExceeded) {
				t.Fatalf("pending read: %v", err)
			}
		case <-time.After(time.Second):
			t.Fatal("deadline change did not wake pending read")
		}
	}
	if err := u.SetReadDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}
	go func() { _, _, _, err := u.ReadFromGTP(make([]byte, 1)); results <- err }()
	if err := u.SetReadDeadline(time.Now().Add(20 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-results:
		if !errors.Is(err, os.ErrDeadlineExceeded) {
			t.Fatalf("timer expiry: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("new deadline did not expire pending read")
	}
	// An expired deadline must not consume queued data; a later deadline resumes reads.
	u.tpduCh = make(chan *tpduSet, 1)
	u.tpduCh <- &tpduSet{teid: 7, payload: []byte{0x45}}
	if _, _, _, err := u.ReadFromGTP(make([]byte, 1)); !errors.Is(err, os.ErrDeadlineExceeded) {
		t.Fatalf("expired deadline: %v", err)
	}
	if err := u.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	p := make([]byte, 1)
	if n, _, teid, err := u.ReadFromGTP(p); err != nil || n != 1 || teid != 7 || p[0] != 0x45 {
		t.Fatalf("resumed read: %d %d %x %v", n, teid, p, err)
	}
}

func TestServeResumesAfterReadDeadline(t *testing.T) {
	u, err := NewUPlaneConnFromUDPConn(listenUDP(t, "127.0.0.1", 0))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = u.Close() })
	u.DisableErrorIndication()
	if err := u.SetReadDeadline(time.Unix(1, 0)); err != nil {
		t.Fatal(err)
	}
	served := make(chan error, 1)
	go func() { served <- u.serve(context.Background()) }()
	if err := u.SetReadDeadline(time.Now().Add(20 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := u.ReadFromGTP(make([]byte, 1)); !errors.Is(err, os.ErrDeadlineExceeded) {
		t.Fatalf("expiry: %v", err)
	}
	if err := u.SetReadDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}
	c, err := net.Dial("udp", u.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	b, err := message.NewTPDU(7, []byte{0x45}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Write(b); err != nil {
		t.Fatal(err)
	}
	if err := u.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	p := make([]byte, 1)
	if n, _, teid, err := u.ReadFromGTP(p); err != nil || n != 1 || teid != 7 || p[0] != 0x45 {
		t.Fatalf("resumed serve: %d %d %x %v", n, teid, p, err)
	}
	if err := u.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-served:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("serve did not stop")
	}
}
