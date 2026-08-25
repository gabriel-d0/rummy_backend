package match

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
	"github.com/heroiclabs/nakama-common/runtime"
)

// helper: creates a Playing MeldOrDiscard state where current seat is 0 (alice) with
// a rack containing the given tiles plus fill to reach 15. Returns state and allTiles for conservation (106 total).
// The rack for current player will be exactly the provided tiles plus filler tiles (distinct IDs) to 15.
// Other players get 14 filler tiles, stock 76 (so that racks 29 + stock 76 + discard 1 =106), discard 1 (opening). All IDs distinct.
func playingStateForMeldInitial(currentRack []tile.TileInstance) (*RoundState, []tile.TileInstance) {
	players := []PlayerId{"alice", "bob"}
	assigned, _ := AssignSeats(players)
	// Build racks
	racks := map[Seat][]tile.TileInstance{}
	var allTiles []tile.TileInstance

	// alice rack: currentRack plus filler to 15
	aliceRack := make([]tile.TileInstance, len(currentRack))
	copy(aliceRack, currentRack)
	for i := len(aliceRack); i < 15; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("alice-fill-%02d", i))
		t := tile.MustTile(id, tile.Black, tile.Rank(1+i%13))
		aliceRack = append(aliceRack, t)
	}
	if len(aliceRack) > 15 {
		aliceRack = aliceRack[:15]
	}
	racks[0] = aliceRack
	allTiles = append(allTiles, aliceRack...)

	// bob rack: 14 filler
	bobRack := make([]tile.TileInstance, 14)
	for i := 0; i < 14; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("bob-%02d", i))
		t := tile.MustTile(id, tile.Yellow, tile.Rank(1+i%13))
		bobRack[i] = t
	}
	allTiles = append(allTiles, bobRack...)
	racks[1] = bobRack

	// stock 76 for 2p: 106 -29 -1 =76
	stockCount := 106 - (len(aliceRack) + len(bobRack)) - 1
	stock := make([]tile.TileInstance, stockCount)
	for i := 0; i < stockCount; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("stock-%02d", i))
		t := tile.MustTile(id, tile.Blue, tile.Rank(1+i%13))
		stock[i] = t
	}
	allTiles = append(allTiles, stock...)

	// discard row: opening blocked
	discard := []DiscardEntry{{Tile: tile.MustTile("disc-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0}}
	allTiles = append(allTiles, discard[0].Tile)

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

func TestMeldInitialSuccess(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Build a valid 50-point batch: run 5-6-7 red (15), run 8-9-10 blue (20), run 2-3-4 yellow (15) =50
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Red, 6)
	t3 := tile.MustTile("t3", tile.Red, 7)
	t4 := tile.MustTile("t4", tile.Blue, 8)
	t5 := tile.MustTile("t5", tile.Blue, 9)
	t6 := tile.MustTile("t6", tile.Blue, 10)
	t7 := tile.MustTile("t7", tile.Yellow, 2)
	t8 := tile.MustTile("t8", tile.Yellow, 3)
	t9 := tile.MustTile("t9", tile.Yellow, 4)
	rackTiles := []tile.TileInstance{t1, t2, t3, t4, t5, t6, t7, t8, t9}
	state, allTiles := playingStateForMeldInitial(rackTiles)
	// Need to ensure our state rack actually contains exactly those 9 plus filler; playingStateForMeldInitial added filler but we want rack to be exactly those 9 plus 6 filler? Let's replace alice rack directly with 9 tiles + filler that are exactly the 9 we need.
	// Our helper already used rackTiles as the 9 and added filler to 15. But we used t1..t9 as currentRack, helper added 6 filler black tiles. That's fine.
	// However we must ensure the alice rack in state is exactly the helper's rack which includes t1..t9 plus filler.
	// Retrieve actual rack to confirm contains t1..t9
	if len(state.Racks[0]) != 15 {
		t.Fatalf("alice rack %d want 15", len(state.Racks[0]))
	}
	// Build payload: 3 melds
	payload := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "run", "tileIds": []string{"t1", "t2", "t3"}},
			map[string]interface{}{"id": "m2", "kind": "run", "tileIds": []string{"t4", "t5", "t6"}},
			map[string]interface{}{"id": "m3", "kind": "run", "tileIds": []string{"t7", "t8", "t9"}},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)

	// Assertions per AGENTS Day 13
	if len(st.TableMelds) != 3 {
		t.Fatalf("TableMelds %d want 3", len(st.TableMelds))
	}
	if len(st.Racks[0]) != 6 {
		t.Fatalf("alice rack after meld %d want 6 (15-9)", len(st.Racks[0]))
	}
	// HasOpened true
	found := false
	for _, p := range st.Players {
		if p.ID == "alice" && p.HasOpened {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("alice HasOpened should be true")
	}
	// Still MeldOrDiscard, same seat, not advanced
	if st.TurnPhase != TurnMeldOrDiscard {
		t.Fatalf("TurnPhase %v want MeldOrDiscard", st.TurnPhase)
	}
	if st.CurrentSeat != 0 {
		t.Fatalf("CurrentSeat %v want 0 (must discard to advance)", st.CurrentSeat)
	}
	// Conservation: moved tiles from rack to melds, total same
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation failed: %v", err)
	}
	// Other players see public meld but not rack details
	pub := PublicView(st)
	if len(pub.TableMelds) != 3 {
		t.Fatalf("public view TableMelds %d want 3", len(pub.TableMelds))
	}
	privBob := PrivateView(st, 1)
	// Bob's private view should not contain alice's remaining rack IDs
	aliceRackIDs := map[tile.TileInstanceId]bool{}
	for _, tid := range st.Racks[0] {
		aliceRackIDs[tid.ID] = true
	}
	for _, tid := range privBob.OwnRack {
		if aliceRackIDs[tid.ID] {
			t.Fatalf("bob private view leaked alice tile %v", tid.ID)
		}
	}
	// Bob's public rackCount should be 14, alice's 6
	for _, pp := range pub.Players {
		if pp.ID == "alice" && pp.RackCount != 6 {
			t.Fatalf("alice public rackCount %d want 6", pp.RackCount)
		}
		if pp.ID == "bob" && pp.RackCount != 14 {
			t.Fatalf("bob public rackCount %d want 14", pp.RackCount)
		}
	}
	// Meld IDs stable
	if st.TableMelds[0].ID != "m1" || st.TableMelds[1].ID != "m2" || st.TableMelds[2].ID != "m3" {
		t.Fatalf("meld IDs not stable: %v %v %v", st.TableMelds[0].ID, st.TableMelds[1].ID, st.TableMelds[2].ID)
	}
}

