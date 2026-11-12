// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// ResumeNotification is a ResumeNotification Header and its IEs.
type ResumeNotification struct {
	*Header
	IMSI                       *ie.IE
	LinkedEPSBearerID          *ie.IE
	OriginatingNode            *ie.IE
	SenderFTEIDForControlPlane *ie.IE
	PrivateExtension           *ie.IE
	AdditionalIEs              []*ie.IE
}

// NewResumeNotification creates a new ResumeNotification.
func NewResumeNotification(teid, seq uint32, ies ...*ie.IE) *ResumeNotification {
	c := &ResumeNotification{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeResumeNotification, teid, seq, nil,
		),
	}

	for _, i := range ies {
		c.addIE(i)
	}

	c.SetLength()
	return c
}

// Marshal serializes ResumeNotification into bytes.
func (c *ResumeNotification) Marshal() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes ResumeNotification into bytes.
func (c *ResumeNotification) MarshalTo(b []byte) error {
	if c == nil || c.Header == nil {
		return ErrNilMessage
	}
	if err := c.Header.validate(MsgTypeResumeNotification, c.MarshalLen()); err != nil {
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
	if i := c.LinkedEPSBearerID; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "LinkedEPSBearerID", iePair{ie.EPSBearerID, 0}); err != nil {
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
	if i := c.SenderFTEIDForControlPlane; i != nil {
		if err := marshalNamedIE(i, c.Payload[offset:], "SenderFTEIDForControlPlane", iePair{ie.FullyQualifiedTEID, 0}); err != nil {
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

// ParseResumeNotification decodes given bytes as ResumeNotification.
func ParseResumeNotification(b []byte) (*ResumeNotification, error) {
	c := &ResumeNotification{}
	if err := c.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return c, err
		}
		return nil, err
	}
	return c, nil
}

// UnmarshalBinary decodes given bytes as ResumeNotification.
func (c *ResumeNotification) UnmarshalBinary(b []byte) error {
	if c == nil {
		return ErrNilMessage
	}
	*c = ResumeNotification{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeResumeNotification {
		return &UnexpectedTypeError{Expected: MsgTypeResumeNotification, Actual: h.Type}
	}
	c.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, c.addIE)
}

// MarshalLen returns the serial length in int.
func (c *ResumeNotification) MarshalLen() int {
	if c == nil || c.Header == nil {
		return 0
	}
	l := c.Header.MarshalLen() - len(c.Header.Payload)
	if i := c.IMSI; i != nil {
		l += i.MarshalLen()
	}
	if i := c.LinkedEPSBearerID; i != nil {
		l += i.MarshalLen()
	}
	if i := c.OriginatingNode; i != nil {
		l += i.MarshalLen()
	}
	if i := c.SenderFTEIDForControlPlane; i != nil {
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
func (c *ResumeNotification) SetLength() {
	if c != nil && c.Header != nil {
		c.Header.setLength(c.MarshalLen() - c.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (c *ResumeNotification) MessageTypeName() string {
	return "Resume Notification"
}

// TEID returns the TEID in uint32.
func (c *ResumeNotification) TEID() uint32 {
	if c == nil {
		return 0
	}
	return c.Header.teid()
}

func (c *ResumeNotification) GetHeader() *Header {
	if c == nil {
		return nil
	}
	return c.Header
}

// TS 29.274 Table 7.4.3-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (c *ResumeNotification) addIE(i *ie.IE) {
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
	case ie.EPSBearerID:
		switch i.Instance() {
		case 0:
			if c.LinkedEPSBearerID == nil {
				c.LinkedEPSBearerID = i
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
	case ie.FullyQualifiedTEID:
		switch i.Instance() {
		case 0:
			if c.SenderFTEIDForControlPlane == nil {
				c.SenderFTEIDForControlPlane = i
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
func (c *ResumeNotification) Version() int { return c.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (c *ResumeNotification) MessageType() uint8 { return c.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (c *ResumeNotification) Sequence() uint32 { return c.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (c *ResumeNotification) IsPiggybacking() bool { return c.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (c *ResumeNotification) HasTEID() bool { return c.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (c *ResumeNotification) HasMessagePriority() bool { return c.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (c *ResumeNotification) MessagePriority() uint8 { return c.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (c *ResumeNotification) SetSequenceNumber(seq uint32) error {
	return c.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (c *ResumeNotification) SetMessagePriority(priority uint8) error {
	return c.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (c *ResumeNotification) SetTEID(teid uint32) { c.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (c *ResumeNotification) SetPiggybacking(value uint8) { c.GetHeader().SetPiggybacking(value) }
