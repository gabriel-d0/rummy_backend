package setup

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestShuffleDeterministicFixedSeed(t *testing.T) {
	deck := NewDeck()
	a := Shuffle(deck, NewSeededRand(42))
	b := Shuffle(deck, NewSeededRand(42))
	if len(a) != len(b) {
		t.Fatalf("len %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i].ID != b[i].ID {
			t.Fatalf("deterministic mismatch at %d: %v vs %v", i, a[i].ID, b[i].ID)
		}
	}
}

func TestShuffleChangesOrder(t *testing.T) {
	deck := NewDeck()
	shuffled := Shuffle(deck, NewSeededRand(123))
	// Very likely not identical; with 106 tiles, probability of identity is 1/106! negligible
	same := true
	for i := range deck {
		if deck[i].ID != shuffled[i].ID {
			same = false
			break
		}
	}
	if same {
		t.Fatalf("Shuffle with seed 123 unexpectedly left order unchanged")
	}
}

func TestShuffleNoLostDuplicated(t *testing.T) {
	deck := NewDeck()
	shuffled := Shuffle(deck, NewSeededRand(99))
	if len(shuffled) != len(deck) {
		t.Fatalf("len %d want %d", len(shuffled), len(deck))
	}
	// No lost/duplicated: map IDs
	seen := map[tile.TileInstanceId]int{}
	for _, tl := range deck {
		seen[tl.ID] = 0
	}
	for _, tl := range shuffled {
		if _, ok := seen[tl.ID]; !ok {
			t.Fatalf("shuffled has unknown ID %v", tl.ID)
		}
		seen[tl.ID]++
	}
	for id, cnt := range seen {
		if cnt != 1 {
			t.Fatalf("ID %v count %d want 1", id, cnt)
		}
	}
	// All shuffled tiles still Validate
	for _, tl := range shuffled {
		if err := tl.Validate(); err != nil {
			t.Fatalf("shuffled tile %v invalid: %v", tl.ID, err)
		}
	}
}

func TestShuffleOriginalUnmodified(t *testing.T) {
	deck := NewDeck()
	origFirst := deck[0].ID
	_ = Shuffle(deck, NewSeededRand(7))
	if deck[0].ID != origFirst {
		t.Fatalf("Shuffle mutated original deck")
	}
}

func TestShuffleDifferentSeedsDiverge(t *testing.T) {
	deck := NewDeck()
	a := Shuffle(deck, NewSeededRand(1))
	b := Shuffle(deck, NewSeededRand(2))
	same := true
	for i := range a {
		if a[i].ID != b[i].ID {
			same = false
			break
		}
	}
	if same {
		t.Fatalf("different seeds gave same shuffle order")
	}
}

func TestShufflePanicsNilRand(t *testing.T) {
	deck := NewDeck()
	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		Shuffle(deck, nil)
	}()
	if !didPanic {
		t.Fatalf("expected panic for nil Rand")
	}
}
