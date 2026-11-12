package gtp_test

import (
	"fmt"

	"github.com/gomaja-labs/go-gtp"
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

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
