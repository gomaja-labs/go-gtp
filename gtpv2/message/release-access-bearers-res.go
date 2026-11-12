// Copyright 2019-2024 go-gtp authors. All rights reservev.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

// ReleaseAccessBearersResponse is a ReleaseAccessBearersResponse Header and its IEs above.
type ReleaseAccessBearersResponse struct {
	*Header
	Cause                         *ie.IE
	Recovery                      *ie.IE
	IndicationFlags               *ie.IE
	SGWNodeLoadControlInformation *ie.IE
	SGWOverloadControlInformation *ie.IE
	PrivateExtension              *ie.IE
	AdditionalIEs                 []*ie.IE
}

// NewReleaseAccessBearersResponse creates a new ReleaseAccessBearersResponse.
func NewReleaseAccessBearersResponse(teid, seq uint32, ies ...*ie.IE) *ReleaseAccessBearersResponse {
	r := &ReleaseAccessBearersResponse{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeReleaseAccessBearersResponse, teid, seq, nil,
		),
	}

	for _, i := range ies {
		r.addIE(i)
	}

	r.SetLength()
	return r
}

// Marshal serializes ReleaseAccessBearersResponse into bytes.
func (r *ReleaseAccessBearersResponse) Marshal() ([]byte, error) {
	b := make([]byte, r.MarshalLen())
	if err := r.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ReleaseAccessBearersResponse into bytes.
func (r *ReleaseAccessBearersResponse) MarshalTo(b []byte) error {
	if r == nil || r.Header == nil {
		return ErrNilMessage
	}
	if err := r.Header.validate(MsgTypeReleaseAccessBearersResponse, r.MarshalLen()); err != nil {
		return err
	}
	if len(b) < r.MarshalLen() {
		return ErrTooShortToMarshal
	}
	r.Header.Payload = nil
	r.Header.Payload = make([]byte, r.MarshalLen()-r.Header.MarshalLen())
	offset := 0
	if i := r.Cause; i != nil {
		if err := marshalNamedIE(i, r.Payload[offset:], "Cause", iePair{ie.Cause, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := r.Recovery; i != nil {
		if err := marshalNamedIE(i, r.Payload[offset:], "Recovery", iePair{ie.Recovery, 0}); err != nil {
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
	if i := r.SGWNodeLoadControlInformation; i != nil {
		if err := marshalNamedIE(i, r.Payload[offset:], "SGWNodeLoadControlInformation", iePair{ie.LoadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := r.SGWOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, r.Payload[offset:], "SGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
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

// ParseReleaseAccessBearersResponse decodes given bytes as ReleaseAccessBearersResponse.
func ParseReleaseAccessBearersResponse(b []byte) (*ReleaseAccessBearersResponse, error) {
	r := &ReleaseAccessBearersResponse{}
	if err := r.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return r, err
		}
		return nil, err
	}
	return r, nil
}

// UnmarshalBinary decodes given bytes as ReleaseAccessBearersResponse.
func (r *ReleaseAccessBearersResponse) UnmarshalBinary(b []byte) error {
	if r == nil {
		return ErrNilMessage
	}
	*r = ReleaseAccessBearersResponse{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeReleaseAccessBearersResponse {
		return &UnexpectedTypeError{Expected: MsgTypeReleaseAccessBearersResponse, Actual: h.Type}
	}
	r.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, r.addIE)
}

// MarshalLen returns the serial length in int.
func (r *ReleaseAccessBearersResponse) MarshalLen() int {
	if r == nil || r.Header == nil {
		return 0
	}
	l := r.Header.MarshalLen() - len(r.Header.Payload)
	if i := r.Cause; i != nil {
		l += i.MarshalLen()
	}
	if i := r.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := r.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := r.SGWNodeLoadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := r.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
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
func (r *ReleaseAccessBearersResponse) SetLength() {
	if r != nil && r.Header != nil {
		r.Header.setLength(r.MarshalLen() - r.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (r *ReleaseAccessBearersResponse) MessageTypeName() string {
	return "Release Access Bearers Response"
}

// TEID returns the TEID in uint32.
func (r *ReleaseAccessBearersResponse) TEID() uint32 {
	if r == nil {
		return 0
	}
	return r.Header.teid()
}

func (r *ReleaseAccessBearersResponse) GetHeader() *Header {
	if r == nil {
		return nil
	}
	return r.Header
}

// TS 29.274 Table 7.2.22-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (r *ReleaseAccessBearersResponse) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		r.Header.err = err
	}
	switch i.Type {
	case ie.Cause:
		switch i.Instance() {
		case 0:
			if r.Cause == nil {
				r.Cause = i
			}
		default:
			r.AdditionalIEs = append(r.AdditionalIEs, i)
		}
	case ie.Recovery:
		switch i.Instance() {
		case 0:
			if r.Recovery == nil {
				r.Recovery = i
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
	case ie.LoadControlInformation:
		switch i.Instance() {
		case 0:
			if r.SGWNodeLoadControlInformation == nil {
				r.SGWNodeLoadControlInformation = i
			}
		default:
			r.AdditionalIEs = append(r.AdditionalIEs, i)
		}
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			if r.SGWOverloadControlInformation == nil {
				r.SGWOverloadControlInformation = i
			}
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
func (r *ReleaseAccessBearersResponse) Version() int { return r.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersResponse) MessageType() uint8 { return r.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersResponse) Sequence() uint32 { return r.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersResponse) IsPiggybacking() bool { return r.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersResponse) HasTEID() bool { return r.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersResponse) HasMessagePriority() bool {
	return r.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (r *ReleaseAccessBearersResponse) MessagePriority() uint8 {
	return r.GetHeader().MessagePriority()
}

// SetSequenceNumber updates the common header safely on a nil receiver.
func (r *ReleaseAccessBearersResponse) SetSequenceNumber(seq uint32) error {
	return r.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (r *ReleaseAccessBearersResponse) SetMessagePriority(priority uint8) error {
	return r.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (r *ReleaseAccessBearersResponse) SetTEID(teid uint32) { r.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (r *ReleaseAccessBearersResponse) SetPiggybacking(value uint8) {
	r.GetHeader().SetPiggybacking(value)
}
