// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"strconv"
)

// IndicationFlag is a bit position counted from octet 5 bit 8.
type IndicationFlag uint8

// TS 29.274 §8.12, Figure 8.12-1 and clause text: octets 5 through 14.
const (
	IndicationDAF IndicationFlag = iota
	IndicationDTF
	IndicationHI
	IndicationDFI
	IndicationOI
	IndicationISRSI
	IndicationISRAI
	IndicationSGWCI
	IndicationSQCI
	IndicationUIMSI
	IndicationCFSI
	IndicationCRSI
	IndicationPS
	IndicationPT
	IndicationSI
	IndicationMSV
	IndicationRetLoc
	IndicationPBIC
	IndicationSRNI
	IndicationS6AF
	IndicationS4AF
	IndicationMBMDT
	IndicationISRAU
	IndicationCCRSI
	IndicationCPRAI
	IndicationARRL
	IndicationPPOFF
	IndicationPPON
	IndicationPPSI
	IndicationCSFBI
	IndicationCLII
	IndicationCPSR
	IndicationNSI
	IndicationUASI
	IndicationDTCI
	IndicationBDWI
	IndicationPSCI
	IndicationPCRI
	IndicationAOSI
	IndicationAOPI
	IndicationROAAI
	IndicationEPCOSI
	IndicationCPOPCI
	IndicationPMTSMI
	IndicationS11TF
	IndicationPNSI
	IndicationUNACCSI
	IndicationWPMSI
	Indication5GSNN26
	IndicationREPREFI
	Indication5GSIWKI
	IndicationEEVRSI
	IndicationLTEMUI
	IndicationLTEMPI
	IndicationENBCRSI
	IndicationTSPCMI
	IndicationCSRMFI
	IndicationMTEDTN
	IndicationMTEDTA
	IndicationN5GNMI
	Indication5GCNRS
	Indication5GCNRI
	Indication5SRHOI
	IndicationETHPDN
	IndicationNSPUSI
	IndicationPGWRNSI
	IndicationRPPCSI
	IndicationPGWCHI
	IndicationSISSME
	IndicationNSENBI
	IndicationIDFUPF
	IndicationEMCI
	indicationOctet14Spare // TS 29.274 §8.12: octet 14 bit 8.
	IndicationUHRRSI
	IndicationDISREI
	IndicationPPRI
	IndicationLAPCOSI
	IndicationLTEMSAI
	IndicationSRTPI
	IndicationUPIPSI
)

// NewIndication creates the shortest Indication containing all requested flags.
func NewIndication(flags ...IndicationFlag) *IE {
	// TS 29.274 Table 8.1-1: at least two fixed octets.
	size := 2
	for _, flag := range flags {
		if flag > IndicationUPIPSI || flag == indicationOctet14Spare {
			return newErrorIE(Indication, ErrMalformed)
		}
		if n := int(flag)/8 + 1; n > size {
			size = n
		}
	}
	payload := make([]byte, size)
	for _, flag := range flags {
		payload[int(flag)/8] |= 0x80 >> (flag % 8)
	}
	return New(Indication, 0, payload)
}

// NewIndicationFromBitSequence creates a new Indication IE from string-formatted
// sequence of bits. The sequence should look the same as Wireshark appearance.
// The input is to be like "101000010000100000010101000100001000100010000001010000001010000000000001".
func NewIndicationFromBitSequence(bitSequence string) *IE {
	if len(bitSequence) < 16 || len(bitSequence)%8 != 0 {
		return newErrorIE(Indication, ErrInvalidLength)
	}

	ie := New(Indication, 0x00, make([]byte, len(bitSequence)/8))
	for i, r := range bitSequence {
		bit, err := strconv.Atoi(string(r))
		if err != nil {
			return newErrorIE(Indication, err)
		}
		if bit > 1 {
			return newErrorIE(Indication, ErrMalformed)
		}

		// index:
		//   0 =< i < 8 : set bits in first octet
		//   8 =< i < 16: set bits in second octet ...
		// bit shift:
		//   0 =< i < 8 : shift i to left
		//   8 =< i < 16: shift i - 8 to left ...
		ie.Payload[i/8] |= uint8(bit << uint8(8*(i/8+1)-i-1))
	}

	if len(ie.Payload) >= 10 {
		ie.Payload[9] &= 0x7f
	}
	return ie
}

