package match

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
	"github.com/heroiclabs/nakama-common/runtime"
)

// TestCentralizedViewProjection ensures PublicView/PrivateView are centralized and versioned.
func TestCentralizedViewProjection(t *testing.T) {
	players, _ := AssignSeats([]PlayerId{"alice", "bob"})
	state := &RoundState{
		Players: players,
		Racks: map[Seat][]tile.TileInstance{
			0: {tile.MustTile("alice-1", tile.Red, 5), tile.MustTile("alice-2", tile.Red, 6)},
			1: {tile.MustTile("bob-1", tile.Blue, 7)},
		},
		Stock:       []tile.TileInstance{tile.MustTile("stock-1", tile.Yellow, 3)},
		DiscardRow:  []DiscardEntry{{Tile: tile.MustTile("disc-1", tile.Red, 7), IsOpeningDiscard: true, Index: 0}},
		TableMelds:  []TableMeld{{ID: "m1", Kind: "run", Tiles: []tile.TileInstance{tile.MustTile("ex1", tile.Red, 1), tile.MustTile("ex2", tile.Red, 2), tile.MustTile("ex3", tile.Red, 3)}, JokerReps: map[tile.TileInstanceId]tile.TileInstance{}, OwnerSeat: 0}},
		CurrentSeat: 0,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	pub := PublicView(state)
	if pub.Version != 1 {
		t.Fatalf("PublicSnapshot version %d want 1", pub.Version)
	}
	if pub.GamePhase != "Playing" || pub.TurnPhase != "MustDraw" {
		t.Fatalf("phase mismatch %v/%v", pub.GamePhase, pub.TurnPhase)
	}
	privAlice := PrivateView(state, 0)
	if privAlice.OwnSeat != 0 || len(privAlice.OwnRack) != 2 {
		t.Fatalf("PrivateView alice OwnSeat %v rack %d", privAlice.OwnSeat, len(privAlice.OwnRack))
	}
	privBob := PrivateView(state, 1)
	if privBob.OwnSeat != 1 || len(privBob.OwnRack) != 1 {
		t.Fatalf("PrivateView bob")
	}
	// Public should not contain private rack IDs
	b, _ := json.Marshal(pub)
	s := string(b)
	if strings.Contains(s, "alice-1") || strings.Contains(s, "bob-1") {
		t.Fatalf("public leaked private rack IDs")
	}
	// Private should contain own but not other
	ba, _ := json.Marshal(privAlice)
	if !strings.Contains(string(ba), "alice-1") || strings.Contains(string(ba), "bob-1") {
		t.Fatalf("private alice should contain own not bob")
	}
}

// TestReconnectionRestoresPrivateRack verifies that a disconnected player
// rejoining receives their private rack plus public state, and others never receive it.
func TestReconnectionRestoresPrivateRack(t *testing.T) {
	m := &RummyMatch{}
	logger := &testLogger{}
	state, _, _ := m.MatchInit(context.Background(), logger, nil, nil, nil)
	dispatcher := &mockDispatcher{}
	// Join alice and bob in Waiting
	state = m.MatchJoin(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.Presence{newPresence("alice"), newPresence("bob")})
	st := state.(*RoundState)
	// Simulate that game has started and alice has a known rack
	st.GamePhase = PhasePlaying
	st.CurrentSeat = 0
	st.TurnPhase = TurnMustDraw
	st.Racks[0] = []tile.TileInstance{tile.MustTile("alice-secret-1", tile.Red, 5), tile.MustTile("alice-secret-2", tile.Red, 6)}
	st.Racks[1] = []tile.TileInstance{tile.MustTile("bob-secret-1", tile.Blue, 7)}
	st.Stock = []tile.TileInstance{tile.MustTile("stock-1", tile.Yellow, 3)}
	st.DiscardRow = []DiscardEntry{{Tile: tile.MustTile("disc-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0}}
	st.TableMelds = []TableMeld{{ID: "m1", Kind: "run", Tiles: []tile.TileInstance{tile.MustTile("ex1", tile.Red, 1), tile.MustTile("ex2", tile.Red, 2), tile.MustTile("ex3", tile.Red, 3)}, JokerReps: map[tile.TileInstanceId]tile.TileInstance{}, OwnerSeat: 0}}

	// Alice disconnects (MatchLeave keeps her)
	state = m.MatchLeave(context.Background(), logger, nil, nil, dispatcher, 0, state, []runtime.Presence{newPresence("alice")})
	st = state.(*RoundState)
	if len(st.Players) != 2 {
		t.Fatalf("players should be kept for reconnect, got %d", len(st.Players))
	}
	if _, ok := st.Racks[0]; !ok {
		t.Fatalf("alice rack should be kept")
	}
	// Alice reconnects: send new presence with same userId but different session
	dispatcher2 := &mockDispatcher{}
	presAlice2 := newPresence("alice")
	presAlice2.sessionId = "sess-alice-2"
	state = m.MatchJoin(context.Background(), logger, nil, nil, dispatcher2, 1, state, []runtime.Presence{presAlice2})
	st = state.(*RoundState)
	// Should still have 2 players, not duplicate
	if len(st.Players) != 2 {
		t.Fatalf("reconnect should not duplicate, got %d", len(st.Players))
	}
	// Dispatcher should have sent PrivateView to alice with her secret rack
	foundSnapshot := false
	for _, b := range dispatcher2.broadcasts {
		if b.Op == 100 { // OpServerState
			s := string(b.Data)
			if strings.Contains(s, "alice-secret-1") && !strings.Contains(s, "bob-secret-1") {
				// Check that it was sent only to alice (To contains alice)
				if len(b.To) == 1 && b.To[0].GetUserId() == "alice" {
					foundSnapshot = true
				}
			}
			// Ensure bob's secret not leaked in alice's snapshot
			if strings.Contains(s, "bob-secret-1") && len(b.To) == 1 && b.To[0].GetUserId() == "alice" {
				t.Fatalf("alice snapshot leaked bob's rack")
			}
		}
	}
	if !foundSnapshot {
		t.Fatalf("reconnection did not send PrivateView to alice with her rack")
	}
}

// TestRedactionRoundComplete ensures winner state does not leak private racks.
func TestRedactionRoundComplete(t *testing.T) {
	state := &RoundState{
		Players: []PlayerState{{ID: "alice", Seat: 0, HasOpened: true}, {ID: "bob", Seat: 1, HasOpened: true}},
		Racks: map[Seat][]tile.TileInstance{
			0: {}, // winner empty
			1: {tile.MustTile("bob-secret-1", tile.Blue, 7), tile.MustTile("bob-secret-2", tile.Red, 9)},
		},
		Stock:       []tile.TileInstance{},
		DiscardRow:  []DiscardEntry{{Tile: tile.MustTile("disc-open", tile.Red, 7), IsOpeningDiscard: true, Index: 0}, {Tile: tile.MustTile("disc-win", tile.Yellow, 3), IsOpeningDiscard: false, Index: 1}},
		TableMelds:  []TableMeld{{ID: "m1", Kind: "run", Tiles: []tile.TileInstance{tile.MustTile("ex1", tile.Red, 1), tile.MustTile("ex2", tile.Red, 2), tile.MustTile("ex3", tile.Red, 3)}, JokerReps: map[tile.TileInstanceId]tile.TileInstance{}, OwnerSeat: 0}},
		CurrentSeat: 0,
		GamePhase:   PhaseRoundComplete,
		TurnPhase:   TurnMeldOrDiscard,
		Winner:      0,
	}
	// Public should not contain bob's secrets
	pub := PublicView(state)
	b, _ := json.Marshal(pub)
	s := string(b)
	if strings.Contains(s, "bob-secret") {
		t.Fatalf("public RoundComplete leaked bob's rack: %s", s)
	}
	if !strings.Contains(s, `"winner":0`) {
		t.Fatalf("public should contain winner 0, got %s", s)
	}
	if pub.Version != 1 {
		t.Fatalf("version %d want 1", pub.Version)
	}
	// Private for winner (empty rack)
	privAlice := PrivateView(state, 0)
	if len(privAlice.OwnRack) != 0 {
		t.Fatalf("winner rack should be empty, got %d", len(privAlice.OwnRack))
	}
	if privAlice.Winner != 0 {
		t.Fatalf("private winner %v want 0", privAlice.Winner)
	}
	// Private for bob should contain his secrets but not winner's (winner empty anyway)
	privBob := PrivateView(state, 1)
	bb, _ := json.Marshal(privBob)
	if !strings.Contains(string(bb), "bob-secret-1") {
		t.Fatalf("bob private should contain his rack")
	}
	if strings.Contains(string(bb), "alice-secret") {
		t.Fatalf("bob private leaked alice")
	}
}

// TestSnapshotVersioning ensures all snapshots are versioned and stable.
func TestSnapshotVersioning(t *testing.T) {
	players, _ := AssignSeats([]PlayerId{"alice", "bob"})
	state := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: {tile.MustTile("a1", tile.Red, 1)}, 1: {tile.MustTile("b1", tile.Blue, 2)}},
		Stock:       []tile.TileInstance{},
		DiscardRow:  []DiscardEntry{},
		TableMelds:  []TableMeld{},
		CurrentSeat: 0,
		GamePhase:   PhaseWaiting,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	pub := PublicView(state)
	if pub.Version != 1 {
		t.Fatalf("version %d", pub.Version)
	}
	priv := PrivateView(state, 0)
	if priv.Version != 1 {
		t.Fatalf("private version %d", priv.Version)
	}
	// Ensure JSON marshals deterministically
	b1, _ := json.Marshal(pub)
	b2, _ := json.Marshal(PublicView(state))
	if string(b1) != string(b2) {
		t.Fatalf("public snapshot not stable")
	}
}
