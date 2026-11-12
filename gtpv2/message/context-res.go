// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// ContextResponse is a ContextResponse Header and its IEs above.
type ContextResponse struct {
	*Header
	Cause                               *ie.IE
	IMSI                                *ie.IE
	UEMMContext                         *ie.IE
	UEPDNConnections                    []*ie.IE
	SenderFTEID                         *ie.IE
	SGWS11S4FTEID                       *ie.IE
	SGWNodeName                         *ie.IE
	IndicationFlags                     *ie.IE
	TraceInformation                    *ie.IE
	S101IPAddress                       *ie.IE
	S102IPAddress                       *ie.IE
	SubscribedRFSPIndex                 *ie.IE
	RFSPIndexInUse                      *ie.IE
	UETimeZone                          *ie.IE
	MMESGSNLDN                          *ie.IE
	MDTConfiguration                    *ie.IE
	SGSNNodeName                        *ie.IE
	MMENodeName                         *ie.IE
	UCI                                 *ie.IE
	MonitoringEventInformation          []*ie.IE
	MonitoringEventExtensionInformation []*ie.IE
	UEUsageType                         *ie.IE
	SCEFPDNConnection                   []*ie.IE
	RATType                             *ie.IE
	ServingPLMNRateControl              *ie.IE
	MOExceptionDataCounter              *ie.IE
	RemainingRunningServiceGapTimer     *ie.IE
	ExtendedTraceInformation            *ie.IE
	SubscribedAdditionalRRMPolicyIndex  *ie.IE
	AdditionalRRMPolicyIndexInUse       *ie.IE
	PrivateExtension                    *ie.IE
	AdditionalIEs                       []*ie.IE
}

// NewContextResponse creates a new ContextResponse.
func NewContextResponse(teid, seq uint32, ies ...*ie.IE) *ContextResponse {
	c := &ContextResponse{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeContextResponse, teid, seq, nil,
		),
	}

	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes ContextResponse into bytes.