// NewIndicationFromOctets creates a new IndicationFromOctets IE from the set of octets.
func NewIndicationFromOctets(octs ...uint8) *IE {
	// TS 29.274 §8.12: octets 5 and 6 are mandatory; later flags extend the IE.
	if len(octs) < 2 {
		return newErrorIE(Indication, ErrInvalidLength)
	}
	payload := append([]byte(nil), octs...)
	if len(payload) >= 10 {
		payload[9] &= 0x7f
	}
	return New(Indication, 0, payload)
}

// Indication returns Indication flags in []byte if the type of IE matches.
func (i *IE) Indication() ([]byte, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	if i.Type != Indication {
		return nil, &InvalidTypeError{i.Type}
	}

	if len(i.Payload) < 2 {
		return nil, ErrInvalidLength
	}
	return i.Payload, nil
}

func (i *IE) hasIndication(flag IndicationFlag) bool {
	v, err := i.Indication()
	return err == nil && int(flag)/8 < len(v) && v[int(flag)/8]&(0x80>>(flag%8)) != 0
}

// HasDAF reports whether the DAF flag is set.
func (i *IE) HasDAF() bool { return i.hasIndication(IndicationDAF) }

// HasDTF reports whether the DTF flag is set.
func (i *IE) HasDTF() bool { return i.hasIndication(IndicationDTF) }

// HasHI reports whether the HI flag is set.
func (i *IE) HasHI() bool { return i.hasIndication(IndicationHI) }

// HasDFI reports whether the DFI flag is set.
func (i *IE) HasDFI() bool { return i.hasIndication(IndicationDFI) }

// HasOI reports whether the OI flag is set.
func (i *IE) HasOI() bool { return i.hasIndication(IndicationOI) }

// HasISRSI reports whether the ISRSI flag is set.
func (i *IE) HasISRSI() bool { return i.hasIndication(IndicationISRSI) }

// HasISRAI reports whether the ISRAI flag is set.
func (i *IE) HasISRAI() bool { return i.hasIndication(IndicationISRAI) }

// HasSGWCI reports whether the SGWCI flag is set.
func (i *IE) HasSGWCI() bool { return i.hasIndication(IndicationSGWCI) }

// HasSQCI reports whether the SQCI flag is set.
func (i *IE) HasSQCI() bool { return i.hasIndication(IndicationSQCI) }

// HasUIMSI reports whether the UIMSI flag is set.
func (i *IE) HasUIMSI() bool { return i.hasIndication(IndicationUIMSI) }

// HasCFSI reports whether the CFSI flag is set.
func (i *IE) HasCFSI() bool { return i.hasIndication(IndicationCFSI) }

// HasCRSI reports whether the CRSI flag is set.
func (i *IE) HasCRSI() bool { return i.hasIndication(IndicationCRSI) }

// HasPS reports whether the PS flag is set.
func (i *IE) HasPS() bool { return i.hasIndication(IndicationPS) }

// HasPT reports whether the PT flag is set.
func (i *IE) HasPT() bool { return i.hasIndication(IndicationPT) }

// HasSI reports whether the SI flag is set.
func (i *IE) HasSI() bool { return i.hasIndication(IndicationSI) }

// HasMSV reports whether the MSV flag is set.
func (i *IE) HasMSV() bool { return i.hasIndication(IndicationMSV) }

// HasRetLoc reports whether the RetLoc flag is set.
func (i *IE) HasRetLoc() bool { return i.hasIndication(IndicationRetLoc) }

// HasPBIC reports whether the PBIC flag is set.
func (i *IE) HasPBIC() bool { return i.hasIndication(IndicationPBIC) }

// HasSRNI reports whether the SRNI flag is set.
func (i *IE) HasSRNI() bool { return i.hasIndication(IndicationSRNI) }

// HasS6AF reports whether the S6AF flag is set.
func (i *IE) HasS6AF() bool { return i.hasIndication(IndicationS6AF) }

// HasS4AF reports whether the S4AF flag is set.
func (i *IE) HasS4AF() bool { return i.hasIndication(IndicationS4AF) }

// HasMBMDT reports whether the MBMDT flag is set.
func (i *IE) HasMBMDT() bool { return i.hasIndication(IndicationMBMDT) }

// HasISRAU reports whether the ISRAU flag is set.
func (i *IE) HasISRAU() bool { return i.hasIndication(IndicationISRAU) }

// HasCCRSI reports whether the CCRSI flag is set.
func (i *IE) HasCCRSI() bool { return i.hasIndication(IndicationCCRSI) }

