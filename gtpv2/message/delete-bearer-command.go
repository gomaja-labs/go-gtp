// Copyright 2019-2024 go-gtp authors. All rights reservev.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

// DeleteBearerCommand is a DeleteBearerCommand Header and its IEs above.
type DeleteBearerCommand struct {
	*Header
	BearerContexts                    []*ie.IE
	ULI                               *ie.IE
	ULITimestamp                      *ie.IE
	UETimeZone                        *ie.IE
	MMESGSNOverloadControlInformation *ie.IE
	SGWOverloadControlInformation     *ie.IE
	SenderFTEIDC                      *ie.IE
	SecondaryRATDataUsageReport       []*ie.IE
	PrivateExtension                  *ie.IE
	AdditionalIEs                     []*ie.IE
}

// NewDeleteBearerCommand creates a new DeleteBearerCommand.
func NewDeleteBearerCommand(teid, seq uint32, ies ...*ie.IE) *DeleteBearerCommand {
	d := &DeleteBearerCommand{
		Header: NewHeader(
			NewHeaderFlags(2, 0, 1),
			MsgTypeDeleteBearerCommand, teid, seq, nil,
		),
	}

	for _, i := range ies {
		d.addIE(i)
	}

	d.SetLength()
	return d
}

// Marshal serializes DeleteBearerCommand into bytes.
func (d *DeleteBearerCommand) Marshal() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes DeleteBearerCommand into bytes.
func (d *DeleteBearerCommand) MarshalTo(b []byte) error {
	if d == nil || d.Header == nil {
		return ErrNilMessage
	}
	if err := d.Header.validate(MsgTypeDeleteBearerCommand, d.MarshalLen()); err != nil {
		return err
	}
	if len(b) < d.MarshalLen() {
		return ErrTooShortToMarshal
	}
	d.Header.Payload = nil
	d.Header.Payload = make([]byte, d.MarshalLen()-d.Header.MarshalLen())
	offset := 0
	for _, i := range d.BearerContexts {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "BearerContexts", iePair{ie.BearerContext, 0}); err != nil {
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
	if i := d.UETimeZone; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "UETimeZone", iePair{ie.UETimeZone, 0}); err != nil {
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
	if i := d.SGWOverloadControlInformation; i != nil {
		if err := marshalNamedIE(i, d.Payload[offset:], "SGWOverloadControlInformation", iePair{ie.OverloadControlInformation, 1}); err != nil {
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
	for _, i := range d.SecondaryRATDataUsageReport {
		if i == nil {
			continue
		}
		if err := marshalNamedIE(i, d.Payload[offset:], "SecondaryRATDataUsageReport", iePair{ie.SecondaryRATUsageDataReport, 0}); err != nil {
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

// ParseDeleteBearerCommand decodes given bytes as DeleteBearerCommand.
func ParseDeleteBearerCommand(b []byte) (*DeleteBearerCommand, error) {
	d := &DeleteBearerCommand{}
	if err := d.UnmarshalBinary(b); err != nil {
		if isIEContentError(err) {
			return d, err
		}
		return nil, err
	}
	return d, nil
}

// UnmarshalBinary decodes given bytes as DeleteBearerCommand.
func (d *DeleteBearerCommand) UnmarshalBinary(b []byte) error {
	if d == nil {
		return ErrNilMessage
	}
	*d = DeleteBearerCommand{}
	h, err := ParseHeader(b)
	if err != nil {
		return err
	}
	if h.Type != MsgTypeDeleteBearerCommand {
		return &UnexpectedTypeError{Expected: MsgTypeDeleteBearerCommand, Actual: h.Type}
	}
	d.Header = h
	decodedIEs, err := parseMessageIEs(h)
	if err != nil {
		return err
	}
	return addDecodedIEs(h, decodedIEs, d.addIE)
}

// MarshalLen returns the serial length in int.
func (d *DeleteBearerCommand) MarshalLen() int {
	if d == nil || d.Header == nil {
		return 0
	}
	l := d.Header.MarshalLen() - len(d.Header.Payload)
	for _, i := range d.BearerContexts {
		if i != nil {
			l += i.MarshalLen()
		}
	}
	if i := d.ULI; i != nil {
		l += i.MarshalLen()
	}
	if i := d.ULITimestamp; i != nil {
		l += i.MarshalLen()
	}
	if i := d.UETimeZone; i != nil {
		l += i.MarshalLen()
	}
	if i := d.MMESGSNOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := d.SGWOverloadControlInformation; i != nil {
		l += i.MarshalLen()
	}
	if i := d.SenderFTEIDC; i != nil {
		l += i.MarshalLen()
	}
	for _, i := range d.SecondaryRATDataUsageReport {
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
func (d *DeleteBearerCommand) SetLength() {
	if d != nil && d.Header != nil {
		d.Header.setLength(d.MarshalLen() - d.Header.piggybackLen() - 4)
	}
}

// MessageTypeName returns the name of protocol.
func (d *DeleteBearerCommand) MessageTypeName() string {
	return "Delete Bearer Command"
}

// TEID returns the TEID in uint32.
func (d *DeleteBearerCommand) TEID() uint32 {
	if d == nil {
		return 0
	}
	return d.Header.teid()
}

func (d *DeleteBearerCommand) GetHeader() *Header {
	if d == nil {
		return nil
	}
	return d.Header
}

// TS 29.274 Table 7.2.17.1-1: identify IEs by type and instance; §7.7.10: keep the first singleton.
func (d *DeleteBearerCommand) addIE(i *ie.IE) {
	if i == nil {
		return
	}
	if err := i.Err(); err != nil {
		d.Header.err = err
	}
	switch i.Type {
	case ie.BearerContext:
		switch i.Instance() {
		case 0:
			d.BearerContexts = append(d.BearerContexts, i)
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
	case ie.UETimeZone:
		switch i.Instance() {
		case 0:
			if d.UETimeZone == nil {
				d.UETimeZone = i
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
			if d.SGWOverloadControlInformation == nil {
				d.SGWOverloadControlInformation = i
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
	case ie.SecondaryRATUsageDataReport:
		switch i.Instance() {
		case 0:
			d.SecondaryRATDataUsageReport = append(d.SecondaryRATDataUsageReport, i)
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
func (d *DeleteBearerCommand) Version() int { return d.GetHeader().Version() }

// MessageType reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerCommand) MessageType() uint8 { return d.GetHeader().MessageType() }

// Sequence reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerCommand) Sequence() uint32 { return d.GetHeader().Sequence() }

// IsPiggybacking reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerCommand) IsPiggybacking() bool { return d.GetHeader().IsPiggybacking() }

// HasTEID reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerCommand) HasTEID() bool { return d.GetHeader().HasTEID() }

// HasMessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerCommand) HasMessagePriority() bool { return d.GetHeader().HasMessagePriority() }

// MessagePriority reads the common header, returning zero for a nil receiver.
func (d *DeleteBearerCommand) MessagePriority() uint8 { return d.GetHeader().MessagePriority() }

// SetSequenceNumber updates the common header safely on a nil receiver.
func (d *DeleteBearerCommand) SetSequenceNumber(seq uint32) error {
	return d.GetHeader().SetSequenceNumber(seq)
}

// SetMessagePriority updates the common header safely on a nil receiver.
func (d *DeleteBearerCommand) SetMessagePriority(priority uint8) error {
	return d.GetHeader().SetMessagePriority(priority)
}

// SetTEID updates the common header safely on a nil receiver.
func (d *DeleteBearerCommand) SetTEID(teid uint32) { d.GetHeader().SetTEID(teid) }

// SetPiggybacking updates the common header safely on a nil receiver.
func (d *DeleteBearerCommand) SetPiggybacking(value uint8) { d.GetHeader().SetPiggybacking(value) }
