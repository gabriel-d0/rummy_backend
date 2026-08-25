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

func playingStateForPrevDiscard(opened bool, discardTiles []tile.TileInstance, isOpening []bool, current Seat) (*RoundState, []tile.TileInstance) {
	players := []PlayerId{"alice", "bob"}
	assigned, _ := AssignSeats(players)
	if opened {
		for i := range assigned {
			if assigned[i].ID == "alice" {
				assigned[i].HasOpened = true
				break
			}
		}
	}
	racks := map[Seat][]tile.TileInstance{}
	var allTiles []tile.TileInstance
	// alice rack 14
	aliceRack := make([]tile.TileInstance, 14)
	for i := 0; i < 14; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("alice-prev-%02d", i))
		t := tile.MustTile(id, tile.Black, tile.Rank(1+i%13))
		aliceRack[i] = t
		allTiles = append(allTiles, aliceRack[i])
	}
	racks[0] = aliceRack
	// bob rack 14
	bobRack := make([]tile.TileInstance, 14)
	for i := 0; i < 14; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("bob-prev-%02d", i))
		t := tile.MustTile(id, tile.Yellow, tile.Rank(1+i%13))
		bobRack[i] = t
		allTiles = append(allTiles, bobRack[i])
	}
	racks[1] = bobRack
	// discard row
	discardRow := make([]DiscardEntry, len(discardTiles))
	for i, tl := range discardTiles {
		isOpen := false
		if i < len(isOpening) {
			isOpen = isOpening[i]
		}
		discardRow[i] = DiscardEntry{Tile: tl, IsOpeningDiscard: isOpen, Index: i}
		allTiles = append(allTiles, tl)
	}
	// stock remainder to 106
	meldCount := 0 // no melds for this test
	stockCount := 106 - (len(aliceRack) + len(bobRack) + len(discardRow) + meldCount)
	stock := make([]tile.TileInstance, stockCount)
	for i := 0; i < stockCount; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("stock-prev-%02d", i))
		t := tile.MustTile(id, tile.Blue, tile.Rank(1+i%13))
		stock[i] = t
		allTiles = append(allTiles, stock[i])
	}
	return &RoundState{
		Players:     assigned,
		Racks:       racks,
		Stock:       stock,
		DiscardRow:  discardRow,
		TableMelds:  []TableMeld{},
		CurrentSeat: current,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}, allTiles
}

func TestDrawPreviousDiscardSuccess(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Discard row: opening + 2 normal discards, latest is tile "disc2"
	d1 := tile.MustTile("disc1", tile.Red, 5)
	d2 := tile.MustTile("disc2", tile.Blue, 9)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), d1, d2}
	isOpening := []bool{true, false, false}
	state, allTiles := playingStateForPrevDiscard(true, discardTiles, isOpening, 0)
	origRackLen := len(state.Racks[0])
	origDiscardLen := len(state.DiscardRow)
	origStockLen := len(state.Stock)

	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawPreviousDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawPreviousDiscard, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)

	if len(st.Racks[0]) != origRackLen+1 {
		t.Fatalf("rack %d want %d", len(st.Racks[0]), origRackLen+1)
	}
	if len(st.DiscardRow) != origDiscardLen-1 {
		t.Fatalf("discard %d want %d", len(st.DiscardRow), origDiscardLen-1)
	}
	if len(st.Stock) != origStockLen {
		t.Fatalf("stock should be unchanged, got %d want %d", len(st.Stock), origStockLen)
	}
	// Last tile should be disc2
	lastRackTile := st.Racks[0][len(st.Racks[0])-1]
	if lastRackTile.ID != "disc2" {
		t.Fatalf("drawn tile %v want disc2", lastRackTile.ID)
	}
	if st.TurnPhase != TurnMeldOrDiscard {
		t.Fatalf("TurnPhase %v want MeldOrDiscard", st.TurnPhase)
	}
	if st.CurrentSeat != 0 {
		t.Fatalf("CurrentSeat should stay 0 after draw")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
	// Ensure discard row now ends with d1, not disc2
	if len(st.DiscardRow) > 0 && st.DiscardRow[len(st.DiscardRow)-1].Tile.ID == "disc2" {
		t.Fatalf("disc2 should be removed from discard row")
	}
}

