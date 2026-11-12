// Copyright 2026 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package gtpc contains checks shared by the EPC examples.
package gtpc

import (
	"errors"
	"net"

	"github.com/gomaja-labs/go-gtp/gtpv2"
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

// BearerContexts checks contexts before an example uses their first bearer.
// TS 29.274 bearer context tables require an EBI in each context.
func BearerContexts(contexts []*ie.IE) error {
	if len(contexts) == 0 {
		return &gtpv2.RequiredIEMissingError{Type: ie.BearerContext}
	}
	for _, ctx := range contexts {
		if ctx == nil {
			return &gtpv2.RequiredIEMissingError{Type: ie.BearerContext}
		}
		if _, err := ctx.EPSBearerID(); err != nil {
			return &gtpv2.RequiredIEMissingError{Type: ie.EPSBearerID}
		}
		for _, child := range ctx.ChildIEs {
			if child == nil {
				continue
			}
			if child.Type == ie.FullyQualifiedTEID {
				if _, err := FTEIDAddress(child); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// FTEIDAddress rejects F-TEIDs without an advertised usable address.
func FTEIDAddress(i *ie.IE) (string, error) {
	if i == nil {
		return "", &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
	}
	ip, err := i.IPAddress()
	if err != nil {
		return "", err
	}
	if net.ParseIP(ip) == nil {
		return "", ie.ErrMalformed
	}
	return ip, nil
}

// CheckRequest replies to incomplete input before an example changes session state.
// TS 29.274 §7.7.6 and §7.7.8 require a Cause and offending type/instance.
func CheckRequest(c *gtpv2.Conn, peer net.Addr, msg message.Message) error {
	var err error
	var typ uint8
	offending := &ie.IE{Type: ie.BearerContext}
	switch req := msg.(type) {
	case *message.CreateSessionRequest:
		typ = message.MsgTypeCreateSessionResponse
		if err = BearerContexts(req.BearerContextsToBeCreated); err == nil {
			offending = &ie.IE{Type: ie.FullyQualifiedTEID}
			_, err = FTEIDAddress(req.SenderFTEIDC)
		}
	case *message.ModifyBearerRequest:
		typ = message.MsgTypeModifyBearerResponse
		err = BearerContexts(req.BearerContextsToBeModified)
		if err == nil {
			for _, ctx := range req.BearerContextsToBeModified {
				found := false
				for _, child := range ctx.ChildIEs {
					if child != nil && child.Type == ie.FullyQualifiedTEID {
						found = true
					}
				}
				if !found {
					err = &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
					break
				}
			}
		}
	case *message.DeleteBearerRequest:
		typ = message.MsgTypeDeleteBearerResponse
		_, _, err = deleteBearerIDs(req)
		offending = &ie.IE{Type: ie.EPSBearerID}
	default:
		return nil
	}
	if err == nil {
		return nil
	}
	cause := gtpv2.CauseMandatoryIEIncorrect
	var missing *gtpv2.RequiredIEMissingError
	if errors.As(err, &missing) {
		cause = gtpv2.CauseMandatoryIEMissing
		offending = &ie.IE{Type: missing.Type}
		offending.WithInstance(missing.Instance)
	}
	res := message.NewGeneric(typ, 0, 0, ie.NewCause(cause, 0, 0, 0, offending))
	if sendErr := c.RespondTo(peer, msg, res); sendErr != nil {
		return sendErr
	}
	return err
}

func deleteBearerIDs(req *message.DeleteBearerRequest) ([]uint8, bool, error) {
	if req.LinkedEBI != nil && len(req.EBIs) > 0 {
		return nil, false, ie.ErrMalformed
	}
	ies := req.EBIs
	linked := req.LinkedEBI != nil
	if linked {
		ies = []*ie.IE{req.LinkedEBI}
	}
	if len(ies) == 0 {
		return nil, false, &gtpv2.RequiredIEMissingError{Type: ie.EPSBearerID}
	}
	ids := make([]uint8, 0, len(ies))
	for _, i := range ies {
		ebi, err := i.EPSBearerID()
		if err != nil {
			return nil, false, err
		}
		ids = append(ids, ebi)
	}
	return ids, linked, nil
}

// HandleDeleteBearerRequest acknowledges deletion and then removes local state.
// TS 29.274 §§7.2.9.2, 7.2.10.2: linked EBI deletes the PDN connection;
// instance-one EBIs delete the listed dedicated bearers.
func HandleDeleteBearerRequest(c *gtpv2.Conn, peer net.Addr, msg message.Message) error {
	req, ok := msg.(*message.DeleteBearerRequest)
	if !ok {
		return &gtpv2.UnexpectedTypeError{Msg: msg}
	}
	if err := CheckRequest(c, peer, msg); err != nil {
		return err
	}
	ids, linked, err := deleteBearerIDs(req)
	if err != nil {
		return err
	}
	s, err := c.GetSessionByTEID(req.TEID())
	if err != nil {
		return c.RespondTo(peer, req, message.NewDeleteBearerResponse(0, 0, ie.NewCause(gtpv2.CauseContextNotFound, 0, 0, 0, nil)))
	}
	teid, err := s.GetTEID(gtpv2.IFTypeS11S4SGWGTPC)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if _, err = s.LookupBearerByEBI(id); err != nil {
			return c.RespondTo(peer, req, message.NewDeleteBearerResponse(teid, 0, ie.NewCause(gtpv2.CauseContextNotFound, 0, 0, 0, nil)))
		}
	}
	ies := []*ie.IE{ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil)}
	if linked {
		ies = append(ies, ie.NewEPSBearerID(ids[0]))
	} else {
		for _, id := range ids {
			ies = append(ies, ie.NewBearerContext(ie.NewEPSBearerID(id), ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil)))
		}
	}
	if err = c.RespondTo(peer, req, message.NewDeleteBearerResponse(teid, 0, ies...)); err != nil {
		return err
	}
	if linked {
		c.RemoveSession(s)
	} else {
		for _, id := range ids {
			s.RemoveBearerByEBI(id)
		}
	}
	return nil
}
