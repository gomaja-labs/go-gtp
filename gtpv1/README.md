# GTPv1

GTPv1-C codecs follow TS 29.060 V19.0.0. GTPv1-U codecs and `UPlaneConn` follow
TS 29.281 V19.2.0; PDU Session Container fields follow TS 38.415 V19.1.0.
GTPv1-C has codecs only; no control-plane connection or PDP context manager is
provided.

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

`message.Parse` rejects descending IE types with `message.ErrIEOutOfSequence`
and returns no message, also for message types it decodes as `Generic`.
Repeated equal types are in order. Missing mandatory IEs
and incorrect mandatory values take precedence over ordering errors (TS 29.060
§11.1; GSM 09.60 §10.1). Typed encoders sort IEs stably (TS 29.060 §7.7.0;
GSM 09.60 §7.9). A caller that must answer a rejected request with "Invalid
message format" or "Mandatory IE missing" can use `ParseGeneric` to recover its
header and IEs.

## UPlaneConn

`DialUPlane` binds a UDP socket, performs an Echo handshake and starts serving.
Its context bounds the handshake; call `Close` to stop the returned connection.
`NewUPlaneConn` plus blocking `ListenAndServe(ctx)` starts a server without that
handshake. `NewUPlaneConnFromUDPConn` takes ownership of an already bound socket.
Register handlers with `AddHandler` before serving; a GTPv1 handler takes
`gtpv1.Conn`, a peer `net.Addr` and a `message.Message`.

This [send/receive example](example_test.go) needs a peer that answers
Echo and sends a T-PDU. `ReadFromGTP` returns the decapsulated payload, peer and
TEID; `WriteToGTP` takes the outgoing TEID, payload and destination. Unknown TEIDs
normally trigger Error Indication; `DisableErrorIndication` permits application
reads instead. Set deadlines for bounded reads.

```go
func ExampleUPlaneConn() {
	// A GTP-U peer must be listening at this address and answer Echo.
	peer := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 2152}
	local := &net.UDPAddr{IP: net.ParseIP("127.0.0.1")}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	u, err := gtpv1.DialUPlane(ctx, local, peer)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = u.Close() }()

	// Enable delivery of unknown TEIDs to ReadFromGTP instead of Error Indication.
	u.DisableErrorIndication()
	if _, err := u.WriteToGTP(0x11223344, []byte("payload"), peer); err != nil {
		fmt.Println(err)
		return
	}
	if err := u.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 65535)
	n, from, teid, err := u.ReadFromGTP(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("peer=%s TEID=%08x payload=%x\n", from, teid, buf[:n])
}
```

### Relay

`incoming.RelayTo(outgoing, incomingTEID, outgoingTEID, peer)` installs a one-way
userland forwarding rule. Matching T-PDUs have their TEID replaced and leave
through `outgoing`; they are not delivered to `ReadFromGTP`. Use `CloseRelay` to
remove a rule and a separate reverse rule for bidirectional forwarding.
The [relay example](example_test.go) binds the outgoing socket before
starting the incoming receive loop:

```go
func ExampleUPlaneConn_RelayTo() {
	// Bind sockets before starting the receive loops.
	inSocket, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 2152})
	if err != nil {
		fmt.Println(err)
		return
	}
	in, err := gtpv1.NewUPlaneConnFromUDPConn(inSocket)
	if err != nil {
		_ = inSocket.Close()
		fmt.Println(err)
		return
	}
	defer func() { _ = in.Close() }()
	outSocket, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	if err != nil {
		fmt.Println(err)
		return
	}
	out, err := gtpv1.NewUPlaneConnFromUDPConn(outSocket)
	if err != nil {
		_ = outSocket.Close()
		fmt.Println(err)
		return
	}
	defer func() { _ = out.Close() }()
	peer := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 32152}
	if err := in.RelayTo(out, 0x11223344, 0x55667788, peer); err != nil {
		fmt.Println(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := in.ListenAndServe(ctx); err != nil {
		fmt.Println(err)
	}
}
```

### Kernel GTP-U and logging

On Linux, `EnableKernelGTP` enables kernel forwarding with `RoleGGSN` or `RoleSGSN`.
It needs privileges and a kernel with CONFIG_GTP. Add tunnels with `AddTunnel` or
`AddTunnelOverride`; remove them with `DelTunnelByITEI` or `DelTunnelByMSAddress`.
The application configures routes. Kernel mode and userland `RelayTo` cannot be
combined; kernel mode cannot be disabled on a running connection. There is no
TUN/TAP integration. Kernel forwarding is not exercised in CI.

Logging uses `log/slog`; see [ExampleSetLogger](example_test.go):

