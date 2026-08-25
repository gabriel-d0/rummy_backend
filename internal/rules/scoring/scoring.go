// Package scoring — Tile scoring model (Day 56).
// Pure scoring for Romanian Tile Rummy opening meld per docs/rules-decisions.md:1.4.
package scoring

import "github.com/gabriel-d0/rummy_backend/internal/rules/tile"

// ScoreTile returns the point value of a tile in a given context.
// For numbered tiles 2-9 =>5, 10-13 =>10.
// For Ace, context matters:
//   - In a low run 1-2-3 (isLowAceRun true) =>5
//   - In a high run 12-13-1 (isHighAceRun true) =>10
//   - In a set of three Aces (isAceSet true) =>25 each (overrides 5/10)
//
// For Joker, jokerRep must be non-nil and ScoreTile is called on the represented tile's value
// (joker's own colour/rank is ignored).
func ScoreTile(t tile.TileInstance, isLowAceRun, isHighAceRun, isAceSet bool, jokerRep *tile.TileInstance) int {
	if t.IsJoker {
		if jokerRep == nil {
			return 0
		}
		// Joker value equals represented tile's value in its context
		return ScoreTile(*jokerRep, isLowAceRun, isHighAceRun, isAceSet, nil)
	}
	// AceSet special: 25 each (all tiles are Aces in a set of 3)
	if isAceSet && t.Rank == tile.RankAce {
		return 25
	}
	if t.Rank == tile.RankAce {
		if isHighAceRun {
			return 10
		}
		if isLowAceRun {
			return 5
		}
		// Bare Ace not in run context? Default to 5 for safety, but callers should specify.
		return 5
	}
	if t.Rank >= 2 && t.Rank <= 9 {
		return 5
	}
	if t.Rank >= 10 && t.Rank <= 13 {
		return 10
	}
	return 0
}

// ScoreTileSimple is a helper for non-Ace, non-joker contexts.
func ScoreTileSimple(t tile.TileInstance) int {
	return ScoreTile(t, false, false, false, nil)
}
