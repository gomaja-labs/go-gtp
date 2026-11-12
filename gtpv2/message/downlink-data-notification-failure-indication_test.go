// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2"
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/internal/testutils"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

func TestDownlinkDataNotificationFailureIndication(t *testing.T) {
	cases := []testutils.TestCase{
		{
			Description: "Normal",
			Structured: message.NewDownlinkDataNotificationFailureIndication(
				testutils.TestBearerInfo.TEID, testutils.TestBearerInfo.Seq,
				ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
				ie.NewNodeType(gtpv2.NodeTypeMME),
				ie.NewIMSI("123451234567890"),
			),
			Serialized: []byte{
				// Header
				0x48, 0x46, 0x00, 0x1f, 0x11, 0x22, 0x33, 0x44, 0x00, 0x00, 0x01, 0x00,
				// Cause
				0x02, 0x00, 0x02, 0x00, 0x10, 0x00,
				// TS 29.274 Table 8.65-1: MME is Node Type 0.
				0x87, 0x00, 0x01, 0x00, 0x00,
				// IMSI
				0x01, 0x00, 0x08, 0x00, 0x21, 0x43, 0x15, 0x32, 0x54, 0x76, 0x98, 0xf0,
			},
		},
	}

	testutils.Run(t, cases, func(b []byte) (testutils.Serializable, error) {
		v, err := message.ParseDownlinkDataNotificationFailureIndication(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}
