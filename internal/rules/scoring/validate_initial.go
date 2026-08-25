// Package scoring — Initial meld validation (Day 65).
package scoring

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// ValidateInitialBatch checks that a batch is a legal opening meld:
// - all tiles owned and no duplicate (ValidateBatchOwnership)
// - total >=50 (ValidateBatchScore)
// - at least one run (ValidateBatchHasRun)
// It does not check HasOpened already (caller must ensure player not yet opened).
func ValidateInitialBatch(batch Batch, rack []tile.TileInstance) error {
	if err := ValidateBatchOwnership(batch, rack); err != nil {
		return err
	}
	if err := ValidateBatchScore(batch); err != nil {
		return fmt.Errorf("initial score %w", err)
	}
	if err := ValidateBatchHasRun(batch); err != nil {
		return err
	}
	return nil
}
