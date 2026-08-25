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

// fullDeck returns 106 distinct valid tiles for conservation checks.
func fullDeck() []tile.TileInstance {
	var deck []tile.TileInstance
	for i := 0; i < 104; i++ {
		colour := []tile.Colour{tile.Red, tile.Yellow, tile.Blue, tile.Black}[i%4]
		rank := tile.Rank(1 + (i % 13))
		id := tile.TileInstanceId(fmt.Sprintf("fd-%03d", i))
		deck = append(deck, tile.MustTile(id, colour, rank))
	}
	deck = append(deck, tile.MustJoker("fd-j1"), tile.MustJoker("fd-j2"))
	return deck
}

// stateForLoop creates a Playing MustDraw state with n players, current seat, and full 106 conservation.
func stateForLoop(n int, current Seat) (*RoundState, []tile.TileInstance) {
	deck := fullDeck()
	// Distribute: need to create a valid RoundState with conservation.
	// For simplicity, use same distribution as NewRoundState: 15 for seat 0, 14 for others, stock remainder.
	// But for this test we want current to be variable, so we will manually slice deck.
	// Use deck[0:15] for seat0, deck[15:29] for seat1 etc., stock remainder, and a single opening discard.
	players := make([]PlayerId, n)
	for i := 0; i < n; i++ {
		players[i] = PlayerId(fmt.Sprintf("p%d", i))
	}
	assigned, _ := AssignSeats(players)
	racks := map[Seat][]tile.TileInstance{}
	offset := 0
	for seat := 0; seat < n; seat++ {
		cnt := 14
		if seat == 0 {
			cnt = 15
		}
		racks[Seat(seat)] = append([]tile.TileInstance(nil), deck[offset:offset+cnt]...)
		offset += cnt
	}
	// Opening discard is deck[offset] as blocked, then stock is remainder
	openTile := deck[offset]
	discard := []DiscardEntry{{Tile: openTile, IsOpeningDiscard: true, Index: 0}}
	offset++
	stock := append([]tile.TileInstance(nil), deck[offset:]...)
	// After opening discard, the opener's rack should be 14 (remove openTile from rack0)
	// But our racks already include openTile as part of rack0's 15; we need to simulate that opening discard already happened:
	// So rack0 should be 14 after discarding openTile, and discard row has it.
	// Remove openTile from rack0
	newRack0 := []tile.TileInstance{}
	for _, t := range racks[0] {
		if t.ID != openTile.ID {
			newRack0 = append(newRack0, t)
		}
	}
	// If openTile was not in rack0 (because we sliced sequentially, openTile is deck[15] for n=2? Actually for n=2, rack0 15 is deck[0:15], rack1 is deck[15:29], openTile is deck[29], which is not in rack0)
	// For this helper we want to simulate post-opening state: current is next after opener, so opener's rack should be 14, with one tile moved to discard.
	// Instead, we should create a state where opener's rack already has 14 after discard, and discard has openTile.
	// Our current racks[0] is 15, but we just removed openTile if it was there; if not, we need to adjust.
	// For simplicity, ensure rack0 has 14 by removing its first tile and using that as openTile
	if len(newRack0) == 15 {
		// openTile was in rack0, now newRack0 is 14 correct
		racks[0] = newRack0
	} else {
		// openTile was not in rack0 (because offset after racks), so we need to make rack0 14 by removing one tile and using that as openTile
		// Use rack0[0] as open and put deck[offset] back to stock
		openTile = racks[0][0]
		racks[0] = racks[0][1:]
		discard[0].Tile = openTile
		// Put the previous openTile (deck[offset]) back to stock top
		stock = append([]tile.TileInstance{deck[offset]}, stock...)
	}
	allTiles := deck
	return &RoundState{
		Players:     assigned,
		Racks:       racks,
		Stock:       stock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{},
		CurrentSeat: current,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}, allTiles
}

func TestTurnLoopOpeningDrawDiscardNextPlayer(t *testing.T) {
	for _, n := range []int{2, 3, 4} {
		m := &RummyMatch{}
		logger := &testLogger{}
		dispatcher := &mockDispatcher{}
		// Start with current 1 for n=2, or 1 for n=3, etc. We'll test a full loop: current draws, discards, next becomes current
		current := Seat(1 % n)
		state, allTiles := stateForLoop(n, current)
		if err := state.Validate(); err != nil {
			t.Fatalf("n=%d validate %v", n, err)
		}
		if err := CheckTileConservation(state, allTiles); err != nil {
			t.Fatalf("n=%d conservation before %v", n, err)
		}
		// Current draws
		initialStock := len(state.Stock)
		initialRack := len(state.Racks[current])
		drawPayload, _ := json.Marshal(map[string]interface{}{})
		drawEnv, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawStock, Payload: drawPayload})
		drawMsg := &mockMatchData{mockPresence: mockPresence{userId: string(state.Players[current].ID)}, opCode: protocol.OpClientDrawStock, data: drawEnv}
		next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{drawMsg})
		st := next.(*RoundState)
		if len(st.Stock) != initialStock-1 {
			t.Fatalf("n=%d after draw stock %d want %d", n, len(st.Stock), initialStock-1)
		}
		if len(st.Racks[current]) != initialRack+1 {
			t.Fatalf("n=%d after draw rack %d want %d", n, len(st.Racks[current]), initialRack+1)
		}
		if st.TurnPhase != TurnMeldOrDiscard {
			t.Fatalf("n=%d after draw TurnPhase %v", n, st.TurnPhase)
		}
		if err := CheckTileConservation(st, allTiles); err != nil {
			t.Fatalf("n=%d after draw conservation %v", n, err)
		}
		// Now discard
		tileId := string(st.Racks[current][0].ID)
		discardPayload, _ := json.Marshal(map[string]string{"tileId": tileId})
		discardEnv, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: discardPayload})
		discardMsg := &mockMatchData{mockPresence: mockPresence{userId: string(st.Players[current].ID)}, opCode: protocol.OpClientDiscard, data: discardEnv}
		next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 1, st, []runtime.MatchData{discardMsg})
		st2 := next2.(*RoundState)
		if len(st2.DiscardRow) != 2 {
			t.Fatalf("n=%d after discard discard len %d want 2", n, len(st2.DiscardRow))
		}
		if st2.DiscardRow[1].Index != 1 || st2.DiscardRow[1].IsOpeningDiscard {
			t.Fatalf("n=%d second discard %+v", n, st2.DiscardRow[1])
		}
		expectedNext, _ := NextSeat(current, n)
		if st2.CurrentSeat != expectedNext {
			t.Fatalf("n=%d after discard CurrentSeat %v want %v", n, st2.CurrentSeat, expectedNext)
		}
		if st2.TurnPhase != TurnMustDraw {
			t.Fatalf("n=%d after discard TurnPhase %v want MustDraw", n, st2.TurnPhase)
		}
		if err := CheckTileConservation(st2, allTiles); err != nil {
			t.Fatalf("n=%d after discard conservation %v", n, err)
		}
		// Ensure turn order is anticlockwise
		if expectedNext != Seat((int(current)+1)%n) {
			t.Fatalf("n=%d anticlockwise failed", n)
		}
	}
}
