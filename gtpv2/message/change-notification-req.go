// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

// ChangeNotificationRequest is a ChangeNotificationRequest Header and its IEs above.
type ChangeNotificationRequest struct {
	*Header
	IMSI                             *ie.IE
	MEI                              *ie.IE
	IndicationFlags                  *ie.IE
	RATType                          *ie.IE
	ULI                              *ie.IE
	UCI                              *ie.IE
	PGWS5S8IPAddressForControlPlane  *ie.IE
	LinkedEBI                        *ie.IE
	PresenceReportingAreaInformation []*ie.IE
	MOExceptionDataCounter           *ie.IE
	SecondaryRATUsageDataReport      []*ie.IE
	PrivateExtension                 *ie.IE
	AdditionalIEs                    []*ie.IE
}

// NewChangeNotificationRequest creates a new ChangeNotificationRequest.
func NewChangeNotificationRequest(teid, seq uint32, ies ...*ie.IE) *ChangeNotificationRequest {
	c := &ChangeNotificationRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeChangeNotificationRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes ChangeNotificationRequest into bytes.
func (c *ChangeNotificationRequest) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ChangeNotificationRequest into bytes.
func (c *ChangeNotificationRequest) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeChangeNotificationRequest, c.MarshalLen()); err != nil {
		return err
	}
	if len(b) < c.MarshalLen() {
		return ErrTooShortToMarshal
	}
	c.Header.Payload = nil
	c.Header.Payload = make([]byte, c.MarshalLen()-c.Header.MarshalLen())
	offset := 0
	if i := c.IMSI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "IMSI", iePair{ie.IMSI, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MEI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MEI", iePair{ie.MobileEquipmentIdentity, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.IndicationFlags; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "IndicationFlags", iePair{ie.Indication, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.RATType; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "RATType", iePair{ie.RATType, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.ULI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ULI", iePair{ie.UserLocationInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.UCI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UCI", iePair{ie.UserCSGInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PGWS5S8IPAddressForControlPlane; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PGWS5S8IPAddressForControlPlane", iePair{ie.IPAddress, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.LinkedEBI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "LinkedEBI", iePair{ie.EPSBearerID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.PresenceReportingAreaInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "PresenceReportingAreaInformation", iePair{ie.PresenceReportingAreaInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MOExceptionDataCounter; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MOExceptionDataCounter", iePair{ie.Counter, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.SecondaryRATUsageDataReport {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "SecondaryRATUsageDataReport", iePair{ie.SecondaryRATUsageDataReport, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PrivateExtension; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PrivateExtension", iePair{ie.PrivateExtension, 255}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.AdditionalIEs {
		if i == nil {
			continue
		}
		if err := i.MarshalTo(c.Payload[offset:]); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	c.Header.SetLength()
	return c.Header.MarshalTo(b)
}

// ParseChangeNotificationRequest decodes given bytes as ChangeNotificationRequest.
func ParseChangeNotificationRequest(b []byte) (*ChangeNotificationRequest, error) {
	c := &ChangeNotificationRequest{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as ChangeNotificationRequest.
func (c *ChangeNotificationRequest) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = ChangeNotificationRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeChangeNotificationRequest {
		return &UnexpectedTypeError{Expected: MsgTypeChangeNotificationRequest, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *ChangeNotificationRequest) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	l := c.Header.MarshalLen() - len(c.Header.Payload)
	if i := c.IMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MEI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := c.RATType; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ULI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UCI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PGWS5S8IPAddressForControlPlane; i != nil {
		l += i.MarshalLen()
	}
	if i := c.LinkedEBI; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.PresenceReportingAreaInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.MOExceptionDataCounter; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.SecondaryRATUsageDataReport {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.PrivateExtension; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.AdditionalIEs {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	return l
}

// SetLength sets the length in Length field.
func (c *ChangeNotificationRequest) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *ChangeNotificationRequest) MessageTypeName() string {
	return "Change Notification Request"
}

// TEID returns the TEID in uint32.
func (c *ChangeNotificationRequest) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *ChangeNotificationRequest) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.3.14-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *ChangeNotificationRequest) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		c.Header.err = err
	}
	switch i.Type {
	case ie.IMSI:
		switch i.Instance() {
		case 0:
			if c.IMSI == nil {
				c.IMSI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.MobileEquipmentIdentity:
		switch i.Instance() {
		case 0:
			if c.MEI == nil {
				c.MEI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.Indication:
		switch i.Instance() {
		case 0:
			if c.IndicationFlags == nil {
				c.IndicationFlags = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.RATType:
		switch i.Instance() {
		case 0:
			if c.RATType == nil {
				c.RATType = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.UserLocationInformation:
		switch i.Instance() {
		case 0:
			if c.ULI == nil {
				c.ULI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.UserCSGInformation:
		switch i.Instance() {
		case 0:
			if c.UCI == nil {
				c.UCI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.IPAddress:
		switch i.Instance() {
		case 0:
			if c.PGWS5S8IPAddressForControlPlane == nil {
				c.PGWS5S8IPAddressForControlPlane = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.EPSBearerID:
		switch i.Instance() {
		case 0:
			if c.LinkedEBI == nil {
				c.LinkedEBI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.PresenceReportingAreaInformation:
		switch i.Instance() {
		case 0:
			c.PresenceReportingAreaInformation = append(c.PresenceReportingAreaInformation, i)
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.Counter:
		switch i.Instance() {
		case 0:
			if c.MOExceptionDataCounter == nil {
				c.MOExceptionDataCounter = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.SecondaryRATUsageDataReport:
		switch i.Instance() {
		case 0:
			c.SecondaryRATUsageDataReport = append(c.SecondaryRATUsageDataReport, i)
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.PrivateExtension:
		if c.PrivateExtension == nil {
			c.PrivateExtension = i
		} else if c.PrivateExtension.Instance() != i.Instance() {
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	default:
		c.AdditionalIEs = append(c.AdditionalIEs, i)
	}
}

// Version reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationRequest) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationRequest) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationRequest) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationRequest) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationRequest) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationRequest) HasMessagePriority() bool {
	return c.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationRequest) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *ChangeNotificationRequest) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *ChangeNotificationRequest) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *ChangeNotificationRequest) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *ChangeNotificationRequest) SetPiggybacking(value uint8) {
	c.GetHeader().SetPiggybacking(value)
}
