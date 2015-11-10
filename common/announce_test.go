// Copyright 2015 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by the BSD 2-Clause license,
// which can be found in the LICENSE file.

package common

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/jzelinskie/otoshi"
)

type mockAnnounceResponseWriter struct {
	resp *otoshi.AnnounceResponse
	err  error
}

func (w *mockAnnounceResponseWriter) WriteAnnounceResponse(r *otoshi.AnnounceResponse) {
	w.resp = r
}

func (w *mockAnnounceResponseWriter) WriteError(err error) {
	w.err = err
}

func TestAnnounceTimer(t *testing.T) {
	handler := AnnounceTimer(func(ctx context.Context, w otoshi.AnnounceResponseWriter, r *otoshi.AnnounceRequest) (context.Context, error) {
		time.Sleep(time.Second)
		return ctx, nil
	})

	ctx := context.Background()
	ctx, _ = handler(ctx, &mockAnnounceResponseWriter{}, &otoshi.AnnounceRequest{})
	duration := ctx.Value("time").(time.Duration)
	if duration < time.Second {
		t.Errorf("failed to properly time 1 second handler")
	}
}
