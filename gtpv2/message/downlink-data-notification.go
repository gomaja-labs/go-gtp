// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// DownlinkDataNotification is a DownlinkDataNotification Header and its IEs above.
type DownlinkDataNotification struct {
	*Header
	Cause                         *ie.IE
	EPSBearerID                   []*ie.IE
	AllocationRetensionPriority   []*ie.IE
	IMSI                          *ie.IE
	SenderFTEIDC                  *ie.IE
	IndicationFlags               *ie.IE
	SGWNodeLoadControlInformation *ie.IE
	SGWOverloadControlInformation *ie.IE
	PagingAndServiceInformation   []*ie.IE
	DLDataPacketsSize             *ie.IE
	PrivateExtension              *ie.IE
	AdditionalIEs                 []*ie.IE
}

// NewDownlinkDataNotification creates a new DownlinkDataNotification.
func NewDownlinkDataNotification(teid, seq uint32, ies ...*ie.IE) *DownlinkDataNotification {
	d := &DownlinkDataNotification{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeDownlinkDataNotification, teid, seq, nil,
		),
	}

	for _, i := range ies {
		d.addIE(i)
	}

	d.SetLength()
	return d
}

// Marshal serializes DownlinkDataNotification into bytes.
func (d *DownlinkDataNotification) Marshal() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes DownlinkDataNotification into bytes.
func (d *DownlinkDataNotification) MarshalTo(b []byte) error {
	if d == nil || d.Header == nil {
		return ErrNilMessage
	}
	if err := d.Header.validate(MsgTypeDownlinkDataNotification, d.MarshalLen()); err != nil {
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
	for _, i := range d.EPSBearerID {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "EPSBearerID", iePair{ie.EPSBearerID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range d.AllocationRetensionPriority {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "AllocationRetensionPriority", iePair{ie.AllocationRetensionPriority, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.IMSI; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "IMSI", iePair{ie.IMSI, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.SenderFTEIDC; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SenderFTEIDC", iePair{ie.FullyQualifiedTEID, 0}); err != nil {
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
	if i := d.SGWNodeLoadControlInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SGWNodeLoadControlInformation", iePair{ie.LoadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.SGWOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range d.PagingAndServiceInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "PagingAndServiceInformation", iePair{ie.PagingAndServiceInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.DLDataPacketsSize; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "DLDataPacketsSize", iePair{ie.IntegerNumber, 0}); err != nil {
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

// ParseDownlinkDataNotification decodes given bytes as DownlinkDataNotification.
func ParseDownlinkDataNotification(b []byte) (*DownlinkDataNotification, error) {
	d := &DownlinkDataNotification{}
	if err := d.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return d, err
		}
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary decodes given bytes as DownlinkDataNotification.
func (d *DownlinkDataNotification) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrNilMessage
	}
	*d = DownlinkDataNotification{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeDownlinkDataNotification {
		return &UnexpectedTypeError{Expected: MsgTypeDownlinkDataNotification, Actual: h.Type}
	}
	d.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, d.addIE)
}

// MarshalLen returns the serial length in int.
func (d *DownlinkDataNotification) MarshalLen() int {
	if d == nil || d.Header == nil {
		return 0
	}
	l := d.Header.MarshalLen() - len(d.Header.Payload)
	if i := d.Cause; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.EPSBearerID {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range d.AllocationRetensionPriority {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := d.IMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := d.SenderFTEIDC; i != nil {
		l += i.MarshalLen()
	}
	if i := d.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := d.SGWNodeLoadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := d.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.PagingAndServiceInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := d.DLDataPacketsSize; i != nil {
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
func (d *DownlinkDataNotification) SetLength() {
	if d != nil && d.Header != nil {
		d.Header.setLength(d.MarshalLen() - d.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (d *DownlinkDataNotification) MessageTypeName() string {
	return "Downlink Data Notification"
}

// TEID returns the TEID in uint32.
func (d *DownlinkDataNotification) TEID() uint32 {
	if d == nil {
		return 0
	}
	return d.Header.teid()
}

func (d *DownlinkDataNotification) GetHeader() *Header {
	if d == nil {
		return nil
	}
	return d.Header
}

// TS 29.274 Table 7.2.11.1-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (d *DownlinkDataNotification) addIE(i *ie.IE) {
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
	case ie.EPSBearerID:
		switch i.Instance() {
		case 0:
			d.EPSBearerID = append(d.EPSBearerID, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.AllocationRetensionPriority:
		switch i.Instance() {
		case 0:
			d.AllocationRetensionPriority = append(d.AllocationRetensionPriority, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.IMSI:
		switch i.Instance() {
		case 0:
			if d.IMSI == nil {
				d.IMSI = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.FullyQualifiedTEID:
		switch i.Instance() {
		case 0:
			if d.SenderFTEIDC == nil {
				d.SenderFTEIDC = i
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
			if d.SGWNodeLoadControlInformation == nil {
				d.SGWNodeLoadControlInformation = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			if d.SGWOverloadControlInformation == nil {
				d.SGWOverloadControlInformation = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.PagingAndServiceInformation:
		switch i.Instance() {
		case 0:
			d.PagingAndServiceInformation = append(d.PagingAndServiceInformation, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.IntegerNumber:
		switch i.Instance() {
		case 0:
			if d.DLDataPacketsSize == nil {
				d.DLDataPacketsSize = i
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
func (d *DownlinkDataNotification) Version() int { return d.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotification) MessageType() uint8 { return d.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotification) Sequence() uint32 { return d.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotification) IsPiggybacking() bool { return d.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotification) HasTEID() bool { return d.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotification) HasMessagePriority() bool {
	return d.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotification) MessagePriority() uint8 { return d.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (d *DownlinkDataNotification) SetSequenceNumber(seq uint32) error {
	return d.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (d *DownlinkDataNotification) SetMessagePriority(priority uint8) error {
	return d.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (d *DownlinkDataNotification) SetTEID(teid uint32) { d.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (d *DownlinkDataNotification) SetPiggybacking(value uint8) { d.GetHeader().SetPiggybacking(value) }
