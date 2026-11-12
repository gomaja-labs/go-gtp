// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/internal/testutils"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

func TestGeneric(t *testing.T) {
	cases := []testutils.TestCase{
		{
			Description: "Normal",
			Structured: message.NewGeneric(
				message.MsgTypeEchoRequest,
				testutils.TestBearerInfo.TEID, testutils.TestBearerInfo.Seq,
				ie.NewRecovery(0x80),
			),
			Serialized: []byte{
				0x40, 0x01, 0x00, 0x09, 0x00, 0x00, 0x01, 0x00, 0x03, 0x00, 0x01, 0x00,
				0x80,
			},
		}, {
			Description: "No-TEID",
			Structured: message.NewGeneric(
				message.MsgTypeEchoRequest,
				testutils.TestBearerInfo.TEID, testutils.TestBearerInfo.Seq,
				ie.NewRecovery(0x80),
			),
			Serialized: []byte{
				0x40, 0x01, 0x00, 0x09, 0x00, 0x00, 0x01, 0x00, 0x03, 0x00, 0x01, 0x00,
				0x80,
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
			if diff := cmp.Diff(tc.Structured, g, cmp.Exporter(func(reflect.Type) bool { return true }), cmpopts.IgnoreFields(message.Header{}, "replay")); diff != "" {
				t.Fatal(diff)
			}
			for _, msg := range []testutils.Serializable{tc.Structured, g} {
				wire, err := msg.Marshal()
				if err != nil {
					t.Fatal(err)
				}
				if diff := cmp.Diff(tc.Serialized, wire); diff != "" {
					t.Fatal(diff)
				}
				if msg.MarshalLen() != len(tc.Serialized) {
					t.Fatalf("MarshalLen = %d, want %d", msg.MarshalLen(), len(tc.Serialized))
				}
			}
		})
	}
}
