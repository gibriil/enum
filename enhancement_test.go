// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum_test

import (
	"testing"

	"github.com/gibriil/enum"
)

type registration uint

const (
	seen registration = 1 << iota
	caught
)

type state struct {
	enum.MemberAs[registration]

	NickName string
	Type     string
}

type tamerLog struct {
	Turtle        state
	Lizard        state
	SeedFrog      state
	ElectricMouse state
}

func TestMemberAsEmbedding(t *testing.T) {

	dex := enum.Define(tamerLog{
		Turtle: state{
			NickName: "Bluey",
			Type:     "water",
		},
		Lizard: state{
			NickName: "Dragard",
			Type:     "fire",
		},
		SeedFrog: state{
			NickName: "Plantit",
			Type:     "grass",
		},
		ElectricMouse: state{
			NickName: "eleberl",
			Type:     "electric",
		},
	})

	var _ enum.Enum = dex.Turtle
	var _ enum.Enum = dex.Lizard
	var _ enum.Enum = dex.SeedFrog
	var _ enum.Enum = dex.ElectricMouse
}
