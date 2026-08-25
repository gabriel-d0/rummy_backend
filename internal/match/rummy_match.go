// Package match — Rummy authoritative match handler skeleton (Day 22-23).
// Day 23: lobby waiting room with seat allocation, join/leave handling.
package match

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

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
	for _, msg := range messages {
		data := msg.GetData()
		senderId := PlayerId(msg.GetUserId())
		// Envelope parse + payload schema validation (Day 27-28). Use ValidateEnvelope which checks version/op/payload.
		env, err := protocol.ValidateEnvelope(data)
		if err != nil {
			// Try to extract requestId/op for error correlation even on parse failure
			var reqId string
			var op int64
			// Attempt to parse partial envelope for requestId/op
			// If data is not JSON, requestId stays empty
			_ = op
			_ = reqId
			if env.RequestId != "" {
				reqId = env.RequestId
				op = env.OpCode
			} else if msg.GetOpCode() != 0 {
				op = msg.GetOpCode()
			}
			code := protocol.ErrCodeBadJSON
			if pe, ok := err.(*protocol.ParseError); ok {
				switch pe.Code {
				case "empty":
					code = protocol.ErrCodeBadJSON
				case "bad_json":
					code = protocol.ErrCodeBadJSON
				case "bad_version":
					code = protocol.ErrCodeBadVersion
				case "unknown_opcode":
					code = protocol.ErrCodeUnknownOpcode
				case "bad_payload":
					code = protocol.ErrCodeBadPayload
				default:
					code = pe.Code
				}
			}
			sendError(dispatcher, msg, code, err.Error(), reqId, op, logger)
			continue
		}
		// Use envelope RequestId and OpCode for correlation; Nakama GetOpCode should match envelope OpCode
		op := env.OpCode
		requestId := env.RequestId
		// Active-player check (Day 33) — only for non-Waiting phases; start is special case below
		if st.GamePhase == PhaseOpeningDiscard || st.GamePhase == PhasePlaying {
			if perr := ValidateActivePlayer(st, senderId); perr != nil {
				perr.RequestId = requestId
				perr.OpCode = op
				sendError(dispatcher, msg, perr.Code, perr.Message, requestId, op, logger)
				continue
			}
		}
		// Phase-op check (Day 32/34)
		if perr := ValidatePhaseOp(st, op); perr != nil {
			perr.RequestId = requestId
			perr.OpCode = op
			sendError(dispatcher, msg, perr.Code, perr.Message, requestId, op, logger)
			continue
		}
		// Payload already validated via ValidateEnvelope, but double-check for safety
		if perr := protocol.ValidatePayload(op, env.Payload); perr != nil {
			if pe, ok := perr.(*protocol.ParseError); ok {
				sendError(dispatcher, msg, protocol.ErrCodeBadPayload, pe.Message, requestId, op, logger)
			} else {
				sendError(dispatcher, msg, protocol.ErrCodeBadPayload, perr.Error(), requestId, op, logger)
			}
			continue
		}

		logger.Debug("MatchLoop tick=%d op=%d sender=%s requestId=%s", tick, op, senderId, requestId)

		// Handle start opcode from host Seat 0 when in Waiting with 2-4 players.
		if op == protocol.OpClientStart {
			if st.GamePhase != PhaseWaiting {
				sendError(dispatcher, msg, protocol.ErrCodeWrongPhase, "game already started", requestId, op, logger)
				continue
			}
			if len(st.Players) < 2 {
				sendError(dispatcher, msg, protocol.ErrCodeBadRequest, "need 2 players", requestId, op, logger)
				continue
			}
			seat := SeatOfPlayer(st.Players, senderId)
			if seat != 0 {
				sendError(dispatcher, msg, protocol.ErrCodeNotYourTurn, "only host may start", requestId, op, logger)
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
			continue
		}

		// Opening discard: must be from Seat 0 with 15 tiles, becomes blocked tile.
		if op == protocol.OpClientDiscard && st.GamePhase == PhaseOpeningDiscard {
			if err := handleOpeningDiscard(st, senderId, env.Payload, requestId, op, dispatcher, msg, logger); err != nil {
				// handleOpeningDiscard already sent error
				continue
			}
			continue
		}

		// Other opcodes (draw/meld) will be handled Day 35+; for now just phase validation.
		// If we reach here, op was allowed by phase but not yet implemented — treat as not implemented
		sendError(dispatcher, msg, protocol.ErrCodeBadRequest, fmt.Sprintf("op %d not implemented yet", op), requestId, op, logger)
	}
	return st
}

