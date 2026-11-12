// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// CreateBearerResponse is a CreateBearerResponse Header and its IEs above.
type CreateBearerResponse struct {
	*Header
	Cause                               *ie.IE
	BearerContexts                      []*ie.IE
	Recovery                            *ie.IE
	MMEFQCSID                           *ie.IE
	SGWFQCSID                           *ie.IE
	EPDGFQCSID                          *ie.IE
	TWANFQCSID                          *ie.IE
	PCO                                 *ie.IE
	UETimeZone                          *ie.IE
	ULI                                 *ie.IE
	TWANIdentifier                      *ie.IE
	MMEOverloadControlInformation       *ie.IE
	SGWOverloadControlInformation       *ie.IE
	PresenceReportingAreaInformation    []*ie.IE
	MMESGSNIdentifierOrUELocalIPAddress *ie.IE
	TWANePDGOverloadControlInformation  *ie.IE
	WLANLocationInformation             *ie.IE
	WLANLocationTimestamp               *ie.IE
	UEUDPPort                           *ie.IE
	NBIFOMContainer                     *ie.IE
	UETCPPort                           *ie.IE
	PrivateExtension                    *ie.IE
	AdditionalIEs                       []*ie.IE
}

// NewCreateBearerResponse creates a new CreateBearerResponse.
func NewCreateBearerResponse(teid, seq uint32, ies ...*ie.IE) *CreateBearerResponse {
	c := &CreateBearerResponse{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeCreateBearerResponse, teid, seq, nil,
		),
	}

	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes CreateBearerResponse into bytes.
