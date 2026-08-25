package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateRunBasic(t *testing.T) {
	// Valid run 5-6-7 red
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
	}, nil)
	if err := ValidateRun(m); err != nil {
		t.Fatalf("valid run 5-6-7: %v", err)
	}
	// Unsorted input should still pass (sorted inside)
	m2, _ := New("m2", KindRun, []tile.TileInstance{
		tile.MustTile("t3", tile.Red, 7),
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
	}, nil)
	if err := ValidateRun(m2); err != nil {
		t.Fatalf("unsorted valid run: %v", err)
	}
	// Longer run 1-2-3-4
	m3, _ := New("m3", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Blue, 1),
		tile.MustTile("t2", tile.Blue, 2),
		tile.MustTile("t3", tile.Blue, 3),
		tile.MustTile("t4", tile.Blue, 4),
	}, nil)
	if err := ValidateRun(m3); err != nil {
		t.Fatalf("1-2-3-4 should pass (low Ace), got %v", err)
	}
}

func TestValidateRunRejectsSameColourAndLength(t *testing.T) {
	// Different colours
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Blue, 6),
		tile.MustTile("t3", tile.Red, 7),
	}, nil)
	if err := ValidateRun(m); err == nil {
		t.Fatalf("different colours should fail")
	}
	// Too short <3
	m2, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
	}, nil)
	m2.Tiles = m2.Tiles[:2]
	if err := ValidateRun(m2); err == nil {
		t.Fatalf("<3 should fail")
	}
	// Gap not consecutive 5-7-8
	m3, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 7),
		tile.MustTile("t3", tile.Red, 8),
	}, nil)
	if err := ValidateRun(m3); err == nil {
		t.Fatalf("gap 5-7-8 should fail")
	}
	// Duplicate rank 5-6-6
	m4, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 6),
	}, nil)
	if err := ValidateRun(m4); err == nil {
		t.Fatalf("duplicate rank should fail")
	}
}

func TestValidateRunRejectsJoker(t *testing.T) {
	// Day 52: joker with valid rep is now allowed
	j := tile.MustJoker("j1")
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("rep", tile.Red, 7)})
	if err := ValidateRun(m); err != nil {
		t.Fatalf("joker with valid rep should now pass (Day 52), got %v", err)
	}
	// Missing rep should still fail
	m2 := Meld{ID: "m2", Kind: KindRun, Tiles: []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 7), j}, JokerReps: map[tile.TileInstanceId]tile.TileInstance{}}
	if err := ValidateRun(m2); err == nil {
		t.Fatalf("joker missing rep should fail")
	}
}

func TestValidateRunKindMismatch(t *testing.T) {
	m, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		tile.MustTile("t3", tile.Blue, 5),
	}, nil)
	if err := ValidateRun(m); err == nil {
		t.Fatalf("KindSet should fail ValidateRun")
	}
}
