package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// Day 49 — Low-ace runs: 1-2-3 and longer low-ace are valid and already handled
// as consecutive 1,2,3... by the basic sorter. This test documents the canonical low-Ace case.
func TestValidateRunLowAce(t *testing.T) {
	// 1-2-3
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Blue, 1),
		tile.MustTile("t2", tile.Blue, 2),
		tile.MustTile("t3", tile.Blue, 3),
	}, nil)
	if err := ValidateRun(m); err != nil {
		t.Fatalf("1-2-3 low Ace should be valid: %v", err)
	}
	// 1-2-3-4-5
	m2, _ := New("m2", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 1),
		tile.MustTile("t2", tile.Red, 2),
		tile.MustTile("t3", tile.Red, 3),
		tile.MustTile("t4", tile.Red, 4),
		tile.MustTile("t5", tile.Red, 5),
	}, nil)
	if err := ValidateRun(m2); err != nil {
		t.Fatalf("1-2-3-4-5 low Ace should be valid: %v", err)
	}
	// Unsorted 3-1-2 should still be valid after sort (1,2,3)
	m3, _ := New("m3", KindRun, []tile.TileInstance{
		tile.MustTile("t3", tile.Yellow, 3),
		tile.MustTile("t1", tile.Yellow, 1),
		tile.MustTile("t2", tile.Yellow, 2),
	}, nil)
	if err := ValidateRun(m3); err != nil {
		t.Fatalf("unsorted 3-1-2 (1,2,3) should be valid: %v", err)
	}
}
