// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv1

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv1/ie"
	"github.com/gomaja-labs/go-gtp/gtpv1/message"
	v2ie "github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type tpduSet struct {
	raddr   net.Addr
	teid    uint32
	seq     uint16
	payload []byte
}

type pktConn interface {
	// WriteToWithDSCPECN writes a packet with payload p to addr using the given DSCP/ECN value.
	// WriteToWithDSCPECN can be made to time out and return
	// an Error with Timeout() == true after a fixed time limit;
	// see SetDeadline and SetWriteDeadline.
	// On packet-oriented connections, write timeouts are rare.
	WriteToWithDSCPECN(p []byte, addr net.Addr, dscpecn int) (n int, err error)

	// File returns a copy of the underlying os.File. It is the caller's responsibility to close f when finished.
	// Closing c does not affect f, and closing f does not affect c.
	// The returned os.File's file descriptor is different from the connection's.
	// Attempting to change properties of the original using this duplicate may or may not have the desired effect.
	File() (f *os.File, err error)

	net.PacketConn
}

type pktConn4 struct {
	// mu is the mutex used before Writing to the PacketConn,
	// to be sure the right DSCP/ECN value
	// is applied before performing the Write.
	mu *sync.Mutex

	// udpConn is the UDPConn used as underlying transport
	udpConn *net.UDPConn

	*ipv4.PacketConn
}

// ReadFrom implements the io.ReaderFrom ReadFrom method.
func (pkt pktConn4) ReadFrom(b []byte) (n int, src net.Addr, err error) {
	n, _, src, err = pkt.PacketConn.ReadFrom(b)
	return n, src, err
}

// WriteTo implements the PacketConn WriteTo method.
func (pkt pktConn4) WriteTo(b []byte, dst net.Addr) (n int, err error) {
	return pkt.PacketConn.WriteTo(b, nil, dst)
}

// setDSCPECN sets the DSCP/ECN value used for next writes.
func (pkt pktConn4) setDSCPECN(dscpecs int) error {
	// With IPv4, DSCP/ECN is using the TOS field
	return pkt.SetTOS(dscpecs)
}

// DSCPECN returns the DSCP/ECN value.
func (pkt pktConn4) DSCPECN() (int, error) {
	// With IPv4, DSCP/ECN is using the TOS field
	return pkt.TOS()
}

// WriteToWithDSCPECN implements the pktConn WriteToWithDSCPECN method.
func (pkt pktConn4) WriteToWithDSCPECN(p []byte, addr net.Addr, dscpecn int) (n int, err error) {
	pkt.mu.Lock()
	defer pkt.mu.Unlock()
	// DSCP/ECN 0 is the socket default, which every marked write restores, so
	// a write with 0 needs no socket option. Plain writes therefore also work
	// where golang.org/x/net cannot read or set the field (IPv6 on Windows,
	// https://github.com/golang/go/issues/7174).
	if dscpecn == 0 {
		return pkt.WriteTo(p, addr)
	}
	oldDSCPECN, err := pkt.DSCPECN()
	if err != nil {
		return 0, err
	}
	// Nothing to set nor restore when the socket already uses the requested
	// value.
	if oldDSCPECN == dscpecn {
		return pkt.WriteTo(p, addr)
	}
	err = pkt.setDSCPECN(dscpecn)
	if err != nil {
		return 0, err
	}
	defer func() {
		// set back DSCP/ECN for next write calls
		_ = pkt.setDSCPECN(oldDSCPECN)
	}()
	return pkt.WriteTo(p, addr)
}

// SyscallConn exposes the socket without changing its nonblocking mode.
func (pkt pktConn4) SyscallConn() (syscall.RawConn, error) {
	return pkt.udpConn.SyscallConn()
}

func (pkt pktConn4) File() (f *os.File, err error) {
	return pkt.udpConn.File()
}

type pktConn6 struct {
	// mu is the mutex used before Writing to the PacketConn,
	// to be sure the right DSCP/ECN value
	// is applied before performing the Write.
	mu *sync.Mutex

	// udpConn is the UDPConn used as underlying transport.
	udpConn *net.UDPConn

	*ipv6.PacketConn
}

// ReadFrom implements the io.ReaderFrom ReadFrom method.
func (pkt pktConn6) ReadFrom(b []byte) (n int, src net.Addr, err error) {
	n, _, src, err = pkt.PacketConn.ReadFrom(b)
	return n, src, err
}