// HasCPRAI reports whether the CPRAI flag is set.
func (i *IE) HasCPRAI() bool { return i.hasIndication(IndicationCPRAI) }

// HasARRL reports whether the ARRL flag is set.
func (i *IE) HasARRL() bool { return i.hasIndication(IndicationARRL) }

// HasPPON reports whether the PPON / PPEI flag is set.
func (i *IE) HasPPON() bool { return i.hasIndication(IndicationPPON) }

// HasPPSI reports whether the PPSI flag is set.
func (i *IE) HasPPSI() bool { return i.hasIndication(IndicationPPSI) }

// HasCSFBI reports whether the CSFBI flag is set.
func (i *IE) HasCSFBI() bool { return i.hasIndication(IndicationCSFBI) }

// HasCLII reports whether the CLII flag is set.
func (i *IE) HasCLII() bool { return i.hasIndication(IndicationCLII) }

// HasCPSR reports whether the CPSR flag is set.
func (i *IE) HasCPSR() bool { return i.hasIndication(IndicationCPSR) }

// HasNSI reports whether the NSI flag is set.
func (i *IE) HasNSI() bool { return i.hasIndication(IndicationNSI) }

// HasUASI reports whether the UASI flag is set.
func (i *IE) HasUASI() bool { return i.hasIndication(IndicationUASI) }

// HasDTCI reports whether the DTCI flag is set.
func (i *IE) HasDTCI() bool { return i.hasIndication(IndicationDTCI) }

// HasBDWI reports whether the BDWI flag is set.
func (i *IE) HasBDWI() bool { return i.hasIndication(IndicationBDWI) }

// HasPSCI reports whether the PSCI flag is set.
func (i *IE) HasPSCI() bool { return i.hasIndication(IndicationPSCI) }

// HasPCRI reports whether the PCRI flag is set.
func (i *IE) HasPCRI() bool { return i.hasIndication(IndicationPCRI) }

// HasAOSI reports whether the AOSI flag is set.
func (i *IE) HasAOSI() bool { return i.hasIndication(IndicationAOSI) }

// HasAOPI reports whether the AOPI flag is set.
func (i *IE) HasAOPI() bool { return i.hasIndication(IndicationAOPI) }

// HasROAAI reports whether the ROAAI flag is set.
func (i *IE) HasROAAI() bool { return i.hasIndication(IndicationROAAI) }

// HasEPCOSI reports whether the EPCOSI flag is set.
func (i *IE) HasEPCOSI() bool { return i.hasIndication(IndicationEPCOSI) }

// HasCPOPCI reports whether the CPOPCI flag is set.
func (i *IE) HasCPOPCI() bool { return i.hasIndication(IndicationCPOPCI) }

// HasPMTSMI reports whether the PMTSMI flag is set.
func (i *IE) HasPMTSMI() bool { return i.hasIndication(IndicationPMTSMI) }

// HasS11TF reports whether the S11TF flag is set.
func (i *IE) HasS11TF() bool { return i.hasIndication(IndicationS11TF) }

// HasPNSI reports whether the PNSI flag is set.
func (i *IE) HasPNSI() bool { return i.hasIndication(IndicationPNSI) }

// HasUNACCSI reports whether the UNACCSI flag is set.
func (i *IE) HasUNACCSI() bool { return i.hasIndication(IndicationUNACCSI) }

// HasWPMSI reports whether the WPMSI flag is set.
func (i *IE) HasWPMSI() bool { return i.hasIndication(IndicationWPMSI) }

// Has5GSNN26 reports whether the 5GSNN26 flag is set.
func (i *IE) Has5GSNN26() bool { return i.hasIndication(Indication5GSNN26) }

// HasREPREFI reports whether the REPREFI flag is set.
func (i *IE) HasREPREFI() bool { return i.hasIndication(IndicationREPREFI) }

// HasEEVRSI reports whether the EEVRSI flag is set.
func (i *IE) HasEEVRSI() bool { return i.hasIndication(IndicationEEVRSI) }

// HasLTEMUI reports whether the LTEMUI flag is set.
func (i *IE) HasLTEMUI() bool { return i.hasIndication(IndicationLTEMUI) }

// HasLTEMPI reports whether the LTEMPI flag is set.
func (i *IE) HasLTEMPI() bool { return i.hasIndication(IndicationLTEMPI) }

