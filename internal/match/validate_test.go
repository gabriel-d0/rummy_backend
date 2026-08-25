package match

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestValidateActivePlayer(t *testing.T) {
	players, _ := AssignSeats([]PlayerId{"alice", "bob", "carol"})
	state := &RoundState{
		Players:     players,
		CurrentSeat: 1,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Racks:       map[Seat][]tile.TileInstance{},
	}
	// bob is seat 1 current
	if err := ValidateActivePlayer(state, "bob"); err != nil {
		t.Fatalf("bob should be active: %v", err)
	}
	// alice not active
	if err := ValidateActivePlayer(state, "alice"); err == nil || err.Code != "not_your_turn" {
		t.Fatalf("alice should be not_your_turn, got %v", err)
	}
	// unknown not member
	if err := ValidateActivePlayer(state, "dave"); err == nil || err.Code != "not_member" {
		t.Fatalf("unknown should be not_member, got %v", err)
	}
	// OpeningDiscard also enforces
	state.GamePhase = PhaseOpeningDiscard
	state.CurrentSeat = 0
	if err := ValidateActivePlayer(state, "alice"); err != nil {
		t.Fatalf("alice opening should be active")
	}
	if err := ValidateActivePlayer(state, "bob"); err == nil {
		t.Fatalf("bob opening should be not_your_turn")
	}
	// Waiting does not enforce (only start via host, handled separately)
	state.GamePhase = PhaseWaiting
	state.CurrentSeat = SeatInvalid
	if err := ValidateActivePlayer(state, "bob"); err != nil {
		t.Fatalf("Waiting should not enforce active, got %v", err)
	}
}

func TestValidatePhaseOp(t *testing.T) {
	state := &RoundState{GamePhase: PhaseWaiting, TurnPhase: TurnMustDraw}
	if err := ValidatePhaseOp(state, 1); err != nil {
		t.Fatalf("Waiting start should be allowed")
	}
	if err := ValidatePhaseOp(state, 2); err == nil {
		t.Fatalf("Waiting discard should be wrong_phase")
	}
	state.GamePhase = PhaseOpeningDiscard
	if err := ValidatePhaseOp(state, 2); err != nil {
		t.Fatalf("OpeningDiscard discard should be allowed")
	}
	if err := ValidatePhaseOp(state, 3); err == nil {
		t.Fatalf("OpeningDiscard draw should be wrong_phase")
	}
	state.GamePhase = PhasePlaying
	state.TurnPhase = TurnMustDraw
	if err := ValidatePhaseOp(state, 3); err != nil {
		t.Fatalf("MustDraw draw should be allowed")
	}
	if err := ValidatePhaseOp(state, 2); err == nil {
		t.Fatalf("MustDraw discard should be wrong_phase")
	}
	state.TurnPhase = TurnMeldOrDiscard
	if err := ValidatePhaseOp(state, 2); err != nil {
		t.Fatalf("MeldOrDiscard discard should be allowed")
	}
	if err := ValidatePhaseOp(state, 3); err == nil {
		t.Fatalf("MeldOrDiscard draw should be wrong_phase")
	}
}
