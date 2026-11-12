module github.com/gomaja-labs/go-gtp/examples

go 1.26.0

require (
	github.com/gomaja-labs/go-gtp v0.0.0
	github.com/prometheus/client_golang v1.25.0
	github.com/vishvananda/netlink v1.3.1
	go.yaml.in/yaml/v3 v3.0.6
	google.golang.org/grpc v1.85.0
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pandatix/go-abnf v0.5.5 // indirect
	github.com/prometheus/client_model v0.6.3 // indirect
	github.com/prometheus/common v0.73.1 // indirect
	github.com/prometheus/procfs v0.23.0 // indirect
	github.com/vishvananda/netns v0.0.5 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.44.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261110000943-09e21b237f2c // indirect
)

// The examples build against the library in the same commit.
replace github.com/gomaja-labs/go-gtp => ../
