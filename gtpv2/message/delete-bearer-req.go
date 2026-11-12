// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// DeleteBearerRequest is a DeleteBearerRequest Header and its IEs above.
type DeleteBearerRequest struct {
	*Header
	LinkedEBI                     *ie.IE
	EBIs                          []*ie.IE
	FailedBearerContext           []*ie.IE
	PTI                           *ie.IE
	PCO                           *ie.IE
	PGWFQCSID                     *ie.IE
	SGWFQCSID                     *ie.IE
	Cause                         *ie.IE
	IndicationFlags               *ie.IE
	PGWNodeLoadControlInformation *ie.IE
	PGWAPNLoadControlInformation  []*ie.IE
	SGWNodeLoadControlInformation *ie.IE
	PGWOverloadControlInformation []*ie.IE
	SGWOverloadControlInformation *ie.IE
	NBIFOMContainer               *ie.IE
	APNRateControlStatus          *ie.IE
	EPCO                          *ie.IE
	PrivateExtension              *ie.IE
	AdditionalIEs                 []*ie.IE
}

// NewDeleteBearerRequest creates a new DeleteBearerRequest.
func NewDeleteBearerRequest(teid, seq uint32, ies ...*ie.IE) *DeleteBearerRequest {
	d := &DeleteBearerRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeDeleteBearerRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		d.addIE(i)
	}

	d.SetLength()
	return d
}

// Marshal serializes DeleteBearerRequest into bytes.
func (d *DeleteBearerRequest) Marshal() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes DeleteBearerRequest into bytes.
func (d *DeleteBearerRequest) MarshalTo(b []byte) error {
	if d == nil || d.Header == nil {
		return ErrNilMessage
	}
	if err := d.Header.validate(MsgTypeDeleteBearerRequest, d.MarshalLen()); err != nil {
		return err
	}
	if len(b) < d.MarshalLen() {
		return ErrTooShortToMarshal
	}
	d.Header.Payload = nil
	d.Header.Payload = make([]byte, d.MarshalLen()-d.Header.MarshalLen())
	offset := 0
	if i := d.LinkedEBI; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "LinkedEBI", iePair{ie.EPSBearerID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range d.EBIs {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "EBIs", iePair{ie.EPSBearerID, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range d.FailedBearerContext {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "FailedBearerContext", iePair{ie.BearerContext, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.PTI; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "PTI", iePair{ie.ProcedureTransactionID, 0}); err != nil {
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
	if i := d.PGWFQCSID; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "PGWFQCSID", iePair{ie.FullyQualifiedCSID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.SGWFQCSID; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SGWFQCSID", iePair{ie.FullyQualifiedCSID, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.Cause; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "Cause", iePair{ie.Cause, 0}); err != nil {
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
	if i := d.NBIFOMContainer; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "NBIFOMContainer", iePair{ie.FContainer, 0}); err != nil {
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
	if i := d.EPCO; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "EPCO", iePair{ie.ExtendedProtocolConfigurationOptions, 0}); err != nil {
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

// ParseDeleteBearerRequest decodes given bytes as DeleteBearerRequest.
func ParseDeleteBearerRequest(b []byte) (*DeleteBearerRequest, error) {
	d := &DeleteBearerRequest{}
	if err := d.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return d, err
		}
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary decodes given bytes as DeleteBearerRequest.
func (d *DeleteBearerRequest) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrNilMessage
	}
	*d = DeleteBearerRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeDeleteBearerRequest {
		return &UnexpectedTypeError{Expected: MsgTypeDeleteBearerRequest, Actual: h.Type}
	}
	d.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, d.addIE)
}

// MarshalLen returns the serial length in int.
func (d *DeleteBearerRequest) MarshalLen() int {
	if d == nil || d.Header == nil {
		return 0
	}
	l := d.Header.MarshalLen() - len(d.Header.Payload)
	if i := d.LinkedEBI; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.EBIs {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range d.FailedBearerContext {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := d.PTI; i != nil {
		l += i.MarshalLen()
	}
	if i := d.PCO; i != nil {
		l += i.MarshalLen()
	}
	if i := d.PGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := d.SGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := d.Cause; i != nil {
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
	if i := d.NBIFOMContainer; i != nil {
		l += i.MarshalLen()
	}
	if i := d.APNRateControlStatus; i != nil {
		l += i.MarshalLen()
	}
	if i := d.EPCO; i != nil {
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
func (d *DeleteBearerRequest) SetLength() {
	if d != nil && d.Header != nil {
		d.Header.setLength(d.MarshalLen() - d.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (d *DeleteBearerRequest) MessageTypeName() string {
	return "Delete Bearer Request"
}

// TEID returns the TEID in uint32.
func (d *DeleteBearerRequest) TEID() uint32 {
	if d == nil {
		return 0
	}
	return d.Header.teid()
}

func (d *DeleteBearerRequest) GetHeader() *Header {
	if d == nil {
		return nil
	}
	return d.Header
}

// TS 29.274 Table 7.2.9.2-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (d *DeleteBearerRequest) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		d.Header.err = err
	}
	switch i.Type {
	case ie.EPSBearerID:
		switch i.Instance() {
		case 0:
			if d.LinkedEBI == nil {
				d.LinkedEBI = i
			}
		case 1:
			d.EBIs = append(d.EBIs, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.BearerContext:
		switch i.Instance() {
		case 0:
			d.FailedBearerContext = append(d.FailedBearerContext, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.ProcedureTransactionID:
		switch i.Instance() {
		case 0:
			if d.PTI == nil {
				d.PTI = i
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
	case ie.FullyQualifiedCSID:
		switch i.Instance() {
		case 0:
			if d.PGWFQCSID == nil {
				d.PGWFQCSID = i
			}
		case 1:
			if d.SGWFQCSID == nil {
				d.SGWFQCSID = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.Cause:
		switch i.Instance() {
		case 0:
			if d.Cause == nil {
				d.Cause = i
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
	case ie.FContainer:
		switch i.Instance() {
		case 0:
			if d.NBIFOMContainer == nil {
				d.NBIFOMContainer = i
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
	case ie.ExtendedProtocolConfigurationOptions:
		switch i.Instance() {
		case 0:
			if d.EPCO == nil {
				d.EPCO = i
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
func (d *DeleteBearerRequest) Version() int { return d.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerRequest) MessageType() uint8 { return d.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerRequest) Sequence() uint32 { return d.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerRequest) IsPiggybacking() bool { return d.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerRequest) HasTEID() bool { return d.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerRequest) HasMessagePriority() bool { return d.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerRequest) MessagePriority() uint8 { return d.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (d *DeleteBearerRequest) SetSequenceNumber(seq uint32) error {
	return d.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (d *DeleteBearerRequest) SetMessagePriority(priority uint8) error {
	return d.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (d *DeleteBearerRequest) SetTEID(teid uint32) { d.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (d *DeleteBearerRequest) SetPiggybacking(value uint8) { d.GetHeader().SetPiggybacking(value) }
