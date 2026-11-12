// Copyright 2019-2024 go-gtp authors. All rights reservev.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

// DeleteBearerFailureIndication is a DeleteBearerFailureIndication Header and its IEs above.
type DeleteBearerFailureIndication struct {
	*Header
	Cause                         *ie.IE
	BearerContexts                []*ie.IE
	Recovery                      *ie.IE
	IndicationFlags               *ie.IE
	PGWOverloadControlInformation []*ie.IE
	SGWOverloadControlInformation *ie.IE
	PrivateExtension              *ie.IE
	AdditionalIEs                 []*ie.IE
}

// NewDeleteBearerFailureIndication creates a new DeleteBearerFailureIndication.
func NewDeleteBearerFailureIndication(teid, seq uint32, ies ...*ie.IE) *DeleteBearerFailureIndication {
	d := &DeleteBearerFailureIndication{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeDeleteBearerFailureIndication, teid, seq, nil,
		),
	}

	for _, i := range ies {
		d.addIE(i)
	}

	d.SetLength()
	return d
}

// Marshal serializes DeleteBearerFailureIndication into bytes.
func (d *DeleteBearerFailureIndication) Marshal() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes DeleteBearerFailureIndication into bytes.
func (d *DeleteBearerFailureIndication) MarshalTo(b []byte) error {
	if d == nil || d.Header == nil {
		return ErrNilMessage
	}
	if err := d.Header.validate(MsgTypeDeleteBearerFailureIndication, d.MarshalLen()); err != nil {
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
	for _, i := range d.BearerContexts {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "BearerContexts", iePair{ie.BearerContext, 0}); err != nil {
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
	if i := d.IndicationFlags; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "IndicationFlags", iePair{ie.Indication, 0}); err != nil {
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

// ParseDeleteBearerFailureIndication decodes given bytes as DeleteBearerFailureIndication.
func ParseDeleteBearerFailureIndication(b []byte) (*DeleteBearerFailureIndication, error) {
	d := &DeleteBearerFailureIndication{}
	if err := d.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return d, err
		}
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary decodes given bytes as DeleteBearerFailureIndication.
func (d *DeleteBearerFailureIndication) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrNilMessage
	}
	*d = DeleteBearerFailureIndication{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeDeleteBearerFailureIndication {
		return &UnexpectedTypeError{Expected: MsgTypeDeleteBearerFailureIndication, Actual: h.Type}
	}
	d.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, d.addIE)
}

// MarshalLen returns the serial length in int.
func (d *DeleteBearerFailureIndication) MarshalLen() int {
	if d == nil || d.Header == nil {
		return 0
	}
	l := d.Header.MarshalLen() - len(d.Header.Payload)
	if i := d.Cause; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.BearerContexts {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := d.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := d.IndicationFlags; i != nil {
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
func (d *DeleteBearerFailureIndication) SetLength() {
	if d != nil && d.Header != nil {
		d.Header.setLength(d.MarshalLen() - d.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (d *DeleteBearerFailureIndication) MessageTypeName() string {
	return "Delete Bearer Failure Indication"
}

// TEID returns the TEID in uint32.
func (d *DeleteBearerFailureIndication) TEID() uint32 {
	if d == nil {
		return 0
	}
	return d.Header.teid()
}

func (d *DeleteBearerFailureIndication) GetHeader() *Header {
	if d == nil {
		return nil
	}
	return d.Header
}

// TS 29.274 Table 7.2.17.2-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (d *DeleteBearerFailureIndication) addIE(i *ie.IE) {
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
	case ie.BearerContext:
		switch i.Instance() {
		case 0:
			d.BearerContexts = append(d.BearerContexts, i)
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
	case ie.Indication:
		switch i.Instance() {
		case 0:
			if d.IndicationFlags == nil {
				d.IndicationFlags = i
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
func (d *DeleteBearerFailureIndication) Version() int { return d.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerFailureIndication) MessageType() uint8 { return d.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerFailureIndication) Sequence() uint32 { return d.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerFailureIndication) IsPiggybacking() bool { return d.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerFailureIndication) HasTEID() bool { return d.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerFailureIndication) HasMessagePriority() bool {
	return d.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerFailureIndication) MessagePriority() uint8 {
	return d.GetHeader().MessagePriority()
}

// SetSequenceNumber updates the common header safely on a nil receiver.
func (d *DeleteBearerFailureIndication) SetSequenceNumber(seq uint32) error {
	return d.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (d *DeleteBearerFailureIndication) SetMessagePriority(priority uint8) error {
	return d.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (d *DeleteBearerFailureIndication) SetTEID(teid uint32) { d.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (d *DeleteBearerFailureIndication) SetPiggybacking(value uint8) {
	d.GetHeader().SetPiggybacking(value)
}
