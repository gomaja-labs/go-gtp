// Copyright 2026 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv2

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

type transactionKey struct {
	peer string
	seq  uint32
	typ  uint8
}
type transaction struct {
	payload  []byte
	peer     net.Addr
	expected []uint8
	attempts int
	interval time.Duration
	retries  int
	timer    *time.Timer
}
type requestKey struct {
	peer string
	seq  uint32
	typ  uint8
}
type receivedRequest struct {
	acceptedPrimary bool
	key             requestKey
	payload         []byte
	expires         time.Time
	order, answered *list.Element
}

type replayPeer struct {
	ip       string
	requests list.List
	answered list.List
	order    *list.Element
}

const (
	maxReplayPeers           = 64
	maxReplayRequestsPerPeer = 64
)

// SetReplayWindow sets retention for future incoming requests, independently of
// local retry settings. Configure it to cover the longest peer retry window
// (TS 29.274 §7.6 permits per-procedure timers). The default is three minutes.
// The cache holds at most 64 peer IPs and 64 requests per IP, including ports.
// At capacity the oldest answered request of that peer is evicted; unanswered
// requests are protected. New peer IPs replace the least recently used peer.
// Replay protection therefore lasts until expiry or capacity eviction.
func (c *Conn) SetReplayWindow(window time.Duration) error {
	if c == nil {
		return net.ErrClosed
	}
	if window <= 0 {
		return errors.New("invalid replay window")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.replayWindow = window
	return nil
}

// SendRawTo validates and sends an encoded request, preserving its sequence and
// bytes while registering the same reply correlation/retries as SendMessageTo.
// The caller must assign a unique sequence (and set the MSB for commands,
// TS 29.274 §7.6). Unregistered WriteTo requests cannot receive replies.
func (c *Conn) SendRawTo(payload []byte, addr net.Addr) error {
	if c == nil {
		return net.ErrClosed
	}
	msg, err := message.Parse(payload)
	if err != nil {
		return err
	}
	if err := msg.GetHeader().PiggybackError; err != nil {
		return err
	}
	if isCommand(msg.MessageType()) && msg.Sequence()&0x800000 == 0 {
		return message.ErrInvalidSequence
	}
	return c.sendReliable(msg, payload, addr)
}

// TransactionError reports a failed or invalid exchange to the application.
type TransactionError struct {
	Peer        net.Addr
	Sequence    uint32
	MessageType uint8
	Err         error
}

func (e *TransactionError) Error() string {
	return fmt.Sprintf("GTP transaction peer %v sequence %d type %d: %v", e.Peer, e.Sequence, e.MessageType, e.Err)
}
func (e *TransactionError) Unwrap() error { return e.Err }

// SetErrorHandler installs the callback for failed transfers and invalid replies.
// It may run on a timer goroutine and must return promptly. A nil callback logs errors.
func (c *Conn) SetErrorHandler(handler func(error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errorHandler = handler
}
func (c *Conn) reportError(err error) {
	c.mu.Lock()
	handler := c.errorHandler
	c.mu.Unlock()
	if handler != nil {
		handler(err)
	} else {
		logf("%v", err)
	}
}

// SetRetransmission configures T3-RESPONSE and N3-REQUESTS for future exchanges.
// retries is the number of retransmissions after the original send.
// TS 29.274 §7.6 leaves these values implementation dependent; defaults are 3s and 3.
func (c *Conn) SetRetransmission(interval time.Duration, retries int) error {
	if interval <= 0 || retries < 0 {
		return errors.New("invalid retransmission interval or retry count")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.retransmitInterval, c.retransmitRetries = interval, retries
	return nil
}

func peerIP(addr net.Addr) string {
	if addr == nil {
		return ""
	}
	if udp, ok := addr.(*net.UDPAddr); ok {
		return udp.IP.String() + "%" + udp.Zone
	}
	host, _, err := net.SplitHostPort(addr.String())
	if err == nil {
		return host
	}
	return addr.String()
}
func samePeer(a, b net.Addr) bool     { return a != nil && b != nil && peerIP(a) == peerIP(b) }
func sameEndpoint(a, b net.Addr) bool { return samePeer(a, b) && a.String() == b.String() }
func copyAddr(addr net.Addr) net.Addr {
	if udp, ok := addr.(*net.UDPAddr); ok {
		cp := *udp
		cp.IP = append(net.IP(nil), udp.IP...)
		return &cp
	}
	return addr
}
func isCommand(typ uint8) bool {
	return typ == message.MsgTypeModifyBearerCommand || typ == message.MsgTypeDeleteBearerCommand || typ == message.MsgTypeBearerResourceCommand
}

// TS 29.274 Table 6.1-1 and §4.2.5: Initial/Triggered/Triggered Reply pairs.
var replies = map[uint8][]uint8{
	message.MsgTypeEchoRequest:                   {message.MsgTypeEchoResponse},
	message.MsgTypeCreateSessionRequest:          {message.MsgTypeCreateSessionResponse},
	message.MsgTypeModifyBearerRequest:           {message.MsgTypeModifyBearerResponse},
	message.MsgTypeDeleteSessionRequest:          {message.MsgTypeDeleteSessionResponse},
	message.MsgTypeCreateBearerRequest:           {message.MsgTypeCreateBearerResponse},
	message.MsgTypeUpdateBearerRequest:           {message.MsgTypeUpdateBearerResponse},
	message.MsgTypeDeleteBearerRequest:           {message.MsgTypeDeleteBearerResponse},
	message.MsgTypeModifyBearerCommand:           {message.MsgTypeUpdateBearerRequest, message.MsgTypeModifyBearerFailureIndication},
	message.MsgTypeDeleteBearerCommand:           {message.MsgTypeDeleteBearerRequest, message.MsgTypeDeleteBearerFailureIndication},
	message.MsgTypeBearerResourceCommand:         {message.MsgTypeCreateBearerRequest, message.MsgTypeUpdateBearerRequest, message.MsgTypeDeleteBearerRequest, message.MsgTypeBearerResourceFailureIndication},
	message.MsgTypeContextRequest:                {message.MsgTypeContextResponse},
	message.MsgTypeContextResponse:               {message.MsgTypeContextAcknowledge},
	message.MsgTypeReleaseAccessBearersRequest:   {message.MsgTypeReleaseAccessBearersResponse},
	message.MsgTypeModifyAccessBearersRequest:    {message.MsgTypeModifyAccessBearersResponse},
	message.MsgTypeDeletePDNConnectionSetRequest: {message.MsgTypeDeletePDNConnectionSetResponse},
	message.MsgTypeUpdatePDNConnectionSetRequest: {message.MsgTypeUpdatePDNConnectionSetResponse},
	message.MsgTypePGWRestartNotification:        {message.MsgTypePGWRestartNotificationAcknowledge},
	message.MsgTypeDetachNotification:            {message.MsgTypeDetachAcknowledge},
	message.MsgTypeResumeNotification:            {message.MsgTypeResumeAcknowledge},
	message.MsgTypeSuspendNotification:           {message.MsgTypeSuspendAcknowledge},
	message.MsgTypeChangeNotificationRequest:     {message.MsgTypeChangeNotificationResponse},
	message.MsgTypeDownlinkDataNotification:      {message.MsgTypeDownlinkDataNotificationAcknowledge},
}

func containsType(types []uint8, typ uint8) bool {
	for _, t := range types {
		if t == typ {
			return true
		}
	}
	return false
}
func isReply(typ uint8) bool {
	for initial, types := range replies {
		if isCommand(initial) {
			continue
		}
		if containsType(types, typ) {
			return true
		}
	}
	return typ == message.MsgTypeVersionNotSupportedIndication || typ == message.MsgTypeModifyBearerFailureIndication || typ == message.MsgTypeDeleteBearerFailureIndication || typ == message.MsgTypeBearerResourceFailureIndication
}

func (c *Conn) sendReliable(msg message.Message, payload []byte, addr net.Addr) error {
	if addr == nil {
		return errors.New("nil peer address")
	}
	expected := replies[msg.MessageType()]
	// A rejected Context Response does not expect an acknowledgement (§7.3.7).
	if msg.MessageType() == message.MsgTypeContextResponse && rejectionCause(msg) {
		expected = nil
	}
	var key transactionKey
	var tx *transaction
	if len(expected) != 0 {
		key = transactionKey{peerIP(addr), msg.Sequence(), msg.MessageType()}
		c.mu.Lock()
		select {
		case <-c.closeCh:
			c.mu.Unlock()
			return net.ErrClosed
		default:
		}
		if _, exists := c.transactions[key]; exists {
			c.mu.Unlock()
			return errors.New("transaction already outstanding")
		}
		tx = &transaction{payload: append([]byte(nil), payload...), peer: copyAddr(addr), expected: expected, interval: c.retransmitInterval, retries: c.retransmitRetries}
		c.transactions[key] = tx
		// Register before writing: a loopback reply can arrive before WriteTo returns.
		tx.timer = time.AfterFunc(tx.interval, func() { c.retry(key, tx) })
		c.mu.Unlock()
	}

	var piggyKey transactionKey
	var piggyTx *transaction
	if piggy := msg.GetHeader().Piggybacked; piggy != nil {
		piggyKey = transactionKey{peerIP(addr), piggy.Sequence(), piggy.MessageType()}
		var err error
		piggyTx, err = c.trackPiggyback(piggy, payload, addr)
		if err != nil {
			c.mu.Lock()
			if tx != nil && c.transactions[key] == tx {
				delete(c.transactions, key)
				tx.timer.Stop()
			}
			c.mu.Unlock()
			return err
		}
	}
	if _, err := c.WriteTo(payload, addr); err != nil {
		c.mu.Lock()
		if tx != nil && c.transactions[key] == tx {
			delete(c.transactions, key)
			tx.timer.Stop()
		}
		if piggyTx != nil && c.transactions[piggyKey] == piggyTx {
			delete(c.transactions, piggyKey)
			piggyTx.timer.Stop()
		}
		c.mu.Unlock()
		return err
	}
	return nil
}

func (c *Conn) trackPiggyback(msg message.Message, payload []byte, addr net.Addr) (*transaction, error) {
	key := transactionKey{peerIP(addr), msg.Sequence(), msg.MessageType()}
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.closeCh:
		return nil, net.ErrClosed
	default:
	}
	if _, exists := c.transactions[key]; exists {
		return nil, errors.New("piggybacked transaction already outstanding")
	}
	tx := &transaction{payload: append([]byte(nil), payload...), peer: copyAddr(addr), expected: replies[msg.MessageType()], interval: c.retransmitInterval, retries: c.retransmitRetries}
	c.transactions[key] = tx
	tx.timer = time.AfterFunc(tx.interval, func() { c.retry(key, tx) })
	return tx, nil
}

func (c *Conn) retry(key transactionKey, tx *transaction) {
	c.mu.Lock()
	if c.transactions[key] != tx {
		c.mu.Unlock()
		return
	}
	if tx.attempts == tx.retries {
		delete(c.transactions, key)
		c.mu.Unlock()
		c.reportError(&TransactionError{Peer: tx.peer, Sequence: key.seq, MessageType: key.typ, Err: ErrTimeout})
		return
	}
	tx.attempts++
	// Keep completion serialized with retry selection; the socket itself is safe for concurrent I/O.
	pc := c.pktConn
	c.mu.Unlock()
	var err error
	if pc == nil {
		err = net.ErrClosed
	} else {
		_, err = pc.WriteTo(tx.payload, tx.peer)
	}
	c.mu.Lock()
	if c.transactions[key] != tx {
		c.mu.Unlock()
		return
	}
	if err != nil {
		delete(c.transactions, key)
		c.mu.Unlock()
		c.reportError(&TransactionError{tx.peer, key.seq, key.typ, err})
		return
	}
	tx.timer = time.AfterFunc(tx.interval, func() { c.retry(key, tx) })
	c.mu.Unlock()
}

func (c *Conn) matchReply(peer net.Addr, seq uint32, typ uint8) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.matchReplyLocked(peer, seq, typ)
}

func (c *Conn) matchReplyLocked(peer net.Addr, seq uint32, typ uint8) bool {
	matched := false
	for key, tx := range c.transactions {
		// TS 29.274 §§4.2.1.2, 7.6: match the endpoint. §5.3: VNSI's
		// sequence is ignored, so stop all exchanges to that endpoint.
		if sameEndpoint(peer, tx.peer) && (typ == message.MsgTypeVersionNotSupportedIndication || key.seq == seq && containsType(tx.expected, typ)) {
			tx.timer.Stop()
			delete(c.transactions, key)
			matched = true
			if typ != message.MsgTypeVersionNotSupportedIndication {
				break
			}
		}
	}
	return matched
}

// replayPeerLocked expires entries in arrival order, never scanning the cache.
func (c *Conn) replayPeerLocked(peer net.Addr, now time.Time) *replayPeer {
	ip := peerIP(peer)
	p := c.replayPeers[ip]
	if p == nil {
		if len(c.replayPeers) == maxReplayPeers {
			oldest := c.replayPeerOrder.Front().Value.(*replayPeer)
			for oldest.requests.Len() != 0 {
				c.removeReplayLocked(oldest, oldest.requests.Front().Value.(*receivedRequest))
			}
			c.replayPeerOrder.Remove(oldest.order)
			delete(c.replayPeers, oldest.ip)
		}
		p = &replayPeer{ip: ip}
		p.order = c.replayPeerOrder.PushBack(p)
		c.replayPeers[ip] = p
	} else {
		c.replayPeerOrder.MoveToBack(p.order)
	}
	for p.requests.Len() != 0 {
		r := p.requests.Front().Value.(*receivedRequest)
		if now.Before(r.expires) {
			break
		}
		c.removeReplayLocked(p, r)
	}
	return p
}

func (c *Conn) removeReplayLocked(p *replayPeer, r *receivedRequest) {
	delete(c.received, r.key)
	p.requests.Remove(r.order)
	if r.answered != nil {
		p.answered.Remove(r.answered)
	}
}

func (c *Conn) isDuplicateOrUnmatched(peer net.Addr, msg message.Message) bool {
	if peer == nil {
		return true
	}
	typ := msg.MessageType()
	key := requestKey{peer.String(), msg.Sequence(), typ}
	c.mu.Lock()
	var p *replayPeer
	if len(replies[typ]) != 0 {
		p = c.replayPeerLocked(peer, time.Now())
		if prev := c.received[key]; prev != nil {
			payload := append([]byte(nil), prev.payload...)
			c.mu.Unlock()
			if len(payload) != 0 {
				if _, err := c.WriteTo(payload, peer); err != nil {
					logf("replaying GTP response: %v", err)
				}
			}
			return true
		}
		if p.requests.Len() == maxReplayRequestsPerPeer {
			// Never complete a command if its triggered request cannot be admitted.
			if p.answered.Len() == 0 {
				c.mu.Unlock()
				return true
			}
			c.removeReplayLocked(p, p.answered.Front().Value.(*receivedRequest))
		}
	}
	matched := false
	if isReply(typ) || msg.Sequence()&0x800000 != 0 {
		matched = c.matchReplyLocked(peer, msg.Sequence(), typ)
	}
	if isReply(typ) && !matched {
		c.mu.Unlock()
		return true
	}
	if p != nil {
		r := &receivedRequest{key: key, expires: time.Now().Add(c.replayWindow)}
		c.insertReplayLocked(p, r)
	}
	c.mu.Unlock()
	return false
}

func (c *Conn) cacheResponse(peer net.Addr, received message.Message, payload []byte) {
	key := requestKey{peer.String(), received.Sequence(), received.MessageType()}
	c.mu.Lock()
	defer c.mu.Unlock()
	if prev := c.received[key]; prev != nil {
		// UDP payloads are bounded; one peer can retain at most 64 datagrams.
		if len(payload) > 65535 {
			return
		}
		prev.payload = append([]byte(nil), payload...)
		if prev.answered == nil {
			prev.answered = c.replayPeers[peerIP(peer)].answered.PushBack(prev)
		}
	}
}

// insertReplayLocked preserves expiration order even if SetReplayWindow changes.
func (c *Conn) insertReplayLocked(p *replayPeer, r *receivedRequest) {
	at := p.requests.Back()
	for at != nil && at.Value.(*receivedRequest).expires.After(r.expires) {
		at = at.Prev()
	}
	if at == nil {
		r.order = p.requests.PushFront(r)
	} else {
		r.order = p.requests.InsertAfter(r, at)
	}
	c.received[r.key] = r
}

// acceptedPrimary remembers successful establishment so retransmission of a
// piggybacked datagram can replay the dedicated reply without repeating the
// primary handler (TS 29.274 §§7.6, 7.7.0). It shares all replay cache bounds.
func (c *Conn) acceptedPrimary(peer net.Addr, msg message.Message, remember bool) bool {
	if peer == nil || msg.MessageType() != message.MsgTypeCreateSessionResponse && msg.MessageType() != message.MsgTypeCreateBearerResponse {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	p := c.replayPeerLocked(peer, time.Now())
	key := requestKey{peer.String(), msg.Sequence(), msg.MessageType()}
	r := c.received[key]
	if remember && r == nil {
		if p.requests.Len() == maxReplayRequestsPerPeer {
			if p.answered.Len() == 0 {
				return false
			}
			c.removeReplayLocked(p, p.answered.Front().Value.(*receivedRequest))
		}
		r = &receivedRequest{key: key, expires: time.Now().Add(c.replayWindow), acceptedPrimary: true}
		c.insertReplayLocked(p, r)
		r.answered = p.answered.PushBack(r)
	}
	return r != nil && r.acceptedPrimary
}
