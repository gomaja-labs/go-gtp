// Copyright 2019-2024 go-gtp authors. All rights reservev.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

// ReleaseAccessBearersRequest is a ReleaseAccessBearersRequest Header and its IEs above.
type ReleaseAccessBearersRequest struct {
	*Header
	ListOfRABs                  []*ie.IE
	OriginatingNode             *ie.IE
	IndicationFlags             *ie.IE
	SecondaryRATUsageDataReport []*ie.IE
	PrivateExtension            *ie.IE
	AdditionalIEs               []*ie.IE
}

// NewReleaseAccessBearersRequest creates a new ReleaseAccessBearersRequest.
func NewReleaseAccessBearersRequest(teid, seq uint32, ies ...*ie.IE) *ReleaseAccessBearersRequest {
	r := &ReleaseAccessBearersRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeReleaseAccessBearersRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		r.addIE(i)
	}

	r.SetLength()
	return r
}

// Marshal serializes ReleaseAccessBearersRequest into bytes.
func (r *ReleaseAccessBearersRequest) Marshal() ([]byte, error) {
	b := make([]byte, r.MarshalLen())
	if err := r.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ReleaseAccessBearersRequest into bytes.
func (r *ReleaseAccessBearersRequest) MarshalTo(b []byte) error {
	if r == nil || r.Header == nil {
		return ErrNilMessage
	}
	if err := r.Header.validate(MsgTypeReleaseAccessBearersRequest, r.MarshalLen()); err != nil {
		return err
	}
	if len(b) < r.MarshalLen() {
		return ErrTooShortToMarshal
	}
	r.Header.Payload = nil
	r.Header.Payload = make([]byte, r.MarshalLen()-r.Header.MarshalLen())
	offset := 0
	for _, i := range r.ListOfRABs {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, r.Payload[offset:], "ListOfRABs", iePair{ie.EPSBearerID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := r.OriginatingNode; i != nil {
		if err := marshalNamedIE(i, r.Payload[offset:], "OriginatingNode", iePair{ie.NodeType, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := r.IndicationFlags; i != nil {
		if err := marshalNamedIE(i, r.Payload[offset:], "IndicationFlags", iePair{ie.Indication, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range r.SecondaryRATUsageDataReport {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, r.Payload[offset:], "SecondaryRATUsageDataReport", iePair{ie.SecondaryRATUsageDataReport, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := r.PrivateExtension; i != nil {
		if err := marshalNamedIE(i, r.Payload[offset:], "PrivateExtension", iePair{ie.PrivateExtension, 255}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range r.AdditionalIEs {
		if i == nil {
			continue
		}
		if err := i.MarshalTo(r.Payload[offset:]); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	r.Header.SetLength()
	return r.Header.MarshalTo(b)
}

// ParseReleaseAccessBearersRequest decodes given bytes as ReleaseAccessBearersRequest.
func ParseReleaseAccessBearersRequest(b []byte) (*ReleaseAccessBearersRequest, error) {
	r := &ReleaseAccessBearersRequest{}
	if err := r.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return r, err
		}
		return nil, err
	}
	return r, nil
}

// UnmarshalBinary decodes given bytes as ReleaseAccessBearersRequest.
func (r *ReleaseAccessBearersRequest) UnmarshalBinary(b []byte) error {
	if r == nil {
		return ErrNilMessage
	}
	*r = ReleaseAccessBearersRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeReleaseAccessBearersRequest {
		return &UnexpectedTypeError{Expected: MsgTypeReleaseAccessBearersRequest, Actual: h.Type}
	}
	r.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, r.addIE)
}

// MarshalLen returns the serial length in int.
func (r *ReleaseAccessBearersRequest) MarshalLen() int {
	if r == nil || r.Header == nil {
		return 0
	}
	l := r.Header.MarshalLen() - len(r.Header.Payload)
	for _, i := range r.ListOfRABs {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := r.OriginatingNode; i != nil {
		l += i.MarshalLen()
	}
	if i := r.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range r.SecondaryRATUsageDataReport {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := r.PrivateExtension; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range r.AdditionalIEs {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	return l
}

// SetLength sets the length in Length field.
func (r *ReleaseAccessBearersRequest) SetLength() {
	if r != nil && r.Header != nil {
		r.Header.setLength(r.MarshalLen() - r.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (r *ReleaseAccessBearersRequest) MessageTypeName() string {
	return "Release Access Bearers Request"
}

// TEID returns the TEID in uint32.
func (r *ReleaseAccessBearersRequest) TEID() uint32 {
	if r == nil {
		return 0
	}
	return r.Header.teid()
}

func (r *ReleaseAccessBearersRequest) GetHeader() *Header {
	if r == nil {
		return nil
	}
	return r.Header
}

// TS 29.274 Table 7.2.21-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (r *ReleaseAccessBearersRequest) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		r.Header.err = err
	}
	switch i.Type {
	case ie.EPSBearerID:
		switch i.Instance() {
		case 0:
			r.ListOfRABs = append(r.ListOfRABs, i)
		default:
			r.AdditionalIEs = append(r.AdditionalIEs, i)
		}
	case ie.NodeType:
		switch i.Instance() {
		case 0:
			if r.OriginatingNode == nil {
				r.OriginatingNode = i
			}
		default:
			r.AdditionalIEs = append(r.AdditionalIEs, i)
		}
	case ie.Indication:
		switch i.Instance() {
		case 0:
			if r.IndicationFlags == nil {
				r.IndicationFlags = i
			}
		default:
			r.AdditionalIEs = append(r.AdditionalIEs, i)
		}
	case ie.SecondaryRATUsageDataReport:
		switch i.Instance() {
		case 0:
			r.SecondaryRATUsageDataReport = append(r.SecondaryRATUsageDataReport, i)
		default:
			r.AdditionalIEs = append(r.AdditionalIEs, i)
		}
	case ie.PrivateExtension:
		if r.PrivateExtension == nil {
			r.PrivateExtension = i
		} else if r.PrivateExtension.Instance() != i.Instance() {
			r.AdditionalIEs = append(r.AdditionalIEs, i)
		}
	default:
		r.AdditionalIEs = append(r.AdditionalIEs, i)
	}
}

// Version reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersRequest) Version() int { return r.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersRequest) MessageType() uint8 { return r.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersRequest) Sequence() uint32 { return r.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersRequest) IsPiggybacking() bool { return r.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersRequest) HasTEID() bool { return r.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersRequest) HasMessagePriority() bool {
	return r.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersRequest) MessagePriority() uint8 { return r.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (r *ReleaseAccessBearersRequest) SetSequenceNumber(seq uint32) error {
	return r.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (r *ReleaseAccessBearersRequest) SetMessagePriority(priority uint8) error {
	return r.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (r *ReleaseAccessBearersRequest) SetTEID(teid uint32) { r.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (r *ReleaseAccessBearersRequest) SetPiggybacking(value uint8) {
	r.GetHeader().SetPiggybacking(value)
}