func TestDrawPreviousDiscardUnopenedRejected(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	d1 := tile.MustTile("disc1", tile.Red, 5)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), d1}
	isOpening := []bool{true, false}
	state, allTiles := playingStateForPrevDiscard(false, discardTiles, isOpening, 0)
	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawPreviousDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawPreviousDiscard, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.Racks[0]) != 14 {
		t.Fatalf("unopened rack should be unchanged")
	}
	if len(st.DiscardRow) != 2 {
		t.Fatalf("discard should be unchanged")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestDrawPreviousDiscardOpeningBlocked(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Only opening discard exists
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7)}
	isOpening := []bool{true}
	state, allTiles := playingStateForPrevDiscard(true, discardTiles, isOpening, 0)
	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawPreviousDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawPreviousDiscard, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.Racks[0]) != 14 {
		t.Fatalf("opening blocked should not be picked, rack %d", len(st.Racks[0]))
	}
	if len(st.DiscardRow) != 1 {
		t.Fatalf("discard should remain 1")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestDrawPreviousDiscardOnlyLatest(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Opening + 2 discards, latest is disc2, earlier is disc1
	d1 := tile.MustTile("disc1", tile.Yellow, 3)
	d2 := tile.MustTile("disc2", tile.Blue, 9)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), d1, d2}
	isOpening := []bool{true, false, false}
	state, _ := playingStateForPrevDiscard(true, discardTiles, isOpening, 0)
	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawPreviousDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawPreviousDiscard, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	// Should have taken disc2, not disc1
	if len(st.DiscardRow) != 2 {
		t.Fatalf("discard len %d want 2", len(st.DiscardRow))
	}
	if st.DiscardRow[1].Tile.ID != "disc1" {
		t.Fatalf("remaining discard should be disc1, got %v", st.DiscardRow[1].Tile.ID)
	}
	// Ensure disc1 still there, disc2 moved to rack
	foundDisc1InDiscard := false
	foundDisc2InDiscard := false
	for _, d := range st.DiscardRow {
		if d.Tile.ID == "disc1" {
			foundDisc1InDiscard = true
		}
		if d.Tile.ID == "disc2" {
			foundDisc2InDiscard = true
		}
	}
	if !foundDisc1InDiscard || foundDisc2InDiscard {
		t.Fatalf("discard row should contain disc1 not disc2")
	}
}

func TestDrawPreviousDiscardCannotDrawTwice(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	d1 := tile.MustTile("disc1", tile.Red, 5)
	d2 := tile.MustTile("disc2", tile.Blue, 9)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), d1, d2}
	isOpening := []bool{true, false, false}
	state, _ := playingStateForPrevDiscard(true, discardTiles, isOpening, 0)
	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawPreviousDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawPreviousDiscard, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if st.TurnPhase != TurnMeldOrDiscard {
		t.Fatalf("after first draw should be MeldOrDiscard")
	}
	// Second draw attempt should be rejected (wrong_phase, since now MeldOrDiscard)
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawPreviousDiscard, Payload: payload})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawPreviousDiscard, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.Racks[0]) != len(st.Racks[0]) {
		t.Fatalf("second draw should be rejected, rack %d vs %d", len(st2.Racks[0]), len(st.Racks[0]))
	}
	if len(st2.DiscardRow) != len(st.DiscardRow) {
		t.Fatalf("second draw discard should be unchanged")
	}
}

func TestDrawPreviousDiscardNotYourTurnAndWrongPhase(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	d1 := tile.MustTile("disc1", tile.Red, 5)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), d1}
	isOpening := []bool{true, false}
	state, _ := playingStateForPrevDiscard(true, discardTiles, isOpening, 0)
	payload, _ := json.Marshal(map[string]interface{}{})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawPreviousDiscard, Payload: payload})
	// Bob tries (not current)
	msgBob := &mockMatchData{mockPresence: mockPresence{userId: "bob"}, opCode: protocol.OpClientDrawPreviousDiscard, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msgBob})
	st := next.(*RoundState)
	if len(st.Racks[1]) != 14 {
		t.Fatalf("bob should be rejected")
	}
	// Wrong phase: MeldOrDiscard cannot draw previous
	state2, _ := playingStateForPrevDiscard(true, discardTiles, isOpening, 0)
	state2.TurnPhase = TurnMeldOrDiscard
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawPreviousDiscard, Payload: payload})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawPreviousDiscard, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state2, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.Racks[0]) != 14 {
		t.Fatalf("wrong phase should be rejected")
	}
}
