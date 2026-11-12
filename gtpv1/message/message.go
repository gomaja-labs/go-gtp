// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
Package message provides encoding/decoding feature of GTPv1 protocol.

Decoding rejects framing errors. Incorrect optional and unexpected IEs are
retained in their named field or AdditionalIEs; Err and typed accessors report
content errors. Applications must treat these IEs as absent (TS 29.060
§§11.1.8, 11.1.11, 11.1.13; GSM 09.60 §§10.1.8, 10.1.11, 10.1.13).
Incorrect mandatory values return the decoded message and an InvalidIEError,
so callers can respond with Mandatory IE incorrect. Conditional presence that
depends on procedure state must be checked by the caller. Encoding validates
all supplied IEs, including optional and unexpected ones.

PDU Session Container reception accepts unknown Future Extension data while
checking the lengths of known fields (TS 38.415 §5.5.1). Sending retains strict
layout validation.
*/
package message

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

// Message Type definitions.
const (
	_ uint8 = iota
	MsgTypeEchoRequest
	MsgTypeEchoResponse
	MsgTypeVersionNotSupported
	MsgTypeNodeAliveRequest
	MsgTypeNodeAliveResponse
	MsgTypeRedirectionRequest
	MsgTypeRedirectionResponse
	_
	_
	_
	_
	_
	_
	_
	_
	MsgTypeCreatePDPContextRequest // 16
	MsgTypeCreatePDPContextResponse
	MsgTypeUpdatePDPContextRequest
	MsgTypeUpdatePDPContextResponse
	MsgTypeDeletePDPContextRequest
	MsgTypeDeletePDPContextResponse
	MsgTypeCreateAAPDPContextRequest
	MsgTypeCreateAAPDPContextResponse
	MsgTypeDeleteAAPDPContextRequest
	MsgTypeDeleteAAPDPContextResponse
	MsgTypeErrorIndication
	MsgTypePDUNotificationRequest
	MsgTypePDUNotificationResponse
	MsgTypePDUNotificationRejectRequest
	MsgTypePDUNotificationRejectResponse
	MsgTypeSupportedExtensionHeaderNotification
	MsgTypeSendRoutingInfoRequest
	MsgTypeSendRoutingInfoResponse
	MsgTypeFailureReportRequest
	MsgTypeFailureReportResponse
	MsgTypeNoteMSPresentRequest
	MsgTypeNoteMSPresentResponse
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	// TS 29.060 §7.1 Table 1: mobility messages start at type 48.
	MsgTypeIdentificationRequest      uint8 = 48
	MsgTypeIdentificationResponse     uint8 = 49
	MsgTypeSGSNContextRequest         uint8 = 50
	MsgTypeSGSNContextResponse        uint8 = 51
	MsgTypeSGSNContextAcknowledge     uint8 = 52
	MsgTypeDataRecordTransferRequest  uint8 = 240
	MsgTypeDataRecordTransferResponse uint8 = 241
	MsgTypeEndMarker                  uint8 = 254
	MsgTypeTPDU                       uint8 = 255
)

// Message is an interface that defines Message message.
type Message interface {
	GetHeader() *Header
	MarshalTo([]byte) error
	UnmarshalBinary(b []byte) error
	MarshalLen() int
	Version() int
	MessageType() uint8
	MessageTypeName() string
	TEID() uint32
	SetTEID(uint32)
	Sequence() uint16
	SetSequenceNumber(uint16)
}

