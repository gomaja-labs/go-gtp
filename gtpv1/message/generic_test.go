// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/gomaja-labs/go-gtp/gtpv1/internal/testutils"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
)

func TestGeneric(t *testing.T) {
	cases := []testutils.TestCase{
		{
			Description: "Normal",
			Structured:  message.NewGeneric(message.MsgTypeEchoRequest, 0, 0),
			Serialized: []byte{
				0x32, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			g, err := message.ParseGeneric(tc.Serialized)
			if err != nil {
				t.Fatal(err)
			}
			g.Payload = nil
			if diff := cmp.Diff(tc.Structured, g, cmp.Exporter(func(reflect.Type) bool { return true }), cmpopts.IgnoreFields(message.Generic{}, "receivedHeader")); diff != "" {
				t.Fatal(diff)
			}
			b, err := tc.Structured.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.Serialized, b); diff != "" {
				t.Fatal(diff)
			}
			if g.MarshalLen() != len(tc.Serialized) {
				t.Fatal("incorrect length")
			}
		})
	}
}
