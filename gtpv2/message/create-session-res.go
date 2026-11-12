// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// CreateSessionResponse is a CreateSessionResponse Header and its IEs above.
type CreateSessionResponse struct {
	*Header
	Cause                         *ie.IE
	ChangeReportingAction         *ie.IE
	CSGInformationReportingAction *ie.IE
	HeNBInformationReporting      *ie.IE
	SenderFTEIDC                  *ie.IE
	PGWS5S8FTEIDC                 *ie.IE
	PAA                           *ie.IE
	APNRestriction                *ie.IE
	AMBR                          *ie.IE
	EBI                           *ie.IE
	PCO                           *ie.IE
	BearerContextsCreated         []*ie.IE
	BearerContextMarkedForRemoval []*ie.IE
	Recovery                      *ie.IE
	ChargingGatewayName           *ie.IE
	ChargingGatewayAddress        *ie.IE
	PGWFQCSID                     *ie.IE
	SGWFQCSID                     *ie.IE
	SGWLDN                        *ie.IE
	PGWLDN                        *ie.IE
	PGWBackOffTime                *ie.IE
	APCO                          *ie.IE
	TrustedTWANIPv4Parameters     *ie.IE
	IndicationFlags               *ie.IE
	PresenceReportingAreaAction   []*ie.IE
	PGWNodeLoadControlInformation *ie.IE
	PGWAPNLoadControlInformation  []*ie.IE
	SGWNodeLoadControlInformation *ie.IE
	PGWOverloadControlInformation []*ie.IE
	SGWOverloadControlInformation *ie.IE
	NBIFOMContainer               *ie.IE
	PDNConnectionChargingID       *ie.IE
	EPCO                          *ie.IE
	PrivateExtension              *ie.IE
	AdditionalIEs                 []*ie.IE
}

// NewCreateSessionResponse creates a new CreateSessionResponse.
func NewCreateSessionResponse(teid, seq uint32, ies ...*ie.IE) *CreateSessionResponse {
	c := &CreateSessionResponse{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeCreateSessionResponse, teid, seq, nil,
		),
	}
	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes CreateSessionResponse into bytes.
