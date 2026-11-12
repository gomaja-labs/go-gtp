// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// SuspendNotification is a SuspendNotification Header and its IEs above.
type SuspendNotification struct {
	*Header
	IMSI                   *ie.IE
	RAI                    *ie.IE
	LinkedEBI              *ie.IE
	PTMSI                  *ie.IE
	OriginatingNode        *ie.IE
	AddressForControlPlane *ie.IE
	UDPSourcePortNumber    *ie.IE
	HopCounter             *ie.IE
	SenderFTEIDC           *ie.IE
	PrivateExtension       *ie.IE
	AdditionalIEs          []*ie.IE
}

// NewSuspendNotification creates a new SuspendNotification.
func NewSuspendNotification(teid, seq uint32, ies ...*ie.IE) *SuspendNotification {
	c := &SuspendNotification{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeSuspendNotification, teid, seq, nil,
		),
	}

	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes SuspendNotification into bytes.
func (c *SuspendNotification) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes SuspendNotification into bytes.
func (c *SuspendNotification) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeSuspendNotification, c.MarshalLen()); err != nil {
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
	if i := c.RAI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "RAI", iePair{ie.UserLocationInformation, 0}); err != nil {
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
	if i := c.PTMSI; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "PTMSI", iePair{ie.PacketTMSI, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.OriginatingNode; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "OriginatingNode", iePair{ie.NodeType, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := c.AddressForControlPlane; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "AddressForControlPlane", iePair{ie.IPAddress, 0}); err != nil {
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
	if i := c.HopCounter; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "HopCounter", iePair{ie.HopCounter, 0}); err != nil {
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

// ParseSuspendNotification decodes given bytes as SuspendNotification.
func ParseSuspendNotification(b []byte) (*SuspendNotification, error) {
	c := &SuspendNotification{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as SuspendNotification.
func (c *SuspendNotification) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = SuspendNotification{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeSuspendNotification {
		return &UnexpectedTypeError{Expected: MsgTypeSuspendNotification, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *SuspendNotification) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	l := c.Header.MarshalLen() - len(c.Header.Payload)
	if i := c.IMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.RAI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.LinkedEBI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.PTMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.OriginatingNode; i != nil {
		l += i.MarshalLen()
	}
	if i := c.AddressForControlPlane; i != nil {
		l += i.MarshalLen()
	}
	if i := c.UDPSourcePortNumber; i != nil {
		l += i.MarshalLen()
	}
	if i := c.HopCounter; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SenderFTEIDC; i != nil {
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
func (c *SuspendNotification) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *SuspendNotification) MessageTypeName() string {
	return "Suspend Notification"
}

// TEID returns the TEID in uint32.
func (c *SuspendNotification) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *SuspendNotification) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.4.1-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *SuspendNotification) addIE(i *ie.IE) {
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
	case ie.UserLocationInformation:
		switch i.Instance() {
		case 0:
			if c.RAI == nil {
				c.RAI = i
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
	case ie.PacketTMSI:
		switch i.Instance() {
		case 0:
			if c.PTMSI == nil {
				c.PTMSI = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.NodeType:
		switch i.Instance() {
		case 0:
			if c.OriginatingNode == nil {
				c.OriginatingNode = i
			}
		default:
			c.AdditionalIEs = append(c.AdditionalIEs, i)
		}
	case ie.IPAddress:
		switch i.Instance() {
		case 0:
			if c.AddressForControlPlane == nil {
				c.AddressForControlPlane = i
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
	case ie.HopCounter:
		switch i.Instance() {
		case 0:
			if c.HopCounter == nil {
				c.HopCounter = i
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
func (c *SuspendNotification) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *SuspendNotification) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *SuspendNotification) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *SuspendNotification) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *SuspendNotification) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *SuspendNotification) HasMessagePriority() bool { return c.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *SuspendNotification) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *SuspendNotification) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *SuspendNotification) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *SuspendNotification) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *SuspendNotification) SetPiggybacking(value uint8) { c.GetHeader().SetPiggybacking(value) }
