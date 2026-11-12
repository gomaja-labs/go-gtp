// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"errors"
	"fmt"
)

// Error definitions.
var (
	ErrTooShortToMarshal = errors.New("too short to marshal IE")
	ErrTooShortToParse   = errors.New("too short to parse IE")
	ErrInvalidLength     = errors.New("invalid IE length")

	ErrInvalidType     = errors.New("invalid type")
	ErrIENotFound      = errors.New("IE not found")
	ErrIEValueNotFound = errors.New("could not find the specified value in an IE")

	ErrMalformed = errors.New("malformed IE")
)

// InvalidTypeError indicates the type of IE is invalid.
type InvalidTypeError struct {
	Type uint8
}

// Error returns message with the invalid type given.
func (e *InvalidTypeError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("got invalid type: %v", e.Type)
}

// Unwrap exposes the invalid-type category.
func (e *InvalidTypeError) Unwrap() error { return ErrInvalidType }
