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

func playingStateForPickup(opened bool, discardTiles []tile.TileInstance, isOpening []bool, rackTiles []tile.TileInstance, current Seat) (*RoundState, []tile.TileInstance) {
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
	// alice rack: rackTiles + filler to 14 (since MustDraw has 14)
	aliceRack := make([]tile.TileInstance, len(rackTiles))
	copy(aliceRack, rackTiles)
	for i := len(aliceRack); i < 14; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("alice-pickup-fill-%02d", i))
		t := tile.MustTile(id, tile.Black, tile.Rank(1+i%13))
		aliceRack = append(aliceRack, t)
	}
	if len(aliceRack) > 14 {
		aliceRack = aliceRack[:14]
	}
	racks[0] = aliceRack
	allTiles = append(allTiles, aliceRack...)
	// bob rack 14
	bobRack := make([]tile.TileInstance, 14)
	for i := 0; i < 14; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("bob-pickup-%02d", i))
		t := tile.MustTile(id, tile.Yellow, tile.Rank(1+i%13))
		bobRack[i] = t
	}
	allTiles = append(allTiles, bobRack...)
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
	// stock remainder
	meldCount := 0
	stockCount := 106 - (len(aliceRack) + len(bobRack) + len(discardRow) + meldCount)
	stock := make([]tile.TileInstance, stockCount)
	for i := 0; i < stockCount; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("stock-pickup-%02d", i))
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