// WriteTo implements the PacketConn WriteTo method.
func (pkt pktConn6) WriteTo(b []byte, dst net.Addr) (n int, err error) {
	return pkt.PacketConn.WriteTo(b, nil, dst)
}

// setDSCPECN sets the DSCP/ECN value used for next writes.
func (pkt pktConn6) setDSCPECN(dscpecs int) error {
	// With IPv6, DSCP/ECN is using the Traffic Class field
	return pkt.SetTrafficClass(dscpecs)
}

// DSCPECN returns the DSCP/ECN value.
func (pkt pktConn6) DSCPECN() (int, error) {
	// With IPv6, DSCP/ECN is using the Traffic Class field
	return pkt.TrafficClass()
}

// WriteToWithDSCPECN implements the pktConn WriteToWithDSCPECN method.
func (pkt pktConn6) WriteToWithDSCPECN(p []byte, addr net.Addr, dscpecn int) (n int, err error) {
	pkt.mu.Lock()
	defer pkt.mu.Unlock()
	// DSCP/ECN 0 is the socket default, which every marked write restores, so
	// a write with 0 needs no socket option. Plain writes therefore also work
	// where golang.org/x/net cannot read or set the field (IPv6 on Windows,
	// https://github.com/golang/go/issues/7174).
	if dscpecn == 0 {
		return pkt.WriteTo(p, addr)
	}
	oldDSCPECN, err := pkt.DSCPECN()
	if err != nil {
		return 0, err
	}
	// Nothing to set nor restore when the socket already uses the requested
	// value.
	if oldDSCPECN == dscpecn {
		return pkt.WriteTo(p, addr)
	}
	err = pkt.setDSCPECN(dscpecn)
	if err != nil {
		return 0, err
	}
	defer func() {
		// set back DSCP/ECN for next write calls
		_ = pkt.setDSCPECN(oldDSCPECN)
	}()
	return pkt.WriteTo(p, addr)
}

// SyscallConn exposes the socket without changing its nonblocking mode.
func (pkt pktConn6) SyscallConn() (syscall.RawConn, error) {
	return pkt.udpConn.SyscallConn()
}

func (pkt pktConn6) File() (f *os.File, err error) {
	return pkt.udpConn.File()
}

// newPktConn creates a new pktConn initialized with a given local UDP address
func newPktConn(laddr net.Addr) (pktConn, error) {
	if laddr == nil {
		return nil, ErrInvalidConnection
	}
	addr, err := net.ResolveUDPAddr(laddr.Network(), laddr.String())
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP(addr.Network(), addr)
	if err != nil {
		return nil, err
	}
	pc, err := newPktConnFromUDPConn(conn, conn.LocalAddr().(*net.UDPAddr).IP)
	if err != nil {
		_ = conn.Close()
	}
	return pc, err
}

// Receiving the destination address of each datagram (IP_PKTINFO and
// IPV6_PKTINFO) lets Error Indication name the address a T-PDU was sent to
// when the socket is bound to a wildcard address (TS 29.281 §7.3.1).
// golang.org/x/net does not implement these control messages on every
// platform; on Windows see https://github.com/golang/go/issues/7175. There
// the connection works without them, and Error Indication falls back to the
// bound address, refusing to name a wildcard one. The variables let tests
// exercise that path.
var (
	receiveDst4 = func(pc *ipv4.PacketConn) error { return pc.SetControlMessage(ipv4.FlagDst, true) }
	receiveDst6 = func(pc *ipv6.PacketConn) error { return pc.SetControlMessage(ipv6.FlagDst, true) }
)

// newPktConnFromUDPConn wraps conn in a pktConn, choosing IPv4 or IPv6 from ip.
func newPktConnFromUDPConn(conn *net.UDPConn, ip net.IP) (pktConn, error) {
	if ip.To4() != nil {
		pc := ipv4.NewPacketConn(conn)
		if err := receiveDst4(pc); err != nil {
			logf("GTP-U destination addresses unavailable on %s: %v", conn.LocalAddr(), err)
		}
		return pktConn4{mu: &sync.Mutex{}, udpConn: conn, PacketConn: pc}, nil
	}
	if ip.To16() != nil || ip == nil {
		pc := ipv6.NewPacketConn(conn)
		if err := receiveDst6(pc); err != nil {
			logf("GTP-U destination addresses unavailable on %s: %v", conn.LocalAddr(), err)
		}
		return pktConn6{mu: &sync.Mutex{}, udpConn: conn, PacketConn: pc}, nil
	}
	return nil, fmt.Errorf("laddr must refer to an IP address")
}

