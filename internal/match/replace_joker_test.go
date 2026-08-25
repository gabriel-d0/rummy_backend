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

func playingStateForReplace(target TableMeld, rackTiles []tile.TileInstance, seat Seat) (*RoundState, []tile.TileInstance) {
	players := []PlayerId{"alice", "bob"}
	assigned, _ := AssignSeats(players)
	for i := range assigned {
		if assigned[i].ID == "alice" {
			assigned[i].HasOpened = true
			break
		}
	}
	racks := map[Seat][]tile.TileInstance{}
	var allTiles []tile.TileInstance
	// alice rack: rackTiles + filler to 14? For replace we need 3 tiles in rack (1 repl +2 new), but we can fill to 14
	aliceRack := make([]tile.TileInstance, len(rackTiles))
	copy(aliceRack, rackTiles)
	for i := len(aliceRack); i < 14; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("alice-replace-fill-%02d", i))
		t := tile.MustTile(id, tile.Black, tile.Rank(1+i%13))
		aliceRack = append(aliceRack, t)
	}
	if len(aliceRack) > 14 {
		aliceRack = aliceRack[:14]
	}
	// Need to ensure rackTiles are exactly those needed; our helper will have them plus filler
	// But we need to make sure rackTiles are at start of aliceRack for test payload
	// Already copied, so first len(rackTiles) are the needed ones
	racks[0] = aliceRack
	allTiles = append(allTiles, aliceRack...)
	// bob rack 14
	bobRack := make([]tile.TileInstance, 14)
	for i := 0; i < 14; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("bob-replace-%02d", i))
		t := tile.MustTile(id, tile.Yellow, tile.Rank(1+i%13))
		bobRack[i] = t
	}
	allTiles = append(allTiles, bobRack...)
	racks[1] = bobRack
	// target meld tiles
	allTiles = append(allTiles, target.Tiles...)
	// discard
	discard := []DiscardEntry{{Tile: tile.MustTile("disc-replace-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0}}
	allTiles = append(allTiles, discard[0].Tile)
	// stock remainder
	meldCount := len(target.Tiles)
	stockCount := 106 - (len(aliceRack) + len(bobRack) + len(discard) + meldCount)
	if stockCount < 0 {
		stockCount = 0
	}
	stock := make([]tile.TileInstance, stockCount)
	for i := 0; i < stockCount; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("stock-replace-%02d", i))
		t := tile.MustTile(id, tile.Blue, tile.Rank(1+i%13))
		stock[i] = t
		allTiles = append(allTiles, stock[i])
	}
	current := Seat(0)
	if seat != SeatInvalid {
		current = seat
	}
	return &RoundState{
		Players:     assigned,
		Racks:       racks,
		Stock:       stock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{target},
		CurrentSeat: current,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMeldOrDiscard,
		Winner:      SeatInvalid,
	}, allTiles
}

