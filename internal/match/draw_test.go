package match

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
	"github.com/heroiclabs/nakama-common/runtime"
)

func playingStateWithStock() (*RoundState, []tile.TileInstance) {
	players, _ := AssignSeats([]PlayerId{"alice", "bob"})
	// Use small deterministic stock and racks: alice 14, bob 14, stock 2, discard 0, current alice MustDraw
	rack0 := []tile.TileInstance{tile.MustTile("a1", tile.Red, 1), tile.MustTile("a2", tile.Red, 2)}
	rack1 := []tile.TileInstance{tile.MustTile("b1", tile.Blue, 3)}
	stock := []tile.TileInstance{tile.MustTile("s1", tile.Yellow, 5), tile.MustTile("s2", tile.Black, 6)}
	// Pad to 14 each for realism
	for len(rack0) < 14 {
		rack0 = append(rack0, tile.MustTile(tile.TileInstanceId(fmt.Sprintf("a-pad-%02d", len(rack0))), tile.Red, tile.Rank(1+len(rack0)%13)))
	}
	for len(rack1) < 14 {
		rack1 = append(rack1, tile.MustTile(tile.TileInstanceId(fmt.Sprintf("b-pad-%02d", len(rack1))), tile.Blue, tile.Rank(1+len(rack1)%13)))
	}
	// Ensure unique IDs for padded tiles by using fmt
	// For test we already have distinct, but we will just use existing
	// Stock already 2
	state := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: rack0, 1: rack1},
		Stock:       stock,
		DiscardRow:  []DiscardEntry{},
		TableMelds:  []TableMeld{},
		CurrentSeat: 0,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	// allTiles for conservation: racks+stock = 14+14+2=30, need to create a 30 deck for check
	allTiles := append([]tile.TileInstance(nil), rack0...)
	allTiles = append(allTiles, rack1...)
	allTiles = append(allTiles, stock...)
	// For conservation we just need those 30 as deck, but CheckTileConservation expects 106.
	// Instead we will just check counts, not full conservation, for this small state.
	return state, allTiles
}

func TestDrawStockSuccess(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _ := playingStateWithStock()
	// Alice is current MustDraw, stock has 2
	initialStock := len(state.Stock)
	initialRack := len(state.Racks[0])
	dispatcher := &mockDispatcher{}
	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawStock, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawStock, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.Stock) != initialStock-1 {
		t.Fatalf("stock %d want %d", len(st.Stock), initialStock-1)
	}
	if len(st.Racks[0]) != initialRack+1 {
		t.Fatalf("rack %d want %d", len(st.Racks[0]), initialRack+1)
	}
	if st.TurnPhase != TurnMeldOrDiscard {
		t.Fatalf("TurnPhase %v want MeldOrDiscard", st.TurnPhase)
	}
	if st.CurrentSeat != 0 {
		t.Fatalf("CurrentSeat should stay same after draw")
	}
	// Ensure drawn tile was top of stock (s2)
	drawnId := st.Racks[0][len(st.Racks[0])-1].ID
	if drawnId != "s2" {
		t.Fatalf("drawn should be top s2, got %v", drawnId)
	}
}

func TestDrawStockRejectsWrongPlayer(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _ := playingStateWithStock()
	dispatcher := &mockDispatcher{}
	// Bob tries to draw while alice is current
	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawStock, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "bob"}, opCode: protocol.OpClientDrawStock, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.Stock) != 2 {
		t.Fatalf("stock should be unchanged on wrong player")
	}
	if len(st.Racks[0]) != 14 {
		t.Fatalf("alice rack should be unchanged")
	}
}

func TestDrawStockRejectsWrongPhase(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _ := playingStateWithStock()
	state.TurnPhase = TurnMeldOrDiscard // already drawn
	dispatcher := &mockDispatcher{}
	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawStock, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawStock, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.Stock) != 2 {
		t.Fatalf("stock should be unchanged when already MeldOrDiscard")
	}
}

func TestDrawStockEmpty(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _ := playingStateWithStock()
	state.Stock = []tile.TileInstance{} // empty
	dispatcher := &mockDispatcher{}
	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawStock, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawStock, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.Stock) != 0 {
		t.Fatalf("empty stock should stay 0")
	}
	if st.TurnPhase != TurnMustDraw {
		t.Fatalf("should stay MustDraw on empty stock error")
	}
}
