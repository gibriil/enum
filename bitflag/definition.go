// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitflag

import "github.com/gibriil/enum/internal"

type Store struct {
	def   *internal.Definition
	state uint64
}

func (def definition) NewState()