// handleOpeningDiscard validates and applies the opening discard. It sends OpServerError on failure and returns error.
func handleOpeningDiscard(st *RoundState, senderId PlayerId, payload []byte, requestId string, op int64, dispatcher runtime.MatchDispatcher, sender runtime.Presence, logger runtime.Logger) error {
	seat := SeatOfPlayer(st.Players, senderId)
	if seat != 0 {
		sendError(dispatcher, sender, protocol.ErrCodeNotYourTurn, "only opening player may discard", requestId, op, logger)
		return fmt.Errorf("not opening player")
	}
	if len(st.Players) < 2 {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, "need 2 players", requestId, op, logger)
		return fmt.Errorf("not enough players")
	}
	if len(st.DiscardRow) != 0 {
		sendError(dispatcher, sender, protocol.ErrCodeWrongPhase, "opening discard already done", requestId, op, logger)
		return fmt.Errorf("already discarded")
	}
	// Parse payload {tileId}
	var p struct {
		TileId string `json:"tileId"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "payload must be {tileId}", requestId, op, logger)
		return err
	}
	tileId := tile.TileInstanceId(p.TileId)
	if tileId == "" {
		sendError(dispatcher, sender, protocol.ErrCodeBadPayload, "tileId required", requestId, op, logger)
		return fmt.Errorf("tileId empty")
	}
	// Find tile in sender's rack (must have 15)
	rack := st.Racks[seat]
	if len(rack) != 15 {
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, fmt.Sprintf("opening rack must have 15, has %d", len(rack)), requestId, op, logger)
		return fmt.Errorf("rack size")
	}
	idx := -1
	var tileToDiscard tile.TileInstance
	for i, t := range rack {
		if t.ID == tileId {
			idx = i
			tileToDiscard = t
			break
		}
	}
	if idx == -1 {
		sendError(dispatcher, sender, protocol.ErrCodeInvalidTile, "tile not in rack", requestId, op, logger)
		return fmt.Errorf("tile not in rack")
	}
	// Remove from rack (preserve order)
	newRack := make([]tile.TileInstance, 0, 14)
	newRack = append(newRack, rack[:idx]...)
	newRack = append(newRack, rack[idx+1:]...)
	st.Racks[seat] = newRack
	// Append to discard row as opening discard
	entry := DiscardEntry{Tile: tileToDiscard, IsOpeningDiscard: true, Index: 0}
	st.DiscardRow = append(st.DiscardRow, entry)
	// Advance to next seat and go to Playing MustDraw
	st.GamePhase = PhasePlaying
	// CurrentSeat is still opening seat (0); AdvanceTurn moves to next anticlockwise
	if err := AdvanceTurn(st); err != nil {
		// Should not happen — opening seat was validated
		sendError(dispatcher, sender, protocol.ErrCodeBadRequest, err.Error(), requestId, op, logger)
		return err
	}
	next := st.CurrentSeat
	logger.Info("Opening discard %v by %s, now current %v phase Playing", tileId, senderId, next)
	if dispatcher != nil {
		_ = dispatcher.BroadcastMessage(protocol.OpServerEvent, []byte(fmt.Sprintf(`{"phase":"Playing","currentSeat":%d,"discard":"%s","isOpening":true}`, int(next), tileId)), nil, nil, true)
	}
	return nil
}

// sendError sends OpServerError to the sender only (not broadcast) with requestId correlation.
func sendError(dispatcher runtime.MatchDispatcher, sender runtime.Presence, code, message, requestId string, op int64, logger runtime.Logger) {
	if dispatcher == nil {
		logger.Warn("sendError %s %s requestId=%s op=%d (no dispatcher)", code, message, requestId, op)
		return
	}
	// Prefer to send only to sender if available via Presence, else broadcast
	errResp := protocol.NewError(code, message, requestId, map[string]string{"op": fmt.Sprintf("%d", op)})
	errResp.OpCode = op
	payload := protocol.EncodeError(errResp)
	// Use presence from MatchData sender if possible
	var presences []runtime.Presence
	if sender != nil {
		presences = []runtime.Presence{sender}
	}
	_ = dispatcher.BroadcastMessage(protocol.OpServerError, payload, presences, nil, true)
	logger.Info("Sent error %s to %v: %s requestId=%s op=%d", code, sender, message, requestId, op)
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
