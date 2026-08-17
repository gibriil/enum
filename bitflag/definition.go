// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitflag

import (
	"reflect"

	"github.com/gibriil/enum/internal"
)

type definition struct {
	identity   reflect.Type // type for closed set
	name       string       // name of closed set
	optionType reflect.Type // Flag Option Type
	length     int          // Number of struct members
	options    []Flag
	names      []string            // Slice of names for each initialized flag option
	lookup     map[string]int      //Lookup map for identifying the values index of a flag option by name
	metadata   []internal.Metadata // Slice of the reflection details of each flag option
}

type Store struct {
	def   *definition
	state uint64
}

func (def definition) NewState()
