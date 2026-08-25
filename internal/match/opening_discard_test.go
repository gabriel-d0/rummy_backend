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

// helper to create a RoundState in OpeningDiscard with 15/14 racks for 2 players, using distinct tile IDs.
func openingStateWithDeal() (*RoundState, []tile.TileInstance) {
	// Create a synthetic 106 deck with distinct IDs and valid tiles
	var deck []tile.TileInstance
	for i := 0; i < 104; i++ {
		colour := []tile.Colour{tile.Red, tile.Yellow, tile.Blue, tile.Black}[i%4]
		rank := tile.Rank(1 + (i % 13))
		id := tile.TileInstanceId(fmt.Sprintf("t-%03d", i))
		t := tile.MustTile(id, colour, rank)
		deck = append(deck, t)
	}
	deck = append(deck, tile.MustJoker("joker-1"), tile.MustJoker("joker-2"))
	// Now deal 15/14 for 2 players
	players, _ := AssignSeats([]PlayerId{"alice", "bob"})
	rack0 := deck[0:15]
	rack1 := deck[15:29]
	stock := deck[29:106] // 77
	// Copy for allTiles
	allTiles := make([]tile.TileInstance, len(deck))
	copy(allTiles, deck)
	return &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: append([]tile.TileInstance(nil), rack0...), 1: append([]tile.TileInstance(nil), rack1...)},
		Stock:       append([]tile.TileInstance(nil), stock...),
		DiscardRow:  []DiscardEntry{},
		TableMelds:  []TableMeld{},
		CurrentSeat: 0,
		GamePhase:   PhaseOpeningDiscard,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}, allTiles
}

func TestOpeningDiscardSuccess(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, allTiles := openingStateWithDeal()
	dispatcher := &mockDispatcher{}
	// Alice discards first tile of her rack
	tileId := string(state.Racks[0][0].ID)
	payload, _ := json.Marshal(map[string]string{"tileId": tileId})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice", sessionId: "sess-alice"}, opCode: protocol.OpClientDiscard, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.Racks[0]) != 14 {
		t.Fatalf("rack0 len %d want 14", len(st.Racks[0]))
	}
	if len(st.DiscardRow) != 1 {
		t.Fatalf("discard len %d want 1", len(st.DiscardRow))
	}
	if !st.DiscardRow[0].IsOpeningDiscard || st.DiscardRow[0].Index != 0 {
		t.Fatalf("opening discard flag/index wrong %+v", st.DiscardRow[0])
	}
	if st.DiscardRow[0].Tile.ID != tile.TileInstanceId(tileId) {
		t.Fatalf("discard tile %v want %v", st.DiscardRow[0].Tile.ID, tileId)
	}
	if st.CurrentSeat != 1 {
		t.Fatalf("CurrentSeat %v want 1", st.CurrentSeat)
	}
	if st.GamePhase != PhasePlaying || st.TurnPhase != TurnMustDraw {
		t.Fatalf("phase %v/%v want Playing/MustDraw", st.GamePhase, st.TurnPhase)
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation after opening discard: %v", err)
	}
	// Ensure the discarded tile is not still in rack
	for _, tl := range st.Racks[0] {
		if tl.ID == tile.TileInstanceId(tileId) {
			t.Fatalf("discarded tile still in rack")
		}
	}
}

func TestOpeningDiscardOnlyCurrentPlayer(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _ := openingStateWithDeal()
	dispatcher := &mockDispatcher{}
	// Bob tries to discard (should be rejected, only alice is current)
	tileId := string(state.Racks[1][0].ID)
	payload, _ := json.Marshal(map[string]string{"tileId": tileId})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "bob"}, opCode: protocol.OpClientDiscard, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.DiscardRow) != 0 {
		t.Fatalf("bob's discard should be rejected, discard len %d", len(st.DiscardRow))
	}
	if len(st.Racks[1]) != 14 {
		t.Fatalf("bob rack should be unchanged")
	}
	// Check that error was sent (dispatcher BroadcastMessage with OpServerError)
	// Our mockDispatcher doesn't record errors, but we can check that state unchanged is enough for now
}

func TestOpeningDiscardTileMustBeOwned(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _ := openingStateWithDeal()
	dispatcher := &mockDispatcher{}
	// Alice tries to discard bob's tile
	bobTileId := string(state.Racks[1][0].ID)
	payload, _ := json.Marshal(map[string]string{"tileId": bobTileId})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDiscard, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.DiscardRow) != 0 {
		t.Fatalf("foreign tile discard should be rejected")
	}
	// Alice tries to discard non-existent tile
	payload2, _ := json.Marshal(map[string]string{"tileId": "nonexistent"})
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload2})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDiscard, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.DiscardRow) != 0 {
		t.Fatalf("nonexistent tile should be rejected")
	}
}
