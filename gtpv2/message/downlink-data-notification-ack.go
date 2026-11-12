// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// DownlinkDataNotificationAcknowledge is a DownlinkDataNotificationAcknowledge Header and its IEs above.
type DownlinkDataNotificationAcknowledge struct {
	*Header
	Cause                           *ie.IE
	DataNotificationDelay           *ie.IE
	Recovery                        *ie.IE
	DLLowPriorityTrafficThrottling  *ie.IE
	IMSI                            *ie.IE
	DLBufferingDuration             *ie.IE
	DLBufferingSuggestedPacketCount *ie.IE
	PrivateExtension                *ie.IE
	AdditionalIEs                   []*ie.IE
}

// NewDownlinkDataNotificationAcknowledge creates a new DownlinkDataNotificationAcknowledge.
func NewDownlinkDataNotificationAcknowledge(teid, seq uint32, ies ...*ie.IE) *DownlinkDataNotificationAcknowledge {
	d := &DownlinkDataNotificationAcknowledge{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeDownlinkDataNotificationAcknowledge, teid, seq, nil,
		),
	}

	for _, i := range ies {
		d.addIE(i)
	}

	d.SetLength()
	return d
}

// Marshal serializes DownlinkDataNotificationAcknowledge into bytes.
func (d *DownlinkDataNotificationAcknowledge) Marshal() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes DownlinkDataNotificationAcknowledge into bytes.
func (d *DownlinkDataNotificationAcknowledge) MarshalTo(b []byte) error {
	if d == nil || d.Header == nil {
		return ErrNilMessage
	}
	if err := d.Header.validate(MsgTypeDownlinkDataNotificationAcknowledge, d.MarshalLen()); err != nil {
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
	if i := d.DataNotificationDelay; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "DataNotificationDelay", iePair{ie.DelayValue, 0}); err != nil {
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
	if i := d.DLLowPriorityTrafficThrottling; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "DLLowPriorityTrafficThrottling", iePair{ie.Throttling, 0}); err != nil {
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
	if i := d.DLBufferingDuration; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "DLBufferingDuration", iePair{ie.EPCTimer, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.DLBufferingSuggestedPacketCount; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "DLBufferingSuggestedPacketCount", iePair{ie.IntegerNumber, 0}); err != nil {
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

// ParseDownlinkDataNotificationAcknowledge decodes given bytes as DownlinkDataNotificationAcknowledge.
func ParseDownlinkDataNotificationAcknowledge(b []byte) (*DownlinkDataNotificationAcknowledge, error) {
	d := &DownlinkDataNotificationAcknowledge{}
	if err := d.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return d, err
		}
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary decodes given bytes as DownlinkDataNotificationAcknowledge.
func (d *DownlinkDataNotificationAcknowledge) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrNilMessage
	}
	*d = DownlinkDataNotificationAcknowledge{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeDownlinkDataNotificationAcknowledge {
		return &UnexpectedTypeError{Expected: MsgTypeDownlinkDataNotificationAcknowledge, Actual: h.Type}
	}
	d.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, d.addIE)
}

// MarshalLen returns the serial length in int.
func (d *DownlinkDataNotificationAcknowledge) MarshalLen() int {
	if d == nil || d.Header == nil {
		return 0
	}
	l := d.Header.MarshalLen() - len(d.Header.Payload)
	if i := d.Cause; i != nil {
		l += i.MarshalLen()
	}
	if i := d.DataNotificationDelay; i != nil {
		l += i.MarshalLen()
	}
	if i := d.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := d.DLLowPriorityTrafficThrottling; i != nil {
		l += i.MarshalLen()
	}
	if i := d.IMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := d.DLBufferingDuration; i != nil {
		l += i.MarshalLen()
	}
	if i := d.DLBufferingSuggestedPacketCount; i != nil {
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
func (d *DownlinkDataNotificationAcknowledge) SetLength() {
	if d != nil && d.Header != nil {
		d.Header.setLength(d.MarshalLen() - d.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (d *DownlinkDataNotificationAcknowledge) MessageTypeName() string {
	return "Downlink Data Notification Acknowledge"
}

// TEID returns the TEID in uint32.
func (d *DownlinkDataNotificationAcknowledge) TEID() uint32 {
	if d == nil {
		return 0
	}
	return d.Header.teid()
}

func (d *DownlinkDataNotificationAcknowledge) GetHeader() *Header {
	if d == nil {
		return nil
	}
	return d.Header
}

// TS 29.274 Table 7.2.11.2-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (d *DownlinkDataNotificationAcknowledge) addIE(i *ie.IE) {
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
	case ie.DelayValue:
		switch i.Instance() {
		case 0:
			if d.DataNotificationDelay == nil {
				d.DataNotificationDelay = i
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
	case ie.Throttling:
		switch i.Instance() {
		case 0:
			if d.DLLowPriorityTrafficThrottling == nil {
				d.DLLowPriorityTrafficThrottling = i
			}
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
	case ie.EPCTimer:
		switch i.Instance() {
		case 0:
			if d.DLBufferingDuration == nil {
				d.DLBufferingDuration = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.IntegerNumber:
		switch i.Instance() {
		case 0:
			if d.DLBufferingSuggestedPacketCount == nil {
				d.DLBufferingSuggestedPacketCount = i
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
func (d *DownlinkDataNotificationAcknowledge) Version() int { return d.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotificationAcknowledge) MessageType() uint8 { return d.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotificationAcknowledge) Sequence() uint32 { return d.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotificationAcknowledge) IsPiggybacking() bool {
	return d.GetHeader().IsPiggybacking()
}

// HasTEID reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotificationAcknowledge) HasTEID() bool { return d.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotificationAcknowledge) HasMessagePriority() bool {
	return d.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (d *DownlinkDataNotificationAcknowledge) MessagePriority() uint8 {
	return d.GetHeader().MessagePriority()
}

// SetSequenceNumber updates the common header safely on a nil receiver.
func (d *DownlinkDataNotificationAcknowledge) SetSequenceNumber(seq uint32) error {
	return d.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (d *DownlinkDataNotificationAcknowledge) SetMessagePriority(priority uint8) error {
	return d.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (d *DownlinkDataNotificationAcknowledge) SetTEID(teid uint32) { d.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (d *DownlinkDataNotificationAcknowledge) SetPiggybacking(value uint8) {
	d.GetHeader().SetPiggybacking(value)
}
