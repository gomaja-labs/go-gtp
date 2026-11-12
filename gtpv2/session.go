// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv2

import (
	"net"
	"sync"
	"time"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
	"github.com/gomaja-labs/go-gtp/gtpv2/message"
)

// Location is a subscriber's location.
type Location struct {
	MCC, MNC               string
	RATType                uint8
	LAC, CI, SAI, RAI, TAI uint16
	ECI, MeNBI, EMeNBI     uint32
}

// Subscriber is a subscriber that belongs to a GTPv2 session.
type Subscriber struct {
	IMSI, MSISDN, IMEI string
	*Location
}

// Session is a GTPv2 Session.
type Session struct {
	AdditionalIEs []*ie.IE
	mu            sync.Mutex
	isActive      bool
	*teidMap
	*bearerMap

	// Messages awaiting their sequence-specific waiter.
	pending   []pendingMessage
	abandoned map[uint32]time.Time
	// abandonedOrder lists abandoned sequences oldest first, for expiry.
	abandonedOrder []abandonedSequence
	changed        chan struct{}

	// peerAddr is protected by mu; getters return an address snapshot.
	peerAddr net.Addr

	// Subscriber is a Subscriber associated with Session.
	*Subscriber
}

// NewSession creates a new Session with subscriber information.
//
// This is expected to be used by server-like nodes. Otherwise, use CreateSession(),
// which sends Create Session Request and returns a new Session.
func NewSession(peerAddr net.Addr, sub *Subscriber) *Session {
	if sub == nil {
		sub = &Subscriber{}
	}
	if sub.Location == nil {
		sub.Location = &Location{}
	}
	s := &Session{
		mu:       sync.Mutex{},
		peerAddr: copyAddr(peerAddr),

		teidMap:    newTeidMap(),
		bearerMap:  newBearerMap("default", NewBearer(0, "", &QoSProfile{})),
		Subscriber: sub,
		changed:    make(chan struct{}),
		abandoned:  make(map[uint32]time.Time),
	}

	return s
}