// Marshal returns the byte sequence generated from a Message instance.
// Better to use MarshalXxx instead if you know the name of message to be serialized.
func Marshal(g Message) ([]byte, error) {
	b := make([]byte, g.MarshalLen())
	if err := g.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// Parse parses the given bytes as a Message.
func Parse(b []byte) (Message, error) {
	if len(b) < 2 {
		return nil, ErrTooShortToParse
	}

	var m Message

	switch b[1] {
	case MsgTypeEchoRequest:
		m = &EchoRequest{}
	case MsgTypeEchoResponse:
		m = &EchoResponse{}
	case MsgTypeCreatePDPContextRequest:
		m = &CreatePDPContextRequest{}
	case MsgTypeCreatePDPContextResponse:
		m = &CreatePDPContextResponse{}
	case MsgTypeUpdatePDPContextRequest:
		m = &UpdatePDPContextRequest{}
	case MsgTypeUpdatePDPContextResponse:
		m = &UpdatePDPContextResponse{}
	case MsgTypeDeletePDPContextRequest:
		m = &DeletePDPContextRequest{}
	case MsgTypeVersionNotSupported:
		m = &VersionNotSupported{}
	case MsgTypeDeletePDPContextResponse:
		m = &DeletePDPContextResponse{}
	/* TODO: Implement!
	case MsgTypeNodeAliveRequest:
		m = &NodeAliveReq{}
	case MsgTypeNodeAliveResponse:
		m = &NodeAliveRes{}
	case MsgTypeRedirectionRequest:
		m = &RedirectionReq{}
	case MsgTypeRedirectionResponse:
		m = &RedirectionRes{}
	case MsgTypeCreateAaPDPContextRequest:
		m = &CreateAaPDPContextReq{}
	case MsgTypeCreateAaPDPContextResponse:
		m = &CreateAaPDPContextRes{}
	case MsgTypeDeleteAaPDPContextRequest:
		m = &DeleteAaPDPContextReq{}
	case MsgTypeDeleteAaPDPContextResponse:
		m = &DeleteAaPDPContextRes{}
	*/
	case MsgTypeErrorIndication:
		m = &ErrorIndication{}
	/* TODO: Implement!
	case MsgTypePduNotificationRequest:
		m = &PduNotificationReq{}
	case MsgTypePduNotificationResponse:
		m = &PduNotificationRes{}
	case MsgTypePduNotificationRejectRequest:
		m = &PduNotificationRejectReq{}
	case MsgTypePduNotificationRejectResponse:
		m = &PduNotificationRejectRes{}
	*/
	case MsgTypeSupportedExtensionHeaderNotification:
		m = &SupportedExtensionHeaderNotification{}
	/* TODO: Implement!
	case MsgTypeSendRoutingInfoRequest:
		m = &SendRoutingInfoReq{}
	case MsgTypeSendRoutingInfoResponse:
		m = &SendRoutingInfoRes{}
	case MsgTypeFailureReportRequest:
		m = &FailureReportReq{}
	case MsgTypeFailureReportResponse:
		m = &FailureReportRes{}
	case MsgTypeNoteMsPresentRequest:
		m = &NoteMsPresentReq{}
	case MsgTypeNoteMsPresentResponse:
		m = &NoteMsPresentRes{}
	case MsgTypeIdentificationRequest:
		m = &IdentificationReq{}
	case MsgTypeIdentificationResponse:
		m = &IdentificationRes{}
	case MsgTypeSgsnContextRequest:
		m = &SgsnContextReq{}
	case MsgTypeSgsnContextResponse:
		m = &SgsnContextRes{}
	case MsgTypeSgsnContextAcknowledge:
		m = &SgsnContextAck{}
	case MsgTypeDataRecordTransferRequest:
		m = &DataRecordTransferReq{}
	case MsgTypeDataRecordTransferResponse:
		m = &DataRecordTransferRes{}
	*/
	case MsgTypeEndMarker:
		m = &EndMarker{}
	case MsgTypeTPDU:
		m = &TPDU{}
	default:
		m = &Generic{}
	}

	if err := m.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return m, err
		}
		return nil, err
	}
	return m, nil
}

// Prettify returns a Message in prettified representation in string.
//
// Note that this relies much on reflect package, and thus the frequent use of
// this function may have a serious impact on the performance of your software.
func Prettify(m Message) string {
	name := m.MessageTypeName()
	header := strings.TrimSuffix(fmt.Sprint(m), "}")

	v := reflect.Indirect(reflect.ValueOf(m))
	n := v.NumField() - 1
	fields := make([]*field, n)
	for i := 1; i < n+1; i++ { // Skip *Header
		fields[i-1] = &field{name: v.Type().Field(i).Name, maybeIE: v.Field(i).Interface()}
	}

	return fmt.Sprintf("{%s: %s, IEs: [%v]}", name, header, strings.Join(prettifyFields(fields), ", "))
}

type field struct {
	name    string
	maybeIE interface{}
}

func prettifyFields(fields []*field) []string {
	ret := []string{}
	for _, field := range fields {
		if field.maybeIE == nil {
			ret = append(ret, prettifyIE(field.name, nil))
			continue
		}

		// TODO: do this recursively?
		v, ok := field.maybeIE.(*ie.IE)
		if !ok {
			// only for AdditionalIEs field
			if ies, ok := field.maybeIE.([]*ie.IE); ok {
				vals := make([]string, len(ies))
				for i, val := range ies {
					vals[i] = fmt.Sprint(val)
				}
				ret = append(ret, fmt.Sprintf("{%s: [%v]}", field.name, strings.Join(vals, ", ")))
			}
			continue
		}

		ret = append(ret, prettifyIE(field.name, v))
	}

	return ret
}

func prettifyIE(name string, i *ie.IE) string {
	return fmt.Sprintf("{%s: %v}", name, i)
}
