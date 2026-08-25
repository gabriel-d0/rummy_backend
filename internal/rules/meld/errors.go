// Package meld — Structured validation errors (Day 45).
package meld

import "fmt"

// ValidationError is a structured reason for an invalid meld, not just bool.
// It allows MatchLoop to map to protocol.ErrCodeBadPayload with field details per AGENTS.md:370.
type ValidationError struct {
	Code    string
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s [%s]: %s", e.Code, e.Field, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common codes for set/run validation.
const (
	ErrCodeInvalidSize     = "invalid_size"
	ErrCodeRankMismatch    = "rank_mismatch"
	ErrCodeDuplicateColour = "duplicate_colour"
	ErrCodeDuplicateTile   = "duplicate_tile"
	ErrCodeInvalidKind     = "invalid_kind"
	ErrCodeJokerNotAllowed = "joker_not_allowed"
)
