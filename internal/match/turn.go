// Package match — Turn advance (Day 37).
package match

import "fmt"

// AdvanceTurn moves CurrentSeat anticlockwise to the next seat and resets to MustDraw.
// It is used after OpeningDiscard and after normal DISCARD in Playing.
// It validates that there are 2..4 players and CurrentSeat is valid.
func AdvanceTurn(state *RoundState) error {
	n := len(state.Players)
	if n < 2 || n > 4 {
		return fmt.Errorf("player count %d must be 2..4", n)
	}
	if !state.CurrentSeat.IsValid(n) {
		return fmt.Errorf("current seat %v invalid for n=%d", state.CurrentSeat, n)
	}
	next, err := NextSeat(state.CurrentSeat, n)
	if err != nil {
		return err
	}
	state.CurrentSeat = next
	state.TurnPhase = TurnMustDraw
	return nil
}
