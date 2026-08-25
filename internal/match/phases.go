// Package match — Turn state machine (Day 32).
// Defines explicit phases and allowed commands per docs/terminology.md and
// docs/rules-decisions.md:5. Pure helpers, no I/O.
package match

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
)

// AllowedOps returns the set of client opcodes allowed in the given phase/tick.
// It is used by MatchLoop to reject out-of-phase actions before any state mutation
// per AGENTS.md:176. Day 32 defines the matrix; enforcement lands Day 33-34.
func AllowedOps(gamePhase GamePhase, turnPhase TurnPhase) map[int64]bool {
	ops := map[int64]bool{}
	switch gamePhase {
	case PhaseWaiting:
		// In Waiting only start is allowed (host). Join/leave handled via MatchJoinAttempt/Join, not loop.
		ops[protocol.OpClientStart] = true
	case PhaseOpeningDiscard:
		// Opening player must DISCARD one tile from 15
		ops[protocol.OpClientDiscard] = true
	case PhasePlaying:
		switch turnPhase {
		case TurnMustDraw:
			ops[protocol.OpClientDrawStock] = true
			ops[protocol.OpClientDrawPreviousDiscard] = true
			ops[protocol.OpClientPickupDiscardForMeld] = true
		case TurnMeldOrDiscard:
			ops[protocol.OpClientDiscard] = true // must end turn
			ops[protocol.OpClientMeldInitial] = true
			ops[protocol.OpClientMeldNew] = true
			ops[protocol.OpClientExtendMeld] = true
			ops[protocol.OpClientReplaceJoker] = true
		}
	case PhaseRoundComplete:
		// No gameplay ops
	}
	return ops
}

// ValidatePhase checks that gamePhase/turnPhase combo is legal.
func ValidatePhase(gamePhase GamePhase, turnPhase TurnPhase) error {
	if !gamePhase.IsValid() {
		return fmt.Errorf("invalid GamePhase %v", gamePhase)
	}
	if gamePhase == PhasePlaying && !turnPhase.IsValid() {
		return fmt.Errorf("invalid TurnPhase %v for Playing", turnPhase)
	}
	if gamePhase != PhasePlaying && turnPhase != TurnMustDraw {
		// For Waiting/OpeningDiscard/RoundComplete we expect TurnMustDraw as placeholder;
		// this keeps RoundState.Validate simple. Day 34 will tighten.
	}
	return nil
}