// UPlaneConn represents a U-Plane Connection of GTPv1.
type UPlaneConn struct {
	mu      sync.Mutex
	laddr   net.Addr
	pktConn pktConn
	*msgHandlerMap
	*iteiMap

	tpduCh    chan *tpduSet
	closeCh   chan struct{}
	closeOnce sync.Once
	closeErr  error
	serving   bool

	readDeadline        time.Time
	readDeadlineChanged chan struct{}

	relayMap map[uint32]*peer

	errIndEnabled bool

	// for Linux kernel GTP with netlink
	KernelGTP
}

// KernelGTP consists of the Linux Kernel GTP-U related objects.
type KernelGTP struct {
	enabled bool
	Link    *netlink.GTP
}

// NewUPlaneConn creates a new UPlaneConn used for server. On client side, use DialUPlane instead.
func NewUPlaneConn(laddr net.Addr) *UPlaneConn {
	return &UPlaneConn{
		mu:            sync.Mutex{},
		msgHandlerMap: newDefaultMsgHandlerMap(),
		iteiMap:       newiteiMap(),
		laddr:         laddr,

		tpduCh:              make(chan *tpduSet, 64),
		closeCh:             make(chan struct{}),
		readDeadlineChanged: make(chan struct{}),

		errIndEnabled: true,
	}
}

// NewUPlaneConnFromUDPConn creates a new UPlaneConn used for server, on top of
// an already bound conn instead of a local address. It lets the caller tune the
// socket before serving, e.g. SetReadBuffer/SetWriteBuffer, or options set
// through conn.SyscallConn().
//
// The returned UPlaneConn takes ownership of conn: it is closed when the
// UPlaneConn is closed or the context given to ListenAndServe is done.
func NewUPlaneConnFromUDPConn(conn *net.UDPConn) (*UPlaneConn, error) {
	if conn == nil {
		return nil, ErrInvalidConnection
	}
	laddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return nil, fmt.Errorf("unexpected local address type %T", conn.LocalAddr())
	}
	pc, err := newPktConnFromUDPConn(conn, laddr.IP)
	if err != nil {
		return nil, err
	}

	u := NewUPlaneConn(laddr)
	u.pktConn = pc
	return u, nil
}

// DialUPlane sends Echo Request to raddr to check if the endpoint is alive and returns UPlaneConn.
// ctx bounds the handshake; Close stops the returned connection.
func DialUPlane(ctx context.Context, laddr, raddr net.Addr) (result *UPlaneConn, resultErr error) {
	// Cancellation can close the socket during any handshake operation, including setting deadlines.
	defer func() {
		if resultErr != nil && ctx.Err() != nil {
			resultErr = ctx.Err()
		}
	}()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if raddr == nil {
		return nil, ErrInvalidConnection
	}
	remote, err := net.ResolveUDPAddr(raddr.Network(), raddr.String())
	if err != nil {
		return nil, err
	}
	u := NewUPlaneConn(laddr)
	u.pktConn, err = newPktConn(laddr)
	if err != nil {
		return nil, err
	}
	success := false
	defer func() {
		if !success {
			_ = u.Close()
		}
	}()
	stop := context.AfterFunc(ctx, func() { _ = u.Close() })
	defer stop()
	deadline := time.Now().Add(5 * time.Second)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	if err := u.pktConn.SetReadDeadline(deadline); err != nil {
		return nil, err
	}
	// TS 29.281 §§7.2.1–7.2.2: only the requested peer's matching sequence is a response.
	seq := uint16(generateRandomUint32())
	request := message.NewEchoRequest(seq)
	raw, err := request.Marshal()
	if err != nil {
		return nil, err
	}
	if _, err = u.WriteTo(raw, remote); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	buf := make([]byte, 65535)
	for {
		n, from, err := u.ReadFrom(buf)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, err
		}
		peer, ok := from.(*net.UDPAddr)
		if !ok || peer.Port != remote.Port || !peer.IP.Equal(remote.IP) || peer.Zone != remote.Zone {
			continue
		}
		raw := append([]byte(nil), buf[:n]...)
		msg, err := message.Parse(raw)
		if err != nil {
			continue
		}
		if _, ok := msg.(*message.EchoResponse); !ok || msg.Sequence() != seq {
			continue
		}
		if err := u.pktConn.SetReadDeadline(time.Time{}); err != nil {
			return nil, err
		}
		break
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !stop() {
		return nil, ctx.Err()
	}
	success = true
	u.serving = true
	go func() {
		if err := u.serveLoop(context.Background()); err != nil {
			logf("error on UPlaneConn: %v", err)
		}
	}()
	return u, nil
}

