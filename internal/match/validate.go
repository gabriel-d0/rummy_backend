// Package match — Active-player validation (Day 33).
package match

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
)

// ValidateActivePlayer checks that sender is the current player per AGENTS.md:176.
// It returns a *protocol.ErrorResponse with code not_member or not_your_turn
// so MatchLoop can send OpServerError without mutating state.
func ValidateActivePlayer(state *RoundState, senderId PlayerId) *protocol.ErrorResponse {
	seat := SeatOfPlayer(state.Players, senderId)
	if seat == SeatInvalid {
		return protocol.NewError(protocol.ErrCodeNotMember, fmt.Sprintf("player %v not in match", senderId), "", map[string]string{"playerId": string(senderId)})
	}
	// In Waiting/RoundComplete there is no active turn — allow only if caller is checking start (handled separately).
	// For OpeningDiscard and Playing, enforce CurrentSeat.
	if state.GamePhase == PhaseOpeningDiscard || state.GamePhase == PhasePlaying {
		if seat != state.CurrentSeat {
			return protocol.NewError(protocol.ErrCodeNotYourTurn, fmt.Sprintf("not your turn: current %v, sender %v", state.CurrentSeat, seat), "", map[string]string{"currentSeat": fmt.Sprintf("%d", state.CurrentSeat), "senderSeat": fmt.Sprintf("%d", seat)})
		}
	}
	return nil
}

// ValidatePhaseOp checks that op is allowed in current phase/turn per AllowedOps.
func ValidatePhaseOp(state *RoundState, op int64) *protocol.ErrorResponse {
	allowed := AllowedOps(state.GamePhase, state.TurnPhase)
	if !allowed[op] {
		return protocol.NewError(protocol.ErrCodeWrongPhase, fmt.Sprintf("op %d not allowed in %v/%v", op, state.GamePhase, state.TurnPhase), "", map[string]string{"op": fmt.Sprintf("%d", op), "gamePhase": state.GamePhase.String(), "turnPhase": state.TurnPhase.String()})
	}
	return nil
}
