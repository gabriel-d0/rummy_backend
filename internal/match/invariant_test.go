package match

import (
	"fmt"
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// syntheticDeck creates 106 unique tiles with distinct IDs for invariant tests.
// It mimics the real 104 numbered + 2 jokers distribution without needing the
// deck factory (Day 15). IDs are tile-000..tile-105; first 104 are numbered,
// last 2 are jokers. This satisfies TileInstance.Validate and is deterministic.
func syntheticDeck() []tile.TileInstance {
	deck := make([]tile.TileInstance, 0, 106)
	// 104 numbered: we cycle colours and ranks to keep each valid, but IDs are what matters for conservation
	colours := []tile.Colour{tile.Red, tile.Yellow, tile.Blue, tile.Black}
	rank := tile.Rank(1)
	colourIdx := 0
	for i := 0; i < 104; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("tile-%03d", i))
		c := colours[colourIdx%len(colours)]
		deck = append(deck, tile.MustTile(id, c, rank))
		colourIdx++
		if colourIdx%len(colours) == 0 {
			rank++
			if rank > tile.RankMax {
				rank = tile.Rank(1)
			}
		}
	}
	deck = append(deck, tile.MustJoker("joker-104"))
	deck = append(deck, tile.MustJoker("joker-105"))
	return deck
}

func TestCheckTileConservationEmptyMelds(t *testing.T) {
	deck := syntheticDeck()
	players, _ := AssignSeats([]PlayerId{"alice", "bob"})
	// Simulate a valid deal: alice 15, bob 14, stock 106-29=77? Actually 106-29=77
	// For this test we just distribute all 106 across racks+stock+discard to satisfy conservation.
	// Use: alice 15, bob 14, discard 1 opening, stock 76 (15+14+1+76=106)
	rack0 := deck[0:15]
	rack1 := deck[15:29]
	discard := []DiscardEntry{{Tile: deck[29], IsOpeningDiscard: true, Index: 0}}
	stock := deck[30:106] // 76

	s := &RoundState{
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
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if err := CheckTileConservation(s, deck); err != nil {
		t.Fatalf("CheckTileConservation failed: %v", err)
	}
	r, stk, disc, melds := CountTiles(s)
	if r != 29 || stk != 76 || disc != 1 || melds != 0 {
		t.Fatalf("counts racks %d stock %d discard %d melds %d", r, stk, disc, melds)
	}
}

func TestCheckTileConservationDuplicateFails(t *testing.T) {
	deck := syntheticDeck()
	players, _ := AssignSeats([]PlayerId{"a", "b"})
	// duplicate deck[0] in both rack and stock
	rack0 := []tile.TileInstance{deck[0], deck[1]}
	rack1 := []tile.TileInstance{deck[2]}
	stock := []tile.TileInstance{deck[0], deck[3]} // dup deck[0]
	// Fill remaining to reach 106 but dup will still be caught early before missing
	// For simplicity, make stock contain dup and omit deck[4] to keep 106 count weird — but duplicate check fires first
	// So we just test duplicate detection, not full 106.
	// Build state with duplicate and also include all other tiles to reach 106, but duplicate will be duplicate.
	// Easiest: create state that has duplicate and then call CheckTileConservation which should report duplicate even if total 106 not met.
	// We will create a state that is otherwise valid 106 but with one tile duplicated and one missing.
	allButOne := deck[5:106] // 101 tiles
	stockFull := append(stock, allButOne...)
	s := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: rack0, 1: rack1},
		Stock:       stockFull,
		DiscardRow:  []DiscardEntry{},
		TableMelds:  []TableMeld{},
		CurrentSeat: 0,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	if err := CheckTileConservation(s, deck); err == nil {
		t.Fatalf("expected duplicate error")
	} else {
		t.Logf("got expected duplicate error: %v", err)
	}
}

func TestCheckTileConservationMissingFails(t *testing.T) {
	deck := syntheticDeck()
	players, _ := AssignSeats([]PlayerId{"a", "b"})
	// Omit last tile deck[105] entirely — have 105 tiles in state
	rack0 := deck[0:15]
	rack1 := deck[15:29]
	discard := []DiscardEntry{{Tile: deck[29], IsOpeningDiscard: true, Index: 0}}
	// stock is deck[30:105] (75) not including deck[105] (last joker) — so missing one
	stock := deck[30:105] // 75, total 15+14+1+75=105
	s := &RoundState{
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
	if err := CheckTileConservation(s, deck); err == nil {
		t.Fatalf("expected missing error")
	}
}

func TestCheckTileConservationWithMelds(t *testing.T) {
	deck := syntheticDeck()
	players, _ := AssignSeats([]PlayerId{"a", "b"})
	// Place 3 tiles into a meld, remove them from rack/stock counts accordingly
	meldTiles := deck[0:3]
	rack0 := deck[3:15]  // 12
	rack1 := deck[15:29] // 14
	discard := []DiscardEntry{{Tile: deck[29], IsOpeningDiscard: true, Index: 0}}
	stock := deck[30:106] // 76
	meld := TableMeld{
		ID:        "m1",
		Tiles:     meldTiles,
		JokerReps: map[tile.TileInstanceId]tile.TileInstance{},
		OwnerSeat: 0,
	}
	s := &RoundState{
		Players:     players,
		Racks:       map[Seat][]tile.TileInstance{0: rack0, 1: rack1},
		Stock:       stock,
		DiscardRow:  discard,
		TableMelds:  []TableMeld{meld},
		CurrentSeat: 1,
		GamePhase:   PhasePlaying,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	// 12+14+76+1+3 = 106
	if err := CheckTileConservation(s, deck); err != nil {
		t.Fatalf("with melds failed: %v", err)
	}
}
