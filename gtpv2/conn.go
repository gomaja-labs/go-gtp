// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv2

import (
	"container/list"
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

// Conn represents a GTPv2-C connection.
//
// Conn provides the automatic handling of message by adding handlers to it with
// AddHandler(s). See AddHandler for detailed usage.
//
// Conn also provides the functions to manage Sessions/Bearers that works over the
// connection(=between a node to another).
// See the docs of CreateSession, AddSession, DeleteSession methods for details.
type Conn struct {
	mu                 sync.Mutex
	laddr              net.Addr
	pktConn            net.PacketConn
	sessionsByTEID     map[uint32]*Session
	sessions           map[*Session]struct{}
	reservedTEIDs      map[uint32]struct{}
	transactions       map[transactionKey]*transaction
	received           map[requestKey]*receivedRequest
	replayPeers        map[string]*replayPeer
	replayPeerOrder    list.List
	replayWindow       time.Duration
	retransmitInterval time.Duration
	retransmitRetries  int
	errorHandler       func(error)
	serving            bool
	localIfType        uint8

	validationEnabled bool

	closeCh chan struct{}
	*msgHandlerMap

	// sequence is the last SequenceNumber used in the request.
	//
	// TS29.274 7.6  Reliable Delivery of Signalling Messages;
	// The Sequence Number shall be unique for each outstanding Initial message sourced
	// from the same IP/UDP endpoint(=Conn).
	sequence uint32

	// RestartCounter is the RestartCounter value in Recovery IE, which represents how many
	// times the GTPv2-C endpoint is restarted.
	RestartCounter uint8
}

// NewConn creates a new Conn used for server. On client side, use Dial instead.
func NewConn(laddr net.Addr, localIfType, counter uint8) *Conn {
	return &Conn{
		mu:                 sync.Mutex{},
		laddr:              laddr,
		sessionsByTEID:     make(map[uint32]*Session),
		sessions:           make(map[*Session]struct{}),
		reservedTEIDs:      make(map[uint32]struct{}),
		transactions:       make(map[transactionKey]*transaction),
		received:           make(map[requestKey]*receivedRequest),
		replayPeers:        make(map[string]*replayPeer),
		replayWindow:       3 * time.Minute,
		retransmitInterval: 3 * time.Second,
		retransmitRetries:  3,
		localIfType:        localIfType,
		validationEnabled:  true,
		closeCh:            make(chan struct{}),
		msgHandlerMap:      newDefaultMsgHandlerMap(),
		sequence:           0,
		RestartCounter:     counter,
	}
}

// Dial sends Echo Request to raddr to check if the endpoint is alive and returns Conn.
// ctx bounds only the handshake. Close the returned connection to stop serving.
//
// It does not bind the raddr to the underlying connection, which enables a Conn to
// send to/receive from multiple peers with single laddr.
//
// If Echo exchange is unnecessary, use NewConn and ListenAndServe instead.
func Dial(ctx context.Context, laddr, raddr net.Addr, localIfType, counter uint8) (_ *Conn, err error) {
	c := NewConn(laddr, localIfType, counter)
	if err = c.Listen(ctx); err != nil {
		return nil, err
	}
	defer func() {
		if err != nil && ctx.Err() != nil {
			err = ctx.Err()
		}
	}()
	success := false
	defer func() {
		if !success {
			_ = c.Close()
		}
	}()
	// Closing the socket interrupts the handshake immediately on cancellation.
	stop := context.AfterFunc(ctx, func() { _ = c.Close() })
	defer stop()
	seq, err := c.EchoRequest(raddr)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(3 * time.Second)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	if err = c.SetReadDeadline(deadline); err != nil {
		return nil, err
	}
	buf := make([]byte, 65535)
	for {
		n, peer, readErr := c.ReadFrom(buf)
		if readErr != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, readErr
		}
		msg, parseErr := message.Parse(buf[:n])
		// TS 29.274 §7.6: only the matching Echo Response establishes liveness.
		if parseErr != nil || !sameEndpoint(peer, raddr) || msg.MessageType() != message.MsgTypeEchoResponse || msg.Sequence() != seq {
			continue
		}
		if err = c.handleMessage(peer, msg); err != nil {
			return nil, err
		}
		break
	}
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	if err = c.SetReadDeadline(time.Time{}); err != nil {
		return nil, err
	}
	if !stop() {
		return nil, ctx.Err()
	}
	success = true
	go func() {
		if serveErr := c.Serve(context.Background()); serveErr != nil {
			c.reportError(serveErr)
		}
	}()
	return c, nil
}

