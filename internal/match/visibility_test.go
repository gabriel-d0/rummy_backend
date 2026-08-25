package match

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestPublicViewHidesRacks(t *testing.T) {
	// Create a round with 2 players and distinct tile IDs
	players, _ := AssignSeats([]PlayerId{"alice", "bob"})
	state := &RoundState{
		Players: players,
		Racks: map[Seat][]tile.TileInstance{
			0: {tile.MustTile("alice-t0", tile.Red, 1), tile.MustTile("alice-t1", tile.Blue, 2)},
			1: {tile.MustTile("bob-t0", tile.Yellow, 3), tile.MustTile("bob-t1", tile.Black, 4)},
		},
		Stock:       []tile.TileInstance{tile.MustTile("stock-1", tile.Red, 5)},
		DiscardRow:  []DiscardEntry{},
		TableMelds:  []TableMeld{},
		CurrentSeat: 0,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}

	pub := PublicView(state)
	b, _ := json.Marshal(pub)
	s := string(b)

	// Public view must not contain private tile IDs from racks (only counts/discard/melds)
	// Check that some rack tile IDs are not leaked via string search
	rack0ID := string(state.Racks[0][0].ID)
	rack1ID := string(state.Racks[1][0].ID)
	// Public JSON should contain discard row? Currently empty, but check racks not leaked
	// Since racks are counts only, tile IDs should not appear as substrings in public JSON
	// However stock count, not IDs, so we can verify via absence of rack-specific IDs
	// But discard row is empty, melds empty, so the only IDs that could leak are rack IDs
	// We check that neither rack ID appears in public JSON
	if strings.Contains(s, rack0ID) {
		t.Fatalf("public view leaked rack0 tile %v", rack0ID)
	}
	if strings.Contains(s, rack1ID) {
		t.Fatalf("public view leaked rack1 tile %v", rack1ID)
	}
	// But public should contain rack counts
	if !strings.Contains(s, `"rackCount"`) {
		t.Fatalf("public missing rackCount")
	}
	// Private view for alice should contain her own rack but not bob's
	privAlice := PrivateView(state, 0)
	bAlice, _ := json.Marshal(privAlice)
	sAlice := string(bAlice)
	if !strings.Contains(sAlice, rack0ID) {
		t.Fatalf("private alice should contain own tile %v", rack0ID)
	}
	if strings.Contains(sAlice, rack1ID) {
		t.Fatalf("private alice leaked bob's tile %v", rack1ID)
	}
	// Private for bob opposite
	privBob := PrivateView(state, 1)
	bBob, _ := json.Marshal(privBob)
	sBob := string(bBob)
	if !strings.Contains(sBob, rack1ID) {
		t.Fatalf("private bob should contain own tile %v", rack1ID)
	}
	if strings.Contains(sBob, rack0ID) {
		t.Fatalf("private bob leaked alice's tile %v", rack0ID)
	}
}

func TestPrivateViewOpponentNotLeakedViaMarshal(t *testing.T) {
	players, _ := AssignSeats([]PlayerId{"alice", "bob", "carol"})
	state := &RoundState{
		Players: players,
		Racks: map[Seat][]tile.TileInstance{
			0: {tile.MustTile("alice-t0", tile.Red, 1)},
			1: {tile.MustTile("bob-t0", tile.Blue, 2)},
			2: {tile.MustTile("carol-t0", tile.Yellow, 3)},
		},
		Stock:       []tile.TileInstance{},
		DiscardRow:  []DiscardEntry{{Tile: tile.MustTile("disc-1", tile.Red, 5), IsOpeningDiscard: true, Index: 0}},
		TableMelds:  []TableMeld{{ID: "m1", Tiles: []tile.TileInstance{tile.MustTile("m1-t1", tile.Blue, 2), tile.MustTile("m1-t2", tile.Blue, 3), tile.MustTile("m1-t3", tile.Blue, 4)}, JokerReps: map[tile.TileInstanceId]tile.TileInstance{}}},
		CurrentSeat: 0,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}

	for _, seat := range []Seat{0, 1, 2} {
		priv := PrivateView(state, seat)
		b, _ := json.Marshal(priv)
		s := string(b)
		// Own rack must be present
		ownID := string(state.Racks[seat][0].ID)
		if !strings.Contains(s, ownID) {
			t.Fatalf("seat %v private missing own %v", seat, ownID)
		}
		// Other racks must not be present
		for otherSeat := range state.Racks {
			if otherSeat == seat {
				continue
			}
			otherID := string(state.Racks[otherSeat][0].ID)
			if strings.Contains(s, otherID) {
				t.Fatalf("seat %v private leaked other seat %v tile %v", seat, otherSeat, otherID)
			}
		}
		// Public discard/meld IDs are allowed to appear (they are public)
		if !strings.Contains(s, "disc-1") {
			t.Fatalf("public discard should be in private view")
		}
	}
}

func TestPublicViewStructure(t *testing.T) {
	players, _ := AssignSeats([]PlayerId{"a", "b"})
	state := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: {tile.MustTile("t0", tile.Red, 1)}, 1: {tile.MustTile("t1", tile.Blue, 2)}},
		Stock:       []tile.TileInstance{tile.MustTile("s1", tile.Yellow, 3)},
		DiscardRow:  []DiscardEntry{},
		TableMelds:  []TableMeld{},
		CurrentSeat: 0,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	pub := PublicView(state)
	if pub.StockCount != 1 {
		t.Fatalf("stockCount %d", pub.StockCount)
	}
	if pub.Players[0].RackCount != 1 || pub.Players[1].RackCount != 1 {
		t.Fatalf("rack counts")
	}
	if pub.Version != 1 {
		t.Fatalf("version")
	}
	priv := PrivateView(state, 0)
	if len(priv.OwnRack) != 1 || priv.OwnRack[0].ID != "t0" {
		t.Fatalf("ownRack")
	}
}