// HasENBCRSI reports whether the ENBCRSI flag is set.
func (i *IE) HasENBCRSI() bool { return i.hasIndication(IndicationENBCRSI) }

// HasTSPCMI reports whether the TSPCMI flag is set.
func (i *IE) HasTSPCMI() bool { return i.hasIndication(IndicationTSPCMI) }

// HasCSRMFI reports whether the CSRMFI flag is set.
func (i *IE) HasCSRMFI() bool { return i.hasIndication(IndicationCSRMFI) }

// HasMTEDTN reports whether the MTEDTN flag is set.
func (i *IE) HasMTEDTN() bool { return i.hasIndication(IndicationMTEDTN) }

// HasMTEDTA reports whether the MTEDTA flag is set.
func (i *IE) HasMTEDTA() bool { return i.hasIndication(IndicationMTEDTA) }

// HasN5GNMI reports whether the N5GNMI flag is set.
func (i *IE) HasN5GNMI() bool { return i.hasIndication(IndicationN5GNMI) }

// Has5GCNRS reports whether the 5GCNRS flag is set.
func (i *IE) Has5GCNRS() bool { return i.hasIndication(Indication5GCNRS) }

// Has5GCNRI reports whether the 5GCNRI flag is set.
func (i *IE) Has5GCNRI() bool { return i.hasIndication(Indication5GCNRI) }

// Has5SRHOI reports whether the 5SRHOI flag is set.
func (i *IE) Has5SRHOI() bool { return i.hasIndication(Indication5SRHOI) }

// HasETHPDN reports whether the ETHPDN flag is set.
func (i *IE) HasETHPDN() bool { return i.hasIndication(IndicationETHPDN) }

// HasNSPUSI reports whether the NSPUSI flag is set.
func (i *IE) HasNSPUSI() bool { return i.hasIndication(IndicationNSPUSI) }

// HasPGWRNSI reports whether the PGWRNSI flag is set.
func (i *IE) HasPGWRNSI() bool { return i.hasIndication(IndicationPGWRNSI) }

// HasRPPCSI reports whether the RPPCSI flag is set.
func (i *IE) HasRPPCSI() bool { return i.hasIndication(IndicationRPPCSI) }

// HasPGWCHI reports whether the PGWCHI flag is set.
func (i *IE) HasPGWCHI() bool { return i.hasIndication(IndicationPGWCHI) }

// HasSISSME reports whether the SISSME flag is set.
func (i *IE) HasSISSME() bool { return i.hasIndication(IndicationSISSME) }

// HasNSENBI reports whether the NSENBI flag is set.
func (i *IE) HasNSENBI() bool { return i.hasIndication(IndicationNSENBI) }

// HasIDFUPF reports whether the IDFUPF flag is set.
func (i *IE) HasIDFUPF() bool { return i.hasIndication(IndicationIDFUPF) }

// HasEMCI reports whether the EMCI flag is set.
func (i *IE) HasEMCI() bool { return i.hasIndication(IndicationEMCI) }

// HasUHRRSI reports whether the UHRRSI flag is set.
func (i *IE) HasUHRRSI() bool { return i.hasIndication(IndicationUHRRSI) }

// HasDISREI reports whether the DISREI flag is set.
func (i *IE) HasDISREI() bool { return i.hasIndication(IndicationDISREI) }

// HasPPRI reports whether the PPRI flag is set.
func (i *IE) HasPPRI() bool { return i.hasIndication(IndicationPPRI) }

// HasLAPCOSI reports whether the LAPCOSI flag is set.
func (i *IE) HasLAPCOSI() bool { return i.hasIndication(IndicationLAPCOSI) }

// HasLTEMSAI reports whether the LTEMSAI flag is set.
func (i *IE) HasLTEMSAI() bool { return i.hasIndication(IndicationLTEMSAI) }

// HasSRTPI reports whether the SRTPI flag is set.
func (i *IE) HasSRTPI() bool { return i.hasIndication(IndicationSRTPI) }

// HasUPIPSI reports whether the UPIPSI flag is set.
func (i *IE) HasUPIPSI() bool { return i.hasIndication(IndicationUPIPSI) }

// HasPPOFF reports whether the PDN Pause Off flag is set.
func (i *IE) HasPPOFF() bool { return i.hasIndication(IndicationPPOFF) }

// Has5GSIWKI reports whether the 5GS Interworking flag is set.
func (i *IE) Has5GSIWKI() bool { return i.hasIndication(Indication5GSIWKI) }
