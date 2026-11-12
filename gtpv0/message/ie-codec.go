package message

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
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
	if messageLen(h, ies)-20 > 65535 {
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

// marshalGenericIEs is structural: preserve order and received invalid contents.
func marshalGenericIEs(h *Header, b []byte, ies []*ie.IE, received []byte) error {
	n := messageLen(h, ies)
	if h == nil || n-20 > 65535 {
		return ErrInvalidLength
	}
	if len(b) < n {
		return ErrTooShortToMarshal
	}
	p := make([]byte, n-(h.MarshalLen()-len(h.Payload)))
	off := 0
	for _, i := range ies {
		if i == nil {
			continue
		}
		var err error
		if received == nil {
			err = i.MarshalTo(p[off:])
		} else {
			err = i.MarshalToLossless(p[off:])
		}
		if err != nil {
			return err
		}
		off += i.MarshalLen()
	}
	copyHeader := *h
	copyHeader.Payload = p
	copyHeader.SetLength()
	return copyHeader.marshalTo(b, received)
}

// ErrIEOutOfSequence identifies descending IE types in a typed message.
var ErrIEOutOfSequence = errors.New("IE out of sequence")

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
	var mandatory []*ie.IE
	require := func(t uint8, n int) error {
		if len(values[t]) < n {
			return fmt.Errorf("message %d IE %d: %w", typ, t, ErrMissingIE)
		}
		mandatory = append(mandatory, values[t][:n]...)
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
		needed = []uint8{ie.QualityOfServiceProfile, ie.SelectionMode, ie.FlowLabelDataI, ie.FlowLabelSignalling, ie.EndUserAddress, ie.AccessPointName, ie.MSISDN}
		if err := require(ie.GSNAddress, 2); err != nil {
			return err
		}
	case MsgTypeUpdatePDPContextRequest:
		needed = []uint8{ie.QualityOfServiceProfile, ie.FlowLabelDataI, ie.FlowLabelSignalling}
		if err := require(ie.GSNAddress, 2); err != nil {
			return err
		}
	}
	for _, t := range needed {
		if err := require(t, 1); err != nil {
			return err
		}
	}
	// TS 29.060 §§11.1, 11.1.5–11.1.7; GSM 09.60 §§10.1, 10.1.5–10.1.7:
	// complete all presence checks before checking mandatory contents.
	occurrences := map[uint8]int{}
	for _, i := range mandatory {
		occurrence := occurrences[i.Type]
		occurrences[i.Type]++
		if err := i.Err(); err != nil {
			return &InvalidIEError{Type: i.Type, Occurrence: occurrence, Err: err}
		}
	}
	// TS 29.060 §11.1; GSM 09.60 §10.1: missing/incorrect mandatory
	// IEs take priority over ordering (§11.1.10 / §10.1.10).
	return checkIEOrder(typ, ies)
}

// checkIEOrder reports IEs that are not in ascending type order (TS 29.060
// §§7.7.0, 11.1.10; GSM 09.60 §§7.9, 10.1.10). Repetitions of a type are in order.
func checkIEOrder(typ uint8, ies []*ie.IE) error {
	for j := 1; j < len(ies); j++ {
		if ies[j].Type < ies[j-1].Type {
			return fmt.Errorf("message %d IE %d after IE %d: %w", typ, ies[j].Type, ies[j-1].Type, ErrIEOutOfSequence)
		}
	}
	return nil
}
