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

func playingStateForExtend(opened bool, existing TableMeld, newTiles []tile.TileInstance, seat Seat) (*RoundState, []tile.TileInstance) {
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

	// alice rack: newTiles + filler to 15
	aliceRack := make([]tile.TileInstance, len(newTiles))
	copy(aliceRack, newTiles)
	for i := len(aliceRack); i < 15; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("alice-ext-fill-%02d", i))
		t := tile.MustTile(id, tile.Black, tile.Rank(1+i%13))
		aliceRack = append(aliceRack, t)
	}
	if len(aliceRack) > 15 {
		aliceRack = aliceRack[:15]
	}
	racks[0] = aliceRack
	allTiles = append(allTiles, aliceRack...)

	// bob rack: 14
	bobRack := make([]tile.TileInstance, 14)
	for i := 0; i < 14; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("bob-ext-%02d", i))
		t := tile.MustTile(id, tile.Yellow, tile.Rank(1+i%13))
		bobRack[i] = t
	}
	allTiles = append(allTiles, bobRack...)
	racks[1] = bobRack

	// existing meld tiles
	allTiles = append(allTiles, existing.Tiles...)

	// discard
	discard := []DiscardEntry{{Tile: tile.MustTile("disc-ext-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0}}
	allTiles = append(allTiles, discard[0].Tile)

	// stock remainder
	stockCount := 106 - (len(aliceRack) + len(bobRack) + len(discard) + len(existing.Tiles))
	if stockCount < 0 {
		stockCount = 0
	}
	stock := make([]tile.TileInstance, stockCount)
	for i := 0; i < stockCount; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("stock-ext-%02d", i))
		t := tile.MustTile(id, tile.Blue, tile.Rank(1+i%13))
		stock[i] = t
	}
	allTiles = append(allTiles, stock...)

	current := Seat(0)
	if seat != SeatInvalid {
		current = seat
	}
	return &RoundState{
		Players:     assigned,
		Racks:       racks,
		Stock:       stock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{existing},
		CurrentSeat: current,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMeldOrDiscard,
		Winner:      SeatInvalid,
	}, allTiles
}

func TestExtendRunAtEnd(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	existing := TableMeld{
		ID:        "run1",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 5), tile.MustTile("ex2", tile.Red, 6), tile.MustTile("ex3", tile.Red, 7)},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{},
		OwnerSeat: 0,
	}
	newTile := tile.MustTile("t1", tile.Red, 8)
	state, allTiles := playingStateForExtend(true, existing, []tile.TileInstance{newTile}, 0)

	payload := map[string]interface{}{"meldId": "run1", "tileIds": []string{"t1"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)

	if len(st.TableMelds) != 1 {
		t.Fatalf("table melds %d want 1", len(st.TableMelds))
	}
	if len(st.TableMelds[0].Tiles) != 4 {
		t.Fatalf("tiles %d want 4", len(st.TableMelds[0].Tiles))
	}
	if len(st.Racks[0]) != 14 {
		t.Fatalf("rack %d want 14", len(st.Racks[0]))
	}
	// Check new tile present
	found := false
	for _, tl := range st.TableMelds[0].Tiles {
		if tl.ID == "t1" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("new tile not in meld")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
	// Extend at low end: 4-5-6-7 -> add 4
	existing2 := TableMeld{
		ID:        "run2",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Blue, 5), tile.MustTile("ex2", tile.Blue, 6), tile.MustTile("ex3", tile.Blue, 7)},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{},
		OwnerSeat: 0,
	}
	newTile2 := tile.MustTile("t2", tile.Blue, 4)
	state2, allTiles2 := playingStateForExtend(true, existing2, []tile.TileInstance{newTile2}, 0)
	payload2 := map[string]interface{}{"meldId": "run2", "tileIds": []string{"t2"}}
	pBytes2, _ := json.Marshal(payload2)
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes2})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state2, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.TableMelds[0].Tiles) != 4 {
		t.Fatalf("low end extend tiles %d want 4", len(st2.TableMelds[0].Tiles))
	}
	if err := CheckTileConservation(st2, allTiles2); err != nil {
		t.Fatalf("conservation2: %v", err)
	}
}

