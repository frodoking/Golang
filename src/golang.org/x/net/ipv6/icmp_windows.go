// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

type sysICMPv6Filter struct {
    // TODO(mikio): implement this
}

func (f *sysICMPv6Filter) set(typ ICMPType, block bool) {
    // TODO(mikio): implement this
}

func (f *sysICMPv6Filter) setAll(block bool) {
    // TODO(mikio): implement this
}

func (f *sysICMPv6Filter) willBlock(typ ICMPType) bool {
    // TODO(mikio): implement this
    return false
}
