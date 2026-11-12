// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// DeleteBearerResponse is a DeleteBearerResponse Header and its IEs above.
type DeleteBearerResponse struct {
	*Header
	Cause                               *ie.IE
	LinkedEBI                           *ie.IE
	BearerContexts                      []*ie.IE
	Recovery                            *ie.IE
	MMEFQCSID                           *ie.IE
	SGWFQCSID                           *ie.IE
	EPDGFQCSID                          *ie.IE
	TWANFQCSID                          *ie.IE
	PCO                                 *ie.IE
	UETimeZone                          *ie.IE
	ULI                                 *ie.IE
	ULITimestamp                        *ie.IE
	TWANIdentifier                      *ie.IE
	TWANIdentifierTimestamp             *ie.IE
	MMEOverloadControlInformation       *ie.IE
	SGWOverloadControlInformation       *ie.IE
	MMESGSNIdentifierOrUELocalIPAddress *ie.IE
	TWANePDGOverloadControlInformation  *ie.IE
	WLANLocationInformation             *ie.IE
	WLANLocationTimestamp               *ie.IE
	UEUDPPort                           *ie.IE
	NBIFOMContainer                     *ie.IE
	UETCPPort                           *ie.IE
	SecondaryRATUsageDataReport         []*ie.IE
	PrivateExtension                    *ie.IE
	AdditionalIEs                       []*ie.IE
}

// NewDeleteBearerResponse creates a new DeleteBearerResponse.
func NewDeleteBearerResponse(teid, seq uint32, ies ...*ie.IE) *DeleteBearerResponse {
	d := &DeleteBearerResponse{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeDeleteBearerResponse, teid, seq, nil,
		),
	}

	for _, i := range ies {
		d.addIE(i)
	}

	d.SetLength()
	return d
}

