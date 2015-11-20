// Copyright 2015 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by the BSD 2-Clause license,
// which can be found in the LICENSE file.

// Package common implements common AnnounceMiddleware and ScrapeMiddleware.
package common

import (
	"time"

	"golang.org/x/net/context"

	"github.com/jzelinskie/dango"
)

// AnnounceTimer times the handling of Announces and returns a context with the
// key "time" set to a time.Duration.
func AnnounceTimer(next dango.AnnounceHandler) dango.AnnounceHandler {
	return func(ctx context.Context, w dango.AnnounceResponseWriter, r *dango.AnnounceRequest) (context.Context, error) {
		var err error
		start := time.Now()
		ctx, err = next(ctx, w, r)
		end := time.Now()
		ctx = context.WithValue(ctx, "time", end.Sub(start))
		return ctx, err
	}
}

// ScrapeTimer times the handling of Scrapes and returns a context with the
// key "time" set to a time.Duration.
func ScrapeTimer(next dango.ScrapeHandler) dango.ScrapeHandler {
	return func(ctx context.Context, w dango.ScrapeResponseWriter, r *dango.ScrapeRequest) (context.Context, error) {
		var err error
		start := time.Now()
		ctx, err = next(ctx, w, r)
		end := time.Now()
		ctx = context.WithValue(ctx, "time", end.Sub(start))
		return ctx, err
	}
}
