// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

// ChangeNotificationResponse is a ChangeNotificationResponse Header and its IEs above.
type ChangeNotificationResponse struct {
	*Header
	IMSI                          *ie.IE
	MEI                           *ie.IE
	Cause                         *ie.IE
	ChangeReportingAction         *ie.IE
	CSGInformationReportingAction *ie.IE
	PresenceReportingAreaAction   []*ie.IE
	PrivateExtension              *ie.IE
	AdditionalIEs                 []*ie.IE
}

// NewChangeNotificationResponse creates a new ChangeNotificationResponse.
func NewChangeNotificationResponse(teid, seq uint32, ies ...*ie.IE) *ChangeNotificationResponse {
	c := &ChangeNotificationResponse{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeChangeNotificationResponse, teid, seq, nil,
		),
	}

	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes ChangeNotificationResponse into bytes.
func (c *ChangeNotificationResponse) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ChangeNotificationResponse into bytes.
func (c *ChangeNotificationResponse) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeChangeNotificationResponse, c.MarshalLen()); err != nil {
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
	if i := c.Cause; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "Cause", iePair{ie.Cause, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.ChangeReportingAction; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ChangeReportingAction", iePair{ie.ChangeReportingAction, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.CSGInformationReportingAction; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "CSGInformationReportingAction", iePair{ie.CSGInformationReportingAction, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.PresenceReportingAreaAction {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "PresenceReportingAreaAction", iePair{ie.PresenceReportingAreaAction, 0}); err != nil {
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

// ParseChangeNotificationResponse decodes given bytes as ChangeNotificationResponse.
func ParseChangeNotificationResponse(b []byte) (*ChangeNotificationResponse, error) {
	c := &ChangeNotificationResponse{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as ChangeNotificationResponse.
func (c *ChangeNotificationResponse) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = ChangeNotificationResponse{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeChangeNotificationResponse {
		return &UnexpectedTypeError{Expected: MsgTypeChangeNotificationResponse, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *ChangeNotificationResponse) MarshalLen() int {
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
	if i := c.Cause; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ChangeReportingAction; i != nil {
		l += i.MarshalLen()
	}
	if i := c.CSGInformationReportingAction; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.PresenceReportingAreaAction {
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
func (c *ChangeNotificationResponse) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *ChangeNotificationResponse) MessageTypeName() string {
	return "Change Notification Response"
}

// TEID returns the TEID in uint32.
func (c *ChangeNotificationResponse) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *ChangeNotificationResponse) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.3.15-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *ChangeNotificationResponse) addIE(i *ie.IE) {
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
	case ie.Cause:
		switch i.Instance() {
		case 0:
			if c.Cause == nil {
				c.Cause = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.ChangeReportingAction:
		switch i.Instance() {
		case 0:
			if c.ChangeReportingAction == nil {
				c.ChangeReportingAction = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.CSGInformationReportingAction:
		switch i.Instance() {
		case 0:
			if c.CSGInformationReportingAction == nil {
				c.CSGInformationReportingAction = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.PresenceReportingAreaAction:
		switch i.Instance() {
		case 0:
			c.PresenceReportingAreaAction = append(c.PresenceReportingAreaAction, i)
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
func (c *ChangeNotificationResponse) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationResponse) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationResponse) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationResponse) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationResponse) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationResponse) HasMessagePriority() bool {
	return c.GetHeader().HasMessagePriority()
}

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *ChangeNotificationResponse) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *ChangeNotificationResponse) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *ChangeNotificationResponse) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *ChangeNotificationResponse) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *ChangeNotificationResponse) SetPiggybacking(value uint8) {
	c.GetHeader().SetPiggybacking(value)
}
