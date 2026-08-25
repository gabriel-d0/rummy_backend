// Package match — Visibility (Day 30). Centralized public/private view projection.
// Ensures hidden information (other players' racks) never leaks per AGENTS.md:173.
package match

import "github.com/gabriel-d0/rummy_backend/internal/rules/tile"

// PublicSnapshot is the view broadcast to all participants.
// It contains counts for racks, not tile IDs, except for discard row and table melds.
type PublicSnapshot struct {
	Version     int            `json:"v"`
	GamePhase   string         `json:"gamePhase"`
	TurnPhase   string         `json:"turnPhase"`
	CurrentSeat Seat           `json:"currentSeat"`
	Players     []PublicPlayer `json:"players"`
	StockCount  int            `json:"stockCount"`
	DiscardRow  []DiscardEntry `json:"discardRow"`
	TableMelds  []TableMeld    `json:"tableMelds"`
	Winner      Seat           `json:"winner"`
}

// PublicPlayer is the public view of a player.
type PublicPlayer struct {
	ID        PlayerId `json:"id"`
	Seat      Seat     `json:"seat"`
	HasOpened bool     `json:"hasOpened"`
	RackCount int      `json:"rackCount"`
}

// PrivateSnapshot is the per-player view: public + own rack.
type PrivateSnapshot struct {
	PublicSnapshot
	OwnRack []tile.TileInstance `json:"ownRack"`
	OwnSeat Seat                `json:"ownSeat"`
}

// PublicView returns the public snapshot for any observer.
func PublicView(state *RoundState) PublicSnapshot {
	players := make([]PublicPlayer, len(state.Players))
	for i, p := range state.Players {
		players[i] = PublicPlayer{
			ID:        p.ID,
			Seat:      p.Seat,
			HasOpened: p.HasOpened,
			RackCount: len(state.Racks[p.Seat]),
		}
	}
	return PublicSnapshot{
		Version:     1,
		GamePhase:   state.GamePhase.String(),
		TurnPhase:   state.TurnPhase.String(),
		CurrentSeat: state.CurrentSeat,
		Players:     players,
		StockCount:  len(state.Stock),
		DiscardRow:  state.DiscardRow,
		TableMelds:  state.TableMelds,
		Winner:      state.Winner,
	}
}

// PrivateView returns the private view for a specific seat.
// It includes OwnRack for that seat; other racks remain counts only.
func PrivateView(state *RoundState, seat Seat) PrivateSnapshot {
	pub := PublicView(state)
	ownRack := []tile.TileInstance{}
	if rack, ok := state.Racks[seat]; ok {
		ownRack = append([]tile.TileInstance(nil), rack...)
	}
	return PrivateSnapshot{
		PublicSnapshot: pub,
		OwnRack:        ownRack,
		OwnSeat:        seat,
	}
}
