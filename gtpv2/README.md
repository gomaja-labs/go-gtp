# GTPv2-C

GTPv2-C codecs and `Conn` follow TS 29.274 V20.0.0. `Conn` handles UDP I/O,
message dispatch, sequence allocation, retransmission and duplicate request
response replay. Sessions and bearers are in memory; applications implement
procedure policy and persistence.

## Codecs

Import this version's `ie` and `message` subpackages. This complete function is
[Example_codecs](example_test.go); it builds a message with an IE, marshals it,
parses a datagram, checks the concrete type and reads the IE with error handling.

```go
func Example_codecs() {
	msg := message.NewEchoResponse(1, ie.NewRecovery(0))
	wire, err := message.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	parsed, err := message.Parse(wire)
	if err != nil {
		fmt.Println(err)
		return
	}
	echo, ok := parsed.(*message.EchoResponse)
	if !ok {
		fmt.Printf("unexpected message %T\n", parsed)
		return
	}
	recovery, err := echo.Recovery.Recovery()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(echo.MessageTypeName(), recovery)
	// Output: Echo Response 0
}
```

IE constructors return an error-carrying `*ie.IE`, not a separate error. Check
`Err()` and use typed accessors, which return `ie.ErrIENotFound` on nil receivers.
[ExampleIE_Err](ie/example_test.go) is an executable example:

```go
func ExampleIE_Err() {
	invalid := ie.NewIMSI("1234567890123456")
	fmt.Println(errors.Is(invalid.Err(), ie.ErrInvalidLength))

	var missing *ie.IE
	_, err := missing.IMSI()
	fmt.Println(errors.Is(err, ie.ErrIENotFound))

	valid := ie.NewIMSI("001010123456789")
	if err := valid.Err(); err != nil {
		fmt.Println(err)
		return
	}
	imsi, err := valid.IMSI()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(imsi)
	// Output:
	// true
	// true
	// 001010123456789
}
```

