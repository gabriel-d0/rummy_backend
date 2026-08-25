// Package scoring — Opening meld validation (Day 61).
package scoring

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// ValidateBatchOwnership checks that every tile in batch.Melds is owned
// in the player's rack and that no TileInstanceId is used more than once
// across melds. It does not yet check 50-point or run requirement (Day 62-63).
func ValidateBatchOwnership(batch Batch, rack []tile.TileInstance) error {
	rackByID := make(map[tile.TileInstanceId]bool, len(rack))
	for _, t := range rack {
		rackByID[t.ID] = true
	}
	seen := map[tile.TileInstanceId]bool{}
	for _, m := range batch.Melds {
		// Meld structural validation
		if err := m.Validate(); err != nil {
			return fmt.Errorf("meld %q invalid: %w", m.ID, err)
		}
		// ValidateSet/Run also ensures joker reps are valid, but we call them for extra
		switch m.Kind {
		case meld.KindRun:
			if err := meld.ValidateRun(m); err != nil {
				return fmt.Errorf("run %q invalid: %w", m.ID, err)
			}
		case meld.KindSet:
			if err := meld.ValidateSet(m); err != nil {
				return fmt.Errorf("set %q invalid: %w", m.ID, err)
			}
		default:
			return fmt.Errorf("meld %q invalid kind %q", m.ID, m.Kind)
		}
		for _, t := range m.Tiles {
			if seen[t.ID] {
				return fmt.Errorf("duplicate tile %v across melds", t.ID)
			}
			seen[t.ID] = true
			if !rackByID[t.ID] {
				return fmt.Errorf("tile %v not in rack", t.ID)
			}
		}
	}
	return nil
}
