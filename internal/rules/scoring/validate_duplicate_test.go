package scoring

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// Day 64 — Duplicate tile prevention: any TileInstanceId used more than once across the batch must be rejected
func TestValidateBatchDuplicateAcrossMelds(t *testing.T) {
	dupTile := tile.MustTile("dup-1", tile.Red, 5)
	// Two melds sharing the same tile ID dup-1
	run, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{
		dupTile,
		tile.MustTile("t2", tile.Red, 6),
		tile.MustTile("t3", tile.Red, 7),
	}, nil)
	set, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{
		dupTile,
		tile.MustTile("t4", tile.Yellow, 5),
		tile.MustTile("t5", tile.Blue, 5),
	}, nil)
	rack := []tile.TileInstance{dupTile, tile.MustTile("t2", tile.Red, 6), tile.MustTile("t3", tile.Red, 7), tile.MustTile("t4", tile.Yellow, 5), tile.MustTile("t5", tile.Blue, 5)}
	batch := Batch{PlayerID: "alice", Melds: []meld.Meld{run, set}}
	if err := ValidateBatchOwnership(batch, rack); err == nil {
		t.Fatalf("duplicate tile across melds should fail")
	}
	// Same tile duplicated across two runs
	run2, _ := meld.New("m3", meld.KindRun, []tile.TileInstance{
		tile.MustTile("a1", tile.Blue, 1),
		tile.MustTile("a2", tile.Blue, 2),
		tile.MustTile("a3", tile.Blue, 3),
	}, nil)
	dupRun, _ := meld.New("m4", meld.KindRun, []tile.TileInstance{
		tile.MustTile("a1", tile.Blue, 1), // same ID as run2's first tile
		tile.MustTile("b2", tile.Blue, 5),
		tile.MustTile("b3", tile.Blue, 6),
	}, nil)
	rack2 := []tile.TileInstance{
		tile.MustTile("a1", tile.Blue, 1),
		tile.MustTile("a2", tile.Blue, 2),
		tile.MustTile("a3", tile.Blue, 3),
		tile.MustTile("b2", tile.Blue, 5),
		tile.MustTile("b3", tile.Blue, 6),
	}
	batch2 := Batch{PlayerID: "alice", Melds: []meld.Meld{run2, dupRun}}
	if err := ValidateBatchOwnership(batch2, rack2); err == nil {
		t.Fatalf("duplicate tile across runs should fail")
	}
	// Valid batch with no duplicate should pass
	setValid, _ := meld.New("m5", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t6", tile.Red, 9),
		tile.MustTile("t7", tile.Yellow, 9),
		tile.MustTile("t8", tile.Blue, 9),
	}, nil)
	rackValid := []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 6), tile.MustTile("t3", tile.Red, 7),
		tile.MustTile("t6", tile.Red, 9), tile.MustTile("t7", tile.Yellow, 9), tile.MustTile("t8", tile.Blue, 9),
	}
	batchValid := Batch{PlayerID: "alice", Melds: []meld.Meld{run, setValid}}
	// Need to use run that doesn't share tiles with setValid: run uses t1,t2,t3? Actually run uses dupTile etc. Let's just use run2 and setValid which are distinct
	batchValid = Batch{PlayerID: "alice", Melds: []meld.Meld{run2, setValid}}
	rackValid = []tile.TileInstance{
		tile.MustTile("a1", tile.Blue, 1), tile.MustTile("a2", tile.Blue, 2), tile.MustTile("a3", tile.Blue, 3),
		tile.MustTile("t6", tile.Red, 9), tile.MustTile("t7", tile.Yellow, 9), tile.MustTile("t8", tile.Blue, 9),
	}
	if err := ValidateBatchOwnership(batchValid, rackValid); err != nil {
		t.Fatalf("valid no-duplicate batch should pass, got %v", err)
	}
}
