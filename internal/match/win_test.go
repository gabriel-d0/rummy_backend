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

func playingStateForWinDiscard() (*RoundState, []tile.TileInstance) {
	players := []PlayerId{"alice", "bob"}
	assigned, _ := AssignSeats(players)
	assigned[0].HasOpened = true
	assigned[1].HasOpened = true
	// alice has 1 tile left
	aliceTile := tile.MustTile("win-tile", tile.Red, 5)
	aliceRack := []tile.TileInstance{aliceTile}
	// Need to fill to at least 1, but for win we want 1, but state.Validate requires racks not empty? It's okay to have 1
	// But we need allTiles to be 106, so we need to pad stock and other racks to reach 106
	racks := map[Seat][]tile.TileInstance{
		0: aliceRack,
		1: {tile.MustTile("bob1", tile.Yellow, 1), tile.MustTile("bob2", tile.Yellow, 2), tile.MustTile("bob3", tile.Yellow, 3), tile.MustTile("bob4", tile.Yellow, 4), tile.MustTile("bob5", tile.Yellow, 5), tile.MustTile("bob6", tile.Yellow, 6), tile.MustTile("bob7", tile.Yellow, 7), tile.MustTile("bob8", tile.Yellow, 8), tile.MustTile("bob9", tile.Yellow, 9), tile.MustTile("bob10", tile.Yellow, 10), tile.MustTile("bob11", tile.Yellow, 11), tile.MustTile("bob12", tile.Yellow, 12), tile.MustTile("bob13", tile.Yellow, 13), tile.MustTile("bob14", tile.Yellow, 1)},
	}
	var allTiles []tile.TileInstance
	allTiles = append(allTiles, aliceRack...)
	allTiles = append(allTiles, racks[1]...)
	// Add existing melds? None for this test
	// Discard row: opening + maybe one
	discard := []DiscardEntry{{Tile: tile.MustTile("disc-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0}}
	allTiles = append(allTiles, discard[0].Tile)
	// Add aliceTile already, bob racks, discard, need stock to reach 106
	// Count: alice 1 + bob 14 =15, discard 1 =16, need 90 stock
	stock := make([]tile.TileInstance, 90)
	for i := 0; i < 90; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("stock-win-%02d", i))
		col := tile.Colour((i % 4) + 1)
		if col < tile.Red || col > tile.Black {
			col = tile.Blue
		}
		rank := tile.Rank(1 + i%13)
		// Ensure not duplicate with existing IDs
		t := tile.MustTile(id, col, rank)
		stock[i] = t
		allTiles = append(allTiles, t)
	}
	// Need to add remaining tiles to reach 106? Let's compute: allTiles currently 1+14+1+90=106 correct
	return &RoundState{
		Players:     assigned,
		Racks:       racks,
		Stock:       stock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{},
		CurrentSeat: 0,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMeldOrDiscard,
		Winner:      SeatInvalid,
	}, allTiles
}

func TestWinAfterDiscardToZero(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, allTiles := playingStateForWinDiscard()
	// Alice discards her last tile win-tile
	payload, _ := json.Marshal(map[string]string{"tileId": "win-tile"})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDiscard, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)

	if st.GamePhase != PhaseRoundComplete {
		t.Fatalf("GamePhase %v want RoundComplete", st.GamePhase)
	}
	if st.Winner != 0 {
		t.Fatalf("Winner %v want 0", st.Winner)
	}
	if len(st.Racks[0]) != 0 {
		t.Fatalf("alice rack %d want 0", len(st.Racks[0]))
	}
	if len(st.DiscardRow) != 2 {
		t.Fatalf("discard len %d want 2", len(st.DiscardRow))
	}
	if st.DiscardRow[1].Tile.ID != "win-tile" {
		t.Fatalf("discard tile %v want win-tile", st.DiscardRow[1].Tile.ID)
	}
	if st.CurrentSeat != 0 {
		t.Fatalf("CurrentSeat should stay winner 0, got %v", st.CurrentSeat)
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
	// Public view should show winner
	pub := PublicView(st)
	if pub.Winner != 0 {
		t.Fatalf("public winner %v want 0", pub.Winner)
	}
	if pub.GamePhase != "RoundComplete" {
		t.Fatalf("public phase %q want RoundComplete", pub.GamePhase)
	}
}

func TestWinAfterMeldWithoutDiscard(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Alice has 3 tiles forming a valid run 5-6-7, will meld them all to win without discard
	t1 := tile.MustTile("win1", tile.Red, 5)
	t2 := tile.MustTile("win2", tile.Red, 6)
	t3 := tile.MustTile("win3", tile.Red, 7)
	players := []PlayerId{"alice", "bob"}
	assigned, _ := AssignSeats(players)
	assigned[0].HasOpened = true
	racks := map[Seat][]tile.TileInstance{
		0: {t1, t2, t3},
		1: {tile.MustTile("bob1", tile.Yellow, 1), tile.MustTile("bob2", tile.Yellow, 2), tile.MustTile("bob3", tile.Yellow, 3), tile.MustTile("bob4", tile.Yellow, 4), tile.MustTile("bob5", tile.Yellow, 5), tile.MustTile("bob6", tile.Yellow, 6), tile.MustTile("bob7", tile.Yellow, 7), tile.MustTile("bob8", tile.Yellow, 8), tile.MustTile("bob9", tile.Yellow, 9), tile.MustTile("bob10", tile.Yellow, 10), tile.MustTile("bob11", tile.Yellow, 11), tile.MustTile("bob12", tile.Yellow, 12), tile.MustTile("bob13", tile.Yellow, 13), tile.MustTile("bob14", tile.Yellow, 1)},
	}
	var allTiles []tile.TileInstance
	allTiles = append(allTiles, t1, t2, t3)
	allTiles = append(allTiles, racks[1]...)
	discard := []DiscardEntry{{Tile: tile.MustTile("disc-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0}}
	allTiles = append(allTiles, discard[0].Tile)
	stock := make([]tile.TileInstance, 88) // 3+14+1+88=106
	for i := 0; i < 88; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("stock-winmeld-%02d", i))
		t := tile.MustTile(id, tile.Blue, tile.Rank(1+i%13))
		stock[i] = t
		allTiles = append(allTiles, t)
	}
	state := &RoundState{
		Players:     assigned,
		Racks:       racks,
		Stock:       stock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{},
		CurrentSeat: 0,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMeldOrDiscard,
		Winner:      SeatInvalid,
	}
	// Meld the 3 tiles as a run
	payload := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "run", "tileIds": []string{"win1", "win2", "win3"}},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldNew, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldNew, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)

	if st.GamePhase != PhaseRoundComplete {
		t.Fatalf("GamePhase %v want RoundComplete after melding out", st.GamePhase)
	}
	if st.Winner != 0 {
		t.Fatalf("Winner %v want 0", st.Winner)
	}
	if len(st.Racks[0]) != 0 {
		t.Fatalf("rack %d want 0", len(st.Racks[0]))
	}
	if len(st.TableMelds) != 1 {
		t.Fatalf("table melds %d want 1", len(st.TableMelds))
	}
	if len(st.DiscardRow) != 1 {
		t.Fatalf("discard should remain 1 (no discard needed), got %d", len(st.DiscardRow))
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestNoGameplayAfterRoundComplete(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, allTiles := playingStateForWinDiscard()
	// Win via discard
	payload, _ := json.Marshal(map[string]string{"tileId": "win-tile"})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDiscard, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if st.GamePhase != PhaseRoundComplete {
		t.Fatalf("should be RoundComplete")
	}
	// Try to draw stock as alice (should be rejected)
	payload2, _ := json.Marshal(map[string]interface{}{})
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawStock, Payload: payload2})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawStock, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if st2.GamePhase != PhaseRoundComplete {
		t.Fatalf("should stay RoundComplete")
	}
	if len(st2.Racks[0]) != 0 {
		t.Fatalf("rack should stay 0")
	}
	if st2.Winner != 0 {
		t.Fatalf("winner should stay 0")
	}
	// Bob tries to discard
	payload3, _ := json.Marshal(map[string]string{"tileId": "bob1"})
	envBytes3, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload3})
	msg3 := &mockMatchData{mockPresence: mockPresence{userId: "bob"}, opCode: protocol.OpClientDiscard, data: envBytes3}
	next3 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st2, []runtime.MatchData{msg3})
	st3 := next3.(*RoundState)
	if st3.GamePhase != PhaseRoundComplete {
		t.Fatalf("bob discard should be rejected, still RoundComplete")
	}
	if err := CheckTileConservation(st3, allTiles); err != nil {
		t.Fatalf("conservation after rejected: %v", err)
	}
	// Ensure no duplicate winner change
	if st3.Winner != 0 {
		t.Fatalf("winner should remain 0")
	}
}

