package match

import "testing"

func TestAdvanceTurnAnticlockwise(t *testing.T) {
	// 2 players: 0→1→0
	players, _ := AssignSeats([]PlayerId{"alice", "bob"})
	state := &RoundState{Players: players, CurrentSeat: 0, GamePhase: PhasePlaying, TurnPhase: TurnMeldOrDiscard}
	if err := AdvanceTurn(state); err != nil {
		t.Fatalf("AdvanceTurn 0→1: %v", err)
	}
	if state.CurrentSeat != 1 || state.TurnPhase != TurnMustDraw {
		t.Fatalf("after 0→1 got %v %v", state.CurrentSeat, state.TurnPhase)
	}
	if err := AdvanceTurn(state); err != nil {
		t.Fatalf("1→0: %v", err)
	}
	if state.CurrentSeat != 0 {
		t.Fatalf("1→0 got %v", state.CurrentSeat)
	}
	// 3 players: 1→2→0
	players3, _ := AssignSeats([]PlayerId{"a", "b", "c"})
	state.Players = players3
	state.CurrentSeat = 1
	if err := AdvanceTurn(state); err != nil {
		t.Fatalf("3p 1→2: %v", err)
	}
	if state.CurrentSeat != 2 {
		t.Fatalf("3p 1→2 got %v", state.CurrentSeat)
	}
	if err := AdvanceTurn(state); err != nil {
		t.Fatalf("3p 2→0: %v", err)
	}
	if state.CurrentSeat != 0 {
		t.Fatalf("3p 2→0 got %v", state.CurrentSeat)
	}
	// 4 players: 3→0
	players4, _ := AssignSeats([]PlayerId{"a", "b", "c", "d"})
	state.Players = players4
	state.CurrentSeat = 3
	if err := AdvanceTurn(state); err != nil {
		t.Fatalf("4p 3→0: %v", err)
	}
	if state.CurrentSeat != 0 {
		t.Fatalf("4p 3→0 got %v", state.CurrentSeat)
	}
}

func TestAdvanceTurnInvalid(t *testing.T) {
	state := &RoundState{Players: []PlayerState{{ID: "a", Seat: 0}}, CurrentSeat: 0}
	if err := AdvanceTurn(state); err == nil {
		t.Fatalf("should fail for n=1")
	}
	state.Players = []PlayerState{{ID: "a", Seat: 0}, {ID: "b", Seat: 1}}
	state.CurrentSeat = SeatInvalid
	if err := AdvanceTurn(state); err == nil {
		t.Fatalf("should fail for invalid CurrentSeat")
	}
}
