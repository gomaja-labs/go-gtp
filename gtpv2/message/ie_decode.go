// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// mandatoryRoles covers the mandatory roles in the implemented message tables.
// Presence and session/interface-dependent conditional requirements remain the
// caller's responsibility; this checks the contents of supplied mandatory IEs.
var mandatoryRoles = map[uint8][]iePair{
	// TS 29.274 Table 7.3.14-1.
	MsgTypeChangeNotificationRequest: {{ie.RATType, 0}},
	// TS 29.274 Table 7.3.15-1.
	MsgTypeChangeNotificationResponse: {{ie.Cause, 0}},
	// TS 29.274 Table 7.3.7-1.
	MsgTypeContextAcknowledge: {{ie.Cause, 0}},
	// TS 29.274 Table 7.3.6-1.
	MsgTypeContextResponse: {{ie.Cause, 0}},
	// TS 29.274 Table 7.2.3-1.
	MsgTypeCreateBearerRequest: {{ie.EPSBearerID, 0}, {ie.BearerContext, 0}},
	// TS 29.274 Table 7.2.4-1.
	MsgTypeCreateBearerResponse: {{ie.Cause, 0}, {ie.BearerContext, 0}},
	// TS 29.274 Table 7.2.1-1.
	MsgTypeCreateSessionRequest: {{ie.RATType, 0}, {ie.FullyQualifiedTEID, 0}, {ie.AccessPointName, 0}, {ie.BearerContext, 0}},
	// TS 29.274 Table 7.2.2-1.
	MsgTypeCreateSessionResponse: {{ie.Cause, 0}, {ie.BearerContext, 0}},
	// TS 29.274 Table 7.2.17.1-1.
	MsgTypeDeleteBearerCommand: {{ie.BearerContext, 0}},
	// TS 29.274 Table 7.2.17.2-1.
	MsgTypeDeleteBearerFailureIndication: {{ie.Cause, 0}, {ie.BearerContext, 0}},
	// TS 29.274 Table 7.2.10.2-1.
	MsgTypeDeleteBearerResponse: {{ie.Cause, 0}},
	// TS 29.274 Table 7.9.2.
	MsgTypeDeletePDNConnectionSetResponse: {{ie.Cause, 0}},
	// TS 29.274 Table 7.2.10.1-1.
	MsgTypeDeleteSessionResponse: {{ie.Cause, 0}},
	// TS 29.274 Table 7.3.13-1.
	MsgTypeDetachAcknowledge: {{ie.Cause, 0}},
	// TS 29.274 Table 7.3.12-1.
	MsgTypeDetachNotification: {{ie.Cause, 0}},
	// TS 29.274 Table 7.2.11.2-1.
	MsgTypeDownlinkDataNotificationAcknowledge: {{ie.Cause, 0}},
	// TS 29.274 Table 7.2.11.3-1.
	MsgTypeDownlinkDataNotificationFailureIndication: {{ie.Cause, 0}},
	// TS 29.274 Table 7.1.2-1.
	MsgTypeEchoResponse: {{ie.Recovery, 0}},
	// TS 29.274 Table 7.2.25-1.
	MsgTypeModifyAccessBearersResponse: {{ie.Cause, 0}},
	// TS 29.274 Table 7.2.14.1-1.
	MsgTypeModifyBearerCommand: {{ie.AggregateMaximumBitRate, 0}, {ie.BearerContext, 0}},
	// TS 29.274 Table 7.2.14.2-1.
	MsgTypeModifyBearerFailureIndication: {{ie.Cause, 0}},
	// TS 29.274 Table 7.2.8-1.
	MsgTypeModifyBearerResponse: {{ie.Cause, 0}},
	// TS 29.274 Table 7.9.5-1.
	MsgTypePGWRestartNotification: {{ie.IPAddress, 0}, {ie.IPAddress, 1}},
	// TS 29.274 Table 7.9.6-1.
	MsgTypePGWRestartNotificationAcknowledge: {{ie.Cause, 0}},
	// TS 29.274 Table 7.2.22-1.
	MsgTypeReleaseAccessBearersResponse: {{ie.Cause, 0}},
	// TS 29.274 Table 7.4.4-1.
	MsgTypeResumeAcknowledge: {{ie.Cause, 0}},
	// TS 29.274 Table 7.4.3-1.
	MsgTypeResumeNotification: {{ie.IMSI, 0}},
	// TS 29.274 Table 7.4.2-1.
	MsgTypeSuspendAcknowledge: {{ie.Cause, 0}},
	// TS 29.274 Table 7.2.15-1.
	MsgTypeUpdateBearerRequest: {{ie.BearerContext, 0}, {ie.AggregateMaximumBitRate, 0}},
	// TS 29.274 Table 7.2.16-1.
	MsgTypeUpdateBearerResponse: {{ie.Cause, 0}, {ie.BearerContext, 0}},
	// TS 29.274 Table 7.9.4-1.
	MsgTypeUpdatePDNConnectionSetResponse: {{ie.Cause, 0}},
}

