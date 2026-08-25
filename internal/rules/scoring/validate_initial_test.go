package scoring

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateInitialBatch(t *testing.T) {
	// Valid: run 5-6-7 (15) + run 8-9-10 (20) + run 2-3-4 (15) =50 with at least one run
	runA, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 6), tile.MustTile("t3", tile.Red, 7)}, nil)
	runB, _ := meld.New("m2", meld.KindRun, []tile.TileInstance{tile.MustTile("t4", tile.Blue, 8), tile.MustTile("t5", tile.Blue, 9), tile.MustTile("t6", tile.Blue, 10)}, nil)
	runC, _ := meld.New("m3", meld.KindRun, []tile.TileInstance{tile.MustTile("t7", tile.Yellow, 2), tile.MustTile("t8", tile.Yellow, 3), tile.MustTile("t9", tile.Yellow, 4)}, nil)
	batchValid := Batch{PlayerID: "alice", Melds: []meld.Meld{runA, runB, runC}}
	rackValid := []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Red, 6), tile.MustTile("t3", tile.Red, 7),
		tile.MustTile("t4", tile.Blue, 8), tile.MustTile("t5", tile.Blue, 9), tile.MustTile("t6", tile.Blue, 10),
		tile.MustTile("t7", tile.Yellow, 2), tile.MustTile("t8", tile.Yellow, 3), tile.MustTile("t9", tile.Yellow, 4),
	}
	if err := ValidateInitialBatch(batchValid, rackValid); err != nil {
		t.Fatalf("valid initial batch should pass: %v", err)
	}
	// 45 should fail (not enough points)
	run30, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{tile.MustTile("t1", tile.Red, 10), tile.MustTile("t2", tile.Red, 11), tile.MustTile("t3", tile.Red, 12)}, nil)  // 30
	set15, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{tile.MustTile("t4", tile.Red, 7), tile.MustTile("t5", tile.Yellow, 7), tile.MustTile("t6", tile.Blue, 7)}, nil) // 15 total 45
	batch45 := Batch{PlayerID: "alice", Melds: []meld.Meld{run30, set15}}
	rack45 := []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 10), tile.MustTile("t2", tile.Red, 11), tile.MustTile("t3", tile.Red, 12),
		tile.MustTile("t4", tile.Red, 7), tile.MustTile("t5", tile.Yellow, 7), tile.MustTile("t6", tile.Blue, 7),
	}
	if err := ValidateInitialBatch(batch45, rack45); err == nil {
		t.Fatalf("45 should fail")
	}
	// No run should fail
	set5, _ := meld.New("m1", meld.KindSet, []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Yellow, 5), tile.MustTile("t3", tile.Blue, 5)}, nil)
	set6, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{tile.MustTile("t4", tile.Red, 6), tile.MustTile("t5", tile.Yellow, 6), tile.MustTile("t6", tile.Blue, 6)}, nil)
	set7, _ := meld.New("m3", meld.KindSet, []tile.TileInstance{tile.MustTile("t7", tile.Red, 7), tile.MustTile("t8", tile.Yellow, 7), tile.MustTile("t9", tile.Blue, 7)}, nil)    // 5+5+5=15 each, total 45 but also no run, but need 50+ to make total pass, so add another set to make 60
	set8, _ := meld.New("m4", meld.KindSet, []tile.TileInstance{tile.MustTile("t10", tile.Red, 8), tile.MustTile("t11", tile.Yellow, 8), tile.MustTile("t12", tile.Blue, 8)}, nil) // 60 total, but no run
	batchNoRun := Batch{PlayerID: "alice", Melds: []meld.Meld{set5, set6, set7, set8}}                                                                                             // 15*4=60
	rackNoRun := []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Yellow, 5), tile.MustTile("t3", tile.Blue, 5),
		tile.MustTile("t4", tile.Red, 6), tile.MustTile("t5", tile.Yellow, 6), tile.MustTile("t6", tile.Blue, 6),
		tile.MustTile("t7", tile.Red, 7), tile.MustTile("t8", tile.Yellow, 7), tile.MustTile("t9", tile.Blue, 7),
		tile.MustTile("t10", tile.Red, 8), tile.MustTile("t11", tile.Yellow, 8), tile.MustTile("t12", tile.Blue, 8),
	}
	if err := ValidateInitialBatch(batchNoRun, rackNoRun); err == nil {
		t.Fatalf("no run should fail")
	}
	// Duplicate across melds should fail
	dupTile := tile.MustTile("dup", tile.Red, 5)
	runDup, _ := meld.New("m1", meld.KindRun, []tile.TileInstance{dupTile, tile.MustTile("t2", tile.Red, 6), tile.MustTile("t3", tile.Red, 7)}, nil)
	setDup, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{dupTile, tile.MustTile("t5", tile.Yellow, 5), tile.MustTile("t6", tile.Blue, 5)}, nil)
	batchDup := Batch{PlayerID: "alice", Melds: []meld.Meld{runDup, setDup}}
	rackDup := []tile.TileInstance{dupTile, tile.MustTile("t2", tile.Red, 6), tile.MustTile("t3", tile.Red, 7), tile.MustTile("t5", tile.Yellow, 5), tile.MustTile("t6", tile.Blue, 5)}
	if err := ValidateInitialBatch(batchDup, rackDup); err == nil {
		t.Fatalf("duplicate across melds should fail")
	}
}