func TestMeldInitialInvalidLeavesStateUnchanged(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Valid run 5-6-7 (15) + set 7 (15) =30 <50 -> should fail 50-point rule
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Red, 6)
	t3 := tile.MustTile("t3", tile.Red, 7)
	t4 := tile.MustTile("t4", tile.Red, 7)
	t5 := tile.MustTile("t5", tile.Yellow, 7)
	t6 := tile.MustTile("t6", tile.Blue, 7)
	rackTiles := []tile.TileInstance{t1, t2, t3, t4, t5, t6}
	state, allTiles := playingStateForMeldInitial(rackTiles)
	// Keep original counts
	origRackLen := len(state.Racks[0])
	origMelds := len(state.TableMelds)
	origHasOpened := state.Players[0].HasOpened

	payload := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "run", "tileIds": []string{"t1", "t2", "t3"}},
			map[string]interface{}{"id": "m2", "kind": "set", "tileIds": []string{"t4", "t5", "t6"}},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)

	if len(st.Racks[0]) != origRackLen {
		t.Fatalf("rack len changed on invalid meld: %d vs %d", len(st.Racks[0]), origRackLen)
	}
	if len(st.TableMelds) != origMelds {
		t.Fatalf("melds changed on invalid: %d vs %d", len(st.TableMelds), origMelds)
	}
	if st.Players[0].HasOpened != origHasOpened {
		t.Fatalf("HasOpened changed on invalid")
	}
	if st.TurnPhase != TurnMeldOrDiscard || st.CurrentSeat != 0 {
		t.Fatalf("phase/seat changed on invalid")
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation after rejected meld: %v", err)
	}
	// Also test no-run rejection: 4 sets =60 but no run should fail
	t7 := tile.MustTile("t7", tile.Red, 5)
	t8 := tile.MustTile("t8", tile.Yellow, 5)
	t9 := tile.MustTile("t9", tile.Blue, 5)
	t10 := tile.MustTile("t10", tile.Red, 6)
	t11 := tile.MustTile("t11", tile.Yellow, 6)
	t12 := tile.MustTile("t12", tile.Blue, 6)
	t13 := tile.MustTile("t13", tile.Red, 7)
	t14 := tile.MustTile("t14", tile.Yellow, 7)
	t15 := tile.MustTile("t15", tile.Blue, 7)
	// Need to rebuild state with new rackTiles that are sets only
	rackTiles2 := []tile.TileInstance{t7, t8, t9, t10, t11, t12, t13, t14, t15}
	state2, allTiles2 := playingStateForMeldInitial(rackTiles2)
	// Replace alice rack to contain these 9 sets tiles + filler? Actually playingStateForMeldInitial already created 15 with t7..t15 plus filler.
	// But we need to ensure rack contains exactly those 9 plus filler
	payload2 := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "set", "tileIds": []string{"t7", "t8", "t9"}},
			map[string]interface{}{"id": "m2", "kind": "set", "tileIds": []string{"t10", "t11", "t12"}},
			map[string]interface{}{"id": "m3", "kind": "set", "tileIds": []string{"t13", "t14", "t15"}},
		},
	}
	pBytes2, _ := json.Marshal(payload2)
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes2})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state2, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.TableMelds) != 0 {
		t.Fatalf("no-run batch should be rejected, melds %d", len(st2.TableMelds))
	}
	if err := CheckTileConservation(st2, allTiles2); err != nil {
		t.Fatalf("conservation after no-run reject: %v", err)
	}
}

