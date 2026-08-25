// Package meld — Set validation (Day 44 + Day 46 joker support).
// Pure validation for terta: 3 or 4 tiles, same rank, distinct colours,
// with explicit joker reps and real>=2*joker ratio per docs/rules-decisions.md:1.3.
package meld

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// ValidateSet checks a meld claimed to be a set.
// It expects meld.Kind == KindSet, 3 or 4 tiles, same rank, distinct colours,
// with jokers allowed via JokerReps (explicit, same rank, distinct colour, ratio).
func ValidateSet(m Meld) error {
	if m.Kind != KindSet {
		return &ValidationError{Code: ErrCodeInvalidKind, Field: "kind", Message: fmt.Sprintf("meld %q not a set (kind %q)", m.ID, m.Kind)}
	}
	if len(m.Tiles) < 3 || len(m.Tiles) > 4 {
		return &ValidationError{Code: ErrCodeInvalidSize, Field: "tiles", Message: fmt.Sprintf("set %q must have 3 or 4 tiles, got %d", m.ID, len(m.Tiles))}
	}
	// Count real vs jokers
	realCount, jokerCount := 0, 0
	for _, t := range m.Tiles {
		if t.IsJoker {
			jokerCount++
		} else {
			realCount++
		}
	}
	if realCount < 2*jokerCount {
		return &ValidationError{Code: ErrCodeInvalidSize, Field: "jokerReps", Message: fmt.Sprintf("set %q ratio real %d must be >=2*joker %d", m.ID, realCount, jokerCount)}
	}
	// Determine rank from first real tile, or from first joker rep if no real (should not happen due to ratio, but handle)
	var rank tile.Rank
	foundReal := false
	for _, t := range m.Tiles {
		if !t.IsJoker {
			rank = t.Rank
			foundReal = true
			break
		}
	}
	if !foundReal {
		// All jokers? Not allowed by ratio, but pick first rep's rank for error messaging
		for _, t := range m.Tiles {
			if t.IsJoker {
				if rep, ok := m.JokerReps[t.ID]; ok {
					rank = rep.Rank
					foundReal = true
					break
				}
			}
		}
		if !foundReal {
			return &ValidationError{Code: ErrCodeRankMismatch, Field: "rank", Message: fmt.Sprintf("set %q has no real tile and no joker rep to determine rank", m.ID)}
		}
	}
	if !rank.IsValid() {
		return &ValidationError{Code: ErrCodeRankMismatch, Field: "rank", Message: fmt.Sprintf("set %q invalid rank %v", m.ID, rank)}
	}
	seenColour := map[tile.Colour]bool{}
	for _, t := range m.Tiles {
		if t.IsJoker {
			rep, ok := m.JokerReps[t.ID]
			if !ok {
				return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "jokerReps", Message: fmt.Sprintf("set %q joker %v missing rep", m.ID, t.ID)}
			}
			if rep.IsJoker {
				return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "jokerReps", Message: fmt.Sprintf("set %q joker %v rep cannot be joker", m.ID, t.ID)}
			}
			if rep.Rank != rank {
				return &ValidationError{Code: ErrCodeRankMismatch, Field: "rank", Message: fmt.Sprintf("set %q joker %v rep rank %v vs set rank %v", m.ID, t.ID, rep.Rank, rank)}
			}
			if !rep.Colour.IsValid() {
				return &ValidationError{Code: ErrCodeRankMismatch, Field: "colour", Message: fmt.Sprintf("set %q joker %v rep invalid colour %v", m.ID, t.ID, rep.Colour)}
			}
			if seenColour[rep.Colour] {
				return &ValidationError{Code: ErrCodeDuplicateColour, Field: "colour", Message: fmt.Sprintf("set %q duplicate colour %v (joker rep)", m.ID, rep.Colour)}
			}
			seenColour[rep.Colour] = true
		} else {
			if t.Rank != rank {
				return &ValidationError{Code: ErrCodeRankMismatch, Field: "rank", Message: fmt.Sprintf("set %q rank mismatch %v vs %v", m.ID, t.Rank, rank)}
			}
			if seenColour[t.Colour] {
				return &ValidationError{Code: ErrCodeDuplicateColour, Field: "colour", Message: fmt.Sprintf("set %q duplicate colour %v", m.ID, t.Colour)}
			}
			seenColour[t.Colour] = true
		}
	}
	// Also check that JokerReps doesn't contain rep for non-joker (already checked via not seen, but also need to ensure non-joker not in map)
	for jid := range m.JokerReps {
		found := false
		for _, t := range m.Tiles {
			if t.ID == jid {
				found = true
				if !t.IsJoker {
					return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "jokerReps", Message: fmt.Sprintf("set %q rep for non-joker %v", m.ID, jid)}
				}
				break
			}
		}
		if !found {
			return &ValidationError{Code: ErrCodeJokerNotAllowed, Field: "jokerReps", Message: fmt.Sprintf("set %q rep for non-existent tile %v", m.ID, jid)}
		}
	}
	return nil
}
