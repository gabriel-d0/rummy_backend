package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestMeldValidateNoJoker(t *testing.T) {
	m, err := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 1),
		tile.MustTile("t2", tile.Red, 2),
		tile.MustTile("t3", tile.Red, 3),
	}, nil)
	if err != nil {
		t.Fatalf("New run no joker: %v", err)
	}
	if m.Kind != KindRun || len(m.Tiles) != 3 {
		t.Fatalf("meld fields")
	}
	// Ensure Tiles copy is immutable
	m.Tiles[0] = tile.MustTile("hack", tile.Blue, 5)
	orig, _ := New("m2", KindRun, []tile.TileInstance{tile.MustTile("t1", tile.Red, 1), tile.MustTile("t2", tile.Red, 2), tile.MustTile("t3", tile.Red, 3)}, nil)
	if orig.Tiles[0].ID == "hack" {
		t.Fatalf("Tiles not copied")
	}
}

func TestMeldWithJokerExplicitRep(t *testing.T) {
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-1", tile.Red, 3)
	m, err := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 1),
		tile.MustTile("t2", tile.Red, 2),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": rep})
	if err != nil {
		t.Fatalf("joker run: %v", err)
	}
	if m.JokerReps["j1"].Rank != 3 {
		t.Fatalf("rep rank")
	}
	// Missing rep should fail
	if _, err := New("m2", KindRun, []tile.TileInstance{tile.MustTile("t1", tile.Red, 1), tile.MustTile("t2", tile.Red, 2), j}, nil); err == nil {
		t.Fatalf("missing rep should fail")
	}
	// Rep for non-joker should fail
	if _, err := New("m3", KindRun, []tile.TileInstance{tile.MustTile("t1", tile.Red, 1), tile.MustTile("t2", tile.Red, 2), tile.MustTile("t3", tile.Red, 3)}, map[tile.TileInstanceId]tile.TileInstance{"t1": rep}); err == nil {
		t.Fatalf("rep for non-joker should fail")
	}
}

func TestMeldIDAndKindValidation(t *testing.T) {
	if _, err := New("", KindRun, []tile.TileInstance{tile.MustTile("t1", tile.Red, 1), tile.MustTile("t2", tile.Red, 2), tile.MustTile("t3", tile.Red, 3)}, nil); err == nil {
		t.Fatalf("empty ID should fail")
	}
	if _, err := New("m1", Kind("bad"), []tile.TileInstance{tile.MustTile("t1", tile.Red, 1), tile.MustTile("t2", tile.Red, 2), tile.MustTile("t3", tile.Red, 3)}, nil); err == nil {
		t.Fatalf("bad kind should fail")
	}
	if _, err := New("m1", KindSet, []tile.TileInstance{tile.MustTile("t1", tile.Red, 1), tile.MustTile("t2", tile.Red, 2)}, nil); err == nil {
		t.Fatalf("<3 tiles should fail")
	}
}

func TestMeldDuplicateTile(t *testing.T) {
	dup := tile.MustTile("dup", tile.Red, 5)
	if _, err := New("m1", KindSet, []tile.TileInstance{dup, dup, tile.MustTile("t3", tile.Blue, 5)}, nil); err == nil {
		t.Fatalf("duplicate tile should fail")
	}
}
