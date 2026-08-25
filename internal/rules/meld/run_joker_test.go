package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateRunJokerGap(t *testing.T) {
	// 5 red, 7 red + joker as 6 red => 5-6-7
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-6", tile.Red, 6)
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 7),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	if err := ValidateRun(m); err != nil {
		t.Fatalf("joker gap 5-6-7 should pass: %v", err)
	}
}

func TestValidateRunJokerAtEnd(t *testing.T) {
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-7", tile.Red, 7)
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	if err := ValidateRun(m); err != nil {
		t.Fatalf("joker at end 5-6-7 should pass: %v", err)
	}
}

func TestValidateRunJokerWrongColour(t *testing.T) {
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-blue-6", tile.Blue, 6) // wrong colour
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 7),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	if err := ValidateRun(m); err == nil {
		t.Fatalf("joker wrong colour should fail")
	}
}

func TestValidateRunJokerRatio(t *testing.T) {
	// 3 tiles with 2 jokers => 1 real vs 2 jokers ratio 1>=4 false
	j1 := tile.MustJoker("j1")
	j2 := tile.MustJoker("j2")
	rep1 := tile.MustTile("r1", tile.Red, 6)
	rep2 := tile.MustTile("r2", tile.Red, 7)
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		j1,
		j2,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep1, "j2": rep2})
	if err := ValidateRun(m); err == nil {
		t.Fatalf("2 jokers in 3 should fail ratio")
	}
	// 5 tiles with 2 jokers => 3 real vs 2 jokers 3>=4 false? Actually 3>=4 false, so 5 with 2 jokers should also fail ratio? Wait ratio real>=2*joker: 3>=4 false, so 5 with 2 jokers (3 real) fails, 6 with 2 jokers (4 real) passes 4>=4 true
	// So 5 tiles with 2 jokers should fail
	m2, _ := New("m2", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
		j1,
		j2,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("r1", tile.Red, 8), "j2": tile.MustTile("r2", tile.Red, 9)})
	if err := ValidateRun(m2); err == nil {
		t.Fatalf("3 real +2 jokers in 5 should fail ratio 3>=4 false")
	}
	// 6 tiles with 2 jokers => 4 real >=4 true, should pass if consecutive 5-10
	j3 := tile.MustJoker("j1")
	j4 := tile.MustJoker("j2")
	// Use non-conflicting IDs
	j3 = tile.MustJoker("j3")
	j4 = tile.MustJoker("j4")
	m3, _ := New("m3", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
		tile.MustTile("t4", tile.Red, 8),
		j3,
		j4,
	}, map[tile.TileInstanceId]tile.TileInstance{"j3": tile.MustTile("r3", tile.Red, 9), "j4": tile.MustTile("r4", tile.Red, 10)})
	if err := ValidateRun(m3); err != nil {
		t.Fatalf("4 real +2 jokers in 6 should pass (4>=4): %v", err)
	}
}

func TestValidateRunJokerMissingRep(t *testing.T) {
	j := tile.MustJoker("j1")
	// No rep provided, New will fail, but ValidateRun should also fail if we bypass New
	m := Meld{ID: "m1", Kind: KindRun, Tiles: []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 7), j}, JokerReps: map[tile.TileInstanceId]tile.TileInstance{}}
	if err := ValidateRun(m); err == nil {
		t.Fatalf("missing rep should fail")
	}
}
