package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// Day 51 — Invalid ace-middle runs: 13-1-2 and similar must be rejected
// per docs/rules-decisions.md:1.3 “Aces cannot appear in the middle of a run”.
func TestValidateRunInvalidAceMiddle(t *testing.T) {
	cases := []struct {
		name  string
		ranks []tile.Rank
	}{
		{"13-1-2", []tile.Rank{13, 1, 2}},
		{"12-1-2", []tile.Rank{12, 1, 2}},
		{"13-1-3", []tile.Rank{13, 1, 3}},
		{"11-1-13", []tile.Rank{11, 1, 13}},
		{"2-13-1", []tile.Rank{2, 13, 1}},
	}
	for _, tc := range cases {
		tiles := make([]tile.TileInstance, len(tc.ranks))
		for i, r := range tc.ranks {
			tiles[i] = tile.MustTile(tile.TileInstanceId("t"+string(rune('0'+i))), tile.Red, r)
		}
		m, _ := New("m1", KindRun, tiles, nil)
		if err := ValidateRun(m); err == nil {
			t.Fatalf("%s should be invalid (Ace in middle)", tc.name)
		}
	}
	// Ensure low and high remain valid (sanity)
	validLow, _ := New("m1", KindRun, []tile.TileInstance{tile.MustTile("t1", tile.Blue, 1), tile.MustTile("t2", tile.Blue, 2), tile.MustTile("t3", tile.Blue, 3)}, nil)
	if err := ValidateRun(validLow); err != nil {
		t.Fatalf("1-2-3 should still be valid: %v", err)
	}
	validHigh, _ := New("m1", KindRun, []tile.TileInstance{tile.MustTile("t1", tile.Red, 12), tile.MustTile("t2", tile.Red, 13), tile.MustTile("t3", tile.Red, 1)}, nil)
	if err := ValidateRun(validHigh); err != nil {
		t.Fatalf("12-13-1 should still be valid: %v", err)
	}
}