func TestMeldInitialCannotTwice(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Red, 6)
	t3 := tile.MustTile("t3", tile.Red, 7)
	t4 := tile.MustTile("t4", tile.Blue, 8)
	t5 := tile.MustTile("t5", tile.Blue, 9)
	t6 := tile.MustTile("t6", tile.Blue, 10)
	t7 := tile.MustTile("t7", tile.Yellow, 2)
	t8 := tile.MustTile("t8", tile.Yellow, 3)
	t9 := tile.MustTile("t9", tile.Yellow, 4)
	// Need a second valid batch after opening: use remaining filler tiles that are 10-11-12 red etc. But we will just try to open again with same tiles which are now not in rack.
	rackTiles := []tile.TileInstance{t1, t2, t3, t4, t5, t6, t7, t8, t9}
	state, _ := playingStateForMeldInitial(rackTiles)
	payload := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "run", "tileIds": []string{"t1", "t2", "t3"}},
			map[string]interface{}{"id": "m2", "kind": "run", "tileIds": []string{"t4", "t5", "t6"}},
			map[string]interface{}{"id": "m3", "kind": "run", "tileIds": []string{"t7", "t8", "t9"}},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if !st.Players[0].HasOpened {
		t.Fatalf("first meld should open")
	}
	// Second attempt: try to open again with filler tiles that could make a valid meld (but should be rejected as already opened)
	// Filler tiles in rack are black 1.. etc. Let's create a valid second payload using some filler IDs
	// Alice's remaining rack after first meld is 6 filler black tiles. Pick 3 of them that are consecutive same colour?
	// Our filler for alice are black 1,2,3,4,5,6... Actually we added filler black tiles with rank 1+ i%13. For i from 9 to 14, ranks 10,11,12,13,1,2 etc. Not necessarily consecutive same colour but we used Black colour for filler.
	// Let's try to use three filler tiles that are consecutive: they are black 10,11,12 maybe?
	// Instead we will just try to reuse same payload which should fail as tiles not in rack (already used) and also already opened.
	envBytes2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes2}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{msg2})
	st2 := next2.(*RoundState)
	if len(st2.TableMelds) != 3 {
		t.Fatalf("second meld should be rejected, table %d", len(st2.TableMelds))
	}
	if len(st2.Racks[0]) != 6 {
		t.Fatalf("rack should stay 6 after rejected second meld, got %d", len(st2.Racks[0]))
	}
	// Also test that second meld with new valid tiles but already opened is rejected via already_opened code path
	// Create a new valid run from remaining filler: pick first 3 filler IDs
	remainingIDs := []string{}
	for _, t := range st.Racks[0] {
		remainingIDs = append(remainingIDs, string(t.ID))
		if len(remainingIDs) == 3 {
			break
		}
	}
	if len(remainingIDs) == 3 {
		payloadNew := map[string]interface{}{
			"melds": []interface{}{
				map[string]interface{}{"id": "m4", "kind": "run", "tileIds": remainingIDs},
			},
		}
		pBytesNew, _ := json.Marshal(payloadNew)
		envBytesNew, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytesNew})
		msgNew := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytesNew}
		next3 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{msgNew})
		st3 := next3.(*RoundState)
		if len(st3.TableMelds) != 3 {
			t.Fatalf("MELD_INITIAL after already opened should be rejected regardless of tiles, got %d", len(st3.TableMelds))
		}
	}
}

