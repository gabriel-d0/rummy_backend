// Package meld — Run validation (Day 48).
// Basic same-colour consecutive runs of length >=3, no jokers (Day 52 will add).
package meld

import (
	"fmt"
	"sort"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// ValidateRun checks a meld claimed to be a run.
// It expects KindRun, >=3 tiles, same colour, consecutive ranks, with
// explicit joker reps (same colour, gap filling) and real>=2*joker ratio
// per docs/rules-decisions.md:1.3. Handles low Ace 1-2-3 and high Ace 12-13-1.
func ValidateRun(m Meld) error {
	if m.Kind != KindRun {
		return &ValidationError{Code: ErrCodeInvalidKind, Field: "kind", Message: fmt.Sprintf("meld %q not a run (kind %q)", m.ID, m.Kind)}
	}
	if len(m.Tiles) < 3 {
		return &ValidationError{Code: ErrCodeInvalidSize, Field: "tiles", Message: fmt.Sprintf("run %q must have >=3 tiles, got %d", m.ID, len(m.Tiles))}
	}
	// Count real vs joker and check ratio real>=2*joker
	realCount, jokerCount := 0, 0
	for _, t := range m.Tiles {
		if t.IsJoker {
			jokerCount++
		} else {
			realCount++
		}
	}
	if realCount < 2*jokerCount {
		return &ValidationError{Code: ErrCodeInvalidSize, Field: "jokerReps", Message: fmt.Sprintf("run %q ratio real %d must be >=2*joker %d", m.ID, realCount, jokerCount)}
	}
	// Determine run colour from first real tile (ratio ensures at least 1 real)
	var colour tile.Colour
	foundReal := false
	for _, t := range m.Tiles {
		if !t.IsJoker {
			colour = t.Colour
			foundReal = true
			break
		}
	}
	if !foundReal {
		// All jokers not allowed by ratio, but pick first rep's colour for error
		for _, t := range m.Tiles {
			if t.IsJoker {
				if rep, ok := m.JokerReps[t.ID]; ok {
					colour = rep.Colour
					foundReal = true
					break
				}
			}
		}
		if !foundReal {
			return &ValidationError{Code: ErrCodeRankMismatch, Field: "colour", Message: fmt.Sprintf("run %q has no real tile to determine colour", m.ID)}
		}
	}
	if !colour.IsValid() {
		return &ValidationError{Code: ErrCodeRankMismatch, Field: "colour", Message: fmt.Sprintf("run %q invalid colour %v", m.ID, colour)}
	}
	ranks := make([]int, 0, len(m.Tiles))
	for _, t := range m.Tiles {
		var effColour tile.Colour
		var effRank tile.Rank
		if t.IsJoker {
			rep, ok := m.JokerReps[t.ID]
			if !ok {
				return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "jokerReps", Message: fmt.Sprintf("run %q joker %v missing rep", m.ID, t.ID)}
			}
			if rep.IsJoker {
				return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "jokerReps", Message: fmt.Sprintf("run %q joker %v rep cannot be joker", m.ID, t.ID)}
			}
			if !rep.Colour.IsValid() || !rep.Rank.IsValid() {
				return &ValidationError{Code: ErrCodeRankMismatch, Field: "jokerReps", Message: fmt.Sprintf("run %q joker %v rep invalid", m.ID, t.ID)}
			}
			effColour = rep.Colour
			effRank = rep.Rank
		} else {
			if !t.Rank.IsValid() {
				return &ValidationError{Code: ErrCodeRankMismatch, Field: "rank", Message: fmt.Sprintf("run %q invalid rank %v", m.ID, t.Rank)}
			}
			effColour = t.Colour
			effRank = t.Rank
		}
		if effColour != colour {
			return &ValidationError{Code: "colour_mismatch", Field: "colour", Message: fmt.Sprintf("run %q colour %v vs %v", m.ID, effColour, colour)}
		}
		ranks = append(ranks, int(effRank))
	}
	// Validate JokerReps doesn't contain rep for non-joker
	for jid := range m.JokerReps {
		found := false
		for _, t := range m.Tiles {
			if t.ID == jid {
				found = true
				if !t.IsJoker {
					return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "jokerReps", Message: fmt.Sprintf("run %q rep for non-joker %v", m.ID, jid)}
				}
				break
			}
		}
		if !found {
			return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "jokerReps", Message: fmt.Sprintf("run %q rep for non-existent tile %v", m.ID, jid)}
		}
	}
	sort.Ints(ranks)
	// Check low-Ace consecutive (1,2,3 ...) — already sorted as 1,2,3
	isLowConsecutive := true
	for i := 1; i < len(ranks); i++ {
		if ranks[i] == ranks[i-1] {
			return &ValidationError{Code: ErrCodeDuplicateTile, Field: "rank", Message: fmt.Sprintf("run %q duplicate rank %d", m.ID, ranks[i])}
		}
		if ranks[i] != ranks[i-1]+1 {
			isLowConsecutive = false
			break
		}
	}
	if isLowConsecutive {
		return nil
	}
	// Check high-Ace: treat Ace (1) as 14, so ranks become n..13,14
	// High-Ace is valid only if ranks contain 1 and the non-Ace ranks are consecutive ending at 13.
	if contains(ranks, 1) {
		// Build high ranks: replace 1 with 14
		highRanks := make([]int, 0, len(ranks))
		for _, r := range ranks {
			if r == 1 {
				highRanks = append(highRanks, 14)
			} else {
				highRanks = append(highRanks, r)
			}
		}
		sort.Ints(highRanks)
		// High-Ace must be consecutive and contain 14, and not also contain 2 with 14 (which would be Ace in middle)
		// Check consecutive
		highConsecutive := true
		for i := 1; i < len(highRanks); i++ {
			if highRanks[i] == highRanks[i-1] {
				highConsecutive = false
				break
			}
			if highRanks[i] != highRanks[i-1]+1 {
				highConsecutive = false
				break
			}
		}
		if highConsecutive {
			// Ensure Ace is at an end, not middle: for high-Ace, Ace as 14 must be last (largest)
			// Our highRanks sorted already has 14 last if consecutive, but check that original low gap wasn't 13-1-2 style
			// For 12-13-1, highRanks is 12,13,14 consecutive true, and low check already failed (1,12,13 gap), so high passes.
			// For 13-1-2, low ranks 1,2,13 gap fails, highRanks is 2,13,14 not consecutive, so fails.
			// For 12-13-1-2, ranks 1,2,12,13 low fails, highRanks 2,12,13,14 not consecutive, fails.
			return nil
		}
	}
	return &ValidationError{Code: ErrCodeRankMismatch, Field: "rank", Message: fmt.Sprintf("run %q not consecutive %v", m.ID, ranks)}
}

func contains(ranks []int, v int) bool {
	for _, r := range ranks {
		if r == v {
			return true
		}
	}
	return false
}
