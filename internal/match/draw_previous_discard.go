// Package match — Previous discard pickup handler (AGENTS.md Day 16).
// Implements DRAW_PREVIOUS_DISCARD: opened active player in MustDraw may take
// the immediately previous discard (not the opening blocked tile) into rack.
package match

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/heroiclabs/nakama-common/runtime"
)

// handleDrawPreviousDiscard validates and applies DRAW_PREVIOUS_DISCARD.
// It is atomic and moves exactly one tile from DiscardRow tail to rack.
func handleDrawPreviousDiscard(st *RoundState, senderId PlayerId, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) error {
	if st.GamePhase != PhasePlaying || st.TurnPhase != TurnMustDraw {
		sendError(dispatcher, sender, protocol.ErrCodeWrongPhase, fmt.Sprintf("previous discard only in Playing MustDraw, got %v/%v", st.GamePhase, st.TurnPhase), requestId, op, logger)
		return fmt.Errorf("wrong phase")
	}
	seat := SeatOfPlayer(st.Players, senderId)
	if seat == SeatInvalid {
		sendError(dispatcher, sender, protocol.ErrCodeNotMember, "not in match", requestId, op, logger)
		return fmt.Errorf("not member")
	}
	if seat != st.CurrentSeat {
		sendError(dispatcher, sender, protocol.ErrCodeNotYourTurn, fmt.Sprintf("not your turn: current %v sender %v", st.CurrentSeat, seat), requestId, op, logger)
		return fmt.Errorf("not your turn")
	}
	// Find player and check HasOpened
	idx := findPlayerIndex(st, senderId)
	if idx == -1 {
		sendError(dispatcher, sender, protocol.ErrCodeNotMember, "player not found", requestId, op, logger)
		return fmt.Errorf("player not found")
	}
	if !st.Players[idx].HasOpened {
		sendError(dispatcher, sender, protocol.ErrCodeNotOpened, "must open with MELD_INITIAL first (50+ with run)", requestId, op, logger)
		return fmt.Errorf("not opened")
	}
	if perr := CanPickupPreviousDiscard(st); perr != nil {
		// Map to appropriate code; CanPickupPreviousDiscard returns bad_request or wrong_phase
		sendError(dispatcher, sender, perr.Code, perr.Message, requestId, op, logger)
		return fmt.Errorf("%s: %s", perr.Code, perr.Message)
	}
	// Move last discard to rack
	lastIdx := len(st.DiscardRow) - 1
	last := st.DiscardRow[lastIdx]
	// Remove from discard row
	st.DiscardRow = st.DiscardRow[:lastIdx]
	// Append to rack
	rack := st.Racks[seat]
	rack = append(rack, last.Tile)
	st.Racks[seat] = rack
	st.TurnPhase = TurnMeldOrDiscard

	logger.Info("DrawPreviousDiscard by %s seat %v tile %v discard %d->%d rack %d", senderId, seat, last.Tile.ID, lastIdx+1, len(st.DiscardRow), len(rack))
	if dispatcher != nil {
		_ = dispatcher.BroadcastMessage(protocol.OpServerEvent, []byte(fmt.Sprintf(`{"op":"drawPreviousDiscard","seat":%d,"tileId":"%s","discardCount":%d}`, int(seat), last.Tile.ID, len(st.DiscardRow))), nil, nil, true)
	}
	return nil
}