func TestExtendSetToFourColours(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	existing := TableMeld{
		ID:        "set1",
		Kind:      "set",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 9), tile.MustTile("ex2", tile.Yellow, 9), tile.MustTile("ex3", tile.Blue, 9)},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{},
		OwnerSeat: 0,
	}
	newTile := tile.MustTile("t1", tile.Black, 9)
	state, allTiles := playingStateForExtend(true, existing, []tile.TileInstance{newTile}, 0)
	payload := map[string]interface{}{"meldId": "set1", "tileIds": []string{"t1"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds[0].Tiles) != 4 {
		t.Fatalf("set extend tiles %d want 4", len(st.TableMelds[0].Tiles))
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
	// Try to extend beyond 4 should fail (set max 4)
	extra := tile.MustTile("t2", tile.Red, 9) // duplicate colour, also 5 tiles invalid
	state2, allTiles2 := playingStateForExtend(true, st.TableMelds[0], []tile.TileInstance{extra}, 0)
	payload2 := map[string]interface{}{"meldId": "set1", "tileIds": []string{"t2"}}
	pBytes2, _ := json.Marshal(payload2)
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes2})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state2, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.TableMelds[0].Tiles) != 4 {
		t.Fatalf("extend beyond 4 should be rejected, got %d", len(st2.TableMelds[0].Tiles))
	}
	if err := CheckTileConservation(st2, allTiles2); err != nil {
		t.Fatalf("conservation2: %v", err)
	}
}

func TestExtendAnotherPlayersMeld(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Bob created meld, alice (opened) extends it
	existing := TableMeld{
		ID:        "bob-run",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 10), tile.MustTile("ex2", tile.Red, 11), tile.MustTile("ex3", tile.Red, 12)},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{},
		OwnerSeat: 1,
	}
	newTile := tile.MustTile("t1", tile.Red, 13)
	state, allTiles := playingStateForExtend(true, existing, []tile.TileInstance{newTile}, 0)
	// Ensure alice is current and opened, bob's meld exists
	payload := map[string]interface{}{"meldId": "bob-run", "tileIds": []string{"t1"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds[0].Tiles) != 4 {
		t.Fatalf("extend other player's meld tiles %d want 4", len(st.TableMelds[0].Tiles))
	}
	if st.TableMelds[0].OwnerSeat != 1 {
		t.Fatalf("OwnerSeat should stay 1, got %v", st.TableMelds[0].OwnerSeat)
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestExtendInvalidDoesNotMutate(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	existing := TableMeld{
		ID:        "run1",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 5), tile.MustTile("ex2", tile.Red, 6), tile.MustTile("ex3", tile.Red, 7)},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{},
		OwnerSeat: 0,
	}
	// Invalid: extend with wrong colour (blue 8 instead of red 8)
	newTile := tile.MustTile("t1", tile.Blue, 8)
	state, allTiles := playingStateForExtend(true, existing, []tile.TileInstance{newTile}, 0)
	payload := map[string]interface{}{"meldId": "run1", "tileIds": []string{"t1"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds[0].Tiles) != 3 {
		t.Fatalf("invalid extend should be rejected, got %d", len(st.TableMelds[0].Tiles))
	}
	if len(st.Racks[0]) != 15 {
		t.Fatalf("rack should be unchanged, got %d", len(st.Racks[0]))
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
	// Invalid: extend run with gap (10 instead of 8)
	newTile2 := tile.MustTile("t2", tile.Red, 10)
	state2, allTiles2 := playingStateForExtend(true, existing, []tile.TileInstance{newTile2}, 0)
	payload2 := map[string]interface{}{"meldId": "run1", "tileIds": []string{"t2"}}
	pBytes2, _ := json.Marshal(payload2)
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes2})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state2, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.TableMelds[0].Tiles) != 3 {
		t.Fatalf("gap extend should be rejected")
	}
	if err := CheckTileConservation(st2, allTiles2); err != nil {
		t.Fatalf("conservation2: %v", err)
	}
}

func TestExtendJokerImmutable(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	j1 := tile.MustJoker("j1")
	rep := tile.MustTile("rep-j1", tile.Red, 7)
	existing := TableMeld{
		ID:        "run-joker",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 5), tile.MustTile("ex2", tile.Red, 6), j1},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{"j1": rep},
		OwnerSeat: 0,
	}
	// Valid extend with 8 should keep joker rep 7
	newTile := tile.MustTile("t1", tile.Red, 8)
	state, allTiles := playingStateForExtend(true, existing, []tile.TileInstance{newTile}, 0)
	payload := map[string]interface{}{"meldId": "run-joker", "tileIds": []string{"t1"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds[0].Tiles) != 4 {
		t.Fatalf("extend with joker preserve tiles %d want 4", len(st.TableMelds[0].Tiles))
	}
	if st.TableMelds[0].JokerReps["j1"].Rank != 7 || st.TableMelds[0].JokerReps["j1"].Colour != tile.Red {
		t.Fatalf("joker rep changed")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
	// Try to extend with tile that would imply joker reinterpretation: existing 5-6-J(7), extend with 4 should be valid (4-5-6-7), but joker rep stays 7.
	// Instead test that trying to mutate joker rep via payload is rejected.
	// Client sends jokerRep for existing joker with different rank
	newTile2 := tile.MustTile("t2", tile.Red, 4)
	state2, _ := playingStateForExtend(true, existing, []tile.TileInstance{newTile2}, 0)
	payload2 := map[string]interface{}{
		"meldId":    "run-joker",
		"tileIds":   []string{"t2"},
		"jokerReps": map[string]interface{}{"j1": map[string]interface{}{"colour": "red", "rank": 8}},
	}
	pBytes2, _ := json.Marshal(payload2)
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes2})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state2, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.TableMelds[0].Tiles) != 3 {
		t.Fatalf("joker rep mutation should be rejected, got %d", len(st2.TableMelds[0].Tiles))
	}
}

