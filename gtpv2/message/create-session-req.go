// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// CreateSessionRequest is a CreateSessionRequest Header and its IEs above.
type CreateSessionRequest struct {
	*Header
	IMSI                               *ie.IE
	MSISDN                             *ie.IE
	MEI                                *ie.IE
	ULI                                *ie.IE
	ServingNetwork                     *ie.IE
	RATType                            *ie.IE
	IndicationFlags                    *ie.IE
	SenderFTEIDC                       *ie.IE
	PGWS5S8FTEIDC                      *ie.IE
	APN                                *ie.IE
	SelectionMode                      *ie.IE
	PDNType                            *ie.IE
	PAA                                *ie.IE
	APNRestriction                     *ie.IE
	AMBR                               *ie.IE
	LinkedEBI                          *ie.IE
	TWMI                               *ie.IE
	PCO                                *ie.IE
	BearerContextsToBeCreated          []*ie.IE
	BearerContextsToBeRemoved          []*ie.IE
	TraceInformation                   *ie.IE
	Recovery                           *ie.IE
	MMEFQCSID                          *ie.IE
	SGWFQCSID                          *ie.IE
	EPDGFQCSID                         *ie.IE
	TWANFQCSID                         *ie.IE
	UETimeZone                         *ie.IE
	UCI                                *ie.IE
	ChargingCharacteristics            *ie.IE
	MMESGSNLDN                         *ie.IE
	SGWLDN                             *ie.IE
	EPDGLDN                            *ie.IE
	TWANLDN                            *ie.IE
	SignallingPriorityIndication       *ie.IE
	UELocalIPAddress                   *ie.IE
	UEUDPPort                          *ie.IE
	APCO                               *ie.IE
	HeNBLocalIPAddress                 *ie.IE
	HeNBUDPPort                        *ie.IE
	MMESGSNIdentifier                  *ie.IE
	TWANIdentifier                     *ie.IE
	EPDGIPAddress                      *ie.IE
	CNOperatorSelectionEntity          *ie.IE
	PresenceReportingAreaInformation   []*ie.IE
	MMESGSNOverloadControlInformation  *ie.IE
	SGWOverloadControlInformation      *ie.IE
	TWANePDGOverloadControlInformation *ie.IE
	OriginationTimeStamp               *ie.IE
	MaximumWaitTime                    *ie.IE
	WLANLocationInformation            *ie.IE
	WLANLocationTimeStamp              *ie.IE
	NBIFOMContainer                    *ie.IE
	RemoteUEContextConnected           []*ie.IE
	TGPPAAAServerIdentifier            *ie.IE
	EPCO                               *ie.IE
	ServingPLMNRateControl             *ie.IE
	MOExceptionDataCounter             *ie.IE
	UETCPPort                          *ie.IE
	MappedUEUsageType                  *ie.IE
	ULIForSGW                          *ie.IE
	SGWUNodeName                       *ie.IE
	SecondaryRATUsageDataReport        []*ie.IE
	UPFunctionSelectionIndicationFlags *ie.IE
	APNRateControlStatus               *ie.IE
	PrivateExtension                   *ie.IE
	AdditionalIEs                      []*ie.IE
}

