package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// Day 53 — Run joker ratio: real >= 2*joker per docs/rules-decisions.md:1.3
func TestValidateRunJokerRatioExplicit(t *testing.T) {
	// 3 tiles with 1 joker: 2 real vs 1 joker => 2>=2 true, should pass
	j := tile.MustJoker("j1")
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("r1", tile.Red, 7)})
	if err := ValidateRun(m); err != nil {
		t.Fatalf("2 real +1 joker in 3 should pass (2>=2): %v", err)
	}
	// 4 tiles with 1 joker: 3 real vs 1 => 3>=2 true
	m2, _ := New("m2", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("r2", tile.Red, 8)})
	if err := ValidateRun(m2); err != nil {
		t.Fatalf("3 real +1 joker in 4 should pass: %v", err)
	}
	// 5 tiles with 2 jokers: 3 real vs 2 => 3>=4 false
	j2 := tile.MustJoker("j2")
	m3, _ := New("m3", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
		j, j2,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("r1", tile.Red, 8), "j2": tile.MustTile("r2", tile.Red, 9)})
	if err := ValidateRun(m3); err == nil {
		t.Fatalf("3 real +2 jokers in 5 should fail 3>=4 false")
	}
	// 6 tiles with 2 jokers: 4 real vs 2 => 4>=4 true
	j3 := tile.MustJoker("j3")
	j4 := tile.MustJoker("j4")
	m4, _ := New("m4", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
		tile.MustTile("t4", tile.Red, 8),
		j3, j4,
	}, map[tile.TileInstanceId]tile.TileInstance{"j3": tile.MustTile("r3", tile.Red, 9), "j4": tile.MustTile("r4", tile.Red, 10)})
	if err := ValidateRun(m4); err != nil {
		t.Fatalf("4 real +2 jokers in 6 should pass 4>=4: %v", err)
	}
	// 7 tiles with 3 jokers: 4 real vs 3 => 4>=6 false
	j5 := tile.MustJoker("j5")
	m5, _ := New("m5", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
		tile.MustTile("t4", tile.Red, 8),
		j, j2, j5,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("r1", tile.Red, 9), "j2": tile.MustTile("r2", tile.Red, 10), "j5": tile.MustTile("r5", tile.Red, 11)})
	if err := ValidateRun(m5); err == nil {
		t.Fatalf("4 real +3 jokers in 7 should fail 4>=6 false")
	}
}
