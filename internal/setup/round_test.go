package setup

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/match"
)

func TestNewRoundStateCounts(t *testing.T) {
	cases := []struct {
		n         int
		wantStock int
	}{
		{2, 77},
		{3, 63},
		{4, 49},
	}
	for _, tc := range cases {
		ids := make([]match.PlayerId, tc.n)
		for i := 0; i < tc.n; i++ {
			ids[i] = match.PlayerId(string(rune('a' + i)))
		}
		state, allTiles, err := NewRoundState(ids, 42)
		if err != nil {
			t.Fatalf("n=%d err %v", tc.n, err)
		}
		if len(allTiles) != 106 {
			t.Fatalf("allTiles %d", len(allTiles))
		}
		if len(state.Players) != tc.n {
			t.Fatalf("players %d", len(state.Players))
		}
		if len(state.Racks[0]) != 15 {
			t.Fatalf("n=%d rack0 %d want 15", tc.n, len(state.Racks[0]))
		}
		for seat := 1; seat < tc.n; seat++ {
			if len(state.Racks[match.Seat(seat)]) != 14 {
				t.Fatalf("n=%d seat %d %d want 14", tc.n, seat, len(state.Racks[match.Seat(seat)]))
			}
		}
		if len(state.Stock) != tc.wantStock {
			t.Fatalf("n=%d stock %d want %d", tc.n, len(state.Stock), tc.wantStock)
		}
		if len(state.DiscardRow) != 0 {
			t.Fatalf("discard should be empty initially")
		}
		if len(state.TableMelds) != 0 {
			t.Fatalf("melds should be empty")
		}
		if state.CurrentSeat != 0 {
			t.Fatalf("CurrentSeat %v want 0", state.CurrentSeat)
		}
		if state.GamePhase != match.PhaseOpeningDiscard {
			t.Fatalf("GamePhase %v want OpeningDiscard", state.GamePhase)
		}
		if state.Winner != match.SeatInvalid {
			t.Fatalf("Winner should be invalid initially")
		}
		if err := state.Validate(); err != nil {
			t.Fatalf("Validate n=%d failed: %v", tc.n, err)
		}
		if err := match.CheckTileConservation(state, allTiles); err != nil {
			t.Fatalf("Conservation n=%d failed: %v", tc.n, err)
		}
	}
}

func TestNewRoundStateDeterministicSeed(t *testing.T) {
	ids := []match.PlayerId{"alice", "bob", "carol"}
	a, deckA, _ := NewRoundState(ids, 123)
	b, deckB, _ := NewRoundState(ids, 123)
	// Same seed → same shuffled deck and same racks
	for i := range deckA {
		if deckA[i].ID != deckB[i].ID {
			t.Fatalf("deck deterministic mismatch at %d", i)
		}
	}
	for seat := 0; seat < 3; seat++ {
		ra := a.Racks[match.Seat(seat)]
		rb := b.Racks[match.Seat(seat)]
		for i := range ra {
			if ra[i].ID != rb[i].ID {
				t.Fatalf("rack seat %d deterministic mismatch at %d", seat, i)
			}
		}
	}
	// Different seed should diverge (at least one deck position differs)
	_, deckC, _ := NewRoundState(ids, 999)
	diff := false
	for i := range deckA {
		if deckA[i].ID != deckC[i].ID {
			diff = true
			break
		}
	}
	if !diff {
		t.Fatalf("different seeds unexpectedly gave identical deck")
	}
}

func TestNewRoundStateInvalidCounts(t *testing.T) {
	if _, _, err := NewRoundState([]match.PlayerId{"solo"}, 1); err == nil {
		t.Fatalf("expected error n=1")
	}
	if _, _, err := NewRoundState([]match.PlayerId{"a", "b", "c", "d", "e"}, 1); err == nil {
		t.Fatalf("expected error n=5")
	}
	// duplicate IDs
	if _, _, err := NewRoundState([]match.PlayerId{"a", "a"}, 1); err == nil {
		t.Fatalf("expected duplicate error")
	}
}

func TestNewRoundStateWithDeck(t *testing.T) {
	deck := NewDeck()
	shuffled := Shuffle(deck, NewSeededRand(7))
	ids := []match.PlayerId{"p1", "p2", "p3", "p4"}
	state, allTiles, err := NewRoundStateWithDeck(ids, shuffled)
	if err != nil {
		t.Fatalf("WithDeck err %v", err)
	}
	if len(state.Players) != 4 || len(state.Racks[0]) != 15 || len(state.Stock) != 49 {
		t.Fatalf("4p counts wrong: players %d rack0 %d stock %d", len(state.Players), len(state.Racks[0]), len(state.Stock))
	}
	if err := match.CheckTileConservation(state, allTiles); err != nil {
		t.Fatalf("conservation: %v", err)
	}
	// Short deck should fail
	if _, _, err := NewRoundStateWithDeck(ids, shuffled[:50]); err == nil {
		t.Fatalf("expected short deck error")
	}
}
