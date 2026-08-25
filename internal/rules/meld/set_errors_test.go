package meld

import (
	"errors"
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateSetStructuredErrors(t *testing.T) {
	// duplicate colour
	mDup, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Red, 5),
		tile.MustTile("t3", tile.Blue, 5),
	}, nil)
	err := ValidateSet(mDup)
	var ve *ValidationError
	if !errors.As(err, &ve) || ve.Code != ErrCodeDuplicateColour {
		t.Fatalf("duplicate colour should be ErrCodeDuplicateColour, got %v", err)
	}
	// rank mismatch
	mRank, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		tile.MustTile("t3", tile.Blue, 6),
	}, nil)
	err = ValidateSet(mRank)
	if !errors.As(err, &ve) || ve.Code != ErrCodeRankMismatch {
		t.Fatalf("rank mismatch should be ErrCodeRankMismatch, got %v", err)
	}
	// invalid size
	mSmall, _ := New("m1", KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		tile.MustTile("t3", tile.Blue, 5),
	}, nil)
	mSmall.Tiles = mSmall.Tiles[:2]
	err = ValidateSet(mSmall)
	if !errors.As(err, &ve) || ve.Code != ErrCodeInvalidSize {
		t.Fatalf("invalid size should be ErrCodeInvalidSize, got %v", err)
	}
	// duplicate tile (caught by Meld.Validate before ValidateSet, but also test via New)
	dup := tile.MustTile("dup", tile.Red, 5)
	if _, err := New("m1", KindSet, []tile.TileInstance{dup, dup, tile.MustTile("t3", tile.Blue, 5)}, nil); err == nil {
		t.Fatalf("duplicate tile should fail at New")
	} else {
		var ve2 *ValidationError
		// New's Validate will return duplicate_tile via Meld.Validate, not Via ValidateSet
		// We just check that New fails (covered)
		_ = ve2
	}
}