func offendingIE(i *ie.IE, err error, framing bool) *IEError {
	cause := uint8(69) // TS 29.274 §7.7.8: Mandatory IE incorrect.
	if framing || errors.Is(err, ie.ErrInvalidLength) || errors.Is(err, ie.ErrTooShortToParse) || errors.Is(err, io.ErrUnexpectedEOF) {
		cause = 67 // TS 29.274 §7.7.7: Invalid length.
	}
	return &IEError{Type: i.Type, Instance: i.Instance(), Cause: cause, Framing: framing, Err: err}
}

func parseMessageIEs(h *Header) ([]*ie.IE, error) {
	var ies []*ie.IE
	for b := h.Payload; len(b) > 0; {
		i, err := ie.Parse(b)
		if err != nil {
			// Recover the offending type/instance only when those bytes exist.
			i = &ie.IE{Type: b[0]}
			if len(b) >= 4 {
				i.Length = binary.BigEndian.Uint16(b[1:3])
				i.SetInstance(b[3] & 15)
				i.Payload = b[4:]
			}
			if h.Type == MsgTypeEchoRequest {
				// TS 29.274 §7.1.2 NOTE: answer Echo despite erroneous IEs.
				// No later IE can be delimited after a framing failure.
				h.err = err // Do not transmit the invalid received representation.
				return append(ies, i), nil
			}
			return nil, offendingIE(i, err, true)
		}
		ies = append(ies, i)
		b = b[4+int(i.Length):]
	}
	return ies, nil
}

func addDecodedIEs(h *Header, ies []*ie.IE, add func(*ie.IE)) error {
	// addIE retains construction failures for encoders. Received content errors
	// belong to the retained IE, allowing callers to remove or replace that IE.
	headerErr := h.err
	for _, i := range ies {
		add(i)
	}
	h.err = headerErr
	return validateMandatoryIEs(h, ies)
}

func validateMandatoryIEs(h *Header, ies []*ie.IE) error {
	rejection := false
	switch h.Type {
	case MsgTypeCreateSessionResponse, MsgTypeCreateBearerResponse, MsgTypeUpdateBearerResponse:
		// TS 29.274 §6.1.1: Bearer Context is optional on rejection.
		for _, i := range ies {
			if i.Type == ie.Cause && i.Instance() == 0 {
				cause, err := i.Cause()
				rejection = err == nil && cause >= 64
				break
			}
		}
	}
	seen := make(map[iePair]bool)
	for _, i := range ies {
		pair := iePair{i.Type, i.Instance()}
		if seen[pair] && !i.IsGrouped() {
			continue // TS 29.274 §7.7.10: only the first singleton is handled.
		}
		seen[pair] = true
		for _, required := range mandatoryRoles[h.Type] {
			if rejection && required.typ == ie.BearerContext {
				continue
			}
			if pair != required {
				continue
			}
			if err := mandatoryContent(i); err != nil {
				return offendingIE(i, err, false)
			}
			if i.Type == ie.BearerContext {
				if err := validateBearerContent(h.Type, i); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func mandatoryContent(i *ie.IE) error {
	if err := i.Err(); err != nil {
		return err
	}
	// TS 29.274 §7.7.7: fixed/extendable IE suffixes are ignored on receipt.
	switch i.Type {
	case ie.Cause:
		_, err := i.Cause()
		return err
	case ie.Recovery:
		_, err := i.Recovery()
		return err
	case ie.RATType:
		_, err := i.RATType()
		return err
	case ie.FullyQualifiedTEID:
		_, err := i.FullyQualifiedTEID()
		return err
	case ie.AccessPointName:
		_, err := i.AccessPointName()
		return err
	case ie.EPSBearerID:
		_, err := i.EPSBearerID()
		return err
	case ie.AggregateMaximumBitRate:
		_, err := i.AggregateMaximumBitRateDown()
		return err
	case ie.IPAddress:
		_, err := i.IPAddress()
		return err
	case ie.IMSI:
		_, err := i.IMSI()
		return err
	case ie.BearerQoS:
		_, err := i.BearerQoS()
		return err
	case ie.BearerTFT:
		_, err := i.BearerTFT()
		return err
	}
	return nil
}

func validateBearerContent(messageType uint8, group *ie.IE) error {
	// TS 29.274 Tables 7.2.1-2, 7.2.2-2, 7.2.3-2, 7.2.4-2,
	// 7.2.14.1-2, 7.2.15-2, 7.2.16-2, 7.2.17.1-2, 7.2.17.2-2.
	seen := make(map[iePair]bool)
	for _, child := range group.ChildIEs {
		if child == nil {
			return offendingIE(group, ie.ErrMalformed, false)
		}
		pair := iePair{child.Type, child.Instance()}
		if seen[pair] || child.Instance() != 0 {
			continue
		}
		seen[pair] = true
		required := child.Type == ie.EPSBearerID
		switch messageType {
		case MsgTypeCreateSessionRequest, MsgTypeCreateBearerRequest:
			required = required || child.Type == ie.BearerQoS
		}
		if messageType == MsgTypeCreateBearerRequest {
			required = required || child.Type == ie.BearerTFT
		}
		switch messageType {
		case MsgTypeCreateSessionResponse, MsgTypeCreateBearerResponse,
			MsgTypeUpdateBearerResponse, MsgTypeDeleteBearerFailureIndication:
			required = required || child.Type == ie.Cause
		}
		if required {
			if err := mandatoryContent(child); err != nil {
				return offendingIE(child, err, false)
			}
		}
	}
	return nil
}