// ListenAndServe creates a new GTPv2-C Conn and start serving background.
func (c *Conn) ListenAndServe(ctx context.Context) error {
	err := c.Listen(ctx)
	if err != nil {
		return err
	}
	return c.listenAndServe(ctx)
}

// Listen creates a new GTPv2-C Conn
func (c *Conn) Listen(ctx context.Context) error {
	if c == nil {
		return net.ErrClosed
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.closeCh:
		return net.ErrClosed
	default:
	}
	if c.pktConn != nil {
		return errors.New("connection already listening")
	}
	if c.laddr == nil {
		return errors.New("nil local address")
	}
	pc, err := net.ListenPacket(c.laddr.Network(), c.laddr.String())
	if err != nil {
		return err
	}
	c.pktConn = pc
	return nil
}

func (c *Conn) listenAndServe(ctx context.Context) error {
	// TODO: this func is left for future enhancement.
	return c.Serve(ctx)
}

func (c *Conn) closed() <-chan struct{} {
	return c.closeCh
}

// Serve starts serving GTPv2 connection. Handlers execute serially on this
// goroutine and must return promptly; blocking handlers delay every peer.
func (c *Conn) Serve(ctx context.Context) error {
	if c == nil {
		return net.ErrClosed
	}
	c.mu.Lock()
	if c.serving || c.pktConn == nil {
		c.mu.Unlock()
		return errors.New("connection is not listening or already serving")
	}
	c.serving = true
	pc := c.pktConn
	c.mu.Unlock()
	stop := context.AfterFunc(ctx, func() { _ = c.Close() })
	defer stop()
	defer func() { _ = c.Close() }()
	// One buffer covers any UDP datagram. Handlers run in receive order on this
	// goroutine and must not block for long (TS 29.274 §7.7.0).
	buf := make([]byte, 65535)
	for {
		n, peer, err := pc.ReadFrom(buf)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return nil
			}
			return fmt.Errorf("error reading from Conn %s: %w", c.LocalAddr(), err)
		}
		c.processDatagram(peer, buf[:n])
	}
}

// ReadFrom reads a packet from the connection,
// copying the payload into p. It returns the number of
// bytes copied into p and the return address that
// was on the packet.
// It returns the number of bytes read (0 <= n <= len(p))
// and any error encountered. Callers should always process
// the n > 0 bytes returned before considering the error err.
// ReadFrom can be made to time out and return
// an Error with Timeout() == true after a fixed time limit;
// see SetDeadline and SetReadDeadline.
func (c *Conn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	pc := c.packetConn()
	if pc == nil {
		return 0, nil, net.ErrClosed
	}
	return pc.ReadFrom(p)
}

// WriteTo sends raw bytes without registering a transaction. Replies to raw
// requests are discarded unless sent through SendRawTo or SendMessageTo.
// WriteTo writes a packet with payload p to addr.
// WriteTo can be made to time out and return
// an Error with Timeout() == true after a fixed time limit;
// see SetDeadline and SetWriteDeadline.
// On packet-oriented connections, write timeouts are rare.
func (c *Conn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	pc := c.packetConn()
	if pc == nil {
		return 0, net.ErrClosed
	}
	return pc.WriteTo(p, addr)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *Conn) Close() error {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	select {
	case <-c.closeCh:
		c.mu.Unlock()
		return nil
	default:
		close(c.closeCh)
	}
	for key, tx := range c.transactions {
		tx.timer.Stop()
		delete(c.transactions, key)
	}
	clear(c.received)
	clear(c.replayPeers)
	c.replayPeerOrder.Init()
	pc := c.pktConn
	c.mu.Unlock()
	if pc != nil {
		return pc.Close()
	}
	return nil
}

func (c *Conn) packetConn() net.PacketConn {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pktConn
}