See the [root conventions](../README.md#api-conventions) for buffer lifetimes,
decode error handling and marshaling. Generic codecs preserve framed data; they
check header and IE boundaries and retain IEs in received order without enforcing
mandatory presence, IE contents or IE order. Decoded Generics re-encode datagrams
byte for byte, including received header spare and unused fields. Every IE,
including appended IEs, uses `MarshalToLossless`, which preserves contents and
spare bits without content validation; construction errors and unencodable
lengths remain errors. `IE.Err()` and accessors expose content errors. Generics
built with `NewGeneric` use normal IE and header encoding checks. Generic codecs
do not implement the procedure rules for unsupported message types.
`NewGenericWithoutTEID` also uses normal encoding checks. `ParseGeneric` retains
both messages of a piggybacked datagram as `*message.Generic`, or preserves an
undecodable secondary for replay. Inspect `PiggybackError` independently; replay
does not make the secondary valid. `Parse` and `ParseHeader` retain typed
secondary decoding. Ordinary IE `Marshal`/`MarshalTo` reject malformed grouped
contents and clear IE header spare bits.

## Conn

`NewConn(localAddr, interfaceType, restartCounter)` creates an unbound server.
Register `HandlerFunc` callbacks before starting blocking `ListenAndServe(ctx)`.
Alternatively use `Listen(ctx)` then `Serve(ctx)` to control startup. Handlers
receive `*gtpv2.Conn`, the peer address and a `message.Message`. Default handlers
answer Echo Request and process Echo Response and Version Not Supported.

[ExampleConn](example_test.go) replaces the Echo Request handler.
`RespondTo` and `EchoResponse` copy the request sequence into the reply:

```go
func ExampleConn() {
	local := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 2123}
	c := gtpv2.NewConn(local, gtpv2.IFTypeS11S4SGWGTPC, 0)
	defer func() { _ = c.Close() }()
	// Echo is handled by default; this shows how to replace a handler.
	c.AddHandler(message.MsgTypeEchoRequest, func(c *gtpv2.Conn, peer net.Addr, msg message.Message) error {
		return c.EchoResponse(peer, msg)
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := c.ListenAndServe(ctx); err != nil {
		fmt.Println(err)
	}
}
```

`Dial` performs an Echo handshake and starts serving. Its context governs only
the handshake; `Close` stops the returned connection. Both client and server can
communicate with several peers through one local UDP endpoint.
[ExampleDial](example_test.go) connects to the server above:

```go
func ExampleDial() {
	// Run the server example first. Dial performs an Echo exchange and starts serving.
	local := &net.UDPAddr{IP: net.ParseIP("127.0.0.1")}
	peer := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 2123}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, err := gtpv2.Dial(ctx, local, peer, gtpv2.IFTypeS11MMEGTPC, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = c.Close() }()
	// The handshake context does not govern the returned connection's lifetime.
	fmt.Println(c.LocalAddr())
}
```

### Sessions and transactions

`CreateSession` sends a request and returns a session and sequence; it does not
wait for success. Supply IEs required by the interface and procedure. Register
response handlers before sending. A handler can use `GetSessionByTEID` and
`PassMessageTo` to deliver a response to a session's `WaitMessage(seq, timeout)`
called from outside the serving goroutine. Check the response type, Cause and
IE errors before activating the session. Keep peer TEIDs distinct from local
TEIDs: the connection registry indexes locally allocated receiving TEIDs.

On the server, use `NewSession`, populate TEIDs/bearers, `RegisterSession` and
`RespondTo` after validating the request. `ParseCreateSession` extracts supported
fields but does not implement all conditional procedure rules. `RemoveSession`
cleans up local state; sending Delete Session alone does not perform that cleanup.
This [ExampleSession](example_test.go) illustrates registry operations
without network I/O or an over-the-wire session establishment:

```go
func ExampleSession() {
	local := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 2123}
	peer := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 32123}
	c := gtpv2.NewConn(local, gtpv2.IFTypeS11S4SGWGTPC, 0)
	defer func() { _ = c.Close() }()
	session := gtpv2.NewSession(peer, &gtpv2.Subscriber{IMSI: "001010123456789"})
	session.AddTEID(gtpv2.IFTypeS11S4SGWGTPC, 0x11223344)
	c.RegisterSession(0x11223344, session)
	found, err := c.GetSessionByTEID(0x11223344)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := found.Activate(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(found.IMSI, found.IsActive(), c.SessionCount())
	c.RemoveSession(found)
	fmt.Println(c.SessionCount())
	// Output:
	// 001010123456789 true 1
	// 0
}
```

Use `SendMessageTo` for structured initial messages and `RespondTo` for replies.
`SendRawTo` registers raw requests for response matching; `WriteTo` does not.
`SetRetransmission` configures retry timing, `SetReplayWindow` configures response
replay retention and `SetErrorHandler` reports background errors. Validation is
enabled by default; application handlers still check interface-specific and
state-dependent requirements and authorize the sender. Handlers run sequentially
on the serving goroutine; never call `WaitMessage` there.

### Logging

[ExampleSetLogger](example_test.go) configures the package logger:

```go
func ExampleSetLogger() {
	previous := slog.Default()
	gtpv2.SetLogger(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
	defer gtpv2.SetLogger(previous)
	gtpv2.SetLogger(nil)
	// Output:
}
```

## Message support

These types have dedicated structs, constructors and parsing dispatch in
[message/message.go](message/message.go). Other types parse as `message.Generic`;
`NewGeneric` constructs one. A named type constant alone is not typed support.

| Type | Typed message | Constructor |
| --- | --- | --- |
| 1 | [EchoRequest](message/echo-req.go) | `NewEchoRequest` |
| 2 | [EchoResponse](message/echo-res.go) | `NewEchoResponse` |
| 3 | [VersionNotSupportedIndication](message/version-not-supported.go) | `NewVersionNotSupportedIndication` |
| 32 | [CreateSessionRequest](message/create-session-req.go) | `NewCreateSessionRequest` |
| 33 | [CreateSessionResponse](message/create-session-res.go) | `NewCreateSessionResponse` |
| 34 | [ModifyBearerRequest](message/modify-bearer-req.go) | `NewModifyBearerRequest` |
| 35 | [ModifyBearerResponse](message/modify-bearer-res.go) | `NewModifyBearerResponse` |
| 36 | [DeleteSessionRequest](message/delete-session-req.go) | `NewDeleteSessionRequest` |
| 37 | [DeleteSessionResponse](message/delete-session-res.go) | `NewDeleteSessionResponse` |
| 38 | [ChangeNotificationRequest](message/change-notification-req.go) | `NewChangeNotificationRequest` |
| 39 | [ChangeNotificationResponse](message/change-notification-res.go) | `NewChangeNotificationResponse` |
| 64 | [ModifyBearerCommand](message/modify-bearer-command.go) | `NewModifyBearerCommand` |
| 65 | [ModifyBearerFailureIndication](message/modify-bearer-failure-indication.go) | `NewModifyBearerFailureIndication` |
| 66 | [DeleteBearerCommand](message/delete-bearer-command.go) | `NewDeleteBearerCommand` |
| 67 | [DeleteBearerFailureIndication](message/delete-bearer-failure-indication.go) | `NewDeleteBearerFailureIndication` |
| 70 | [DownlinkDataNotificationFailureIndication](message/downlink-data-notification-failure-indication.go) | `NewDownlinkDataNotificationFailureIndication` |
| 73 | [StopPagingIndication](message/stop-paging-indication.go) | `NewStopPagingIndication` |
| 95 | [CreateBearerRequest](message/create-bearer-req.go) | `NewCreateBearerRequest` |
| 96 | [CreateBearerResponse](message/create-bearer-res.go) | `NewCreateBearerResponse` |
| 97 | [UpdateBearerRequest](message/update-bearer-req.go) | `NewUpdateBearerRequest` |
| 98 | [UpdateBearerResponse](message/update-bearer-res.go) | `NewUpdateBearerResponse` |
| 99 | [DeleteBearerRequest](message/delete-bearer-req.go) | `NewDeleteBearerRequest` |
| 100 | [DeleteBearerResponse](message/delete-bearer-res.go) | `NewDeleteBearerResponse` |
| 101 | [DeletePDNConnectionSetRequest](message/delete-pdn-connection-set-req.go) | `NewDeletePDNConnectionSetRequest` |
| 102 | [DeletePDNConnectionSetResponse](message/delete-pdn-connection-set-res.go) | `NewDeletePDNConnectionSetResponse` |
| 130 | [ContextRequest](message/context-req.go) | `NewContextRequest` |
| 131 | [ContextResponse](message/context-res.go) | `NewContextResponse` |
| 132 | [ContextAcknowledge](message/context-ack.go) | `NewContextAcknowledge` |
| 149 | [DetachNotification](message/detach-notification.go) | `NewDetachNotification` |
| 150 | [DetachAcknowledge](message/detach-acknowledge.go) | `NewDetachAcknowledge` |
| 162 | [SuspendNotification](message/suspend-notification.go) | `NewSuspendNotification` |
| 163 | [SuspendAcknowledge](message/suspend-acknowledge.go) | `NewSuspendAcknowledge` |
| 164 | [ResumeNotification](message/resume-notification.go) | `NewResumeNotification` |
| 165 | [ResumeAcknowledge](message/resume-acknowledge.go) | `NewResumeAcknowledge` |
| 170 | [ReleaseAccessBearersRequest](message/release-access-bearers-req.go) | `NewReleaseAccessBearersRequest` |
| 171 | [ReleaseAccessBearersResponse](message/release-access-bearers-res.go) | `NewReleaseAccessBearersResponse` |
| 176 | [DownlinkDataNotification](message/downlink-data-notification.go) | `NewDownlinkDataNotification` |
| 177 | [DownlinkDataNotificationAcknowledge](message/downlink-data-notification-ack.go) | `NewDownlinkDataNotificationAcknowledge` |
| 179 | [PGWRestartNotification](message/pgw-restart-notification.go) | `NewPGWRestartNotification` |
| 180 | [PGWRestartNotificationAcknowledge](message/pgw-restart-notification-acknowledge.go) | `NewPGWRestartNotificationAcknowledge` |
| 200 | [UpdatePDNConnectionSetRequest](message/update-pdn-connection-set-req.go) | `NewUpdatePDNConnectionSetRequest` |
| 201 | [UpdatePDNConnectionSetResponse](message/update-pdn-connection-set-res.go) | `NewUpdatePDNConnectionSetResponse` |
| 211 | [ModifyAccessBearersRequest](message/modify-access-bearers-req.go) | `NewModifyAccessBearersRequest` |
| 212 | [ModifyAccessBearersResponse](message/modify-access-bearers-res.go) | `NewModifyAccessBearersResponse` |

## IE support

These source files provide typed constructors returning `*ie.IE`; associated
field structs and accessors are documented in the linked files. Types without
such helpers can use `ie.New` with raw payloads and the generic scalar helpers
in [ie/ie.go](ie/ie.go). This table lists helper availability, not complete
procedure or feature coverage. Constants cover IE types 214–223 and SRVCC
Cause (56) from TS 29.274 Table 8.1-1.

| IE family (source) | Typed constructors returning `*IE` |
| --- | --- |
| [ambr](ie/ambr.go) | `NewAggregateMaximumBitRate` |
| [apn-restriction](ie/apn-restriction.go) | `NewAPNRestriction` |
| [apn](ie/apn.go) | `NewAccessPointName` |
| [arp](ie/arp.go) | `NewAllocationRetensionPriority` |
| [bearer-context](ie/bearer-context.go) | `NewBearerContext`, `NewBearerContextWithinCreateBearerRequest`, `NewBearerContextWithinCreateBearerResponse`, `NewBearerContextWithinDeleteBearerRequest`, `NewBearerContextWithinDeleteBearerResponse`, `NewBearerContextWithinModifyBearerCommand`, `NewBearerContextWithinUpdateBearerRequest`, `NewBearerContextWithinUpdateBearerResponse`, `NewBearerContextWithinDeleteBearerCommand`, `NewBearerContextWithinDeleteBearerFailureIndication`, `NewBearerContextWithinCreateIndirectDataForwardingTunnelRequest`, `NewBearerContextWithinCreateIndirectDataForwardingTunnelResponse`, `NewBearerContextWithinForwardRelocationRequest`, `NewBearerContextWithinContextResponse`, `NewBearerContextWithinContextAcknowledge` |
| [bearer-flags](ie/bearer-flags.go) | `NewBearerFlags(uhrrei, asi, vInd, vb, ppc)`; `HasUHRREI` reads UHRREI (TS 29.274 §8.32) |
| [bearer-qos](ie/bearer-qos.go) | `NewBearerQoS` |
| [bearer-tft](ie/bearer-tft.go) | `NewBearerTFT`, `NewBearerTFTCreateNewTFT`, `NewBearerTFTAddPacketFilters`, `NewBearerTFTReplacePacketFilters`, `NewBearerTFTDeletePacketFilters`, `NewBearerTFTDeleteExistingTFT`, `NewBearerTFTNoTFTOperation` |
| [cause](ie/cause.go) | `NewCause` |
| [charging-characteristics](ie/charging-characteristics.go) | `NewChargingCharacteristics` |
| [charging-id](ie/charging-id.go) | `NewChargingID` |
| [cmi](ie/cmi.go) | `NewCSGMembershipIndication` |
| [csg-id](ie/csg-id.go) | `NewCSGID` |
| [delay-value](ie/delay-value.go) | `NewDelayValue`, `NewDelayValueRaw` |
| [detach-type](ie/detach-type.go) | `NewDetachType` |
| [ebi](ie/ebi.go) | `NewEPSBearerID` |
| [epc-timer](ie/epc-timer.go) | `NewEPCTimer`, `NewEPCTimerRaw` |
| [f-teid](ie/f-teid.go) | `NewFullyQualifiedTEID`, `NewFullyQualifiedTEIDNetIP`; `gtpv2.IFTypeN26AMFGTPC` (40) and `gtpv2.IFTypeN19mbUPFGTPU` (41) identify N26 and N19mb endpoints (TS 29.274 §8.22) |
| [flow-qos](ie/flow-qos.go) | `NewFlowQoS` |
| [fq-csid](ie/fq-csid.go) | `NewFullyQualifiedCSID` |
| [fqdn](ie/fqdn.go) | `NewFullyQualifiedDomainName` |
| [global-cn-id](ie/global-cn-id.go) | `NewGlobalCNID` |
| [guti](ie/guti.go) | `NewGUTI` |
| [hop-counter](ie/hop-counter.go) | `NewHopCounter` |
| [ie_grouped](ie/ie_grouped.go) | `NewGroupedIE` |
| [imsi](ie/imsi.go) | `NewIMSI` |
| [indication](ie/indication.go) | `NewIndication(flags ...IndicationFlag)`, `NewIndicationFromBitSequence`, `NewIndicationFromOctets`; named flags use the shortest encoding, at least two octets. `Has<X>` follows `Indication<X>`, e.g. `HasPMTSMI`, `HasRetLoc`, `HasPPOFF`, `Has5GSIWKI`; `HasPPON` reads PPON / PPEI (TS 29.274 §8.12) |
| [integer-number](ie/integer-number.go) | `NewIntegerNumber` |
| [ip-addr](ie/ip-addr.go) | `NewIPAddress`, `NewIPAddressNetIP` |
| [ldn](ie/ldn.go) | `NewLocalDistinguishedName` |
| [mbms-flags](ie/mbms-flags.go) | `NewMBMSFlags` |
| [mei](ie/mei.go) | `NewMobileEquipmentIdentity` |
| [msisdn](ie/msisdn.go) | `NewMSISDN` |
| [node-features](ie/node-features.go) | `NewNodeFeatures`; `HasPSSET` reads PGW-C/SMF Set support (TS 29.274 §8.83) |
| [node-type](ie/node-type.go) | `NewNodeType` |
| [p-tmsi-signature](ie/p-tmsi-signature.go) | `NewPTMSISignature` |
| [p-tmsi](ie/p-tmsi.go) | `NewPacketTMSI` |
| [paa](ie/paa.go) | `NewPDNAddressAllocation`, `NewPDNAddressAllocationIPv6`, `NewPDNAddressAllocationDual`, `NewPDNAddressAllocationNetIP`, `NewPDNAddressAllocationDualNetIP` |
| [paging-and-service-information](ie/paging-and-service-information.go) | `NewPagingAndServiceInformation` |
| [pco](ie/pco.go) | `NewProtocolConfigurationOptions`, `NewAdditionalProtocolConfigurationOptions`, `NewExtendedProtocolConfigurationOptions`, `NewLargeAdditionalProtocolConfigurationOptions` |
| [pdn-type](ie/pdn-type.go) | `NewPDNType` |
| [pgw-change-info](ie/pgw-change-info.go) | `NewPGWChangeInfo(ies ...*IE)` constructs a grouped IE (TS 29.274 §8.145) |
| [plmn-id](ie/plmn-id.go) | `NewPLMNID` |
| [port-number](ie/port-number.go) | `NewPortNumber` |
| [private-extension](ie/private-extension.go) | `NewPrivateExtension` |
| [pti](ie/pti.go) | `NewProcedureTransactionID` |
| [ran-nas-cause](ie/ran-nas-cause.go) | `NewRANNASCause` |
| [rat-type](ie/rat-type.go) | `NewRATType` |
| [recovery](ie/recovery.go) | `NewRecovery` |
| [rfsp-index](ie/rfsp-index.go) | `NewRFSPIndex` |
| [s103pdf](ie/s103pdf.go) | `NewS103PDNDataForwardingInfo`, `NewS103PDNDataForwardingInfoNetIP` |
| [s1udf](ie/s1udf.go) | `NewS1UDataForwarding`, `NewS1UDataForwardingNetIP` |
| [secondary-rat-usage-data-report](ie/secondary-rat-usage-data-report.go) | `NewSecondaryRATUsageDataReport` |
| [selection-mode](ie/selection-mode.go) | `NewSelectionMode` |
| [service-indicator](ie/service-indicator.go) | `NewServiceIndicator` |
| [serving-nw](ie/serving-nw.go) | `NewServingNetwork` |
| [tad](ie/tad.go) | `NewTrafficAggregateDescription`, `NewTrafficAggregateDescriptionCreateNewTFT`, `NewTrafficAggregateDescriptionAddPacketFilters`, `NewTrafficAggregateDescriptionReplacePacketFilters`, `NewTrafficAggregateDescriptionDeletePacketFilters`, `NewTrafficAggregateDescriptionDeleteExistingTFT`, `NewTrafficAggregateDescriptionNoTFTOperation` |
| [throttling](ie/throttling.go) | `NewThrottling` |
| [tmsi](ie/tmsi.go) | `NewTMSI` |
| [trace-reference](ie/trace-reference.go) | `NewTraceReference` |
| [type-extension](ie/type-extension.go) | `NewWithTypeExtension(ext uint16, instance uint8, content []byte)`; `TypeExtension() (uint16, error)` reads the extended type, range 256–65535 (TS 29.274 §8.2.1A) |
| [uci](ie/uci.go) | `NewUserCSGInformation` |
| [ue-timezone](ie/ue-timezone.go) | `NewUETimeZone` |
| [uli-timestamp](ie/uli-timestamp.go) | `NewULITimestamp` |
| [uli](ie/uli.go) | `NewUserLocationInformation` |
