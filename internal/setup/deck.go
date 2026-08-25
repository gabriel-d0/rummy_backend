// Package setup creates the Romanian Tile Rummy deck.
// Day 15 — Full deck factory: 104 numbered (4 colours ×13 ranks ×2 copies) + 2 jokers =106
// with unique TileInstanceIds per docs/rules-decisions.md:1.1.
package setup

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// NewDeck returns a fresh 106-tile deck in deterministic order:
// Red 1-13 ×2, Yellow 1-13 ×2, Blue 1-13 ×2, Black 1-13 ×2, then 2 jokers.
// IDs are unique and immutable, e.g. "red-01-1", "red-01-2", ..., "joker-1".
// Order is deterministic; shuffle is separate (Day 17-18).
func NewDeck() []tile.TileInstance {
	colours := []tile.Colour{tile.Red, tile.Yellow, tile.Blue, tile.Black}
	deck := make([]tile.TileInstance, 0, 106)
	for _, c := range colours {
		for r := tile.RankMin; r <= tile.RankMax; r++ {
			for copyIdx := 1; copyIdx <= 2; copyIdx++ {
				id := tile.TileInstanceId(fmt.Sprintf("%s-%02d-%d", c.String(), int(r), copyIdx))
				t := tile.MustTile(id, c, r)
				deck = append(deck, t)
			}
		}
	}
	// 2 jokers
	for i := 1; i <= 2; i++ {
		id := tile.TileInstanceId(fmt.Sprintf("joker-%d", i))
		j := tile.MustJoker(id)
		deck = append(deck, j)
	}
	return deck
}

// NewDeckIDs returns the TileInstanceIds of a fresh deck — convenient for tests.
func NewDeckIDs() []tile.TileInstanceId {
	deck := NewDeck()
	ids := make([]tile.TileInstanceId, len(deck))
	for i, t := range deck {
		ids[i] = t.ID
	}
	return ids
}
