// Package match — Win detection (AGENTS.md Day 19, docs/rules-decisions.md:6.1).
// MVP decision: rack == 0 after any legal state-mutating command in MeldOrDiscard
// (DISCARD or MELD/EXTEND/REPLACE/PICKUP) is a win, whether via final discard
// or via melding that empties the rack without discard. Checked after every
// successful mutation per docs/rules-decisions.md:6.1.
package match

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
)

// checkWinAndComplete checks if the actor's rack is empty and, if so, transitions
// to RoundComplete with Winner set. It broadcasts a final OpServerEvent and returns true if the round ended.
// It must be called after the successful atomic mutation that removed tiles from rack.
func checkWinAndComplete(st *RoundState, seat Seat, dispatcher runtime.MatchDispatcher, logger runtime.Logger) bool {
	if st.GamePhase != PhasePlaying {
		return false
	}
	rack := st.Racks[seat]
	if len(rack) != 0 {
		return false
	}
	// Ensure player is still in match and has opened? For win, opening must have happened, but we allow any empty rack as win for MVP
	// Per docs/rules-decisions.md:6.1, we allow win without final discard
	st.GamePhase = PhaseRoundComplete
	st.Winner = seat
	// Keep TurnPhase as is, but no further gameplay allowed via AllowedOps (RoundComplete has none)
	// CurrentSeat stays as winner for final snapshot
	st.CurrentSeat = seat
	logger.Info("RoundComplete winner seat %v (%s) table %d", seat, st.Players[seatForIdx(st, seat)].ID, len(st.TableMelds))
	if dispatcher != nil {
		_ = dispatcher.BroadcastMessage(101, []byte(fmt.Sprintf(`{"phase":"RoundComplete","winner":%d,"seat":%d}`, int(seat), int(seat))), nil, nil, true)
		_ = dispatcher.BroadcastMessage(103, []byte(fmt.Sprintf(`{"op":"roundComplete","winner":%d}`, int(seat))), nil, nil, true)
	}
	return true
}

func seatForIdx(st *RoundState, seat Seat) int {
	for i, p := range st.Players {
		if p.Seat == seat {
			return i
		}
	}
	return 0
}
