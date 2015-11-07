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

// Infohash is the hash of a set of files that are to be downloaded by an
// individual Swarm.
type Infohash string

// Swarm represents the metadata of collection of Peers sharing a torrent.
type Swarm struct{ Complete, Incomplete, Downloaded uint64 }

// AnnounceIntervals represents the intervals at which a BitTorrent client is
// intended to announce to a Tracker.
type AnnounceIntervals struct{ AnnounceInterval, MinAnnounceInterval time.Duration }

// AnnounceRequest represents the transport-agnostic Announce request sent
// from a BitTorrent client to a Tracker.
type AnnounceRequest struct {
	Peer
	Infohash
	URL        string
	Event      string
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

// AnnounceMiddleware is any function that operates on an Announce before a
// response has been written.
type AnnounceMiddleware func(context.Context, *AnnounceRequest, *AnnounceResponse) error

// ScrapeRequest represents a transport-agnostic Scrape request sent from a
// BitTorrent client to a tracker.
type ScrapeRequest []Infohash

// ScrapeResponse represents the transport-agnostic Scrape response sent from
// a tracker to a BitTorrent client.
type ScrapeResponse map[Infohash]Swarm

// ScrapeMiddleware is any function that operates on a Scrape before a response
// has been written.
type ScrapeMiddleware func(*ScrapeRequest, *ScrapeResponse) error

// TransportWriter abstracts the details of writing a response over a specific
// protocol.
type TransportWriter interface {
	WriteError(error) error
	WriteAnnounceResponse(*AnnounceRequest, *AnnounceResponse) error
	WriteScrapeResponse(*ScrapeRequest, *ScrapeResponse) error
}

// Tracker represents a BitTorrent tracker as the composition of
// AnnounceMiddleware and ScrapeMiddleware.
type Tracker struct {
	AnnnounceMiddleware []AnnounceMiddleware
	ScrapeMiddleWare    []ScrapeMiddleware
}

// ServeAnnounce runs a parsed AnnounceRequest through all of the tracker
// middleware and writes a response.
func (t *Tracker) ServeAnnounce(req AnnounceRequest, w TransportWriter) error {
	var resp AnnounceResponse
	for _, middleware := range t.AnnounceMiddleware {
		err := middleware(&req, &resp)
		if IsClientError(err) {
			w.WriteError(err)
			return nil
		} else if err != nil {
			return err
		}
	}

	err = w.WriteAnnounceResponse(&req, &resp)
	if err != nil {
		return err
	}
}

// ServeScrape runs a parsed ScrapeRequest through all of the tracker
// middleware and writes a response.
func (t *Tracker) ServeScrape(req ScrapeRequest, w TransportWriter) error {
	var resp ScrapeResponse
	for _, middleware := range t.ScrapeMiddleware {
		err := middleware(&req, &resp)
		if IsClientError(err) {
			w.WriteError(err)
			return nil
		} else if err != nil {
			return err
		}
	}

	err = w.WriteScrapeResponse(&req, &resp)
	if err != nil {
		return err
	}
}