func TestWinnerCorrectlyRecorded(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Bob will win: set current to bob with 1 tile
	players := []PlayerId{"alice", "bob"}
	assigned, _ := AssignSeats(players)
	assigned[0].HasOpened = true
	assigned[1].HasOpened = true
	winTile := tile.MustTile("bob-win", tile.Blue, 9)
	racks := map[Seat][]tile.TileInstance{
		0: {tile.MustTile("alice1", tile.Red, 1), tile.MustTile("alice2", tile.Red, 2)},
		1: {winTile},
	}
	var allTiles []tile.TileInstance
	allTiles = append(allTiles, racks[0]...)
	allTiles = append(allTiles, racks[1]...)
	discard := []DiscardEntry{{Tile: tile.MustTile("disc-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0}}
	allTiles = append(allTiles, discard[0].Tile)
	stock := make([]tile.TileInstance, 89) // 2+1+1+? => 2+1=3, +1 discard=4, need 102 stock? Wait 106-3-1=102? Actually racks 2+1=3, discard 1=4, need 102 stock
	// Let's compute correctly: we have alice 2, bob 1 =3, discard 1=4, need 102 stock to reach 106
	stock = make([]tile.TileInstance, 102)
	for i := 0; i < 102; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("stock-bobwin-%02d", i))
		t := tile.MustTile(id, tile.Black, tile.Rank(1+i%13))
		stock[i] = t
		allTiles = append(allTiles, t)
	}
	state := &RoundState{
		Players:     assigned,
		Racks:       racks,
		Stock:       stock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{},
		CurrentSeat: 1,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMeldOrDiscard,
		Winner:      SeatInvalid,
	}
	payload, _ := json.Marshal(map[string]string{"tileId": "bob-win"})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "bob"}, opCode: protocol.OpClientDiscard, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if st.Winner != 1 {
		t.Fatalf("Winner %v want 1 (bob)", st.Winner)
	}
	if st.GamePhase != PhaseRoundComplete {
		t.Fatalf("phase %v want RoundComplete", st.GamePhase)
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}
