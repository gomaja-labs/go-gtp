# go-gtp

[![ci](https://github.com/gomaja-labs/go-gtp/actions/workflows/ci.yml/badge.svg)](https://github.com/gomaja-labs/go-gtp/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/gomaja-labs/go-gtp.svg)](https://pkg.go.dev/github.com/gomaja-labs/go-gtp)

Go codecs and UDP connection APIs for the GPRS Tunneling Protocol.

| Package | Codecs | Connection layer |
| --- | --- | --- |
| [gtpv0](gtpv0/README.md) | GTPv0 messages and IEs | None |
| [gtpv1](gtpv1/README.md) | GTPv1-C and GTPv1-U messages, IEs and extension headers | GTP-U `UPlaneConn`, userland relay, optional Linux kernel GTP-U |
| [gtpv2](gtpv2/README.md) | GTPv2-C messages and IEs, including grouped IEs and piggybacking | `Conn`, handlers, transactions, sessions and bearers |

Support tables in each version README list the typed helpers in this tree.
Generic messages and raw IEs cover types without dedicated helpers; they do not
supply procedure validation for those types. GTP prime is not implemented.

## Requirements and installation

Use Go 1.26 or newer. Install with `go get github.com/gomaja-labs/go-gtp@latest`.
The examples have a [separate module](examples/go.mod) whose replace directive
builds against the library in the same checkout. Run their build commands from
`examples/`, rather than installing the examples through the library module.

## Reference specifications

The codecs follow these editions; support for a typed message or IE does not
mean that every procedure or optional feature in its specification is implemented.

| Specification | Edition | Subject |
| --- | --- | --- |
| 3GPP TS 29.274 | V20.0.0 | GTPv2-C |
| 3GPP TS 29.060 | V19.0.0 | GTPv1-C |
| 3GPP TS 29.281 | V19.2.0 | GTPv1-U |
| 3GPP TS 38.415 | V19.1.0 | PDU Session Container |
| 3GPP TS 24.008 | V20.1.0 | Shared information element formats |
| 3GPP TS 23.003 | V20.1.0 | Identifiers and addressing |
| GSM 09.60 | V7.10.0 | GTPv0 |

## Building and parsing

Import the root package as `gtp`, and the GTPv2 subpackages as `ie` and `message`.
This is [ExampleParse](example_test.go); `gtp.Parse` dispatches by wire version.
When the version is known, use its `message.Parse` instead. The version READMEs
show those calls and safe IE access.

```go
func ExampleParse() {
	msg := message.NewEchoRequest(1, ie.NewRecovery(0))
	wire, err := gtp.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	parsed, err := gtp.Parse(wire)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(parsed.Version(), parsed.MessageTypeName())
	// Output: 2 Echo Request
}
```

## API conventions

- `Marshal` allocates; `MarshalTo` writes into caller storage and returns an error
  for invalid values or insufficient space. Allocate using `MarshalLen` and check
  the marshal error. Constructors returning `*ie.IE` carry construction failures:
  call `Err()` before use. Marshaling propagates these errors through messages.
- Typed IE accessors return errors; a nil IE returns `ie.ErrIENotFound`. `Must*`
  and `Has*` helpers return zero values on error and discard that distinction.
- Decoding checks lengths and framing, with sentinel and typed errors inspectable
  through `errors.Is` and `errors.As`. Content handling follows each message's
  rules: malformed optional or unexpected IEs may be retained and ignored, while
  an erroneous mandatory IE can return both a message and an error. Check errors
  before using values. The version message package documentation describes the
  exceptions and response policy; decoding alone is not procedure validation.
- Parsed payloads and IE values can reference the input buffer. Copy the datagram
  before parsing if the result must survive reuse or modification of that buffer.
  Handlers must copy retained data before returning. GTPv2 piggybacked messages
  and `GetHeader().PiggybackError` need independent inspection.
- Repeated IEs have slice fields where modeled. Unknown or unexpected IEs remain
  in `AdditionalIEs`; decoding into a reused message resets its fields.
- Handlers run on the serving goroutine in receive order. Keep them short; do not
  wait inside one for a response that the same receive loop must dispatch. Send
  replies with `RespondTo` to preserve the request sequence. Register handlers
  before serving; `Dial` already starts serving after its Echo handshake.
- `gtpv1.SetLogger` and `gtpv2.SetLogger` accept `*slog.Logger`; nil disables the
  respective package's background logging. See their compiled examples.

There is no GTPv1-C connection layer: it has codecs only. The library does not
provide TUN/TAP integration, a complete mobile core, or session persistence.
Applications own routing, procedure-specific mandatory/conditional IE checks,
peer authorization and durable state. Kernel GTP-U is Linux only and requires
kernel support and privileges; userland GTP-U does not require kernel GTP.

## Migrating from wmnsk/go-gtp

The API targets the latest main and does not retain compatibility wrappers.

- Change imports to `github.com/gomaja-labs/go-gtp` (and its subpackages), and use
  Go 1.26 or newer. Build the standalone programs from the `examples/` module.
- Replace removed `Serialize`, `SerializeTo`, `Decode`, `DecodeFromBytes` and
  `Len` wrappers with `Marshal`, `MarshalTo`, `Parse`, `UnmarshalBinary` and
  `MarshalLen`. Check errors from encoding and sequence setters. Message
  interfaces no longer include the deprecated methods.
- Replace `gtpv1.Relay` with `UPlaneConn.RelayTo`. Use `DialUPlane` for GTP-U.
  GTPv2 handlers take `*gtpv2.Conn`; there is no GTPv1-C `Conn` implementation.
- Replace the old logger interface with `SetLogger(*slog.Logger)` in each
  connection package; pass nil to disable logging.
- Check IE `Err()` and accessor errors. Invalid digits, addresses, lengths and
  overflowing fields now fail instead of being silently truncated. Shared
  TBCD/BCD and FQDN helpers return errors; use `utils.DecodePLMN` instead of the
  removed `DecodeMCC`/`DecodeMNC` helpers.
- Build GTPv2 ULI using `NewUserLocationInformation` with typed location structs.
  The old positional/lazy constructors, `NewUserLocationInformationStruct` and
  `UserLocationInfo` alias are removed. `NewRAI` takes a `uint8` RAC. The removed
  ARP aliases are replaced by `HasPCI`/`HasPVI`.
- GTPv1 `NewTEIDDataII` now takes NSAPI and TEID. GTPv2 `NewRFSPIndex` and its
  accessors use `uint16`; `NewIntegerNumber` takes a `uint32` value and an explicit
  width. Its accessors return `uint32`. Extended PCO constructors/accessors take
  a `PCODirection`; secondary RAT usage reports use typed fields. Use
  `TimestampNear` when resolving a timestamp's NTP era against a known time.
- Update callers of retyped message fields: repeated monitoring, load and
  overload IEs use slices. In bearer responses, use
  `MMESGSNIdentifierOrUELocalIPAddress`; its IP Address instance 0 meaning
  depends on the interface. Inspect `AdditionalIEs` for unmodeled pairs.

## Programs and validation

[GW Tester](examples/gw-tester/README.md) runs a pseudo eNB/MME and S-GW/P-GW.
Its CI check verifies attachment counts; it does not exercise kernel GTP-U.
Other standalone programs are under [examples/](examples/).

Run `go test ./... -run Example -count=1` for deterministic documentation examples,
then `go test ./... -count=1` for the full library suite. Network examples without
an Output comment are compiled but not executed by `go test`. In `examples/`,
run `go build ./...` and `go vet ./...` to check the standalone programs.

## See also

[go-pfcp](https://github.com/gomaja-labs/go-pfcp), a PFCP implementation in Go.

## License

[MIT](LICENSE). See the license for copyright holders.
