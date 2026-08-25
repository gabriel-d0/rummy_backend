package scoring

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateBatchOwnership(t *testing.T) {
	rack := []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
		tile.MustTile("t4", tile.Yellow, 7),
		tile.MustTile("t5", tile.Blue, 7),
	}
	run, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{rack[0], rack[1], rack[2]}, nil)
	set, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{rack[2], rack[3], rack[4]}, nil) // shares t3 duplicate
	batchDup := Batch{PlayerID: "alice", Melds: []meld.Meld{run, set}}
	if err := ValidateBatchOwnership(batchDup, rack); err == nil {
		t.Fatalf("duplicate tile across melds should fail")
	}
	// Valid: run 5-6-7 and set 7 with distinct tiles (need 6 tiles, but rack has 5, so we need 5 distinct)
	// Use run 5-6-7 (t1,t2,t3) and set with 7 red-> need another 7? Actually set needs 3 tiles of same rank 7, we have t3 is 7 red, t4 7 yellow, t5 7 blue — but t3 is already in run.
	// For valid batch, use run 5-6-7 (t1,t2) plus need 3 tiles for run, but we have t1,t2 and need a third not used in set.
	// Let's create valid batch with no overlap: run 1-2-3 and set 5
	runValid, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 1),
		tile.MustTile("t2", tile.Red, 2),
		tile.MustTile("t3", tile.Red, 3),
	}, nil)
	setValid, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t4", tile.Red, 5),
		tile.MustTile("t5", tile.Yellow, 5),
		tile.MustTile("t6", tile.Blue, 5),
	}, nil)
	rack2 := []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 1),
		tile.MustTile("t2", tile.Red, 2),
		tile.MustTile("t3", tile.Red, 3),
		tile.MustTile("t4", tile.Red, 5),
		tile.MustTile("t5", tile.Yellow, 5),
		tile.MustTile("t6", tile.Blue, 5),
	}
	batchValid := Batch{PlayerID: "alice", Melds: []meld.Meld{runValid, setValid}}
	if err := ValidateBatchOwnership(batchValid, rack2); err != nil {
		t.Fatalf("valid batch ownership should pass: %v", err)
	}
	// Tile not in rack
	mForeign, _ := meld.New("m3", meld.KindRun, []tile.TileInstance{
		tile.MustTile("foreign", tile.Red, 9),
		tile.MustTile("t2", tile.Red, 2),
		tile.MustTile("t3", tile.Red, 3),
	}, nil)
	batchForeign := Batch{PlayerID: "alice", Melds: []meld.Meld{mForeign}}
	if err := ValidateBatchOwnership(batchForeign, rack2); err == nil {
		t.Fatalf("foreign tile should fail")
	}
}
