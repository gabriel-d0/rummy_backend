package setup

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestNewDeckTotalCount(t *testing.T) {
	deck := NewDeck()
	if len(deck) != 106 {
		t.Fatalf("deck len %d want 106", len(deck))
	}
}

func TestNewDeckCounts(t *testing.T) {
	deck := NewDeck()
	numbered, jokers := 0, 0
	for _, tl := range deck {
		if err := tl.Validate(); err != nil {
			t.Fatalf("tile %v invalid: %v", tl.ID, err)
		}
		if tl.IsJoker {
			jokers++
		} else {
			numbered++
		}
	}
	if numbered != 104 {
		t.Fatalf("numbered %d want 104", numbered)
	}
	if jokers != 2 {
		t.Fatalf("jokers %d want 2", jokers)
	}
}

func TestNewDeckTwoPerColourRank(t *testing.T) {
	deck := NewDeck()
	// count per colour+rank
	type key struct {
		c tile.Colour
		r tile.Rank
	}
	m := map[key]int{}
	for _, tl := range deck {
		if tl.IsJoker {
			continue
		}
		k := key{tl.Colour, tl.Rank}
		m[k]++
	}
	colours := []tile.Colour{tile.Red, tile.Yellow, tile.Blue, tile.Black}
	for _, c := range colours {
		for r := tile.RankMin; r <= tile.RankMax; r++ {
			k := key{c, r}
			if m[k] != 2 {
				t.Fatalf("colour %v rank %v count %d want 2", c, r, m[k])
			}
		}
	}
	if len(m) != 52 {
		t.Fatalf("distinct colour/rank combos %d want 52", len(m))
	}
}

func TestNewDeckAllIDsUnique(t *testing.T) {
	deck := NewDeck()
	seen := map[tile.TileInstanceId]bool{}
	for _, tl := range deck {
		if seen[tl.ID] {
			t.Fatalf("duplicate ID %v", tl.ID)
		}
		seen[tl.ID] = true
		if !tl.ID.IsValid() {
			t.Fatalf("empty ID")
		}
	}
	if len(seen) != 106 {
		t.Fatalf("unique IDs %d want 106", len(seen))
	}
}

func TestNewDeckIDsMatch(t *testing.T) {
	deck := NewDeck()
	ids := NewDeckIDs()
	if len(ids) != len(deck) {
		t.Fatalf("NewDeckIDs len %d != deck %d", len(ids), len(deck))
	}
	for i, id := range ids {
		if id != deck[i].ID {
			t.Fatalf("ids[%d] %v != deck[%d].ID %v", i, id, i, deck[i].ID)
		}
	}
}

func TestNewDeckDeterministic(t *testing.T) {
	a := NewDeck()
	b := NewDeck()
	if len(a) != len(b) {
		t.Fatalf("len mismatch")
	}
	for i := range a {
		if a[i].ID != b[i].ID || a[i].Colour != b[i].Colour || a[i].Rank != b[i].Rank || a[i].IsJoker != b[i].IsJoker {
			t.Fatalf("deck not deterministic at %d: %+v vs %+v", i, a[i], b[i])
		}
	}
	// Check exact order snapshot: first few and last two are jokers
	if a[0].ID != "red-01-1" || a[1].ID != "red-01-2" {
		t.Fatalf("first IDs %v %v want red-01-1 red-01-2", a[0].ID, a[1].ID)
	}
	if !a[104].IsJoker || a[104].ID != "joker-1" || !a[105].IsJoker || a[105].ID != "joker-2" {
		t.Fatalf("last two should be jokers joker-1/2, got %v %v", a[104], a[105])
	}
}

func TestNewDeckTilesValidate(t *testing.T) {
	deck := NewDeck()
	for _, tl := range deck {
		if err := tl.Validate(); err != nil {
			t.Fatalf("tile %v Validate failed: %v", tl, err)
		}
	}
}
