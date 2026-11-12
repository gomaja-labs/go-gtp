// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// DeleteSessionRequest is a DeleteSessionRequest Header and its IEs above.
type DeleteSessionRequest struct {
	*Header
	Cause                             *ie.IE
	LinkedEBI                         *ie.IE
	ULI                               *ie.IE
	IndicationFlags                   *ie.IE
	PCO                               *ie.IE
	OriginatingNode                   *ie.IE
	SenderFTEIDC                      *ie.IE
	UETimeZone                        *ie.IE
	ULITimestamp                      *ie.IE
	RANNASReleaseCause                *ie.IE
	TWANIdentifier                    *ie.IE
	TWANIdentifierTimestamp           *ie.IE
	MMESGSNOverloadControlInformation *ie.IE
	SGWOverloadControlInformaion      *ie.IE
	TWANePDGOverloadControlInformaion *ie.IE
	WLANLocationInformation           *ie.IE
	WLANLocationTimeStamp             *ie.IE
	UELocalIPAddress                  *ie.IE
	UEUDPPort                         *ie.IE
	EPCO                              *ie.IE
	UETCPPort                         *ie.IE
	SecondaryRATUsageDataReport       []*ie.IE
	PrivateExtension                  *ie.IE
	AdditionalIEs                     []*ie.IE
}

// NewDeleteSessionRequest creates a new DeleteSessionRequest.
func NewDeleteSessionRequest(teid, seq uint32, ies ...*ie.IE) *DeleteSessionRequest {
	d := &DeleteSessionRequest{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeDeleteSessionRequest, teid, seq, nil,
		),
	}

	for _, i := range ies {
		d.addIE(i)
	}

	d.SetLength()
	return d
}

