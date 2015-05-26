// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package triegen

// This file defines Compacter and its implementations.

import "io"

// A Compacter generates an alternative, more space-efficient way to store a
// trie value block. A trie value block holds all possible values for the last
// byte of a UTF-8 encoded rune. Excluding ASCII characters, a trie value block
// always has 64 values, as a UTF-8 encoding ends with a byte in [0x80, 0xC0).
type Compacter interface {
    // Size returns whether the Compacter could encode the given block as well
    // as its size in case it can. len(v) is always 64.
    Size(v []uint64) (sz int, ok bool)

    // Store stores the block using the Compacter's compression method.
    // It returns a handle with which the block can be retrieved.
    // len(v) is always 64.
    Store(v []uint64) uint32

    // Print writes the data structures associated to the given store to w.
    Print(w io.Writer) error

    // Handler returns the name of a function that gets called during trie
    // lookup for blocks generated by the Compacter. The function should be of
    // the form func (n uint32, b byte) uint64, where n is the index returned by
    // the Compacter's Store method and b is the last byte of the UTF-8
    // encoding, where 0x80 <= b < 0xC0, for which to do the lookup in the
    // block.
    Handler() string
}

// simpleCompacter is the default Compacter used by builder. It implements a
// normal trie block.
type simpleCompacter builder

func (b *simpleCompacter) Size([]uint64) (sz int, ok bool) {
    return blockSize * b.ValueSize, true
}

func (b *simpleCompacter) Store(v []uint64) uint32 {
    h := uint32(len(b.ValueBlocks) - blockOffset)
    b.ValueBlocks = append(b.ValueBlocks, v)
    return h
}

func (b *simpleCompacter) Print(io.Writer) error {
    // Structures are printed in print.go.
    return nil
}

func (b *simpleCompacter) Handler() string {
    panic("Handler should be special-cased for this Compacter")
}