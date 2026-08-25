package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateSetJokerValid3(t *testing.T) {
	// 7 red, 7 yellow + joker as 7 blue
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-blue-7", tile.Blue, 7)
	m, err := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 7),
		tile.MustTile("t2", tile.Yellow, 7),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	if err != nil {
		t.Fatalf("New joker 3-set: %v", err)
	}
	if err := ValidateSet(m); err != nil {
		t.Fatalf("valid joker 3-set should pass: %v", err)
	}
}

func TestValidateSetJokerValid4(t *testing.T) {
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-black-10", tile.Black, 10)
	m, err := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 10),
		tile.MustTile("t2", tile.Yellow, 10),
		tile.MustTile("t3", tile.Blue, 10),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	if err != nil {
		t.Fatalf("New joker 4-set: %v", err)
	}
	if err := ValidateSet(m); err != nil {
		t.Fatalf("valid joker 4-set should pass: %v", err)
	}
}

func TestValidateSetJokerDuplicateColour(t *testing.T) {
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-red-5", tile.Red, 5) // duplicate red
	m, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	if err := ValidateSet(m); err == nil {
		t.Fatalf("joker duplicate colour should fail")
	}
}

func TestValidateSetJokerRankMismatch(t *testing.T) {
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-blue-6", tile.Blue, 6) // rank 6 vs set 5
	m, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	if err := ValidateSet(m); err == nil {
		t.Fatalf("joker rank mismatch should fail")
	}
}

func TestValidateSetJokerRatio(t *testing.T) {
	// 3 tiles with 2 jokers should fail ratio real>=2*joker (1>=4 false)
	j1 := tile.MustJoker("j1")
	j2 := tile.MustJoker("j2")
	rep1 := tile.MustTile("r1", tile.Blue, 5)
	rep2 := tile.MustTile("r2", tile.Black, 5)
	m, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		j1,
		j2,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep1, "j2": rep2})
	if err := ValidateSet(m); err == nil {
		t.Fatalf("2 jokers in 3-set should fail ratio")
	}
	// 4 tiles with 2 jokers also fails (2>=4 false)
	m2, _ := New("m2", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		j1,
		j2,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep1, "j2": rep2})
	if err := ValidateSet(m2); err == nil {
		t.Fatalf("2 jokers in 4-set should fail ratio")
	}
}
