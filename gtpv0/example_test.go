package gtpv0_test

import (
	"fmt"

	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
	"github.com/gomaja-labs/go-gtp/gtpv0/message"
)

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
