// Package match — Rummy authoritative match handler skeleton (Day 22).
// Minimal registration; full lobby/match logic lands Day 23-31.
package match

import (
	"context"
	"database/sql"

	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
	"github.com/heroiclabs/nakama-common/runtime"
)

// RummyMatch is the authoritative match handler for Romanian Tile Rummy.
// It implements runtime.Match. State is kept as *RoundState (or nil until MatchInit).
type RummyMatch struct{}

// NewRummyMatch is the factory used by RegisterMatch. Nakama calls it once per match creation.
func NewRummyMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
	logger.Info("RummyMatch factory called")
	return &RummyMatch{}, nil
}

func (m *RummyMatch) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	logger.Info("Rummy MatchInit params=%v", params)
	// Tick rate 5 is comfortable for turn-based; label is filterable in match listing.
	// State is nil until MatchJoin creates players — or could init empty RoundState.
	// For skeleton we start with empty map and let joins populate.
	tickRate := 5
	label := "rummy"
	// Use a minimal placeholder state: a RoundState in Waiting with no players yet.
	// Real init will happen when first player joins or via RPC signal.
	initialState := &RoundState{
		Players:     []PlayerState{},
		Racks:       map[Seat][]tile.TileInstance{},
		Stock:       []tile.TileInstance{},
		DiscardRow:  []DiscardEntry{},
		TableMelds:  []TableMeld{},
		CurrentSeat: SeatInvalid,
		GamePhase:   PhaseWaiting,
		TurnPhase:   TurnMustDraw,
		Winner:      SeatInvalid,
	}
	// Need import for TileInstance; use fully qualified via tile package.
	// But TileInstance is from internal/rules/tile; we already have it in match via state.
	// The placeholder above uses TileInstance from tile via type alias? Actually RoundState uses tile.TileInstance,
	// but we used bare TileInstance — need alias. Instead use tile.TileInstance via import.
	// Fix: import tile and use tile.TileInstance — but RoundState already typed as tile.TileInstance.
	// So the literal above should be []tile.TileInstance{} which is same as TileInstance alias.
	// We already did.
	return initialState, tickRate, label
}

func (m *RummyMatch) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	logger.Info("MatchJoinAttempt user=%s session=%s metadata=%v", presence.GetUserId(), presence.GetSessionId(), metadata)
	// Allow up to 4 players; real logic will check state.Players length.
	if st, ok := state.(*RoundState); ok {
		if len(st.Players) >= 4 {
			return state, false, "match full"
		}
	}
	return state, true, ""
}

func (m *RummyMatch) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	logger.Info("MatchJoin presences=%d", len(presences))
	// For skeleton we just log; Day 24 will allocate seats via AssignSeats.
	return state
}

func (m *RummyMatch) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	logger.Info("MatchLeave presences=%d", len(presences))
	return state
}

func (m *RummyMatch) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	// Turn-based: no per-tick logic needed; messages are client opcodes handled via dispatcher.
	// Day 26+ will route opcodes here.
	if len(messages) > 0 {
		logger.Debug("MatchLoop tick=%d messages=%d", tick, len(messages))
	}
	return state
}

func (m *RummyMatch) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Info("MatchTerminate grace=%d", graceSeconds)
	return state
}

func (m *RummyMatch) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	logger.Info("MatchSignal data=%s", data)
	return state, ""
}