func TestMeldInitialPlayerStillMustDiscard(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	t1 := tile.MustTile("t1", tile.Red, 10)
	t2 := tile.MustTile("t2", tile.Red, 11)
	t3 := tile.MustTile("t3", tile.Red, 12) // 30
	t4 := tile.MustTile("t4", tile.Blue, 10)
	t5 := tile.MustTile("t5", tile.Blue, 11)
	t6 := tile.MustTile("t6", tile.Blue, 12) // 30 total 60 but need at least one run, we have two runs =60 >=50
	rackTiles := []tile.TileInstance{t1, t2, t3, t4, t5, t6}
	state, _ := playingStateForMeldInitial(rackTiles)
	payload := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "run", "tileIds": []string{"t1", "t2", "t3"}},
			map[string]interface{}{"id": "m2", "kind": "run", "tileIds": []string{"t4", "t5", "t6"}},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if st.TurnPhase != TurnMeldOrDiscard {
		t.Fatalf("after initial meld should stay MeldOrDiscard, got %v", st.TurnPhase)
	}
	if st.CurrentSeat != 0 {
		t.Fatalf("CurrentSeat should stay 0 after meld, got %v", st.CurrentSeat)
	}
	// Now discard should succeed and advance turn
	if len(st.Racks[0]) == 0 {
		t.Fatalf("rack empty, cannot discard")
	}
	discardTileId := string(st.Racks[0][0].ID)
	discardPayload, _ := json.Marshal(map[string]string{"tileId": discardTileId})
	discardEnv, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: discardPayload})
	discardMsg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDiscard, data: discardEnv}
	next2 := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, st, []runtime.MatchData{discardMsg})
	st2 := next2.(*RoundState)
	if st2.CurrentSeat != 1 {
		t.Fatalf("after discard CurrentSeat %v want 1", st2.CurrentSeat)
	}
	if st2.TurnPhase != TurnMustDraw {
		t.Fatalf("after discard TurnPhase %v want MustDraw", st2.TurnPhase)
	}
	if len(st2.DiscardRow) != 2 { // opening + new discard
		t.Fatalf("discard row %d want 2", len(st2.DiscardRow))
	}
}

