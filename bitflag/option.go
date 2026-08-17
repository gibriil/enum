// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitflag

type Option struct {
	def   *definition
	index int
	mask  uint64
}

// initializer is a non-exported interface for reflection type safety
type initializer interface {
	initialize(*definition, int, uint64)
}

// initialize initializes the flag option with its closed set definition
// and sets its position index in the list
func (f *Option) initialize(def *definition, index int, mask uint64) {
	f.def = def
	f.index = index
	f.mask = mask
}

// flag marks Option as a valid flag implementation.
// It intentionally has no behavior; it seals the Flag interface.
func (f Option) flag() {}

type Flag interface {
	flag()
}
