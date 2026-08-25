package tile

import (
	"testing"
)

func TestColourStringAndParse(t *testing.T) {
	cases := []struct {
		c Colour
		s string
	}{
		{Red, "red"},
		{Yellow, "yellow"},
		{Blue, "blue"},
		{Black, "black"},
	}
	for _, tc := range cases {
		if got := tc.c.String(); got != tc.s {
			t.Fatalf("Colour %d String()=%q want %q", tc.c, got, tc.s)
		}
		gotC, ok := ParseColour(tc.s)
		if !ok || gotC != tc.c {
			t.Fatalf("ParseColour(%q)=%v %v want %v true", tc.s, gotC, ok, tc.c)
		}
		if !tc.c.IsValid() {
			t.Fatalf("Colour %v should be valid", tc.c)
		}
	}
	if ColourInvalid.IsValid() {
		t.Fatalf("ColourInvalid should be invalid")
	}
	if Colour(99).IsValid() {
		t.Fatalf("out of range colour should be invalid")
	}
}

func TestRankValidation(t *testing.T) {
	for r := RankMin; r <= RankMax; r++ {
		if !r.IsValid() {
			t.Fatalf("Rank %v should be valid", r)
		}
	}
	if Rank(0).IsValid() || Rank(14).IsValid() {
		t.Fatalf("out of range ranks should be invalid")
	}
	if RankAce.String() != "A" {
		t.Fatalf("Ace String = %q want A", RankAce.String())
	}
}

func TestNewTileValid(t *testing.T) {
	tile, err := NewTile("t-001", Red, Rank(5))
	if err != nil {
		t.Fatalf("NewTile failed: %v", err)
	}
	if tile.IsJoker || !tile.IsNumbered() {
		t.Fatalf("expected numbered tile")
	}
	if err := tile.Validate(); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
}

func TestNewTileInvalidRank(t *testing.T) {
	if _, err := NewTile("t-002", Red, Rank(0)); err == nil {
		t.Fatalf("expected error for rank 0")
	}
	if _, err := NewTile("t-003", ColourInvalid, Rank(5)); err == nil {
		t.Fatalf("expected error for invalid colour")
	}
	if _, err := NewTile("", Red, Rank(5)); err == nil {
		t.Fatalf("expected error for empty ID")
	}
}

func TestNewJoker(t *testing.T) {
	j, err := NewJoker("j-001")
	if err != nil {
		t.Fatalf("NewJoker failed: %v", err)
	}
	if !j.IsJoker || j.IsNumbered() {
		t.Fatalf("expected joker")
	}
	if j.Colour != ColourInvalid || j.Rank != 0 {
		t.Fatalf("joker should have invalid colour and rank 0")
	}
	if err := j.Validate(); err != nil {
		t.Fatalf("joker Validate failed: %v", err)
	}
	// joker with colour should fail
	bad := TileInstance{ID: "j-002", Colour: Red, Rank: 0, IsJoker: true}
	if err := bad.Validate(); err == nil {
		t.Fatalf("joker with colour should be invalid")
	}
}

func TestTileIDUniqueness(t *testing.T) {
	ids := map[TileInstanceId]bool{}
	tiles := []TileInstance{}
	for i := 0; i < 10; i++ {
		id := TileInstanceId(string(rune('a' + i)))
		// Use distinct IDs
		tile := MustTile(TileInstanceId("id-"+string(rune('0'+i))), Red, Rank(1))
		if ids[tile.ID] {
			t.Fatalf("duplicate ID %v", tile.ID)
		}
		ids[tile.ID] = true
		tiles = append(tiles, tile)
		_ = id
	}
	// Ensure 10 unique
	if len(ids) != 10 {
		t.Fatalf("expected 10 unique IDs, got %d", len(ids))
	}
}

func TestMustHelpersPanicOnInvalid(t *testing.T) {
	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		MustTile("", Red, Rank(1))
	}()
	if !didPanic {
		t.Fatalf("MustTile should panic on invalid ID")
	}
}

func TestTileString(t *testing.T) {
	tile := MustTile("t-123", Blue, RankAce)
	if tile.String() != "blue-A{t-123}" {
		t.Fatalf("String got %q", tile.String())
	}
	joker := MustJoker("j-999")
	if joker.String() != "Joker{j-999}" {
		t.Fatalf("Joker String got %q", joker.String())
	}
}
