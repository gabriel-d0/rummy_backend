// Package meld — Run validation (Day 48).
// Basic same-colour consecutive runs of length >=3, no jokers (Day 52 will add).
package meld

import (
	"fmt"
	"sort"
)

// ValidateRun checks a meld claimed to be a run.
// It expects KindRun, >=3 tiles, same colour, consecutive ranks 1..13
// without jokers. Ace handling for Day 49-51 (low/high) is deferred: here
// 1-2-3 consecutive is naturally valid via sorted ranks, but 12-13-1 and
// 13-1-2 are not yet specially handled — 12-13-1 will be considered
// non-consecutive here and will be allowed only after Day 50.
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
	// Check consecutive without gaps and no duplicates
	for i := 1; i < len(ranks); i++ {
		if ranks[i] == ranks[i-1] {
			return &ValidationError{Code: ErrCodeDuplicateTile, Field: "rank", Message: fmt.Sprintf("run %q duplicate rank %d", m.ID, ranks[i])}
		}
		if ranks[i] != ranks[i-1]+1 {
			return &ValidationError{Code: ErrCodeRankMismatch, Field: "rank", Message: fmt.Sprintf("run %q not consecutive %v", m.ID, ranks)}
		}
	}
	return nil
}
