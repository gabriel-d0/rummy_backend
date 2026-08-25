// Package setup — Round initialization (Day 20).
// Builds createRoundState with seats, opening player, racks, stock, empty
// table melds, empty discard row with opening-discard marker ready, and
// opening turn phase per docs/rules-decisions.md:2 and docs/terminology.md.
package setup

import (
	"fmt"

	"github.com/gabriel-d0/rummy_backend/internal/match"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

// NewRoundState creates a legal initial round state for 2–4 players.
// Deck is NewDeck() shuffled deterministically with seed; playerIds order
// determines seats (first joiner Seat 0 is opening player with 15 tiles).
// Returned allTiles is the shuffled deck for CheckTileConservation.
// Opening player must discard first; state is in PhaseOpeningDiscard with
// CurrentSeat 0, empty DiscardRow and TableMelds.
func NewRoundState(playerIds []match.PlayerId, seed int64) (*match.RoundState, []tile.TileInstance, error) {
	deck := NewDeck()
	shuffled := Shuffle(deck, NewSeededRand(seed))
	return NewRoundStateWithDeck(playerIds, shuffled)
}

// NewRoundStateWithDeck builds round state from an already-shuffled 106-tile deck.
// Useful for tests with fixed deck ordering.
func NewRoundStateWithDeck(playerIds []match.PlayerId, shuffled []tile.TileInstance) (*match.RoundState, []tile.TileInstance, error) {
	if len(playerIds) < 2 || len(playerIds) > 4 {
		return nil, nil, fmt.Errorf("player count %d must be 2..4", len(playerIds))
	}
	if len(shuffled) != 106 {
		return nil, nil, fmt.Errorf("shuffled deck must have 106 tiles, got %d", len(shuffled))
	}
	n := len(playerIds)
	players, err := match.AssignSeats(playerIds)
	if err != nil {
		return nil, nil, err
	}
	racks, stock, err := Deal(shuffled, n)
	if err != nil {
		return nil, nil, err
	}

	state := &match.RoundState{
		Players:     players,
		Racks:       racks,
		Stock:       stock,
		DiscardRow:  []match.DiscardEntry{},
		TableMelds:  []match.TableMeld{},
		CurrentSeat: 0, // opening player
		GamePhase:   match.PhaseOpeningDiscard,
		TurnPhase:   match.TurnMustDraw, // not used in OpeningDiscard but set for completeness
		Winner:      match.SeatInvalid,
	}
	if err := state.Validate(); err != nil {
		return nil, nil, fmt.Errorf("round state validate failed: %w", err)
	}
	if err := match.CheckTileConservation(state, shuffled); err != nil {
		return nil, nil, fmt.Errorf("conservation failed: %w", err)
	}
	return state, shuffled, nil
}
