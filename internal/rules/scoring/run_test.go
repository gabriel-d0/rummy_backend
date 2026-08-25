package scoring

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestScoreRunSimple(t *testing.T) {
	// 5-6-7 red => 5+5+5=15
	m, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
	}, nil)
	if v, err := ScoreRun(m); err != nil || v != 15 {
		t.Fatalf("5-6-7 =>15 got %d err %v", v, err)
	}
	// 10-11-12 => 10+10+10=30
	m2, _ := meld.New("m2", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Blue, 10),
		tile.MustTile("t2", tile.Blue, 11),
		tile.MustTile("t3", tile.Blue, 12),
	}, nil)
	if v, _ := ScoreRun(m2); v != 30 {
		t.Fatalf("10-11-12 =>30 got %d", v)
	}
}

func TestScoreRunLowAndHighAce(t *testing.T) {
	// 1-2-3 low => Ace 5 +5+5=15
	mLow, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 1),
		tile.MustTile("t2", tile.Red, 2),
		tile.MustTile("t3", tile.Red, 3),
	}, nil)
	if v, _ := ScoreRun(mLow); v != 15 {
		t.Fatalf("1-2-3 low =>15 got %d", v)
	}
	// 12-13-1 high => 10+10+10=30
	mHigh, _ := meld.New("m2", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Yellow, 12),
		tile.MustTile("t2", tile.Yellow, 13),
		tile.MustTile("t3", tile.Yellow, 1),
	}, nil)
	if v, _ := ScoreRun(mHigh); v != 30 {
		t.Fatalf("12-13-1 high =>30 got %d", v)
	}
}

func TestScoreRunJoker(t *testing.T) {
	// 5-7 + joker 6 => 5+5+5=15
	j := tile.MustJoker("j1")
	m, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 7),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("rep", tile.Red, 6)})
	if v, _ := ScoreRun(m); v != 15 {
		t.Fatalf("joker 6 =>15 got %d", v)
	}
	// Joker as high Ace 1 =>10
	j2 := tile.MustJoker("j1")
	m2, _ := meld.New("m2", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Blue, 12),
		tile.MustTile("t2", tile.Blue, 13),
		j2,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("rep", tile.Blue, 1)})
	if v, _ := ScoreRun(m2); v != 30 {
		t.Fatalf("joker Ace high =>30 got %d", v)
	}
}
