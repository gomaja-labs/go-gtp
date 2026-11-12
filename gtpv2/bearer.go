// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv2

import (
	"net"
	"sync"

	"github.com/gomaja-labs/go-gtp/gtpv2/ie"
)

// QoSProfile represents a QoS-related information that belongs to a Bearer.
type QoSProfile struct {
	PCI, PVI bool
	PL, QCI  uint8
	// Max bit rate for Uplink and Downlink
	MBRUL, MBRDL uint64
	// Guaranteed bit rate for Uplink and Downlink
	GBRUL, GBRDL uint64
}

// Bearer represents a GTPv2 bearer.
type Bearer struct {
	*teidMap
	AdditionalIEs   []*ie.IE
	mu              sync.RWMutex
	raddr           net.Addr
	teidIn, teidOut uint32

	EBI               uint8
	SubscriberIP, APN string
	ChargingID        uint32
	*QoSProfile
}

// NewBearer creates a new Bearer.
func NewBearer(ebi uint8, apn string, qos *QoSProfile) *Bearer {
	if qos == nil {
		qos = &QoSProfile{}
	}
	return &Bearer{
		EBI: ebi, APN: apn, QoSProfile: qos, teidMap: newTeidMap(),
	}
}

// RemoteAddress returns the remote address associated with Bearer.
func (b *Bearer) RemoteAddress() net.Addr {
	if b == nil {
		return nil
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	return copyAddr(b.raddr)
}

// SetRemoteAddress sets the remote address associated with Bearer.
func (b *Bearer) SetRemoteAddress(raddr net.Addr) {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.raddr = copyAddr(raddr)
}

// IncomingTEID returns the incoming TEID associated with Bearer.
func (b *Bearer) IncomingTEID() uint32 {
	if b == nil {
		return 0
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.teidIn
}

// SetIncomingTEID sets the incoming TEID associated with Bearer.
func (b *Bearer) SetIncomingTEID(teid uint32) {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.teidIn = teid
}

// OutgoingTEID returns the outgoing TEID associated with Bearer.
func (b *Bearer) OutgoingTEID() uint32 {
	if b == nil {
		return 0
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.teidOut
}

// SetOutgoingTEID sets the outgoing TEID associated with Bearer.
func (b *Bearer) SetOutgoingTEID(teid uint32) {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.teidOut = teid
}

// AddTEID records the F-TEID of one interface within this bearer.
func (b *Bearer) AddTEID(ifType uint8, teid uint32) {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.teidMap == nil {
		b.teidMap = newTeidMap()
	}
	b.teidMap.store(ifType, teid)
}

// GetTEID returns the F-TEID for an interface within this bearer.
func (b *Bearer) GetTEID(ifType uint8) (uint32, error) {
	if b == nil {
		return 0, ErrTEIDNotFound
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.teidMap != nil {
		if teid, ok := b.teidMap.load(ifType); ok {
			return teid, nil
		}
	}
	return 0, ErrTEIDNotFound
}
