// Package setup — Shuffle (Day 18). Fisher–Yates via injectable Rand.
package setup

import "github.com/gabriel-d0/rummy_backend/internal/rules/tile"

// Shuffle returns a shuffled copy of deck using r. It is Fisher–Yates via
// Rand.Shuffle (which itself is Fisher–Yates) and is deterministic for a
// given seed. Original deck is not modified; the returned slice has the same
// 106 TileInstances permuted with no loss or duplication.
func Shuffle(deck []tile.TileInstance, r *Rand) []tile.TileInstance {
	if r == nil {
		panic("Shuffle requires non-nil Rand")
	}
	shuffled := make([]tile.TileInstance, len(deck))
	copy(shuffled, deck)
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}
