package meld

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateRunHighAce(t *testing.T) {
	// 12-13-1 valid
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 12),
		tile.MustTile("t2", tile.Red, 13),
		tile.MustTile("t3", tile.Red, 1),
	}, nil)
	if err := ValidateRun(m); err != nil {
		t.Fatalf("12-13-1 should be valid: %v", err)
	}
	// 11-12-13-1 valid
	m2, _ := New("m2", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Blue, 11),
		tile.MustTile("t2", tile.Blue, 12),
		tile.MustTile("t3", tile.Blue, 13),
		tile.MustTile("t4", tile.Blue, 1),
	}, nil)
	if err := ValidateRun(m2); err != nil {
		t.Fatalf("11-12-13-1 should be valid: %v", err)
	}
	// 10-11-12-13-1 valid
	m3, _ := New("m3", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Yellow, 10),
		tile.MustTile("t2", tile.Yellow, 11),
		tile.MustTile("t3", tile.Yellow, 12),
		tile.MustTile("t4", tile.Yellow, 13),
		tile.MustTile("t5", tile.Yellow, 1),
	}, nil)
	if err := ValidateRun(m3); err != nil {
		t.Fatalf("10-11-12-13-1 should be valid: %v", err)
	}
}

func TestValidateRunHighAceInvalid(t *testing.T) {
	// 13-1-2 invalid per docs
	m, _ := New("m1", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 13),
		tile.MustTile("t2", tile.Red, 1),
		tile.MustTile("t3", tile.Red, 2),
	}, nil)
	if err := ValidateRun(m); err == nil {
		t.Fatalf("13-1-2 should be invalid")
	}
	// 13-1 alone not enough (len 2) already fails, but test 13-1-3 (gap)
	m2, _ := New("m2", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Blue, 13),
		tile.MustTile("t2", tile.Blue, 1),
		tile.MustTile("t3", tile.Blue, 3),
	}, nil)
	if err := ValidateRun(m2); err == nil {
		t.Fatalf("13-1-3 should be invalid")
	}
	// 12-13-1-2 invalid (contains both low and high Ace? would be 1,2,12,13 not consecutive nor high)
	m3, _ := New("m3", KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 12),
		tile.MustTile("t2", tile.Red, 13),
		tile.MustTile("t3", tile.Red, 1),
		tile.MustTile("t4", tile.Red, 2),
	}, nil)
	if err := ValidateRun(m3); err == nil {
		t.Fatalf("12-13-1-2 should be invalid (Ace in middle)")
	}
}
