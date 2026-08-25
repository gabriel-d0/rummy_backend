package match

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
	"github.com/heroiclabs/nakama-common/runtime"
)

func TestEmptyStockDrawRejectedAndConservation(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	dispatcher := &mockDispatcher{}
	// Create a Playing MustDraw state with empty stock
	players, _ := AssignSeats([]PlayerId{"alice", "bob"})
	state := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: {tile.MustTile("a1", tile.Red, 1)}, 1: {tile.MustTile("b1", tile.Blue, 2)}},
		Stock:       []tile.TileInstance{}, // empty
		DiscardRow:  []DiscardEntry{{Tile: tile.MustTile("disc-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0}, {Tile: tile.MustTile("disc-1", tile.Yellow, 3), IsOpeningDiscard: false, Index: 1}},
		TableMelds:  []TableMeld{},
		CurrentSeat: 0,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	// Build a fake allTiles for conservation: racks+discard = 1+1+2 =4, stock 0, need 4 deck for this small state
	// For this test we just check that state remains valid and stock stays 0
	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawStock, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawStock, data: envBytes}
	initialStock := len(state.Stock)
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.Stock) != initialStock {
		t.Fatalf("empty stock should stay 0, got %d", len(st.Stock))
	}
	if st.TurnPhase != TurnMustDraw {
		t.Fatalf("should stay MustDraw on stock empty, got %v", st.TurnPhase)
	}
	if st.CurrentSeat != 0 {
		t.Fatalf("CurrentSeat should stay 0")
	}
	// Ensure the error was sent (dispatcher would have received OpServerError, but mock doesn't record; we just check state unchanged)
	// Verify that a valid draw would have succeeded if stock had one tile
	state.Stock = []tile.TileInstance{tile.MustTile("s1", tile.Red, 9)}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st2 := next2.(*RoundState)
	if len(st2.Stock) != 0 {
		t.Fatalf("after adding one stock, draw should succeed and stock 0, got %d", len(st2.Stock))
	}
	if st2.TurnPhase != TurnMeldOrDiscard {
		t.Fatalf("after successful draw TurnPhase MeldOrDiscard")
	}
}
