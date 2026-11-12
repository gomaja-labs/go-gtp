// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
	"github.com/gomaja-labs/go-gtp/gtpv1/internal/testutils"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
)

func TestSupportedExtensionHeaderNotification(t *testing.T) {
	cases := []testutils.TestCase{
		{
			Description: "Normal",
			Structured: message.NewSupportedExtensionHeaderNotification(
				testutils.TestBearerInfo.TEID, 0,
				ie.NewExtensionHeaderTypeList(
					message.ExtHeaderTypePDUSessionContainer,
					message.ExtHeaderTypeUDPPort,
				),
			),
			Serialized: []byte{
				// Header
				0x32, 0x1f, 0x00, 0x08, 0x11, 0x22, 0x33, 0x44,
				0, 0, 0, 0,
				// ExtensionHeaderTypeList
				0x8d, 0x02, 0x85, 0x40,
			},
		},
	}

	testutils.Run(t, cases, func(b []byte) (testutils.Serializable, error) {
		v, err := message.ParseSupportedExtensionHeaderNotification(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}
