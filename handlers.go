package dango

import "golang.org/x/net/context"

// AnnounceHandler is a function that operates on an Announce before a response
// has been written.
type AnnounceHandler func(context.Context, AnnounceResponseWriter, AnnounceRequest) (context.Context, error)

// AnnounceMiddleware enables the extension of an AnnounceHandler via a closure.
type AnnounceMiddleware func(AnnounceHandler) AnnounceHandler

// ScrapeHandler is a function that operates on a Scrape before a response
// has been written.
type ScrapeHandler func(context.Context, ScrapeResponseWriter, ScrapeRequest) (context.Context, error)

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
func (h AnnounceHandler) ServeAnnounce(w AnnounceResponseWriter, r AnnounceRequest) {
	ctx := context.Background()
	h(ctx, w, r)
}

// ServeScrape creates an empty context and calls itself.
func (h ScrapeHandler) ServeScrape(w ScrapeResponseWriter, r ScrapeRequest) {
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
	ServeAnnounce(AnnounceResponseWriter, AnnounceRequest)
	ServeScrape(ScrapeResponseWriter, ScrapeRequest)
}
