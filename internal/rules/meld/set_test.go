package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateSetValid3and4(t *testing.T) {
	// 3-colour set: 7 red, 7 yellow, 7 blue
	m3, err := New("m3", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 7),
		tile.MustTile("t2", tile.Yellow, 7),
		tile.MustTile("t3", tile.Blue, 7),
	}, nil)
	if err != nil {
		t.Fatalf("New m3: %v", err)
	}
	if err := ValidateSet(m3); err != nil {
		t.Fatalf("valid 3-set should pass: %v", err)
	}
	// 4-colour set: 10 red, yellow, blue, black
	m4, err := New("m4", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 10),
		tile.MustTile("t2", tile.Yellow, 10),
		tile.MustTile("t3", tile.Blue, 10),
		tile.MustTile("t4", tile.Black, 10),
	}, nil)
	if err != nil {
		t.Fatalf("New m4: %v", err)
	}
	if err := ValidateSet(m4); err != nil {
		t.Fatalf("valid 4-set should pass: %v", err)
	}
}

func TestValidateSetRejectsDuplicateColour(t *testing.T) {
	m, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 5),
		tile.MustTile("t3", tile.Blue, 5),
	}, nil)
	if err := ValidateSet(m); err == nil {
		t.Fatalf("duplicate colour should fail")
	}
}

func TestValidateSetRejectsRankMismatch(t *testing.T) {
	m, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		tile.MustTile("t3", tile.Blue, 6),
	}, nil)
	if err := ValidateSet(m); err == nil {
		t.Fatalf("rank mismatch should fail")
	}
}

func TestValidateSetRejectsInvalidSize(t *testing.T) {
	// 2 tiles
	m2, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		tile.MustTile("t3", tile.Blue, 5),
	}, nil)
	// Manually truncate to 2 for test of ValidateSet size check (New would reject <3, so we bypass New)
	m2.Tiles = m2.Tiles[:2]
	if err := ValidateSet(m2); err == nil {
		t.Fatalf("2 tiles should fail")
	}
	// 5 tiles
	m5, _ := New("m5", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		tile.MustTile("t3", tile.Blue, 5),
	}, nil)
	m5.Tiles = append(m5.Tiles, tile.MustTile("t4", tile.Black, 5), tile.MustTile("t5", tile.Red, 5))
	if err := ValidateSet(m5); err == nil {
		t.Fatalf("5 tiles should fail")
	}
}

func TestValidateSetRejectsJoker(t *testing.T) {
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep", tile.Black, 5)
	m, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	if err := ValidateSet(m); err == nil {
		t.Fatalf("joker should be rejected by basic ValidateSet")
	}
}

func TestValidateSetKindMismatch(t *testing.T) {
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 1),
		tile.MustTile("t2", tile.Red, 2),
		tile.MustTile("t3", tile.Red, 3),
	}, nil)
	if err := ValidateSet(m); err == nil {
		t.Fatalf("KindRun should fail ValidateSet")
	}
}
