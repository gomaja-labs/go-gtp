// Copyright 2026 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv2

import (
	"encoding/binary"
	"errors"
	"net"
	"reflect"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

type ieRequirement struct{ typ, instance uint8 }

// Mandatory presence only: Conditional IEs require procedure-specific state.
// TS 29.274 V20.0.0 §7.7.6 and the individual message tables below.
var mandatoryIEs = map[uint8][]ieRequirement{
	message.MsgTypeEchoRequest:                               {{ie.Recovery, 0}},                                                                            // Table 7.1.1-1
	message.MsgTypeEchoResponse:                              {{ie.Recovery, 0}},                                                                            // Table 7.1.2-1
	message.MsgTypeVersionNotSupportedIndication:             {},                                                                                            // §5.3
	message.MsgTypeCreateSessionRequest:                      {{ie.RATType, 0}, {ie.FullyQualifiedTEID, 0}, {ie.AccessPointName, 0}, {ie.BearerContext, 0}}, // Table 7.2.1-1
	message.MsgTypeCreateSessionResponse:                     {{ie.Cause, 0}, {ie.BearerContext, 0}},                                                        // Table 7.2.2-1
	message.MsgTypeCreateBearerRequest:                       {{ie.EPSBearerID, 0}, {ie.BearerContext, 0}},                                                  // Table 7.2.3-1
	message.MsgTypeCreateBearerResponse:                      {{ie.Cause, 0}, {ie.BearerContext, 0}},                                                        // Table 7.2.4-1
	message.MsgTypeModifyBearerRequest:                       {},                                                                                            // Table 7.2.7-1
	message.MsgTypeModifyBearerResponse:                      {{ie.Cause, 0}},                                                                               // Table 7.2.8-1
	message.MsgTypeDeleteSessionRequest:                      {},                                                                                            // Table 7.2.9.1-1
	message.MsgTypeDeleteBearerRequest:                       {},                                                                                            // Table 7.2.9.2-1
	message.MsgTypeDeleteSessionResponse:                     {{ie.Cause, 0}},                                                                               // Table 7.2.10.1-1
	message.MsgTypeDeleteBearerResponse:                      {{ie.Cause, 0}},                                                                               // Table 7.2.10.2-1
	message.MsgTypeDownlinkDataNotification:                  {},                                                                                            // Table 7.2.11.1-1
	message.MsgTypeDownlinkDataNotificationAcknowledge:       {{ie.Cause, 0}},                                                                               // Table 7.2.11.2-1
	message.MsgTypeDownlinkDataNotificationFailureIndication: {{ie.Cause, 0}},                                                                               // Table 7.2.11.3-1
	message.MsgTypeBearerResourceCommand:                     {{ie.EPSBearerID, 0}, {ie.ProcedureTransactionID, 0}, {ie.TrafficAggregateDescription, 0}},    // Table 7.2.5-1
	message.MsgTypeBearerResourceFailureIndication:           {{ie.Cause, 0}, {ie.EPSBearerID, 0}, {ie.ProcedureTransactionID, 0}},                          // Table 7.2.6-1

	message.MsgTypeModifyBearerCommand:               {{ie.AggregateMaximumBitRate, 0}, {ie.BearerContext, 0}}, // Table 7.2.14.1-1
	message.MsgTypeModifyBearerFailureIndication:     {{ie.Cause, 0}},                                          // Table 7.2.14.2-1
	message.MsgTypeUpdateBearerRequest:               {{ie.BearerContext, 0}, {ie.AggregateMaximumBitRate, 0}}, // Table 7.2.15-1
	message.MsgTypeUpdateBearerResponse:              {{ie.Cause, 0}, {ie.BearerContext, 0}},                   // Table 7.2.16-1
	message.MsgTypeDeleteBearerCommand:               {{ie.BearerContext, 0}},                                  // Table 7.2.17.1-1
	message.MsgTypeDeleteBearerFailureIndication:     {{ie.Cause, 0}, {ie.BearerContext, 0}},                   // Table 7.2.17.2-1
	message.MsgTypeReleaseAccessBearersRequest:       {},                                                       // Table 7.2.21-1
	message.MsgTypeReleaseAccessBearersResponse:      {{ie.Cause, 0}},                                          // Table 7.2.22-1
	message.MsgTypeStopPagingIndication:              {},                                                       // Table 7.2.23-1
	message.MsgTypeModifyAccessBearersRequest:        {},                                                       // Table 7.2.24-1
	message.MsgTypeModifyAccessBearersResponse:       {{ie.Cause, 0}},                                          // Table 7.2.25-1
	message.MsgTypeContextRequest:                    {},                                                       // Table 7.3.5-1
	message.MsgTypeContextResponse:                   {{ie.Cause, 0}},                                          // Table 7.3.6-1
	message.MsgTypeContextAcknowledge:                {{ie.Cause, 0}},                                          // Table 7.3.7-1
	message.MsgTypeDetachNotification:                {{ie.Cause, 0}},                                          // Table 7.3.12-1
	message.MsgTypeDetachAcknowledge:                 {{ie.Cause, 0}},                                          // Table 7.3.13-1
	message.MsgTypeChangeNotificationRequest:         {{ie.RATType, 0}},                                        // Table 7.3.14-1
	message.MsgTypeChangeNotificationResponse:        {{ie.Cause, 0}},                                          // Table 7.3.15-1
	message.MsgTypeSuspendNotification:               {},                                                       // Table 7.4.1-1
	message.MsgTypeSuspendAcknowledge:                {{ie.Cause, 0}},                                          // Table 7.4.2-1
	message.MsgTypeResumeNotification:                {{ie.IMSI, 0}},                                           // Table 7.4.3-1
	message.MsgTypeResumeAcknowledge:                 {{ie.Cause, 0}},                                          // Table 7.4.4-1
	message.MsgTypeDeletePDNConnectionSetRequest:     {},                                                       // Table 7.9.1-1
	message.MsgTypeDeletePDNConnectionSetResponse:    {{ie.Cause, 0}},                                          // Table 7.9.2
	message.MsgTypeUpdatePDNConnectionSetRequest:     {},                                                       // Table 7.9.3-1
	message.MsgTypeUpdatePDNConnectionSetResponse:    {{ie.Cause, 0}},                                          // Table 7.9.4-1
	message.MsgTypePGWRestartNotification:            {{ie.IPAddress, 0}, {ie.IPAddress, 1}},                   // Table 7.9.5-1
	message.MsgTypePGWRestartNotificationAcknowledge: {{ie.Cause, 0}},                                          // Table 7.9.6-1
}

// knownType accepts codec models and application-defined handlers (TS 29.274 §7.7.4).
func (c *Conn) knownType(typ uint8) bool {
	if _, ok := mandatoryIEs[typ]; ok {
		return true
	}
	handler, ok := c.msgHandlerMap.load(typ)
	return ok && handler != nil
}

// Bearer Context mandatory children, including optional contexts when present.
// TS 29.274 Tables 7.2.1-2/3, 7.2.2-2/3, 7.2.3-2, 7.2.4-2,
// 7.2.7-2/3, 7.2.8-2/3, 7.2.10.2-2, 7.2.14.1-2, 7.2.15-2,
// 7.2.16-2, 7.2.17.1-2, 7.2.17.2-2, 7.2.24-2/3, 7.2.25-2/3.
func bearerRequirements(typ, instance uint8) []ieRequirement {
	ebi := ieRequirement{ie.EPSBearerID, 0}
	switch typ {
	case message.MsgTypeCreateSessionRequest:
		if instance == 0 {
			return []ieRequirement{ebi, {ie.BearerQoS, 0}}
		}
		if instance == 1 {
			return []ieRequirement{ebi}
		}
	case message.MsgTypeCreateBearerRequest:
		if instance == 0 {
			return []ieRequirement{ebi, {ie.BearerTFT, 0}, {ie.BearerQoS, 0}}
		}
	case message.MsgTypeCreateSessionResponse, message.MsgTypeModifyBearerResponse, message.MsgTypeModifyAccessBearersResponse:
		if instance <= 1 {
			return []ieRequirement{ebi, {ie.Cause, 0}}
		}
	case message.MsgTypeCreateBearerResponse, message.MsgTypeUpdateBearerResponse, message.MsgTypeDeleteBearerResponse, message.MsgTypeDeleteBearerFailureIndication:
		if instance == 0 {
			return []ieRequirement{ebi, {ie.Cause, 0}}
		}
	case message.MsgTypeModifyBearerRequest, message.MsgTypeModifyAccessBearersRequest:
		if instance <= 1 {
			return []ieRequirement{ebi}
		}
	case message.MsgTypeModifyBearerCommand, message.MsgTypeUpdateBearerRequest, message.MsgTypeDeleteBearerCommand:
		if instance == 0 {
			return []ieRequirement{ebi}
		}
	case message.MsgTypeContextAcknowledge:
		// TS 29.274 Table 7.3.7-2.
		if instance == 0 {
			return []ieRequirement{ebi, {ie.FullyQualifiedTEID, 0}}
		}
	}
	return nil
}

// messageIEs reads modeled IE fields without marshaling (which mutates headers).
func messageIEs(msg message.Message) []*ie.IE {
	v := reflect.ValueOf(msg)
	if !v.IsValid() || v.IsNil() {
		return nil
	}
	v = v.Elem()
	var ies []*ie.IE
	for i := 0; i < v.NumField(); i++ {
		switch field := v.Field(i).Interface().(type) {
		case *ie.IE:
			if field != nil {
				ies = append(ies, field)
			}
		case []*ie.IE:
			for _, item := range field {
				if item != nil {
					ies = append(ies, item)
				}
			}
		}
	}
	return ies
}
func rejectionCause(msg message.Message) bool {
	for _, i := range messageIEs(msg) {
		if i.Type == ie.Cause && i.Instance() == 0 {
			cause, err := i.Cause()
			return err == nil && cause >= 64
		}
	}
	return false
}
func missingIE(ies []*ie.IE, requirements []ieRequirement) *RequiredIEMissingError {
	for _, req := range requirements {
		found := false
		for _, i := range ies {
			if i != nil && i.Type == req.typ && i.Instance() == req.instance {
				found = true
				break
			}
		}
		if !found {
			return &RequiredIEMissingError{Type: req.typ, Instance: req.instance}
		}
	}
	return nil
}
func (c *Conn) validateMandatory(peer net.Addr, msg message.Message) error {
	ies := messageIEs(msg)
	requirements := mandatoryIEs[msg.MessageType()]
	// TS 29.274 §§6.1.1, 7.7.6: rejection responses need only Cause.
	rejected := isReply(msg.MessageType()) && rejectionCause(msg)
	if rejected {
		requirements = []ieRequirement{{ie.Cause, 0}}
	}
	missing := missingIE(ies, requirements)
	if missing == nil && !rejected {
		missing = validateGrouped(msg.MessageType(), ies, false)
	}

	if missing == nil {
		return nil
	}
	if isReply(msg.MessageType()) {
		c.reportError(&TransactionError{peer, msg.Sequence(), msg.MessageType(), missing})
	} else if msg.MessageType() != message.MsgTypeEchoRequest {
		if err := c.respondWithCause(peer, msg, CauseMandatoryIEMissing, ie.New(missing.Type, missing.Instance, nil)); err != nil {
			return err
		}
	}
	return missing
}

func (c *Conn) respondWithCause(peer net.Addr, req message.Message, cause uint8, offending *ie.IE) error {
	types := replies[req.MessageType()]
	if len(types) == 0 || req.MessageType() == message.MsgTypeEchoRequest || isReply(req.MessageType()) {
		return nil
	}
	typ := types[0]
	if isCommand(req.MessageType()) {
		typ = types[len(types)-1]
	}
	// TS 29.274 §5.5.2 permits TEID zero for protocol-error responses and
	// requires it for Context not found. §7.7.6 includes the missing type/instance.
	res := message.NewGeneric(typ, 0, req.Sequence(), ie.NewCause(cause, 0, 0, 0, offending))
	return c.RespondTo(peer, req, res)
}

func (c *Conn) processDatagram(peer net.Addr, raw []byte) {
	// The codec retains slices, so every datagram needs its own backing bytes.
	raw = append([]byte(nil), raw...)
	msg, err := message.Parse(raw)
	if err != nil {
		c.processDecodeError(peer, raw, err)
		return
	}
	if !c.knownType(msg.MessageType()) {
		return
	}
	handled, handleErr := c.dispatchMessage(peer, msg)
	if handleErr != nil {
		logf("error handling GTP message: %v", handleErr)
	}
	// TS 29.274 §7.7.0: an unmatched, invalid or rejected primary cannot
	// establish the default bearer. Secondary errors never undo the primary.
	if handleErr != nil || rejectionCause(msg) {
		return
	}
	if handled {
		c.acceptedPrimary(peer, msg, true)
	} else if !c.acceptedPrimary(peer, msg, false) {
		return
	}
	h := msg.GetHeader()
	end := int(binary.BigEndian.Uint16(raw[2:4])) + 4
	if h.PiggybackError != nil {
		c.processDecodeError(peer, raw[end:], h.PiggybackError)
	} else if h.Piggybacked != nil {
		if err := c.handleMessage(peer, h.Piggybacked); err != nil {
			logf("error handling piggybacked message: %v", err)
		}
	}
}

func (c *Conn) processDecodeError(peer net.Addr, raw []byte, err error) {
	if err != nil {
		// TS 29.274 §7.7.1: even an erroneous reply stops retransmissions
		// and its parsing error is reported to the application.
		if h := transactionHeader(raw); h != nil && c.matchReply(peer, h.Sequence(), h.MessageType()) {
			c.reportError(&TransactionError{peer, h.Sequence(), h.MessageType(), err})
		}
		var version *message.UnsupportedVersionError
		var length *message.LengthError
		var offending *message.IEError
		switch {
		case errors.As(err, &version):
			// TS 29.274 §7.7.2: reply only to versions higher than two.
			if version.Version > 2 {
				seq := uint32(0)
				off := 4
				if len(raw) > 0 && raw[0]&8 != 0 {
					off = 8
				}
				if len(raw) >= off+4 {
					seq = uint32(raw[off])<<16 | uint32(raw[off+1])<<8 | uint32(raw[off+2])
				}
				res := message.NewVersionNotSupportedIndication(0, seq)
				if b, e := message.Marshal(res); e == nil {
					if _, e = c.WriteTo(b, peer); e != nil {
						c.reportError(e)
					}
				}
			}
		case errors.As(err, &offending):
			// TS 29.274 §§7.7.7–7.7.8: echo the offending IE type/instance.
			if req := transactionHeader(raw); req != nil && c.knownType(req.MessageType()) {
				if e := c.respondWithCause(peer, req, offending.Cause, ie.New(offending.Type, offending.Instance, nil)); e != nil {
					c.reportError(e)
				}
			}
		case errors.As(err, &length):
			// TS 29.274 §7.7.3: inconsistent requests get Invalid Length;
			// responses stop retransmissions and are discarded. Short headers never enter here.
			if req := piggybackHeader(raw); req != nil {
				if e := c.respondWithCause(peer, req, CauseInvalidOverallLengthOfTheTriggeredResponseMessageAndAPiggybackedInitialMessage, nil); e != nil {
					c.reportError(e)
				}
				return
			}
			req := message.NewGeneric(length.MessageType, length.TEID, length.SequenceNumber)
			if c.knownType(length.MessageType) {
				if e := c.respondWithCause(peer, req, CauseInvalidLength, nil); e != nil {
					c.reportError(e)
				}
			}
		}
		return
	}
}

// validateGrouped checks mandatory children of optional grouped IEs too.
// TS 29.274 Tables 7.2.1-4/5, 7.2.2-4/5, 7.2.3-3/4,
// 7.3.6-2/3/4/5 and the corresponding load/overload tables of each message.
func validateGrouped(typ uint8, ies []*ie.IE, inPDN bool) *RequiredIEMissingError {
	for _, i := range ies {
		if i == nil {
			continue
		}
		if inPDN && i.Type != ie.BearerContext && i.Type != ie.RemoteUEContext {
			continue
		}
		if !inPDN {
			known := false
			for _, g := range messageGroups[typ] {
				if g.typ == i.Type && g.instance == i.Instance() {
					known = true
					break
				}
			}
			if !known {
				continue
			}
		}
		var requirements []ieRequirement
		switch i.Type {
		case ie.BearerContext:
			if typ == message.MsgTypeContextResponse && inPDN {
				requirements = []ieRequirement{{ie.EPSBearerID, 0}, {ie.BearerQoS, 0}}
			} else {
				requirements = bearerRequirements(typ, i.Instance())
			}
		case ie.OverloadControlInformation:
			requirements = []ieRequirement{{ie.SequenceNumber, 0}, {ie.Metric, 0}, {ie.EPCTimer, 0}}
		case ie.LoadControlInformation:
			requirements = []ieRequirement{{ie.SequenceNumber, 0}, {ie.Metric, 0}}
		case ie.RemoteUEContext:
			if i.Instance() == 0 {
				requirements = []ieRequirement{{ie.RemoteUserID, 0}, {ie.RemoteUEIPinformation, 0}}
			}
		case ie.PDNConnection:
			if typ == message.MsgTypeContextResponse {
				requirements = []ieRequirement{{ie.AccessPointName, 0}, {ie.EPSBearerID, 0}, {ie.FullyQualifiedTEID, 0}, {ie.BearerContext, 0}, {ie.AggregateMaximumBitRate, 0}}
			}
		case ie.SCEFPDNConnection:
			if typ == message.MsgTypeContextResponse {
				requirements = []ieRequirement{{ie.AccessPointName, 0}, {ie.EPSBearerID, 0}, {ie.NodeIdentifier, 0}}
			}
		}
		if len(requirements) == 0 {
			continue
		}
		if missing := missingIE(i.ChildIEs, requirements); missing != nil {
			return missing
		}
		if i.Type == ie.PDNConnection {
			if missing := validateGrouped(typ, i.ChildIEs, true); missing != nil {
				return missing
			}
		}
	}
	return nil
}

// Group instances defined by each message table; unexpected groups are ignored (§7.7.9).
var messageGroups = map[uint8][]ieRequirement{
	message.MsgTypeCreateSessionRequest:          {{ie.BearerContext, 0}, {ie.BearerContext, 1}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}, {ie.OverloadControlInformation, 2}, {ie.RemoteUEContext, 0}},                                    // Table 7.2.1-1
	message.MsgTypeCreateSessionResponse:         {{ie.BearerContext, 0}, {ie.BearerContext, 1}, {ie.LoadControlInformation, 0}, {ie.LoadControlInformation, 1}, {ie.LoadControlInformation, 2}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}}, // Table 7.2.2-1
	message.MsgTypeCreateBearerRequest:           {{ie.BearerContext, 0}, {ie.LoadControlInformation, 0}, {ie.LoadControlInformation, 1}, {ie.LoadControlInformation, 2}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}},                        // Table 7.2.3-1
	message.MsgTypeCreateBearerResponse:          {{ie.BearerContext, 0}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}, {ie.OverloadControlInformation, 2}},                                                                                    // Table 7.2.4-1
	message.MsgTypeModifyBearerRequest:           {{ie.BearerContext, 0}, {ie.BearerContext, 1}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}, {ie.OverloadControlInformation, 2}},                                                             // Table 7.2.7-1
	message.MsgTypeModifyBearerResponse:          {{ie.BearerContext, 0}, {ie.BearerContext, 1}, {ie.LoadControlInformation, 0}, {ie.LoadControlInformation, 1}, {ie.LoadControlInformation, 2}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}}, // Table 7.2.8-1
	message.MsgTypeDeleteSessionRequest:          {{ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}, {ie.OverloadControlInformation, 2}},                                                                                                           // Table 7.2.9.1-1
	message.MsgTypeDeleteBearerRequest:           {{ie.BearerContext, 0}, {ie.LoadControlInformation, 0}, {ie.LoadControlInformation, 1}, {ie.LoadControlInformation, 2}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}},                        // Table 7.2.9.2-1
	message.MsgTypeDeleteSessionResponse:         {{ie.LoadControlInformation, 0}, {ie.LoadControlInformation, 1}, {ie.LoadControlInformation, 2}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}},                                               // Table 7.2.10.1-1
	message.MsgTypeDeleteBearerResponse:          {{ie.BearerContext, 0}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}, {ie.OverloadControlInformation, 2}},                                                                                    // Table 7.2.10.2-1
	message.MsgTypeDownlinkDataNotification:      {{ie.LoadControlInformation, 0}, {ie.OverloadControlInformation, 0}},                                                                                                                                                   // Table 7.2.11.1-1
	message.MsgTypeModifyBearerCommand:           {{ie.BearerContext, 0}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}, {ie.OverloadControlInformation, 2}},                                                                                    // Table 7.2.14.1-1
	message.MsgTypeModifyBearerFailureIndication: {{ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}},                                                                                                                                               // Table 7.2.14.2-1
	message.MsgTypeUpdateBearerRequest:           {{ie.BearerContext, 0}, {ie.LoadControlInformation, 0}, {ie.LoadControlInformation, 1}, {ie.LoadControlInformation, 2}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}},                        // Table 7.2.15-1
	message.MsgTypeUpdateBearerResponse:          {{ie.BearerContext, 0}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}, {ie.OverloadControlInformation, 2}},                                                                                    // Table 7.2.16-1
	message.MsgTypeDeleteBearerCommand:           {{ie.BearerContext, 0}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}},                                                                                                                        // Table 7.2.17.1-1
	message.MsgTypeDeleteBearerFailureIndication: {{ie.BearerContext, 0}, {ie.OverloadControlInformation, 0}, {ie.OverloadControlInformation, 1}},                                                                                                                        // Table 7.2.17.2-1
	message.MsgTypeReleaseAccessBearersResponse:  {{ie.LoadControlInformation, 0}, {ie.OverloadControlInformation, 0}},                                                                                                                                                   // Table 7.2.22-1
	message.MsgTypeModifyAccessBearersRequest:    {{ie.BearerContext, 0}, {ie.BearerContext, 1}},                                                                                                                                                                         // Table 7.2.24-1
	message.MsgTypeModifyAccessBearersResponse:   {{ie.BearerContext, 0}, {ie.BearerContext, 1}, {ie.LoadControlInformation, 0}, {ie.OverloadControlInformation, 0}},                                                                                                     // Table 7.2.25-1
	message.MsgTypeContextResponse:               {{ie.PDNConnection, 0}, {ie.SCEFPDNConnection, 0}},                                                                                                                                                                     // Table 7.3.6-1
	message.MsgTypeContextAcknowledge:            {{ie.BearerContext, 0}},                                                                                                                                                                                                // Table 7.3.7-1
}

