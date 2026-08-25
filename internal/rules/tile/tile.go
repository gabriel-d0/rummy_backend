// Package tile defines the core Tile domain model for Romanian Tile Rummy.
// It is the Day 10 artefact of Phase 2 — pure types with no I/O or Nakama deps.
// Terminology per docs/terminology.md and docs/rules-decisions.md:1.1.
// Every TileInstance has a unique immutable ID (TileInstanceId) even when face-identical.
package tile

import "fmt"

// Colour is the tile colour. Joker has no colour until represented.
type Colour int

const (
	ColourInvalid Colour = iota
	Red
	Yellow
	Blue
	Black
)

var colourNames = map[Colour]string{
	Red:    "red",
	Yellow: "yellow",
	Blue:   "blue",
	Black:  "black",
}

var colourByName = map[string]Colour{
	"red":    Red,
	"yellow": Yellow,
	"blue":   Blue,
	"black":  Black,
}

// IsValid reports whether c is a real colour (red/yellow/blue/black).
func (c Colour) IsValid() bool {
	return c == Red || c == Yellow || c == Blue || c == Black
}

func (c Colour) String() string {
	if s, ok := colourNames[c]; ok {
		return s
	}
	return "invalid"
}

// ParseColour parses a colour name (lowercase) to Colour.
func ParseColour(s string) (Colour, bool) {
	c, ok := colourByName[s]
	return c, ok
}

// Rank is 1..13 where 1 = Ace, 11=J, 12=Q, 13=K.
// Ace dual low/high handling is in the run validator, not here.
type Rank int

const (
	RankAce Rank = 1
	RankMin Rank = 1
	RankMax Rank = 13
)

// IsValid reports whether r is 1..13.
func (r Rank) IsValid() bool {
	return r >= RankMin && r <= RankMax
}

func (r Rank) String() string {
	switch r {
	case 1:
		return "A"
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

// TileInstanceId is a unique immutable identifier for a tile copy.
// It is opaque — consumers must not parse it for colour/rank.
type TileInstanceId string

// IsValid reports whether id is non-empty.
func (id TileInstanceId) IsValid() bool {
	return id != ""
}

func (id TileInstanceId) String() string {
	return string(id)
}

// TileInstance is a concrete tile at a location (shuffle, rack, meld, discard).
// For numbered tiles: IsJoker==false, Colour valid, Rank valid.
// For jokers: IsJoker==true, Colour==ColourInvalid, Rank==0, ID unique.
type TileInstance struct {
	ID      TileInstanceId
	Colour  Colour
	Rank    Rank
	IsJoker bool
}

// IsNumbered reports whether this is a numbered 1–13 tile.
func (t TileInstance) IsNumbered() bool {
	return !t.IsJoker
}

// Validate checks internal invariants (colour/rank vs joker flag, ID, range).
func (t TileInstance) Validate() error {
	if !t.ID.IsValid() {
		return fmt.Errorf("tile ID must be non-empty")
	}
	if t.IsJoker {
		if t.Colour != ColourInvalid {
			return fmt.Errorf("joker %q must have invalid colour, got %v", t.ID, t.Colour)
		}
		if t.Rank != 0 {
			return fmt.Errorf("joker %q must have rank 0, got %v", t.ID, t.Rank)
		}
		return nil
	}
	// numbered
	if !t.Colour.IsValid() {
		return fmt.Errorf("tile %q invalid colour %v", t.ID, t.Colour)
	}
	if !t.Rank.IsValid() {
		return fmt.Errorf("tile %q invalid rank %v", t.ID, t.Rank)
	}
	return nil
}

func (t TileInstance) String() string {
	if t.IsJoker {
		return fmt.Sprintf("Joker{%s}", t.ID)
	}
	return fmt.Sprintf("%s-%s{%s}", t.Colour, t.Rank, t.ID)
}

// NewTile creates a numbered tile and validates it.
func NewTile(id TileInstanceId, colour Colour, rank Rank) (TileInstance, error) {
	t := TileInstance{ID: id, Colour: colour, Rank: rank, IsJoker: false}
	if err := t.Validate(); err != nil {
		return TileInstance{}, err
	}
	return t, nil
}

// NewJoker creates a joker tile with the given unique ID.
func NewJoker(id TileInstanceId) (TileInstance, error) {
	t := TileInstance{ID: id, Colour: ColourInvalid, Rank: 0, IsJoker: true}
	if err := t.Validate(); err != nil {
		return TileInstance{}, err
	}
	return t, nil
}

// MustTile is a helper for tests — panics on error.
func MustTile(id TileInstanceId, colour Colour, rank Rank) TileInstance {
	t, err := NewTile(id, colour, rank)
	if err != nil {
		panic(err)
	}
	return t
}

// MustJoker is a helper for tests.
func MustJoker(id TileInstanceId) TileInstance {
	t, err := NewJoker(id)
	if err != nil {
		panic(err)
	}
	return t
}
