// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import "reflect"

// Metadata holds the reflection information for an enum
type Metadata struct {
	Name  string
	Field reflect.StructField
	Type  reflect.Type
	Value reflect.Value
}
