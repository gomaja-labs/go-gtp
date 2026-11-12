package message

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
)

func messageLen(h *Header, ies []*ie.IE) int {
	if h == nil {
		return 0
	}
	n := h.MarshalLen() - len(h.Payload)
	for _, i := range ies {
		if i != nil {
			n += i.MarshalLen()
		}
	}
	return n
}

func marshalIEs(h *Header, b []byte, ies []*ie.IE) error {
	if h == nil {
		return ErrInvalidLength
	}
	if messageLen(h, ies)-8 > 65535 {
		return ErrInvalidLength
	}
	if len(b) < messageLen(h, ies) {
		return ErrTooShortToMarshal
	}
	// TS 29.060 §7.7.0; GSM 09.60 §7.9: ascending type order, stable repetitions.
	ordered := make([]*ie.IE, 0, len(ies))
	for _, i := range ies {
		if i != nil {
			if err := i.Err(); err != nil {
				return err
			}
			ordered = append(ordered, i)
		}
	}
	sort.SliceStable(ordered, func(a, b int) bool { return ordered[a].Type < ordered[b].Type })
	p := make([]byte, messageLen(h, ordered)-(h.MarshalLen()-len(h.Payload)))
	off := 0
	for _, i := range ordered {
		if err := i.MarshalTo(p[off:]); err != nil {
			return err
		}
		off += i.MarshalLen()
	}
	if err := validateMandatory(h.Type, p); err != nil {
		return err
	}
	copyHeader := *h
	copyHeader.Payload = p
	copyHeader.SetLength()
	return copyHeader.MarshalTo(b)
}

// ErrMissingIE identifies a message without a mandatory information element.
var ErrMissingIE = errors.New("missing mandatory IE")

// InvalidIEError identifies a mandatory IE whose contents are incorrect.
// Occurrence is zero-based among IEs of Type. GTPv0/v1 IEs have no instance field.
// The decoded message remains available so a caller can send the appropriate Cause.
type InvalidIEError struct {
	Type       uint8
	Occurrence int
	Err        error
}

func (e *InvalidIEError) Error() string {
	return fmt.Sprintf("mandatory IE %d occurrence %d: %v", e.Type, e.Occurrence, e.Err)
}

func (e *InvalidIEError) Unwrap() error { return e.Err }

func isIEContentError(err error) bool {
	var invalid *InvalidIEError
	return errors.As(err, &invalid)
}

func validateMandatory(typ uint8, p []byte) error {
	ies, err := ie.ParseMultiIEs(p)
	if err != nil {
		return err
	}
	values := map[uint8][]*ie.IE{}
	for _, i := range ies {
		values[i.Type] = append(values[i.Type], i)
	}
	require := func(t uint8, n int) error {
		if len(values[t]) < n {
			return fmt.Errorf("message %d IE %d: %w", typ, t, ErrMissingIE)
		}
		// TS 29.060 §11.1.7; GSM 09.60 §10.1.7: incorrect mandatory
		// values require a Cause response, while optional errors do not.
		for occurrence, i := range values[t][:n] {
			if err := i.Err(); err != nil {
				return &InvalidIEError{Type: t, Occurrence: occurrence, Err: err}
			}
		}
		return nil
	}
	var needed []uint8
	// TS 29.060 §7 Tables 3–12; TS 29.281 §§7.2.2, 7.2.3, 7.3.1;
	// GSM 09.60 §7 Tables 3–9: mandatory presence is checked before dispatch.
	switch typ {
	case MsgTypeEchoResponse:
		needed = []uint8{ie.Recovery}
	case MsgTypeCreatePDPContextResponse, MsgTypeUpdatePDPContextResponse, MsgTypeDeletePDPContextResponse:
		needed = []uint8{ie.Cause}
	case MsgTypeCreatePDPContextRequest:
		needed = []uint8{ie.TEIDDataI, ie.NSAPI, ie.QoSProfile}
		if err := require(ie.GSNAddress, 2); err != nil {
			return err
		}
	case MsgTypeUpdatePDPContextRequest, MsgTypeDeletePDPContextRequest:
		needed = []uint8{ie.NSAPI}
	case MsgTypeErrorIndication:
		needed = []uint8{ie.TEIDDataI, ie.GSNAddress}
	case MsgTypeSupportedExtensionHeaderNotification:
		needed = []uint8{ie.ExtensionHeaderTypeList}
	}
	for _, t := range needed {
		if err := require(t, 1); err != nil {
			return err
		}
	}
	return nil
}
