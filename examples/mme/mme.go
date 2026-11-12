// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package main

import (
	"fmt"
	"net"

	"github.com/gomaja-labs/go-gtp/examples/internal/gtpc"
	"github.com/gomaja-labs/go-gtp/gtpv2"
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

func handleCreateSessionResponse(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
	loggerCh <- fmt.Sprintf("Received %s from %s", msg.MessageTypeName(), sgwAddr)

	// find the session associated with TEID
	session, err := c.GetSessionByTEID(msg.TEID())
	if err != nil {
		c.RemoveSession(session)
		return err
	}
	bearer := session.GetDefaultBearer()

	// assert type to refer to the struct field specific to the message.
	// in general, no need to check if it can be type-asserted, as long as the MessageType is
	// specified correctly in AddHandler().
	csRspFromSGW := msg.(*message.CreateSessionResponse)

	// check Cause value first.
	if ie := csRspFromSGW.Cause; ie != nil {
		cause, err := ie.Cause()
		if err != nil {
			return err
		}
		if cause != gtpv2.CauseRequestAccepted {
			c.RemoveSession(session)
			return &gtpv2.CauseNotOKError{
				MsgType: csRspFromSGW.MessageTypeName(),
				Cause:   cause,
				Msg:     fmt.Sprintf("subscriber: %s", session.IMSI),
			}
		}
	} else {
		return &gtpv2.RequiredIEMissingError{Type: msg.MessageType()}
	}

	if ie := csRspFromSGW.PAA; ie != nil {
		bearer.SubscriberIP, err = ie.IPAddress()
		if err != nil {
			return err
		}
	}
	if senderIE := csRspFromSGW.SenderFTEIDC; senderIE != nil {
		teid, err := senderIE.TEID()
		if err != nil {
			return err
		}
		session.AddTEID(gtpv2.IFTypeS11S4SGWGTPC, teid)
	} else {
		return &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
	}

	s11sgwTEID, err := session.GetTEID(gtpv2.IFTypeS11S4SGWGTPC)
	if err != nil {
		c.RemoveSession(session)
		return err
	}
	s11mmeTEID, err := session.GetTEID(gtpv2.IFTypeS11MMEGTPC)
	if err != nil {
		c.RemoveSession(session)
		return err
	}

	if err := gtpc.BearerContexts(csRspFromSGW.BearerContextsCreated); err != nil {
		return err
	}
	if brCtxIE := csRspFromSGW.BearerContextsCreated; brCtxIE != nil {
		for _, childIE := range brCtxIE[0].ChildIEs {
			if childIE == nil {
				continue
			}
			switch childIE.Type {
			case ie.EPSBearerID:
				bearer.EBI, err = childIE.EPSBearerID()
				if err != nil {
					return err
				}
			case ie.FullyQualifiedTEID:
				if childIE.Instance() != 0 {
					continue
				}
				it, err := childIE.InterfaceType()
				if err != nil {
					return err
				}
				teid, err := childIE.TEID()
				if err != nil {
					return err
				}
				if it == gtpv2.IFTypeS1USGWGTPU {
					// TS 29.274 Table 7.2.2-2 carries the S1-U SGW F-TEID.
					ip, err := gtpc.FTEIDAddress(childIE)
					if err != nil {
						return err
					}
					bearer.SetRemoteAddress(&net.UDPAddr{IP: net.ParseIP(ip), Port: 2152})
					bearer.SetOutgoingTEID(teid)
				}

				session.AddTEID(it, teid)
			}
		}
	} else {
		return &gtpv2.RequiredIEMissingError{Type: ie.BearerContext}
	}

	if bearer.RemoteAddress() == nil {
		return &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
	}
	if err := session.Activate(); err != nil {
		c.RemoveSession(session)
		return err
	}

	createdCh <- session.Subscriber.IMSI
	loggerCh <- fmt.Sprintf(
		"Session created with S-GW for Subscriber: %s;\n\tS11 S-GW: %s, TEID->: %#x, TEID<-: %#x",
		session.Subscriber.IMSI, sgwAddr, s11sgwTEID, s11mmeTEID,
	)
	return nil
}

func handleModifyBearerResponse(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
	loggerCh <- fmt.Sprintf("Received %s from %s", msg.MessageTypeName(), sgwAddr)

	session, err := c.GetSessionByTEID(msg.TEID())
	if err != nil {
		return err
	}

	mbRspFromSGW := msg.(*message.ModifyBearerResponse)
	if causeIE := mbRspFromSGW.Cause; causeIE != nil {
		cause, err := causeIE.Cause()
		if err != nil {
			return err
		}
		if cause != gtpv2.CauseRequestAccepted {
			return &gtpv2.CauseNotOKError{
				MsgType: msg.MessageTypeName(),
				Cause:   cause,
				Msg:     fmt.Sprintf("subscriber: %s", session.IMSI),
			}
		}
	} else {
		return &gtpv2.RequiredIEMissingError{Type: ie.Cause}
	}

	mock := &mockUEeNB{
		subscriberIP: session.GetDefaultBearer().SubscriberIP,
		payload:      payload,
	}
	mock.raddr = session.GetDefaultBearer().RemoteAddress()
	mock.teidOut = session.GetDefaultBearer().OutgoingTEID()
	if mock.raddr == nil {
		return &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
	}

	go mock.run(errCh)

	loggerCh <- fmt.Sprintf("Bearer modified with S-GW for Subscriber: %s", session.IMSI)
	return nil
}

func handleDeleteSessionResponse(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
	loggerCh <- fmt.Sprintf("Received %s from %s", msg.MessageTypeName(), sgwAddr)

	session, err := c.GetSessionByTEID(msg.TEID())
	if err != nil {
		return err
	}

	c.RemoveSession(session)
	delWG.Done()
	loggerCh <- fmt.Sprintf("Session deleted with S-GW for Subscriber: %s", session.IMSI)
	return nil
}
