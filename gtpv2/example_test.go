package gtpv2_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv2"
	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

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

func ExampleSetLogger() {
	previous := slog.Default()
	gtpv2.SetLogger(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
	defer gtpv2.SetLogger(previous)
	gtpv2.SetLogger(nil)
	// Output:
}
