// Package match — Rummy authoritative match handler skeleton (Day 22-23).
// Day 23: lobby waiting room with seat allocation, join/leave handling.
package match

import (
	"context"
	"database/sql"

	"github.com/gabriel-d0/rummy_backend/internal/protocol"
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
	st, ok := state.(*RoundState)
	if !ok {
		logger.Error("MatchJoinAttempt bad state type %T", state)
		return state, false, "bad state"
	}
	// Only allow joins in Waiting phase
	if st.GamePhase != PhaseWaiting {
		return state, false, "game already started"
	}
	if len(st.Players) >= 4 {
		return state, false, "match full"
	}
	// Reject duplicate user
	for _, p := range st.Players {
		if string(p.ID) == presence.GetUserId() {
			return state, false, "already joined"
		}
	}
	return state, true, ""
}

func (m *RummyMatch) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	logger.Info("MatchJoin presences=%d", len(presences))
	st, ok := state.(*RoundState)
	if !ok {
		logger.Error("MatchJoin bad state type %T", state)
		return state
	}
	// Allocate seats deterministically in join order for each presence
	for _, pres := range presences {
		pid := PlayerId(pres.GetUserId())
		// Double-check not already present (defensive)
		found := false
		for _, p := range st.Players {
			if p.ID == pid {
				found = true
				break
			}
		}
		if found {
			continue
		}
		seat := Seat(len(st.Players))
		st.Players = append(st.Players, PlayerState{ID: pid, Seat: seat, HasOpened: false})
		if st.Racks == nil {
			st.Racks = map[Seat][]tile.TileInstance{}
		}
		st.Racks[seat] = []tile.TileInstance{}
		logger.Info("Player %s joined as %v (total %d)", pid, seat, len(st.Players))
	}
	// Update label with player count for match listing (e.g. "rummy:2")
	if nk != nil && dispatcher != nil {
		_ = dispatcher.MatchLabelUpdate(st.GamePhase.String() + ":" + string(rune('0'+len(st.Players))))
	}
	return st
}

func (m *RummyMatch) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	logger.Info("MatchLeave presences=%d", len(presences))
	st, ok := state.(*RoundState)
	if !ok {
		logger.Error("MatchLeave bad state type %T", state)
		return state
	}
	// Remove leaving players; keep seat numbers stable for MVP (no re-shuffle)
	// This keeps Racks entries but removes Player; Day 24 will handle dealer rotation.
	for _, pres := range presences {
		pid := PlayerId(pres.GetUserId())
		newPlayers := make([]PlayerState, 0, len(st.Players))
		for _, p := range st.Players {
			if p.ID != pid {
				newPlayers = append(newPlayers, p)
			} else {
				logger.Info("Player %s (seat %v) left", pid, p.Seat)
				delete(st.Racks, p.Seat)
			}
		}
		st.Players = newPlayers
	}
	// If no players left, Nakama will terminate via MatchTerminate after grace
	return st
}

func (m *RummyMatch) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	st, ok := state.(*RoundState)
	if !ok {
		logger.Error("MatchLoop bad state type %T", state)
		return state
	}
	// Handle start opcode from host Seat 0 when in Waiting with 2-4 players.
	// Day 24-26: OpClientStart is stable 1 per protocol.Version 1.
	for _, msg := range messages {
		op := msg.GetOpCode()
		senderId := msg.GetUserId()
		logger.Debug("MatchLoop tick=%d op=%d sender=%s len=%d", tick, op, senderId, len(messages))
		if op == protocol.OpClientStart {
			if st.GamePhase != PhaseWaiting {
				logger.Warn("Start rejected: game already started phase %v", st.GamePhase)
				continue
			}
			if len(st.Players) < 2 {
				logger.Warn("Start rejected: need 2 players, have %d", len(st.Players))
				continue
			}
			seat := SeatOfPlayer(st.Players, PlayerId(senderId))
			if seat != 0 {
				logger.Warn("Start rejected: only host Seat 0 may start, sender seat %v", seat)
				continue
			}
			// Transition to OpeningDiscard
			st.GamePhase = PhaseOpeningDiscard
			st.CurrentSeat = 0
			st.TurnPhase = TurnMustDraw
			logger.Info("Match started by host %s with %d players, phase OpeningDiscard seat 0", senderId, len(st.Players))
			if dispatcher != nil {
				_ = dispatcher.BroadcastMessage(protocol.OpServerEvent, []byte(`{"phase":"OpeningDiscard","currentSeat":0}`), nil, nil, true)
			}
		}
	}
	return st
}

func (m *RummyMatch) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Info("MatchTerminate grace=%d", graceSeconds)
	return state
}

func (m *RummyMatch) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	logger.Info("MatchSignal data=%s", data)
	st, ok := state.(*RoundState)
	if !ok {
		return state, "bad state"
	}
	// Also allow signal-based start for local dev (e.g. nk.matchSignal or RPC)
	// Expected data: "start:<hostUserId>" or "start"
	if data == "start" || len(data) > 6 && data[:6] == "start:" {
		if st.GamePhase != PhaseWaiting {
			return state, "already started"
		}
		if len(st.Players) < 2 {
			return state, "need 2 players"
		}
		// If host specified, verify it is Seat 0
		if len(data) > 6 {
			hostId := PlayerId(data[6:])
			if SeatOfPlayer(st.Players, hostId) != 0 {
				return state, "only host may start"
			}
		}
		st.GamePhase = PhaseOpeningDiscard
		st.CurrentSeat = 0
		st.TurnPhase = TurnMustDraw
		logger.Info("Match started via signal with %d players", len(st.Players))
		return st, "started"
	}
	return st, ""
}