// ListenAndServe binds the GTP-U socket and serves until ctx is done or Close is called.
// This blocks, and returns error only if it face the fatal one. Non-fatal errors are logged
// with the configured logger. See SetLogger for configuration.
func (u *UPlaneConn) ListenAndServe(ctx context.Context) error {
	u.mu.Lock()
	select {
	case <-u.closeCh:
		u.mu.Unlock()
		return net.ErrClosed
	default:
	}
	if u.serving {
		u.mu.Unlock()
		return errors.New("connection already serving")
	}
	if u.pktConn == nil {
		pc, err := newPktConn(u.laddr)
		if err != nil {
			u.mu.Unlock()
			return err
		}
		u.pktConn = pc
	}
	u.mu.Unlock()
	return u.listenAndServe(ctx)
}

func (u *UPlaneConn) listenAndServe(ctx context.Context) error {
	// TODO: this func is left for future enhancement.
	return u.serve(ctx)
}

func (u *UPlaneConn) serve(ctx context.Context) error {
	u.mu.Lock()
	if u.serving {
		u.mu.Unlock()
		return errors.New("connection already serving")
	}
	u.serving = true
	u.mu.Unlock()
	return u.serveLoop(ctx)
}

func (u *UPlaneConn) serveLoop(ctx context.Context) error {
	defer func() { _ = u.Close() }()
	stop := context.AfterFunc(ctx, func() { _ = u.Close() })
	defer stop()
	// TS 29.281 §5.1: a complete UDP datagram can exceed the inner packet MTU.
	buf := make([]byte, 65535)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-u.closed():
			return nil
		default:
		}
		pc := u.packetConn()
		if pc == nil {
			return net.ErrClosed
		}
		var n int
		var raddr net.Addr
		var dst net.IP
		var err error
		if c, ok := pc.(interface {
			readPacket([]byte) (int, net.Addr, net.IP, error)
		}); ok {
			n, raddr, dst, err = c.readPacket(buf)
		} else {
			n, raddr, err = pc.ReadFrom(buf)
		}
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) || ctx.Err() != nil {
				return nil
			}
			var timeout net.Error
			if errors.As(err, &timeout) && timeout.Timeout() {
				continue
			}
			return fmt.Errorf("error reading from UPlaneConn: %w", err)
		}
		// Validate before any indexing, including the relay fast path (TS 29.281 §§5.1, 9.1).
		raw := append([]byte(nil), buf[:n]...)
		msg, err := message.Parse(raw)
		if err != nil {
			logf("error parsing GTP: %v", err)
			continue
		}
		var relay *peer
		if msg.MessageType() == message.MsgTypeTPDU {
			u.mu.Lock()
			relay = u.relayMap[msg.TEID()]
			u.mu.Unlock()
		}
		if handled, err := u.checkExtensionsForRole(raddr, msg, dst, relay != nil); handled {
			if err != nil {
				logf("extension notification: %v", err)
			}
			continue
		}
		if relay != nil {
			raw = relayExtensions(raw, msg.GetHeader())
			binary.BigEndian.PutUint32(raw[4:8], relay.teid)
			if _, err := relay.srcConn.WriteToWithDSCPECN(raw, relay.addr, 0); err != nil {
				logf("error relaying GTP: %v", err)
			}
			continue
		}
		// Dispatch synchronously to preserve receive order and bound resource usage.
		if err := u.handleMessageAt(raddr, msg, dst); err != nil {
			logf("error handling GTP: %v", err)
		}
	}
}

// ReadFrom reads a packet from the connection,
// copying the payload into p. It returns the number of
// bytes copied into p and the return address that
// was on the packet.
// It returns the number of bytes read (0 <= n <= len(p))
// and any error encountered. Callers should always process
// the n > 0 bytes returned before considering the error err.
// ReadFrom reads the raw socket and must not run concurrently with serving.
// Use ReadFromGTP and SetReadDeadline for payload reads while serving.
//
// Note that valid GTP-U packets handled by Kernel can NOT be retrieved by this.
func (u *UPlaneConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	pc := u.packetConn()
	if pc == nil {
		return 0, nil, net.ErrClosed
	}
	return pc.ReadFrom(p)
}