// NewCreateSessionRequest creates a new CreateSessionRequest.
func NewCreateSessionRequest(teid, seq uint32, ies ...*ie.IE) *CreateSessionRequest {
	c := &CreateSessionRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeCreateSessionRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes CreateSessionRequest into bytes.
func (c *CreateSessionRequest) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes CreateSessionRequest into bytes.
func (c *CreateSessionRequest) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeCreateSessionRequest, c.MarshalLen()); err != nil {
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
	if i := c.MSISDN; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MSISDN", iePair{ie.MSISDN, 0}); err != nil {
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
	if i := c.ULI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ULI", iePair{ie.UserLocationInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.ServingNetwork; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ServingNetwork", iePair{ie.ServingNetwork, 0}); err != nil {
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
	if i := c.IndicationFlags; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "IndicationFlags", iePair{ie.Indication, 0}); err != nil {
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
	if i := c.APN; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "APN", iePair{ie.AccessPointName, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SelectionMode; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SelectionMode", iePair{ie.SelectionMode, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PDNType; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PDNType", iePair{ie.PDNType, 0}); err != nil {
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
	if i := c.LinkedEBI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "LinkedEBI", iePair{ie.EPSBearerID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.TWMI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TWMI", iePair{ie.TrustedWLANModeIndication, 0}); err != nil {
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
	for _, i := range c.BearerContextsToBeCreated {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "BearerContextsToBeCreated", iePair{ie.BearerContext, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.BearerContextsToBeRemoved {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "BearerContextsToBeRemoved", iePair{ie.BearerContext, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.TraceInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TraceInformation", iePair{ie.TraceInformation, 0}); err != nil {
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
	if i := c.UETimeZone; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UETimeZone", iePair{ie.UETimeZone, 0}); err != nil {
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
	if i := c.ChargingCharacteristics; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ChargingCharacteristics", iePair{ie.ChargingCharacteristics, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MMESGSNLDN; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMESGSNLDN", iePair{ie.LocalDistinguishedName, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGWLDN; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGWLDN", iePair{ie.LocalDistinguishedName, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.EPDGLDN; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "EPDGLDN", iePair{ie.LocalDistinguishedName, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.TWANLDN; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TWANLDN", iePair{ie.LocalDistinguishedName, 3}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SignallingPriorityIndication; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SignallingPriorityIndication", iePair{ie.SignallingPriorityIndication, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.UELocalIPAddress; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UELocalIPAddress", iePair{ie.IPAddress, 0}); err != nil {
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
	if i := c.APCO; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "APCO", iePair{ie.AdditionalProtocolConfigurationOptions, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.HeNBLocalIPAddress; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "HeNBLocalIPAddress", iePair{ie.IPAddress, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.HeNBUDPPort; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "HeNBUDPPort", iePair{ie.PortNumber, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MMESGSNIdentifier; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMESGSNIdentifier", iePair{ie.IPAddress, 2}); err != nil {
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
	if i := c.EPDGIPAddress; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "EPDGIPAddress", iePair{ie.IPAddress, 3}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.CNOperatorSelectionEntity; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "CNOperatorSelectionEntity", iePair{ie.CNOperatorSelectionEntity, 0}); err != nil {
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
	if i := c.MMESGSNOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMESGSNOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
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
	if i := c.TWANePDGOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TWANePDGOverloadControlInformation", iePair{ie.OverloadControlInformation, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.OriginationTimeStamp; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "OriginationTimeStamp", iePair{ie.MillisecondTimeStamp, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MaximumWaitTime; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MaximumWaitTime", iePair{ie.IntegerNumber, 0}); err != nil {
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
	if i := c.WLANLocationTimeStamp; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "WLANLocationTimeStamp", iePair{ie.TWANIdentifierTimestamp, 0}); err != nil {
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
	for _, i := range c.RemoteUEContextConnected {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "RemoteUEContextConnected", iePair{ie.RemoteUEContext, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.TGPPAAAServerIdentifier; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TGPPAAAServerIdentifier", iePair{ie.NodeIdentifier, 0}); err != nil {
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
	if i := c.ServingPLMNRateControl; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ServingPLMNRateControl", iePair{ie.ServingPLMNRateControl, 0}); err != nil {
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
	if i := c.UETCPPort; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UETCPPort", iePair{ie.PortNumber, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MappedUEUsageType; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MappedUEUsageType", iePair{ie.MappedUEUsageType, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.ULIForSGW; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ULIForSGW", iePair{ie.UserLocationInformation, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGWUNodeName; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGWUNodeName", iePair{ie.FullyQualifiedDomainName, 0}); err != nil {
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
	if i := c.UPFunctionSelectionIndicationFlags; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UPFunctionSelectionIndicationFlags", iePair{ie.UPFunctionSelectionIndicationFlags, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.APNRateControlStatus; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "APNRateControlStatus", iePair{ie.APNRateControlStatus, 0}); err != nil {
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

// ParseCreateSessionRequest decodes given bytes as CreateSessionRequest.
func ParseCreateSessionRequest(b []byte) (*CreateSessionRequest, error) {
	c := &CreateSessionRequest{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as CreateSessionRequest.
func (c *CreateSessionRequest) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = CreateSessionRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeCreateSessionRequest {
		return &UnexpectedTypeError{Expected: MsgTypeCreateSessionRequest, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *CreateSessionRequest) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	l := c.Header.MarshalLen() - len(c.Header.Payload)
	if i := c.IMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MSISDN; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MEI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ULI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ServingNetwork; i != nil {
		l += i.MarshalLen()
	}
	if i := c.RATType; i != nil {
		l += i.MarshalLen()
	}
	if i := c.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SenderFTEIDC; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PGWS5S8FTEIDC; i != nil {
		l += i.MarshalLen()
	}
	if i := c.APN; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SelectionMode; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PDNType; i != nil {
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
	if i := c.LinkedEBI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.TWMI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PCO; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.BearerContextsToBeCreated {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range c.BearerContextsToBeRemoved {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.TraceInformation; i != nil {
		l += i.MarshalLen()
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
	if i := c.UETimeZone; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UCI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ChargingCharacteristics; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMESGSNLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGWLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := c.EPDGLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := c.TWANLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SignallingPriorityIndication; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UELocalIPAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UEUDPPort; i != nil {
		l += i.MarshalLen()
	}
	if i := c.APCO; i != nil {
		l += i.MarshalLen()
	}
	if i := c.HeNBLocalIPAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := c.HeNBUDPPort; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMESGSNIdentifier; i != nil {
		l += i.MarshalLen()
	}
	if i := c.TWANIdentifier; i != nil {
		l += i.MarshalLen()
	}
	if i := c.EPDGIPAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := c.CNOperatorSelectionEntity; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.PresenceReportingAreaInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.MMESGSNOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := c.TWANePDGOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := c.OriginationTimeStamp; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MaximumWaitTime; i != nil {
		l += i.MarshalLen()
	}
	if i := c.WLANLocationInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := c.WLANLocationTimeStamp; i != nil {
		l += i.MarshalLen()
	}
	if i := c.NBIFOMContainer; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.RemoteUEContextConnected {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.TGPPAAAServerIdentifier; i != nil {
		l += i.MarshalLen()
	}
	if i := c.EPCO; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ServingPLMNRateControl; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MOExceptionDataCounter; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UETCPPort; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MappedUEUsageType; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ULIForSGW; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGWUNodeName; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.SecondaryRATUsageDataReport {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.UPFunctionSelectionIndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := c.APNRateControlStatus; i != nil {
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
func (c *CreateSessionRequest) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *CreateSessionRequest) MessageTypeName() string {
	return "Create Session Request"
}

// TEID returns the TEID in uint32.
func (c *CreateSessionRequest) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *CreateSessionRequest) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.2.1-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *CreateSessionRequest) addIE(i *ie.IE) {
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
	case ie.MSISDN:
		switch i.Instance() {
		case 0:
			if c.MSISDN == nil {
				c.MSISDN = i
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
	case ie.UserLocationInformation:
		switch i.Instance() {
		case 0:
			if c.ULI == nil {
				c.ULI = i
			}
		case 1:
			if c.ULIForSGW == nil {
				c.ULIForSGW = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.ServingNetwork:
		switch i.Instance() {
		case 0:
			if c.ServingNetwork == nil {
				c.ServingNetwork = i
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
	case ie.Indication:
		switch i.Instance() {
		case 0:
			if c.IndicationFlags == nil {
				c.IndicationFlags = i
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
	case ie.AccessPointName:
		switch i.Instance() {
		case 0:
			if c.APN == nil {
				c.APN = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.SelectionMode:
		switch i.Instance() {
		case 0:
			if c.SelectionMode == nil {
				c.SelectionMode = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.PDNType:
		switch i.Instance() {
		case 0:
			if c.PDNType == nil {
				c.PDNType = i
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
			if c.LinkedEBI == nil {
				c.LinkedEBI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.TrustedWLANModeIndication:
		switch i.Instance() {
		case 0:
			if c.TWMI == nil {
				c.TWMI = i
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
			c.BearerContextsToBeCreated = append(c.BearerContextsToBeCreated, i)
		case 1:
			c.BearerContextsToBeRemoved = append(c.BearerContextsToBeRemoved, i)
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
	case ie.TraceInformation:
		switch i.Instance() {
		case 0:
			if c.TraceInformation == nil {
				c.TraceInformation = i
			}
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
	case ie.UETimeZone:
		switch i.Instance() {
		case 0:
			if c.UETimeZone == nil {
				c.UETimeZone = i
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
	case ie.ChargingCharacteristics:
		switch i.Instance() {
		case 0:
			if c.ChargingCharacteristics == nil {
				c.ChargingCharacteristics = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.LocalDistinguishedName:
		switch i.Instance() {
		case 0:
			if c.MMESGSNLDN == nil {
				c.MMESGSNLDN = i
			}
		case 1:
			if c.SGWLDN == nil {
				c.SGWLDN = i
			}
		case 2:
			if c.EPDGLDN == nil {
				c.EPDGLDN = i
			}
		case 3:
			if c.TWANLDN == nil {
				c.TWANLDN = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.SignallingPriorityIndication:
		switch i.Instance() {
		case 0:
			if c.SignallingPriorityIndication == nil {
				c.SignallingPriorityIndication = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.IPAddress:
		switch i.Instance() {
		case 0:
			if c.UELocalIPAddress == nil {
				c.UELocalIPAddress = i
			}
		case 1:
			if c.HeNBLocalIPAddress == nil {
				c.HeNBLocalIPAddress = i
			}
		case 2:
			if c.MMESGSNIdentifier == nil {
				c.MMESGSNIdentifier = i
			}
		case 3:
			if c.EPDGIPAddress == nil {
				c.EPDGIPAddress = i
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
			if c.HeNBUDPPort == nil {
				c.HeNBUDPPort = i
			}
		case 2:
			if c.UETCPPort == nil {
				c.UETCPPort = i
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
	case ie.CNOperatorSelectionEntity:
		switch i.Instance() {
		case 0:
			if c.CNOperatorSelectionEntity == nil {
				c.CNOperatorSelectionEntity = i
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
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			if c.MMESGSNOverloadControlInformation == nil {
				c.MMESGSNOverloadControlInformation = i
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
	case ie.MillisecondTimeStamp:
		switch i.Instance() {
		case 0:
			if c.OriginationTimeStamp == nil {
				c.OriginationTimeStamp = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.IntegerNumber:
		switch i.Instance() {
		case 0:
			if c.MaximumWaitTime == nil {
				c.MaximumWaitTime = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.TWANIdentifierTimestamp:
		switch i.Instance() {
		case 0:
			if c.WLANLocationTimeStamp == nil {
				c.WLANLocationTimeStamp = i
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
	case ie.RemoteUEContext:
		switch i.Instance() {
		case 0:
			c.RemoteUEContextConnected = append(c.RemoteUEContextConnected, i)
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.NodeIdentifier:
		switch i.Instance() {
		case 0:
			if c.TGPPAAAServerIdentifier == nil {
				c.TGPPAAAServerIdentifier = i
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
	case ie.ServingPLMNRateControl:
		switch i.Instance() {
		case 0:
			if c.ServingPLMNRateControl == nil {
				c.ServingPLMNRateControl = i
			}
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
	case ie.MappedUEUsageType:
		switch i.Instance() {
		case 0:
			if c.MappedUEUsageType == nil {
				c.MappedUEUsageType = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.FullyQualifiedDomainName:
		switch i.Instance() {
		case 0:
			if c.SGWUNodeName == nil {
				c.SGWUNodeName = i
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
	case ie.UPFunctionSelectionIndicationFlags:
		switch i.Instance() {
		case 0:
			if c.UPFunctionSelectionIndicationFlags == nil {
				c.UPFunctionSelectionIndicationFlags = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.APNRateControlStatus:
		switch i.Instance() {
		case 0:
			if c.APNRateControlStatus == nil {
				c.APNRateControlStatus = i
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
func (c *CreateSessionRequest) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *CreateSessionRequest) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *CreateSessionRequest) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *CreateSessionRequest) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *CreateSessionRequest) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *CreateSessionRequest) HasMessagePriority() bool { return c.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *CreateSessionRequest) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *CreateSessionRequest) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *CreateSessionRequest) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *CreateSessionRequest) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *CreateSessionRequest) SetPiggybacking(value uint8) { c.GetHeader().SetPiggybacking(value) }
