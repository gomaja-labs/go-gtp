// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// ContextAcknowledge is a ContextAcknowledge Header and its IEs above.
type ContextAcknowledge struct {
	*Header
	Cause                  *ie.IE
	IndicationFlags        *ie.IE
	ForwardingFTEID        *ie.IE
	BearerContexts         []*ie.IE
	SGSNNumber             *ie.IE
	MMENumberForMTSMS      *ie.IE
	SGSNIdentifierForMTSMS *ie.IE
	MMEIdentifierForMTSMS  *ie.IE
	PrivateExtension       *ie.IE
	AdditionalIEs          []*ie.IE
}

// NewContextAcknowledge creates a new ContextAcknowledge.
func NewContextAcknowledge(teid, seq uint32, ies ...*ie.IE) *ContextAcknowledge {
	c := &ContextAcknowledge{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeContextAcknowledge, teid, seq, nil,
		),
	}

	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes ContextAcknowledge into bytes.
func (c *ContextAcknowledge) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ContextAcknowledge into bytes.
func (c *ContextAcknowledge) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeContextAcknowledge, c.MarshalLen()); err != nil {
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
	if i := c.IndicationFlags; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "IndicationFlags", iePair{ie.Indication, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.ForwardingFTEID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "ForwardingFTEID", iePair{ie.FullyQualifiedTEID, 0}); err != nil {
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
	if i := c.SGSNNumber; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGSNNumber", iePair{ie.NodeNumber, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MMENumberForMTSMS; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMENumberForMTSMS", iePair{ie.NodeNumber, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.SGSNIdentifierForMTSMS; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SGSNIdentifierForMTSMS", iePair{ie.NodeIdentifier, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.MMEIdentifierForMTSMS; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "MMEIdentifierForMTSMS", iePair{ie.NodeIdentifier, 1}); err != nil {
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

// ParseContextAcknowledge decodes given bytes as ContextAcknowledge.
func ParseContextAcknowledge(b []byte) (*ContextAcknowledge, error) {
	c := &ContextAcknowledge{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as ContextAcknowledge.
func (c *ContextAcknowledge) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = ContextAcknowledge{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeContextAcknowledge {
		return &UnexpectedTypeError{Expected: MsgTypeContextAcknowledge, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *ContextAcknowledge) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	l := c.Header.MarshalLen() - len(c.Header.Payload)
	if i := c.Cause; i != nil {
		l += i.MarshalLen()
	}
	if i := c.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := c.ForwardingFTEID; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range c.BearerContexts {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := c.SGSNNumber; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMENumberForMTSMS; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SGSNIdentifierForMTSMS; i != nil {
		l += i.MarshalLen()
	}
	if i := c.MMEIdentifierForMTSMS; i != nil {
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
func (c *ContextAcknowledge) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *ContextAcknowledge) MessageTypeName() string {
	return "Context Acknowledge"
}

// TEID returns the TEID in uint32.
func (c *ContextAcknowledge) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *ContextAcknowledge) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.3.7-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *ContextAcknowledge) addIE(i *ie.IE) {
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
			if c.ForwardingFTEID == nil {
				c.ForwardingFTEID = i
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
	case ie.NodeNumber:
		switch i.Instance() {
		case 0:
			if c.SGSNNumber == nil {
				c.SGSNNumber = i
			}
		case 1:
			if c.MMENumberForMTSMS == nil {
				c.MMENumberForMTSMS = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.NodeIdentifier:
		switch i.Instance() {
		case 0:
			if c.SGSNIdentifierForMTSMS == nil {
				c.SGSNIdentifierForMTSMS = i
			}
		case 1:
			if c.MMEIdentifierForMTSMS == nil {
				c.MMEIdentifierForMTSMS = i
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
func (c *ContextAcknowledge) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *ContextAcknowledge) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *ContextAcknowledge) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *ContextAcknowledge) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *ContextAcknowledge) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *ContextAcknowledge) HasMessagePriority() bool { return c.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *ContextAcknowledge) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *ContextAcknowledge) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *ContextAcknowledge) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *ContextAcknowledge) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *ContextAcknowledge) SetPiggybacking(value uint8) { c.GetHeader().SetPiggybacking(value) }
