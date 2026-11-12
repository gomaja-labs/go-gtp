// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

import "github.com/gomaja-labs/go-gtp/gtpv2/ie"

type iePair struct{ typ, instance uint8 }

func marshalNamedIE(i *ie.IE, b []byte, field string, allowed ...iePair) error {
	if err := i.Err(); err != nil {
		return err
	}
	if i.Type == ie.IntegerNumber {
		// TS 29.274 §8.118 and Table 7.3.6-1: role-specific widths.
		width := len(i.Payload)
		valid := true
		switch field {
		case "UEUsageType":
			// Table 7.3.6-1 NOTE 6: zero length means subscription unavailable.
			valid = width == 0 || width == 1
		case "RemainingRunningServiceGapTimer":
			valid = width == 4
		case "MaximumWaitTime", "DLDataPacketsSize":
			valid = width == 2
		case "DLBufferingSuggestedPacketCount":
			valid = width == 1 || width == 2
		}
		if !valid {
			return ie.ErrInvalidLength
		}
	}
	// TS 29.274 §7.7.9: a named IE role is identified by both type and instance.
	for _, pair := range allowed {
		if i.Type == pair.typ && (pair.instance == 255 || i.Instance() == pair.instance) {
			return i.MarshalTo(b)
		}
	}
	return &UnexpectedIEError{Field: field, Type: i.Type, Instance: i.Instance()}
}