// LocalAddr returns the bound address, or the configured address before Listen.
func (c *Conn) LocalAddr() net.Addr {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.pktConn == nil {
		return c.laddr
	}
	return c.pktConn.LocalAddr()
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
func (c *Conn) SetDeadline(t time.Time) error {
	pc := c.packetConn()
	if pc == nil {
		return net.ErrClosed
	}
	return pc.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *Conn) SetReadDeadline(t time.Time) error {
	pc := c.packetConn()
	if pc == nil {
		return net.ErrClosed
	}
	return pc.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	pc := c.packetConn()
	if pc == nil {
		return net.ErrClosed
	}
	return pc.SetWriteDeadline(t)
}

// AddHandler adds a message handler to Conn.
//
// By adding HandlerFunc, Conn(and Session, Bearer created over the Conn) will handle
// the specified type of message with it's paired HandlerFunc when receiving.
// Messages without registered handlers are just ignored and logged.
//
// This should be performed just after creating Conn, otherwise the user cannot retrieve
// any values, which is in most cases vital to continue working as a node, from the incoming
// message.
//
// The error returned from handler is just logged. Any important due should be done inside
// the HandlerFunc before returning. This behavior might change in the future.
//
// HandlerFunc for EchoResponse and VersionNotSupportedIndication are registered by default.
// These HandlerFunc can be overridden by specifying message.MsgTypeEchoResponse and/or
// message.MsgTypeVersionNotSupportedIndication as msgType parameter.
func (c *Conn) AddHandler(msgType uint8, fn HandlerFunc) {
	if c == nil {
		return
	}
	c.msgHandlerMap.store(msgType, fn)
}

// AddHandlers adds multiple handler funcs at a time, using a map.
// The key of the map is message type of the GTPv2-C message. You can use MsgTypeFooBar
// constants defined in this package as well as any raw uint8 values.
//
// See AddHandler for how the given handlers behave.
func (c *Conn) AddHandlers(funcs map[uint8]HandlerFunc) {
	if c == nil {
		return
	}
	for msgType, fn := range funcs {
		c.msgHandlerMap.store(msgType, fn)
	}
}

func (c *Conn) handleMessage(senderAddr net.Addr, msg message.Message) error {
	_, err := c.dispatchMessage(senderAddr, msg)
	return err
}

// dispatchMessage distinguishes a discarded primary from a successfully handled one.
func (c *Conn) dispatchMessage(senderAddr net.Addr, msg message.Message) (bool, error) {
	if msg == nil || msg.GetHeader() == nil {
		return false, message.ErrNilMessage
	}
	if !c.knownType(msg.MessageType()) {
		return false, nil
	}
	c.mu.Lock()
	validationEnabled := c.validationEnabled
	c.mu.Unlock()
	if c.isDuplicateOrUnmatched(senderAddr, msg) {
		return false, nil
	}
	if validationEnabled {
		if err := c.validate(senderAddr, msg); err != nil {
			if isReply(msg.MessageType()) {
				c.reportError(&TransactionError{senderAddr, msg.Sequence(), msg.MessageType(), err})
			}
			return false, fmt.Errorf("failed to validate %s: %w", msg.MessageTypeName(), err)
		}
	}

	if validationEnabled {
		if err := c.validateMandatory(senderAddr, msg); err != nil {
			return false, err
		}
	}
	handle, ok := c.msgHandlerMap.load(msg.MessageType())
	if !ok || handle == nil {
		return false, &HandlerNotFoundError{MsgType: msg.MessageTypeName()}
	}

	if err := handle(c, senderAddr, msg); err != nil {
		return false, fmt.Errorf("failed to handle %s: %w", msg.MessageTypeName(), err)
	}

	return true, nil
}

// EnableValidation turns on automatic validation of incoming message.
// This is expected to be used only after DisableValidation() is used, as the validation
// is enabled by default.
//
// Conn checks if;
//
// GTP Version is 2
// TEID is known to Conn
// Mandatory IEs are present according to the modeled message tables.
//
// Even the validation is failed, it does not return error to user. Instead, it just logs
// and discards the packets so that the HandlerFunc won't get the invalid message.
// Extra validations should be done in HandlerFunc.
func (c *Conn) EnableValidation() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.validationEnabled = true
}

