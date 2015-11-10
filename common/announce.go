// Copyright 2015 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by the BSD 2-Clause license,
// which can be found in the LICENSE file.

package common

import (
	"time"

	"golang.org/x/net/context"

	"github.com/jzelinskie/otoshi"
)

func AnnounceTimer(next otoshi.AnnounceHandler) otoshi.AnnounceHandler {
	return func(ctx context.Context, w otoshi.AnnounceResponseWriter, r *otoshi.AnnounceRequest) (context.Context, error) {
		var err error
		start := time.Now()
		ctx, err = next(ctx, w, r)
		end := time.Now()
		ctx = context.WithValue(ctx, "time", end.Sub(start))
		return ctx, err
	}
}
