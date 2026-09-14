// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitflag

import "github.com/gibriil/enum/internal"

type Store[T any] struct {
	*internal.Definition[T]
	state uint64
}

func (s Store[T]) NewState() Store[T] {
	return Store[T]{
		Definition: s.Definition,
	}
}
