// Copyright 2026 Peter James Beard. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum_test

import (
	"database/sql/driver"
	"encoding"
	"strings"
	"testing"

	"github.com/gibriil/enum"
)

type itemTier struct {
	enum.Member

	// Categories defines the allowed item types or equipment slots for this tier.
	Categories []string

	// Effects maps item identifiers to their stat bonuses or special abilities.
	Effects map[string]string
}

// Ensures itemTier type satisfies Enum interface
var (
	_ enum.Enum              = itemTier{}
	_ encoding.TextMarshaler = itemTier{}
	_ driver.Valuer          = (*itemTier)(nil)
)

// RarityTiers groups item configurations by rarity level
var rarityTiers = struct {
	Common    itemTier
	Uncommon  itemTier
	Rare      itemTier
	Epic      itemTier
	Legendary itemTier
}{
	// Common: basic starter equipment
	Common: itemTier{
		Categories: []string{"weapon", "armor", "consumable"},
		Effects: map[string]string{
			"iron_sword":    "damage+2",
			"leather_armor": "defense+1",
			"health_potion": "heal+25",
			"torch":         "light",
		},
	},

	// Uncommon: improved gear with minor stat bumps
	Uncommon: itemTier{
		Categories: []string{"weapon", "armor", "accessory", "consumable"},
		Effects: map[string]string{
			"steel_sword":   "damage+5",
			"chain_mail":    "defense+3",
			"health_potion": "heal+50",
			"mana_potion":   "mana+30",
			"ring_of_sight": "detection+5",
		},
	},

	// Rare: powerful equipment with combined effects
	Rare: itemTier{
		Categories: []string{"weapon", "armor", "accessory", "consumable", "reagent"},
		Effects: map[string]string{
			"flame_blade":      "damage+10,fire_damage+5",
			"plate_armor":      "defense+8",
			"health_potion":    "heal+100",
			"elixir":           "regeneration+2",
			"amulet_of_wisdom": "mp_regen+1",
		},
	},

	// Epic: master-tier gear with advanced abilities
	Epic: itemTier{
		Categories: []string{"weapon", "armor", "accessory", "reagent", "quest_item"},
		Effects: map[string]string{
			"storm_edge":      "damage+20,lightning_damage+10",
			"mythic_plate":    "defense+15,resist_magic+5",
			"potion_of_aging": "heal+200,age_backward",
			"rune_of_power":   "all_stats+3",
		},
	},

	// Legendary: game-changing artifacts
	Legendary: itemTier{
		Categories: []string{"weapon", "armor", "accessory", "quest_item"},
		Effects: map[string]string{
			"excalibur":       "damage+50,authenticity",
			"dragon_armor":    "defense+30,fire_resist_max",
			"ring_of_wish":    "wish_count+1",
			"phoenix_feather": "revival",
		},
	},
}

func init() {
	rarityTiers = enum.Define(rarityTiers)
}

func TestCollectionsByNameAs(t *testing.T) {
	tests := []struct {
		Name    string
		Lookup  string
		Want    itemTier
		wantErr bool
	}{
		{"Common", "Common", rarityTiers.Common, false},
		{"Uncommon", "Uncommon", rarityTiers.Uncommon, false},
		{"Rare", "Rare", rarityTiers.Rare, false},
		{"Epic", "Epic", rarityTiers.Epic, false},
		{"Legendary", "Legendary", rarityTiers.Legendary, false},
		{"Missing", "Mythic", itemTier{}, true},
	}

	for _, test := range tests {
		variations := []struct {
			Name    string
			Lookup  string
			wantErr bool
		}{
			{Name: "exact", Lookup: test.Lookup, wantErr: test.wantErr},
			{Name: "lowercase", Lookup: strings.ToLower(test.Lookup), wantErr: true},
			{Name: "uppercase", Lookup: strings.ToUpper(test.Lookup), wantErr: true},
		}

		for _, v := range variations {
			t.Run(test.Name+"/"+v.Name, func(t *testing.T) {
				got, err := enum.ByNameAs[itemTier](rarityTiers, v.Lookup)

				if err != nil && !v.wantErr {
					t.Fatalf(
						"ByNameAs[%s](%q): err=%v, wantErr=%v",
						rarityTiers.Rare.Type().Name(),
						v.Lookup,
						err,
						v.wantErr,
					)
				}

				if v.wantErr {
					if got != nil {
						t.Errorf("ByNameAs[%s](%q): got %v, want nil",
							rarityTiers.Rare.Type().Name(),
							v.Lookup,
							got,
						)
					}
					return
				}

				if got == nil {
					t.Fatal("ByNameAs returned nil member without an error")
				}

				if !enum.Equal(*got, test.Want) {
					t.Errorf(
						"ByNameAs[%s](%q): got %v, want %v",
						got.Type().Name(),
						v.Lookup,
						got,
						test.Want,
					)
				}

				if got.Name() != v.Lookup {
					t.Errorf("Name() = %q, want %q", got.Name(), v.Lookup)
				}
			})
		}
	}
}

