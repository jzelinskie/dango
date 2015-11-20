// Copyright 2015 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by the BSD 2-Clause license,
// which can be found in the LICENSE file.

package common

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/jzelinskie/dango"
	"github.com/jzelinskie/dango/mock"
)

func TestAnnounceTimer(t *testing.T) {
	handler := AnnounceTimer(func(ctx context.Context, w dango.AnnounceResponseWriter, r *dango.AnnounceRequest) (context.Context, error) {
		time.Sleep(time.Second)
		return ctx, nil
	})

	ctx := context.Background()
	ctx, _ = handler(ctx, &mock.AnnounceResponseWriter{}, &dango.AnnounceRequest{})
	duration := ctx.Value("time").(time.Duration)
	if duration < time.Second {
		t.Errorf("failed to properly time 1 second handler")
	}
}

func TestScrapeTimer(t *testing.T) {
	handler := ScrapeTimer(func(ctx context.Context, w dango.ScrapeResponseWriter, r *dango.ScrapeRequest) (context.Context, error) {
		time.Sleep(time.Second)
		return ctx, nil
	})

	ctx := context.Background()
	ctx, _ = handler(ctx, &mock.ScrapeResponseWriter{}, &dango.ScrapeRequest{})
	duration := ctx.Value("time").(time.Duration)
	if duration < time.Second {
		t.Errorf("failed to properly time 1 second handler")
	}
}
