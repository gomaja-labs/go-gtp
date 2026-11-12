# GW Tester

GW Tester runs a pseudo eNB and MME with S-GW and P-GW programs. It tests GTPv2-C
session establishment; it is not a full radio or mobile core implementation.
The eNB/MME interface uses protobuf/gRPC instead of S1AP. Subscribers in the eNB
configuration attach without authentication. Gateway selection and subscriber
addresses come from YAML configuration rather than DNS or address allocation.

## Build and run

Use Go 1.26 or newer. These programs belong to the separate [examples
module](../go.mod), which builds against the library in this checkout. From
`examples/`, run `go build ./...` and `go vet ./...`.

For the automated attachment check, run `./gw-tester/e2e.sh` from `examples/` on
Linux. It requires Bash, Go, curl, awk and root access (it invokes sudo when
needed). The nodes configure loopback addresses and routes. Run on a disposable
test host or container with network administration privileges.

[e2e.sh](e2e.sh) builds all four nodes, starts S-GW, P-GW, MME and eNB using the
checked-in YAML files, then polls each node's Prometheus endpoint. It succeeds
when every node's active session count equals the number of subscribers in
[enb/enb.yml](enb/enb.yml). `E2E_TIMEOUT` sets the wait in seconds (default 60).
On timeout it prints node logs. On exit it stops the processes and removes its
temporary directory. This check runs in CI. It checks attachment counts, not
packet delivery, HTTP success, detach or kernel GTP-U forwarding.

To run nodes individually, start them in the same order from `examples/`:
`sudo go run ./gw-tester/sgw -config ./gw-tester/sgw/sgw.yml`, then the corresponding
`pgw`, `mme` and `enb` commands with their own config paths. Each program also
accepts `-config` when built as a binary; its default path is `./<node>.yml`.
The checked-in metrics addresses are 127.0.10.1 through 127.0.10.4 on TCP 58080.
Allow the configured gRPC and GTP ports (normally TCP 36412, UDP 2123 and 2152).

## Configuration and limits

| Node | Configuration | Main settings |
| --- | --- | --- |
| eNB | [enb/enb.yml](enb/enb.yml) | Local S1 addresses, MME address, subscribers, incoming TEIDs, traffic source |
| MME | [mme/mme.yml](mme/mme.yml) | S1/S11 addresses, MCC/MNC, APN, gateway addresses |
| S-GW | [sgw/sgw.yml](sgw/sgw.yml) | S11, S1-U, S5-C and S5-U addresses |
| P-GW | [pgw/pgw.yml](pgw/pgw.yml) | S5 addresses, SGi interface/address, subscriber route subnet |

The eNB subscriber `src_ip` is configured locally; the P-GW response does not
assign that local address. An `i_tei` of zero requests random allocation.
`type: external` uses traffic from the configured interface; `type: http_get`
uses the configured URL and source interface. Inspect the YAML files and each
node's `config.go` for the complete configuration surface.

All checked-in configurations set `use_kernel_gtp: false`. The CI attachment
check therefore does not test the kernel data path. For kernel GTP-U traffic,
enable it on the relevant nodes and use Linux with CONFIG_GTP (built in or the
`gtp` module), root/network administration privileges, forwarding enabled and
suitable interfaces/routes. Kernel GTP-U mode is not exercised in CI; successful
attachment does not establish user-plane connectivity. The example eNB and P-GW
do not provide a working user-plane path with kernel mode disabled.

The Detach RPC returns gRPC Unimplemented. A caller must perform a Delete Session
exchange and local cleanup separately. eNB and MME handle SIGHUP for config
reload; S-GW and P-GW do not implement reload. Metrics are exposed when
`prom_addr` is configured; inspect each node's `metrics.go` for exact names.

## S1-MME bindings

From `examples/`, install the buf, protoc-gen-go and protoc-gen-go-grpc tools,
put their installation directory on PATH and run `go generate ./gw-tester/s1mme`.
[s1mme/buf.gen.yaml](s1mme/buf.gen.yaml) selects the Go opaque API and gRPC outputs
for [s1mme/s1mme.proto](s1mme/s1mme.proto).