func TestPickupDiscardForMeldValidSet(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Discard row: opening + disc1 (6 red) + disc2 (extra) ; we will pickup disc1 (index 1) with 2 rack tiles to form set
	disc1 := tile.MustTile("disc1", tile.Red, 7)
	disc2 := tile.MustTile("disc2", tile.Blue, 9)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), disc1, disc2}
	isOpening := []bool{true, false, false}
	// Rack has 7 yellow and 7 blue to form set with disc1 (7 red)
	t1 := tile.MustTile("t1", tile.Yellow, 7)
	t2 := tile.MustTile("t2", tile.Blue, 7)
	state, allTiles := playingStateForPickup(true, discardTiles, isOpening, []tile.TileInstance{t1, t2}, 0)
	origRackLen := len(state.Racks[0])
	origDiscardLen := len(state.DiscardRow)
	origMelds := len(state.TableMelds)

	payload := map[string]interface{}{"discardIndex": 1, "tileIds": []string{"t1", "t2"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientPickupDiscardForMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientPickupDiscardForMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)

	if len(st.TableMelds) != origMelds+1 {
		t.Fatalf("melds %d want %d", len(st.TableMelds), origMelds+1)
	}
	// Rack: -2 + later discards (disc2) => orig 14 -2 +1 =13, and DiscardRow becomes 1 (only opening)
	if len(st.Racks[0]) != origRackLen-2+1 {
		t.Fatalf("rack %d want %d", len(st.Racks[0]), origRackLen-2+1)
	}
	if len(st.DiscardRow) != 1 {
		t.Fatalf("discard len %d want 1", len(st.DiscardRow))
	}
	if st.DiscardRow[0].IsOpeningDiscard != true {
		t.Fatalf("remaining discard should be opening")
	}
	// Check discard row no longer contains disc1 or disc2
	for _, d := range st.DiscardRow {
		if d.Tile.ID == "disc1" || d.Tile.ID == "disc2" {
			t.Fatalf("discard row should not contain picked discards")
		}
	}
	// New meld should contain disc1 + t1 + t2
	foundDisc1InMeld := false
	for _, meld := range st.TableMelds {
		for _, tl := range meld.Tiles {
			if tl.ID == "disc1" {
				foundDisc1InMeld = true
			}
		}
	}
	if !foundDisc1InMeld {
		t.Fatalf("disc1 should be in new meld")
	}
	// Later tile disc2 should be in rack
	foundDisc2InRack := false
	for _, tl := range st.Racks[0] {
		if tl.ID == "disc2" {
			foundDisc2InRack = true
		}
	}
	if !foundDisc2InRack {
		t.Fatalf("later discard disc2 should be in rack")
	}
	if st.TurnPhase != TurnMeldOrDiscard {
		t.Fatalf("TurnPhase %v want MeldOrDiscard", st.TurnPhase)
	}
	if st.CurrentSeat != 0 {
		t.Fatalf("CurrentSeat should stay 0")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
	// Ensure original discard not in discard row but in meld, and later swept correctly
	_ = origDiscardLen
}

func TestPickupDiscardForMeldValidRun(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Run: discard is 6 red, rack has 5 red and 7 red => 5-6-7
	disc1 := tile.MustTile("disc1", tile.Red, 6)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), disc1}
	isOpening := []bool{true, false}
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Red, 7)
	state, allTiles := playingStateForPickup(true, discardTiles, isOpening, []tile.TileInstance{t1, t2}, 0)
	payload := map[string]interface{}{"discardIndex": 1, "tileIds": []string{"t1", "t2"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientPickupDiscardForMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientPickupDiscardForMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 1 {
		t.Fatalf("melds %d want 1", len(st.TableMelds))
	}
	// Check meld contains 5-6-7
	meldTiles := st.TableMelds[0].Tiles
	if len(meldTiles) != 3 {
		t.Fatalf("meld tiles %d want 3", len(meldTiles))
	}
	// Ensure meld is run (kind should be run)
	if st.TableMelds[0].Kind != "run" {
		t.Fatalf("kind %q want run", st.TableMelds[0].Kind)
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestPickupDiscardForMeldInvalidRejected(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	disc1 := tile.MustTile("disc1", tile.Red, 7)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), disc1}
	isOpening := []bool{true, false}
	// Invalid: rack tiles 5 red and 6 blue (different colour, not consecutive, not same rank)
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Blue, 6)
	state, allTiles := playingStateForPickup(true, discardTiles, isOpening, []tile.TileInstance{t1, t2}, 0)
	payload := map[string]interface{}{"discardIndex": 1, "tileIds": []string{"t1", "t2"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientPickupDiscardForMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientPickupDiscardForMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 0 {
		t.Fatalf("invalid meld should be rejected, got %d", len(st.TableMelds))
	}
	if len(st.Racks[0]) != 14 {
		t.Fatalf("rack should be unchanged")
	}
	if len(st.DiscardRow) != 2 {
		t.Fatalf("discard should be unchanged")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestPickupDiscardForMeldOpeningBlocked(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	disc1 := tile.MustTile("disc-open", tile.Red, 7) // opening
	discardTiles := []tile.TileInstance{disc1}
	isOpening := []bool{true}
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Red, 6)
	state, allTiles := playingStateForPickup(true, discardTiles, isOpening, []tile.TileInstance{t1, t2}, 0)
	payload := map[string]interface{}{"discardIndex": 0, "tileIds": []string{"t1", "t2"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientPickupDiscardForMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientPickupDiscardForMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 0 {
		t.Fatalf("opening discard should be rejected")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
}

func TestPickupDiscardForMeldLaterSweep(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Discard row: opening, disc1 (6 red to pickup), disc2 (extra1), disc3 (extra2)
	disc1 := tile.MustTile("disc1", tile.Red, 6)
	disc2 := tile.MustTile("disc2", tile.Yellow, 3)
	disc3 := tile.MustTile("disc3", tile.Blue, 9)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), disc1, disc2, disc3}
	isOpening := []bool{true, false, false, false}
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Red, 7)
	state, _ := playingStateForPickup(true, discardTiles, isOpening, []tile.TileInstance{t1, t2}, 0)
	origRackLen := len(state.Racks[0])
	payload := map[string]interface{}{"discardIndex": 1, "tileIds": []string{"t1", "t2"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientPickupDiscardForMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientPickupDiscardForMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	// After pickup index 1, later discards disc2, disc3 (2 tiles) should be in rack
	// New rack = 14 -2 +2 =14, discard becomes 1 (opening)
	if len(st.Racks[0]) != origRackLen {
		t.Fatalf("rack %d want %d (sweep later 2)", len(st.Racks[0]), origRackLen)
	}
	if len(st.DiscardRow) != 1 {
		t.Fatalf("discard len %d want 1", len(st.DiscardRow))
	}
	// Check rack contains disc2 and disc3
	found2, found3 := false, false
	for _, tl := range st.Racks[0] {
		if tl.ID == "disc2" {
			found2 = true
		}
		if tl.ID == "disc3" {
			found3 = true
		}
	}
	if !found2 || !found3 {
		t.Fatalf("later discards should be in rack, found2 %v found3 %v", found2, found3)
	}
	// Ensure order preserved: disc2 before disc3 in rack tail? Our implementation appends in order.
	// Check that disc2 appears before disc3 in rack's tail
	idx2, idx3 := -1, -1
	for i, tl := range st.Racks[0] {
		if tl.ID == "disc2" {
			idx2 = i
		}
		if tl.ID == "disc3" {
			idx3 = i
		}
	}
	if idx2 == -1 || idx3 == -1 || idx2 > idx3 {
		t.Fatalf("later discards order wrong idx2 %d idx3 %d", idx2, idx3)
	}
}

func TestPickupDiscardForMeldAtomic(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	disc1 := tile.MustTile("disc1", tile.Red, 7)
	disc2 := tile.MustTile("disc2", tile.Blue, 9)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), disc1, disc2}
	isOpening := []bool{true, false, false}
	t1 := tile.MustTile("t1", tile.Yellow, 7)
	t2 := tile.MustTile("t2", tile.Blue, 7) // valid set would be disc1(7 red) + t1(7 yellow)+t2(7 blue) = valid, but we will use invalid to test atomic
	// Use invalid: t1 5 red, t2 6 blue not forming meld with disc1 7 red
	t1bad := tile.MustTile("t1bad", tile.Red, 5)
	t2bad := tile.MustTile("t2bad", tile.Blue, 6)
	state, allTiles := playingStateForPickup(true, discardTiles, isOpening, []tile.TileInstance{t1, t2, t1bad, t2bad}, 0)
	// Actually rackTiles includes t1,t2 but we will try invalid with t1bad,t2bad which are also in rack (we need them in rack)
	// Our helper rackTiles is [t1,t2] plus filler, but we need t1bad,t2bad in rack too, so we should include them
	// For this test, recreate state with t1bad,t2bad as the 2 tiles
	stateBad, allTilesBad := playingStateForPickup(true, discardTiles, isOpening, []tile.TileInstance{t1bad, t2bad}, 0)
	payload := map[string]interface{}{"discardIndex": 1, "tileIds": []string{"t1bad", "t2bad"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientPickupDiscardForMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientPickupDiscardForMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, stateBad, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 0 {
		t.Fatalf("invalid should be rejected")
	}
	if len(st.Racks[0]) != 14 {
		t.Fatalf("rack unchanged want 14 got %d", len(st.Racks[0]))
	}
	if len(st.DiscardRow) != 3 {
		t.Fatalf("discard unchanged want 3 got %d", len(st.DiscardRow))
	}
	if err := CheckTileConservation(st, allTilesBad); err != nil {
		t.Fatalf("conservation: %v", err)
	}
	_ = state
	_ = allTiles
	_ = t1
	_ = t2
}

func TestPickupDiscardForMeldRequiresOpenedAndMustDraw(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	disc1 := tile.MustTile("disc1", tile.Red, 7)
	discardTiles := []tile.TileInstance{tile.MustTile("disc-open", tile.Red, 7), disc1}
	isOpening := []bool{true, false}
	t1 := tile.MustTile("t1", tile.Yellow, 7)
	t2 := tile.MustTile("t2", tile.Blue, 7)
	// Not opened
	state, _ := playingStateForPickup(false, discardTiles, isOpening, []tile.TileInstance{t1, t2}, 0)
	payload := map[string]interface{}{"discardIndex": 1, "tileIds": []string{"t1", "t2"}}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientPickupDiscardForMeld, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientPickupDiscardForMeld, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 0 {
		t.Fatalf("not opened should be rejected")
	}
	// Wrong phase MeldOrDiscard
	state2, _ := playingStateForPickup(true, discardTiles, isOpening, []tile.TileInstance{t1, t2}, 0)
	state2.TurnPhase = TurnMeldOrDiscard
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientPickupDiscardForMeld, Payload: pBytes})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientPickupDiscardForMeld, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state2, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.TableMelds) != 0 {
		t.Fatalf("wrong phase should be rejected")
	}
}
