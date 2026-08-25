package scoring

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestScoreJokerExplicit(t *testing.T) {
	joker := tile.MustJoker("j1")
	// Joker as 5 =>5
	rep5 := tile.MustTile("rep5", tile.Red, 5)
	if v := ScoreTile(joker, false, false, false, &rep5); v != 5 {
		t.Fatalf("joker 5 =>5 got %d", v)
	}
	// Joker as 10 =>10
	rep10 := tile.MustTile("rep10", tile.Blue, 10)
	if v := ScoreTile(joker, false, false, false, &rep10); v != 10 {
		t.Fatalf("joker 10 =>10 got %d", v)
	}
	// Joker as low Ace 1-2-3 =>5
	repLow := tile.MustTile("repLow", tile.Yellow, 1)
	if v := ScoreTile(joker, true, false, false, &repLow); v != 5 {
		t.Fatalf("joker low Ace 5 got %d", v)
	}
	// Joker as high Ace 12-13-1 =>10
	repHigh := tile.MustTile("repHigh", tile.Black, 1)
	if v := ScoreTile(joker, false, true, false, &repHigh); v != 10 {
		t.Fatalf("joker high Ace 10 got %d", v)
	}
	// Joker as Ace in Ace set =>25
	repAceSet := tile.MustTile("repAceSet", tile.Red, 1)
	if v := ScoreTile(joker, false, false, true, &repAceSet); v != 25 {
		t.Fatalf("joker Ace set 25 got %d", v)
	}
}