// ReadFromGTP reads a packet from the connection, copying the payload without
// GTP header into p. Up to 64 unread packets are queued; additional packets
// are dropped so signalling and relaying continue. It returns the bytes copied, the return
// address that was on the packet, TEID in the GTP header.
//
// Note that valid GTP-U packets handled by Kernel can NOT be retrieved by this.
func (u *UPlaneConn) ReadFromGTP(p []byte) (n int, addr net.Addr, teid uint32, err error) {
	for {
		deadline, changed := u.readDeadlineState()
		if !deadline.IsZero() && !time.Now().Before(deadline) {
			return 0, nil, 0, os.ErrDeadlineExceeded
		}
		var timer *time.Timer
		var expired <-chan time.Time
		if !deadline.IsZero() {
			timer = time.NewTimer(time.Until(deadline))
			expired = timer.C
		}
		select {
		case <-u.closed():
			err = net.ErrClosed
		case <-changed:
			if timer != nil {
				timer.Stop()
			}
			continue
		case <-expired:
			// Re-read the deadline: it may have been extended concurrently.
			continue
		case tpdu, ok := <-u.tpduCh:
			if !ok {
				err = ErrConnNotOpened
				break
			}
			n = copy(p, tpdu.payload)
			addr = tpdu.raddr
			teid = tpdu.teid
		}
		if timer != nil {
			timer.Stop()
		}
		return
	}
}

func (u *UPlaneConn) readDeadlineState() (time.Time, <-chan struct{}) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.readDeadlineChanged == nil {
		u.readDeadlineChanged = make(chan struct{})
	}
	return u.readDeadline, u.readDeadlineChanged
}

// WriteTo writes a packet with payload p to addr.
// WriteTo can be made to time out and return
// an Error with Timeout() == true after a fixed time limit;
// see SetDeadline and SetWriteDeadline.
// On packet-oriented connections, write timeouts are rare.
func (u *UPlaneConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return u.writePacket(p, addr, nil)
}

// WriteToWithDSCPECN writes a packet with payload p to addr using the given DSCP/ECN value.
// WriteToWithDSCPECN can be made to time out and return
// an Error with Timeout() == true after a fixed time limit;
// see SetDeadline and SetWriteDeadline.
// On packet-oriented connections, write timeouts are rare.
func (u *UPlaneConn) WriteToWithDSCPECN(p []byte, addr net.Addr, dscpecn int) (n int, err error) {
	pc := u.packetConn()
	if pc == nil {
		return 0, net.ErrClosed
	}
	return pc.WriteToWithDSCPECN(p, addr, dscpecn)
}

// WriteToGTP writes a packet with TEID and payload to addr.
func (u *UPlaneConn) WriteToGTP(teid uint32, p []byte, addr net.Addr) (n int, err error) {
	b, err := Encapsulate(teid, p).Marshal()
	if err != nil {
		return
	}

	if _, err = u.WriteTo(b, addr); err != nil {
		return
	}
	return len(b), nil
}

