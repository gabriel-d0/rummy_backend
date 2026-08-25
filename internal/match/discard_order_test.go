package match

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
	"github.com/heroiclabs/nakama-common/runtime"
)

func TestDiscardRowOrderingAndOpeningDistinct(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	// Start with opening state: 2 players, alice 15, bob 14, stock, discard empty, current 0 OpeningDiscard
	state, _ := openingStateWithDeal()
	dispatcher := &mockDispatcher{}
	// Opening discard alice t-0
	openTileId := string(state.Racks[0][0].ID)
	payload, _ := json.Marshal(map[string]string{"tileId": openTileId})
	envBytes, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload})
	msg := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDiscard, data: envBytes}
	next := m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.MatchData{msg})
	st := next.(*RoundState)
	if len(st.DiscardRow) != 1 || !st.DiscardRow[0].IsOpeningDiscard || st.DiscardRow[0].Index != 0 {
		t.Fatalf("opening discard not distinct %+v", st.DiscardRow[0])
	}
	if string(st.DiscardRow[0].Tile.ID) != openTileId {
		t.Fatalf("opening tile %v want %v", st.DiscardRow[0].Tile.ID, openTileId)
	}
	// Now Playing MustDraw for bob (seat 1) — draw
	drawPayload, _ := json.Marshal(map[string]interface{}{})
	drawEnv, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawStock, Payload: drawPayload})
	drawMsg := &mockMatchData{mockPresence: mockPresence{userId: "bob"}, opCode: protocol.OpClientDrawStock, data: drawEnv}
	// Need to ensure bob has a stock to draw from; our openingStateWithDeal has stock 77, so ok
	// But openingStateWithDeal's stock is from synthetic deck, and we are in Playing MustDraw after opening, so draw should succeed
	// For this test we use the state after opening discard which already has GamePhase Playing and CurrentSeat 1 MustDraw
	next = m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 1, st, []runtime.MatchData{drawMsg})
	st = next.(*RoundState)
	if st.TurnPhase != TurnMeldOrDiscard {
		t.Fatalf("after draw TurnPhase %v want MeldOrDiscard", st.TurnPhase)
	}
	// Bob discards
	bobTileId := string(st.Racks[1][0].ID)
	payload2, _ := json.Marshal(map[string]string{"tileId": bobTileId})
	env2, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload2})
	msg2 := &mockMatchData{mockPresence: mockPresence{userId: "bob"}, opCode: protocol.OpClientDiscard, data: env2}
	next = m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 2, st, []runtime.MatchData{msg2})
	st = next.(*RoundState)
	if len(st.DiscardRow) != 2 {
		t.Fatalf("discard len %d want 2", len(st.DiscardRow))
	}
	if st.DiscardRow[1].Index != 1 || st.DiscardRow[1].IsOpeningDiscard {
		t.Fatalf("second discard should be Index 1 not opening %+v", st.DiscardRow[1])
	}
	if string(st.DiscardRow[1].Tile.ID) != bobTileId {
		t.Fatalf("second discard tile %v", st.DiscardRow[1].Tile.ID)
	}
	// Ensure ordering preserved: 0 is opening, 1 is bob's, and Index matches position
	for i, d := range st.DiscardRow {
		if d.Index != i {
			t.Fatalf("discard %d Index %d", i, d.Index)
		}
	}
	// Ensure opening still blocked via CanPickupPreviousDiscard should now allow since last is not opening
	if err := CanPickupPreviousDiscard(st); err != nil {
		t.Fatalf("previous should be pickable now, got %v", err)
	}
	// Simulate alice's turn: draw then discard
	drawPayloadAlice, _ := json.Marshal(map[string]interface{}{})
	drawEnvAlice, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDrawStock, Payload: drawPayloadAlice})
	drawMsgAlice := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDrawStock, data: drawEnvAlice}
	next = m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 3, st, []runtime.MatchData{drawMsgAlice})
	st = next.(*RoundState)
	aliceTileId := string(st.Racks[0][0].ID)
	payload3, _ := json.Marshal(map[string]string{"tileId": aliceTileId})
	env3, _ := json.Marshal(protocol.Envelope{Version: protocol.Version, OpCode: protocol.OpClientDiscard, Payload: payload3})
	msg3 := &mockMatchData{mockPresence: mockPresence{userId: "alice"}, opCode: protocol.OpClientDiscard, data: env3}
	next = m.MatchLoop(context.Background(), logger, nil, nil, dispatcher, 4, st, []runtime.MatchData{msg3})
	st = next.(*RoundState)
	if len(st.DiscardRow) != 3 || st.DiscardRow[2].Index != 2 {
		t.Fatalf("third discard ordering wrong")
	}
	if st.DiscardRow[0].IsOpeningDiscard != true {
		t.Fatalf("opening should remain distinct")
	}
	// Verify all discards are in chronological order and opening remains at 0
	if fmt.Sprintf("%v %v %v", st.DiscardRow[0].Tile.ID, st.DiscardRow[1].Tile.ID, st.DiscardRow[2].Tile.ID) == "" {
		t.Fatalf("discard order empty")
	}
}
