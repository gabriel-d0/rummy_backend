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

func playingStateForDiscard(n int, current Seat) (*RoundState, []tile.TileInstance) {
	// Create a state with Playing MeldOrDiscard, current seat, and racks with 15 each for current, 14 for others
	// Use distinct IDs for conservation check
	players := make([]PlayerId, n)
	for i := 0; i < n; i++ {
		players[i] = PlayerId(fmt.Sprintf("p%d", i))
	}
	assigned, _ := AssignSeats(players)
	racks := map[Seat][]tile.TileInstance{}
	var allTiles []tile.TileInstance
	// Create 15 for current, 14 for others, plus stock and discard
	for seat := 0; seat < n; seat++ {
		count := 14
		if Seat(seat) == current {
			count = 15
		}
		rack := make([]tile.TileInstance, count)
		for i := 0; i < count; i++ {
			id := tile.TileInstanceId(fmt.Sprintf("rack-%d-%02d", seat, i))
			rack[i] = tile.MustTile(id, tile.Red, tile.Rank(1+i%13))
			allTiles = append(allTiles, rack[i])
		}
		racks[Seat(seat)] = rack
	}
	// Stock 10 tiles
	stock := make([]tile.TileInstance, 10)
	for i := 0; i < 10; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("stock-%02d", i))
		stock[i] = tile.MustTile(id, tile.Blue, tile.Rank(1+i%13))
		allTiles = append(allTiles, stock[i])
	}
	// Discard row with opening blocked + maybe one normal
	discard := []DiscardEntry{
		{Tile: tile.MustTile("disc-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0},
		{Tile: tile.MustTile("disc-1", tile.Yellow, 3), IsOpeningDiscard: false, Index: 1},
	}
	allTiles = append(allTiles, discard[0].Tile, discard[1].Tile)
	// Add remaining stock tiles to allTiles already, but need to ensure allTiles includes discard and stock
	// For conservation we need a full 106 deck, but for this test we just use this small allTiles as deck for CheckTileConservation
	// Instead we will not check full 106, just validate counts and that our small deck's conservation holds if we pass it as allTiles
	// For simplicity, we will just return state with small deck and test via counts, not full 106
	return &RoundState{
		Players:     assigned,
		Racks:       racks,
		Stock:       stock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{},
		CurrentSeat: current,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMeldOrDiscard,
		Winner:      SeatInvalid,
	}, allTiles
}

func TestNormalDiscardSuccess(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _ := playingStateForDiscard(2, 0)
	dispatcher := &mockDispatcher{}
	initialRackLen := len(state.Racks[0])
	initialDiscardLen := len(state.DiscardRow)
	tileId := string(state.Racks[0][0].ID)
	payload, _ := json.Marshal(map[string]string{"tileId": tileId})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "p0"}, opCode: protocol.OpClientDiscard, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.Racks[0]) != initialRackLen-1 {
		t.Fatalf("rack len %d want %d", len(st.Racks[0]), initialRackLen-1)
	}
	if len(st.DiscardRow) != initialDiscardLen+1 {
		t.Fatalf("discard len %d want %d", len(st.DiscardRow), initialDiscardLen+1)
	}
	last := st.DiscardRow[len(st.DiscardRow)-1]
	if last.Tile.ID != tile.TileInstanceId(tileId) {
		t.Fatalf("discard tile %v want %v", last.Tile.ID, tileId)
	}
	if last.IsOpeningDiscard {
		t.Fatalf("normal discard should not be opening")
	}
	if last.Index != initialDiscardLen {
		t.Fatalf("discard index %d want %d", last.Index, initialDiscardLen)
	}
	// Current should advance 0→1 for 2p
	if st.CurrentSeat != 1 {
		t.Fatalf("CurrentSeat %v want 1", st.CurrentSeat)
	}
	if st.TurnPhase != TurnMustDraw {
		t.Fatalf("TurnPhase %v want MustDraw", st.TurnPhase)
	}
	// Ensure discarded tile not still in rack
	for _, tl := range st.Racks[0] {
		if tl.ID == tile.TileInstanceId(tileId) {
			t.Fatalf("discarded tile still in rack")
		}
	}
}

func TestNormalDiscardRejectsWrongPlayer(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _ := playingStateForDiscard(2, 0)
	dispatcher := &mockDispatcher{}
	// p1 tries to discard while p0 is current
	tileId := string(state.Racks[1][0].ID)
	payload, _ := json.Marshal(map[string]string{"tileId": tileId})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "p1"}, opCode: protocol.OpClientDiscard, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.DiscardRow) != 2 {
		t.Fatalf("wrong player discard should be rejected, discard len %d", len(st.DiscardRow))
	}
}

func TestNormalDiscardRejectsBeforeDraw(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _ := playingStateForDiscard(2, 0)
	state.TurnPhase = TurnMustDraw // must draw first
	dispatcher := &mockDispatcher{}
	tileId := string(state.Racks[0][0].ID)
	payload, _ := json.Marshal(map[string]string{"tileId": tileId})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "p0"}, opCode: protocol.OpClientDiscard, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.DiscardRow) != 2 {
		t.Fatalf("discard before draw should be rejected (wrong phase)")
	}
}

func TestNormalDiscardTurnOrder(t *testing.T) {
	for _, n := range []int{2, 3, 4} {
		m := &RummyMatch{}
		logger := &testLogger{}
		state, _ := playingStateForDiscard(n, Seat(n-1))
		// Current is last seat, next should be 0
		tileId := string(state.Racks[Seat(n-1)][0].ID)
		payload, _ := json.Marshal(map[string]string{"tileId": tileId})
		envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
		senderId := fmt.Sprintf("p%d", n-1)
		msg := &mockMatchData{mockPresence: mockPresence{userId: senderId}, opCode: protocol.OpClientDiscard, data: envBytes}
		next := m.MatchLoop(context.Background(), logger, nil, nil, nil, 0, state, []runtime.MatchData{msg})
		st := next.(*RoundState)
		if st.CurrentSeat != 0 {
			t.Fatalf("n=%d turn order %d→0 failed, got %v", n, n-1, st.CurrentSeat)
		}
		if st.TurnPhase != TurnMustDraw {
			t.Fatalf("n=%d TurnPhase not MustDraw", n)
		}
		// Check discard row ordering: last discard should be the discarded tile
		last := st.DiscardRow[len(st.DiscardRow)-1]
		if string(last.Tile.ID) != tileId {
			t.Fatalf("n=%d discard order wrong", n)
		}
		if last.Index != 2 { // was 0,1 then new is 2
			t.Fatalf("n=%d discard index %d want 2", n, last.Index)
		}
	}
}
