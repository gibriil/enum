// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

import (
	"database/sql/driver"
	"encoding"
	"testing"
)

var (
	_ Enum                     = color{}
	_ encoding.TextMarshaler   = color{}
	_ encoding.TextUnmarshaler = (*color)(nil)
	_ driver.Valuer            = (*color)(nil)
)

type color struct {
	Member
}

type colors struct {
	Red   color
	Green color
	Blue  color
}

func TestColorsIndexes(t *testing.T) {
	var Colors = Define(colors{
		Red:   color{},
		Green: color{},
		Blue:  color{},
	})

	if r := Colors.Red.Index(); r != 0 {
		t.Errorf("Red.Index(): got %d, want %d", r, 0)
	}

	if g := Colors.Green.Index(); g != 1 {
		t.Errorf("Green.Index(): got %d, want %d", g, 1)
	}

	if b := Colors.Blue.Index(); b != 2 {
		t.Errorf("Blue.Index(): got %d, want %d", b, 2)
	}
}

func TestColorsNames(t *testing.T) {
	var Colors = Define(colors{
		Red:   color{},
		Green: color{},
		Blue:  color{},
	})

	if r := Colors.Red.Name(); r != "Red" {
		t.Errorf("Red.Name(): got %s, want %s", r, "Red")
	}

	if g := Colors.Green.Name(); g != "Green" {
		t.Errorf("Green.Name(): got %s, want %s", g, "Green")
	}

	if b := Colors.Blue.Name(); b != "Blue" {
		t.Errorf("Blue.Name(): got %s, want %s", b, "Blue")
	}
}

func TestColorValuesAreDistinct(t *testing.T) {
	var Colors = Define(colors{
		Red:   color{},
		Green: color{},
		Blue:  color{},
	})

	if Colors.Red == Colors.Green {
		t.Error("Red and Green should not be equal")
	}

	if Colors.Red == Colors.Blue {
		t.Error("Red and Blue should not be equal")
	}

	if Colors.Green == Colors.Blue {
		t.Error("Green and Blue should not be equal")
	}
}
