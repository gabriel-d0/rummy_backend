package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// Day 55 — Meld test matrix: regression corpus for sets/runs/ace/joker
func TestMeldMatrix(t *testing.T) {
	type tc struct {
		name  string
		meld  Meld
		valid bool // true if Validate should pass (run or set accordingly)
		isRun bool
	}
	// Helper to create set melds
	newSet := func(id string, tiles []tile.TileInstance, reps map[tile.TileInstanceId]tile.TileInstance) Meld {
		m, _ := New(MeldID(id), KindSet, tiles, reps)
		return m
	}
	newRun := func(id string, tiles []tile.TileInstance, reps map[tile.TileInstanceId]tile.TileInstance) Meld {
		m, _ := New(MeldID(id), KindRun, tiles, reps)
		return m
	}

	cases := []tc{
		// Sets
		{"set 3 valid", newSet("m1", []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Yellow, 5), tile.MustTile("t3", tile.Blue, 5)}, nil), true, false},
		{"set 4 valid", newSet("m2", []tile.TileInstance{tile.MustTile("t1", tile.Red, 10), tile.MustTile("t2", tile.Yellow, 10), tile.MustTile("t3", tile.Blue, 10), tile.MustTile("t4", tile.Black, 10)}, nil), true, false},
		{"set duplicate colour", newSet("m3", []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 5), tile.MustTile("t3", tile.Blue, 5)}, nil), false, false},
		{"set rank mismatch", newSet("m4", []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Yellow, 5), tile.MustTile("t3", tile.Blue, 6)}, nil), false, false},
		{"set joker valid 3 with 1", func() Meld {
			j := tile.MustJoker("j1")
			return newSet("m5", []tile.TileInstance{tile.MustTile("t1", tile.Red, 7), tile.MustTile("t2", tile.Yellow, 7), j}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("rep", tile.Blue, 7)})
		}(), true, false},
		{"set joker ratio 3 with 2", func() Meld {
			j1, j2 := tile.MustJoker("j1"), tile.MustJoker("j2")
			return newSet("m6", []tile.TileInstance{tile.MustTile("t1", tile.Red, 7), j1, j2}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("r1", tile.Yellow, 7), "j2": tile.MustTile("r2", tile.Blue, 7)})
		}(), false, false},

		// Runs
		{"run 5-6-7 valid", newRun("m7", []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 6), tile.MustTile("t3", tile.Red, 7)}, nil), true, true},
		{"run unsorted 7-5-6 valid", newRun("m8", []tile.TileInstance{tile.MustTile("t3", tile.Red, 7), tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 6)}, nil), true, true},
		{"run different colours", newRun("m9", []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Blue, 6), tile.MustTile("t3", tile.Red, 7)}, nil), false, true},
		{"run gap 5-7-8", newRun("m10", []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 7), tile.MustTile("t3", tile.Red, 8)}, nil), false, true},
		{"run duplicate rank", newRun("m11", []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 6), tile.MustTile("t3", tile.Red, 6)}, nil), false, true},

		// Ace
		{"low Ace 1-2-3", newRun("m12", []tile.TileInstance{tile.MustTile("t1", tile.Blue, 1), tile.MustTile("t2", tile.Blue, 2), tile.MustTile("t3", tile.Blue, 3)}, nil), true, true},
		{"high Ace 12-13-1", newRun("m13", []tile.TileInstance{tile.MustTile("t1", tile.Red, 12), tile.MustTile("t2", tile.Red, 13), tile.MustTile("t3", tile.Red, 1)}, nil), true, true},
		{"invalid Ace 13-1-2", newRun("m14", []tile.TileInstance{tile.MustTile("t1", tile.Red, 13), tile.MustTile("t2", tile.Red, 1), tile.MustTile("t3", tile.Red, 2)}, nil), false, true},
		{"invalid Ace middle 11-1-13", newRun("m15", []tile.TileInstance{tile.MustTile("t1", tile.Red, 11), tile.MustTile("t2", tile.Red, 1), tile.MustTile("t3", tile.Red, 13)}, nil), false, true},

		// Joker runs
		{"run joker gap 5-7 joker 6", func() Meld {
			j := tile.MustJoker("j1")
			return newRun("m16", []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 7), j}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("rep", tile.Red, 6)})
		}(), true, true},
		{"run joker wrong colour", func() Meld {
			j := tile.MustJoker("j1")
			return newRun("m17", []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 7), j}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("rep", tile.Blue, 6)})
		}(), false, true},
		{"run joker ratio 5 tiles 2 jokers 3 real fail", func() Meld {
			j1, j2 := tile.MustJoker("j1"), tile.MustJoker("j2")
			return newRun("m18", []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 6), tile.MustTile("t3", tile.Red, 7), j1, j2}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("r1", tile.Red, 8), "j2": tile.MustTile("r2", tile.Red, 9)})
		}(), false, true},
	}

	for _, tc := range cases {
		var err error
		if tc.isRun {
			err = ValidateRun(tc.meld)
		} else {
			err = ValidateSet(tc.meld)
		}
		if tc.valid && err != nil {
			t.Errorf("%s should be valid, got %v", tc.name, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("%s should be invalid, got nil", tc.name)
		}
	}
}
