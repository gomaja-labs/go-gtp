// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import (
	"errors"
	"fmt"
)

// Error definitions.
var (
	ErrInvalidLength   = errors.New("length value is invalid")
	ErrTooShortToParse = errors.New("too short to decode as GTP")
)

// Errors returned by message encoders and header validation.
var (
	ErrTooShortToMarshal  = errors.New("destination too short to encode GTP")
	ErrNilMessage         = errors.New("nil GTP message or header")
	ErrUnsupportedVersion = errors.New("unsupported GTP version")
	ErrUnexpectedType     = errors.New("unexpected GTP message type")
	ErrInvalidTEIDFlag    = errors.New("invalid TEID flag for message type")
	ErrInvalidPiggyback   = errors.New("invalid piggybacked message")
	ErrInvalidSequence    = errors.New("sequence number exceeds 24 bits")
	ErrInvalidPriority    = errors.New("message priority exceeds four bits")
)

// UnsupportedVersionError reports the received version so a node can apply
// TS 29.274 §7.7.2 (reply for newer versions, discard older versions).
type UnsupportedVersionError struct{ Version uint8 }

func (e *UnsupportedVersionError) Error() string {
	return fmt.Sprintf("%v: %d", ErrUnsupportedVersion, e.Version)
}
func (e *UnsupportedVersionError) Unwrap() error { return ErrUnsupportedVersion }

// LengthError reports inconsistent datagram framing after a complete header was
// decoded. Its transaction fields support the response required by TS 29.274 §7.7.3.
type LengthError struct {
	MessageType    uint8
	TEID           uint32
	SequenceNumber uint32
	Declared       int
	Actual         int
}

func (e *LengthError) Error() string {
	return fmt.Sprintf("%v: type %d, declared %d bytes, received %d", ErrInvalidLength, e.MessageType, e.Declared, e.Actual)
}
func (e *LengthError) Unwrap() error { return ErrInvalidLength }

// UnexpectedTypeError reports a typed decoder used with another message type.
type UnexpectedTypeError struct{ Expected, Actual uint8 }

func (e *UnexpectedTypeError) Error() string {
	return fmt.Sprintf("%v: expected %d, received %d", ErrUnexpectedType, e.Expected, e.Actual)
}
func (e *UnexpectedTypeError) Unwrap() error { return ErrUnexpectedType }

// ErrUnexpectedIE identifies an IE assigned to the wrong named message field.
var ErrUnexpectedIE = errors.New("unexpected IE type or instance for field")

// UnexpectedIEError describes a named field containing an IE for another role.
type UnexpectedIEError struct {
	Field          string
	Type, Instance uint8
}

func (e *UnexpectedIEError) Error() string {
	return fmt.Sprintf("%v %s: type %d instance %d", ErrUnexpectedIE, e.Field, e.Type, e.Instance)
}
func (e *UnexpectedIEError) Unwrap() error { return ErrUnexpectedIE }

// IEError identifies an offending IE for the response policy in TS 29.274
// §§7.7.7–7.7.8. Framing errors reject the message; content errors accompany
// the decoded message. Cause is the GTP Cause value (67 or 69).
type IEError struct {
	Type, Instance uint8
	Cause          uint8
	Framing        bool
	Err            error
}

func (e *IEError) Error() string {
	return fmt.Sprintf("IE type %d instance %d (Cause %d): %v", e.Type, e.Instance, e.Cause, e.Err)
}
func (e *IEError) Unwrap() error { return e.Err }

func isIEContentError(err error) bool {
	var content *IEError
	return errors.As(err, &content) && !content.Framing
}
