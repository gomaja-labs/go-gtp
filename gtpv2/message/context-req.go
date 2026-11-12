// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// ContextRequest is a ContextRequest Header and its IEs above.
type ContextRequest struct {
	*Header
	IMSI                               *ie.IE
	GUTI                               *ie.IE
	RAI                                *ie.IE
	PTMSI                              *ie.IE
	PTMSISignature                     *ie.IE
	CompleteTAURequestMessage          *ie.IE
	AddressAndTEIDForCPlane            *ie.IE
	UDPSourcePortNumber                *ie.IE
	RATType                            *ie.IE
	Indication                         *ie.IE
	HopCounter                         *ie.IE
	TargetPLMNID                       *ie.IE
	MMESGSNLDN                         *ie.IE
	SGSNNodeName                       *ie.IE
	MMENodeName                        *ie.IE
	SGSNNumber                         *ie.IE
	SGSNIdentifier                     *ie.IE
	MMEIdentifier                      *ie.IE
	CIoTOptimizationsSupportIndication *ie.IE
	PrivateExtension                   *ie.IE
	AdditionalIEs                      []*ie.IE
}

// NewContextRequest creates a new ContextRequest.
func NewContextRequest(teid, seq uint32, ies ...*ie.IE) *ContextRequest {
	c := &ContextRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeContextRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes ContextRequest into bytes.
func (c *ContextRequest) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ContextRequest into bytes.
func (c *ContextRequest) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeContextRequest, c.MarshalLen()); err != nil {
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
	if i := c.GUTI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "GUTI", iePair{ie.GUTI, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.RAI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "RAI", iePair{ie.UserLocationInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PTMSI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PTMSI", iePair{ie.PacketTMSI, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.PTMSISignature; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PTMSISignature", iePair{ie.PTMSISignature, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.CompleteTAURequestMessage; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "CompleteTAURequestMessage", iePair{ie.CompleteRequestMessage, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.AddressAndTEIDForCPlane; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "AddressAndTEIDForCPlane", iePair{ie.FullyQualifiedTEID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.UDPSourcePortNumber; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "UDPSourcePortNumber", iePair{ie.PortNumber, 0}); err != nil {
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
	if i := c.Indication; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "Indication", iePair{ie.Indication, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.HopCounter; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "HopCounter", iePair{ie.HopCounter, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.TargetPLMNID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "TargetPLMNID", iePair{ie.ServingNetwork, 0}); err != nil {
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
	if i := c.SGSNNodeName; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGSNNodeName", iePair{ie.FullyQualifiedDomainName, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MMENodeName; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMENodeName", iePair{ie.FullyQualifiedDomainName, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGSNNumber; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGSNNumber", iePair{ie.NodeNumber, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGSNIdentifier; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGSNIdentifier", iePair{ie.NodeIdentifier, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MMEIdentifier; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMEIdentifier", iePair{ie.NodeIdentifier, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.CIoTOptimizationsSupportIndication; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "CIoTOptimizationsSupportIndication", iePair{ie.CIoTOptimizationsSupportIndication, 0}); err != nil {
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

// ParseContextRequest decodes given bytes as ContextRequest.
func ParseContextRequest(b []byte) (*ContextRequest, error) {
	c := &ContextRequest{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as ContextRequest.
func (c *ContextRequest) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = ContextRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeContextRequest {
		return &UnexpectedTypeError{Expected: MsgTypeContextRequest, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *ContextRequest) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	l := c.Header.MarshalLen() - len(c.Header.Payload)
	if i := c.IMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.GUTI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.RAI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PTMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PTMSISignature; i != nil {
		l += i.MarshalLen()
	}
	if i := c.CompleteTAURequestMessage; i != nil {
		l += i.MarshalLen()
	}
	if i := c.AddressAndTEIDForCPlane; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UDPSourcePortNumber; i != nil {
		l += i.MarshalLen()
	}
	if i := c.RATType; i != nil {
		l += i.MarshalLen()
	}
	if i := c.Indication; i != nil {
		l += i.MarshalLen()
	}
	if i := c.HopCounter; i != nil {
		l += i.MarshalLen()
	}
	if i := c.TargetPLMNID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMESGSNLDN; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGSNNodeName; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMENodeName; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGSNNumber; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGSNIdentifier; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMEIdentifier; i != nil {
		l += i.MarshalLen()
	}
	if i := c.CIoTOptimizationsSupportIndication; i != nil {
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
func (c *ContextRequest) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *ContextRequest) MessageTypeName() string {
	return "Context Request"
}

// TEID returns the TEID in uint32.
func (c *ContextRequest) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *ContextRequest) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.3.5-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *ContextRequest) addIE(i *ie.IE) {
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
	case ie.GUTI:
		switch i.Instance() {
		case 0:
			if c.GUTI == nil {
				c.GUTI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.UserLocationInformation:
		switch i.Instance() {
		case 0:
			if c.RAI == nil {
				c.RAI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.PacketTMSI:
		switch i.Instance() {
		case 0:
			if c.PTMSI == nil {
				c.PTMSI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.PTMSISignature:
		switch i.Instance() {
		case 0:
			if c.PTMSISignature == nil {
				c.PTMSISignature = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.CompleteRequestMessage:
		switch i.Instance() {
		case 0:
			if c.CompleteTAURequestMessage == nil {
				c.CompleteTAURequestMessage = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.FullyQualifiedTEID:
		switch i.Instance() {
		case 0:
			if c.AddressAndTEIDForCPlane == nil {
				c.AddressAndTEIDForCPlane = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.PortNumber:
		switch i.Instance() {
		case 0:
			if c.UDPSourcePortNumber == nil {
				c.UDPSourcePortNumber = i
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
			if c.Indication == nil {
				c.Indication = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.HopCounter:
		switch i.Instance() {
		case 0:
			if c.HopCounter == nil {
				c.HopCounter = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.ServingNetwork:
		switch i.Instance() {
		case 0:
			if c.TargetPLMNID == nil {
				c.TargetPLMNID = i
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
	case ie.FullyQualifiedDomainName:
		switch i.Instance() {
		case 0:
			if c.SGSNNodeName == nil {
				c.SGSNNodeName = i
			}
		case 1:
			if c.MMENodeName == nil {
				c.MMENodeName = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.NodeNumber:
		switch i.Instance() {
		case 0:
			if c.SGSNNumber == nil {
				c.SGSNNumber = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.NodeIdentifier:
		switch i.Instance() {
		case 0:
			if c.SGSNIdentifier == nil {
				c.SGSNIdentifier = i
			}
		case 1:
			if c.MMEIdentifier == nil {
				c.MMEIdentifier = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.CIoTOptimizationsSupportIndication:
		switch i.Instance() {
		case 0:
			if c.CIoTOptimizationsSupportIndication == nil {
				c.CIoTOptimizationsSupportIndication = i
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
func (c *ContextRequest) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *ContextRequest) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *ContextRequest) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *ContextRequest) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *ContextRequest) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *ContextRequest) HasMessagePriority() bool { return c.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *ContextRequest) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *ContextRequest) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *ContextRequest) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *ContextRequest) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *ContextRequest) SetPiggybacking(value uint8) { c.GetHeader().SetPiggybacking(value) }