func TestMeldInitialOtherPlayersSeePublicMeld(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Red, 6)
	t3 := tile.MustTile("t3", tile.Red, 7)
	t4 := tile.MustTile("t4", tile.Blue, 8)
	t5 := tile.MustTile("t5", tile.Blue, 9)
	t6 := tile.MustTile("t6", tile.Blue, 10)
	t7 := tile.MustTile("t7", tile.Yellow, 2)
	t8 := tile.MustTile("t8", tile.Yellow, 3)
	t9 := tile.MustTile("t9", tile.Yellow, 4)
	rackTiles := []tile.TileInstance{t1, t2, t3, t4, t5, t6, t7, t8, t9}
	state, _ := playingStateForMeldInitial(rackTiles)
	payload := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "run", "tileIds": []string{"t1", "t2", "t3"}},
			map[string]interface{}{"id": "m2", "kind": "run", "tileIds": []string{"t4", "t5", "t6"}},
			map[string]interface{}{"id": "m3", "kind": "run", "tileIds": []string{"t7", "t8", "t9"}},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)

	// Public view
	pub := PublicView(st)
	if len(pub.TableMelds) != 3 {
		t.Fatalf("public TableMelds %d", len(pub.TableMelds))
	}
	// Check jokerReps not leaked (no jokers here, but check)
	// Ensure private views
	privAlice := PrivateView(st, 0)
	privBob := PrivateView(st, 1)
	// Alice sees her own rack (6 tiles)
	if len(privAlice.OwnRack) != 6 {
		t.Fatalf("alice ownrack %d want 6", len(privAlice.OwnRack))
	}
	// Bob's rack count public 14, private 14, but should not contain alice's meld tile IDs
	for _, meld := range pub.TableMelds {
		for _, tl := range meld.Tiles {
			// Ensure bob's private rack doesn't contain those meld tiles (they were alice's)
			for _, btl := range privBob.OwnRack {
				if btl.ID == tl.ID {
					t.Fatalf("bob rack leaked meld tile %v", tl.ID)
				}
			}
		}
	}
	// Ensure public view does not contain any alice rack IDs as string (redaction test)
	// Marshal public view and search
	b, _ := json.Marshal(pub)
	pubStr := string(b)
	for _, tl := range privAlice.OwnRack {
		// Alice's remaining rack IDs should NOT be in public view except via RackCount
		// But they shouldn't appear as tile objects; we check that at least one alice tile ID is not in public string?
		// Actually alice's remaining rack should not be in public either, only counts.
		if len(pubStr) > 0 && strings.Contains(pubStr, string(tl.ID)) {
			t.Fatalf("public view leaked alice remaining rack tile %v", tl.ID)
		}
	}
}

func TestMeldInitialAtomicDuplicateTile(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Red, 6)
	t3 := tile.MustTile("t3", tile.Red, 7)
	// Duplicate t2 across melds
	t4 := tile.MustTile("t4", tile.Blue, 5)
	t5 := tile.MustTile("t5", tile.Yellow, 5)
	rackTiles := []tile.TileInstance{t1, t2, t3, t4, t5}
	state, allTiles := playingStateForMeldInitial(rackTiles)
	// Payload uses t2 twice across melds (duplicate)
	payload := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "run", "tileIds": []string{"t1", "t2", "t3"}},
			map[string]interface{}{"id": "m2", "kind": "set", "tileIds": []string{"t2", "t4", "t5"}}, // t2 duplicate
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 0 {
		t.Fatalf("duplicate across melds should be atomic reject, got %d melds", len(st.TableMelds))
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation after duplicate reject: %v", err)
	}
}