func TestCollectionsByIndexAs(t *testing.T) {
	tests := []struct {
		Name    string
		Index   int
		Want    itemTier
		wantErr bool
	}{
		{"Common", 0, rarityTiers.Common, false},
		{"Uncommon", 1, rarityTiers.Uncommon, false},
		{"Rare", 2, rarityTiers.Rare, false},
		{"Epic", 3, rarityTiers.Epic, false},
		{"Legendary", 4, rarityTiers.Legendary, false},
		{"TooHigh", len(enum.Values(rarityTiers)), itemTier{}, true},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := enum.ByIndexAs[itemTier](rarityTiers, test.Index)

			if err != nil && !test.wantErr {
				t.Fatalf(
					"ByIndexAs[%s](%d): err=%v, wantErr=%v",
					rarityTiers.Rare.Type().Name(),
					test.Index,
					err,
					test.wantErr,
				)
			}

			if test.wantErr {
				if got != nil {
					t.Errorf("ByIndexAs[%s](%d): got %v, want nil",
						rarityTiers.Rare.Type().Name(),
						test.Index,
						got,
					)
				}
				return
			}

			if got == nil {
				t.Fatal("ByIndexAs returned nil member without an error")
			}

			if !enum.Equal(*got, test.Want) {
				t.Errorf(
					"ByIndexAs[%s](%d): got %v, want %v",
					got.Type().Name(),
					test.Index,
					got,
					test.Want,
				)
			}

			if got.Index() != test.Index {
				t.Errorf("Index() = %d, want %d", got.Index(), test.Index)
			}
		})
	}
}

func TestCollectionsValuesAs(t *testing.T) {
	tests := []struct {
		Name    string
		Index   int
		Want    itemTier
		wantErr bool
	}{
		{"Common", 0, rarityTiers.Common, false},
		{"Uncommon", 1, rarityTiers.Uncommon, false},
		{"Rare", 2, rarityTiers.Rare, false},
		{"Epic", 3, rarityTiers.Epic, false},
		{"Legendary", 4, rarityTiers.Legendary, false},
	}

	members, err := enum.ValuesAs[itemTier](rarityTiers)

	if err != nil {
		t.Fatalf(
			"ValuesAs[%s](): err=%v",
			rarityTiers.Rare.Type().Name(),
			err,
		)
	}

	if len(members) != len(tests) {
		t.Fatalf("unexpected Values length: got %d, want %d", len(tests), len(members))
	}

	for i, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			tier := members[i]
			if !enum.Equal(tier, test.Want) {
				t.Errorf("Value did not match expected enum: got %s, want %s", tier, test.Want)
			}
		})
	}
}

func TestCollectionsAllAs(t *testing.T) {
	tests := []struct {
		Name    string
		Index   int
		Want    itemTier
		wantErr bool
	}{
		{"Common", 0, rarityTiers.Common, false},
		{"Uncommon", 1, rarityTiers.Uncommon, false},
		{"Rare", 2, rarityTiers.Rare, false},
		{"Epic", 3, rarityTiers.Epic, false},
		{"Legendary", 4, rarityTiers.Legendary, false},
	}

	for e := range enum.AllAs[itemTier](rarityTiers) {
		t.Run(e.Name(), func(t *testing.T) {
			if !enum.Equal(*e, tests[e.Index()].Want) {
				t.Errorf("iteration did not yield expected enum: got %s, want %s", e, tests[e.Index()].Want)
			}
		})
	}
}