func (c *CreateSessionResponse) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes CreateSessionResponse into bytes.
func (c *CreateSessionResponse) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeCreateSessionResponse, c.MarshalLen()); err != nil {
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
	if i := c.SenderFTEIDC; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SenderFTEIDC", iePair{ie.FullyQualifiedTEID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PGWS5S8FTEIDC; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PGWS5S8FTEIDC", iePair{ie.FullyQualifiedTEID, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PAA; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PAA", iePair{ie.PDNAddressAllocation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.APNRestriction; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "APNRestriction", iePair{ie.APNRestriction, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.AMBR; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "AMBR", iePair{ie.AggregateMaximumBitRate, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.EBI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "EBI", iePair{ie.EPSBearerID, 0}); err != nil {
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
	for _, i := range c.BearerContextsCreated {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "BearerContextsCreated", iePair{ie.BearerContext, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.BearerContextMarkedForRemoval {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "BearerContextMarkedForRemoval", iePair{ie.BearerContext, 1}); err != nil {
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
	if i := c.ChargingGatewayName; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ChargingGatewayName", iePair{ie.FullyQualifiedDomainName, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.ChargingGatewayAddress; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ChargingGatewayAddress", iePair{ie.IPAddress, 0}); err != nil {
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
	if i := c.SGWLDN; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGWLDN", iePair{ie.LocalDistinguishedName, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PGWLDN; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PGWLDN", iePair{ie.LocalDistinguishedName, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PGWBackOffTime; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PGWBackOffTime", iePair{ie.EPCTimer, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.APCO; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "APCO", iePair{ie.AdditionalProtocolConfigurationOptions, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.TrustedTWANIPv4Parameters; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TrustedTWANIPv4Parameters", iePair{ie.IPv4ConfigurationParameters, 0}); err != nil {
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
	for _, i := range c.PresenceReportingAreaAction {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "PresenceReportingAreaAction", iePair{ie.PresenceReportingAreaAction, 0}); err != nil {
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
	if i := c.PDNConnectionChargingID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PDNConnectionChargingID", iePair{ie.ChargingID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.EPCO; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "EPCO", iePair{ie.ExtendedProtocolConfigurationOptions, 0}); err != nil {
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

// ParseCreateSessionResponse decodes given bytes as CreateSessionResponse.
func ParseCreateSessionResponse(b []byte) (*CreateSessionResponse, error) {
	c := &CreateSessionResponse{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as CreateSessionResponse.
func (c *CreateSessionResponse) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = CreateSessionResponse{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeCreateSessionResponse {
		return &UnexpectedTypeError{Expected: MsgTypeCreateSessionResponse, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *CreateSessionResponse) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	l := c.Header.MarshalLen() - len(c.Header.Payload)
	if i := c.Cause; i != nil {
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
	if i := c.SenderFTEIDC; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PGWS5S8FTEIDC; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PAA; i != nil {
		l += i.MarshalLen()
	}
	if i := c.APNRestriction; i != nil {
		l += i.MarshalLen()
	}
	if i := c.AMBR; i != nil {
		l += i.MarshalLen()
	}
	if i := c.EBI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PCO; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.BearerContextsCreated {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range c.BearerContextMarkedForRemoval {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ChargingGatewayName; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ChargingGatewayAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGWLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PGWLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PGWBackOffTime; i != nil {
		l += i.MarshalLen()
	}
	if i := c.APCO; i != nil {
		l += i.MarshalLen()
	}
	if i := c.TrustedTWANIPv4Parameters; i != nil {
		l += i.MarshalLen()
	}
	if i := c.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.PresenceReportingAreaAction {
		if i != nil {
			l += i.MarshalLen()
		}
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
	if i := c.PDNConnectionChargingID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.EPCO; i != nil {
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
func (c *CreateSessionResponse) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *CreateSessionResponse) MessageTypeName() string {
	return "Create Session Response"
}

// TEID returns the TEID in uint32.
func (c *CreateSessionResponse) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *CreateSessionResponse) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.2.2-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *CreateSessionResponse) addIE(i *ie.IE) {
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
	case ie.FullyQualifiedTEID:
		switch i.Instance() {
		case 0:
			if c.SenderFTEIDC == nil {
				c.SenderFTEIDC = i
			}
		case 1:
			if c.PGWS5S8FTEIDC == nil {
				c.PGWS5S8FTEIDC = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.PDNAddressAllocation:
		switch i.Instance() {
		case 0:
			if c.PAA == nil {
				c.PAA = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.APNRestriction:
		switch i.Instance() {
		case 0:
			if c.APNRestriction == nil {
				c.APNRestriction = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.AggregateMaximumBitRate:
		switch i.Instance() {
		case 0:
			if c.AMBR == nil {
				c.AMBR = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.EPSBearerID:
		switch i.Instance() {
		case 0:
			if c.EBI == nil {
				c.EBI = i
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
			c.BearerContextsCreated = append(c.BearerContextsCreated, i)
		case 1:
			c.BearerContextMarkedForRemoval = append(c.BearerContextMarkedForRemoval, i)
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
	case ie.FullyQualifiedDomainName:
		switch i.Instance() {
		case 0:
			if c.ChargingGatewayName == nil {
				c.ChargingGatewayName = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.IPAddress:
		switch i.Instance() {
		case 0:
			if c.ChargingGatewayAddress == nil {
				c.ChargingGatewayAddress = i
			}
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
	case ie.LocalDistinguishedName:
		switch i.Instance() {
		case 1:
			if c.PGWLDN == nil {
				c.PGWLDN = i
			}
		case 0:
			if c.SGWLDN == nil {
				c.SGWLDN = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.EPCTimer:
		switch i.Instance() {
		case 0:
			if c.PGWBackOffTime == nil {
				c.PGWBackOffTime = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.AdditionalProtocolConfigurationOptions:
		switch i.Instance() {
		case 0:
			if c.APCO == nil {
				c.APCO = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.IPv4ConfigurationParameters:
		switch i.Instance() {
		case 0:
			if c.TrustedTWANIPv4Parameters == nil {
				c.TrustedTWANIPv4Parameters = i
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
	case ie.PresenceReportingAreaAction:
		switch i.Instance() {
		case 0:
			c.PresenceReportingAreaAction = append(c.PresenceReportingAreaAction, i)
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
	case ie.ChargingID:
		switch i.Instance() {
		case 0:
			if c.PDNConnectionChargingID == nil {
				c.PDNConnectionChargingID = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.ExtendedProtocolConfigurationOptions:
		switch i.Instance() {
		case 0:
			if c.EPCO == nil {
				c.EPCO = i
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
func (c *CreateSessionResponse) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *CreateSessionResponse) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *CreateSessionResponse) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *CreateSessionResponse) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *CreateSessionResponse) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *CreateSessionResponse) HasMessagePriority() bool { return c.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *CreateSessionResponse) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *CreateSessionResponse) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *CreateSessionResponse) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *CreateSessionResponse) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *CreateSessionResponse) SetPiggybacking(value uint8) { c.GetHeader().SetPiggybacking(value) }