func (c *CreateBearerResponse) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes CreateBearerResponse into bytes.
func (c *CreateBearerResponse) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeCreateBearerResponse, c.MarshalLen()); err != nil {
		return err
	}
	if len(b) < c.MarshalLen() {
		return ErrTooShortToMarshal
	}
	c.Header.Payload = nil
	c.Header.Payload = make([]byte, c.MarshalLen()-c.Header.MarshalLen())
	offset := 0
	if i := c.Cause; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "Cause", iePair{ie.Cause, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.BearerContexts {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "BearerContexts", iePair{ie.BearerContext, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.Recovery; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "Recovery", iePair{ie.Recovery, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MMEFQCSID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMEFQCSID", iePair{ie.FullyQualifiedCSID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGWFQCSID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGWFQCSID", iePair{ie.FullyQualifiedCSID, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.EPDGFQCSID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "EPDGFQCSID", iePair{ie.FullyQualifiedCSID, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.TWANFQCSID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TWANFQCSID", iePair{ie.FullyQualifiedCSID, 3}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PCO; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PCO", iePair{ie.ProtocolConfigurationOptions, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.UETimeZone; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UETimeZone", iePair{ie.UETimeZone, 0}); err != nil {
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
	if i := c.TWANIdentifier; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TWANIdentifier", iePair{ie.TWANIdentifier, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MMEOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMEOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGWOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 1}); err != nil {
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
	if i := c.MMESGSNIdentifierOrUELocalIPAddress; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMESGSNIdentifierOrUELocalIPAddress", iePair{ie.IPAddress, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.TWANePDGOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TWANePDGOverloadControlInformation", iePair{ie.OverloadControlInformation, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.WLANLocationInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "WLANLocationInformation", iePair{ie.TWANIdentifier, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.WLANLocationTimestamp; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "WLANLocationTimestamp", iePair{ie.TWANIdentifierTimestamp, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.UEUDPPort; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UEUDPPort", iePair{ie.PortNumber, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.NBIFOMContainer; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "NBIFOMContainer", iePair{ie.FContainer, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.UETCPPort; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UETCPPort", iePair{ie.PortNumber, 1}); err != nil {
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

// ParseCreateBearerResponse decodes given bytes as CreateBearerResponse.
func ParseCreateBearerResponse(b []byte) (*CreateBearerResponse, error) {
	c := &CreateBearerResponse{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as CreateBearerResponse.
func (c *CreateBearerResponse) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = CreateBearerResponse{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeCreateBearerResponse {
		return &UnexpectedTypeError{Expected: MsgTypeCreateBearerResponse, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *CreateBearerResponse) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	l := c.Header.MarshalLen() - len(c.Header.Payload)
	if i := c.Cause; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.BearerContexts {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMEFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.EPDGFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.TWANFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PCO; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UETimeZone; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ULI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.TWANIdentifier; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMEOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.PresenceReportingAreaInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.MMESGSNIdentifierOrUELocalIPAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := c.TWANePDGOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := c.WLANLocationInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := c.WLANLocationTimestamp; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UEUDPPort; i != nil {
		l += i.MarshalLen()
	}
	if i := c.NBIFOMContainer; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UETCPPort; i != nil {
		l += i.MarshalLen()
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
func (c *CreateBearerResponse) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *CreateBearerResponse) MessageTypeName() string {
	return "Create Bearer Response"
}

// TEID returns the TEID in uint32.
func (c *CreateBearerResponse) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *CreateBearerResponse) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.2.4-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *CreateBearerResponse) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		c.Header.err = err
	}
	switch i.Type {
	case ie.Cause:
		switch i.Instance() {
		case 0:
			if c.Cause == nil {
				c.Cause = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.BearerContext:
		switch i.Instance() {
		case 0:
			c.BearerContexts = append(c.BearerContexts, i)
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.Recovery:
		switch i.Instance() {
		case 0:
			if c.Recovery == nil {
				c.Recovery = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.FullyQualifiedCSID:
		switch i.Instance() {
		case 0:
			if c.MMEFQCSID == nil {
				c.MMEFQCSID = i
			}
		case 1:
			if c.SGWFQCSID == nil {
				c.SGWFQCSID = i
			}
		case 2:
			if c.EPDGFQCSID == nil {
				c.EPDGFQCSID = i
			}
		case 3:
			if c.TWANFQCSID == nil {
				c.TWANFQCSID = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.ProtocolConfigurationOptions:
		switch i.Instance() {
		case 0:
			if c.PCO == nil {
				c.PCO = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.UETimeZone:
		switch i.Instance() {
		case 0:
			if c.UETimeZone == nil {
				c.UETimeZone = i
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
	case ie.TWANIdentifier:
		switch i.Instance() {
		case 0:
			if c.TWANIdentifier == nil {
				c.TWANIdentifier = i
			}
		case 1:
			if c.WLANLocationInformation == nil {
				c.WLANLocationInformation = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.TWANIdentifierTimestamp:
		switch i.Instance() {
		case 1:
			if c.WLANLocationTimestamp == nil {
				c.WLANLocationTimestamp = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			if c.MMEOverloadControlInformation == nil {
				c.MMEOverloadControlInformation = i
			}
		case 1:
			if c.SGWOverloadControlInformation == nil {
				c.SGWOverloadControlInformation = i
			}
		case 2:
			if c.TWANePDGOverloadControlInformation == nil {
				c.TWANePDGOverloadControlInformation = i
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
	case ie.IPAddress:
		switch i.Instance() {
		case 0:
			if c.MMESGSNIdentifierOrUELocalIPAddress == nil {
				c.MMESGSNIdentifierOrUELocalIPAddress = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.PortNumber:
		switch i.Instance() {
		case 0:
			if c.UEUDPPort == nil {
				c.UEUDPPort = i
			}
		case 1:
			if c.UETCPPort == nil {
				c.UETCPPort = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.FContainer:
		switch i.Instance() {
		case 0:
			if c.NBIFOMContainer == nil {
				c.NBIFOMContainer = i
			}
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
func (c *CreateBearerResponse) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *CreateBearerResponse) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *CreateBearerResponse) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *CreateBearerResponse) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *CreateBearerResponse) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *CreateBearerResponse) HasMessagePriority() bool { return c.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *CreateBearerResponse) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *CreateBearerResponse) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *CreateBearerResponse) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *CreateBearerResponse) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *CreateBearerResponse) SetPiggybacking(value uint8) { c.GetHeader().SetPiggybacking(value) }