```go
func ExampleSetLogger() {
	previous := slog.Default()
	gtpv1.SetLogger(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	defer gtpv1.SetLogger(previous)
	// Pass nil to disable this package's background logging.
	gtpv1.SetLogger(nil)
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
| 3 | [VersionNotSupported](message/version-not-supported.go) | `NewVersionNotSupported` |
| 16 | [CreatePDPContextRequest](message/create-pdp-context-req.go) | `NewCreatePDPContextRequest` |
| 17 | [CreatePDPContextResponse](message/create-pdp-context-res.go) | `NewCreatePDPContextResponse` |
| 18 | [UpdatePDPContextRequest](message/update-pdp-context-req.go) | `NewUpdatePDPContextRequest` |
| 19 | [UpdatePDPContextResponse](message/update-pdp-context-res.go) | `NewUpdatePDPContextResponse` |
| 20 | [DeletePDPContextRequest](message/delete-pdp-context-req.go) | `NewDeletePDPContextRequest` |
| 21 | [DeletePDPContextResponse](message/delete-pdp-context-res.go) | `NewDeletePDPContextResponse` |
| 26 | [ErrorIndication](message/error-indication.go) | `NewErrorIndication` |
| 31 | [SupportedExtensionHeaderNotification](message/supported-extension-header-notification.go) | `NewSupportedExtensionHeaderNotification` |
| 254 | [EndMarker](message/end-marker.go) | `NewEndMarker` |
| 255 | [TPDU](message/t-pdu.go) | `NewTPDU` |
| 255 | [TPDU](message/t-pdu.go) | `NewTPDUWithExtentionHeader` |
| 255 | [TPDU](message/t-pdu.go) | `NewTPDUWithSequence` |

## IE support

These source files provide typed constructors returning `*ie.IE`; associated
field structs and accessors are documented in the linked files. Types without
such helpers can use `ie.New` with raw payloads and the generic scalar helpers
in [ie/ie.go](ie/ie.go). This table lists helper availability, not complete
procedure or feature coverage.

| IE family (source) | Typed constructors returning `*IE` |
| --- | --- |
| [apn-restriction](ie/apn-restriction.go) | `NewAPNRestriction` |
| [apn](ie/apn.go) | `NewAccessPointName` |
| [authentication-quintuplet](ie/authentication-quintuplet.go) | `NewAuthenticationQuintuplet` |
| [authentication-triplet](ie/authentication-triplet.go) | `NewAuthenticationTriplet` |
| [cause](ie/cause.go) | `NewCause` |
| [charging-id](ie/charging-id.go) | `NewChargingID` |
| [common-flags](ie/common-flags.go) | `NewCommonFlags` |
| [end-user-address](ie/end-user-address.go) | `NewEndUserAddress`, `NewEndUserAddressByIP`, `NewEndUserAddressIPv4`, `NewEndUserAddressIPv6`, `NewEndUserAddressPPP` |
| [extended-common-flags-ii](ie/extended-common-flags-ii.go) | `NewExtendedCommonFlagsII` |
| [extended-common-flags](ie/extended-common-flags.go) | `NewExtendedCommonFlags` |
| [extension-header-type-list](ie/extension-header-type-list.go) | `NewExtensionHeaderTypeList` |
| [gsn-address](ie/gsn-address.go) | `NewGSNAddress`, `NewGSNAddressByIP` |
| [imei](ie/imei.go) | `NewIMEISV` |
| [imsi](ie/imsi.go) | `NewIMSI` |
| [map-cause](ie/map-cause.go) | `NewMAPCause` |
| [ms-timezone](ie/ms-timezone.go) | `NewMSTimeZone` |
| [ms-validated](ie/ms-validated.go) | `NewMSValidated` |
| [msisdn](ie/msisdn.go) | `NewMSISDN` |
| [nsapi](ie/nsapi.go) | `NewNSAPI` |
| [p-tmsi-signature](ie/p-tmsi-signature.go) | `NewPTMSISignature` |
| [p-tmsi](ie/p-tmsi.go) | `NewPacketTMSI` |
| [pco](ie/pco.go) | `NewProtocolConfigurationOptions` |
| [private-extension](ie/private-extension.go) | `NewPrivateExtension` |
| [qos-profile](ie/qos-profile.go) | `NewQoSProfile` |
| [rai](ie/rai.go) | `NewRouteingAreaIdentity` |
| [ranap-cause](ie/ranap-cause.go) | `NewRANAPCause` |
| [rat-type](ie/rat-type.go) | `NewRATType` |
| [recovery](ie/recovery.go) | `NewRecovery` |
| [reordering-required](ie/reordering-required.go) | `NewReorderingRequired` |
| [selection-mode](ie/selection-mode.go) | `NewSelectionMode` |
| [teardown-ind](ie/teardown-ind.go) | `NewTeardownInd` |
| [teid](ie/teid.go) | `NewTEIDDataI`, `NewTEIDCPlane`, `NewTEIDDataII` |
| [uli-timestamp](ie/uli-timestamp.go) | `NewULITimestamp` |
| [uli](ie/uli.go) | `NewUserLocationInformationWithCGI`, `NewUserLocationInformationWithSAI`, `NewUserLocationInformationWithRAI` |
