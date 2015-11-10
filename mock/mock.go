// Copyright 2015 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by the BSD 2-Clause license,
// which can be found in the LICENSE file.

// Package mock implements mocks useful for testing otoshi middleware.
package mock

import "github.com/jzelinskie/otoshi"

// AnnounceResponseWriter implements an otoshi.AnnounceResponseWriter.
type AnnounceResponseWriter struct {
	Resp *otoshi.AnnounceResponse
	Err  error
}

// WriteAnnounceResponse saves the provided response to the Resp field of the
// writer.
func (w *AnnounceResponseWriter) WriteAnnounceResponse(r *otoshi.AnnounceResponse) {
	w.Resp = r
}

// WriteError saves the provided error to the Err field of the writer.
func (w *AnnounceResponseWriter) WriteError(err error) {
	w.Err = err
}
