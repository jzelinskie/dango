// Copyright 2015 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by the BSD 2-Clause license,
// which can be found in the LICENSE file.

// Package otoshi defines a middleware abstraction for the composition of
// BitTorrent tracker functionality.
package otoshi

import (
	"net"
	"time"

	"golang.org/x/net/context"
)

// Peer represents one instance of a BitTorrent client currently participating
// in a Swarm.
type Peer struct {
	ID   string
	IP   net.IP
	Port uint32
	Left uint64
}

// PeerIterator represents a stream of Peers.
type PeerIterator interface {
	Next() bool
	Peer() Peer
}

// Event represents an event done by a BitTorrent client.
type Event string

const (
	// Started is the event sent by a BitTorrent client when it joins a Swarm.
	Started Event = "started"

	// Stopped is the event sent by a BitTorrent client when it leaves a Swarm.
	Stopped Event = "stopped"

	// Completed is the event sent by a BitTorrent client when it finishes
	// downloading all of the required chunks.
	Completed Event = "completed"
)

// Infohash is the hash of a set of files that are to be downloaded by a client
// participating in a Swarm.
type Infohash string

// Swarm represents the metadata of collection of Peers sharing a torrent.
type Swarm struct{ Complete, Incomplete, Downloaded uint64 }

// AnnounceIntervals represents the intervals at which a BitTorrent client is
// intended to announce to a Tracker.
type AnnounceIntervals struct{ AnnounceInterval, MinAnnounceInterval time.Duration }

// AnnounceRequest represents the transport-agnostic Announce request sent
// from a BitTorrent client to a Tracker.
type AnnounceRequest struct {
	Peer       Peer
	Infohash   Infohash
	Event      Event
	URL        string
	Downloaded uint64
	Uploaded   uint64
	Compact    bool
	NumWant    uint16
}

// AnnounceResponse represents the transport-agnostic Announce response sent
// from a Tracker to a BitTorrent client.
type AnnounceResponse struct {
	IPv4Peers, IPv6Peers PeerIterator
	AnnounceIntervals
	Swarm
}

// AnnounceResponseWriter is used by an AnnounceHandler to construct a response.
type AnnounceResponseWriter interface {
	WriteAnnounceResponse(*AnnounceResponse)
	WriteError(error)
}

// AnnounceHandler is a function that operates on an Announce before a response
// has been written.
type AnnounceHandler func(context.Context, AnnounceResponseWriter, *AnnounceRequest) (context.Context, error)

// AnnounceMiddleware enables the extension of an AnnounceHandler via a closure.
type AnnounceMiddleware func(AnnounceHandler) AnnounceHandler

// ScrapeRequest represents a transport-agnostic Scrape request sent from a
// BitTorrent client to a tracker.
type ScrapeRequest []Infohash

// ScrapeResponse represents the transport-agnostic Scrape response sent from
// a tracker to a BitTorrent client.
type ScrapeResponse map[Infohash]Swarm

// ScrapeResponseWriter is used by a ScrapeHandler to construct a response.
type ScrapeResponseWriter interface {
	WriteScrapeResponse(*ScrapeResponse)
	WriteError(error)
}

// ScrapeHandler is a function that operates on a Scrape before a response
// has been written.
type ScrapeHandler func(context.Context, ScrapeResponseWriter, *ScrapeRequest) (context.Context, error)

// ScrapeMiddleware enables the extension of an ScrapeHandler via a closure.
type ScrapeMiddleware func(ScrapeHandler) ScrapeHandler

// AnnounceChain represents a composition of AnnounceMiddleware.
type AnnounceChain struct{ mw []AnnounceMiddleware }

// Append creates a new AnnounceChain by appending the provided middleware to
// the current chain.
func (c AnnounceChain) Append(mw ...AnnounceMiddleware) (nc AnnounceChain) {
	nc.mw = make([]AnnounceMiddleware, len(c.mw)+len(mw))
	copy(nc.mw[:len(c.mw)], c.mw)
	copy(nc.mw[len(c.mw):], mw)
	return
}

// Finalize creates an AnnounceHandler that is the composition of the chain.
func (c AnnounceChain) Finalize(final AnnounceHandler) AnnounceHandler {
	for i := len(c.mw) - 1; i >= 0; i-- {
		final = c.mw[i](final)
	}
	return final
}

// ScrapeChain represents a composition of ScrapeMiddleware.
type ScrapeChain struct{ mw []ScrapeMiddleware }

// Append creates a new ScrapeChain by appending the provided middleware to the
// current chain.
func (c ScrapeChain) Append(mw ...ScrapeMiddleware) (nc ScrapeChain) {
	nc.mw = make([]ScrapeMiddleware, len(c.mw)+len(mw))
	copy(nc.mw[:len(c.mw)], c.mw)
	copy(nc.mw[len(c.mw):], mw)
	return
}

// Finalize creates an ScrapeHandler that is the composition of the chain.
func (c ScrapeChain) Finalize(final ScrapeHandler) ScrapeHandler {
	for i := len(c.mw) - 1; i >= 0; i-- {
		final = c.mw[i](final)
	}
	return final
}

// ServeAnnounce creates an empty context and calls itself.
func (h AnnounceHandler) ServeAnnounce(w AnnounceResponseWriter, r *AnnounceRequest) {
	ctx := context.Background()
	h(ctx, w, r)
}

// ServeScrape creates an empty context and calls itself.
func (h ScrapeHandler) ServeScrape(w ScrapeResponseWriter, r *ScrapeRequest) {
	ctx := context.Background()
	h(ctx, w, r)
}

// NewTracker returns a Tracker given an AnnounceHandler and a ScrapeHandler.
func NewTracker(ah AnnounceHandler, sh ScrapeHandler) Tracker {
	return struct {
		AnnounceHandler
		ScrapeHandler
	}{
		AnnounceHandler: ah,
		ScrapeHandler:   sh,
	}
}

// Tracker represents a BitTorrent tracker.
type Tracker interface {
	ServeAnnounce(AnnounceResponseWriter, *AnnounceRequest)
	ServeScrape(ScrapeResponseWriter, *ScrapeRequest)
}