func TestExtendWithJokerNewTile(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	existing := TableMeld{
		ID:        "run1",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 5), tile.MustTile("ex2", tile.Red, 6), tile.MustTile("ex3", tile.Red, 7)},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{},
		OwnerSeat: 0,
	}
	j1 := tile.MustJoker("j1")
	state, allTiles := playingStateForExtend(true, existing, []tile.TileInstance{j1}, 0)
	payload := map[string]interface{}{
		"meldId":    "run1",
		"tileIds":   []string{"j1"},
		"jokerReps": map[string]interface{}{"j1": map[string]interface{}{"colour": "red", "rank": 8}},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds[0].Tiles) != 4 {
		t.Fatalf("extend with joker tiles %d want 4", len(st.TableMelds[0].Tiles))
	}
	if _, ok := st.TableMelds[0].JokerReps["j1"]; !ok {
		t.Fatalf("new joker rep not stored")
	}
	if st.TableMelds[0].JokerReps["j1"].Rank != 8 {
		t.Fatalf("joker rep rank %v want 8", st.TableMelds[0].JokerReps["j1"].Rank)
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestExtendRejectsUnopenedAndWrongPhase(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	existing := TableMeld{
		ID:        "run1",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 5), tile.MustTile("ex2", tile.Red, 6), tile.MustTile("ex3", tile.Red, 7)},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{},
		OwnerSeat: 0,
	}
	newTile := tile.MustTile("t1", tile.Red, 8)
	// Unopened
	state, _ := playingStateForExtend(false, existing, []tile.TileInstance{newTile}, 0)
	payload := map[string]interface{}{"meldId": "run1", "tileIds": []string{"t1"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds[0].Tiles) != 3 {
		t.Fatalf("unopened extend should be rejected")
	}
	// Wrong phase MustDraw
	state2, _ := playingStateForExtend(true, existing, []tile.TileInstance{newTile}, 0)
	state2.TurnPhase = TurnMustDraw
	pBytes2, _ := json.Marshal(payload)
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes2})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state2, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.TableMelds[0].Tiles) != 3 {
		t.Fatalf("wrong phase extend should be rejected")
	}
	// Not your turn
	state3, _ := playingStateForExtend(true, existing, []tile.TileInstance{newTile}, 1) // bob's turn but alice tries
	state3.CurrentSeat = 1
	pBytes3, _ := json.Marshal(payload)
	envBytes3, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes3})
	msg3 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes3}
	next3 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state3, []runtime.MatchData{msg3})
	st3 := next3.(*RoundState)
	if len(st3.TableMelds[0].Tiles) != 3 {
		t.Fatalf("not your turn should be rejected")
	}
}

func TestExtendTileConservationAndUniqueness(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	existing := TableMeld{
		ID:        "set1",
		Kind:      "set",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 5), tile.MustTile("ex2", tile.Yellow, 5), tile.MustTile("ex3", tile.Blue, 5)},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{},
		OwnerSeat: 0,
	}
	newTile := tile.MustTile("t1", tile.Black, 5)
	state, allTiles := playingStateForExtend(true, existing, []tile.TileInstance{newTile}, 0)
	payload := map[string]interface{}{"meldId": "set1", "tileIds": []string{"t1"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	// Ensure tile not in two melds: new tile should not be in rack and only in meld
	for _, tl := range st.Racks[0] {
		if tl.ID == "t1" {
			t.Fatalf("tile t1 still in rack after extend")
		}
	}
	// Ensure no duplicate across melds: create second meld and try to extend with same tile should fail
	// Already tested duplicate, but check conservation
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
	// Also ensure second extend with same tile fails (tile not in rack)
	payload2 := map[string]interface{}{"meldId": "set1", "tileIds": []string{"t1"}}
	pBytes2, _ := json.Marshal(payload2)
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientExtendMeld, Payload: pBytes2})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientExtendMeld, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.TableMelds[0].Tiles) != 4 {
		t.Fatalf("second extend same tile should be rejected (already used)")
	}
}
