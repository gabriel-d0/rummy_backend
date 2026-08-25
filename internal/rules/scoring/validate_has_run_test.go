package scoring

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateBatchHasRun(t *testing.T) {
	// No run — only sets should fail
	set, _ := meld.New("m1", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		tile.MustTile("t3", tile.Blue, 5),
	}, nil)
	batchNoRun := Batch{PlayerID: "alice", Melds: []meld.Meld{set, set}}
	if err := ValidateBatchHasRun(batchNoRun); err == nil {
		t.Fatalf("no run should fail")
	}
	// With one run should pass
	run, _ := meld.New("m2", meld.KindRun, []tile.TileInstance{
		tile.MustTile("t4", tile.Red, 1),
		tile.MustTile("t5", tile.Red, 2),
		tile.MustTile("t6", tile.Red, 3),
	}, nil)
	batchWithRun := Batch{PlayerID: "alice", Melds: []meld.Meld{set, run}}
	if err := ValidateBatchHasRun(batchWithRun); err != nil {
		t.Fatalf("with run should pass, got %v", err)
	}
	// Empty batch should fail
	if err := ValidateBatchHasRun(Batch{PlayerID: "alice", Melds: []meld.Meld{}}); err == nil {
		t.Fatalf("empty batch should fail")
	}
}
