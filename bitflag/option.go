// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitflag

import (
	"fmt"

	"github.com/gibriil/enum/internal"
)

type Option[T any] struct {
	internal.Entry[T]
	mask uint64
}

// flag marks Option as a valid flag implementation.
// It intentionally has no behavior; it seals the Flag interface.
func (f Option[T]) flag() {}

type Flag[T any] interface {
	flag()

	fmt.Stringer

	Name() string
	Index() int
	Valid() bool
	OptionSet() Store[T]
}

func (f Option[T]) OptionSet() Set[T] {
	if
}