// DisableValidation turns off automatic validation of incoming message.
// It is not recommended to use this except the node is in debugging mode.
//
// See EnableValidation for what are validated.
func (c *Conn) DisableValidation() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.validationEnabled = false
}

func (c *Conn) validate(senderAddr net.Addr, msg message.Message) error {
	// check GTP version
	if msg.Version() != 2 {
		if msg.Version() > 2 {
			if err := c.VersionNotSupportedIndication(senderAddr, msg); err != nil {
				return fmt.Errorf("failed to respond with VersionNotSupportedIndication: %w", err)
			}
		}
		return fmt.Errorf("received an invalid version(%d) of message: %v", msg.Version(), msg)
	}

	// check if TEID is known or not
	if teid := msg.TEID(); teid != 0 {
		if _, err := c.GetSessionByTEID(teid); err != nil {
			if responseErr := c.respondWithCause(senderAddr, msg, CauseContextNotFound, nil); responseErr != nil {
				return responseErr
			}
			return &InvalidTEIDError{TEID: teid}
		}
	}
	return nil
}

// SendMessageTo assigns a sequence and sends a message with reply correlation
// and retransmission. For an already encoded request use SendRawTo.
func (c *Conn) SendMessageTo(msg message.Message, addr net.Addr) (uint32, error) {
	if c == nil {
		return 0, net.ErrClosed
	}
	if msg == nil || msg.GetHeader() == nil {
		return 0, message.ErrNilMessage
	}
	seq := c.IncSequence()
	// TS 29.274 §7.6: commands occupy the MSB-set sequence space.
	if isCommand(msg.MessageType()) {
		seq |= 0x800000
	}
	if err := msg.SetSequenceNumber(seq); err != nil {
		return seq, err
	}
	payload, err := message.Marshal(msg)
	if err != nil {
		return seq, fmt.Errorf("failed to send %T: %w", msg, err)
	}
	return seq, c.sendReliable(msg, payload, addr)
}

