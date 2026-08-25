// Package setup — Deal (Day 19). MVP dealing: opening player 15, others 14, remainder stock.
package setup

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/match"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// Deal distributes shuffled deck to n players (2..4) per docs/rules-decisions.md:2.
// Opening player (Seat 0) gets 15 tiles, all others get 14, remainder is stock.
// Deck must have 106 tiles; it is not shuffled here — caller shuffles via Shuffle first if desired.
// Racks are map[Seat][]TileInstance indexed by Seat 0..n-1 deterministically.
func Deal(deck []tile.TileInstance, n int) (map[match.Seat][]tile.TileInstance, []tile.TileInstance, error) {
	if len(deck) != 106 {
		return nil, nil, fmt.Errorf("deck must have 106 tiles, got %d", len(deck))
	}
	if n < 2 || n > 4 {
		return nil, nil, fmt.Errorf("player count %d must be 2..4", n)
	}
	needed := 15 + (n-1)*14 // 29/43/57
	if needed > 106 {
		return nil, nil, fmt.Errorf("not enough tiles: need %d have 106", needed)
	}
	racks := make(map[match.Seat][]tile.TileInstance, n)
	offset := 0
	for seat := 0; seat < n; seat++ {
		count := 14
		if seat == 0 {
			count = 15
		}
		rack := make([]tile.TileInstance, count)
		copy(rack, deck[offset:offset+count])
		racks[match.Seat(seat)] = rack
		offset += count
	}
	stock := make([]tile.TileInstance, 106-offset)
	copy(stock, deck[offset:])
	return racks, stock, nil
}

// DealForPlayers is a convenience that uses playerIds length for n and returns
// racks keyed by Seat plus stock. Seats are assigned via AssignSeats elsewhere;
// this helper just deals.
func DealForPlayers(deck []tile.TileInstance, playerIds []match.PlayerId) (map[match.Seat][]tile.TileInstance, []tile.TileInstance, error) {
	return Deal(deck, len(playerIds))
}
