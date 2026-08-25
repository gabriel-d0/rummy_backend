// Package meld — Set validation (Day 44).
// Pure validation for terta: 3 or 4 tiles, same rank, distinct colours.
// Joker support is Day 46; here we validate no-joker sets and reject
// jokers with a clear error so Day 46 can extend.
package meld

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// ValidateSet checks a meld claimed to be a set.
// It expects meld.Kind == KindSet, 3 or 4 tiles, same rank, distinct colours,
// no jokers (jokers rejected here; Day 46 will allow with reps).
// It returns *ValidationError with Code/Field for structured handling per Day 45.
func ValidateSet(m Meld) error {
	if m.Kind != KindSet {
		return &ValidationError{Code: ErrCodeInvalidKind, Field: "kind", Message: fmt.Sprintf("meld %q not a set (kind %q)", m.ID, m.Kind)}
	}
	if len(m.Tiles) < 3 || len(m.Tiles) > 4 {
		return &ValidationError{Code: ErrCodeInvalidSize, Field: "tiles", Message: fmt.Sprintf("set %q must have 3 or 4 tiles, got %d", m.ID, len(m.Tiles))}
	}
	if len(m.JokerReps) != 0 {
		return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "jokerReps", Message: fmt.Sprintf("set %q has jokers, use ValidateSetWithJoker (Day 46)", m.ID)}
	}
	// Check all tiles are numbered and distinct colours, same rank
	if len(m.Tiles) == 0 {
		return &ValidationError{Code: ErrCodeInvalidSize, Field: "tiles", Message: fmt.Sprintf("set %q empty", m.ID)}
	}
	rank := m.Tiles[0].Rank
	if !rank.IsValid() {
		return &ValidationError{Code: ErrCodeRankMismatch, Field: "rank", Message: fmt.Sprintf("set %q invalid rank %v", m.ID, rank)}
	}
	seenColour := map[tile.Colour]bool{}
	for _, t := range m.Tiles {
		if t.IsJoker {
			return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "tiles", Message: fmt.Sprintf("set %q has joker %v, not allowed in basic ValidateSet", m.ID, t.ID)}
		}
		if t.Rank != rank {
			return &ValidationError{Code: ErrCodeRankMismatch, Field: "rank", Message: fmt.Sprintf("set %q rank mismatch %v vs %v", m.ID, t.Rank, rank)}
		}
		if seenColour[t.Colour] {
			return &ValidationError{Code: ErrCodeDuplicateColour, Field: "colour", Message: fmt.Sprintf("set %q duplicate colour %v", m.ID, t.Colour)}
		}
		seenColour[t.Colour] = true
	}
	return nil
}