// IncSequence increments the SequenceNumber associated with Conn.
func (c *Conn) IncSequence() uint32 {
	if c == nil {
		return 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for {
		c.sequence = (c.sequence + 1) & 0x7fffff
		used := false
		for key := range c.transactions {
			if key.seq&0x7fffff == c.sequence {
				used = true
				break
			}
		}
		if !used {
			break
		}
	}

	return c.sequence
}

// SequenceNumber returns the current(=last used) SequenceNumber associated with Conn.
func (c *Conn) SequenceNumber() uint32 {
	if c == nil {
		return 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.sequence
}

// EchoRequest sends a EchoRequest.
func (c *Conn) EchoRequest(raddr net.Addr) (uint32, error) {
	msg := message.NewEchoRequest(0, ie.NewRecovery(c.RestartCounter))

	seq, err := c.SendMessageTo(msg, raddr)
	if err != nil {
		return 0, err
	}
	return seq, nil
}

// EchoResponse sends a EchoResponse in response to the EchoRequest.
func (c *Conn) EchoResponse(raddr net.Addr, req message.Message) error {
	res := message.NewEchoResponse(0, ie.NewRecovery(c.RestartCounter))

	if err := c.RespondTo(raddr, req, res); err != nil {
		return err
	}
	return nil
}

// VersionNotSupportedIndication sends VersionNotSupportedIndication message
// in response to any kind of message.Message.
func (c *Conn) VersionNotSupportedIndication(raddr net.Addr, req message.Message) error {
	res := message.NewVersionNotSupportedIndication(0, req.Sequence())

	if err := c.RespondTo(raddr, req, res); err != nil {
		return err
	}
	return nil
}

// ParseCreateSession iterates through the ie and returns a session
func (c *Conn) ParseCreateSession(raddr net.Addr, ies ...*ie.IE) (*Session, error) {
	if c == nil {
		return nil, net.ErrClosed
	}
	// retrieve values from IEs given.
	sess := NewSession(raddr, &Subscriber{Location: &Location{}})
	br := sess.GetDefaultBearer()
	var err error
	bearerIndex := 0
	for _, i := range ies {
		if i == nil {
			continue
		}
		switch i.Type {
		case ie.IMSI:
			sess.IMSI, err = i.IMSI()
			if err != nil {
				return nil, err
			}
		case ie.MSISDN:
			sess.MSISDN, err = i.MSISDN()
			if err != nil {
				return nil, err
			}
		case ie.MobileEquipmentIdentity:
			sess.IMEI, err = i.MobileEquipmentIdentity()
			if err != nil {
				return nil, err
			}
		case ie.ServingNetwork:
			sess.MCC, err = i.MCC()
			if err != nil {
				return nil, err
			}
			sess.MNC, err = i.MNC()
			if err != nil {
				return nil, err
			}
		case ie.AccessPointName:
			br.APN, err = i.AccessPointName()
			if err != nil {
				return nil, err
			}
		case ie.RATType:
			sess.RATType, err = i.RATType()
			if err != nil {
				return nil, err
			}
		case ie.FullyQualifiedTEID:
			it, err := i.InterfaceType()
			if err != nil {
				return nil, err
			}
			teid, err := i.TEID()
			if err != nil {
				return nil, err
			}
			sess.AddTEID(it, teid)

		case ie.BearerContext:
			switch i.Instance() {
			case 0:
				// TS 29.274 Table 7.2.1-2: each context describes its own bearer.
				if bearerIndex > 0 {
					br = NewBearer(0, sess.GetDefaultBearer().APN, &QoSProfile{})
					sess.AddBearer(fmt.Sprintf("bearer-%d", bearerIndex), br)
				}
				bearerIndex++
				for _, child := range i.ChildIEs {
					if child == nil {
						continue
					}
					switch child.Type {
					case ie.EPSBearerID:
						br.EBI, err = child.EPSBearerID()
						if err != nil {
							return nil, err
						}
					case ie.BearerQoS:
						br.PL, err = child.PriorityLevel()
						if err != nil {
							return nil, err
						}
						br.QCI, err = child.QCILabel()
						if err != nil {
							return nil, err
						}
						br.PCI = child.HasPCI()
						br.PVI = child.HasPVI()

						br.MBRUL, err = child.MBRForUplink()
						if err != nil {
							return nil, err
						}
						br.MBRDL, err = child.MBRForDownlink()
						if err != nil {
							return nil, err
						}
						br.GBRUL, err = child.GBRForUplink()
						if err != nil {
							return nil, err
						}
						// TS 29.274 §8.15: downlink GBR is distinct from uplink GBR.
						br.GBRDL, err = child.GBRForDownlink()
						if err != nil {
							return nil, err
						}
					case ie.FullyQualifiedTEID:
						it, err := child.InterfaceType()
						if err != nil {
							return nil, err
						}
						teid, err := child.TEID()
						if err != nil {
							return nil, err
						}
						br.AddTEID(it, teid)
					default:
						br.AdditionalIEs = append(br.AdditionalIEs, child)
					}
				}
			case 1:
				sess.AdditionalIEs = append(sess.AdditionalIEs, i)
			}
		default:
			sess.AdditionalIEs = append(sess.AdditionalIEs, i)
		}
	}
	for _, bearer := range sess.Bearers() {
		for _, i := range ies {
			if i != nil && i.Type == ie.AccessPointName {
				bearer.APN, err = i.AccessPointName()
				if err != nil {
					return nil, err
				}
				break
			}
		}
	}
	if teid, e := sess.GetTEID(c.localIfType); e == nil {
		c.RegisterSession(teid, sess)
	}
	return sess, nil
}

// CreateSession sends a CreateSessionRequest and stores information given with IE
// in the Session returned.
//
// After using this method, users don't need to call AddSession with the session
// returned.
//
// By creating a Session with this method, the values in IEs given, such as TEID in F-TEID
// are stored with "best effort". See the source code to see what kind information is
// handled automatically in this method.
//
// Also, a Bearer named "default" is also created to be used as default bearer.
// The default bearer can be retrieved by using GetDefaultBearer() or LookupBearerByName("default").
//
// Note that this method doesn't care IEs given are sufficient or not, as the required IE
// varies much depending on the context in which the Create Session Request is used.
// In other words, any kind of IE can be put on the Create Session Request message using
// this method.
func (c *Conn) CreateSession(raddr net.Addr, ie ...*ie.IE) (*Session, uint32, error) {

	sess, err := c.ParseCreateSession(raddr, ie...)
	if err != nil {
		return nil, 0, err
	}

	// set IEs into CreateSessionRequest.
	msg := message.NewCreateSessionRequest(0, 0, ie...)

	seq, err := c.SendMessageTo(msg, raddr)
	if err != nil {
		c.RemoveSession(sess)
		return nil, 0, err
	}
	return sess, seq, nil
}

// DeleteSession sends a DeleteSessionRequest with TEID and IEs given.
func (c *Conn) DeleteSession(teid uint32, sess *Session, ie ...*ie.IE) (uint32, error) {
	if sess == nil {
		return 0, &InvalidSessionError{}
	}
	msg := message.NewDeleteSessionRequest(teid, 0, ie...)

	seq, err := c.SendMessageTo(msg, sess.PeerAddr())
	if err != nil {
		return 0, err
	}
	return seq, nil
}

// ModifyBearer sends a ModifyBearerRequest with TEID and IEs given..
func (c *Conn) ModifyBearer(teid uint32, sess *Session, ie ...*ie.IE) (uint32, error) {
	if sess == nil {
		return 0, &InvalidSessionError{}
	}
	msg := message.NewModifyBearerRequest(teid, 0, ie...)

	seq, err := c.SendMessageTo(msg, sess.PeerAddr())
	if err != nil {
		return 0, err
	}
	return seq, nil
}

// DeleteBearer sends a DeleteBearerRequest TEID and with IEs given.
func (c *Conn) DeleteBearer(teid uint32, sess *Session, ie ...*ie.IE) (uint32, error) {
	if sess == nil {
		return 0, &InvalidSessionError{}
	}
	msg := message.NewDeleteBearerRequest(teid, 0, ie...)

	seq, err := c.SendMessageTo(msg, sess.PeerAddr())
	if err != nil {
		return 0, err
	}
	return seq, nil
}

// RespondTo sends a message(specified with "toBeSent" param) in response to a message
// (specified with "received" param).
//
// This exists to make it easier to handle SequenceNumber.
func (c *Conn) RespondTo(raddr net.Addr, received, toBeSent message.Message) error {
	if c == nil {
		return net.ErrClosed
	}
	if received == nil || toBeSent == nil || received.GetHeader() == nil || toBeSent.GetHeader() == nil {
		return message.ErrNilMessage
	}
	if err := toBeSent.SetSequenceNumber(received.Sequence()); err != nil {
		return err
	}
	// TS 29.274 §7.6: the sender assigns a distinct initial sequence.
	if piggy := toBeSent.GetHeader().Piggybacked; piggy != nil && piggy.Sequence() == 0 {
		if err := piggy.SetSequenceNumber(c.IncSequence()); err != nil {
			return err
		}
	}
	b := make([]byte, toBeSent.MarshalLen())

	if err := toBeSent.MarshalTo(b); err != nil {
		return err
	}

	if err := c.sendReliable(toBeSent, b, raddr); err != nil {
		return err
	}
	// TS 29.274 §7.6: only replay responses that were successfully sent.
	c.cacheResponse(raddr, received, b)
	return nil
}

// GetSessionByTEID returns the session for a locally allocated TEID.
// Applications may authorize the source separately using Session.PeerAddr.
func (c *Conn) GetSessionByTEID(teid uint32) (*Session, error) {
	if c == nil {
		return nil, &InvalidTEIDError{TEID: teid}
	}
	c.mu.Lock()
	session := c.sessionsByTEID[teid]
	c.mu.Unlock()
	// TS 29.274 §§5.5.2, 4.2.2.1: the receiver allocates TEID-C;
	// relocation can change the sender's address, including its IP family.
	if session == nil {
		return nil, &InvalidTEIDError{TEID: teid}
	}
	return session, nil
}

// GetSessionsByIMSI returns every PDN connection associated with the subscriber.
func (c *Conn) GetSessionsByIMSI(imsi string) []*Session {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	var sessions []*Session
	for s := range c.sessions {
		if s.IMSI == imsi {
			sessions = append(sessions, s)
		}
	}
	return sessions
}

// GetSessionByIMSI returns the unique session, or an error if the IMSI is ambiguous.
func (c *Conn) GetSessionByIMSI(imsi string) (*Session, error) {
	if c == nil {
		return nil, &UnknownIMSIError{IMSI: imsi}
	}
	sessions := c.GetSessionsByIMSI(imsi)
	if len(sessions) == 0 {
		return nil, &UnknownIMSIError{IMSI: imsi}
	}
	if len(sessions) > 1 {
		return nil, ErrAmbiguousIMSI
	}
	return sessions[0], nil
}

// GetIMSIByTEID returns the IMSI associated with the local TEID.
func (c *Conn) GetIMSIByTEID(teid uint32) (string, error) {
	if c == nil {
		return "", &InvalidTEIDError{TEID: teid}
	}
	sess, err := c.GetSessionByTEID(teid)
	if err != nil {
		return "", err
	}
	return sess.IMSI, nil
}

// RegisterSession publishes the session and atomically remaps its incoming TEID.
func (c *Conn) RegisterSession(itei uint32, session *Session) {
	if c == nil {
		return
	}
	if session == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for teid, s := range c.sessionsByTEID {
		if s == session {
			delete(c.sessionsByTEID, teid)
		}
	}
	if old := c.sessionsByTEID[itei]; old != nil && old != session {
		delete(c.sessions, old)
	}
	c.sessionsByTEID[itei] = session
	delete(c.reservedTEIDs, itei)
	c.sessions[session] = struct{}{}
	session.AddTEID(c.localIfType, itei)
}

// RemoveSession removes only indexes that still refer to this session.
func (c *Conn) RemoveSession(session *Session) {
	if c == nil {
		return
	}
	if session == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for teid, s := range c.sessionsByTEID {
		if s == session {
			delete(c.sessionsByTEID, teid)
		}
	}
	delete(c.sessions, session)
}

// RemoveSessionByIMSI removes all PDN connections associated with the IMSI.
func (c *Conn) RemoveSessionByIMSI(imsi string) {
	if c == nil {
		return
	}
	for _, s := range c.GetSessionsByIMSI(imsi) {
		c.RemoveSession(s)
	}
}

// NewSenderFTEID creates a new F-TEID with random TEID value that is unique within Conn.
// To ensure the uniqueness, don't create in the other way if you once use this method.
// This is meant to be used for creating F-TEID IE only for local interface type that is
// specified at the creation of Conn.
//
// Note that in the case there's a lot of Session on the Conn, it may take a long
// time to find a new unique value.
//
// TODO: optimize performance...
func (c *Conn) NewSenderFTEID(v4, v6 string) (fteidIE *ie.IE) {
	if c == nil {
		return nil
	}
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
		c.mu.Lock()
		_, reserved := c.reservedTEIDs[t]
		_, registered := c.sessionsByTEID[t]
		if !reserved && !registered {
			c.reservedTEIDs[t] = struct{}{}
		}
		c.mu.Unlock()
		if reserved || registered {
			continue
		}

		teid = t
		break
	}

	if teid == 0 {
		return nil
	}
	fteid := ie.NewFullyQualifiedTEID(c.localIfType, teid, v4, v6)
	if fteid.Err() != nil {
		c.mu.Lock()
		delete(c.reservedTEIDs, teid)
		c.mu.Unlock()
		return nil
	}
	return fteid
}

func generateRandomUint32() uint32 {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return 0
	}

	return binary.BigEndian.Uint32(b)
}

// Sessions returns all the sessions registered in Conn.
func (c *Conn) Sessions() []*Session {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	ss := make([]*Session, 0, len(c.sessions))
	for s := range c.sessions {
		ss = append(ss, s)
	}
	return ss
}

// SessionCount returns the number of active PDN connections.
func (c *Conn) SessionCount() int {
	if c == nil {
		return 0
	}
	count := 0
	for _, s := range c.Sessions() {
		if s.IsActive() {
			count++
		}
	}
	return count
}

// BearerCount returns the number of bearers in active PDN connections.
func (c *Conn) BearerCount() int {
	if c == nil {
		return 0
	}
	count := 0
	for _, s := range c.Sessions() {
		if s.IsActive() {
			count += s.BearerCount()
		}
	}
	return count
}
