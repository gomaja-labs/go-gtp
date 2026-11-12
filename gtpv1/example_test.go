package gtpv1_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv1"
	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
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

func ExampleSetLogger() {
	previous := slog.Default()
	gtpv1.SetLogger(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	defer gtpv1.SetLogger(previous)
	// Pass nil to disable this package's background logging.
	gtpv1.SetLogger(nil)
	// Output:
}
