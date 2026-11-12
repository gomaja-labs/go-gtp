package main

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2"
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

func TestSGWS1UAddressFromCreateSessionResponse(t *testing.T) {
	peer := &net.UDPAddr{IP: net.IPv4(192, 0, 2, 1), Port: 2123}
	c := gtpv2.NewConn(&net.UDPAddr{}, gtpv2.IFTypeS11MMEGTPC, 0)
	s := gtpv2.NewSession(peer, &gtpv2.Subscriber{IMSI: "001010000000001"})
	c.RegisterSession(10, s)
	m := &mme{created: make(chan struct{}, 1), modified: make(chan struct{}, 1)}
	response := message.NewCreateSessionResponse(10, 1,
		ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
		ie.NewFullyQualifiedTEID(gtpv2.IFTypeS11S4SGWGTPC, 20, "192.0.2.1", ""),
		ie.NewBearerContext(ie.NewEPSBearerID(5), ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
			ie.NewFullyQualifiedTEID(gtpv2.IFTypeS1USGWGTPU, 30, "192.0.2.2", "")))
	if err := m.handleCreateSessionResponse(c, peer, response); err != nil {
		t.Fatal(err)
	}
	if m.sgw.s1uIP != "192.0.2.2" {
		t.Fatalf("S1-U address=%q, want Create Session Response address", m.sgw.s1uIP)
	}
	if err := m.handleModifyBearerResponse(c, peer, message.NewModifyBearerResponse(10, 2, ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil))); err != nil {
		t.Fatalf("Modify Bearer Response without F-TEID: %v", err)
	}
	if m.sgw.s1uIP != "192.0.2.2" {
		t.Fatal("S1-U address lost")
	}
}
func TestCreateResponseIncompleteBearer(t *testing.T) {
	for _, contexts := range [][]*ie.IE{{}, {nil}, {ie.NewBearerContext(ie.NewEPSBearerID(5), ie.NewFullyQualifiedTEID(gtpv2.IFTypeS1USGWGTPU, 3, "", ""))}} {
		t.Run("incomplete", func(t *testing.T) {
			peer := &net.UDPAddr{IP: net.IPv4(192, 0, 2, 1), Port: 2123}
			c := gtpv2.NewConn(&net.UDPAddr{}, gtpv2.IFTypeS11MMEGTPC, 0)
			s := gtpv2.NewSession(peer, &gtpv2.Subscriber{IMSI: "001010000000001"})
			c.RegisterSession(10, s)
			m := &mme{created: make(chan struct{}, 1)}
			rsp := message.NewCreateSessionResponse(10, 1, ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil), ie.NewFullyQualifiedTEID(gtpv2.IFTypeS11S4SGWGTPC, 2, "192.0.2.1", ""))
			rsp.BearerContextsCreated = contexts
			if err := m.handleCreateSessionResponse(c, peer, rsp); err == nil {
				t.Fatal("incomplete bearer accepted")
			}
		})
	}
}

func TestDetachExplicitlyUnimplemented(t *testing.T) {
	m := &mme{}
	response, err := m.Detach(context.Background(), nil)
	if response != nil || status.Code(err) != codes.Unimplemented {
		t.Fatalf("Detach: %v %v", response, err)
	}
}
