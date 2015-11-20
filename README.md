![](https://i.imgur.com/svvnsus.jpg)

[![GoDoc](https://godoc.org/github.com/jzelinskie/dango?status.svg)](https://godoc.org/github.com/jzelinskie/dango)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://en.wikipedia.org/wiki/BSD_licenses#2-clause_license_.28.22Simplified_BSD_License.22_or_.22FreeBSD_License.22.29)

**dango** (named after the Japanese food *mitarashi dango*) is an experimental framework for building BitTorrent trackers with customizable behavior.
Conceptually, dango is modeled after the layers of middleware often seen implementing reusable behavior in HTTP servers; a tracker is simply the composition of two lists of middleware: one for Announces and one for Scrapes.

This library uses terminology defined in the [BitTorrent protocol specification] and the [popularized jargon] used on an Internet near you.

[BitTorrent protocol specification]: http://www.bittorrent.org/beps/bep_0003.html
[popularized jargon]: https://en.wikipedia.org/wiki/Glossary_of_BitTorrent_terms