// transactionHeader exposes fields only after the complete header is present.
// It intentionally does not validate Length: these fields correlate parse errors.
func transactionHeader(raw []byte) message.Message {
	if len(raw) < 8 || raw[0]>>5 != 2 {
		return nil
	}
	off := 4
	teid := uint32(0)
	if raw[0]&8 != 0 {
		if len(raw) < 12 {
			return nil
		}
		off = 8
		teid = binary.BigEndian.Uint32(raw[4:8])
	}
	seq := uint32(raw[off])<<16 | uint32(raw[off+1])<<8 | uint32(raw[off+2])
	return message.NewGeneric(raw[1], teid, seq)
}

// TS 29.274 §7.7.3: inconsistent concatenated lengths use Cause 105
// in the response to the piggybacked initial message.
func piggybackHeader(raw []byte) message.Message {
	first := transactionHeader(raw)
	if first == nil || raw[0]&0x10 == 0 {
		return nil
	}
	end := int(binary.BigEndian.Uint16(raw[2:4])) + 4
	if end < 12 || end > len(raw) {
		return nil
	}
	child := transactionHeader(raw[end:])
	if child == nil {
		return nil
	}
	if first.MessageType() == message.MsgTypeCreateSessionResponse && child.MessageType() == message.MsgTypeCreateBearerRequest ||
		first.MessageType() == message.MsgTypeCreateBearerResponse && child.MessageType() == message.MsgTypeModifyBearerRequest {
		return child
	}
	return nil
}