func (c *ContextResponse) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ContextResponse into bytes.
func (c *ContextResponse) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeContextResponse, c.MarshalLen()); err != nil {
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
	if i := c.IMSI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "IMSI", iePair{ie.IMSI, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.UEMMContext; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UEMMContext", iePair{ie.MMContextEPSSecurityContextQuadrupletsAndQuintuplets, 0}, iePair{ie.MMContextGSMKeyAndTriplets, 0}, iePair{ie.MMContextGSMKeyUsedCipherAndQuintuplets, 0}, iePair{ie.MMContextUMTSKeyAndQuintuplets, 0}, iePair{ie.MMContextUMTSKeyQuadrupletsAndQuintuplets, 0}, iePair{ie.MMContextUMTSKeyUsedCipherAndQuintuplets, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.UEPDNConnections {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "UEPDNConnections", iePair{ie.PDNConnection, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SenderFTEID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SenderFTEID", iePair{ie.FullyQualifiedTEID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGWS11S4FTEID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGWS11S4FTEID", iePair{ie.FullyQualifiedTEID, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGWNodeName; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGWNodeName", iePair{ie.FullyQualifiedDomainName, 0}); err != nil {
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
	if i := c.TraceInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TraceInformation", iePair{ie.TraceInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.S101IPAddress; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "S101IPAddress", iePair{ie.IPAddress, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.S102IPAddress; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "S102IPAddress", iePair{ie.IPAddress, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SubscribedRFSPIndex; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SubscribedRFSPIndex", iePair{ie.RFSPIndex, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.RFSPIndexInUse; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "RFSPIndexInUse", iePair{ie.RFSPIndex, 1}); err != nil {
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
	if i := c.MMESGSNLDN; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMESGSNLDN", iePair{ie.LocalDistinguishedName, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MDTConfiguration; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MDTConfiguration", iePair{ie.MDTConfiguration, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGSNNodeName; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGSNNodeName", iePair{ie.FullyQualifiedDomainName, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MMENodeName; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMENodeName", iePair{ie.FullyQualifiedDomainName, 2}); err != nil {
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
	for _, i := range c.MonitoringEventInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "MonitoringEventInformation", iePair{ie.MonitoringEventInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.MonitoringEventExtensionInformation {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "MonitoringEventExtensionInformation", iePair{ie.MonitoringEventExtensionInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.UEUsageType; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UEUsageType", iePair{ie.IntegerNumber, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range c.SCEFPDNConnection {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, c.Payload[offset:], "SCEFPDNConnection", iePair{ie.SCEFPDNConnection, 0}); err != nil {
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
	if i := c.RemainingRunningServiceGapTimer; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "RemainingRunningServiceGapTimer", iePair{ie.IntegerNumber, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.ExtendedTraceInformation; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ExtendedTraceInformation", iePair{ie.ExtendedTraceInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SubscribedAdditionalRRMPolicyIndex; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SubscribedAdditionalRRMPolicyIndex", iePair{ie.AdditionalRRMPolicyIndex, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.AdditionalRRMPolicyIndexInUse; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "AdditionalRRMPolicyIndexInUse", iePair{ie.AdditionalRRMPolicyIndex, 1}); err != nil {
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

// ParseContextResponse decodes given bytes as ContextResponse.
func ParseContextResponse(b []byte) (*ContextResponse, error) {
	c := &ContextResponse{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as ContextResponse.
func (c *ContextResponse) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = ContextResponse{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeContextResponse {
		return &UnexpectedTypeError{Expected: MsgTypeContextResponse, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *ContextResponse) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	l := c.Header.MarshalLen() - len(c.Header.Payload)
	if i := c.Cause; i != nil {
		l += i.MarshalLen()
	}
	if i := c.IMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UEMMContext; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.UEPDNConnections {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.SenderFTEID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGWS11S4FTEID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGWNodeName; i != nil {
		l += i.MarshalLen()
	}
	if i := c.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := c.TraceInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := c.S101IPAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := c.S102IPAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SubscribedRFSPIndex; i != nil {
		l += i.MarshalLen()
	}
	if i := c.RFSPIndexInUse; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UETimeZone; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMESGSNLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MDTConfiguration; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGSNNodeName; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMENodeName; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UCI; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.MonitoringEventInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	for _, i := range c.MonitoringEventExtensionInformation {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.UEUsageType; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.SCEFPDNConnection {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.RATType; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ServingPLMNRateControl; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MOExceptionDataCounter; i != nil {
		l += i.MarshalLen()
	}
	if i := c.RemainingRunningServiceGapTimer; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ExtendedTraceInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SubscribedAdditionalRRMPolicyIndex; i != nil {
		l += i.MarshalLen()
	}
	if i := c.AdditionalRRMPolicyIndexInUse; i != nil {
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
func (c *ContextResponse) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *ContextResponse) MessageTypeName() string {
	return "Context Response"
}

// TEID returns the TEID in uint32.
func (c *ContextResponse) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *ContextResponse) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.3.6-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *ContextResponse) addIE(i *ie.IE) {
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
	case ie.IMSI:
		switch i.Instance() {
		case 0:
			if c.IMSI == nil {
				c.IMSI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.MMContextEPSSecurityContextQuadrupletsAndQuintuplets:
		switch i.Instance() {
		case 0:
			if c.UEMMContext == nil {
				c.UEMMContext = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.MMContextGSMKeyAndTriplets:
		switch i.Instance() {
		case 0:
			if c.UEMMContext == nil {
				c.UEMMContext = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.MMContextGSMKeyUsedCipherAndQuintuplets:
		switch i.Instance() {
		case 0:
			if c.UEMMContext == nil {
				c.UEMMContext = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.MMContextUMTSKeyAndQuintuplets:
		switch i.Instance() {
		case 0:
			if c.UEMMContext == nil {
				c.UEMMContext = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.MMContextUMTSKeyQuadrupletsAndQuintuplets:
		switch i.Instance() {
		case 0:
			if c.UEMMContext == nil {
				c.UEMMContext = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.MMContextUMTSKeyUsedCipherAndQuintuplets:
		switch i.Instance() {
		case 0:
			if c.UEMMContext == nil {
				c.UEMMContext = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.PDNConnection:
		switch i.Instance() {
		case 0:
			c.UEPDNConnections = append(c.UEPDNConnections, i)
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.FullyQualifiedTEID:
		switch i.Instance() {
		case 0:
			if c.SenderFTEID == nil {
				c.SenderFTEID = i
			}
		case 1:
			if c.SGWS11S4FTEID == nil {
				c.SGWS11S4FTEID = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.FullyQualifiedDomainName:
		switch i.Instance() {
		case 0:
			if c.SGWNodeName == nil {
				c.SGWNodeName = i
			}
		case 1:
			if c.SGSNNodeName == nil {
				c.SGSNNodeName = i
			}
		case 2:
			if c.MMENodeName == nil {
				c.MMENodeName = i
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
	case ie.TraceInformation:
		switch i.Instance() {
		case 0:
			if c.TraceInformation == nil {
				c.TraceInformation = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.IPAddress:
		switch i.Instance() {
		case 0:
			if c.S101IPAddress == nil {
				c.S101IPAddress = i
			}
		case 1:
			if c.S102IPAddress == nil {
				c.S102IPAddress = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.RFSPIndex:
		switch i.Instance() {
		case 0:
			if c.SubscribedRFSPIndex == nil {
				c.SubscribedRFSPIndex = i
			}
		case 1:
			if c.RFSPIndexInUse == nil {
				c.RFSPIndexInUse = i
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
	case ie.LocalDistinguishedName:
		switch i.Instance() {
		case 0:
			if c.MMESGSNLDN == nil {
				c.MMESGSNLDN = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.MDTConfiguration:
		switch i.Instance() {
		case 0:
			if c.MDTConfiguration == nil {
				c.MDTConfiguration = i
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
	case ie.MonitoringEventInformation:
		switch i.Instance() {
		case 0:
			c.MonitoringEventInformation = append(c.MonitoringEventInformation, i)
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.MonitoringEventExtensionInformation:
		switch i.Instance() {
		case 0:
			c.MonitoringEventExtensionInformation = append(c.MonitoringEventExtensionInformation, i)
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.IntegerNumber:
		switch i.Instance() {
		case 0:
			if c.UEUsageType == nil {
				c.UEUsageType = i
			}
		case 1:
			if c.RemainingRunningServiceGapTimer == nil {
				c.RemainingRunningServiceGapTimer = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.SCEFPDNConnection:
		switch i.Instance() {
		case 0:
			c.SCEFPDNConnection = append(c.SCEFPDNConnection, i)
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
	case ie.ExtendedTraceInformation:
		switch i.Instance() {
		case 0:
			if c.ExtendedTraceInformation == nil {
				c.ExtendedTraceInformation = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.AdditionalRRMPolicyIndex:
		switch i.Instance() {
		case 0:
			if c.SubscribedAdditionalRRMPolicyIndex == nil {
				c.SubscribedAdditionalRRMPolicyIndex = i
			}
		case 1:
			if c.AdditionalRRMPolicyIndexInUse == nil {
				c.AdditionalRRMPolicyIndexInUse = i
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
func (c *ContextResponse) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *ContextResponse) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *ContextResponse) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *ContextResponse) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *ContextResponse) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *ContextResponse) HasMessagePriority() bool { return c.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *ContextResponse) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *ContextResponse) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *ContextResponse) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *ContextResponse) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *ContextResponse) SetPiggybacking(value uint8) { c.GetHeader().SetPiggybacking(value) }
