// Copyright 2015 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by the BSD 2-Clause license,
// which can be found in the LICENSE file.

// Package mock implements mocks useful for testing dango middleware.
package mock

import "github.com/jzelinskie/dango"

// AnnounceResponseWriter implements an dango.AnnounceResponseWriter in the
// simplest possible way.
type AnnounceResponseWriter struct {
	Resp *dango.AnnounceResponse
	Err  error
}

// WriteAnnounceResponse saves the provided response to the Resp field of the
// writer.
func (w *AnnounceResponseWriter) WriteAnnounceResponse(r *dango.AnnounceResponse) {
	w.Resp = r
}

// WriteError saves the provided error to the Err field of the writer.
func (w *AnnounceResponseWriter) WriteError(err error) {
	w.Err = err
}

// ScrapeResponseWriter implements an dango.ScrapeResponseWriter in the
// simplest possible way.
type ScrapeResponseWriter struct {
	Resp *dango.ScrapeResponse
	Err  error
}

// WriteScrapeResponse saves the provided response to the Resp field of the
// writer.
func (w *ScrapeResponseWriter) WriteScrapeResponse(r *dango.ScrapeResponse) {
	w.Resp = r
}

// WriteError saves the provided error to the Err field of the writer.
func (w *ScrapeResponseWriter) WriteError(err error) {
	w.Err = err
}
