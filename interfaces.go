// Copyright 2015 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by the BSD 2-Clause license,
// which can be found in the LICENSE file.

// Package dango defines a middleware abstraction for the composition of
// BitTorrent tracker functionality.
package dango

import (
	"errors"
	"net"
	"time"
)

var (
	// ErrMalformedRequest is returned when a request does not contain the
	// required parameters needed to create a model.
	ErrMalformedRequest = errors.New("malformed request")
)

// Peer represents one instance of a BitTorrent client currently participating
// in a Swarm.
type Peer interface {
	IP() net.IP
	Port() uint32
}

// PrivatePeer represents one instance of a BitTorrent client currently
// participating in a private tracker's Swarm.
type PrivatePeer interface {
	ID() string
	Peer
}

// PeerIterator represents a stream of Peers.
type PeerIterator interface {
	Next() bool
	Peer() Peer
}

// Event represents an event done by a BitTorrent client.
type Event uint8

const (
	// None is the event when a BitTorrent client announces due to time lapsed
	// since the previous announce.
	None Event = iota

	// Started is the event sent by a BitTorrent client when it joins a Swarm.
	Started

	// Stopped is the event sent by a BitTorrent client when it leaves a Swarm.
	Stopped

	// Completed is the event sent by a BitTorrent client when it finishes
	// downloading all of the required chunks.
	Completed
)

var (
	eventToString map[Event]string
	stringToEvent map[string]Event
)

func init() {
	eventToString = make(map[Event]string)
	eventToString[None] = "none"
	eventToString[Started] = "started"
	eventToString[Stopped] = "stopped"
	eventToString[Completed] = "completed"
	stringToEvent = make(map[string]Event)

	stringToEvent[""] = None
	stringToEvent["none"] = None
	stringToEvent["started"] = Started
	stringToEvent["stopped"] = Stopped
	stringToEvent["completed"] = Completed
}

// NewEvent returns the proper Event given a string.
func NewEvent(event string) (Event, error) {
	if e, ok := stringToEvent[event]; ok {
		return e, nil
	}

	return None, errors.New("dango: unknown event")
}

// String implements Stringer for an event.
func (e Event) String() string {
	if name, ok := eventToString[e]; ok {
		return name
	}

	panic("dango: Event has no associated name")
}

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
type AnnounceRequest interface {
	Peer() Peer
	Infohash() Infohash
	Event() Event
	URL() string
	Downloaded() uint64
	Uploaded() uint64
	Left() uint64
	Compact() bool
	NumWant() uint16
}

// PrivateAnnounceRequest represents an Announce request with extra metadata
// used for private BitTorrent trackers.
type PrivateAnnounceRequest interface {
	Passkey() string
	AnnounceRequest
}

// AnnounceResponse represents the transport-agnostic Announce response sent
// from a Tracker to a BitTorrent client.
type AnnounceResponse interface {
	IPv4Peers() PeerIterator
	IPv6Peers() PeerIterator
	AnnounceIntervals() AnnounceIntervals
	Swarm() Swarm
	Compact() bool
}

// ScrapeRequest represents a transport-agnostic Scrape request sent from a
// BitTorrent client to a tracker.
type ScrapeRequest interface {
	Infohashes() []Infohash
}

// PrivateScrapeRequest represents a Scrape request with extra metadata used
// for private BitTorrent trackers.
type PrivateScrapeRequest interface {
	Passkey() string
	ScrapeRequest
}

// ScrapeResponse represents the transport-agnostic Scrape response sent from
// a tracker to a BitTorrent client.
type ScrapeResponse interface {
	Files() map[Infohash]Swarm
}

// ErrorWriter is used to construct a response for a failure.
type ErrorWriter interface {
	WriteError(error) error
}

// AnnounceResponseWriter is used by an AnnounceHandler to construct a response.
type AnnounceResponseWriter interface {
	WriteAnnounceResponse(AnnounceResponse) error
	ErrorWriter
}

// ScrapeResponseWriter is used by a ScrapeHandler to construct a response.
type ScrapeResponseWriter interface {
	WriteScrapeResponse(ScrapeResponse) error
	ErrorWriter
}