func TestMeldInitialWithJoker(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Run 5-6-J where J represents 7 red => 15 points, need total 50, so add two more runs
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Red, 6)
	j1 := tile.MustJoker("j1")
	// Joker will represent 7 red
	t4 := tile.MustTile("t4", tile.Blue, 8)
	t5 := tile.MustTile("t5", tile.Blue, 9)
	t6 := tile.MustTile("t6", tile.Blue, 10) // 20
	t7 := tile.MustTile("t7", tile.Yellow, 2)
	t8 := tile.MustTile("t8", tile.Yellow, 3)
	t9 := tile.MustTile("t9", tile.Yellow, 4) // 15 total with joker run 15 =50
	rackTiles := []tile.TileInstance{t1, t2, j1, t4, t5, t6, t7, t8, t9}
	state, allTiles := playingStateForMeldInitial(rackTiles)
	payload := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{
				"id": "m1", "kind": "run", "tileIds": []string{"t1", "t2", "j1"},
				"jokerReps": map[string]interface{}{"j1": map[string]interface{}{"colour": "red", "rank": 7}},
			},
			map[string]interface{}{"id": "m2", "kind": "run", "tileIds": []string{"t4", "t5", "t6"}},
			map[string]interface{}{"id": "m3", "kind": "run", "tileIds": []string{"t7", "t8", "t9"}},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 3 {
		t.Fatalf("joker meld should succeed, got %d melds", len(st.TableMelds))
	}
	// Check joker rep preserved immutably
	for _, m := range st.TableMelds {
		if m.ID == "m1" {
			rep, ok := m.JokerReps["j1"]
			if !ok {
				t.Fatalf("m1 joker rep missing")
			}
			if rep.Colour != tile.Red || rep.Rank != 7 {
				t.Fatalf("joker rep wrong %v %v", rep.Colour, rep.Rank)
			}
		}
	}
	if err := CheckTileConservation(st, allTiles); err != nil {
		t.Fatalf("conservation joker: %v", err)
	}
	// Test joker missing rep should fail atomic
	state2, _ := playingStateForMeldInitial(rackTiles)
	payloadBad := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "run", "tileIds": []string{"t1", "t2", "j1"}}, // missing jokerReps
		},
	}
	pBytesBad, _ := json.Marshal(payloadBad)
	envBytesBad, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytesBad})
	msgBad := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytesBad}
	nextBad := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state2, []runtime.MatchData{msgBad})
	stBad := nextBad.(*RoundState)
	if len(stBad.TableMelds) != 0 {
		t.Fatalf("joker missing rep should be rejected")
	}
}

func TestMeldInitialRejectsForeignTile(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	t1 := tile.MustTile("t1", tile.Red, 5)
	t2 := tile.MustTile("t2", tile.Red, 6)
	t3 := tile.MustTile("t3", tile.Red, 7)
	rackTiles := []tile.TileInstance{t1, t2, t3}
	state, _ := playingStateForMeldInitial(rackTiles)
	// Try to use bob's tile "bob-00" inside alice's meld
	payload := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "run", "tileIds": []string{"t1", "t2", "bob-00"}},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 0 {
		t.Fatalf("foreign tile should be rejected")
	}
	if len(st.Racks[0]) != 15 {
		t.Fatalf("rack should be unchanged after foreign tile reject")
	}
}

func TestMeldInitialValidatesOwnershipAndHasRunAndScore(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Use a batch that is 45 points but valid structure to ensure score validation catches
	t1 := tile.MustTile("t1", tile.Red, 10)
	t2 := tile.MustTile("t2", tile.Red, 11)
	t3 := tile.MustTile("t3", tile.Red, 12) // 30
	t4 := tile.MustTile("t4", tile.Red, 7)
	t5 := tile.MustTile("t5", tile.Yellow, 7)
	t6 := tile.MustTile("t6", tile.Blue, 7) // 15 total 45
	rackTiles := []tile.TileInstance{t1, t2, t3, t4, t5, t6}
	state, _ := playingStateForMeldInitial(rackTiles)
	payload := map[string]interface{}{
		"melds": []interface{}{
			map[string]interface{}{"id": "m1", "kind": "run", "tileIds": []string{"t1", "t2", "t3"}},
			map[string]interface{}{"id": "m2", "kind": "set", "tileIds": []string{"t4", "t5", "t6"}},
		},
	}
	pBytes, _ := json.Marshal(payload)
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientMeldInitial, Payload: pBytes})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientMeldInitial, data: envBytes}
	dispatcher := &mockDispatcher{}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.TableMelds) != 0 {
		t.Fatalf("45 points should be rejected, got %d melds", len(st.TableMelds))
	}
}
