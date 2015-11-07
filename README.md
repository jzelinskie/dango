![](http://i.imgur.com/NFYaeoK.jpg)

[![GoDoc](https://godoc.org/github.com/jzelinskie/otoshi?status.svg)](https://godoc.org/github.com/jzelinskie/otoshi)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://en.wikipedia.org/wiki/BSD_licenses#2-clause_license_.28.22Simplified_BSD_License.22_or_.22FreeBSD_License.22.29)

**otoshi** (named after the Japanese game *daruma otoshi*) is an experimental framework for building BitTorrent trackers with customizable behavior.
Conceptually, otoshi is modeled after the layers of middleware often seen implementing reusable behavior in HTTP servers; a tracker is simply the composition of two lists of middleware: one for Announces and one for Scrapes.
