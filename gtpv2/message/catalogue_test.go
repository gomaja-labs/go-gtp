package message

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

type catalogueField struct {
	typ, instance uint8
	name          string
}

// These type/instance roles follow the cited top-level message tables in
// TS 29.274 V20.0.0. Keep the fixture separate from constructor dispatch.
var catalogue = map[string][]catalogueField{
	// TS 29.274 Table 7.3.14-1.
	"ChangeNotificationRequest": {
		{ie.IMSI, 0, "IMSI"},
		{ie.MobileEquipmentIdentity, 0, "MEI"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.RATType, 0, "RATType"},
		{ie.UserLocationInformation, 0, "ULI"},
		{ie.UserCSGInformation, 0, "UCI"},
		{ie.IPAddress, 0, "PGWS5S8IPAddressForControlPlane"},
		{ie.EPSBearerID, 0, "LinkedEBI"},
		{ie.PresenceReportingAreaInformation, 0, "PresenceReportingAreaInformation"},
		{ie.Counter, 0, "MOExceptionDataCounter"},
		{ie.SecondaryRATUsageDataReport, 0, "SecondaryRATUsageDataReport"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.3.15-1.
	"ChangeNotificationResponse": {
		{ie.IMSI, 0, "IMSI"},
		{ie.MobileEquipmentIdentity, 0, "MEI"},
		{ie.Cause, 0, "Cause"},
		{ie.ChangeReportingAction, 0, "ChangeReportingAction"},
		{ie.CSGInformationReportingAction, 0, "CSGInformationReportingAction"},
		{ie.PresenceReportingAreaAction, 0, "PresenceReportingAreaAction"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.3.7-1.
	"ContextAcknowledge": {
		{ie.Cause, 0, "Cause"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.FullyQualifiedTEID, 0, "ForwardingFTEID"},
		{ie.BearerContext, 0, "BearerContexts"},
		{ie.NodeNumber, 0, "SGSNNumber"},
		{ie.NodeNumber, 1, "MMENumberForMTSMS"},
		{ie.NodeIdentifier, 0, "SGSNIdentifierForMTSMS"},
		{ie.NodeIdentifier, 1, "MMEIdentifierForMTSMS"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.3.5-1.
	"ContextRequest": {
		{ie.IMSI, 0, "IMSI"},
		{ie.GUTI, 0, "GUTI"},
		{ie.UserLocationInformation, 0, "RAI"},
		{ie.PacketTMSI, 0, "PTMSI"},
		{ie.PTMSISignature, 0, "PTMSISignature"},
		{ie.CompleteRequestMessage, 0, "CompleteTAURequestMessage"},
		{ie.FullyQualifiedTEID, 0, "AddressAndTEIDForCPlane"},
		{ie.PortNumber, 0, "UDPSourcePortNumber"},
		{ie.RATType, 0, "RATType"},
		{ie.Indication, 0, "Indication"},
		{ie.HopCounter, 0, "HopCounter"},
		{ie.ServingNetwork, 0, "TargetPLMNID"},
		{ie.LocalDistinguishedName, 0, "MMESGSNLDN"},
		{ie.FullyQualifiedDomainName, 0, "SGSNNodeName"},
		{ie.FullyQualifiedDomainName, 1, "MMENodeName"},
		{ie.NodeNumber, 0, "SGSNNumber"},
		{ie.NodeIdentifier, 0, "SGSNIdentifier"},
		{ie.NodeIdentifier, 1, "MMEIdentifier"},
		{ie.CIoTOptimizationsSupportIndication, 0, "CIoTOptimizationsSupportIndication"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.3.6-1.
	"ContextResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.IMSI, 0, "IMSI"},
		{ie.MMContextEPSSecurityContextQuadrupletsAndQuintuplets, 0, "UEMMContext"},
		{ie.MMContextGSMKeyAndTriplets, 0, "UEMMContext"},
		{ie.MMContextGSMKeyUsedCipherAndQuintuplets, 0, "UEMMContext"},
		{ie.MMContextUMTSKeyAndQuintuplets, 0, "UEMMContext"},
		{ie.MMContextUMTSKeyQuadrupletsAndQuintuplets, 0, "UEMMContext"},
		{ie.MMContextUMTSKeyUsedCipherAndQuintuplets, 0, "UEMMContext"},
		{ie.PDNConnection, 0, "UEPDNConnections"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEID"},
		{ie.FullyQualifiedTEID, 1, "SGWS11S4FTEID"},
		{ie.FullyQualifiedDomainName, 0, "SGWNodeName"},
		{ie.FullyQualifiedDomainName, 1, "SGSNNodeName"},
		{ie.FullyQualifiedDomainName, 2, "MMENodeName"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.TraceInformation, 0, "TraceInformation"},
		{ie.IPAddress, 0, "S101IPAddress"},
		{ie.IPAddress, 1, "S102IPAddress"},
		{ie.RFSPIndex, 0, "SubscribedRFSPIndex"},
		{ie.RFSPIndex, 1, "RFSPIndexInUse"},
		{ie.UETimeZone, 0, "UETimeZone"},
		{ie.LocalDistinguishedName, 0, "MMESGSNLDN"},
		{ie.MDTConfiguration, 0, "MDTConfiguration"},
		{ie.UserCSGInformation, 0, "UCI"},
		{ie.MonitoringEventInformation, 0, "MonitoringEventInformation"},
		{ie.MonitoringEventExtensionInformation, 0, "MonitoringEventExtensionInformation"},
		{ie.IntegerNumber, 0, "UEUsageType"},
		{ie.IntegerNumber, 1, "RemainingRunningServiceGapTimer"},
		{ie.SCEFPDNConnection, 0, "SCEFPDNConnection"},
		{ie.RATType, 0, "RATType"},
		{ie.ServingPLMNRateControl, 0, "ServingPLMNRateControl"},
		{ie.Counter, 0, "MOExceptionDataCounter"},
		{ie.ExtendedTraceInformation, 0, "ExtendedTraceInformation"},
		{ie.AdditionalRRMPolicyIndex, 0, "SubscribedAdditionalRRMPolicyIndex"},
		{ie.AdditionalRRMPolicyIndex, 1, "AdditionalRRMPolicyIndexInUse"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.3-1.
	"CreateBearerRequest": {
		{ie.ProcedureTransactionID, 0, "PTI"},
		{ie.EPSBearerID, 0, "LinkedEBI"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.BearerContext, 0, "BearerContexts"},
		{ie.FullyQualifiedCSID, 0, "PGWFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.ChangeReportingAction, 0, "ChangeReportingAction"},
		{ie.CSGInformationReportingAction, 0, "CSGInformationReportingAction"},
		{ie.HeNBInformationReporting, 0, "HeNBInformationReporting"},
		{ie.PresenceReportingAreaAction, 0, "PresenceReportingAreaAction"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.LoadControlInformation, 0, "PGWNodeLoadControlInformation"},
		{ie.LoadControlInformation, 1, "PGWAPNLoadControlInformation"},
		{ie.LoadControlInformation, 2, "SGWNodeLoadControlInformation"},
		{ie.OverloadControlInformation, 0, "PGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.FContainer, 0, "NBIFOMContainer"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.4-1.
	"CreateBearerResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.BearerContext, 0, "BearerContexts"},
		{ie.Recovery, 0, "Recovery"},
		{ie.FullyQualifiedCSID, 0, "MMEFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.FullyQualifiedCSID, 2, "EPDGFQCSID"},
		{ie.FullyQualifiedCSID, 3, "TWANFQCSID"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.UETimeZone, 0, "UETimeZone"},
		{ie.UserLocationInformation, 0, "ULI"},
		{ie.TWANIdentifier, 0, "TWANIdentifier"},
		{ie.TWANIdentifier, 1, "WLANLocationInformation"},
		{ie.TWANIdentifierTimestamp, 1, "WLANLocationTimestamp"},
		{ie.OverloadControlInformation, 0, "MMEOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 2, "TWANePDGOverloadControlInformation"},
		{ie.PresenceReportingAreaInformation, 0, "PresenceReportingAreaInformation"},
		{ie.IPAddress, 0, "MMESGSNIdentifierOrUELocalIPAddress"},
		{ie.PortNumber, 0, "UEUDPPort"},
		{ie.PortNumber, 1, "UETCPPort"},
		{ie.FContainer, 0, "NBIFOMContainer"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.1-1.
	"CreateSessionRequest": {
		{ie.IMSI, 0, "IMSI"},
		{ie.MSISDN, 0, "MSISDN"},
		{ie.MobileEquipmentIdentity, 0, "MEI"},
		{ie.UserLocationInformation, 0, "ULI"},
		{ie.UserLocationInformation, 1, "ULIForSGW"},
		{ie.ServingNetwork, 0, "ServingNetwork"},
		{ie.RATType, 0, "RATType"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEIDC"},
		{ie.FullyQualifiedTEID, 1, "PGWS5S8FTEIDC"},
		{ie.AccessPointName, 0, "APN"},
		{ie.SelectionMode, 0, "SelectionMode"},
		{ie.PDNType, 0, "PDNType"},
		{ie.PDNAddressAllocation, 0, "PAA"},
		{ie.APNRestriction, 0, "APNRestriction"},
		{ie.AggregateMaximumBitRate, 0, "AMBR"},
		{ie.EPSBearerID, 0, "LinkedEBI"},
		{ie.TrustedWLANModeIndication, 0, "TWMI"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.BearerContext, 0, "BearerContextsToBeCreated"},
		{ie.BearerContext, 1, "BearerContextsToBeRemoved"},
		{ie.FullyQualifiedCSID, 0, "MMEFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.FullyQualifiedCSID, 2, "EPDGFQCSID"},
		{ie.FullyQualifiedCSID, 3, "TWANFQCSID"},
		{ie.TraceInformation, 0, "TraceInformation"},
		{ie.Recovery, 0, "Recovery"},
		{ie.UETimeZone, 0, "UETimeZone"},
		{ie.UserCSGInformation, 0, "UCI"},
		{ie.ChargingCharacteristics, 0, "ChargingCharacteristics"},
		{ie.LocalDistinguishedName, 0, "MMESGSNLDN"},
		{ie.LocalDistinguishedName, 1, "SGWLDN"},
		{ie.LocalDistinguishedName, 2, "EPDGLDN"},
		{ie.LocalDistinguishedName, 3, "TWANLDN"},
		{ie.SignallingPriorityIndication, 0, "SignallingPriorityIndication"},
		{ie.IPAddress, 0, "UELocalIPAddress"},
		{ie.IPAddress, 1, "HeNBLocalIPAddress"},
		{ie.IPAddress, 2, "MMESGSNIdentifier"},
		{ie.IPAddress, 3, "EPDGIPAddress"},
		{ie.PortNumber, 0, "UEUDPPort"},
		{ie.PortNumber, 1, "HeNBUDPPort"},
		{ie.PortNumber, 2, "UETCPPort"},
		{ie.AdditionalProtocolConfigurationOptions, 0, "APCO"},
		{ie.TWANIdentifier, 0, "TWANIdentifier"},
		{ie.TWANIdentifier, 1, "WLANLocationInformation"},
		{ie.CNOperatorSelectionEntity, 0, "CNOperatorSelectionEntity"},
		{ie.PresenceReportingAreaInformation, 0, "PresenceReportingAreaInformation"},
		{ie.OverloadControlInformation, 0, "MMESGSNOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 2, "TWANePDGOverloadControlInformation"},
		{ie.MillisecondTimeStamp, 0, "OriginationTimeStamp"},
		{ie.IntegerNumber, 0, "MaximumWaitTime"},
		{ie.TWANIdentifierTimestamp, 0, "WLANLocationTimeStamp"},
		{ie.FContainer, 0, "NBIFOMContainer"},
		{ie.RemoteUEContext, 0, "RemoteUEContextConnected"},
		{ie.NodeIdentifier, 0, "TGPPAAAServerIdentifier"},
		{ie.ExtendedProtocolConfigurationOptions, 0, "EPCO"},
		{ie.ServingPLMNRateControl, 0, "ServingPLMNRateControl"},
		{ie.Counter, 0, "MOExceptionDataCounter"},
		{ie.MappedUEUsageType, 0, "MappedUEUsageType"},
		{ie.FullyQualifiedDomainName, 0, "SGWUNodeName"},
		{ie.SecondaryRATUsageDataReport, 0, "SecondaryRATUsageDataReport"},
		{ie.UPFunctionSelectionIndicationFlags, 0, "UPFunctionSelectionIndicationFlags"},
		{ie.APNRateControlStatus, 0, "APNRateControlStatus"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.2-1.
	"CreateSessionResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.ChangeReportingAction, 0, "ChangeReportingAction"},
		{ie.CSGInformationReportingAction, 0, "CSGInformationReportingAction"},
		{ie.HeNBInformationReporting, 0, "HeNBInformationReporting"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEIDC"},
		{ie.FullyQualifiedTEID, 1, "PGWS5S8FTEIDC"},
		{ie.PDNAddressAllocation, 0, "PAA"},
		{ie.APNRestriction, 0, "APNRestriction"},
		{ie.AggregateMaximumBitRate, 0, "AMBR"},
		{ie.EPSBearerID, 0, "EBI"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.BearerContext, 0, "BearerContextsCreated"},
		{ie.BearerContext, 1, "BearerContextMarkedForRemoval"},
		{ie.Recovery, 0, "Recovery"},
		{ie.FullyQualifiedDomainName, 0, "ChargingGatewayName"},
		{ie.IPAddress, 0, "ChargingGatewayAddress"},
		{ie.FullyQualifiedCSID, 0, "PGWFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.LocalDistinguishedName, 1, "PGWLDN"},
		{ie.LocalDistinguishedName, 0, "SGWLDN"},
		{ie.EPCTimer, 0, "PGWBackOffTime"},
		{ie.AdditionalProtocolConfigurationOptions, 0, "APCO"},
		{ie.IPv4ConfigurationParameters, 0, "TrustedTWANIPv4Parameters"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.PresenceReportingAreaAction, 0, "PresenceReportingAreaAction"},
		{ie.LoadControlInformation, 0, "PGWNodeLoadControlInformation"},
		{ie.LoadControlInformation, 1, "PGWAPNLoadControlInformation"},
		{ie.LoadControlInformation, 2, "SGWNodeLoadControlInformation"},
		{ie.OverloadControlInformation, 0, "PGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.FContainer, 0, "NBIFOMContainer"},
		{ie.ChargingID, 0, "PDNConnectionChargingID"},
		{ie.ExtendedProtocolConfigurationOptions, 0, "EPCO"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.17.1-1.
	"DeleteBearerCommand": {
		{ie.BearerContext, 0, "BearerContexts"},
		{ie.UserLocationInformation, 0, "ULI"},
		{ie.ULITimestamp, 0, "ULITimestamp"},
		{ie.UETimeZone, 0, "UETimeZone"},
		{ie.OverloadControlInformation, 0, "MMESGSNOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEIDC"},
		{ie.SecondaryRATUsageDataReport, 0, "SecondaryRATDataUsageReport"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.17.2-1.
	"DeleteBearerFailureIndication": {
		{ie.Cause, 0, "Cause"},
		{ie.BearerContext, 0, "BearerContexts"},
		{ie.Recovery, 0, "Recovery"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.OverloadControlInformation, 0, "PGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.9.2-1.
	"DeleteBearerRequest": {
		{ie.EPSBearerID, 0, "LinkedEBI"},
		{ie.EPSBearerID, 1, "EBIs"},
		{ie.BearerContext, 0, "FailedBearerContext"},
		{ie.ProcedureTransactionID, 0, "PTI"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.FullyQualifiedCSID, 0, "PGWFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.Cause, 0, "Cause"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.LoadControlInformation, 0, "PGWNodeLoadControlInformation"},
		{ie.LoadControlInformation, 1, "PGWAPNLoadControlInformation"},
		{ie.LoadControlInformation, 2, "SGWNodeLoadControlInformation"},
		{ie.OverloadControlInformation, 0, "PGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.FContainer, 0, "NBIFOMContainer"},
		{ie.APNRateControlStatus, 0, "APNRateControlStatus"},
		{ie.ExtendedProtocolConfigurationOptions, 0, "EPCO"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.10.2-1.
	"DeleteBearerResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.EPSBearerID, 0, "LinkedEBI"},
		{ie.BearerContext, 0, "BearerContexts"},
		{ie.Recovery, 0, "Recovery"},
		{ie.FullyQualifiedCSID, 0, "MMEFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.FullyQualifiedCSID, 2, "EPDGFQCSID"},
		{ie.FullyQualifiedCSID, 3, "TWANFQCSID"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.UETimeZone, 0, "UETimeZone"},
		{ie.UserLocationInformation, 0, "ULI"},
		{ie.ULITimestamp, 0, "ULITimestamp"},
		{ie.TWANIdentifier, 0, "TWANIdentifier"},
		{ie.TWANIdentifier, 1, "WLANLocationInformation"},
		{ie.TWANIdentifierTimestamp, 0, "TWANIdentifierTimestamp"},
		{ie.TWANIdentifierTimestamp, 1, "WLANLocationTimestamp"},
		{ie.OverloadControlInformation, 0, "MMEOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 2, "TWANePDGOverloadControlInformation"},
		{ie.IPAddress, 0, "MMESGSNIdentifierOrUELocalIPAddress"},
		{ie.PortNumber, 0, "UEUDPPort"},
		{ie.PortNumber, 1, "UETCPPort"},
		{ie.FContainer, 0, "NBIFOMContainer"},
		{ie.SecondaryRATUsageDataReport, 0, "SecondaryRATUsageDataReport"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.9.1-1.
	"DeletePDNConnectionSetRequest": {
		{ie.FullyQualifiedCSID, 0, "MMEFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.FullyQualifiedCSID, 2, "PGWFQCSID"},
		{ie.FullyQualifiedCSID, 3, "EPDGFQCSID"},
		{ie.FullyQualifiedCSID, 4, "TWANFQCSID"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.9.2.
	"DeletePDNConnectionSetResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.Recovery, 0, "Recovery"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.9.1-1.
	"DeleteSessionRequest": {
		{ie.Cause, 0, "Cause"},
		{ie.EPSBearerID, 0, "LinkedEBI"},
		{ie.UserLocationInformation, 0, "ULI"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.NodeType, 0, "OriginatingNode"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEIDC"},
		{ie.UETimeZone, 0, "UETimeZone"},
		{ie.ULITimestamp, 0, "ULITimestamp"},
		{ie.RANNASCause, 0, "RANNASReleaseCause"},
		{ie.TWANIdentifier, 0, "TWANIdentifier"},
		{ie.TWANIdentifier, 1, "WLANLocationInformation"},
		{ie.TWANIdentifierTimestamp, 0, "TWANIdentifierTimestamp"},
		{ie.TWANIdentifierTimestamp, 1, "WLANLocationTimeStamp"},
		{ie.OverloadControlInformation, 0, "MMESGSNOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformaion"},
		{ie.OverloadControlInformation, 2, "TWANePDGOverloadControlInformaion"},
		{ie.IPAddress, 0, "UELocalIPAddress"},
		{ie.PortNumber, 0, "UEUDPPort"},
		{ie.PortNumber, 1, "UETCPPort"},
		{ie.ExtendedProtocolConfigurationOptions, 0, "EPCO"},
		{ie.SecondaryRATUsageDataReport, 0, "SecondaryRATUsageDataReport"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.10.1-1.
	"DeleteSessionResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.Recovery, 0, "Recovery"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.LoadControlInformation, 0, "PGWNodeLoadControlInformation"},
		{ie.LoadControlInformation, 1, "PGWAPNLoadControlInformation"},
		{ie.LoadControlInformation, 2, "SGWNodeLoadControlInformation"},
		{ie.OverloadControlInformation, 0, "PGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.ExtendedProtocolConfigurationOptions, 0, "EPCO"},
		{ie.APNRateControlStatus, 0, "APNRateControlStatus"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.3.13-1.
	"DetachAcknowledge": {
		{ie.Cause, 0, "Cause"},
		{ie.Recovery, 0, "Recovery"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.3.12-1.
	"DetachNotification": {
		{ie.Cause, 0, "Cause"},
		{ie.DetachType, 0, "DetachType"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.11.1-1.
	"DownlinkDataNotification": {
		{ie.Cause, 0, "Cause"},
		{ie.EPSBearerID, 0, "EPSBearerID"},
		{ie.AllocationRetensionPriority, 0, "AllocationRetensionPriority"},
		{ie.IMSI, 0, "IMSI"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEIDC"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.LoadControlInformation, 0, "SGWNodeLoadControlInformation"},
		{ie.OverloadControlInformation, 0, "SGWOverloadControlInformation"},
		{ie.PagingAndServiceInformation, 0, "PagingAndServiceInformation"},
		{ie.IntegerNumber, 0, "DLDataPacketsSize"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.11.2-1.
	"DownlinkDataNotificationAcknowledge": {
		{ie.Cause, 0, "Cause"},
		{ie.DelayValue, 0, "DataNotificationDelay"},
		{ie.Recovery, 0, "Recovery"},
		{ie.Throttling, 0, "DLLowPriorityTrafficThrottling"},
		{ie.IMSI, 0, "IMSI"},
		{ie.EPCTimer, 0, "DLBufferingDuration"},
		{ie.IntegerNumber, 0, "DLBufferingSuggestedPacketCount"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.11.3-1.
	"DownlinkDataNotificationFailureIndication": {
		{ie.Cause, 0, "Cause"},
		{ie.NodeType, 0, "OriginatingNode"},
		{ie.IMSI, 0, "IMSI"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.1.1-1.
	"EchoRequest": {
		{ie.Recovery, 0, "Recovery"},
		{ie.NodeFeatures, 0, "SendingNodeFeatures"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.1.2-1.
	"EchoResponse": {
		{ie.Recovery, 0, "Recovery"},
		{ie.NodeFeatures, 0, "SendingNodeFeatures"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.24-1.
	"ModifyAccessBearersRequest": {
		{ie.Indication, 0, "IndicationFlags"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEIDC"},
		{ie.DelayValue, 0, "DelayDownlinkPacketNotificationRequest"},
		{ie.BearerContext, 0, "BearerContextsToBeModified"},
		{ie.BearerContext, 1, "BearerContextsToBeRemoved"},
		{ie.Recovery, 0, "Recovery"},
		{ie.SecondaryRATUsageDataReport, 0, "SecondaryRATUsageDataReport"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.25-1.
	"ModifyAccessBearersResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.BearerContext, 0, "BearerContextsModified"},
		{ie.BearerContext, 1, "BearerContextsMarkedForRemoval"},
		{ie.Recovery, 0, "Recovery"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.LoadControlInformation, 0, "SGWNodeLoadControlInformation"},
		{ie.OverloadControlInformation, 0, "SGWOverloadControlInformation"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.14.1-1.
	"ModifyBearerCommand": {
		{ie.AggregateMaximumBitRate, 0, "APNAMBR"},
		{ie.BearerContext, 0, "BearerContext"},
		{ie.OverloadControlInformation, 0, "MMESGSNOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 2, "TWANePDGOverloadControlInformation"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEIDC"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.14.2-1.
	"ModifyBearerFailureIndication": {
		{ie.Cause, 0, "Cause"},
		{ie.Recovery, 0, "Recovery"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.OverloadControlInformation, 0, "PGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.7-1.
	"ModifyBearerRequest": {
		{ie.MobileEquipmentIdentity, 0, "MEI"},
		{ie.UserLocationInformation, 0, "ULI"},
		{ie.UserLocationInformation, 1, "ULIForSGW"},
		{ie.ServingNetwork, 0, "ServingNetwork"},
		{ie.RATType, 0, "RATType"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEIDC"},
		{ie.AggregateMaximumBitRate, 0, "AMBR"},
		{ie.DelayValue, 0, "DelayDownlinkPacketNotificationRequest"},
		{ie.BearerContext, 0, "BearerContextsToBeModified"},
		{ie.BearerContext, 1, "BearerContextsToBeRemoved"},
		{ie.Recovery, 0, "Recovery"},
		{ie.UETimeZone, 0, "UETimeZone"},
		{ie.FullyQualifiedCSID, 0, "MMEFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.UserCSGInformation, 0, "UCI"},
		{ie.IPAddress, 1, "UELocalIPAddress"},
		{ie.IPAddress, 0, "HeNBLocalIPAddress"},
		{ie.PortNumber, 1, "UEUDPPort"},
		{ie.PortNumber, 0, "HeNBUDPPort"},
		{ie.IPAddress, 2, "MMESGSNIdentifier"},
		{ie.LocalDistinguishedName, 0, "MMESGSNLDN"},
		{ie.LocalDistinguishedName, 1, "SGWLDN"},
		{ie.CNOperatorSelectionEntity, 0, "CNOperatorSelectionEntity"},
		{ie.PresenceReportingAreaInformation, 0, "PresenceReportingAreaInformation"},
		{ie.OverloadControlInformation, 0, "MMESGSNOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 2, "EPDGOverloadControlInformation"},
		{ie.ServingPLMNRateControl, 0, "ServingPLMNRateControl"},
		{ie.Counter, 0, "MOExceptionDataCounter"},
		{ie.IMSI, 0, "IMSI"},
		{ie.TWANIdentifier, 0, "WLANLocationInformation"},
		{ie.TWANIdentifierTimestamp, 0, "WLANLocationTimeStamp"},
		{ie.SecondaryRATUsageDataReport, 0, "SecondaryRATUsageDataReport"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.8-1.
	"ModifyBearerResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.MSISDN, 0, "MSISDN"},
		{ie.EPSBearerID, 0, "LinkedEBI"},
		{ie.APNRestriction, 0, "APNRestriction"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.BearerContext, 0, "BearerContextsModified"},
		{ie.BearerContext, 1, "BearerContextsMarkedForRemoval"},
		{ie.ChangeReportingAction, 0, "ChangeReportingAction"},
		{ie.CSGInformationReportingAction, 0, "CSGInformationReportingAction"},
		{ie.HeNBInformationReporting, 0, "HeNBInformationReporting"},
		{ie.FullyQualifiedDomainName, 0, "ChargingGatewayName"},
		{ie.IPAddress, 0, "ChargingGatewayAddress"},
		{ie.FullyQualifiedCSID, 0, "PGWFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.Recovery, 0, "Recovery"},
		{ie.LocalDistinguishedName, 0, "SGWLDN"},
		{ie.LocalDistinguishedName, 1, "PGWLDN"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.PresenceReportingAreaAction, 0, "PresenceReportingAreaAction"},
		{ie.LoadControlInformation, 0, "PGWNodeLoadControlInformation"},
		{ie.LoadControlInformation, 1, "PGWAPNLoadControlInformation"},
		{ie.LoadControlInformation, 2, "SGWNodeLoadControlInformation"},
		{ie.OverloadControlInformation, 0, "PGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.ChargingID, 0, "PDNConnectionChargingID"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.9.5-1.
	"PGWRestartNotification": {
		{ie.IPAddress, 0, "PGWS5S8IPAddressForControlPlaneOrPMIP"},
		{ie.IPAddress, 1, "SGWS11S4IPAddressForControlPlane"},
		{ie.Cause, 0, "Cause"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.9.6-1.
	"PGWRestartNotificationAcknowledge": {
		{ie.Cause, 0, "Cause"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.21-1.
	"ReleaseAccessBearersRequest": {
		{ie.EPSBearerID, 0, "ListOfRABs"},
		{ie.NodeType, 0, "OriginatingNode"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.SecondaryRATUsageDataReport, 0, "SecondaryRATUsageDataReport"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.22-1.
	"ReleaseAccessBearersResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.Recovery, 0, "Recovery"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.LoadControlInformation, 0, "SGWNodeLoadControlInformation"},
		{ie.OverloadControlInformation, 0, "SGWOverloadControlInformation"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.4.4-1.
	"ResumeAcknowledge": {
		{ie.Cause, 0, "Cause"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.4.3-1.
	"ResumeNotification": {
		{ie.IMSI, 0, "IMSI"},
		{ie.EPSBearerID, 0, "LinkedEPSBearerID"},
		{ie.NodeType, 0, "OriginatingNode"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEIDForControlPlane"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.23-1.
	"StopPagingIndication": {
		{ie.IMSI, 0, "IMSI"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.4.2-1.
	"SuspendAcknowledge": {
		{ie.Cause, 0, "Cause"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.4.1-1.
	"SuspendNotification": {
		{ie.IMSI, 0, "IMSI"},
		{ie.UserLocationInformation, 0, "RAI"},
		{ie.EPSBearerID, 0, "LinkedEBI"},
		{ie.PacketTMSI, 0, "PTMSI"},
		{ie.NodeType, 0, "OriginatingNode"},
		{ie.IPAddress, 0, "AddressForControlPlane"},
		{ie.PortNumber, 0, "UDPSourcePortNumber"},
		{ie.HopCounter, 0, "HopCounter"},
		{ie.FullyQualifiedTEID, 0, "SenderFTEIDC"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.15-1.
	"UpdateBearerRequest": {
		{ie.BearerContext, 0, "BearerContexts"},
		{ie.ProcedureTransactionID, 0, "PTI"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.AggregateMaximumBitRate, 0, "APNAMBR"},
		{ie.ChangeReportingAction, 0, "ChangeReportingAction"},
		{ie.CSGInformationReportingAction, 0, "CSGInformationReportingAction"},
		{ie.HeNBInformationReporting, 0, "HeNBInformationReporting"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.FullyQualifiedCSID, 0, "PGWFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.PresenceReportingAreaAction, 0, "PresenceReportingAction"},
		{ie.LoadControlInformation, 0, "PGWNodeLoadControlInformation"},
		{ie.LoadControlInformation, 1, "PGWAPNLoadControlInformation"},
		{ie.LoadControlInformation, 2, "SGWNodeLoadControlInformation"},
		{ie.OverloadControlInformation, 0, "PGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.FContainer, 0, "NBIFOMContainer"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.2.16-1.
	"UpdateBearerResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.BearerContext, 0, "BearerContexts"},
		{ie.ProtocolConfigurationOptions, 0, "PCO"},
		{ie.Recovery, 0, "Recovery"},
		{ie.FullyQualifiedCSID, 0, "MMEFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.FullyQualifiedCSID, 2, "EPDGFQCSID"},
		{ie.FullyQualifiedCSID, 3, "TWANFQCSID"},
		{ie.Indication, 0, "IndicationFlags"},
		{ie.UETimeZone, 0, "UETimeZone"},
		{ie.UserLocationInformation, 0, "ULI"},
		{ie.TWANIdentifier, 0, "TWANIdentifier"},
		{ie.TWANIdentifier, 1, "WLANLocationInformation"},
		{ie.OverloadControlInformation, 0, "MMESGSNOverloadControlInformation"},
		{ie.OverloadControlInformation, 1, "SGWOverloadControlInformation"},
		{ie.OverloadControlInformation, 2, "TWANePDGOverloadControlInformation"},
		{ie.PresenceReportingAreaInformation, 0, "PresenceReportingAreaInformation"},
		{ie.IPAddress, 0, "MMESGSNIdentifierOrUELocalIPAddress"},
		{ie.TWANIdentifierTimestamp, 1, "WLANLocationTimeStamp"},
		{ie.PortNumber, 0, "UEUDPPort"},
		{ie.PortNumber, 1, "UETCPPort"},
		{ie.FContainer, 0, "NBIFOMContainer"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.9.3-1.
	"UpdatePDNConnectionSetRequest": {
		{ie.FullyQualifiedCSID, 0, "MMEFQCSID"},
		{ie.FullyQualifiedCSID, 1, "SGWFQCSID"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 7.9.4-1.
	"UpdatePDNConnectionSetResponse": {
		{ie.Cause, 0, "Cause"},
		{ie.FullyQualifiedCSID, 0, "PGWFQCSID"},
		{ie.Recovery, 0, "Recovery"},
		{ie.PrivateExtension, 0, "PrivateExtension"},
	},
	// TS 29.274 Table 6.1-1.
	"VersionNotSupportedIndication": {},
}

func TestCatalogueFields(t *testing.T) {
	for _, factory := range codecFactories {
		if factory.name == "Generic" {
			continue
		}
		t.Run(factory.name, func(t *testing.T) {
			fields := catalogue[factory.name]
			all := factory.make()
			v := reflect.ValueOf(all).Elem()
			seen := make(map[string]bool)
			for _, field := range fields {
				if seen[field.name] {
					continue
				} // MM Context has several alternative types.
				seen[field.name] = true
				i := codecRoleIE(field.typ, field.instance, field.name)
				dst := v.FieldByName(field.name)
				if dst.Kind() == reflect.Slice {
					dst.Set(reflect.ValueOf([]*ie.IE{i}))
				} else {
					dst.Set(reflect.ValueOf(i))
				}
			}
			for j := 1; j < v.NumField(); j++ {
				name := v.Type().Field(j).Name
				if name != "AdditionalIEs" && !seen[name] {
					t.Errorf("field %s has no catalogue role", name)
				}
			}
			parsed, err := Parse(codecWire(t, all))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(codecFields(all), codecFields(parsed)) {
				t.Fatal("complete catalogue does not round trip")
			}
			for _, field := range fields {
				t.Run(field.name, func(t *testing.T) {
					i := codecRoleIE(field.typ, field.instance, field.name)
					constructed := factory.make(i)
					raw, _ := i.Marshal()
					if !bytes.Equal(codecFields(constructed)[field.name], raw) {
						t.Fatal("constructor misrouted field")
					}
					wire := codecPacket(t, constructed.MessageType(), i)
					decoded, err := Parse(wire)
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(codecFields(decoded)[field.name], raw) {
						t.Fatal("decoder misrouted field")
					}
					// Distinct duplicate payloads prove the first singleton survives, while
					// every repeated field keeps both occurrences in order.
					second := codecRoleIE(field.typ, field.instance, field.name)
					if second.IsGrouped() {
						second = ie.NewGroupedIE(field.typ, ie.NewEPSBearerID(7)).WithInstance(field.instance)
					} else {
						second.Payload = bytes.Repeat([]byte{0x42}, len(second.Payload))
					}
					repeated := factory.make(i, second)
					output := codecFields(repeated)[field.name]
					repeatedField := reflect.ValueOf(repeated).Elem().FieldByName(field.name).Kind() == reflect.Slice
					want := append([]byte{}, raw...)
					if repeatedField {
						b, marshalErr := second.Marshal()
						if marshalErr != nil {
							t.Fatal(marshalErr)
						}
						want = append(want, b...)
					}
					if !bytes.Equal(output, want) {
						t.Fatal("constructor repetition policy")
					}
					decoded, err = Parse(codecPacket(t, constructed.MessageType(), i, second))
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(codecFields(decoded)[field.name], want) {
						t.Fatal("decoder repetition policy")
					}
				})
			}
		})
	}
}

func TestParseErrorDetails(t *testing.T) {
	for _, version := range []uint8{0, 1, 3, 7} {
		_, err := Parse([]byte{version << 5, MsgTypeCreateSessionRequest})
		var ve *UnsupportedVersionError
		if !errors.Is(err, ErrUnsupportedVersion) || !errors.As(err, &ve) || ve.Version != version {
			t.Fatalf("version %d: %v", version, err)
		}
	}
	for n := 0; n < 12; n++ {
		b := make([]byte, n)
		if n > 0 {
			b[0] = 0x48
		}
		if _, err := Parse(b); !errors.Is(err, ErrTooShortToParse) {
			t.Fatalf("size %d: %v", n, err)
		}
	}
	b := codecWire(t, NewCreateSessionRequest(0x12345678, 0xabcdef))
	b[3]++
	_, err := Parse(b)
	var le *LengthError
	if !errors.Is(err, ErrInvalidLength) || !errors.As(err, &le) || le.MessageType != MsgTypeCreateSessionRequest || le.TEID != 0x12345678 || le.SequenceNumber != 0xabcdef || le.Declared != 13 || le.Actual != 12 {
		t.Fatalf("length details: %#v %v", le, err)
	}
	_, err = ParseEchoRequest(codecWire(t, NewEchoResponse(1)))
	var te *UnexpectedTypeError
	if !errors.Is(err, ErrUnexpectedType) || !errors.As(err, &te) || te.Expected != MsgTypeEchoRequest || te.Actual != MsgTypeEchoResponse {
		t.Fatalf("type details: %#v %v", te, err)
	}
}

func TestMessageIEConstructionError(t *testing.T) {
	bad := ie.New(ie.IMSI, 0, make([]byte, 65536))
	if bad.Err() == nil {
		t.Fatal("fixture must carry a construction error")
	}
	for _, f := range codecFactories {
		t.Run(f.name, func(t *testing.T) {
			_, err := Marshal(f.make(bad))
			if !errors.Is(err, ie.ErrInvalidLength) {
				t.Fatalf("construction error lost: %v", err)
			}
		})
	}
}

func TestNilMessageAccessors(t *testing.T) {
	for _, factory := range codecFactories {
		t.Run(factory.name, func(t *testing.T) {
			m := reflect.Zero(reflect.TypeOf(factory.make())).Interface().(Message)
			if m.TEID() != 0 || m.Sequence() != 0 || m.Version() != 0 || m.MessageType() != 0 || m.MarshalLen() != 0 || m.GetHeader() != nil {
				t.Fatal("nonzero nil accessor")
			}
			if err := m.UnmarshalBinary(nil); !errors.Is(err, ErrNilMessage) {
				t.Fatalf("nil decode: %v", err)
			}
			if err := m.SetSequenceNumber(1); !errors.Is(err, ErrNilMessage) {
				t.Fatalf("nil sequence setter: %v", err)
			}
			m.SetTEID(1)
			if _, err := Marshal(m); !errors.Is(err, ErrNilMessage) {
				t.Fatalf("nil marshal: %v", err)
			}
		})
	}
	if _, err := Marshal(nil); !errors.Is(err, ErrNilMessage) {
		t.Fatal(err)
	}
}

func TestNamedFieldValidation(t *testing.T) {
	for _, factory := range codecFactories {
		if factory.name == "Generic" || factory.name == "VersionNotSupportedIndication" {
			continue
		}
		for _, field := range catalogue[factory.name] {
			t.Run(factory.name+"/"+field.name, func(t *testing.T) {
				m := factory.make()
				dst := reflect.ValueOf(m).Elem().FieldByName(field.name)
				bad := ie.New(250, 15, nil)
				if dst.Kind() == reflect.Slice {
					dst.Set(reflect.ValueOf([]*ie.IE{bad}))
				} else {
					dst.Set(reflect.ValueOf(bad))
				}
				if _, err := Marshal(m); !errors.Is(err, ErrUnexpectedIE) {
					t.Fatalf("invalid field accepted: %v", err)
				}
			})
		}
	}
}