// Activate marks a Session active.
func (s *Session) Activate() error {
	if s == nil {
		return &InvalidSessionError{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.IMSI == "" {
		return &RequiredParameterMissingError{"IMSI", "Session must have IMSI set"}
	}

	s.isActive = true
	return nil
}

// Deactivate marks a Session inactive.
func (s *Session) Deactivate() error {
	if s == nil {
		return &InvalidSessionError{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.isActive = false
	return nil
}

// IsActive reports whether a Session is active or not.
func (s *Session) IsActive() bool {
	if s == nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.isActive
}

// PeerAddr returns the address of the peer node associated with Session.
func (s *Session) PeerAddr() net.Addr {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return copyAddr(s.peerAddr)
}

// UpdatePeerAddr updates the address of the peer node associated with Session.
func (s *Session) UpdatePeerAddr(peer net.Addr) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.peerAddr = copyAddr(peer)
}

// AddTEID adds TEID to session with InterfaceType.
//
// This is used to keep TEIDs of any interface types that may be used later,
// including the ones that are assigned to U-Plane.
//
// For incoming TEID of local interface type, (*Conn).RegisterSession does that
// instead of users but it is safe to call it.
func (s *Session) AddTEID(ifType uint8, teid uint32) {
	if s == nil {
		return
	}
	s.teidMap.store(ifType, teid)
}

// GetTEID returns TEID associated with InterfaceType given.
func (s *Session) GetTEID(ifType uint8) (uint32, error) {
	if s == nil {
		return 0, ErrTEIDNotFound
	}
	if teid, ok := s.teidMap.load(ifType); ok {
		return teid, nil
	}
	return 0, ErrTEIDNotFound
}

// pendingMessage is a message passed to a Session and not yet claimed by
// WaitMessage.
type pendingMessage struct {
	msg     message.Message
	arrived time.Time
}

// abandonedSequence records when the WaitMessage for seq timed out.
type abandonedSequence struct {
	seq uint32
	at  time.Time
}

const maxPendingMessages = 1000

// pendingMessageTTL bounds how long a Session keeps a message nobody waits for,
// and how long a sequence whose WaitMessage timed out keeps dropping its late
// reply. It covers a reply that arrives just before WaitMessage is called.
var pendingMessageTTL = 30 * time.Second

// PassMessageTo hands msg to the WaitMessage call waiting for its sequence
// number. It never blocks: a reply whose waiter already timed out is dropped,
// messages nobody claims expire after a while, and when a Session holds
// maxPendingMessages unclaimed messages the oldest one is discarded.
func PassMessageTo(s *Session, msg message.Message) error {
	if s == nil {
		return &InvalidSessionError{}
	}
	if msg == nil || msg.GetHeader() == nil {
		return message.ErrNilMessage
	}
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pruneLocked(now)
	if _, ok := s.abandoned[msg.Sequence()]; ok {
		return nil // late reply for a waiter that timed out
	}
	if len(s.pending) == maxPendingMessages {
		s.pending[0] = pendingMessage{}
		s.pending = s.pending[1:]
	}
	s.pending = append(s.pending, pendingMessage{msg: msg, arrived: now})
	close(s.changed)
	s.changed = make(chan struct{})
	return nil
}

// pruneLocked drops expired unclaimed messages and abandoned sequences. Both
// lists are ordered by time, so only expired entries are visited.
// s.mu must be held.
func (s *Session) pruneLocked(now time.Time) {
	for len(s.pending) > 0 && now.Sub(s.pending[0].arrived) >= pendingMessageTTL {
		s.pending[0] = pendingMessage{}
		s.pending = s.pending[1:]
	}
	for len(s.abandoned) > 0 && len(s.abandonedOrder) > 0 && now.Sub(s.abandonedOrder[0].at) >= pendingMessageTTL {
		a := s.abandonedOrder[0]
		if at, ok := s.abandoned[a.seq]; ok && at.Equal(a.at) {
			delete(s.abandoned, a.seq)
		}
		s.abandonedOrder = s.abandonedOrder[1:]
	}
	if len(s.abandoned) == 0 {
		s.abandonedOrder = nil
	}
}

// takeLocked removes and returns the pending message with sequence seq.
// s.mu must be held.
func (s *Session) takeLocked(seq uint32) (message.Message, bool) {
	for i, p := range s.pending {
		if p.msg.Sequence() != seq {
			continue
		}
		copy(s.pending[i:], s.pending[i+1:])
		s.pending[len(s.pending)-1] = pendingMessage{}
		s.pending = s.pending[:len(s.pending)-1]
		return p.msg, true
	}
	return nil, false
}

// WaitMessage waits for the message with sequence number seq, leaving messages
// of other transactions queued (TS 29.274 §7.6 permits several outstanding
// messages on one session). After a timeout, a late reply for seq is dropped
// when it arrives instead of occupying the queue.
func (s *Session) WaitMessage(seq uint32, timeout time.Duration) (message.Message, error) {
	if s == nil {
		return nil, &InvalidSessionError{}
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	s.mu.Lock()
	delete(s.abandoned, seq) // a new waiter for a reused sequence
	s.mu.Unlock()
	for {
		s.mu.Lock()
		if msg, ok := s.takeLocked(seq); ok {
			s.mu.Unlock()
			return msg, nil
		}
		changed := s.changed
		s.mu.Unlock()
		select {
		case <-changed:
		case <-timer.C:
			s.mu.Lock()
			defer s.mu.Unlock()
			if msg, ok := s.takeLocked(seq); ok {
				return msg, nil
			}
			// Only this Session's own timed-out waits are recorded here, never
			// remote input, and pruneLocked expires them after pendingMessageTTL.
			now := time.Now()
			s.abandoned[seq] = now
			s.abandonedOrder = append(s.abandonedOrder, abandonedSequence{seq: seq, at: now})
			return nil, ErrTimeout
		}
	}
}

// AddBearer adds a Bearer to Session with arbitrary name given.
//
// In the single-bearer environment it is not used, as a bearer named "default" is
// always available after created a Session.
func (s *Session) AddBearer(name string, br *Bearer) {
	if s == nil {
		return
	}
	s.bearerMap.store(name, br)
}

// RemoveBearer removes a Bearer looked up by name.
func (s *Session) RemoveBearer(name string) {
	if s == nil {
		return
	}
	s.bearerMap.delete(name)
}

// RemoveBearerByEBI removes a Bearer looked up by name.
func (s *Session) RemoveBearerByEBI(ebi uint8) {
	if s == nil {
		return
	}
	name, err := s.LookupBearerNameByEBI(ebi)
	if err != nil {
		return
	}
	s.bearerMap.delete(name)
}

// GetDefaultBearer returns the default bearer.
func (s *Session) GetDefaultBearer() *Bearer {
	if s == nil {
		return nil
	}
	// it is not expected that the default bearer cannot be found.
	bearer, ok := s.bearerMap.load("default")
	if !ok {
		return nil
	}

	return bearer
}

// SetDefaultBearer sets given bearer as the default bearer.
func (s *Session) SetDefaultBearer(bearer *Bearer) {
	if s == nil {
		return
	}
	// it is not expected that the default bearer cannot be found.
	s.bearerMap.store("default", bearer)
}

// LookupBearerByName looks up Bearer registered in Session by name.
func (s *Session) LookupBearerByName(name string) (*Bearer, error) {
	if s == nil {
		return nil, &BearerNotFoundError{}
	}
	if br, ok := s.bearerMap.load(name); ok {
		return br, nil
	}

	return nil, &BearerNotFoundError{IMSI: s.IMSI}
}

// LookupBearerByEBI looks up Bearer registered in Session by EBI.
func (s *Session) LookupBearerByEBI(ebi uint8) (*Bearer, error) {
	if s == nil {
		return nil, &BearerNotFoundError{}
	}
	var bearer *Bearer
	s.bearerMap.rangeWithFunc(func(name, br interface{}) bool {
		b := br.(*Bearer)
		if ebi == b.EBI {
			bearer = b
			return false
		}
		return true
	})

	if bearer == nil {
		return nil, &BearerNotFoundError{IMSI: s.IMSI}

	}
	return bearer, nil
}

// LookupBearerNameByEBI looks up name of Bearer by EBI and returns
// its name.
func (s *Session) LookupBearerNameByEBI(ebi uint8) (string, error) {
	if s == nil {
		return "", &BearerNotFoundError{}
	}
	var name string
	found := false
	s.bearerMap.rangeWithFunc(func(n, br interface{}) bool {
		bearer := br.(*Bearer)
		if ebi == bearer.EBI {
			name = n.(string)
			found = true
			return false
		}
		return true
	})

	if !found {
		return "", &BearerNotFoundError{IMSI: s.IMSI}

	}
	return name, nil
}

// LookupEBIByName returns EBI associated with name.
//
// If no EBI found, it returns 0(=invalid value for EBI).
func (s *Session) LookupEBIByName(name string) uint8 {
	if s == nil {
		return 0
	}
	if br, ok := s.bearerMap.load(name); ok {
		return br.EBI
	}

	return 0
}

// LookupEBIByTEID returns EBI associated with TEID.
//
// If no EBI found, it returns 0(=invalid value for EBI).
func (s *Session) LookupEBIByTEID(teid uint32) uint8 {
	if s == nil {
		return 0
	}
	var ebi uint8
	s.bearerMap.rangeWithFunc(func(name, bearer interface{}) bool {
		br := bearer.(*Bearer)
		if teid == br.IncomingTEID() || teid == br.OutgoingTEID() {
			ebi = br.EBI
			return false
		}
		return true
	})

	return ebi
}

type teidMap struct {
	syncMap sync.Map
}

func newTeidMap() *teidMap {
	return &teidMap{}
}

func (t *teidMap) store(ifType uint8, teid uint32) {
	t.syncMap.Store(ifType, teid)
}

func (t *teidMap) load(ifType uint8) (uint32, bool) {
	teid, ok := t.syncMap.Load(ifType)
	if !ok {
		return 0, false
	}

	return teid.(uint32), true
}

type bearerMap struct {
	syncMap sync.Map
}

func newBearerMap(name string, bearer *Bearer) *bearerMap {
	b := &bearerMap{}
	b.store(name, bearer)

	return b
}
func (b *bearerMap) store(name string, bearer *Bearer) {
	if bearer != nil {
		b.syncMap.Store(name, bearer)
	}
}

func (b *bearerMap) load(name string) (*Bearer, bool) {
	bearer, ok := b.syncMap.Load(name)
	if !ok {
		return nil, false
	}

	return bearer.(*Bearer), true
}

func (b *bearerMap) delete(name string) {
	b.syncMap.Delete(name)
}

func (b *bearerMap) rangeWithFunc(fn func(name, bearer interface{}) bool) {
	b.syncMap.Range(fn)
}

// Bearers returns all the bearers registered in Session.
func (s *Session) Bearers() []*Bearer {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	var bs []*Bearer
	s.bearerMap.rangeWithFunc(func(k, v interface{}) bool {
		bs = append(bs, v.(*Bearer))
		return true
	})

	return bs
}

// BearerCount returns the number of bearers registered in Session.
func (s *Session) BearerCount() int {
	if s == nil {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	s.bearerMap.rangeWithFunc(func(k, v interface{}) bool {
		count++
		return true
	})

	return count
}
