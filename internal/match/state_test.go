package match

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestRoundStateValidateMinimal(t *testing.T) {
	players, _ := AssignSeats([]PlayerId{"alice", "bob"})
	// minimal racks: 14 each with valid tiles (use distinct IDs)
	rack0 := []tile.TileInstance{
		tile.MustTile("t0-1", tile.Red, 1),
		tile.MustTile("t0-2", tile.Red, 2),
	}
	rack1 := []tile.TileInstance{
		tile.MustTile("t1-1", tile.Blue, 5),
		tile.MustTile("t1-2", tile.Black, 13),
	}
	s := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: rack0, 1: rack1},
		Stock:       []tile.TileInstance{tile.MustTile("s1", tile.Yellow, 3)},
		DiscardRow:  []DiscardEntry{{Tile: tile.MustTile("d1", tile.Red, 7), IsOpeningDiscard: true, Index: 0}},
		TableMelds:  []TableMeld{},
		CurrentSeat: 1,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate minimal failed: %v", err)
	}
}

func TestRoundStateValidateRejectsDuplicateRackTile(t *testing.T) {
	players, _ := AssignSeats([]PlayerId{"a", "b"})
	dupTile := tile.MustTile("dup", tile.Red, 5)
	s := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: {dupTile, dupTile}, 1: {tile.MustTile("t2", tile.Blue, 2)}},
		Stock:       []tile.TileInstance{},
		DiscardRow:  []DiscardEntry{},
		TableMelds:  []TableMeld{},
		CurrentSeat: 0,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	if err := s.Validate(); err == nil {
		t.Fatalf("expected duplicate tile error")
	}
}

func TestRoundStateOpeningDiscardOnlyFirst(t *testing.T) {
	players, _ := AssignSeats([]PlayerId{"a", "b"})
	s := &RoundState{
		Players: players,
		Racks: map[Seat][]tile.TileInstance{
			0: {tile.MustTile("t0", tile.Red, 1)},
			1: {tile.MustTile("t1", tile.Blue, 2)},
		},
		Stock: []tile.TileInstance{},
		DiscardRow: []DiscardEntry{
			{Tile: tile.MustTile("d0", tile.Red, 1), IsOpeningDiscard: true, Index: 0},
			{Tile: tile.MustTile("d1", tile.Red, 2), IsOpeningDiscard: true, Index: 1}, // invalid second opening
		},
		TableMelds:  []TableMeld{},
		CurrentSeat: 1,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	if err := s.Validate(); err == nil {
		t.Fatalf("expected opening discard only first error")
	}
}

func TestTableMeldValidateJokerRep(t *testing.T) {
	j := tile.MustJoker("j1")
	rep := tile.MustTile("rep-7-red", tile.Red, 7)
	m := TableMeld{
		ID:        "m1",
		Tiles:     []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Yellow, 5), j},
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{"j1": rep},
		OwnerSeat: 0,
	}
	if err := m.Validate(); err != nil {
		t.Fatalf("joker meld Validate failed: %v", err)
	}
	// missing rep should fail
	m2 := TableMeld{ID: "m2", Tiles: []tile.TileInstance{tile.MustTile("t1", tile.Red, 5), tile.MustTile("t2", tile.Yellow, 5), j}, JokerReps: map[tile.TileInstanceId]tile.TileInstance{}, OwnerSeat: 0}
	if err := m2.Validate(); err == nil {
		t.Fatalf("expected missing rep error")
	}
}
