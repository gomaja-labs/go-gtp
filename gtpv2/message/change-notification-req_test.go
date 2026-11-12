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

func TestChangeNotificationRequest(t *testing.T) {
	cases := []testutils.TestCase{
		{
			Description: "Normal",
			Structured: message.NewChangeNotificationRequest(
				testutils.TestBearerInfo.TEID, testutils.TestBearerInfo.Seq,
				ie.NewIMSI("123451234567890"),
				ie.NewMobileEquipmentIdentity("123450123456789"),
				ie.NewIndication(ie.IndicationDAF, ie.IndicationHI, ie.IndicationSGWCI, ie.IndicationPS, ie.IndicationS6AF, ie.IndicationMBMDT, ie.IndicationCCRSI, ie.IndicationPPON, ie.IndicationNSI, ie.IndicationPSCI, ie.IndicationROAAI, ie.IndicationWPMSI, ie.IndicationREPREFI, ie.IndicationCSRMFI, ie.IndicationMTEDTA, ie.IndicationEMCI),
				ie.NewRATType(gtpv2.RATTypeEUTRAN),
				ie.NewUserLocationInformation(
					ie.NewCGI("123", "45", 0x1111, 0x2222),
					ie.NewSAI("123", "45", 0x1111, 0x3333),
					ie.NewRAI("123", "45", 0x1111, 0x44),
					ie.NewTAI("123", "45", 0x5555),
					ie.NewECGI("123", "45", 0x06666666),
					ie.NewLAI("123", "45", 0x1111),
					ie.NewMENBI("123", "45", 0x011111),
					nil,
				),
				ie.NewUserCSGInformation("123", "45", 0x00ffffff, gtpv2.AccessModeHybrid, 0, gtpv2.UCICSG),
				ie.NewIPAddress("1.1.1.1"),
				ie.NewEPSBearerID(0x05),
				ie.NewPrivateExtension(10415, []byte{0xde, 0xad, 0xbe, 0xef}),
			),
			Serialized: []byte{
				// Header
				0x48, 0x26, 0x00, 0x86, 0x11, 0x22, 0x33, 0x44, 0x00, 0x00, 0x01, 0x00,
				// IMSI
				0x01, 0x00, 0x08, 0x00, 0x21, 0x43, 0x15, 0x32, 0x54, 0x76, 0x98, 0xf0,
				// MAI
				0x4b, 0x00, 0x08, 0x00, 0x21, 0x43, 0x05, 0x21, 0x43, 0x65, 0x87, 0xf9,
				// IndicationFlags
				0x4d, 0x00, 0x09, 0x00, 0xa1, 0x08, 0x15, 0x10, 0x88, 0x81, 0x40, 0xa0, 0x01,
				// RATType
				0x52, 0x00, 0x01, 0x00, 0x06,
				// ULI
				// --------------------------------------
				0x56, 0x00, 0x2d, 0x00,
				// TS 29.274 Figure 8.21-1: Macro and Extended Macro IDs are exclusive.
				0x7f,
				// CGI
				0x21, 0xf3, 0x54, 0x11, 0x11, 0x22, 0x22,
				// SAI
				0x21, 0xf3, 0x54, 0x11, 0x11, 0x33, 0x33,
				// TS 29.274 Figure 8.21.3-1: one-octet RAC followed by 0xff.
				0x21, 0xf3, 0x54, 0x11, 0x11, 0x44, 0xff,
				// TAI
				0x21, 0xf3, 0x54, 0x55, 0x55,
				// TS 29.274 Figure 8.21.5-1: 28-bit ECI.
				0x21, 0xf3, 0x54, 0x06, 0x66, 0x66, 0x66,
				// LAI
				0x21, 0xf3, 0x54, 0x11, 0x11,
				// TS 29.274 Figure 8.21.7-1: 20-bit Macro eNB ID.
				0x21, 0xf3, 0x54, 0x01, 0x11, 0x11,
				// --------------------------------------
				// UCI
				0x91, 0x00, 0x08, 0x00, 0x21, 0xf3, 0x54, 0x00, 0xff, 0xff, 0xff, 0x41,
				// PGWS5S8IPAddressForControlPlane
				0x4a, 0x00, 0x04, 0x00, 0x01, 0x01, 0x01, 0x01,
				// LinkedEBI
				0x49, 0x00, 0x01, 0x00, 0x05,
				// PrivateExtension
				0xff, 0x00, 0x06, 0x00, 0x28, 0xaf, 0xde, 0xad, 0xbe, 0xef,
			},
		},
	}

	testutils.Run(t, cases, func(b []byte) (testutils.Serializable, error) {
		v, err := message.ParseChangeNotificationRequest(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}
