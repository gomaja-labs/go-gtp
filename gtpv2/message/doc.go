// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package message encodes and decodes GTPv2-C messages using TS 29.274 V20.0.0.
//
// Framing errors (a truncated header or an IE crossing its container) reject
// the message. Erroneous optional or unexpected IEs with intact framing are
// retained in their named field or AdditionalIEs; their typed accessors report
// content errors. TS 29.274 §§7.7.1, 7.7.7–7.7.10 require continuing processing
// without using such content. Fixed/extendable IE suffixes are ignored on receipt.
// An invalid supplied mandatory IE returns both the message and an *IEError,
// which identifies the offending type, instance and response Cause. Conditional
// requirements depending on the interface, session or feature negotiation must
// still be checked by the caller. Echo Request ignores erroneous IEs, including
// an undelimitable tail, so its transaction can be answered (§7.1.2 NOTE).
// Encoders retain strict input validation; a tolerated received IE may need to
// be removed or replaced before a message can be transmitted.
//
// Parse and the typed ParseXxx functions consume a complete datagram. A valid
// piggybacked message is available through Message.GetHeader().Piggybacked.
// A secondary decode failure is recorded in GetHeader().PiggybackError and does
// not reject the primary message (§7.7.0). The overall datagram lengths must
// still agree (§7.7.3). Both decoding results must be inspected independently.
// MarshalLen and MarshalTo include both messages; each wire Length covers only
// its own message. Process the enclosing message before its piggybacked message
// (TS 29.274 §§5.5.1, 7.7.0).
//
// Parse distinguishes ErrTooShortToParse, UnsupportedVersionError, LengthError,
// and UnexpectedTypeError through errors.Is and errors.As. LengthError includes
// message type, TEID and sequence number for the invalid-length response policy
// in TS 29.274 §7.7.3. UnsupportedVersionError.Version enables a connection to
// reply to newer versions and discard older versions as required by §7.7.2.
// The codec does not send responses or enforce session-dependent IE presence.
//
// Named fields use exact IE type/instance pairs. Repeated fields are slices;
// singleton fields retain the first occurrence (§7.7.10). Unknown pairs remain
// available in AdditionalIEs. UnmarshalBinary resets its receiver. Encoders
// skip nil IE entries and propagate IE construction errors. SetSequenceNumber
// validates the 24-bit range; MarshalTo also validates directly assigned fields.
//
// CreateBearerResponse, UpdateBearerResponse and DeleteBearerResponse expose
// MMESGSNIdentifierOrUELocalIPAddress for IP Address instance 0. Its meaning
// depends on the interface; these roles cannot coexist in a message (Tables
// 7.2.4-1 NOTE 2, 7.2.16-1 and 7.2.10.2-1).
package message
