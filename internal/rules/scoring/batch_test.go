package scoring

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestTotalScoreBatch(t *testing.T) {
	// Run 5-6-7 red =>15, Set 7 red/yellow/blue =>15, total 30
	run, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
	}, nil)
	set, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t4", tile.Red, 7),
		tile.MustTile("t5", tile.Yellow, 7),
		tile.MustTile("t6", tile.Blue, 7),
	}, nil)
	batch := Batch{PlayerID: "alice", Melds: []meld.Meld{run, set}}
	if v, err := TotalScore(batch); err != nil || v != 30 {
		t.Fatalf("TotalScore 30 got %d err %v", v, err)
	}
	// Ace set 25*3=75 + low run 1-2-3 15 =>90
	aceSet, _ := meld.New("m3", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t7", tile.Red, 1),
		tile.MustTile("t8", tile.Yellow, 1),
		tile.MustTile("t9", tile.Blue, 1),
	}, nil)
	runLow, _ := meld.New("m4", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t10", tile.Red, 1),
		tile.MustTile("t11", tile.Red, 2),
		tile.MustTile("t12", tile.Red, 3),
	}, nil)
	batch2 := Batch{PlayerID: "bob", Melds: []meld.Meld{aceSet, runLow}}
	if v, _ := TotalScore(batch2); v != 90 {
		t.Fatalf("Ace set + low run 90 got %d", v)
	}
}
