// Package meld — Run validation (Day 48).
// Basic same-colour consecutive runs of length >=3, no jokers (Day 52 will add).
package meld

import (
	"fmt"
	"sort"
)

// ValidateRun checks a meld claimed to be a run.
// It expects KindRun, >=3 tiles, same colour, consecutive ranks.
// Handles low Ace 1-2-3 (as 1,2,3) and high Ace 12-13-1 (as 12,13,14)
// per docs/rules-decisions.md:1.3; 13-1-2 is invalid (Ace in middle).
func ValidateRun(m Meld) error {
	if m.Kind != KindRun {
		return &ValidationError{Code: ErrCodeInvalidKind, Field: "kind", Message: fmt.Sprintf("meld %q not a run (kind %q)", m.ID, m.Kind)}
	}
	if len(m.Tiles) < 3 {
		return &ValidationError{Code: ErrCodeInvalidSize, Field: "tiles", Message: fmt.Sprintf("run %q must have >=3 tiles, got %d", m.ID, len(m.Tiles))}
	}
	if len(m.JokerReps) != 0 {
		return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "jokerReps", Message: fmt.Sprintf("run %q has jokers, use ValidateRunWithJoker (Day 52)", m.ID)}
	}
	// All tiles must be numbered and same colour
	colour := m.Tiles[0].Colour
	if !colour.IsValid() {
		return &ValidationError{Code: ErrCodeRankMismatch, Field: "colour", Message: fmt.Sprintf("run %q invalid colour %v", m.ID, colour)}
	}
	ranks := make([]int, 0, len(m.Tiles))
	for _, t := range m.Tiles {
		if t.IsJoker {
			return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "tiles", Message: fmt.Sprintf("run %q has joker %v", m.ID, t.ID)}
		}
		if t.Colour != colour {
			return &ValidationError{Code: "colour_mismatch", Field: "colour", Message: fmt.Sprintf("run %q colour %v vs %v", m.ID, t.Colour, colour)}
		}
		if !t.Rank.IsValid() {
			return &ValidationError{Code: ErrCodeRankMismatch, Field: "rank", Message: fmt.Sprintf("run %q invalid rank %v", m.ID, t.Rank)}
		}
		ranks = append(ranks, int(t.Rank))
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
