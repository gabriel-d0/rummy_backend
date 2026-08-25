package setup

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/match"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestDealCounts(t *testing.T) {
	deck := NewDeck()
	cases := []struct {
		n           int
		wantOpening int
		wantOther   int
		wantStock   int
	}{
		{2, 15, 14, 77},
		{3, 15, 14, 63},
		{4, 15, 14, 49},
	}
	for _, tc := range cases {
		racks, stock, err := Deal(deck, tc.n)
		if err != nil {
			t.Fatalf("Deal n=%d err %v", tc.n, err)
		}
		if len(racks[0]) != tc.wantOpening {
			t.Fatalf("n=%d seat0 %d want %d", tc.n, len(racks[0]), tc.wantOpening)
		}
		for seat := 1; seat < tc.n; seat++ {
			if len(racks[match.Seat(seat)]) != tc.wantOther {
				t.Fatalf("n=%d seat %d %d want %d", tc.n, seat, len(racks[match.Seat(seat)]), tc.wantOther)
			}
		}
		if len(stock) != tc.wantStock {
			t.Fatalf("n=%d stock %d want %d", tc.n, len(stock), tc.wantStock)
		}
		// No duplicate IDs across racks+stock
		seen := map[tile.TileInstanceId]bool{}
		for seat := 0; seat < tc.n; seat++ {
			for _, tl := range racks[match.Seat(seat)] {
				if seen[tl.ID] {
					t.Fatalf("duplicate %v n=%d", tl.ID, tc.n)
				}
				seen[tl.ID] = true
			}
		}
		for _, tl := range stock {
			if seen[tl.ID] {
				t.Fatalf("duplicate stock %v n=%d", tl.ID, tc.n)
			}
			seen[tl.ID] = true
		}
		if len(seen) != 106 {
			t.Fatalf("n=%d seen %d want 106", tc.n, len(seen))
		}
	}
}

func TestDealInvalidCounts(t *testing.T) {
	deck := NewDeck()
	if _, _, err := Deal(deck, 1); err == nil {
		t.Fatalf("expected error n=1")
	}
	if _, _, err := Deal(deck, 5); err == nil {
		t.Fatalf("expected error n=5")
	}
	if _, _, err := Deal(deck[:10], 2); err == nil {
		t.Fatalf("expected error deck 10")
	}
}

func TestDealDeterministic(t *testing.T) {
	deck := NewDeck()
	shuffled := Shuffle(deck, NewSeededRand(99))
	racks1, stock1, _ := Deal(shuffled, 4)
	racks2, stock2, _ := Deal(shuffled, 4)
	for seat := 0; seat < 4; seat++ {
		a, b := racks1[match.Seat(seat)], racks2[match.Seat(seat)]
		if len(a) != len(b) {
			t.Fatalf("rack len mismatch")
		}
		for i := range a {
			if a[i].ID != b[i].ID {
				t.Fatalf("deterministic mismatch seat %d idx %d", seat, i)
			}
		}
	}
	if len(stock1) != len(stock2) {
		t.Fatalf("stock len")
	}
	for i := range stock1 {
		if stock1[i].ID != stock2[i].ID {
			t.Fatalf("stock deterministic mismatch %d", i)
		}
	}
}

func TestDealConservationWithState(t *testing.T) {
	deck := NewDeck()
	shuffled := Shuffle(deck, NewSeededRand(7))
	n := 3
	racks, stock, _ := Deal(shuffled, n)
	players, _ := match.AssignSeats([]match.PlayerId{"a", "b", "c"})
	s := &match.RoundState{
		Players:     players,
		Racks:       racks,
		Stock:       stock,
		DiscardRow:  []match.DiscardEntry{},
		TableMelds:  []match.TableMeld{},
		CurrentSeat: 0,
		GamePhase:   match.PhaseOpeningDiscard,
		TurnPhase:   match.TurnMustDraw, // not used in OpeningDiscard but set
		Winner:      match.SeatInvalid,
	}
	// Shuffle preserves deck multiset, Deal distributes all 106, so conservation should hold even though shuffled order differs from NewDeck order.
	// For invariant we need the exact shuffled deck as allTiles, not NewDeck original order — but IDs same, so NewDeck pool is equivalent.
	if err := match.CheckTileConservation(s, shuffled); err != nil {
		t.Fatalf("conservation failed: %v", err)
	}
}