// Marshal serializes DeleteSessionRequest into bytes.
func (d *DeleteSessionRequest) Marshal() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes DeleteSessionRequest into bytes.
func (d *DeleteSessionRequest) MarshalTo(b []byte) error {
	if d == nil || d.Header == nil {
		return ErrNilMessage
	}
	if err := d.Header.validate(MsgTypeDeleteSessionRequest, d.MarshalLen()); err != nil {
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
	if i := d.ULI; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "ULI", iePair{ie.UserLocationInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.IndicationFlags; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "IndicationFlags", iePair{ie.Indication, 0}); err != nil {
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
	if i := d.OriginatingNode; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "OriginatingNode", iePair{ie.NodeType, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.SenderFTEIDC; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SenderFTEIDC", iePair{ie.FullyQualifiedTEID, 0}); err != nil {
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
	if i := d.ULITimestamp; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "ULITimestamp", iePair{ie.ULITimestamp, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.RANNASReleaseCause; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "RANNASReleaseCause", iePair{ie.RANNASCause, 0}); err != nil {
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
	if i := d.MMESGSNOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "MMESGSNOverloadControlInformation", iePair{ie.OverloadControlInformation, 0}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.SGWOverloadControlInformaion; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SGWOverloadControlInformaion", iePair{ie.OverloadControlInformation, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.TWANePDGOverloadControlInformaion; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "TWANePDGOverloadControlInformaion", iePair{ie.OverloadControlInformation, 2}); err != nil {
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
	if i := d.WLANLocationTimeStamp; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "WLANLocationTimeStamp", iePair{ie.TWANIdentifierTimestamp, 1}); err != nil {
			return err
		}
		offset += i.MarshalLen()
	}
	if i := d.UELocalIPAddress; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "UELocalIPAddress", iePair{ie.IPAddress, 0}); err != nil {
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
	if i := d.EPCO; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "EPCO", iePair{ie.ExtendedProtocolConfigurationOptions, 0}); err != nil {
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

// ParseDeleteSessionRequest decodes given bytes as DeleteSessionRequest.
func ParseDeleteSessionRequest(b []byte) (*DeleteSessionRequest, error) {
	d := &DeleteSessionRequest{}
	if err := d.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return d, err
		}
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary decodes given bytes as DeleteSessionRequest.
func (d *DeleteSessionRequest) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrNilMessage
	}
	*d = DeleteSessionRequest{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeDeleteSessionRequest {
		return &UnexpectedTypeError{Expected: MsgTypeDeleteSessionRequest, Actual: h.Type}
	}
	d.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, d.addIE)
}

// MarshalLen returns the serial length in int.
func (d *DeleteSessionRequest) MarshalLen() int {
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
	if i := d.ULI; i != nil {
		l += i.MarshalLen()
	}
	if i := d.IndicationFlags; i != nil {
		l += i.MarshalLen()
	}
	if i := d.PCO; i != nil {
		l += i.MarshalLen()
	}
	if i := d.OriginatingNode; i != nil {
		l += i.MarshalLen()
	}
	if i := d.SenderFTEIDC; i != nil {
		l += i.MarshalLen()
	}
	if i := d.UETimeZone; i != nil {
		l += i.MarshalLen()
	}
	if i := d.ULITimestamp; i != nil {
		l += i.MarshalLen()
	}
	if i := d.RANNASReleaseCause; i != nil {
		l += i.MarshalLen()
	}
	if i := d.TWANIdentifier; i != nil {
		l += i.MarshalLen()
	}
	if i := d.TWANIdentifierTimestamp; i != nil {
		l += i.MarshalLen()
	}
	if i := d.MMESGSNOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := d.SGWOverloadControlInformaion; i != nil {
		l += i.MarshalLen()
	}
	if i := d.TWANePDGOverloadControlInformaion; i != nil {
		l += i.MarshalLen()
	}
	if i := d.WLANLocationInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := d.WLANLocationTimeStamp; i != nil {
		l += i.MarshalLen()
	}
	if i := d.UELocalIPAddress; i != nil {
		l += i.MarshalLen()
	}
	if i := d.UEUDPPort; i != nil {
		l += i.MarshalLen()
	}
	if i := d.EPCO; i != nil {
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
func (d *DeleteSessionRequest) SetLength() {
	if d != nil && d.Header != nil {
		d.Header.setLength(d.MarshalLen() - d.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (d *DeleteSessionRequest) MessageTypeName() string {
	return "Delete Session Request"
}

// TEID returns the TEID in uint32.
func (d *DeleteSessionRequest) TEID() uint32 {
	if d == nil {
		return 0
	}
	return d.Header.teid()
}

func (d *DeleteSessionRequest) GetHeader() *Header {
	if d == nil {
		return nil
	}
	return d.Header
}

// TS 29.274 Table 7.2.9.1-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (d *DeleteSessionRequest) addIE(i *ie.IE) {
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
	case ie.UserLocationInformation:
		switch i.Instance() {
		case 0:
			if d.ULI == nil {
				d.ULI = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.Indication:
		switch i.Instance() {
		case 0:
			if d.IndicationFlags == nil {
				d.IndicationFlags = i
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
	case ie.NodeType:
		switch i.Instance() {
		case 0:
			if d.OriginatingNode == nil {
				d.OriginatingNode = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.FullyQualifiedTEID:
		switch i.Instance() {
		case 0:
			if d.SenderFTEIDC == nil {
				d.SenderFTEIDC = i
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
	case ie.ULITimestamp:
		switch i.Instance() {
		case 0:
			if d.ULITimestamp == nil {
				d.ULITimestamp = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.RANNASCause:
		switch i.Instance() {
		case 0:
			if d.RANNASReleaseCause == nil {
				d.RANNASReleaseCause = i
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
			if d.WLANLocationTimeStamp == nil {
				d.WLANLocationTimeStamp = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.OverloadControlInformation:
		switch i.Instance() {
		case 0:
			if d.MMESGSNOverloadControlInformation == nil {
				d.MMESGSNOverloadControlInformation = i
			}
		case 1:
			if d.SGWOverloadControlInformaion == nil {
				d.SGWOverloadControlInformaion = i
			}
		case 2:
			if d.TWANePDGOverloadControlInformaion == nil {
				d.TWANePDGOverloadControlInformaion = i
			}
		default:
			d.AdditionalIEs = append(d.AdditionalIEs, i)
		}
	case ie.IPAddress:
		switch i.Instance() {
		case 0:
			if d.UELocalIPAddress == nil {
				d.UELocalIPAddress = i
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
	case ie.ExtendedProtocolConfigurationOptions:
		switch i.Instance() {
		case 0:
			if d.EPCO == nil {
				d.EPCO = i
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
func (d *DeleteSessionRequest) Version() int { return d.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionRequest) MessageType() uint8 { return d.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionRequest) Sequence() uint32 { return d.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionRequest) IsPiggybacking() bool { return d.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionRequest) HasTEID() bool { return d.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionRequest) HasMessagePriority() bool { return d.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteSessionRequest) MessagePriority() uint8 { return d.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (d *DeleteSessionRequest) SetSequenceNumber(seq uint32) error {
	return d.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (d *DeleteSessionRequest) SetMessagePriority(priority uint8) error {
	return d.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (d *DeleteSessionRequest) SetTEID(teid uint32) { d.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (d *DeleteSessionRequest) SetPiggybacking(value uint8) { d.GetHeader().SetPiggybacking(value) }
