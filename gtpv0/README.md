# GTPv0

GTPv0 message and IE codecs follow GSM 09.60 V7.10.0. There is no connection
layer or PDP context state machine in this package. Use your own datagram I/O.

## Codecs

Import this version's `ie` and `message` subpackages. This complete function is
[Example_codecs](example_test.go); it builds a message with an IE, marshals it,
parses a datagram, checks the concrete type and reads the IE with error handling.

```go
func Example_codecs() {
	msg := message.NewEchoResponse(1, 0, 0, ie.NewRecovery(0))
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

## Message support

These types have dedicated structs, constructors and parsing dispatch in
[message/message.go](message/message.go). Other types parse as `message.Generic`;
`NewGeneric` constructs one. A named type constant alone is not typed support.

| Type | Typed message | Constructor |
| --- | --- | --- |
| 1 | [EchoRequest](message/echo-req.go) | `NewEchoRequest` |
| 2 | [EchoResponse](message/echo-res.go) | `NewEchoResponse` |
| 16 | [CreatePDPContextRequest](message/create-pdp-context-req.go) | `NewCreatePDPContextRequest` |
| 17 | [CreatePDPContextResponse](message/create-pdp-context-res.go) | `NewCreatePDPContextResponse` |
| 18 | [UpdatePDPContextRequest](message/update-pdp-context-req.go) | `NewUpdatePDPContextRequest` |
| 19 | [UpdatePDPContextResponse](message/update-pdp-context-res.go) | `NewUpdatePDPContextResponse` |
| 20 | [DeletePDPContextRequest](message/delete-pdp-context-req.go) | `NewDeletePDPContextRequest` |
| 21 | [DeletePDPContextResponse](message/delete-pdp-context-res.go) | `NewDeletePDPContextResponse` |
| 255 | [TPDU](message/t-pdu.go) | `NewTPDU` |

## IE support

These source files provide typed constructors returning `*ie.IE`; associated
field structs and accessors are documented in the linked files. Types without
such helpers can use `ie.New` with raw payloads and the generic scalar helpers
in [ie/ie.go](ie/ie.go). This table lists helper availability, not complete
procedure or feature coverage.

| IE family (source) | Typed constructors returning `*IE` |
| --- | --- |
| [apn](ie/apn.go) | `NewAccessPointName` |
| [cause](ie/cause.go) | `NewCause` |
| [charging-gateway-address](ie/charging-gateway-address.go) | `NewChargingGatewayAddress` |
| [charging-id](ie/charging-id.go) | `NewChargingID` |
| [end-user-address](ie/end-user-address.go) | `NewEndUserAddress`, `NewEndUserAddressIPv4`, `NewEndUserAddressIPv6`, `NewEndUserAddressPPP` |
| [flow-label](ie/flow-label.go) | `NewFlowLabelDataI`, `NewFlowLabelSignalling`, `NewFlowLabelDataII` |
| [gsn-address](ie/gsn-address.go) | `NewGSNAddress` |
| [imsi](ie/imsi.go) | `NewIMSI` |
| [ms-not-reachable-reason](ie/ms-not-reachable-reason.go) | `NewMSNotReachableReason` |
| [msisdn](ie/msisdn.go) | `NewMSISDN` |
| [p-tmsi-signature](ie/p-tmsi-signature.go) | `NewPTMSISignature` |
| [p-tmsi](ie/p-tmsi.go) | `NewPacketTMSI` |
| [private-extension](ie/private-extension.go) | `NewPrivateExtension` |
| [qos-profile](ie/qos-profile.go) | `NewQualityOfServiceProfile` |
| [rai](ie/rai.go) | `NewRouteingAreaIdentity` |
| [recovery](ie/recovery.go) | `NewRecovery` |
| [reordering-required](ie/reordering-required.go) | `NewReorderingRequired` |
| [selection-mode](ie/selection-mode.go) | `NewSelectionMode` |
| [tlli](ie/tlli.go) | `NewTemporaryLogicalLinkIdentity` |
