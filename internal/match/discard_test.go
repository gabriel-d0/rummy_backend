package match

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func openingStateWithOneBlocked() *RoundState {
	players, _ := AssignSeats([]PlayerId{"alice", "bob"})
	// Alice has 14 after opening discard, bob 14, discard has opening, stock some
	return &RoundState{
		Players: players,
		Racks: map[Seat][]tile.TileInstance{
			0: {tile.MustTile("a1", tile.Red, 1), tile.MustTile("a2", tile.Red, 2)},
			1: {tile.MustTile("b1", tile.Blue, 3)},
		},
		Stock: []tile.TileInstance{tile.MustTile("s1", tile.Yellow, 3)},
		DiscardRow: []DiscardEntry{
			{Tile: tile.MustTile("disc-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0},
		},
		TableMelds:  []TableMeld{},
		CurrentSeat: 1,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
}

func TestCanPickupPreviousDiscardBlocked(t *testing.T) {
	state := openingStateWithOneBlocked()
	if err := CanPickupPreviousDiscard(state); err == nil {
		t.Fatalf("should reject opening discard as previous")
	}
	// Add a normal discard on top — now previous is not opening
	state.DiscardRow = append(state.DiscardRow, DiscardEntry{Tile: tile.MustTile("disc-normal", tile.Blue, 5), IsOpeningDiscard: false, Index: 1})
	if err := CanPickupPreviousDiscard(state); err != nil {
		t.Fatalf("normal previous should be allowed, got %v", err)
	}
	// Wrong phase
	state.GamePhase = PhaseOpeningDiscard
	if err := CanPickupPreviousDiscard(state); err == nil {
		t.Fatalf("should reject wrong phase")
	}
}

func TestCanPickupDiscardForMeldBlocked(t *testing.T) {
	state := openingStateWithOneBlocked()
	// Index 0 is opening — should be rejected
	if err := CanPickupDiscardForMeld(state, 0); err == nil {
		t.Fatalf("should reject opening index 0")
	}
	// Add normal at 1
	state.DiscardRow = append(state.DiscardRow, DiscardEntry{Tile: tile.MustTile("disc-1", tile.Yellow, 3), IsOpeningDiscard: false, Index: 1})
	if err := CanPickupDiscardForMeld(state, 1); err != nil {
		t.Fatalf("index 1 should be allowed, got %v", err)
	}
	if err := CanPickupDiscardForMeld(state, 0); err == nil {
		t.Fatalf("index 0 still blocked")
	}
	// Out of range
	if err := CanPickupDiscardForMeld(state, 5); err == nil {
		t.Fatalf("out of range should fail")
	}
	// Wrong phase
	state.GamePhase = PhaseWaiting
	if err := CanPickupDiscardForMeld(state, 1); err == nil {
		t.Fatalf("wrong phase should fail")
	}
}
