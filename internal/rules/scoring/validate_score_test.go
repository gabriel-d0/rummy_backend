package scoring

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateBatchScore50(t *testing.T) {
	run30, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 10),
		tile.MustTile("t2", tile.Red, 11),
		tile.MustTile("t3", tile.Red, 12),
	}, nil)
	set15, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t4", tile.Red, 7),
		tile.MustTile("t5", tile.Yellow, 7),
		tile.MustTile("t6", tile.Blue, 7),
	}, nil)
	batch45 := Batch{PlayerID: "alice", Melds: []meld.Meld{run30, set15}}
	if err := ValidateBatchScore(batch45); err == nil {
		t.Fatalf("45 should fail")
	}
	runA, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
	}, nil)
	runB, _ := meld.New("m2", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t4", tile.Blue, 8),
		tile.MustTile("t5", tile.Blue, 9),
		tile.MustTile("t6", tile.Blue, 10),
	}, nil)
	runC, _ := meld.New("m3", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t7", tile.Yellow, 2),
		tile.MustTile("t8", tile.Yellow, 3),
		tile.MustTile("t9", tile.Yellow, 4),
	}, nil)
	batch50 := Batch{PlayerID: "alice", Melds: []meld.Meld{runA, runB, runC}}
	if err := ValidateBatchScore(batch50); err != nil {
		t.Fatalf("50 should pass, got %v", err)
	}
	if v, _ := TotalScore(batch50); v != 50 {
		t.Fatalf("TotalScore 50 got %d", v)
	}
}
