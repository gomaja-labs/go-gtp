// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// CreateBearerRequest is a CreateBearerRequest Header and its IEs above.
type CreateBearerRequest struct {
	*Header
	PTI                           *ie.IE
	LinkedEBI                     *ie.IE
	PCO                           *ie.IE
	BearerContexts                []*ie.IE
	PGWFQCSID                     *ie.IE
	SGWFQCSID                     *ie.IE
	ChangeReportingAction         *ie.IE
	CSGInformationReportingAction *ie.IE
	HeNBInformationReporting      *ie.IE
	PresenceReportingAreaAction   []*ie.IE
	IndicationFlags               *ie.IE
	PGWNodeLoadControlInformation *ie.IE
	PGWAPNLoadControlInformation  []*ie.IE
	SGWNodeLoadControlInformation *ie.IE
	PGWOverloadControlInformation []*ie.IE
	SGWOverloadControlInformation *ie.IE
	NBIFOMContainer               *ie.IE
	PrivateExtension              *ie.IE
	AdditionalIEs                 []*ie.IE
}

// NewCreateBearerRequest creates a new CreateBearerRequest.
func NewCreateBearerRequest(teid, seq uint32, ies ...*ie.IE) *CreateBearerRequest {
	c := &CreateBearerRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeCreateBearerRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes CreateBearerRequest into bytes.
func (c *CreateBearerRequest) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes CreateBearerRequest into bytes.
func (c *CreateBearerRequest) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeCreateBearerRequest, c.MarshalLen()); err != nil {
		return err
	}
	if len(b) < c.MarshalLen() {
		return ErrTooShortToMarshal
	}
	c.Header.Payload = nil
	c.Header.Payload = make([]byte, c.MarshalLen()-c.Header.MarshalLen())
	offset := 0
	if i := c.PTI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PTI", iePair{ie.ProcedureTransactionID, 0}); err != nil {
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
	if i := c.PCO; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PCO", iePair{ie.ProtocolConfigurationOptions, 0}); err != nil {
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
	if i := c.PGWFQCSID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PGWFQCSID", iePair{ie.FullyQualifiedCSID, 0}); err != nil {
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
	if i := c.HeNBInformationReporting; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "HeNBInformationReporting", iePair{ie.HeNBInformationReporting, 0}); err != nil {
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
	if i := c.IndicationFlags; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "IndicationFlags", iePair{ie.Indication, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PGWNodeLoadControlInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PGWNodeLoadControlInformation", iePair{ie.LoadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.PGWAPNLoadControlInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "PGWAPNLoadControlInformation", iePair{ie.LoadControlInformation, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGWNodeLoadControlInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGWNodeLoadControlInformation", iePair{ie.LoadControlInformation, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.PGWOverloadControlInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "PGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
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
	if i := c.NBIFOMContainer; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "NBIFOMContainer", iePair{ie.FContainer, 0}); err != nil {
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

// ParseCreateBearerRequest decodes given bytes as CreateBearerRequest.
func ParseCreateBearerRequest(b []byte) (*CreateBearerRequest, error) {
	c := &CreateBearerRequest{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as CreateBearerRequest.
func (c *CreateBearerRequest) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = CreateBearerRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeCreateBearerRequest {
		return &UnexpectedTypeError{Expected: MsgTypeCreateBearerRequest, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *CreateBearerRequest) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	l := c.Header.MarshalLen() - len(c.Header.Payload)
	if i := c.PTI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.LinkedEBI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PCO; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.BearerContexts {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.PGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ChangeReportingAction; i != nil {
		l += i.MarshalLen()
	}
	if i := c.CSGInformationReportingAction; i != nil {
		l += i.MarshalLen()
	}
	if i := c.HeNBInformationReporting; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.PresenceReportingAreaAction {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PGWNodeLoadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.PGWAPNLoadControlInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.SGWNodeLoadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.PGWOverloadControlInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := c.NBIFOMContainer; i != nil {
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
func (c *CreateBearerRequest) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *CreateBearerRequest) MessageTypeName() string {
	return "Create Bearer Request"
}

// TEID returns the TEID in uint32.
func (c *CreateBearerRequest) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *CreateBearerRequest) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.2.3-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *CreateBearerRequest) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		c.Header.err = err
	}
	switch i.Type {
	case ie.ProcedureTransactionID:
		switch i.Instance() {
		case 0:
			if c.PTI == nil {
				c.PTI = i
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
	case ie.ProtocolConfigurationOptions:
		switch i.Instance() {
		case 0:
			if c.PCO == nil {
				c.PCO = i
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
	case ie.FullyQualifiedCSID:
		switch i.Instance() {
		case 0:
			if c.PGWFQCSID == nil {
				c.PGWFQCSID = i
			}
		case 1:
			if c.SGWFQCSID == nil {
				c.SGWFQCSID = i
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
	case ie.HeNBInformationReporting:
		switch i.Instance() {
		case 0:
			if c.HeNBInformationReporting == nil {
				c.HeNBInformationReporting = i
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
	case ie.Indication:
		switch i.Instance() {
		case 0:
			if c.IndicationFlags == nil {
				c.IndicationFlags = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.LoadControlInformation:
		switch i.Instance() {
		case 0:
			if c.PGWNodeLoadControlInformation == nil {
				c.PGWNodeLoadControlInformation = i
			}
		case 1:
			c.PGWAPNLoadControlInformation = append(c.PGWAPNLoadControlInformation, i)
		case 2:
			if c.SGWNodeLoadControlInformation == nil {
				c.SGWNodeLoadControlInformation = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			c.PGWOverloadControlInformation = append(c.PGWOverloadControlInformation, i)
		case 1:
			if c.SGWOverloadControlInformation == nil {
				c.SGWOverloadControlInformation = i
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
func (c *CreateBearerRequest) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *CreateBearerRequest) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *CreateBearerRequest) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *CreateBearerRequest) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *CreateBearerRequest) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *CreateBearerRequest) HasMessagePriority() bool { return c.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *CreateBearerRequest) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *CreateBearerRequest) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *CreateBearerRequest) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *CreateBearerRequest) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *CreateBearerRequest) SetPiggybacking(value uint8) { c.GetHeader().SetPiggybacking(value) }