// closed would be used in multiple goroutines.
// never send struct{}{} to it; instead, use close(u.closeCh).
func (u *UPlaneConn) closed() <-chan struct{} {
	return u.closeCh
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (u *UPlaneConn) Close() error {
	u.closeOnce.Do(func() {
		u.mu.Lock()
		close(u.closeCh)
		pc := u.pktConn
		k := u.KernelGTP
		u.mu.Unlock()
		if k.enabled && k.Link != nil {
			u.closeErr = errors.Join(u.closeErr, netlink.LinkDel(k.Link))
		}
		if pc != nil {
			u.closeErr = errors.Join(u.closeErr, pc.Close())
		}
	})
	return u.closeErr
}

// LocalAddr returns the local network address.
func (u *UPlaneConn) LocalAddr() net.Addr {
	pc := u.packetConn()
	if pc == nil {
		return u.laddr
	}
	return pc.LocalAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or
// Write. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (u *UPlaneConn) SetDeadline(t time.Time) error {
	return u.setReadDeadline(t, true)
}

// SetReadDeadline sets the deadline for future and pending ReadFromGTP calls.
// It does not interrupt the serving socket or raw ReadFrom calls.
// A zero value for t means Read will not time out.
func (u *UPlaneConn) SetReadDeadline(t time.Time) error {
	return u.setReadDeadline(t, false)
}

func (u *UPlaneConn) setReadDeadline(t time.Time, both bool) error {
	if u == nil {
		return net.ErrClosed
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	select {
	case <-u.closeCh:
		return net.ErrClosed
	default:
	}
	pc := u.pktConn
	if pc == nil {
		return net.ErrClosed
	}
	// ReadFromGTP owns the application deadline; the serving socket must
	// keep receiving signalling and relay traffic after that deadline expires.
	if both {
		if err := pc.SetWriteDeadline(t); err != nil {
			return err
		}
	}
	u.readDeadline = t
	if u.readDeadlineChanged != nil {
		close(u.readDeadlineChanged)
	}
	u.readDeadlineChanged = make(chan struct{})
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (u *UPlaneConn) SetWriteDeadline(t time.Time) error {
	pc := u.packetConn()
	if pc == nil {
		return net.ErrClosed
	}
	return pc.SetWriteDeadline(t)
}

// AddHandler adds a message handler to *UPlaneConn.
//
// By adding HandlerFuncs, *UPlaneConn (and *Session, *Bearer created by the *UPlaneConn) will handle
// the specified type of message with it's paired HandlerFunc when receiving.
// Messages without registered handlers are just ignored and discarded and the user will
// get ErrNoHandlersFound error.
//
// This should be performed just after creating *UPlaneConn, otherwise the user cannot retrieve
// any values, which is in most cases vital to continue working as a node, from the incoming
// message.
//
// HandlerFuncs for EchoResponse and ErrorIndication are registered by default.
// These HandlerFuncs can be overwritten by specifying message.MsgTypeEchoResponse and/or
// message.MsgTypeErrorIndication as msgType parameter.
func (u *UPlaneConn) AddHandler(msgType uint8, fn HandlerFunc) {
	u.msgHandlerMap.store(msgType, fn)
}

// AddHandlers adds multiple handler funcs at a time.
//
// See AddHandler for detailed usage.
func (u *UPlaneConn) AddHandlers(funcs map[uint8]HandlerFunc) {
	for msgType, fn := range funcs {
		u.msgHandlerMap.store(msgType, fn)
	}
}

func (u *UPlaneConn) handleMessage(senderAddr net.Addr, msg message.Message) error {
	return u.handleMessageAt(senderAddr, msg, nil)
}

// EchoRequest sends a EchoRequest.
func (u *UPlaneConn) EchoRequest(raddr net.Addr) error {
	b, err := message.NewEchoRequest(0).Marshal()
	if err != nil {
		return err
	}

	if _, err := u.WriteTo(b, raddr); err != nil {
		return err
	}
	return nil
}

// EchoResponse sends a EchoResponse.
func (u *UPlaneConn) EchoResponse(raddr net.Addr) error {
	b, err := message.NewEchoResponse(0, ie.NewRecovery(0)).Marshal()
	if err != nil {
		return err
	}

	if _, err := u.WriteTo(b, raddr); err != nil {
		return err
	}
	return nil
}

// ErrorIndication just sends ErrorIndication message.
func (u *UPlaneConn) ErrorIndication(raddr net.Addr, received message.Message, destination ...net.IP) error {
	// TS 29.281 §7.3.1: identify the triggering datagram's destination, never a wildcard.
	var ip net.IP
	if len(destination) > 0 {
		ip = destination[0]
	}
	if ip == nil {
		host, _, err := net.SplitHostPort(u.LocalAddr().String())
		if err != nil {
			return err
		}
		ip = net.ParseIP(host)
	}
	if ip == nil || ip.IsUnspecified() {
		return fmt.Errorf("GTP-U destination address unavailable: %w", ErrInvalidConnection)
	}
	// TS 29.281 §4.4.2.4: Error Indication always targets port 2152.
	addr, err := userPlaneAddr(raddr)
	if err != nil {
		return err
	}

	errInd, err := message.NewErrorIndication(
		0, received.Sequence(),
		ie.NewTEIDDataI(received.TEID()),
		ie.NewGSNAddressByIP(ip),
	).Marshal()
	if err != nil {
		return err
	}

	if _, err := u.writePacket(errInd, addr, ip); err != nil {
		return err
	}
	return nil
}

// RespondTo sends a message(specified with "toBeSent" param) in response to
// a message(specified with "received" param).
//
// This is to make it easier to handle SequenceNumber.
func (u *UPlaneConn) RespondTo(raddr net.Addr, received, toBeSent message.Message) error {
	toBeSent.SetSequenceNumber(received.Sequence())
	b := make([]byte, toBeSent.MarshalLen())
	if err := toBeSent.MarshalTo(b); err != nil {
		return err
	}

	if _, err := u.WriteTo(b, raddr); err != nil {
		return err
	}
	return nil
}

// Restarts returns the number of restarts in uint8.
func (u *UPlaneConn) Restarts() uint8 {
	return 0
}

// NewFTEID creates a new GTPv2 F-TEID with random TEID value that is unique within UPlaneConn.
// To ensure the uniqueness, don't create in the other way if you once use this method.
// This is meant to be used for creating F-TEID IE for non-local interface type, such as
// the ones that are used in U-Plane. For local interface, use (*Conn).NewSenderFTEID instead.
//
// Note that in the case there's a lot of Session on the Conn, it may take a long
// time to find a new unique value.
//
// TODO: optimize performance...
func (u *UPlaneConn) NewFTEID(ifType uint8, v4, v6 string) (fteidIE *v2ie.IE) {
	var teid uint32
	for try := uint32(0); try < 0xffff; try++ {
		const logEvery = 0xff
		if try&logEvery == logEvery {
			logf("Generating NewSenderFTEID crossed tries:%d", try)
		}

		t := generateRandomUint32()
		if t == 0 {
			continue
		}

		// Try to mark TEID as taken. Fails if something exists
		if ok := u.iteiMap.tryStore(t, time.Now()); !ok {
			logf("TEID-U: %#08x has already been taken, trying to generate another one...", t)
			continue
		}

		teid = t
		break
	}

	if teid == 0 {
		return nil
	}
	fteidIE = v2ie.NewFullyQualifiedTEID(ifType, teid, v4, v6)
	if fteidIE.Err() != nil {
		u.iteiMap.delete(teid)
	}
	return fteidIE
}

func generateRandomUint32() uint32 {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return 0
	}

	return binary.BigEndian.Uint32(b)
}

type iteiMap struct {
	syncMap sync.Map
}

func newiteiMap() *iteiMap {
	return &iteiMap{}
}

func (t *iteiMap) tryStore(itei uint32, ts time.Time) bool {
	_, loaded := t.syncMap.LoadOrStore(itei, ts)
	return !loaded
}

func (t *iteiMap) delete(itei uint32) {
	t.syncMap.Delete(itei)
}

// EnableErrorIndication re-enables automatic sending of
// Error Indication to unknown messages, which is enabled by
// default.
//
// See also: DisableErrorIndication.
func (u *UPlaneConn) EnableErrorIndication() {
	u.mu.Lock()
	u.errIndEnabled = true
	u.mu.Unlock()
}

// DisableErrorIndication makes default T-PDU handler stop
// responding with Error Indication in case of receiving T-PDU
// with unknown TEID.
//
// When disabled, it passes the unhandled T-PDU to user who calls
// ReadFromGTP instead.
func (u *UPlaneConn) DisableErrorIndication() {
	u.mu.Lock()
	u.errIndEnabled = false
	u.mu.Unlock()
}

func (pkt pktConn4) readPacket(b []byte) (int, net.Addr, net.IP, error) {
	n, cm, src, err := pkt.PacketConn.ReadFrom(b)
	var dst net.IP
	if cm != nil {
		dst = cm.Dst
	}
	return n, src, dst, err
}
func (pkt pktConn6) readPacket(b []byte) (int, net.Addr, net.IP, error) {
	n, cm, src, err := pkt.PacketConn.ReadFrom(b)
	var dst net.IP
	if cm != nil {
		dst = cm.Dst
	}
	return n, src, dst, err
}
func (pkt pktConn4) writePacket(b []byte, dst net.Addr, src net.IP) (int, error) {
	pkt.mu.Lock()
	defer pkt.mu.Unlock()
	return pkt.PacketConn.WriteTo(b, &ipv4.ControlMessage{Src: src}, dst)
}
func (pkt pktConn6) writePacket(b []byte, dst net.Addr, src net.IP) (int, error) {
	pkt.mu.Lock()
	defer pkt.mu.Unlock()
	return pkt.PacketConn.WriteTo(b, &ipv6.ControlMessage{Src: src}, dst)
}

func (u *UPlaneConn) packetConn() pktConn { u.mu.Lock(); defer u.mu.Unlock(); return u.pktConn }

func (u *UPlaneConn) writePacket(b []byte, dst net.Addr, src net.IP) (int, error) {
	pc := u.packetConn()
	if pc == nil {
		return 0, net.ErrClosed
	}
	if src != nil {
		if c, ok := pc.(interface {
			writePacket([]byte, net.Addr, net.IP) (int, error)
		}); ok {
			return c.writePacket(b, dst, src)
		}
	}
	return pc.WriteToWithDSCPECN(b, dst, 0)
}

func userPlaneAddr(addr net.Addr) (*net.UDPAddr, error) {
	if addr == nil {
		return nil, ErrInvalidConnection
	}
	dst, err := net.ResolveUDPAddr(addr.Network(), addr.String())
	if err != nil {
		return nil, err
	}
	dst.Port = 2152
	return dst, nil
}

func (u *UPlaneConn) handleMessageAt(senderAddr net.Addr, msg message.Message, dst net.IP) error {
	if handled, err := u.checkExtensions(senderAddr, msg, dst); handled {
		return err
	}
	handle, ok := u.msgHandlerMap.load(msg.MessageType())
	if !ok {
		switch msg.MessageType() {
		case message.MsgTypeTPDU:
			return handleTPDUAt(u, senderAddr, msg, dst)
		case message.MsgTypeEchoRequest:
			if dst == nil {
				return handleEchoRequest(u, senderAddr, msg)
			}
			// TS 29.281 §§4.4.3.2, 7.2.2: echo the sequence and reply from the request's destination.
			reply := message.NewEchoResponse(msg.Sequence(), ie.NewRecovery(0))
			b, err := reply.Marshal()
			if err != nil {
				return err
			}
			_, err = u.writePacket(b, senderAddr, dst)
			return err
		}
		return &HandlerNotFoundError{MsgType: msg.MessageTypeName()}
	}
	if err := handle(u, senderAddr, msg); err != nil {
		return fmt.Errorf("failed to handle %s: %w", msg.MessageTypeName(), err)
	}
	return nil
}

func (u *UPlaneConn) checkExtensions(sender net.Addr, msg message.Message, dst net.IP) (bool, error) {
	return u.checkExtensionsForRole(sender, msg, dst, false)
}

func (u *UPlaneConn) checkExtensionsForRole(sender net.Addr, msg message.Message, dst net.IP, intermediate bool) (bool, error) {
	// TS 29.281 §5.2.1: unknown comprehension-required extensions are discarded and notified.
	header := msg.GetHeader()
	if header == nil {
		return true, message.ErrInvalidLength
	}
	for _, e := range header.ExtensionHeaders {
		switch e.Type {
		case message.ExtHeaderTypeUDPPort, message.ExtHeaderTypePDCPPDUNumber, message.ExtHeaderTypeLongPDCPPDUNumber, message.ExtHeaderTypeLongPDCPPDUNumberRequired, message.ExtHeaderTypeServiceClassIndicator, message.ExtHeaderTypePDUSessionContainer:
			continue
		}
		if !e.IsComprehensionRequired() {
			continue
		}
		// TS 29.281 §5.2.1 Figure 5.2.1-2: 10xxxxxx is required
		// only at the endpoint; intermediate nodes forward it unchanged.
		if intermediate && e.Type&0xc0 == 0x80 {
			continue
		}
		// TS 29.281 §5.2.1: only requests and G-PDUs trigger notification;
		// all other messages with an unknown required extension are discarded.
		if msg.MessageType() != message.MsgTypeEchoRequest && msg.MessageType() != message.MsgTypeTPDU {
			return true, nil
		}
		addr, err := userPlaneAddr(sender)
		if err != nil {
			return true, err
		}
		notification := message.NewSupportedExtensionHeaderNotification(0, 0, ie.NewExtensionHeaderTypeList(0x40, 0xc0, 0x03, 0x82, 0x20, 0x85))
		b, err := notification.Marshal()
		if err != nil {
			return true, err
		}
		_, err = u.writePacket(b, addr, dst)
		return true, err
	}
	return false, nil
}

// relayExtensions removes 01xxxxxx headers (TS 29.281 §5.2.1 Figure
// 5.2.1-2), preserving other contents byte for byte, including future fields.
// raw and h must have passed message.Parse's framing checks.
func relayExtensions(raw []byte, h *message.Header) []byte {
	if len(h.ExtensionHeaders) == 0 {
		return raw
	}
	out := append([]byte(nil), raw[:12]...)
	previousNext := 11
	off := 12
	for _, e := range h.ExtensionHeaders {
		n := e.MarshalLen()
		if e.Type&0xc0 != 0x40 {
			out[previousNext] = e.Type
			out = append(out, raw[off:off+n]...)
			previousNext = len(out) - 1
		}
		off += n
	}
	out[previousNext] = 0
	if previousNext == 11 {
		out[0] &^= 4
		if out[0]&3 == 0 {
			out = out[:8]
		}
	}
	out = append(out, raw[off:]...)
	binary.BigEndian.PutUint16(out[2:4], uint16(len(out)-8))
	return out
}
