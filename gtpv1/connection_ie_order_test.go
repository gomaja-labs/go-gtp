package gtpv1

import (
	"context"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
)

type parseLogHandler struct{ records chan string }

func (h parseLogHandler) Enabled(context.Context, slog.Level) bool { return true }
func (h parseLogHandler) Handle(_ context.Context, r slog.Record) error {
	h.records <- r.Message
	return nil
}
func (h parseLogHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h parseLogHandler) WithGroup(string) slog.Handler      { return h }

func TestUPlaneDiscardsAndLogsOutOfSequenceIEs(t *testing.T) {
	logs := make(chan string, 16)
	previous := logger.Load()
	SetLogger(slog.New(parseLogHandler{records: logs}))
	t.Cleanup(func() { SetLogger(previous) })
	for _, typ := range []byte{message.MsgTypeEchoRequest, message.MsgTypeErrorIndication, message.MsgTypeEndMarker} {
		t.Run(message.NewGeneric(typ, 0, 1).MessageTypeName(), func(t *testing.T) {
			handled := make(chan uint16, 4)
			sock := listenUDP(t, "127.0.0.1", 0)
			u := serveUDP(t, sock, func(u *UPlaneConn) {
				u.AddHandler(typ, func(_ Conn, _ net.Addr, m message.Message) error { handled <- m.Sequence(); return nil })
			})
			sender := listenUDP(t, "127.0.0.1", 0)
			ies := []*ie.IE{ie.NewRecovery(0), ie.NewCause(128)}
			ordered := []*ie.IE{ies[1], ies[0]}
			if typ == message.MsgTypeErrorIndication {
				ies = []*ie.IE{ie.NewGSNAddress("192.0.2.1"), ie.NewTEIDDataI(1)}
				ordered = []*ie.IE{ies[1], ies[0]}
			}
			// Generic emits a deliberately descending sequence; typed receive must discard it.
			sendUDP(t, sender, u.LocalAddr(), message.NewGeneric(typ, 0, 1, ies...))
			select {
			case got := <-handled:
				t.Fatalf("dispatched invalid datagram, sequence %d", got)
			case entry := <-logs:
				if !strings.Contains(entry, "error parsing GTP:") || !strings.Contains(entry, message.ErrIEOutOfSequence.Error()) {
					t.Fatalf("log = %s", entry)
				}
			case <-time.After(time.Second):
				t.Fatal("no parse error logged")
			}
			// A subsequent valid message proves the receiver continues and gives
			// the handler a deterministic barrier for the discard assertion.
			sendUDP(t, sender, u.LocalAddr(), message.NewGeneric(typ, 0, 2, ordered...))
			select {
			case seq := <-handled:
				if seq != 2 {
					t.Fatalf("dispatched sequence %d, want 2", seq)
				}
			case <-time.After(time.Second):
				t.Fatal("valid message not dispatched")
			}
			if err := sender.SetReadDeadline(time.Now().Add(30 * time.Millisecond)); err != nil {
				t.Fatal(err)
			}
			if _, _, err := sender.ReadFromUDP(make([]byte, 1024)); err == nil {
				t.Fatal("discarded datagram generated a response")
			}
		})
	}
}
