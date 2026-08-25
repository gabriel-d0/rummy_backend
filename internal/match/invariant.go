// Package match — Tile conservation invariant (Day 13).
// Ensures every TileInstanceId appears exactly once across racks, stock,
// discard row, and table melds — per docs/rules-decisions.md:8.
package match

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// CheckTileConservation verifies that the multiset of TileInstanceIds in
// state equals the expected full 106-tile deck. It is pure and returns a
// descriptive error pinpointing duplicate or missing IDs.
// allTiles is the authoritative 106 unique TileInstances (e.g. from deck factory);
// only their IDs are compared — colour/rank are not re-checked here beyond
// TileInstance.Validate already done in state.Validate.
func CheckTileConservation(state *RoundState, allTiles []tile.TileInstance) error {
	if len(allTiles) != 106 {
		return fmt.Errorf("allTiles must have 106 entries, got %d", len(allTiles))
	}
	expected := make(map[tile.TileInstanceId]tile.TileInstance, 106)
	for _, t := range allTiles {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("allTiles tile %v invalid: %w", t.ID, err)
		}
		if _, dup := expected[t.ID]; dup {
			return fmt.Errorf("allTiles duplicate ID %v", t.ID)
		}
		expected[t.ID] = t
	}

	seen := make(map[tile.TileInstanceId]string, 106) // id -> location description
	// helper to record a tile at a location
	record := func(t tile.TileInstance, loc string) error {
		if _, dup := seen[t.ID]; dup {
			return fmt.Errorf("duplicate tile %v: already seen at %v, again at %v", t.ID, seen[t.ID], loc)
		}
		if _, ok := expected[t.ID]; !ok {
			return fmt.Errorf("tile %v at %v not in expected deck", t.ID, loc)
		}
		seen[t.ID] = loc
		return nil
	}

	// Racks
	for _, p := range state.Players {
		rack := state.Racks[p.Seat]
		for i, t := range rack {
			if err := record(t, fmt.Sprintf("rack seat %v index %d", p.Seat, i)); err != nil {
				return err
			}
		}
	}
	// Stock
	for i, t := range state.Stock {
		if err := record(t, fmt.Sprintf("stock index %d", i)); err != nil {
			return err
		}
	}
	// DiscardRow
	for i, d := range state.DiscardRow {
		if err := record(d.Tile, fmt.Sprintf("discard index %d (IsOpening=%v)", i, d.IsOpeningDiscard)); err != nil {
			return err
		}
	}
	// TableMelds
	for _, m := range state.TableMelds {
		for _, t := range m.Tiles {
			if err := record(t, fmt.Sprintf("meld %q", m.ID)); err != nil {
				return err
			}
		}
	}

	if len(seen) != 106 {
		// find missing
		for id := range expected {
			if _, ok := seen[id]; !ok {
				return fmt.Errorf("missing tile %v: expected in deck but not found in state (seen %d/106)", id, len(seen))
			}
		}
		return fmt.Errorf("saw %d tiles, expected 106", len(seen))
	}
	return nil
}

// CountTiles returns counts per location — useful for tests/logs.
func CountTiles(state *RoundState) (racks, stock, discard, melds int) {
	for _, p := range state.Players {
		racks += len(state.Racks[p.Seat])
	}
	stock = len(state.Stock)
	discard = len(state.DiscardRow)
	for _, m := range state.TableMelds {
		melds += len(m.Tiles)
	}
	return
}
