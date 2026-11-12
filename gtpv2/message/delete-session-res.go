// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// DeleteSessionResponse is a DeleteSessionResponse Header and its IEs above.
type DeleteSessionResponse struct {
	*Header
	Cause                         *ie.IE
	Recovery                      *ie.IE
	PCO                           *ie.IE
	IndicationFlags               *ie.IE
	PGWNodeLoadControlInformation *ie.IE
	PGWAPNLoadControlInformation  []*ie.IE
	SGWNodeLoadControlInformation *ie.IE
	PGWOverloadControlInformation []*ie.IE
	SGWOverloadControlInformation *ie.IE
	EPCO                          *ie.IE
	APNRateControlStatus          *ie.IE
	PrivateExtension              *ie.IE
	AdditionalIEs                 []*ie.IE
}

// NewDeleteSessionResponse creates a new DeleteSessionResponse.
func NewDeleteSessionResponse(teid, seq uint32, ies ...*ie.IE) *DeleteSessionResponse {
	d := &DeleteSessionResponse{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeDeleteSessionResponse, teid, seq, nil,
		),
	}

	for _, i := range ies {
		d.addIE(i)
	}

	d.SetLength()
	return d
}

// Marshal serializes DeleteSessionResponse into bytes.
func (d *DeleteSessionResponse) Marshal() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes DeleteSessionResponse into bytes.
func (d *DeleteSessionResponse) MarshalTo(b []byte) error {
	if d == nil || d.Header == nil {
		return ErrNilMessage
	}
	if err := d.Header.validate(MsgTypeDeleteSessionResponse, d.MarshalLen()); err != nil {
		return err
	}
	if len(b) < d.MarshalLen() {
		return ErrTooShortToMarshal
	}
	d.Header.Payload = nil
	d.Header.Payload = make([]byte, d.MarshalLen()-d.Header.MarshalLen())
	offset := 0
	if i := d.Cause; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "Cause", iePair{ie.Cause, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.Recovery; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "Recovery", iePair{ie.Recovery, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.PCO; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "PCO", iePair{ie.ProtocolConfigurationOptions, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.IndicationFlags; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "IndicationFlags", iePair{ie.Indication, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.PGWNodeLoadControlInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "PGWNodeLoadControlInformation", iePair{ie.LoadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range d.PGWAPNLoadControlInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "PGWAPNLoadControlInformation", iePair{ie.LoadControlInformation, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.SGWNodeLoadControlInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SGWNodeLoadControlInformation", iePair{ie.LoadControlInformation, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range d.PGWOverloadControlInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "PGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.SGWOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.EPCO; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "EPCO", iePair{ie.ExtendedProtocolConfigurationOptions, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.APNRateControlStatus; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "APNRateControlStatus", iePair{ie.APNRateControlStatus, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.PrivateExtension; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "PrivateExtension", iePair{ie.PrivateExtension, 255}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range d.AdditionalIEs {
		if i == nil {
			continue
		}
		if err := i.MarshalTo(d.Payload[offset:]); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	d.Header.SetLength()
	return d.Header.MarshalTo(b)
}

// ParseDeleteSessionResponse decodes given bytes as DeleteSessionResponse.
func ParseDeleteSessionResponse(b []byte) (*DeleteSessionResponse, error) {
	d := &DeleteSessionResponse{}
	if err := d.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return d, err
		}
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary decodes given bytes as DeleteSessionResponse.
func (d *DeleteSessionResponse) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrNilMessage
	}
	*d = DeleteSessionResponse{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeDeleteSessionResponse {
		return &UnexpectedTypeError{Expected: MsgTypeDeleteSessionResponse, Actual: h.Type}
	}
	d.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, d.addIE)
}

// MarshalLen returns the serial length in int.
func (d *DeleteSessionResponse) MarshalLen() int {
	if d == nil || d.Header == nil {
		return 0
	}
	l := d.Header.MarshalLen() - len(d.Header.Payload)
	if i := d.Cause; i != nil {
		l += i.MarshalLen()
	}
	if i := d.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := d.PCO; i != nil {
		l += i.MarshalLen()
	}
	if i := d.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := d.PGWNodeLoadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.PGWAPNLoadControlInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := d.SGWNodeLoadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.PGWOverloadControlInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := d.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := d.EPCO; i != nil {
		l += i.MarshalLen()
	}
	if i := d.APNRateControlStatus; i != nil {
		l += i.MarshalLen()
	}
	if i := d.PrivateExtension; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.AdditionalIEs {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	return l
}

// SetLength sets the length in Length field.
func (d *DeleteSessionResponse) SetLength() {
	if d != nil && d.Header != nil {
		d.Header.setLength(d.MarshalLen() - d.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (d *DeleteSessionResponse) MessageTypeName() string {
	return "Delete Session Response"
}

// TEID returns the TEID in uint32.
func (d *DeleteSessionResponse) TEID() uint32 {
	if d == nil {
		return 0
	}
	return d.Header.teid()
}

func (d *DeleteSessionResponse) GetHeader() *Header {
	if d == nil {
		return nil
	}
	return d.Header
}

// TS 29.274 Table 7.2.10.1-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (d *DeleteSessionResponse) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		d.Header.err = err
	}
	switch i.Type {
	case ie.Cause:
		switch i.Instance() {
		case 0:
			if d.Cause == nil {
				d.Cause = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.Recovery:
		switch i.Instance() {
		case 0:
			if d.Recovery == nil {
				d.Recovery = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.ProtocolConfigurationOptions:
		switch i.Instance() {
		case 0:
			if d.PCO == nil {
				d.PCO = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.Indication:
		switch i.Instance() {
		case 0:
			if d.IndicationFlags == nil {
				d.IndicationFlags = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.LoadControlInformation:
		switch i.Instance() {
		case 0:
			if d.PGWNodeLoadControlInformation == nil {
				d.PGWNodeLoadControlInformation = i
			}
		case 1:
			d.PGWAPNLoadControlInformation = append(d.PGWAPNLoadControlInformation, i)
		case 2:
			if d.SGWNodeLoadControlInformation == nil {
				d.SGWNodeLoadControlInformation = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			d.PGWOverloadControlInformation = append(d.PGWOverloadControlInformation, i)
		case 1:
			if d.SGWOverloadControlInformation == nil {
				d.SGWOverloadControlInformation = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.ExtendedProtocolConfigurationOptions:
		switch i.Instance() {
		case 0:
			if d.EPCO == nil {
				d.EPCO = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.APNRateControlStatus:
		switch i.Instance() {
		case 0:
			if d.APNRateControlStatus == nil {
				d.APNRateControlStatus = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.PrivateExtension:
		if d.PrivateExtension == nil {
			d.PrivateExtension = i
		} else if d.PrivateExtension.Instance() != i.Instance() {
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	default:
		d.AdditionalIEs = append(d.AdditionalIEs, i)
	}
}

// Version reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionResponse) Version() int { return d.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionResponse) MessageType() uint8 { return d.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionResponse) Sequence() uint32 { return d.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionResponse) IsPiggybacking() bool { return d.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionResponse) HasTEID() bool { return d.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionResponse) HasMessagePriority() bool { return d.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionResponse) MessagePriority() uint8 { return d.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (d *DeleteSessionResponse) SetSequenceNumber(seq uint32) error {
	return d.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (d *DeleteSessionResponse) SetMessagePriority(priority uint8) error {
	return d.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (d *DeleteSessionResponse) SetTEID(teid uint32) { d.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (d *DeleteSessionResponse) SetPiggybacking(value uint8) { d.GetHeader().SetPiggybacking(value) }