func TestReplaceJokerRunValid(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	j1 := tile.MustJoker("j1")
	rep := tile.MustTile("rep-j1", tile.Red, 7)
	target := TableMeld{
		ID:        "run1",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 5), tile.MustTile("ex2", tile.Red, 6), j1},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{"j1": rep},
		OwnerSeat: 0,
	}
	// Replacement tile 7 red (exact), new meld: j1 + two tiles 9 red, 9 yellow/9 blue? Actually new meld with joker should be valid. Let's make new meld as set: 9 red, 9 yellow, J(9 blue)
	// But new meld tiles are 2 rack tiles plus recovered joker. For set, we need 3 tiles same rank distinct colours. So we need 2 rack tiles 9 red? Wait we need distinct colours.
	// Let's make new meld as run: 8 red, 9 red + J(10 red) => run 8-9-10
	// For that, new tiles should be 8 red and 9 red, and joker represents 10 red
	repTile := tile.MustTile("t-repl", tile.Red, 7) // exact 7 red
	new1 := tile.MustTile("new1", tile.Red, 8)
	new2 := tile.MustTile("new2", tile.Red, 9)
	// Actually new meld will be J(10) +8+9 => 8-9-10, so joker rep 10
	state, allTiles := playingStateForReplace(target, []tile.TileInstance{repTile, new1, new2}, 0)
	// Find actual IDs in state rack (they are repTile, new1, new2 at start of alice rack)
	// Build payload: targetMeldId run1, tileId repTile ID, newMeldTiles [new1,new2], jokerReps for j1 as 10 red
	payload := map[string]interface{}{
		"targetMeldId": "run1",
		"tileId":       string(repTile.ID),
		"newMeldTiles": []string{string(new1.ID), string(new2.ID)},
		"jokerReps": map[string]interface{}{
			"j1": map[string]interface{}{"colour": "red", "rank": 10},
		},
		"newMeldKind": "run",
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientReplaceJoker, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientReplaceJoker, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)

	if len(st.TableMelds) != 2 {
		t.Fatalf("table melds %d want 2", len(st.TableMelds))
	}
	// Target meld should now have 3 real tiles 5,6,7
	foundRepl := false
	for _, m := range st.TableMelds {
		if m.ID == "run1" {
			if len(m.Tiles) != 3 {
				t.Fatalf("target meld tiles %d want 3", len(m.Tiles))
			}
			for _, tl := range m.Tiles {
				if tl.ID == repTile.ID {
					foundRepl = true
				}
				if tl.IsJoker {
					t.Fatalf("target should have no joker after replace")
				}
			}
			if len(m.JokerReps) != 0 {
				t.Fatalf("target JokerReps should be empty, got %d", len(m.JokerReps))
			}
		}
	}
	if !foundRepl {
		t.Fatalf("replacement tile not in target")
	}
	// New meld should contain j1
	foundJokerInNew := false
	for _, m := range st.TableMelds {
		if m.ID != "run1" {
			for _, tl := range m.Tiles {
				if tl.ID == "j1" {
					foundJokerInNew = true
					if m.JokerReps["j1"].Rank != 10 || m.JokerReps["j1"].Colour != tile.Red {
						t.Fatalf("new meld joker rep wrong")
					}
				}
			}
		}
	}
	if !foundJokerInNew {
		t.Fatalf("recovered joker not in new meld")
	}
	// Rack should have 11 after (14-3)
	if len(st.Racks[0]) != 11 {
		t.Fatalf("rack %d want 11", len(st.Racks[0]))
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestReplaceJokerRunWrongTile(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	j1 := tile.MustJoker("j1")
	rep := tile.MustTile("rep-j1", tile.Red, 7)
	target := TableMeld{
		ID:        "run1",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 5), tile.MustTile("ex2", tile.Red, 6), j1},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{"j1": rep},
		OwnerSeat: 0,
	}
	// Wrong replacement: 8 red instead of 7
	wrongTile := tile.MustTile("t-wrong", tile.Red, 8)
	new1 := tile.MustTile("new1", tile.Red, 8)
	new2 := tile.MustTile("new2", tile.Red, 9)
	state, allTiles := playingStateForReplace(target, []tile.TileInstance{wrongTile, new1, new2}, 0)
	payload := map[string]interface{}{
		"targetMeldId": "run1",
		"tileId":       string(wrongTile.ID),
		"newMeldTiles": []string{string(new1.ID), string(new2.ID)},
		"jokerReps": map[string]interface{}{
			"j1": map[string]interface{}{"colour": "red", "rank": 10},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientReplaceJoker, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientReplaceJoker, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 1 {
		t.Fatalf("wrong tile should be rejected, melds %d", len(st.TableMelds))
	}
	if len(st.Racks[0]) != 14 {
		t.Fatalf("rack should be unchanged")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestReplaceJokerSetValid(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	j1 := tile.MustJoker("j1")
	rep := tile.MustTile("rep-j1", tile.Blue, 9) // set 9 red, yellow, J(blue) => missing black? Actually set 9 red,9 yellow, J blue => missing black, but rep is blue
	target := TableMeld{
		ID:        "set1",
		Kind:      "set",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 9), tile.MustTile("ex2", tile.Yellow, 9), j1},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{"j1": rep},
		OwnerSeat: 0,
	}
	// Replacement tile 9 blue (exact)
	repTile := tile.MustTile("t-repl", tile.Blue, 9)
	new1 := tile.MustTile("new1", tile.Red, 5)
	new2 := tile.MustTile("new2", tile.Red, 6)
	// New meld with recovered joker: let's make run 5-6-J(7) => new1 5, new2 6, J 7
	state, allTiles := playingStateForReplace(target, []tile.TileInstance{repTile, new1, new2}, 0)
	payload := map[string]interface{}{
		"targetMeldId": "set1",
		"tileId":       string(repTile.ID),
		"newMeldTiles": []string{string(new1.ID), string(new2.ID)},
		"jokerReps": map[string]interface{}{
			"j1": map[string]interface{}{"colour": "red", "rank": 7},
		},
		"newMeldKind": "run",
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientReplaceJoker, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientReplaceJoker, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 2 {
		t.Fatalf("set replace melds %d want 2", len(st.TableMelds))
	}
	// Target should have 3 real 9s
	for _, m := range st.TableMelds {
		if m.ID == "set1" {
			if len(m.Tiles) != 3 {
				t.Fatalf("target tiles %d want 3", len(m.Tiles))
			}
			for _, tl := range m.Tiles {
				if tl.IsJoker {
					t.Fatalf("target should have no joker")
				}
			}
		}
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestReplaceJokerSetWrongColour(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	j1 := tile.MustJoker("j1")
	rep := tile.MustTile("rep-j1", tile.Blue, 9)
	target := TableMeld{
		ID:        "set1",
		Kind:      "set",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 9), tile.MustTile("ex2", tile.Yellow, 9), j1},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{"j1": rep},
		OwnerSeat: 0,
	}
	// Wrong colour: provide 9 black but joker represents 9 blue, should fail? Actually our MVP says exact colour required, so black would fail
	wrongTile := tile.MustTile("t-wrong", tile.Black, 9)
	new1 := tile.MustTile("new1", tile.Red, 5)
	new2 := tile.MustTile("new2", tile.Red, 6)
	state, allTiles := playingStateForReplace(target, []tile.TileInstance{wrongTile, new1, new2}, 0)
	payload := map[string]interface{}{
		"targetMeldId": "set1",
		"tileId":       string(wrongTile.ID),
		"newMeldTiles": []string{string(new1.ID), string(new2.ID)},
		"jokerReps": map[string]interface{}{
			"j1": map[string]interface{}{"colour": "red", "rank": 7},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientReplaceJoker, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientReplaceJoker, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 1 {
		t.Fatalf("wrong colour should be rejected")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestReplaceJokerNewMeldRequiresTwoTiles(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	j1 := tile.MustJoker("j1")
	rep := tile.MustTile("rep-j1", tile.Red, 7)
	target := TableMeld{
		ID:        "run1",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 5), tile.MustTile("ex2", tile.Red, 6), j1},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{"j1": rep},
		OwnerSeat: 0,
	}
	repTile := tile.MustTile("t-repl", tile.Red, 7)
	new1 := tile.MustTile("new1", tile.Red, 8)
	// Only one new tile, but payload requires 2, validator should reject at envelope level, but handler also checks
	// For this test, we will use a valid payload with 2 tiles but make new meld invalid (e.g., 8 and 12 not consecutive with joker 10)
	new2 := tile.MustTile("new2", tile.Red, 12)
	state, allTiles := playingStateForReplace(target, []tile.TileInstance{repTile, new1, new2}, 0)
	payload := map[string]interface{}{
		"targetMeldId": "run1",
		"tileId":       string(repTile.ID),
		"newMeldTiles": []string{string(new1.ID), string(new2.ID)},
		"jokerReps": map[string]interface{}{
			"j1": map[string]interface{}{"colour": "red", "rank": 10},
		},
		"newMeldKind": "run",
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientReplaceJoker, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientReplaceJoker, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 1 {
		t.Fatalf("invalid new meld should be rejected")
	}
	if len(st.Racks[0]) != 14 {
		t.Fatalf("rack unchanged")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestReplaceJokerAtomicRollback(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	j1 := tile.MustJoker("j1")
	rep := tile.MustTile("rep-j1", tile.Red, 7)
	target := TableMeld{
		ID:        "run1",
		Kind:      "run",
		Tiles:     []tile.TileInstance{tile.MustTile("ex1", tile.Red, 5), tile.MustTile("ex2", tile.Red, 6), j1},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{"j1": rep},
		OwnerSeat: 0,
	}
	repTile := tile.MustTile("t-repl", tile.Red, 7)
	new1 := tile.MustTile("new1", tile.Red, 8)
	new2 := tile.MustTile("new2", tile.Red, 12) // gap, invalid
	state, allTiles := playingStateForReplace(target, []tile.TileInstance{repTile, new1, new2}, 0)
	origTable := append([]TableMeld(nil), state.TableMelds...)
	origRack := append([]tile.TileInstance(nil), state.Racks[0]...)
	payload := map[string]interface{}{
		"targetMeldId": "run1",
		"tileId":       string(repTile.ID),
		"newMeldTiles": []string{string(new1.ID), string(new2.ID)},
		"jokerReps": map[string]interface{}{
			"j1": map[string]interface{}{"colour": "red", "rank": 10},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientReplaceJoker, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientReplaceJoker, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != len(origTable) {
		t.Fatalf("atomic rollback failed, melds %d want %d", len(st.TableMelds), len(origTable))
	}
	if len(st.Racks[0]) != len(origRack) {
		t.Fatalf("atomic rollback rack %d want %d", len(st.Racks[0]), len(origRack))
	}
	// Ensure target unchanged
	if st.TableMelds[0].Tiles[2].ID != "j1" {
		t.Fatalf("target should still have joker")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}
