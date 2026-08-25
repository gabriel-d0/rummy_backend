package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// Day 54 — Immutable joker mapping: a tabled joker's represented tile must not be silently reinterpreted.
func TestMeldImmutability(t *testing.T) {
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-6", tile.Red, 6)
	tiles := []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 7), j}
	reps := map[tile.TileInstanceId]tile.TileInstance{"j1": rep}
	m, _ := New("m1", KindRun, tiles, reps)

	// Mutate original slices/maps after New — meld should be unaffected
	tiles[0] = tile.MustTile("hack", tile.Blue, 1)
	reps["j1"] = tile.MustTile("hack2", tile.Blue, 1)
	if m.Tiles[0].ID == "hack" {
		t.Fatalf("Meld Tiles not copied — mutation leaked")
	}
	if m.JokerReps["j1"].ID == "hack2" {
		t.Fatalf("Meld JokerReps not copied")
	}
	// Direct mutation of returned Meld's map should not affect original (but caller could still mutate returned Meld's map if we didn't copy on Validate)
	// Our New copies, so mutating m.JokerReps after should not affect a second validation of same logical meld
	m.JokerReps["j1"] = tile.MustTile("hack3", tile.Blue, 2)
	// Create a fresh meld with original rep to ensure original rep still valid
	m2, _ := New("m2", KindRun, []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 7), j}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	if err := ValidateRun(m2); err != nil {
		t.Fatalf("fresh meld with original rep should still be valid: %v", err)
	}
	// The mutated m1 now has rep hack3 (Blue 2) which is wrong colour vs Red run, so it should fail
	if err := ValidateRun(m); err == nil {
		t.Fatalf("mutated rep should cause ValidateRun to fail (colour mismatch)")
	}
}

func TestJokerRepImmutabilityAfterValidate(t *testing.T) {
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-7", tile.Red, 7)
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	// Validate once
	if err := ValidateRun(m); err != nil {
		t.Fatalf("initial valid: %v", err)
	}
	// Simulate illegal silent reinterpretation: change rep to 8
	m.JokerReps["j1"] = tile.MustTile("rep-8", tile.Red, 8)
	// Now the run would be 5,6,8 not consecutive 5,6,7 → should fail
	if err := ValidateRun(m); err == nil {
		t.Fatalf("silent reinterpretation to 8 should fail consecutive check")
	}
}
