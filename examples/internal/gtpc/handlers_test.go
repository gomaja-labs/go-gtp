package gtpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv2"
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

func TestDeleteBearerRequestCleanup(t *testing.T) {
	for _, linked := range []bool{false, true} {
		t.Run(map[bool]string{true: "PDN connection", false: "dedicated bearers"}[linked], func(t *testing.T) {
			peer, err := net.ListenPacket("udp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = peer.Close() }()
			c := gtpv2.NewConn(&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, gtpv2.IFTypeS11MMEGTPC, 0)
			if err = c.Listen(context.Background()); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = c.Close() }()
			s := gtpv2.NewSession(peer.LocalAddr(), &gtpv2.Subscriber{IMSI: "001010000000001"})
			s.GetDefaultBearer().EBI = 5
			s.AddBearer("dedicated-6", gtpv2.NewBearer(6, "internet", nil))
			s.AddBearer("dedicated-7", gtpv2.NewBearer(7, "internet", nil))
			s.AddTEID(gtpv2.IFTypeS11S4SGWGTPC, 20)
			c.RegisterSession(10, s)
			ies := []*ie.IE{ie.NewEPSBearerID(6).WithInstance(1), ie.NewEPSBearerID(7).WithInstance(1)}
			if linked {
				ies = []*ie.IE{ie.NewEPSBearerID(5)}
			}
			if err = HandleDeleteBearerRequest(c, peer.LocalAddr(), message.NewDeleteBearerRequest(10, 99, ies...)); err != nil {
				t.Fatal(err)
			}
			if err = peer.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
				t.Fatal(err)
			}
			raw := make([]byte, 1500)
			n, _, err := peer.ReadFrom(raw)
			if err != nil {
				t.Fatal(err)
			}
			msg, err := message.Parse(raw[:n])
			if err != nil {
				t.Fatal(err)
			}
			rsp := msg.(*message.DeleteBearerResponse)
			if rsp.Sequence() != 99 || rsp.TEID() != 20 || rsp.Cause.MustCause() != gtpv2.CauseRequestAccepted {
				t.Fatalf("wrong reply: %v", rsp)
			}
			if linked {
				if len(c.Sessions()) != 0 {
					t.Fatal("PDN connection retained")
				}
			} else {
				if s.BearerCount() != 1 {
					t.Fatal("dedicated bearers retained")
				}
				if _, err = c.GetSessionByTEID(10); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestIncompleteRequestsReturnCause(t *testing.T) {
	for _, req := range []message.Message{
		message.NewCreateSessionRequest(0, 1),
		message.NewCreateSessionRequest(0, 1, ie.NewBearerContext(ie.NewEPSBearerID(5))),
		message.NewModifyBearerRequest(0, 1),
		message.NewModifyBearerRequest(0, 1, ie.NewBearerContext()),
		message.NewModifyBearerRequest(0, 1, ie.NewBearerContext(ie.NewEPSBearerID(5), ie.NewFullyQualifiedTEID(0, 1, "", ""))),
		message.NewDeleteBearerRequest(0, 1),
		message.NewDeleteBearerRequest(0, 1, ie.NewEPSBearerID(5), ie.NewEPSBearerID(6).WithInstance(1)),
	} {
		t.Run(req.MessageTypeName(), func(t *testing.T) {
			peer, err := net.ListenPacket("udp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = peer.Close() }()
			c := gtpv2.NewConn(&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)}, 10, 0)
			if err = c.Listen(context.Background()); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = c.Close() }()
			if err = CheckRequest(c, peer.LocalAddr(), req); err == nil {
				t.Fatal("incomplete request accepted")
			}
			if err = peer.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
				t.Fatal(err)
			}
			raw := make([]byte, 1500)
			n, _, err := peer.ReadFrom(raw)
			if err != nil {
				t.Fatal(err)
			}
			header, err := message.ParseHeader(raw[:n])
			if err != nil {
				t.Fatal(err)
			}
			ies, err := ie.ParseMultiIEs(header.Payload)
			if err != nil {
				t.Fatal(err)
			}
			if header.Sequence() != 1 || len(ies) != 1 || ies[0].Type != ie.Cause || ies[0].MustCause() < 64 {
				t.Fatalf("wrong failure reply %x", raw[:n])
			}
		})
	}
}
