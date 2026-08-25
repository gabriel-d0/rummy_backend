package match

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// TestTileIdentity verifies that two copies of the same colour/rank are
// distinct instances with unique IDs but equal face values — per docs/rules-decisions.md:1.1
// and AGENTS.md:41.
func TestTileIdentity(t *testing.T) {
	a := tile.MustTile("id-same-1", tile.Red, 5)
	b := tile.MustTile("id-same-2", tile.Red, 5)
	if a.ID == b.ID {
		t.Fatalf("distinct copies must have different IDs")
	}
	if a.Colour != b.Colour || a.Rank != b.Rank {
		t.Fatalf("same face expected")
	}
	if a.ID == b.ID {
		t.Fatalf("ID equality should not be assumed from face")
	}
	// Joker vs numbered same ID not allowed — validated
	j := tile.MustJoker("j-1")
	if j.IsJoker == a.IsJoker {
		t.Fatalf("joker vs numbered IsJoker mismatch")
	}
	if err := j.Validate(); err != nil {
		t.Fatalf("joker validate: %v", err)
	}
}

// TestSeatsTurnDirection verifies deterministic anticlockwise turn order
// per docs/terminology.md Seat and docs/rules-decisions.md:1.2.
func TestSeatsTurnDirection(t *testing.T) {
	// 2 players: 0→1→0
	if s, _ := NextSeat(0, 2); s != 1 {
		t.Fatalf("2p 0→%v want 1", s)
	}
	if s, _ := NextSeat(1, 2); s != 0 {
		t.Fatalf("2p 1→%v want 0", s)
	}
	// 3 players: 0→1→2→0
	seq3 := []Seat{0, 1, 2, 0, 1}
	for i := 0; i < len(seq3)-1; i++ {
		next, _ := NextSeat(seq3[i], 3)
		if next != seq3[i+1] {
			t.Fatalf("3p seat %v next %v want %v", seq3[i], next, seq3[i+1])
		}
	}
	// 4 players: full cycle
	for n := 2; n <= 4; n++ {
		for seat := 0; seat < n; seat++ {
			next, _ := NextSeat(Seat(seat), n)
			prev, _ := PrevSeat(next, n)
			if prev != Seat(seat) {
				t.Fatalf("n=%d seat %d next %v prev %v", n, seat, next, prev)
			}
		}
	}
	// Deterministic AssignSeats join order
	players, _ := AssignSeats([]PlayerId{"p1", "p2", "p3"})
	if players[0].Seat != 0 || players[1].Seat != 1 || players[2].Seat != 2 {
		t.Fatalf("AssignSeats not join-order deterministic: %+v", players)
	}
}

// TestBaseStateConstruction verifies a minimal but valid RoundState can be
// constructed and passes Validate and CheckTileConservation — the Day 12/13
// contract before deck factory (Day 15) exists.
func TestBaseStateConstruction(t *testing.T) {
	deck := syntheticDeck() // 106
	players, _ := AssignSeats([]PlayerId{"alice", "bob", "carol"})
	// Simulate deal: opening player 15, others 14, remainder stock, 1 opening discard flag ready
	// For this domain test we just ensure a valid distribution exists and validates.
	rack0 := deck[0:15]
	rack1 := deck[15:29]
	rack2 := deck[29:43]
	// Put one tile as opening discard, rest stock
	discard := []DiscardEntry{{Tile: deck[43], IsOpeningDiscard: true, Index: 0}}
	stock := deck[44:106] // 62

	s := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: rack0, 1: rack1, 2: rack2},
		Stock:       stock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{},
		CurrentSeat: 1, // after opening discard, next seat 1
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if err := CheckTileConservation(s, deck); err != nil {
		t.Fatalf("Conservation: %v", err)
	}
	r, stk, disc, melds := CountTiles(s)
	if r != 43 || stk != 62 || disc != 1 || melds != 0 {
		t.Fatalf("counts racks %d stock %d disc %d melds %d", r, stk, disc, melds)
	}
	if r+stk+disc+melds != 106 {
		t.Fatalf("total not 106")
	}
}

// TestInvariantFailures verifies that conservation catches duplicates and
// missing tiles — the load-bearing property for future shuffle/deal.
func TestInvariantFailures(t *testing.T) {
	deck := syntheticDeck()
	players, _ := AssignSeats([]PlayerId{"a", "b"})
	// Valid baseline
	rack0 := deck[0:15]
	rack1 := deck[15:29]
	discard := []DiscardEntry{{Tile: deck[29], IsOpeningDiscard: true, Index: 0}}
	stock := deck[30:106]
	valid := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: rack0, 1: rack1},
		Stock:       stock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{},
		CurrentSeat: 1,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	if err := CheckTileConservation(valid, deck); err != nil {
		t.Fatalf("valid should pass: %v", err)
	}

	// Duplicate: copy deck[0] into both rack0 and stock
	dupStock := append([]tile.TileInstance{deck[0]}, stock[1:]...)
	dupState := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: rack0, 1: rack1},
		Stock:       dupStock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{},
		CurrentSeat: 1,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	if err := CheckTileConservation(dupState, deck); err == nil {
		t.Fatalf("expected duplicate error")
	}

	// Missing: omit deck[105]
	missingStock := stock[:len(stock)-1] // drop last joker
	missingState := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: rack0, 1: rack1},
		Stock:       missingStock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{},
		CurrentSeat: 1,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	if err := CheckTileConservation(missingState, deck); err == nil {
		t.Fatalf("expected missing error")
	}
}
