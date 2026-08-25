package setup

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/match"
)

// TestSetupInvariantsAllPlayersAndSeeds verifies tile conservation and stock
// counts for many seeds and player counts — the formal Day 21 regression.
func TestSetupInvariantsAllPlayersAndSeeds(t *testing.T) {
	seeds := []int64{0, 1, 42, 123, 999, 2026}
	names := []match.PlayerId{"alice", "bob", "carol", "dana"}
	for _, n := range []int{2, 3, 4} {
		for _, seed := range seeds {
			ids := names[:n]

			state, allTiles, err := NewRoundState(ids, seed)
			if err != nil {
				t.Fatalf("n=%d seed=%d err %v", n, seed, err)
			}
			wantStock := 77
			if n == 3 {
				wantStock = 63
			} else if n == 4 {
				wantStock = 49
			}
			if len(state.Stock) != wantStock {
				t.Fatalf("n=%d seed=%d stock %d want %d", n, seed, len(state.Stock), wantStock)
			}
			if len(state.Racks[0]) != 15 {
				t.Fatalf("n=%d seed=%d rack0 %d want 15", n, seed, len(state.Racks[0]))
			}
			for seat := 1; seat < n; seat++ {
				if len(state.Racks[match.Seat(seat)]) != 14 {
					t.Fatalf("n=%d seed=%d seat %d %d want 14", n, seed, seat, len(state.Racks[match.Seat(seat)]))
				}
			}
			if err := state.Validate(); err != nil {
				t.Fatalf("n=%d seed=%d Validate %v", n, seed, err)
			}
			if err := match.CheckTileConservation(state, allTiles); err != nil {
				t.Fatalf("n=%d seed=%d conservation %v", n, seed, err)
			}
			r, stk, disc, melds := match.CountTiles(state)
			if r+stk+disc+melds != 106 {
				t.Fatalf("n=%d seed=%d total %d want 106", n, seed, r+stk+disc+melds)
			}
			if disc != 0 || melds != 0 {
				t.Fatalf("n=%d seed=%d disc %d melds %d want 0", n, seed, disc, melds)
			}
			// Also check that opening player is always Seat 0
			if state.CurrentSeat != 0 {
				t.Fatalf("n=%d seed=%d CurrentSeat %v want 0", n, seed, state.CurrentSeat)
			}
			if state.GamePhase != match.PhaseOpeningDiscard {
				t.Fatalf("n=%d seed=%d GamePhase %v want OpeningDiscard", n, seed, state.GamePhase)
			}
		}
	}
}
