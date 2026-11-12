package ie_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

func TestIndicationSpecBits(t *testing.T) {
	flags := []struct {
		flag  ie.IndicationFlag
		octet int
		mask  byte
		has   func(*ie.IE) bool
	}{
		{ie.IndicationDAF, 5, 0x80, (*ie.IE).HasDAF},
		{ie.IndicationDTF, 5, 0x40, (*ie.IE).HasDTF},
		{ie.IndicationHI, 5, 0x20, (*ie.IE).HasHI},
		{ie.IndicationDFI, 5, 0x10, (*ie.IE).HasDFI},
		{ie.IndicationOI, 5, 0x08, (*ie.IE).HasOI},
		{ie.IndicationISRSI, 5, 0x04, (*ie.IE).HasISRSI},
		{ie.IndicationISRAI, 5, 0x02, (*ie.IE).HasISRAI},
		{ie.IndicationSGWCI, 5, 0x01, (*ie.IE).HasSGWCI},
		{ie.IndicationSQCI, 6, 0x80, (*ie.IE).HasSQCI},
		{ie.IndicationUIMSI, 6, 0x40, (*ie.IE).HasUIMSI},
		{ie.IndicationCFSI, 6, 0x20, (*ie.IE).HasCFSI},
		{ie.IndicationCRSI, 6, 0x10, (*ie.IE).HasCRSI},
		{ie.IndicationPS, 6, 0x08, (*ie.IE).HasPS},
		{ie.IndicationPT, 6, 0x04, (*ie.IE).HasPT},
		{ie.IndicationSI, 6, 0x02, (*ie.IE).HasSI},
		{ie.IndicationMSV, 6, 0x01, (*ie.IE).HasMSV},
		{ie.IndicationRetLoc, 7, 0x80, (*ie.IE).HasRetLoc},
		{ie.IndicationPBIC, 7, 0x40, (*ie.IE).HasPBIC},
		{ie.IndicationSRNI, 7, 0x20, (*ie.IE).HasSRNI},
		{ie.IndicationS6AF, 7, 0x10, (*ie.IE).HasS6AF},
		{ie.IndicationS4AF, 7, 0x08, (*ie.IE).HasS4AF},
		{ie.IndicationMBMDT, 7, 0x04, (*ie.IE).HasMBMDT},
		{ie.IndicationISRAU, 7, 0x02, (*ie.IE).HasISRAU},
		{ie.IndicationCCRSI, 7, 0x01, (*ie.IE).HasCCRSI},
		{ie.IndicationCPRAI, 8, 0x80, (*ie.IE).HasCPRAI},
		{ie.IndicationARRL, 8, 0x40, (*ie.IE).HasARRL},
		{ie.IndicationPPOFF, 8, 0x20, (*ie.IE).HasPPOFF},
		{ie.IndicationPPON, 8, 0x10, (*ie.IE).HasPPON},
		{ie.IndicationPPSI, 8, 0x08, (*ie.IE).HasPPSI},
		{ie.IndicationCSFBI, 8, 0x04, (*ie.IE).HasCSFBI},
		{ie.IndicationCLII, 8, 0x02, (*ie.IE).HasCLII},
		{ie.IndicationCPSR, 8, 0x01, (*ie.IE).HasCPSR},
		{ie.IndicationNSI, 9, 0x80, (*ie.IE).HasNSI},
		{ie.IndicationUASI, 9, 0x40, (*ie.IE).HasUASI},
		{ie.IndicationDTCI, 9, 0x20, (*ie.IE).HasDTCI},
		{ie.IndicationBDWI, 9, 0x10, (*ie.IE).HasBDWI},
		{ie.IndicationPSCI, 9, 0x08, (*ie.IE).HasPSCI},
		{ie.IndicationPCRI, 9, 0x04, (*ie.IE).HasPCRI},
		{ie.IndicationAOSI, 9, 0x02, (*ie.IE).HasAOSI},
		{ie.IndicationAOPI, 9, 0x01, (*ie.IE).HasAOPI},
		{ie.IndicationROAAI, 10, 0x80, (*ie.IE).HasROAAI},
		{ie.IndicationEPCOSI, 10, 0x40, (*ie.IE).HasEPCOSI},
		{ie.IndicationCPOPCI, 10, 0x20, (*ie.IE).HasCPOPCI},
		{ie.IndicationPMTSMI, 10, 0x10, (*ie.IE).HasPMTSMI},
		{ie.IndicationS11TF, 10, 0x08, (*ie.IE).HasS11TF},
		{ie.IndicationPNSI, 10, 0x04, (*ie.IE).HasPNSI},
		{ie.IndicationUNACCSI, 10, 0x02, (*ie.IE).HasUNACCSI},
		{ie.IndicationWPMSI, 10, 0x01, (*ie.IE).HasWPMSI},
		{ie.Indication5GSNN26, 11, 0x80, (*ie.IE).Has5GSNN26},
		{ie.IndicationREPREFI, 11, 0x40, (*ie.IE).HasREPREFI},
		{ie.Indication5GSIWKI, 11, 0x20, (*ie.IE).Has5GSIWKI},
		{ie.IndicationEEVRSI, 11, 0x10, (*ie.IE).HasEEVRSI},
		{ie.IndicationLTEMUI, 11, 0x08, (*ie.IE).HasLTEMUI},
		{ie.IndicationLTEMPI, 11, 0x04, (*ie.IE).HasLTEMPI},
		{ie.IndicationENBCRSI, 11, 0x02, (*ie.IE).HasENBCRSI},
		{ie.IndicationTSPCMI, 11, 0x01, (*ie.IE).HasTSPCMI},
		{ie.IndicationCSRMFI, 12, 0x80, (*ie.IE).HasCSRMFI},
		{ie.IndicationMTEDTN, 12, 0x40, (*ie.IE).HasMTEDTN},
		{ie.IndicationMTEDTA, 12, 0x20, (*ie.IE).HasMTEDTA},
		{ie.IndicationN5GNMI, 12, 0x10, (*ie.IE).HasN5GNMI},
		{ie.Indication5GCNRS, 12, 0x08, (*ie.IE).Has5GCNRS},
		{ie.Indication5GCNRI, 12, 0x04, (*ie.IE).Has5GCNRI},
		{ie.Indication5SRHOI, 12, 0x02, (*ie.IE).Has5SRHOI},
		{ie.IndicationETHPDN, 12, 0x01, (*ie.IE).HasETHPDN},
		{ie.IndicationNSPUSI, 13, 0x80, (*ie.IE).HasNSPUSI},
		{ie.IndicationPGWRNSI, 13, 0x40, (*ie.IE).HasPGWRNSI},
		{ie.IndicationRPPCSI, 13, 0x20, (*ie.IE).HasRPPCSI},
		{ie.IndicationPGWCHI, 13, 0x10, (*ie.IE).HasPGWCHI},
		{ie.IndicationSISSME, 13, 0x08, (*ie.IE).HasSISSME},
		{ie.IndicationNSENBI, 13, 0x04, (*ie.IE).HasNSENBI},
		{ie.IndicationIDFUPF, 13, 0x02, (*ie.IE).HasIDFUPF},
		{ie.IndicationEMCI, 13, 0x01, (*ie.IE).HasEMCI},
		{ie.IndicationUHRRSI, 14, 0x40, (*ie.IE).HasUHRRSI},
		{ie.IndicationDISREI, 14, 0x20, (*ie.IE).HasDISREI},
		{ie.IndicationPPRI, 14, 0x10, (*ie.IE).HasPPRI},
		{ie.IndicationLAPCOSI, 14, 0x08, (*ie.IE).HasLAPCOSI},
		{ie.IndicationLTEMSAI, 14, 0x04, (*ie.IE).HasLTEMSAI},
		{ie.IndicationSRTPI, 14, 0x02, (*ie.IE).HasSRTPI},
		{ie.IndicationUPIPSI, 14, 0x01, (*ie.IE).HasUPIPSI},
	}
	for j, tc := range flags {
		if int(tc.flag) != (tc.octet-5)*8+7-bitsIndex(tc.mask) {
			t.Errorf("constant %d has wrong position", j)
		}
		payload := make([]byte, 10)
		payload[tc.octet-5] = tc.mask
		i := ie.NewIndicationFromOctets(payload...)
		for k, other := range flags {
			if other.has(i) != (j == k) {
				t.Errorf("bit %d: accessor %d", j, k)
			}
		}
		built := ie.NewIndication(tc.flag, tc.flag)
		n := max(2, tc.octet-4)
		if !bytes.Equal(built.Payload, payload[:n]) {
			t.Errorf("flag %d: %x want %x", j, built.Payload, payload[:n])
		}
		for k, other := range flags {
			if other.has(built) != (j == k) {
				t.Errorf("constructed bit %d: accessor %d", j, k)
			}
		}
		for _, bad := range []*ie.IE{nil, ie.NewRecovery(0), ie.New(ie.Indication, 0, []byte{255})} {
			if tc.has(bad) {
				t.Errorf("invalid receiver for flag %d", j)
			}
		}
	}
	// Octet 14 bit 8 is spare; all flag accessors ignore it.
	payload := make([]byte, 10)
	payload[9] = 0x80
	i := ie.New(ie.Indication, 0, payload)
	for _, tc := range flags {
		if tc.has(i) {
			t.Fatal("spare bit interpreted as flag")
		}
	}
	if !bytes.Equal(ie.NewIndication().Payload, []byte{0, 0}) {
		t.Fatal("minimum length")
	}
	if !ie.NewIndication(ie.Indication5GSIWKI).Has5GSIWKI() || !ie.NewIndication(ie.IndicationPMTSMI).HasPMTSMI() || !ie.NewIndication(ie.IndicationPPON).HasPPON() || !ie.NewIndication(ie.IndicationPPOFF).HasPPOFF() {
		t.Fatal("named aliases")
	}
	for _, bad := range []ie.IndicationFlag{72, 80, 255} {
		if _, err := ie.NewIndication(bad).Marshal(); !errors.Is(err, ie.ErrMalformed) {
			t.Fatalf("undefined flag %d: %v", bad, err)
		}
	}
}
func bitsIndex(b byte) int {
	for j := 0; j < 8; j++ {
		if b == 1<<j {
			return j
		}
	}
	return -1
}
