// Copyright 2015 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by the BSD 2-Clause license,
// which can be found in the LICENSE file.

package otoshi

// ClientError is an error that should be written as a response to a BitTorrent
// client.
type ClientError string

func (e ClientError) Error() string { return string(e) }

// IsClientError determines if an error should be propogated to the client.
func IsClientError(err error) bool {
	if err == nil {
		return false
	}

	_, cl := err.(ClientError)
	return cl
}
