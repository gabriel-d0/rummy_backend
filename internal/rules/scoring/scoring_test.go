package scoring

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestScoreTileSimple(t *testing.T) {
	if v := ScoreTileSimple(tile.MustTile("t2", tile.Red, 2)); v != 5 {
		t.Fatalf("2 =>5 got %d", v)
	}
	if v := ScoreTileSimple(tile.MustTile("t9", tile.Blue, 9)); v != 5 {
		t.Fatalf("9 =>5 got %d", v)
	}
	if v := ScoreTileSimple(tile.MustTile("t10", tile.Yellow, 10)); v != 10 {
		t.Fatalf("10 =>10 got %d", v)
	}
	if v := ScoreTileSimple(tile.MustTile("t13", tile.Black, 13)); v != 10 {
		t.Fatalf("13 =>10 got %d", v)
	}
}

func TestScoreAceLowAndHigh(t *testing.T) {
	ace := tile.MustTile("a1", tile.Red, 1)
	if v := ScoreTile(ace, true, false, false, nil); v != 5 {
		t.Fatalf("Ace low 1-2-3 =>5 got %d", v)
	}
	if v := ScoreTile(ace, false, true, false, nil); v != 10 {
		t.Fatalf("Ace high 12-13-1 =>10 got %d", v)
	}
	// Bare Ace defaults to 5
	if v := ScoreTile(ace, false, false, false, nil); v != 5 {
		t.Fatalf("bare Ace =>5 got %d", v)
	}
}

func TestScoreAceSet(t *testing.T) {
	ace := tile.MustTile("a1", tile.Red, 1)
	if v := ScoreTile(ace, false, false, true, nil); v != 25 {
		t.Fatalf("Ace set 25 got %d", v)
	}
	// Non-Ace in AceSet context still scores normally (should not happen, but check)
	five := tile.MustTile("t5", tile.Red, 5)
	if v := ScoreTile(five, false, false, true, nil); v != 5 {
		t.Fatalf("5 in aceSet context should still be 5, got %d", v)
	}
}

func TestScoreJoker(t *testing.T) {
	joker := tile.MustJoker("j1")
	repLow := tile.MustTile("rep-low", tile.Red, 1)
	repHigh := tile.MustTile("rep-high", tile.Red, 1)
	repFive := tile.MustTile("rep-5", tile.Blue, 5)
	repTen := tile.MustTile("rep-10", tile.Yellow, 10)

	if v := ScoreTile(joker, true, false, false, &repLow); v != 5 {
		t.Fatalf("joker low Ace 5 got %d", v)
	}
	if v := ScoreTile(joker, false, true, false, &repHigh); v != 10 {
		t.Fatalf("joker high Ace 10 got %d", v)
	}
	if v := ScoreTile(joker, false, false, false, &repFive); v != 5 {
		t.Fatalf("joker 5 =>5 got %d", v)
	}
	if v := ScoreTile(joker, false, false, false, &repTen); v != 10 {
		t.Fatalf("joker 10 =>10 got %d", v)
	}
	if v := ScoreTile(joker, false, false, true, &repLow); v != 25 {
		t.Fatalf("joker Ace set 25 got %d", v)
	}
	// Joker without rep =>0
	if v := ScoreTile(joker, false, false, false, nil); v != 0 {
		t.Fatalf("joker nil rep =>0 got %d", v)
	}
}