// Marshal serializes DeleteBearerResponse into bytes.
func (d *DeleteBearerResponse) Marshal() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes DeleteBearerResponse into bytes.
func (d *DeleteBearerResponse) MarshalTo(b []byte) error {
	if d == nil || d.Header == nil {
		return ErrNilMessage
	}
	if err := d.Header.validate(MsgTypeDeleteBearerResponse, d.MarshalLen()); err != nil {
		return err
	}
	if len(b) < d.MarshalLen() {
		return ErrTooShortToMarshal
	}
	d.Header.Payload = nil
	d.Header.Payload = make([]byte, d.MarshalLen()-d.Header.MarshalLen())
	offset := 0
	if i := d.Cause; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "Cause", iePair{ie.Cause, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.LinkedEBI; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "LinkedEBI", iePair{ie.EPSBearerID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range d.BearerContexts {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "BearerContexts", iePair{ie.BearerContext, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.Recovery; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "Recovery", iePair{ie.Recovery, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.MMEFQCSID; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "MMEFQCSID", iePair{ie.FullyQualifiedCSID, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.SGWFQCSID; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SGWFQCSID", iePair{ie.FullyQualifiedCSID, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.EPDGFQCSID; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "EPDGFQCSID", iePair{ie.FullyQualifiedCSID, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.TWANFQCSID; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "TWANFQCSID", iePair{ie.FullyQualifiedCSID, 3}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.PCO; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "PCO", iePair{ie.ProtocolConfigurationOptions, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.UETimeZone; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "UETimeZone", iePair{ie.UETimeZone, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.ULI; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "ULI", iePair{ie.UserLocationInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.ULITimestamp; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "ULITimestamp", iePair{ie.ULITimestamp, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.TWANIdentifier; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "TWANIdentifier", iePair{ie.TWANIdentifier, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.TWANIdentifierTimestamp; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "TWANIdentifierTimestamp", iePair{ie.TWANIdentifierTimestamp, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.MMEOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "MMEOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.SGWOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.MMESGSNIdentifierOrUELocalIPAddress; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "MMESGSNIdentifierOrUELocalIPAddress", iePair{ie.IPAddress, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.TWANePDGOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "TWANePDGOverloadControlInformation", iePair{ie.OverloadControlInformation, 2}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.WLANLocationInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "WLANLocationInformation", iePair{ie.TWANIdentifier, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.WLANLocationTimestamp; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "WLANLocationTimestamp", iePair{ie.TWANIdentifierTimestamp, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.UEUDPPort; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "UEUDPPort", iePair{ie.PortNumber, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.NBIFOMContainer; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "NBIFOMContainer", iePair{ie.FContainer, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.UETCPPort; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "UETCPPort", iePair{ie.PortNumber, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range d.SecondaryRATUsageDataReport {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "SecondaryRATUsageDataReport", iePair{ie.SecondaryRATUsageDataReport, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.PrivateExtension; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "PrivateExtension", iePair{ie.PrivateExtension, 255}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	for _, i := range d.AdditionalIEs {
		if i == nil {
			continue
		}
		if err := i.MarshalTo(d.Payload[offset:]); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	d.Header.SetLength()
	return d.Header.MarshalTo(b)
}

// ParseDeleteBearerResponse decodes given bytes as DeleteBearerResponse.
func ParseDeleteBearerResponse(b []byte) (*DeleteBearerResponse, error) {
	d := &DeleteBearerResponse{}
	if err := d.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return d, err
		}
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary decodes given bytes as DeleteBearerResponse.
func (d *DeleteBearerResponse) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrNilMessage
	}
	*d = DeleteBearerResponse{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeDeleteBearerResponse {
		return &UnexpectedTypeError{Expected: MsgTypeDeleteBearerResponse, Actual: h.Type}
	}
	d.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, d.addIE)
}

// MarshalLen returns the serial length in int.
func (d *DeleteBearerResponse) MarshalLen() int {
	if d == nil || d.Header == nil {
		return 0
	}
	l := d.Header.MarshalLen() - len(d.Header.Payload)
	if i := d.Cause; i != nil {
		l += i.MarshalLen()
	}
	if i := d.LinkedEBI; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.BearerContexts {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := d.Recovery; i != nil {
		l += i.MarshalLen()
	}
	if i := d.MMEFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := d.SGWFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := d.EPDGFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := d.TWANFQCSID; i != nil {
		l += i.MarshalLen()
	}
	if i := d.PCO; i != nil {
		l += i.MarshalLen()
	}
	if i := d.UETimeZone; i != nil {
		l += i.MarshalLen()
	}
	if i := d.ULI; i != nil {
		l += i.MarshalLen()
	}
	if i := d.ULITimestamp; i != nil {
		l += i.MarshalLen()
	}
	if i := d.TWANIdentifier; i != nil {
		l += i.MarshalLen()
	}
	if i := d.TWANIdentifierTimestamp; i != nil {
		l += i.MarshalLen()
	}
	if i := d.MMEOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := d.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := d.MMESGSNIdentifierOrUELocalIPAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := d.TWANePDGOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := d.WLANLocationInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := d.WLANLocationTimestamp; i != nil {
		l += i.MarshalLen()
	}
	if i := d.UEUDPPort; i != nil {
		l += i.MarshalLen()
	}
	if i := d.NBIFOMContainer; i != nil {
		l += i.MarshalLen()
	}
	if i := d.UETCPPort; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.SecondaryRATUsageDataReport {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := d.PrivateExtension; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.AdditionalIEs {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	return l
}

// SetLength sets the length in Length field.
func (d *DeleteBearerResponse) SetLength() {
	if d != nil && d.Header != nil {
		d.Header.setLength(d.MarshalLen() - d.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (d *DeleteBearerResponse) MessageTypeName() string {
	return "Delete Bearer Response"
}

// TEID returns the TEID in uint32.
func (d *DeleteBearerResponse) TEID() uint32 {
	if d == nil {
		return 0
	}
	return d.Header.teid()
}

func (d *DeleteBearerResponse) GetHeader() *Header {
	if d == nil {
		return nil
	}
	return d.Header
}

// TS 29.274 Table 7.2.10.2-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (d *DeleteBearerResponse) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		d.Header.err = err
	}
	switch i.Type {
	case ie.Cause:
		switch i.Instance() {
		case 0:
			if d.Cause == nil {
				d.Cause = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.EPSBearerID:
		switch i.Instance() {
		case 0:
			if d.LinkedEBI == nil {
				d.LinkedEBI = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.BearerContext:
		switch i.Instance() {
		case 0:
			d.BearerContexts = append(d.BearerContexts, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.Recovery:
		switch i.Instance() {
		case 0:
			if d.Recovery == nil {
				d.Recovery = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.FullyQualifiedCSID:
		switch i.Instance() {
		case 0:
			if d.MMEFQCSID == nil {
				d.MMEFQCSID = i
			}
		case 1:
			if d.SGWFQCSID == nil {
				d.SGWFQCSID = i
			}
		case 2:
			if d.EPDGFQCSID == nil {
				d.EPDGFQCSID = i
			}
		case 3:
			if d.TWANFQCSID == nil {
				d.TWANFQCSID = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.ProtocolConfigurationOptions:
		switch i.Instance() {
		case 0:
			if d.PCO == nil {
				d.PCO = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.UETimeZone:
		switch i.Instance() {
		case 0:
			if d.UETimeZone == nil {
				d.UETimeZone = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.UserLocationInformation:
		switch i.Instance() {
		case 0:
			if d.ULI == nil {
				d.ULI = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.ULITimestamp:
		switch i.Instance() {
		case 0:
			if d.ULITimestamp == nil {
				d.ULITimestamp = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.TWANIdentifier:
		switch i.Instance() {
		case 0:
			if d.TWANIdentifier == nil {
				d.TWANIdentifier = i
			}
		case 1:
			if d.WLANLocationInformation == nil {
				d.WLANLocationInformation = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.TWANIdentifierTimestamp:
		switch i.Instance() {
		case 0:
			if d.TWANIdentifierTimestamp == nil {
				d.TWANIdentifierTimestamp = i
			}
		case 1:
			if d.WLANLocationTimestamp == nil {
				d.WLANLocationTimestamp = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			if d.MMEOverloadControlInformation == nil {
				d.MMEOverloadControlInformation = i
			}
		case 1:
			if d.SGWOverloadControlInformation == nil {
				d.SGWOverloadControlInformation = i
			}
		case 2:
			if d.TWANePDGOverloadControlInformation == nil {
				d.TWANePDGOverloadControlInformation = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.IPAddress:
		switch i.Instance() {
		case 0:
			if d.MMESGSNIdentifierOrUELocalIPAddress == nil {
				d.MMESGSNIdentifierOrUELocalIPAddress = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.PortNumber:
		switch i.Instance() {
		case 0:
			if d.UEUDPPort == nil {
				d.UEUDPPort = i
			}
		case 1:
			if d.UETCPPort == nil {
				d.UETCPPort = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.FContainer:
		switch i.Instance() {
		case 0:
			if d.NBIFOMContainer == nil {
				d.NBIFOMContainer = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.SecondaryRATUsageDataReport:
		switch i.Instance() {
		case 0:
			d.SecondaryRATUsageDataReport = append(d.SecondaryRATUsageDataReport, i)
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.PrivateExtension:
		if d.PrivateExtension == nil {
			d.PrivateExtension = i
		} else if d.PrivateExtension.Instance() != i.Instance() {
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	default:
		d.AdditionalIEs = append(d.AdditionalIEs, i)
	}
}

// Version reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerResponse) Version() int { return d.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerResponse) MessageType() uint8 { return d.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerResponse) Sequence() uint32 { return d.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerResponse) IsPiggybacking() bool { return d.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerResponse) HasTEID() bool { return d.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerResponse) HasMessagePriority() bool { return d.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerResponse) MessagePriority() uint8 { return d.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (d *DeleteBearerResponse) SetSequenceNumber(seq uint32) error {
	return d.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (d *DeleteBearerResponse) SetMessagePriority(priority uint8) error {
	return d.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (d *DeleteBearerResponse) SetTEID(teid uint32) { d.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (d *DeleteBearerResponse) SetPiggybacking(value uint8) { d.GetHeader().SetPiggybacking(value) }
