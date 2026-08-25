// Package match — Discard pickup protection (Day 36).
// Ensures the opening blocked tile is never selectable.
package match

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
)

// CanPickupPreviousDiscard reports whether the previous discard (last entry)
// may be taken via DRAW_PREVIOUS_DISCARD per docs/rules-decisions.md:1.5 and
// Day 36 protection. It checks that a discard exists and is not the opening
// blocked tile.
func CanPickupPreviousDiscard(state *RoundState) *protocol.ErrorResponse {
	if len(state.DiscardRow) == 0 {
		return protocol.NewError(protocol.ErrCodeBadRequest, "no discard to pick up", "", nil)
	}
	last := state.DiscardRow[len(state.DiscardRow)-1]
	if last.IsOpeningDiscard {
		return protocol.NewError(protocol.ErrCodeBadRequest, "opening discard is blocked and cannot be picked up", "", nil)
	}
	if state.GamePhase != PhasePlaying || state.TurnPhase != TurnMustDraw {
		return protocol.NewError(protocol.ErrCodeWrongPhase, fmt.Sprintf("previous discard only in Playing MustDraw, got %v/%v", state.GamePhase, state.TurnPhase), "", nil)
	}
	return nil
}

// CanPickupDiscardForMeld checks whether discard at index may be taken via
// PICKUP_DISCARD_FOR_MELD (needs exactly two rack tiles, immediate meld).
// It rejects IsOpeningDiscard and out-of-range indices.
func CanPickupDiscardForMeld(state *RoundState, discardIndex int) *protocol.ErrorResponse {
	if discardIndex < 0 || discardIndex >= len(state.DiscardRow) {
		return protocol.NewError(protocol.ErrCodeBadPayload, fmt.Sprintf("discardIndex %d out of range 0..%d", discardIndex, len(state.DiscardRow)-1), "", nil)
	}
	entry := state.DiscardRow[discardIndex]
	if entry.IsOpeningDiscard {
		return protocol.NewError(protocol.ErrCodeBadPayload, "opening discard cannot be selected for meld pickup", "", nil)
	}
	if state.GamePhase != PhasePlaying || state.TurnPhase != TurnMustDraw {
		return protocol.NewError(protocol.ErrCodeWrongPhase, fmt.Sprintf("pickup only in Playing MustDraw, got %v/%v", state.GamePhase, state.TurnPhase), "", nil)
	}
	return nil
}
